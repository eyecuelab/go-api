module github.com/eyecuelab/go-api

go 1.16

// replace github.com/eyecuelab/kit v0.0.0-20220218184705-7be9103a1e31 => ../kit

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/eyecuelab/kit v0.0.0-20220218231547-ea0ee21d6634
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/google/jsonapi v1.0.0
	github.com/jinzhu/gorm v1.9.16
	github.com/joho/godotenv v1.4.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/lib/pq v1.10.4
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
)
