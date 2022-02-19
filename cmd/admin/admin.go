package admin

import (
	"github.com/eyecuelab/kit/web/server"
	"github.com/eyecuelab/go-api/cmd/admin/routes"
	"github.com/eyecuelab/go-api/cmd/middleware"
	"github.com/spf13/cobra"
)

var ApiCmd *cobra.Command

func Init() {
	routes.Init()
	cobra.OnInitialize(func() {
		server.AddMiddleWare(middleware.Authed())
		server.AddMiddleWare(middleware.Admin())
		server.AddMiddleWare(middleware.Cors())
	})
}
