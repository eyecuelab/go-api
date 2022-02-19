package middleware

import (
	"github.com/labstack/echo"
	emw "github.com/labstack/echo/middleware"
)

var corsCfg = emw.CORSConfig{
	AllowHeaders: []string{
		echo.HeaderAuthorization,
		echo.HeaderContentType,
		"dataType",
		"X-XSRF-TOKEN",
		"X-Requested-With"},
	ExposeHeaders: []string{echo.HeaderAuthorization},
}

// Cors api middleware to handle preflight OPTIONS requests
func Cors() echo.MiddlewareFunc {
	return emw.CORSWithConfig(corsCfg)
}
