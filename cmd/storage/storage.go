package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/kit/db/psql"
	"github.com/google/jsonapi"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/eyecuelab/kit/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StorageCommand main data storage managing command
var StorageCommand = &cobra.Command{
	Use:   "storage",
	Short: "data storage managing",
}

// Init ...
func Init() {
	cmd.Add(StorageCommand)
	StorageCommand.AddCommand(StorageClearCommand)
	StorageCommand.AddCommand(StorageSeedCommand)
}

// StorageClearCommand remove data from db
var StorageClearCommand = &cobra.Command{
	Use:   "clear",
	Short: "clear data storage",
	Run:   clearCmd,
}

func clearCmd(cmd *cobra.Command, args []string) {
	Clear()
}

type clearQuery struct {
	TruncateQ string
	SequenceQ string
}

// Clear clear fixtures
func Clear() {
	var qs []*clearQuery
	tableNamesSelect := "select 'public.' || table_name AS input_table_name from information_schema.tables where table_schema = 'public' and table_type = 'BASE TABLE' and table_name not like 'pg_%' and table_name != 'schema_migrations'"
	sql := fmt.Sprintf("SELECT 'TRUNCATE ' || input_table_name || ' CASCADE;' AS truncate_q, 'ALTER SEQUENCE IF EXISTS ' || input_table_name || '_id_seq RESTART WITH 1;' AS sequence_q FROM(%s) AS information;", tableNamesSelect)
	if err := psql.DB.Raw(sql).Scan(&qs).Error; err != nil {
		fmt.Printf("Error: %+v\n\n", err)
	}

	for _, q := range qs {
		if err := psql.DB.Exec(fmt.Sprintf("%s%s", q.TruncateQ, q.SequenceQ)).Error; err != nil {
			fmt.Printf("Error: %+v\n\n", err)
		}
	}
}

func Migrate() {
	url := viper.GetString("database_url")

	m, err := migrate.New("file:///app/data/migrations", url)
	if err != nil {
		log.Fatal("Migrate init error: ", err)
	}
	if err := m.Up(); err != nil {
		log.Printf("Migrate up: %v", err)
	}
}

// StorageSeedCommand seed data storage with initial data
var StorageSeedCommand = &cobra.Command{
	Use:   "seed",
	Short: "seed database with initial data",
	Run:   seedCmd,
}

func seedCmd(cmd *cobra.Command, args []string) {
	Seed()
}

// Seed default fixtures
func Seed() {
	jsonAPI("companies", new(models.Company))
	users()
}

func jsonAPI(name string, i interface{}) {
	data := seedData(name)
	items, err := jsonapi.UnmarshalManyPayload(bytes.NewReader(data), reflect.TypeOf(i))
	if err != nil {
		fmt.Printf("Error: %+v\n\n", err)
	}
	for _, item := range items {
		if err := psql.DB.Create(item).Error; err != nil {
			fmt.Printf("Error: %+v\n\n", err)
		}
	}
}

func users() {
	var items []models.User
	if err := json.Unmarshal(seedData("users"), &items); err != nil {
		fmt.Printf("Error: %+v\n\n", err)
	}
	for _, item := range items {
		if err := item.RegisterWithPassword(); err != nil {
			fmt.Printf("Error: %+v\n\n", err)
		}
	}
	// sessions := []string{
	// 	"insert into sessions (id, user_id) values (1, 1)",
	// 	"insert into sessions (id, user_id) values (1, 1)",
	// 	"insert into sessions (id, user_id) values (1, 2)",
	// }
	// for _, q := range sessions {
	// 	if err := psql.DB.Exec(q).Error; err != nil {
	// 		panic(err)
	// 	}
	// }
}

func seedData(t string) []byte {
	// ex, err := os.Executable()
	// if err != nil {
	// 	panic(err)
	// }
	// exPath := filepath.Dir(ex)

	// gopath := os.Getenv("GOPATH")
	// if gopath == "" {
	// 	gopath = build.Default.GOPATH
	// }
	// path := fmt.Sprintf("%s/src/github.com/eyecuelab/go-api/cmd/storage/fixtures/%s.json", gopath, t)
	path := fmt.Sprintf("/app/cmd/storage/fixtures/%s.json", t)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: %+v\n\n", err)
	}

	return raw
}

func clearTable(name string) {
	if err := psql.DB.Exec(fmt.Sprintf("DELETE FROM %s; ALTER SEQUENCE IF EXISTS %s_id_seq RESTART WITH 1;", name, name)).Error; err != nil {
		panic(err)
	}
}

func ensureENV(name string) {
	if viper.GetString(name) == "" {
		panic(fmt.Sprintf("Error: %s env var is missing", name))
	}
}
