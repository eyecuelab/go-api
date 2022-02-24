# Go template API [![CircleCI](https://circleci.com/gh/eyecuelab/go-api/tree/main.svg?style=svg&circle-token=0788ba6ad5257362f3736f4546dd9dde43697bb7)](https://circleci.com/gh/eyecuelab/go-api/tree/main)

Eyecue Lab golang api template, could serve as a starting point when developing a Web API.

## Features

- Live reload with docker compose for local development
- API with examples for creating routes and handlers
- JSONAPI schema with custom actions definitions for the client app to avoid hardcoding
- Optional setup to have a separate content management API
- Reusable models between different APIs
- Examples of how to preload related data
- Postgres database migrations
- Data seeding functionality for development and test database prep
- Cron jobs setup for periodic tasks
- Integrations with 3rd party apis: airbrake, sendgrid
- Automated testing with Circleci
- API endpoints integration tests examples
- Models unit tests examples

## Usage

Make commands rely on Docker Engine, please make sure you have it [installed](https://docs.docker.com/desktop/mac/install/)

To start API locally: `make api` and `curl http://localhost:8000 | json_pp`

To use all available endpoints, including sign in etc... you'll need to migrate and seed database, refer to make commands below:

<details>
<summary><b>Show more make commands</b></summary>

* When adding new imports within the app, to update go.mod and go.sum, run:
```sh
make dep
```

* Run database migrations:
```sh
make db_migrate
```

* There are predefined fixtures like users, companies etc... inside `cmd/storage`, to seed database with that data:
```sh
make db_seed
```

* Create a new database migration file:
```sh
make db_create_migration file=create_bananas
```

* Drop dev database. This will stop docker volumes and delete the database volume, after running this command to get your database restored to use API, you'll need to run a migration and seed commands
```sh
make db_drop
```

* Open psql console inside postgres container:
```sh
make psql_console
```

* Run bash console inside API container:
```sh
make bash
```

* Run API integration tests:
```sh
make test_api
```

* Run models unit tests:
```sh
make test_models
```

* There is a standalone Admin CMS API `cmd/admin` that uses the same Models layer as a main API, to start admin API:
```sh
make admin
```

* Cron jobs example is in `cmd/cron`, to run cron jobs:
```sh
make cron
```

</details>

## API Endpoints

List of existing API endpoints and curl examples

<details>
<summary><b>Show endpoints</b></summary>

* Health check
```sh
curl http://localhost:8000/health
```

* Versions
```sh
curl http://localhost:8000/version
```

* Anon user init data. Provides data in jsonapi schema for a guest user
```sh
curl http://localhost:8000/ | json_pp
```

* Sign In. Response headers will hold authentication jwt that can be used to call auth endpoints
```sh
curl -H "Content-Type: application/json" -X POST -d '{"email": "user1@example.com", "password": "goapi123"}' http://localhost:8000/login | json_pp
```

* Authenticated user init data. Provides data in jsonapi schema for an authenticated user
```sh
curl -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoxfQ.VEy5T2jj4rIr2Sfs_mif0uKTi9GyX2eAi4_QYcL673o" http://localhost:8000/ | json_pp
```

</details>

## TODO

* Queue setup example
* Protocol Buffer API
* Generate docs, godoc
* ...
