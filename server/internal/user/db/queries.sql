-- name: UpdateUser :one
update users
set name = $2, username = $3
where id = $1
returning name, username;
