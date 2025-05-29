-- name: GetUser :one
select * from users where id = $1 limit 1;

-- name: GetUserByUsername :one
select * from users where username = $1 limit 1;

-- name: GetUserByEmail :one
select * from users where email = $1 limit 1;

-- name: UserWithUsernameExists :one
select exists(select 1 from users where username = $1);

-- name: CreateUser :one
insert into users (name, username, email, hashed_password)
values ($1, $2, $3, $4)
returning *;

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

-- name: InitializePasswordUpdate :exec
update users
set update_hashed_password = $2,
    password_update_otp = $3,
    password_update_requested_at = now()
where id = $1;

-- name: UpdatePassword :exec
with updated as (
    select id, update_hashed_password
    from users updated_users
    where updated_users.id = $1 and updated_users.password_update_requested_at is not null
)
update users u
set
    hashed_password = updated.update_hashed_password,
    update_hashed_password = null,
    password_update_otp = null,
    password_update_requested_at = null,
    password_last_updated_at = now()
from updated
where u.id = updated.id;

-- name: InitialiseUpdateEmail :exec
update users
set update_email = $2,
    email_update_otp = $3,
    email_last_updated_at = now()
where id = $1;

-- name: UpdateEmail :exec
with updated as (
    select id, update_email
    from users updated_users
    where updated_users.id = $1 and updated_users.email_update_requested_at is not null
)
update users u
set email = updated.update_email,
    update_email = null,
    email_update_otp = null,
    email_update_requested_at = null,
    email_last_updated_at = now()
where u.id = updated.id;