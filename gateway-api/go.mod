module github.com/thisisthemurph/beerbux/gateway-api

go 1.24.0

replace (
	github.com/thisisthemurph/beerbux/auth-service => ../auth-service
	github.com/thisisthemurph/beerbux/shared => ../shared
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/thisisthemurph/beerbux/auth-service v0.0.0-20250314234135-7e9581d5a62d
	google.golang.org/grpc v1.71.0
)

require (
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250313205543-e70fdf4c4cb4 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
