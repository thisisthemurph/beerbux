# user-service

## Running locally

```py
GOOSE_DRIVER=sqlite
GOOSE_DBSTRING=./users.sqlite
GOOSE_MIGRATION_DIR=./internal/db/migrations

DB_DRIVER=sqlite
DB_URI=./users.sqlite

ENVIRONMENT=development
USER_SERVER_ADDRESS=:50051
```

## Manual testing

```shell
grpcurl -plaintext -d '{"username": "mike", "bio": "this is the murph"}' localhost:50051 User.CreateUser
```