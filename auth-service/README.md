# auth-service

## Manual testing

**Register a new user**

A user can be registered by providing the name, username, password, and verification password.

```shell
grpcurl -plaintext -d '{"name": "Mike", "username": "thisisthemurph", "password": "password", "verification_password": "password"}' localhost:50054 Auth.Signup
```

- Creates a new user in the `users` table
- Sends a message to the `auth.user.registered` Kafka topic
- Returns a JWT token and refresh token

**Login**

A user can log in by providing the username and password.

```shell
grpcurl -plaintext -d '{"username": "thisisthemurph", "password": "password"}' localhost:50054 Auth.Login
```

**Refresh token**

A user can refresh their token by providing the refresh token and user ID.

```shell
grpcurl -plaintext -d '{"user_id": "460e1637-8c7d-48c4-9e3f-58e880f77fde", "refresh_token": "..." }' localhost:50054 Auth.RefreshToken
```

- Returns a new JWT token and refresh token