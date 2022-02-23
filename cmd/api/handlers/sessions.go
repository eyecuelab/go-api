package handlers

import (
	"fmt"
	"net/http"

	"github.com/eyecuelab/go-api/cmd/api/serializers"
	"github.com/eyecuelab/go-api/cmd/middleware"
	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/go-api/internal/notifications"
	"github.com/eyecuelab/kit/db/psql"
	"github.com/eyecuelab/kit/web"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// GetSession session data endpoint
func GetSession(c middleware.AuthedContext) error {
	var sess interface{}
	if c.LoggedIn() {
		sess = &serializers.AuthSession{
			User:           c.User(),
			SomethingExtra: "abc",
		}
	} else {
		sess = &serializers.AnonSession{
			SomethingExtra: "abc",
		}
	}

	return c.JsonApiOK(sess)
}

// Login login route handler, expects email and password
func Login(c web.ApiContext) error {
	u := new(models.User)

	if err := c.BindAndValidate(u); err != nil {
		return err
	}

	authedUser, err := u.Login()
	if err != nil {
		return err
	} else if authedUser == nil {
		return c.ApiError("Invalid email or password", http.StatusUnauthorized)
	}

	token, err := authedUser.JwtToken()
	if err != nil {
		return err
	}
	c.Response().Header().Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	sess := &serializers.AuthSession{
		User:           authedUser,
		SomethingExtra: "abc",
	}
	return c.JsonApiOK(sess)
}

// Logout logout endpoint
func Logout(c middleware.AuthedContext) error {
	sess := new(serializers.AnonSession)
	// TODO: expire auth token

	return c.JsonApiOK(sess)
}

// ForgotPassword forgot password handler
func ForgotPassword(c web.ApiContext) error {
	u := new(models.User)
	if err := c.BindAndValidate(u); err != nil {
		return err
	}

	db := psql.DB.First(u, "email = ?", u.Email)
	if db.RecordNotFound() {
		return echo.ErrNotFound
	} else if db.Error != nil {
		return db.Error
	}

	if err := u.SetResetPasswordToken(); err != nil {
		return err
	}
	if err := notifications.ForgotPassword(u); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// ResetPassword reset password handler
func ResetPassword(c web.ApiContext) error {
	u := new(models.User)

	if err := c.BindAndValidate(u); err != nil {
		return err
	}

	if err := u.UpdatePasswordByResetToken(); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}
	return c.NoContent(http.StatusOK)
}
