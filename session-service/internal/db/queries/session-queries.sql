-- name: GetSession :one
select * from sessions where id = ? limit 1;

-- name: ListSessionsForUser :many
select
    s.*,
    m.id as member_id,
    m.name as member_name,
    m.username as member_username
from sessions s
join session_members sm_target on s.id = sm_target.session_id
join session_members sm on s.id = sm.session_id
join members m on sm.member_id = m.id
where sm_target.member_id = :member_id
and (cast(coalesce(:page_token, '') as text) = '' or :page_token < s.id)
order by s.updated_at desc, s.id desc
limit case when :page_size = 0 then -1 else :page_size end;

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
