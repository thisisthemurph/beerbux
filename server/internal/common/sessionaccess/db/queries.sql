-- name: GetSessionByID :one
select s.id, s.name, s.is_active, s.created_at, s.updated_at, coalesce(sum(l.amount), 0)::float8 as total
from sessions s
         left join session_transactions t on s.id = t.session_id
         left join session_transaction_lines l on t.id = l.transaction_id
where s.id = sqlc.arg(session_id)::uuid
group by s.id, s.name, s.is_active, s.created_at, s.updated_at;

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

-- name: ListSessionMembers :many
select u.id, u.username, u.name, u.created_at, u.updated_at, sm.is_admin, sm.is_deleted
from users u
join session_members sm on u.id = sm.member_id
where sm.session_id = $1;
