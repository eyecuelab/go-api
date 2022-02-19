package serializers

import (
	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/kit/web/meta"
	"github.com/google/jsonapi"
)

// User user serializer
type User struct {
	models.User
}

// Users users list serializer
type Users []*User

// JSONAPILinks jsonapi links
func (i Users) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": meta.APIURL("users"),
	}
}

// JSONAPIMeta jsonapi meta
func (i Users) JSONAPIMeta() *jsonapi.Meta {
	var ah meta.ActionHolder
	ah.AddAction(meta.GET, "index", "users")

	return ah.RenderActions()
}
