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

func sessionLinks() jsonapi.Links {
	return jsonapi.Links{
		"self": meta.APIURL(""),
		// "users": meta.APIURL("anon/users")
	}
}

func addCommonSessionActions(ah *meta.ActionHolder) {
	ah.AddAction(meta.POST, "some_other_common_action", "something").
		Field("field1", meta.InputText, nil, true).
		Field("field2", meta.InputText, nil, true)
}

// JSONAPILinks links for the session json
func (sess AnonSession) JSONAPILinks() *jsonapi.Links {
	links := sessionLinks()
	return &links
}

// JSONAPIMeta links for the session json
func (sess AnonSession) JSONAPIMeta() *jsonapi.Meta {
	var ah meta.ActionHolder
	ah.AddAction(meta.POST, "login", "login").
		Field("email", meta.InputText, nil, true).
		Field("password", meta.InputText, nil, true)

	addCommonSessionActions(&ah)

	return ah.RenderActions()
}

// AuthSession authenticated user session data
type AuthSession struct {
	ID   string       `jsonapi:"primary,session"`
	User *models.User `jsonapi:"relation,user"`
}

// JSONAPILinks links for the session json
func (sess AuthSession) JSONAPILinks() *jsonapi.Links {
	links := sessionLinks()
	links["users"] = meta.APIURL("users")
	links["profile"] = meta.APIURL("profile")
	return &links
}

// JSONAPIMeta links for the session json
func (sess AuthSession) JSONAPIMeta() *jsonapi.Meta {
	var ah meta.ActionHolder
	ah.AddAction(meta.DELETE, "logout", "logout")

	ah.AddAction(meta.PATCH, "update_profile", "profile").
		Field("first_name", meta.InputText, sess.User.FirstName, true).
		Field("last_name", meta.InputText, sess.User.LastName, true)

	addCommonSessionActions(&ah)

	return ah.RenderActions()
}
