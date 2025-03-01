-- name: CreateUser :exec
insert into user_details (id, username)
values (?, ?);

-- name: GetUser :one
select * from user_details where id = ?;
