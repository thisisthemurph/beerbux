-- name: Create :one
insert into users (name, username, hashed_password)
values ($1, $2, $2)
returning *;

-- name: UpdatePassword :exec
update users
set hashed_password = $1
where id = $2;
