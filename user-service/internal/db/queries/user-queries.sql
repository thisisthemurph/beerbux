-- name: GetUser :one
select * from users where id = ? limit 1;

-- name: GetUserByUsername :one
select * from users where username = ? limit 1;

-- name: CreateUser :one
insert into users (id, name, username, bio) values (?, ?, ?, ?)
returning *;

-- name: UpdateUser :one
update users
set name = ?, username = ?, bio = ? where id = ?
returning *;

-- name: UpdateUserTotals :exec
update users
set credit = ?, debit = ?, net = ?
where id = ?;

-- name: GetUserBalances :one
select credit, debit, net from users where id = ? limit 1;
