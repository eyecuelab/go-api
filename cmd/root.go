package cmd

import (
	"fmt"
	"os"

	"github.com/eyecuelab/go-api/cmd/admin"
	"github.com/eyecuelab/go-api/cmd/api"
	"github.com/eyecuelab/go-api/cmd/cron"
	"github.com/eyecuelab/go-api/cmd/storage"
	"github.com/eyecuelab/go-api/internal/notifications/email"
	"github.com/eyecuelab/kit/cmd"
	"github.com/eyecuelab/kit/db/psql"
	"github.com/eyecuelab/kit/log"
	"github.com/spf13/cobra"
)

// Root root command
var Root = &cobra.Command{
	Use:   "go-api",
	Short: "Universal binary for back-end go api processes",
}

// Init ...
func Init() {
	apiMode := fmt.Sprint(os.Getenv("API_MODE"))
	if apiMode == "admin" {
		admin.Init()
	} else {
		api.Init()
	}

	storage.Init()
	cron.Init()

	cmd.Use("api")
	if err := cmd.Init(apiMode, Root); err != nil {
		log.Fatalf("cmd.Init: %v", err)
	}
	// cobra.OnInitialize(CheckDb, initEmail)
	cobra.OnInitialize(CheckDb)
}

func Exec() {
	if err := Root.Execute(); err != nil {
		log.Fatalf("Root.Execute: %v", err)
	}
}

func CheckDb() {
	if psql.DBError != nil && !cmd.NoDb {
		log.Fatal("CheckDb(): postgres: ", psql.DBError)
	}
	// if mongo.Error != nil && !cmd.NoDb {
	// 	log.Fatal("CheckDB(): mongo:", mongo.Error)
	// }
}

func initEmail() {
	if err := email.InitEmail(); err != nil {
		log.Fatalf("email.InitEmail: %v", err)
	}
}
