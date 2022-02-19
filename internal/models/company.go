package models

import (
	"time"
)

// Company company data structure
type Company struct {
	ID          int       `jsonapi:"primary,company" gorm:"primary_key"`
	Name        string    `jsonapi:"attr,name"`
	Slug        string    `jsonapi:"attr,slug"`
	Description string    `jsonapi:"attr,description"`
	Users       []*User   `jsonapi:"relation,users,omitempty" gorm:"save_associations:false"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
