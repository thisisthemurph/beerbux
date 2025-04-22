-- name: GetSessionHistory :many
select * from session_history where session_id = ?;

-- name: CreateSessionHistory :exec
insert into session_history (session_id, member_id, event_type, event_data)
values (?, ?, ?, ?);
