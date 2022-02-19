package handlers

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type fileImport struct {
	FileURL string `jsonapi:"attr,file_url"`
}

func applyFilter(c echo.Context, filter string, scope *gorm.DB) *gorm.DB {
	q := c.QueryParam(filter)
	if q != "" {
		return scope.Where(fmt.Sprintf("%s = ?", filter), q)
	}
	return scope
}
