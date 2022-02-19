package models

import (
	"github.com/jinzhu/gorm"
)

// UsersWithRole users by role scope
func UsersWithRole(role string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("? = ANY(users.roles)", role)
	}
}

// Unscoped ...
func Unscoped(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}
