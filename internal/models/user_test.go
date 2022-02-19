package models_test

import (
	"testing"

	"github.com/eyecuelab/go-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserFullName(t *testing.T) {
	user := new(models.User)
	fetchByID(t, user, 1)

	fullName := user.FullName()

	assert.True(t, fullName == "Billy Bob")
}
