package routes

import (
	"fmt"

	"github.com/eyecuelab/go-api/cmd/admin/handlers"
	"github.com/eyecuelab/go-api/cmd/middleware"
	"github.com/eyecuelab/kit/web"
)

// Init ...
func Init() {
	get("health", web.Healthz)

	getAuthed("", handlers.GetSession)
	post("login", handlers.Login)
	deleteAuthed("logout", handlers.Logout)

	getAuthed("users", handlers.ListUsers)
	postAuthed("users", handlers.CreateUser)
	getAuthed("users/:user_id", handlers.GetUser)
	patchAuthed("users/:user_id", handlers.UpdateUser)
	deleteAuthed("users/:user_id", handlers.DestroyUser)
}

func get(path string, hf web.HandlerFunc) *web.Route {
	return addRoute(path, "GET", hf)
}

func post(path string, hf web.HandlerFunc) *web.Route {
	return addRoute(path, "POST", hf)
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

// AuthedHandlerFunc ...
type AuthedHandlerFunc func(middleware.AuthedContext) error

func wrapAuthedHandler(f AuthedHandlerFunc) web.HandlerFunc {
	return func(c web.ApiContext) error { return f(c.(middleware.AuthedContext)) }
}
