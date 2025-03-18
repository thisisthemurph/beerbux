module github.com/thisisthemurph/beerbux/gateway-api

go 1.24.0

replace (
	github.com/thisisthemurph/beerbux/auth-service => ../auth-service
	github.com/thisisthemurph/beerbux/shared => ../shared
	github.com/thisisthemurph/beerbux/user-service => ../user-service
)

require (
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.10.0
	github.com/thisisthemurph/beerbux/auth-service v0.0.0-00010101000000-000000000000
	github.com/thisisthemurph/beerbux/user-service v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.71.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250313205543-e70fdf4c4cb4 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
