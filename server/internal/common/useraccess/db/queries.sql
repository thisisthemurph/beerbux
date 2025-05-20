-- name: GetUserByID :one
select
    u.id, u.username, u.name, u.created_at, u.updated_at,
    coalesce(t.debit, 0) as debit,
    coalesce(t.credit, 0) as credit
from users u
left join user_totals t on u.id = t.user_id
where u.id = $1
limit 1;

-- name: GetByUsername :one
select
    u.id, u.username, u.name, u.created_at, u.updated_at,
    coalesce(t.debit, 0) as debit,
    coalesce(t.credit, 0) as credit
from users u
         left join user_totals t on u.id = t.user_id
where u.username = $1
limit 1;
