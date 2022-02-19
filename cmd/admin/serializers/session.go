package serializers

import (
	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/kit/web/meta"
	"github.com/google/jsonapi"
)

// AnonSession anon session data
type AnonSession struct {
	ID string `jsonapi:"primary,session"`
}

// JSONAPILinks links for the session json
func (sess AnonSession) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": meta.APIURL("session"),
	}
}

// JSONAPIMeta links for the session json
func (sess AnonSession) JSONAPIMeta() *jsonapi.Meta {
	var ah meta.ActionHolder
	ah.AddAction(meta.POST, "login", "login").
		Field("email", meta.InputText, "", true).
		Field("password", meta.InputPass, "", true)

	return ah.RenderActions()
}

// AuthSession authenticated user session data
type AuthSession struct {
	ID string `jsonapi:"primary,session"`
	// User *User  `jsonapi:"relation,user"`
	User *models.User `jsonapi:"relation,user"`
}

// JSONAPILinks links for the session json
func (sess AuthSession) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self":  meta.APIURL("session"),
		"users": meta.APIURL("users"),
	}
}

// JSONAPIMeta links for the authed session json
func (sess AuthSession) JSONAPIMeta() *jsonapi.Meta {
	var ah meta.ActionHolder
	ah.AddAction(meta.DELETE, "logout", "logout")

	return ah.RenderActions()
}
