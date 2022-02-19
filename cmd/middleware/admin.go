package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

// Admin admin routes middleware
func Admin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ac := c.(AuthedContext)

			if ac.LoggedIn() && !ac.User().IsAdmin() {
				return echo.NewHTTPError(http.StatusUnauthorized, "Must have admin access")
			}

			return next(ac)
		}
	}
}
