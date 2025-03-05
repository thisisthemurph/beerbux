module github.com/thisisthemurph/beerbux/transaction-service

go 1.24.0

replace (
	github.com/thisisthemurph/beerbux/session-service => ../session-service
	github.com/thisisthemurph/beerbux/user-service => ../user-service
)

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/nats-io/nats.go v1.39.1
	github.com/thisisthemurph/beerbux/session-service v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
)
