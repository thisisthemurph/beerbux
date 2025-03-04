-- name: GetSession :one
select * from sessions where id = ? limit 1;

-- name: CreateSession :one
insert into sessions (id, name)
values (?, ?)
returning *;

-- name: GetMember :one
select * from members where id = ? limit 1;

-- name: ListMembers :many
select m.*
from members m
join session_members sm on m.id = sm.member_id
where sm.session_id = ?;

-- name: UpsertMember :exec
insert into members (id, name, username)
values (?, ?, ?)
on conflict(id) do update
set name = excluded.name,
    username = excluded.username,
    updated_at = current_timestamp;

-- name: UpdateMember :exec
update members
set name = ?,
    username = ?,
    updated_at = current_timestamp
where id = ?;

-- name: AddSessionMember :exec
insert into session_members (session_id, member_id, is_owner)
values (?, ?, ?)
on conflict do nothing;

-- name: UpdateSession :one
update sessions
set name = ?,
    updated_at = current_timestamp
where id = ?
returning *;
