package handlers

import (
	"net/http"

	"github.com/eyecuelab/go-api/cmd/admin/serializers"
	"github.com/eyecuelab/go-api/cmd/middleware"
	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/kit/db/psql"
	"github.com/eyecuelab/kit/web/pagination"
	"github.com/labstack/echo"
)

// ListUsers admin list of users
func ListUsers(c middleware.AuthedContext) error {
	var users serializers.Users
	serializers.QueryParams = c.QueryParams()
	role := c.QueryParam("role")

	scope := psql.DB.
		Preload("Company").
		Order("id DESC")

	if role == "admin" {
		scope = scope.Scopes(models.UsersWithRole("admin"))
	} else if role == "user" {
		scope = scope.Where("roles = '{}'")
	}

	// if searchTerm := c.QueryParam("search"); searchTerm != "" {
	// 	scope = scope.Scopes(models.UserMatch(searchTerm))
	// }

	if err := pagination.Apply(c, scope, &serializers.User{}, &users, 20); err != nil {
		return err
	}

	return c.JsonApiOK(users)
}

// GetUser admin fetch user by id
func GetUser(c middleware.AuthedContext) error {
	user, err := fetchUser(c)
	if err != nil {
		return err
	}

	return c.JsonApiOK(user)
}

// CreateUser create a user
func CreateUser(c middleware.AuthedContext) error {
	user := new(serializers.User)
	if err := c.BindAndValidate(user); err != nil {
		return err
	}

	attrs := userPermittedAttrs(c)
	attrs["roles"] = user.RolesPayload

	unscopedUser, dbErr := fetchUnscoped(user.Email)

	if dbErr != nil {
		return dbErr
	} else if unscopedUser != nil {
		// format := "3 04 PM"
		// t, e := time.Parse(format, "0 00 PM")
		// if e != nil {
		// 	return e
		// }

		if err := c.BindAndValidate(unscopedUser); err != nil {
			return err
		}

		attrs := userPermittedAttrs(c)
		attrs["roles"] = user.RolesPayload

		tx := psql.DB.Begin()
		// if err := tx.Exec("update users set deleted_at = NULL where email = ?", user.Email).Error; err != nil {
		// 	tx.Rollback()
		// 	return err
		// }
		if err := tx.Unscoped().Model(&unscopedUser).Updates(attrs).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit().Error; err != nil {
			return err
		}

		if err := unscopedUser.UpdatePassword(); err != nil {
			return nil
		}

		return c.JsonApi(unscopedUser, http.StatusCreated)

	} else {

		if err := psql.DB.Create(user).Updates(attrs).Error; err != nil {
			return err
		}

		if err := user.UpdatePassword(); err != nil {
			return nil
		}

		return c.JsonApi(user, http.StatusCreated)

	}

}

// UpdateUser admin update user
func UpdateUser(c middleware.AuthedContext) error {
	user, err := fetchUser(c)
	if err != nil {
		return err
	}

	if err := c.BindAndValidate(user); err != nil {
		return err
	}

	attrs := userPermittedAttrs(c)
	attrs["roles"] = user.RolesPayload

	if err := psql.DB.Model(&user).Updates(attrs).Error; err != nil {
		return err
	}

	if err := user.UpdatePassword(); err != nil {
		return nil
	}

	return c.JsonApiOK(user)
}

// DestroyUser admin delete user by id
func DestroyUser(c middleware.AuthedContext) error {
	if err := psql.DB.Where("id = ?", c.Param("user_id")).Delete(serializers.User{}).Error; err != nil {
		return err
	}

	return c.JsonApiOK(nil)
}

func userPermittedAttrs(c middleware.AuthedContext) map[string]interface{} {
	return c.Attrs("first_name", "last_name", "email", "roles")
}

func fetchUser(c middleware.AuthedContext) (*serializers.User, error) {
	id := c.Param("user_id")
	user := new(serializers.User)

	db := psql.DB.Preload("Company").First(user, "id = ?", id)
	if db.RecordNotFound() {
		return nil, echo.ErrNotFound
	} else if db.Error != nil {
		return nil, db.Error
	}

	return user, nil
}

func fetchUnscoped(email string) (*serializers.User, error) {
	user := new(serializers.User)
	db := psql.DB.Unscoped().Where("email =  ?", email).Find(&user)
	if db.RecordNotFound() {
		return nil, nil
	} else if db.Error != nil {
		return nil, db.Error
	}

	return user, nil
}
