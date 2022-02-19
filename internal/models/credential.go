package models

import (
	"database/sql/driver"
	"time"
)

// SourceType ...
type SourceType string

const (
	// Password ...
	Password SourceType = "password"
	// ConfirmToken ...
	ConfirmToken SourceType = "confirm_token"
)

// Credential credential data structure
type Credential struct {
	ID        int  `jsonapi:"primary,credential"`
	User      User `jsonapi:"relation,user"`
	UserID    int
	Source    SourceType `sql:"not null;type:ENUM('password', 'facebook', 'google', 'confirm_token')"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Scan ...
func (u *SourceType) Scan(value interface{}) error {
	*u = SourceType(value.([]byte))
	return nil
}

// Value ...
func (u SourceType) Value() (driver.Value, error) {
	return string(u), nil
}
