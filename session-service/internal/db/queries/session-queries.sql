-- name: GetSession :one
select * from sessions where id = ? limit 1;

-- name: CreateSession :one
insert into sessions (id, name, owner_id)
values (?, ?, ?)
returning *;

-- name: AddMemberToSession :exec
insert into session_members (session_id, user_id)
values (?, ?);

-- name: UpdateSession :one
update sessions
set name = ?
where id = ?
returning *;
