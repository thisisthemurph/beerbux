-- name: GetSession :one
select * from sessions where id = ? limit 1;

-- name: ListSessionsForUser :many
WITH paged_sessions AS (
    SELECT s.*
    FROM sessions s
             JOIN session_members sm_target ON s.id = sm_target.session_id
    WHERE sm_target.member_id = :member_id
      AND (CAST(COALESCE(:page_token, '') AS TEXT) = '' OR :page_token < s.id)
    ORDER BY s.updated_at DESC, s.id DESC
    LIMIT CASE WHEN :page_size = 0 THEN -1 ELSE :page_size END
)
SELECT
    ps.*,
    m.id AS member_id,
    m.name AS member_name,
    m.username AS member_username
FROM paged_sessions ps
         JOIN session_members sm ON ps.id = sm.session_id
         JOIN members m ON sm.member_id = m.id
ORDER BY ps.updated_at DESC, ps.id DESC;

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
