-- name: CreateSession :one
insert into sessions (name, creator_id) values ($1, $2) returning *;

-- name: SessionExists :one
select exists(select 1 from sessions where id = $1);

-- name: GetSessionDetailsByID :one
-- GetSessionDetailsByID returns the basic session data for the given session ID.
select s.id, s.name, s.is_active, s.created_at, s.updated_at
from sessions s
where s.id = $1;

-- name: GetSessionByID :one
select s.id, s.name, s.is_active, s.created_at, s.updated_at, coalesce(sum(l.amount), 0)::float8 as total
from sessions s
left join session_transactions t on s.id = t.session_id
left join session_transaction_lines l on t.id = l.transaction_id
where s.id = $1
group by s.id, s.name, s.is_active, s.created_at, s.updated_at;

-- name: ListSessionsForUser :many
with user_sessions as (
    select
        s.id,
        s.name,
        s.is_active,
        s.created_at,
        s.updated_at,
        coalesce(sum(tl.amount), 0)::float8 as total_amount
    from sessions s
        join session_members sm_target on s.id = sm_target.session_id
        left join session_transactions t on s.id = t.session_id
        left join session_transaction_lines tl on t.id = tl.transaction_id
    where sm_target.member_id = $1
      and sm_target.is_deleted = false
    group by s.id, s.name, s.is_active, s.created_at, s.updated_at
    order by s.updated_at desc, s.id desc
    limit nullif($2, 0)
)
select
    st.*,
    m.id as member_id,
    m.name as member_name,
    m.username as member_username
from user_sessions st
join session_members sm on st.id = sm.session_id
join users m on sm.member_id = m.id
order by st.updated_at desc, st.id desc;

-- name: GetSessionTransactionLines :many
select
    t.id as transaction_id,
    t.session_id,
    t.member_id as creator_id,
    t.created_at,
    tl.member_id,
    tl.amount::float8 as amount
from session_transactions t
join session_transaction_lines tl on t.id = tl.transaction_id
where t.session_id = $1;

-- name: UpdateSessionActiveState :exec
update sessions set is_active = $2 where id = $1;

-- name: SetSessionUpdatedNow :exec
update sessions set updated_at = now() where id = $1;

-- name: GetSessionMember :one
select u.id, u.username, u.name, u.created_at, u.updated_at, sm.is_admin
from users u
join session_members sm on u.id = sm.member_id
where sm.session_id = $1 and u.id = $2;

-- name: ListSessionMembers :many
select u.id, u.username, u.name, u.created_at, u.updated_at, sm.is_admin, sm.is_deleted
from users u
join session_members sm on u.id = sm.member_id
where sm.session_id = $1;

-- name: ListSessionMemberIDs :many
select member_id
from session_members
where session_id = $1 and is_deleted = false;

-- name: AddMemberToSession :exec
insert into session_members (session_id, member_id, is_admin)
values ($1, $2, $3)
on conflict(session_id, member_id) do update
set is_deleted = false;

-- name: DeleteSessionMember :exec
update session_members
set is_deleted = true, is_admin = false
where session_id = $1 and member_id = $2;

-- name: UpdateSessionMemberAdminState :exec
update session_members
set is_admin = $3
where session_id = $1 and member_id = $2;

-- name: CountSessionMembers :one
select count(*) from session_members where session_id = $1 and is_deleted = false;

-- name: CountSessionMembersIncludingDeleted :one
select count(*) from session_members where session_id = $1;

-- name: CountSessionAdminMembers :one
select count(*) from session_members where session_id = $1 and is_admin = true;

-- name: GetSessionHistory :many
select * from session_history where session_id = $1;

-- name: CreateSessionHistory :exec
insert into session_history (session_id, member_id, event_type, event_data)
values ($1, $2, $3, $4);

-- name: CreateTransaction :one
insert into session_transactions (session_id, member_id)
values ($1, $2)
on conflict do nothing
returning *;

-- name: CreateTransactionLine :one
insert into session_transaction_lines (transaction_id, member_id, amount)
values ($1, $2, $3)
on conflict (transaction_id, member_id)
    do update set amount = excluded.amount
returning *;

-- name: CreateLedgerEntry :exec
insert into ledger (transaction_id, user_id, amount) values ($1, $2, $3);