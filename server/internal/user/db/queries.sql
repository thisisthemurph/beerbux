-- name: UpdateUser :exec
update users
set name = $2, username = $3
where id = $1;

-- name: GetUserIDByUsername :one
select id from users where username = $1 limit 1;