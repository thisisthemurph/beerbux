-- name: Create :one
insert into sessions (name, creator_id) values ($1, $2) returning *;

-- name: Exists :one
select exists(select 1 from sessions where id = $1);

-- name: Get :one
select s.id, s.name, s.is_active, s.created_at, s.updated_at, coalesce(sum(l.amount), 0)::float8 as total
from sessions s
left join session_transactions t on s.id = t.session_id
left join session_transaction_lines l on t.id = l.transaction_id
where s.id = $1
group by s.id, s.name, s.is_active, s.created_at, s.updated_at;

-- name: ListForUser :many
with paged_sessions as (
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
from paged_sessions st
join session_members sm on st.id = sm.session_id
join users m on sm.member_id = m.id
order by st.updated_at desc, st.id desc;

-- name: GetTransactionLines :many
select
    t.id as transaction_id,
    t.session_id,
    t.member_id as creator_id,
    t.created_at,
    tl.member_id,
    tl.amount
from session_transactions t
join session_transaction_lines tl on t.id = tl.transaction_id
where t.session_id = $1;

-- name: UpdateActiveState :exec
update sessions set is_active = $2 where id = $1;

-- name: SetUpdatedNow :exec
update sessions set updated_at = now() where id = $1;

-- name: GetMember :one
select u.id, u.username, u.name, u.created_at, u.updated_at, sm.is_admin
from users u
join session_members sm on u.id = sm.member_id
where sm.session_id = $1 and u.id = $2;

-- name: ListMembers :many
select u.id, u.username, u.name, u.created_at, u.updated_at, sm.is_admin
from users u
join session_members sm on u.id = sm.member_id
where sm.session_id = $1;

-- name: AddMember :exec
insert into session_members (session_id, member_id, is_admin)
values ($1, $2, $3)
on conflict(session_id, member_id) do nothing;

-- name: DeleteMember :exec
update session_members
set is_deleted = true, is_admin = false
where session_id = $1 and member_id = $2;

-- name: UpdateMemberAdminState :exec
update session_members
set is_admin = $3
where session_id = $1 and member_id = $2;

-- name: CountMembers :one
select count(*) from session_members where session_id = $1 and is_deleted = false;

-- name: CountMembersIncludingDeleted :one
select count(*) from session_members where session_id = $1;

-- name: CountAdminMembers :one
select count(*) from session_members where session_id = $1 and is_admin = true;