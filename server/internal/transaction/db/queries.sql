-- name: GetSessionByID :one
select s.id, s.name, s.is_active, s.created_at, s.updated_at
from sessions s
where s.id = $1;

-- name: GetSessionMemberIDs :many
select member_id
from session_members
where session_id = $1 and is_deleted = false;

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
