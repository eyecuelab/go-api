package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// TimeZeroSQL ...
func TimeZeroSQL(field string) string {
	return fmt.Sprintf("%s is null or %s <= 'epoch'", field, field)
}

// TimeSetSQL ...
func TimeSetSQL(field string) string {
	return fmt.Sprintf("%s is not null and %s > 'epoch'", field, field)
}

// Now ...
func Now() *time.Time {
	now := time.Now()
	return &now
}

// GenerateSlug generate a slug from a string
func GenerateSlug(s string) string {
	regexpCharsOut := regexp.MustCompile("[^a-z0-9-_]")
	regexpMultiDash := regexp.MustCompile("-+")
	slug := strings.ToLower(s)
	slug = regexpCharsOut.ReplaceAllString(slug, "-")
	slug = regexpMultiDash.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}
