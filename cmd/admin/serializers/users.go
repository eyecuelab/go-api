package serializers

import (
	"fmt"
	"time"

	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/kit/web/meta"
	"github.com/eyecuelab/kit/web/pagination"
	"github.com/google/jsonapi"
	"github.com/lib/pq"
)

// User user data structure for admin
type User struct {
	ID           int             `jsonapi:"primary,user" gorm:"primary_key"`
	FirstName    string          `jsonapi:"attr,first_name" json:"first_name"`
	LastName     string          `jsonapi:"attr,last_name" json:"last_name"`
	Email        string          `jsonapi:"attr,email" valid:"email"`
	LastSigninAt time.Time       `jsonapi:"attr,last_signin_at,iso8601"`
	CreatedAt    time.Time       `jsonapi:"attr,created_at,iso8601"`
	UpdatedAt    time.Time       `jsonapi:"attr,updated_at,iso8601"`
	Source       string          `jsonapi:"attr,source,omitempty" gorm:"-" sql:"-"`
	Token        string          `jsonapi:"attr,token,omitempty" gorm:"-" sql:"-"`
	URL          string          `gorm:"-" sql:"-" jsonapi:"attr,url,omitempty"`
	Roles        pq.StringArray  `jsonapi:"attr,roles"`
	CompanyID    int             `jsonapi:"attr,company_id,omitempty" json:"company_id"`
	Company      *models.Company `jsonapi:"relation,company,omitempty" gorm:"save_associations:false"`
	models.User
}

// Users admin users list serializer
type Users []*User

func (u User) path() string {
	return fmt.Sprintf("users/%d", u.ID)
}

// JSONAPIMeta meta data
func (u User) JSONAPIMeta() *jsonapi.Meta {
	var ah meta.ActionHolder
	roles := u.Roles
	if roles == nil {
		roles = []string{}
	}

	ah.AddAction(meta.PATCH, "update", u.path()).
		Field("first_name", meta.InputText, u.FirstName, true).
		Field("last_name", meta.InputText, u.LastName, true).
		Field("email", meta.InputText, u.Email, true).
		FieldWithOpts("roles_payload", meta.InputSelect, roles, false,
			[]meta.FieldOption{
				meta.FieldOption{Label: "Admin", Value: "admin"},
			},
		).
		Field("password", meta.InputPass, nil, false)

	ah.AddAction(meta.DELETE, "delete", u.path())

	return ah.RenderActions()
}

// JSONAPILinks jsonapi links
func (i Users) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": meta.APIURL("users"),
	}
}

// JSONAPIMeta jsonapi meta
func (i Users) JSONAPIMeta() *jsonapi.Meta {
	var ah meta.ActionHolder
	ah.AddAction(meta.GET, "index", "users").
		Field("search", meta.InputText, QueryParams.Get("search"), false).
		Pagination(&pagination.Data)

	ah.AddAction(meta.POST, "create", "users").
		Field("first_name", meta.InputText, nil, true).
		Field("last_name", meta.InputText, nil, true).
		Field("email", meta.InputText, nil, true).
		FieldWithOpts("roles_payload", meta.InputSelect, []string{}, false,
			[]meta.FieldOption{
				meta.FieldOption{Label: "Admin", Value: "admin"},
			},
		)

	return ah.RenderActions()
}
