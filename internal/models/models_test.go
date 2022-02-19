package models_test

import (
	"os"
	"testing"

	"github.com/eyecuelab/go-api/cmd/storage"
	"github.com/eyecuelab/kit/config"
	"github.com/eyecuelab/kit/db/psql"
	"github.com/eyecuelab/kit/log"
)

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatalf("Models test setup error: %v", err)
	}
	r := m.Run()
	cleanup()
	os.Exit(r)
}

func setup() error {
	// env.Load("test")
	//
	// assets.Manager = &assets.AssetManager{Get: data.Asset, Dir: data.AssetDir}
	if err := config.Load("test", "/app/config"); err != nil {
		// if err := config.Load("test", "../../config"); err != nil {
		return err
	}

	psql.ConnectDB()
	if psql.DBError != nil {
		log.Fatal("DBError: ", psql.DBError)
	}

	storage.Migrate()
	storage.Seed()

	return nil
}

func cleanup() {
	storage.Clear()
}

func fetchByID(t *testing.T, i interface{}, id int) {
	db := psql.DB.First(i, "id = ?", id)
	if db.RecordNotFound() {
		t.Errorf("Can't find record with ID: %d\n", id)
	} else if db.Error != nil {
		t.Errorf("Error looking up a record: %+v\n", db.Error)
	}
}
