package routes

import (
	"fmt"

	"github.com/eyecuelab/go-api/cmd/api/handlers"
	"github.com/eyecuelab/go-api/cmd/middleware"
	"github.com/eyecuelab/kit/web"
)

// Init ...
func Init() {
	get("health", web.Healthz)
	get("version", handlers.Version)

	getAuthed("", handlers.GetSession)
	post("login", handlers.Login)
	postAuthed("logout", handlers.Logout)
	getAuthed("profile", handlers.GetUserProfile)
	patchAuthed("profile", handlers.UpdateUserProfile)
}

func get(path string, hf web.HandlerFunc) *web.Route {
	return addRoute(path, "GET", hf)
}

func post(path string, hf web.HandlerFunc) *web.Route {
	return addRoute(path, "POST", hf)
}

func patch(path string, hf web.HandlerFunc) *web.Route {
	return addRoute(path, "PATCH", hf)
}

func del(path string, hf web.HandlerFunc) *web.Route {
	return addRoute(path, "DELETE", hf)
}

func getAuthed(path string, hf AuthedHandlerFunc) *web.Route {
	return addAuthedRoute(path, "GET", hf)
}

func postAuthed(path string, hf AuthedHandlerFunc) *web.Route {
	return addAuthedRoute(path, "POST", hf)
}

func patchAuthed(path string, hf AuthedHandlerFunc) *web.Route {
	return addAuthedRoute(path, "PATCH", hf)
}

func deleteAuthed(path string, hf AuthedHandlerFunc) *web.Route {
	return addAuthedRoute(path, "DELETE", hf)
}

func addRoute(path string, t string, hf web.HandlerFunc) *web.Route {
	return web.AddRoute(fmt.Sprintf("/%s", path)).Handle(string(t), hf)
}

func addAuthedRoute(path string, t string, ahf AuthedHandlerFunc) *web.Route {
	return web.AddRoute(fmt.Sprintf("/%s", path)).Handle(string(t), wrapAuthedHandler(ahf))
}

// AuthedHandlerFunc authed handler type
type AuthedHandlerFunc func(middleware.AuthedContext) error

func wrapAuthedHandler(f AuthedHandlerFunc) web.HandlerFunc {
	return func(c web.ApiContext) error { return f(c.(middleware.AuthedContext)) }
}
