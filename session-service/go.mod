module github.com/thisisthemurph/beerbux/session-service

go 1.24.0

//replace github.com/thisisthemurph/beerbux/user-service => ../user-service
replace github.com/thisisthemurph/beerbux/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/pressly/goose/v3 v3.24.1
	github.com/segmentio/kafka-go v0.4.47
	github.com/stretchr/testify v1.10.0
	github.com/thisisthemurph/beerbux/shared v0.0.0-00010101000000-000000000000
	github.com/thisisthemurph/beerbux/user-service v0.0.0-20250308085229-299b031fb2d1
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
	modernc.org/sqlite v1.36.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.61.13 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.8.2 // indirect
)
