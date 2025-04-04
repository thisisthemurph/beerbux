# user-service

## Manual testing

**Get user**

Returns the user information and the calculated net balance (the number of beers they owe or are owed).

```shell
grpcurl -plaintext -d '{"user_id": "10473635-01d4-4e2a-b809-8fce66031ace"}' localhost:50051 user.service.User.GetUser
```

**Get user**

Returns the user information and the calculated net balance (the number of beers they owe or are owed).

```shell
grpcurl -plaintext -d '{"user_id": "10473635-01d4-4e2a-b809-8fce66031ace"}' localhost:50051 user.service.User.GetUserBalance
```

**Update user**

Updates the user and returns the updated user.

```shell
grpcurl -plaintext -d '{"user_id": "10473635-01d4-4e2a-b809-8fce66031ace", "username": "michael.murphy", "name": "Michael"}' localhost:50051 user.service.User.UpdateUser
```