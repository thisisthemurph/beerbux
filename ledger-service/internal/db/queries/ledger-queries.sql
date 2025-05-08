-- name: InsertLedger :exec
insert into ledger (id, transaction_id, session_id, user_id, amount)
values (?, ?, ?, ?, ?);

-- name: CalculateUserTotal :one
select
    user_id,
    cast(coalesce(sum(case when amount > 0 then amount else 0 end), 0) as real) as debit,
    cast(coalesce(sum(case when amount < 0 then amount * -1 else 0 end), 0) as real) as credit,
    cast(coalesce(sum(amount), 0) as real) as net
from ledger
where user_id = ?
group by user_id;
