api: # start main api
	docker compose up api

admin: # start admin api
	docker compose up admin

db_create_migration: # create a new migration file: make db_create_migration file=create_users
	docker compose --profile tools run create-migration $(file)

db_migrate: # run db migrations
	docker compose --profile tools run migrate

seed: # seed database with data
	docker compose run --rm api go run main.go storage seed

psql_console: # run psql console
	docker compose exec db psql -U local-dev -d go_api_dev

dep: # install missing dependencies
	docker compose run --rm api go mod tidy

bash: # run bash inside the api container
	docker compose run --rm api bash

test_api: # run main api integration tests
	docker compose run --rm test go test ./cmd/api/...

test_models: # run models unit tests
	docker compose run --rm test go test ./internal/models/...

cron: # run cron jobs
	docker compose up cron

### Plain (commands to run plain go commands without docker compose)

plain_api:
	API_MODE=api go run main.go api --port=8000

plain_admin_api:
	API_MODE=admin go run main.go api --port=8010

plain_migrate:
	API_MODE=api go run main.go migrate

plain_rollback:
	API_MODE=api go run main.go migrate --down

plain_clear:
	API_MODE=api go run main.go storage clear

plain_dev_db_prepare:
	psql postgres -c "drop database if exists go_api_dev;"
	psql postgres -c "create database go_api_dev;"
	make plain_migrate


### Misc

compose_build:
	docker compose build

compose_air_init:
	docker compose run --rm api air init
