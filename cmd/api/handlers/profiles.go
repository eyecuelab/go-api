package handlers

import (
	"github.com/eyecuelab/kit/brake"

	"github.com/eyecuelab/go-api/cmd/middleware"
	"github.com/eyecuelab/kit/db/psql"
)

// GetUserProfile ...
func GetUserProfile(c middleware.AuthedContext) error {
	user := c.User()

	return c.JsonApiOK(user)
}

// UpdateUserProfile ...
func UpdateUserProfile(c middleware.AuthedContext) error {
	user := c.User()

	if err := c.BindAndValidate(user); err != nil {
		return err
	}

	permittedAttrs := c.Attrs("first_name", "last_name")
	if err := psql.DB.Model(&user).Updates(permittedAttrs).Error; err != nil {
		brake.Notify(err, nil, brake.SeverityCritical)
		return err
	}

	return c.JsonApiOK(user)
}
