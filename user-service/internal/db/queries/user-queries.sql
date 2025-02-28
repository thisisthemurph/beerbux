-- name: GetUser :one
select * from users where id = ? limit 1;

-- name: CreateUser :one
insert into users (id, username, bio) values (?, ?, ?)
returning *;

-- name: UpdateUser :one
update users
set username = ?, bio = ? where id = ?
returning *;
