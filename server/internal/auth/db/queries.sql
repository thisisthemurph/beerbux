-- name: GetUser :one
select * from users where id = $1 limit 1;

-- name: GetUserByUsername :one
select * from users where username = $1 limit 1;

-- name: UserWithUsernameExists :one
select exists(select 1 from users where username = $1);

-- name: CreateUser :one
insert into users (name, username, hashed_password)
values ($1, $2, $3)
returning *;

-- name: UpdatePassword :exec
update users
set hashed_password = $1
where id = $2;

-- name: RegisterRefreshToken :exec
insert into refresh_tokens (user_id, hashed_token, expires_at)
values ($1, $2, $3);

-- name: GetRefreshTokensByUserID :many
select *
from refresh_tokens
where user_id = $1
  and revoked = false
  and expires_at > now();

-- name: DeleteRefreshToken :exec
delete from refresh_tokens where id = $1;

-- name: InvalidateRefreshToken :exec
update refresh_tokens
set revoked = true
where id = $1;