-- name: GetSession :one
select * from sessions where id = ? limit 1;

-- name: CreateSession :one
insert into sessions (id, name)
values (?, ?)
returning *;

-- name: UpsertMember :exec
insert into members (id, name, username)
values (?, ?, ?)
on conflict do nothing;

-- name: AddSessionMember :exec
insert into session_members (session_id, member_id, is_owner)
values (?, ?, ?)
on conflict do nothing;

-- name: UpdateSession :one
update sessions
set name = ?
where id = ?
returning *;
