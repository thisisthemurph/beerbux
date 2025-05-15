-- name: AddTransaction :one
insert into session_transactions (id, session_id, member_id)
values ($1, $2, $3)
on conflict do nothing
returning *;

-- name: AddTransactionLine :one
insert into session_transaction_lines (transaction_id, member_id, amount)
values ($1, $2, $3)
on conflict (transaction_id, member_id)
do update set amount = excluded.amount
returning *;
