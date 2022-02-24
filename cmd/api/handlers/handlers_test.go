package handlers_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/eyecuelab/go-api/cmd/api/handlers"
	"github.com/eyecuelab/go-api/cmd/api/routes"

	"github.com/eyecuelab/go-api/cmd/middleware"
	"github.com/eyecuelab/go-api/cmd/storage"
	"github.com/eyecuelab/go-api/internal/models"

	// "github.com/eyecuelab/go-api/internal/notifications/email"
	"github.com/spf13/cobra"

	api "github.com/eyecuelab/kit/cmd/api"
	"github.com/eyecuelab/kit/config"
	"github.com/eyecuelab/kit/db/psql"
	"github.com/eyecuelab/kit/log"
	"github.com/eyecuelab/kit/web/server"
)

func TestMain(m *testing.M) {
	if err := setupTests(); err != nil {
		log.Fatalf("setupTests: %v", err)
	}
	code := m.Run()
	if err := cleanupTests(); err != nil {
		log.Infof("cleanupTests: %v", err)
	}
	os.Exit(code)
}

func setupTests() error {
	if err := initCobra(); err != nil {
		return fmt.Errorf("initCobra: %v", err)
	}
	return nil
}

func initCobra() error {
	const waitTimeForCobra = 50 * time.Millisecond

	routes.Init()

	cobra.OnInitialize(
		func() {
			server.AddMiddleWare(middleware.Authed())
		},
		// initEmail
	)

	workingDirPath := os.Getenv("APP_WORKING_DIR")
	if workingDirPath == "" {
		workingDirPath = "/app"
	}
	err := config.Load("test", fmt.Sprintf("%s/config", workingDirPath))
	if err != nil {
		return err
	}

	go func() {
		api.ApiCmd.Execute()
	}()
	time.Sleep(waitTimeForCobra)

	psql.ConnectDB()
	if psql.DBError != nil {
		log.Fatal("DBError: ", psql.DBError)
	}

	storage.Migrate()
	storage.Seed()

	return nil
}

func cleanupTests() error {
	storage.Clear()
	return nil
}

func authToken(t *testing.T, email string) string {
	user := new(models.User)

	db := psql.DB.First(user, "email = ?", email)
	if db.RecordNotFound() {
		t.Errorf("Error generating a token, can't find user by email")
	} else if db.Error != nil {
		return ""
	}

	token, err := user.JwtToken()
	if err != nil {
		t.Errorf("Error generating a token")
	}
	return token
}
