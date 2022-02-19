package handlers

import (
	"fmt"

	"net/http"

	"github.com/eyecuelab/go-api/cmd/api/serializers"
	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/kit/db/psql"
	"github.com/eyecuelab/kit/web"
	"github.com/labstack/echo"
)

// Register register handler
func Register(c web.ApiContext) error {
	u := new(models.User)
	if err := c.BindAndValidate(u); err != nil {
		return err
	}
	if err := u.RegisterWithPassword(); err != nil {
		return err
	}
	// if err := notifications.Welcome(u); err != nil {
	// 	return err
	// }
	return authedUserResponse(&c, u)
}

// GetUser retrieve email associated with confirmation token
func GetUser(c web.ApiContext) error {
	user := new(serializers.User)
	db := psql.DB.First(user, "confirm_email_token = ? and (user_id is null or user_id = 0)", c.Param("token"))
	if db.Error != nil {
		return echo.ErrNotFound
	}

	return c.JsonApiOK(user)
}

// RegisterUser register user
func RegisterUser(c web.ApiContext) error {
	user := new(serializers.User)
	db := psql.DB.First(user, "confirm_email_token = ? and (user_id is null or user_id = 0)", c.Param("token"))
	if db.Error != nil {
		return echo.ErrNotFound
	}

	unscopedUser, dbErr := fetchUnscoped(user.Email)

	if dbErr != nil {
		return dbErr
	} else if unscopedUser != nil {
		if err := c.BindAndValidate(unscopedUser); err != nil {
			return err
		}

		tx := psql.DB.Begin()
		if err := tx.Commit().Error; err != nil {
			return err
		}

		if err := unscopedUser.UpdatePassword(); err != nil {
			return nil
		}

		return authedUserResponse(&c, unscopedUser)

	} else {

		u := new(models.User)
		u.Email = user.Email
		if err := c.BindAndValidate(u); err != nil {
			return err
		}
		if err := u.RegisterWithPassword(); err != nil {
			return err
		}

		return authedUserResponse(&c, u)
	}

	// if err := notifications.Welcome(u); err != nil {
	// 	return err
	// }
}

func authedUserResponse(c *web.ApiContext, authedUser *models.User) error {
	if authedUser == nil {
		return (*c).ApiError("Invalid", http.StatusUnauthorized)
	} else if token, err := authedUser.JwtToken(); err != nil {
		return err
	} else {
		(*c).Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
		return (*c).JsonApiOK(authedUser)
	}
}

func fetchUnscoped(email string) (*models.User, error) {
	user := new(models.User)
	db := psql.DB.Unscoped().Where("email =  ?", email).Find(&user)
	if db.RecordNotFound() {
		return nil, nil
	} else if db.Error != nil {
		return nil, db.Error
	}

	return user, nil
}
