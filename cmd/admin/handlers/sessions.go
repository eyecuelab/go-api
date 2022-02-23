package handlers

import (
	"fmt"
	"net/http"

	"github.com/eyecuelab/go-api/cmd/admin/serializers"
	"github.com/eyecuelab/go-api/cmd/middleware"
	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/kit/web"
	"github.com/labstack/echo"
)

// GetSession admin session data endpoint
func GetSession(c middleware.AuthedContext) error {
	var sess interface{}
	if c.LoggedIn() && c.User().IsAdmin() {
		// user := serializers.User{User: *c.User()}
		sess = &serializers.AuthSession{
			// User: &user,
			User: c.User(),
		}
	} else {
		sess = new(serializers.AnonSession)
	}

	return c.JsonApiOK(sess)
}

// Login admin login route handler, expects email and password
func Login(c web.ApiContext) error {
	u := new(models.User)

	if err := c.BindAndValidate(u); err != nil {
		return err
	}

	authedUser, err := u.Login()
	if err != nil {
		return err
	} else if authedUser == nil || !authedUser.IsAdmin() {
		return c.JsonAPIError("Invalid email or password", http.StatusUnauthorized, "email")
	}

	token, err := authedUser.JwtToken()
	if err != nil {
		return err
	}
	c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	// user := serializers.User{User: *authedUser}

	sess := &serializers.AuthSession{
		// User: &user,
		User: authedUser,
	}
	return c.JsonApiOK(sess)
}

// Logout admin logout endpoint
func Logout(c middleware.AuthedContext) error {
	sess := new(serializers.AnonSession)
	// TODO: expire auth token

	return c.JsonApiOK(sess)
}
