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

**Get user**

Returns the user information and the calculated net balance (the number of beers they owe or are owed).

```shell
grpcurl -plaintext -d '{"user_id": "460e1637-8c7d-48c4-9e3f-58e880f77fde"}' localhost:50051 User.GetUser
```

**Create user**

Creates and returns the newly created user.

```shell
grpcurl -plaintext -d '{"username": "user.name", "name": "User Name"}' localhost:50051 User.CreateUser
```

**Update user**

Updates the user and returns the updated user.

```shell
grpcurl -plaintext -d '{"user_id": "460e1637-8c7d-48c4-9e3f-58e880f77fde", "username": "michael.murphy", "name": "Michael"}' localhost:50051 User.UpdateUser
```