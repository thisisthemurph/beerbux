-- name: InsertLedger :exec
insert into ledger (id, transaction_id, session_id, user_id, amount)
values (?, ?, ?, ?, ?);
