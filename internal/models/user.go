package models

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/eyecuelab/kit/db/psql"
	"github.com/eyecuelab/kit/functools"
	"github.com/eyecuelab/kit/random"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	insertPasswordSQL        = "insert into credentials (user_id, value) values (?, crypt(?, gen_salt('bf', 8)));"
	insertPasswordResetToken = "update users set reset_token = ?, reset_sent_at = now() where id = ?;"

	loginJoinSQL            = "join credentials on credentials.user_id = users.id"
	loginWhereSQL           = "users.email = ? and credentials.source = 'password' and credentials.value = crypt(?, credentials.value)"
	loginThridPartyWhereSQL = "credentials.source = ? and credentials.value = ?"

	userByResetTokenQuery = `reset_token = ? and (reset_sent_at + (interval '1 hours' * ?)) > now()`

	resetPasswordSQL = `insert into credentials as c (user_id, value) values (?, crypt(?, gen_salt('bf', 8)))
                        on conflict (user_id, source) do update
	                      set value = crypt(?, gen_salt('bf', 8))
	                      where c.user_id = ?`
	updatePasswordSQL    = "update credentials set value = crypt(?, gen_salt('bf', 8)) where user_id = ?;"
	userCompaniesJoinSQL = "join companies on companies.id = users.company_id"
	setLastSigninAtSQL   = "update users set last_signin_at = now() where id = ?"
)

// User type
type User struct {
	ID           int       `jsonapi:"primary,user" gorm:"primary_key"`
	FirstName    string    `jsonapi:"attr,first_name" json:"first_name"`
	LastName     string    `jsonapi:"attr,last_name" json:"last_name"`
	Email        string    `jsonapi:"attr,email" valid:"email"`
	Password     string    `jsonapi:"attr,password,omitempty" gorm:"-" sql:"-" valid:"length(6|100)~Password must be at least 6 characters"`
	LastSigninAt time.Time `jsonapi:"attr,last_signin_at,iso8601"`
	CreatedAt    time.Time `jsonapi:"attr,created_at,iso8601"`
	UpdatedAt    time.Time `jsonapi:"attr,updated_at,iso8601"`
	ResetToken   string
	Source       string         `jsonapi:"attr,source,omitempty" gorm:"-" sql:"-"`
	Token        string         `jsonapi:"attr,token,omitempty" gorm:"-" sql:"-"`
	URL          string         `gorm:"-" sql:"-" jsonapi:"attr,url,omitempty"`
	Roles        pq.StringArray `jsonapi:"attr,roles"`
	RolesPayload []string       `jsonapi:"attr,roles_payload" gorm:"-" sql:"-"`
	CompanyID    int            `jsonapi:"attr,company_id,omitempty" json:"company_id"`
	Company      *Company       `jsonapi:"relation,company,omitempty" gorm:"save_associations:false"`
}

// FullName ...
func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// MergeVars ...
func (u *User) MergeVars() map[string]string {
	return map[string]string{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"full_name":  u.FullName(),
		"email":      u.Email,
	}
}

// IsAdmin ...
func (u *User) IsAdmin() bool {
	return u.HasRole("admin")
}

// HasRole ...
func (u *User) HasRole(r string) bool {
	return functools.StringSliceContains(u.Roles, string(r))
}

// ForgotMergeVars ...
func (u *User) ForgotMergeVars() map[string]string {
	vars := u.MergeVars()
	vars["token"] = u.ResetToken

	return vars
}

// RegisterWithPassword ...
func (u *User) RegisterWithPassword() error {
	tx := psql.DB.Begin()

	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Exec(insertPasswordSQL, u.ID, u.Password).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Login ...
func (u *User) Login() (*User, error) {
	authed := new(User)
	result := psql.DB.Joins(loginJoinSQL).Where(loginWhereSQL, u.Email, u.Password).
		Find(&authed)
	if result.RecordNotFound() {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	if err := psql.DB.Exec(setLastSigninAtSQL, authed.ID).Error; err != nil {
		return nil, err
	}

	return authed, nil
}

// SetResetPasswordToken ...
func (u *User) SetResetPasswordToken() error {
	token, err := random.RandomString(60)
	if err != nil {
		return err
	}

	if err := psql.DB.Exec(insertPasswordResetToken, token, u.ID).Error; err != nil {
		return err
	}
	u.ResetToken = token
	return nil
}

// UpdatePassword ...
func (u *User) UpdatePassword() error {
	if u.Password == "" {
		return nil
	}

	userCred := new(Credential)
	if err := psql.DB.Where(Credential{UserID: u.ID, Source: "password"}).FirstOrCreate(&userCred).Error; err != nil {
		return err
	}

	if err := psql.DB.Exec(updatePasswordSQL, u.Password, u.ID).Error; err != nil {
		return err
	}

	return nil
}

// JwtToken ...
func (u *User) JwtToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": u.ID})
	tokenString, err := token.SignedString([]byte(viper.GetString("secret")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Find find user by ID
func (u *User) Find(preloads ...string) error {
	db := psql.DB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db.First(u, u.ID).Error
}

// UpdatePasswordByResetToken ...
func (u *User) UpdatePasswordByResetToken() error {
	expires := viper.GetString("reset_expires")
	db := psql.DB.First(u, userByResetTokenQuery, u.Token, expires)
	if db.RecordNotFound() {
		return echo.ErrNotFound
	} else if db.Error != nil {
		return db.Error
	}

	result := psql.DB.Exec(resetPasswordSQL, u.ID, u.Password, u.Password, u.ID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func authedRecordCheck(primedDb *gorm.DB, authed *User) error {
	result := primedDb.Find(&authed)
	if result.RecordNotFound() {
		return nil
	} else if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
