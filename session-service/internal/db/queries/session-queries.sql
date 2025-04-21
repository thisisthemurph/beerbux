-- name: GetSession :one
select s.id, s.name, s.is_active, s.created_at, s.updated_at, cast(coalesce(sum(l.amount), 0) as real) as total
from sessions s
left join transactions t on s.id = t.session_id
left join transaction_lines l on t.id = l.transaction_id
where s.id = ?
group by s.id, s.name, s.is_active, s.created_at, s.updated_at;

-- name: ListSessionsForUser :many
with paged_sessions AS (
    select s.*, cast(coalesce(sum(l.amount), 0) as real) as total_amount
    from sessions s
    join session_members sm_target on s.id = sm_target.session_id
    left join transactions t on s.id = t.session_id
    left join transaction_lines l on t.id = l.transaction_id
    where sm_target.member_id = :member_id and sm_target.is_deleted = false
    group by s.id, s.name, s.is_active, s.created_at, s.updated_at
    order by s.updated_at desc, s.id desc
    limit case when :page_size = 0 then -1 else :page_size end
)
select
    ps.*,
    m.id as member_id,
    m.name as member_name,
    m.username as member_username
from paged_sessions ps
join session_members sm on ps.id = sm.session_id
join members m on sm.member_id = m.id
order by ps.updated_at desc, ps.id desc;

-- name: GetSessionTransactionLines :many
select
    t.id as transaction_id,
    t.session_id,
    t.member_id as creator_id,
    t.created_at,
    l.member_id,
    l.amount
from transactions t
join transaction_lines l on t.id = l.transaction_id
where t.session_id = ?;

-- name: CreateSession :one
insert into sessions (id, name)
values (?, ?)
returning *;

-- name: SetSessionUpdatedAtNow :exec
update sessions
set updated_at = current_timestamp
where id = ?;

-- name: GetMember :one
select * from members where id = ? limit 1;

-- name: GetSessionMember :one
select m.*, sm.is_owner, sm.is_admin
from members m
join session_members sm on m.id = sm.member_id
where m.id = ? and sm.session_id = ?;

-- name: ListSessionMembers :many
select m.*, sm.is_owner, sm.is_admin, sm.is_deleted
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
insert into session_members (session_id, member_id, is_owner, is_admin)
values (?, ?, ?, ?)
on conflict(session_id, member_id) do update
set is_deleted = false, updated_at = current_timestamp;

-- name: DeleteSessionMember :exec
update session_members
set is_deleted = true,
    is_admin = false
where session_id = ? and member_id = ?;

-- name: UpdateSessionMemberAdmin :exec
update session_members
set is_admin = ?,
    updated_at = current_timestamp
where session_id = ?
    and member_id = ?;

-- name: CountSessionMembers :one
select count(*) from session_members
where session_id = ?;

-- name: CountSessionAdmins :one
select count(*) from session_members
where session_id = ? and is_admin = true;

-- name: AddTransaction :one
insert into transactions (id, session_id, member_id, created_at)
values (?, ?, ?, ?)
on conflict do nothing
returning *;

-- name: AddTransactionLine :one
insert into transaction_lines (transaction_id, member_id, amount)
values (?, ?, ?)
on conflict (transaction_id, member_id)
    do update set amount = excluded.amount
returning *;
