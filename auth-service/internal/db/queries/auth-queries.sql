-- name: RegisterUser :one
insert into users (id, username, hashed_password)
values (?, ?, ?)
returning *;

-- name: GetUserByID :one
select * from users where id = ?;

-- name: GetUserByUsername :one
select * from users where username = ?;

-- name: UpdateUserUsername :one
update users
set username = ?,
    updated_at = current_timestamp
where id = ?
returning *;

-- name: UpdateUserPassword :exec
update users
set hashed_password = ?,
    updated_at = current_timestamp
where id = ?
returning *;

-- name: UserWithUsernameExists :one
select exists(select 1 from users where username = ?);
