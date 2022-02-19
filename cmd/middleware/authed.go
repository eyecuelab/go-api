package middleware

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/kit/web"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// AuthedContext web api authed context extension
	AuthedContext interface {
		web.ApiContext
		User() *models.User
		LoggedIn() bool
	}

	authedContext struct {
		web.ApiContext

		user *models.User
	}

	authedContextLookup struct {
		web.AuthedContextLookup
	}
)

func (c *authedContext) User() *models.User {
	return c.user
}

func (c *authedContext) LoggedIn() bool {
	return c.user != nil
}

func (cl *authedContextLookup) Lookup(c echo.Context) (echo.Context, error) {
	ac := &authedContext{c.(web.ApiContext), nil}

	token, ok := ac.Get("user").(*jwt.Token)
	if !ok {
		return nil, errors.New("Failed to find token")
	}
	claims, k := token.Claims.(jwt.MapClaims)
	if !k {
		return nil, errors.New("invalid claims")
	}

	id := claims["user"].(float64)
	ac.user = &models.User{ID: int(id)}

	return ac, ac.user.Find()
}

func (cl *authedContextLookup) Context(c echo.Context) echo.Context {
	return &authedContext{c.(web.ApiContext), nil}
}

// Authed returns a middleware for authed routes
func Authed() echo.MiddlewareFunc {
	middleware.ErrJWTMissing = echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
	return web.AuthedWithConfig(web.AuthedConfig{Skipper: web.AuthedSkipper()}, &authedContextLookup{})
}
