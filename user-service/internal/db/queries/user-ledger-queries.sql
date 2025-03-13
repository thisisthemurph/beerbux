-- name: InsertUserLedger :exec
insert into user_ledger (user_id, participant_id, amount, type)
values (:user_id, :participant_id, :amount, case when :amount < 0 then 'debit' else 'credit' end);

-- name: CalculateUserNetBalance :one
select cast(coalesce(sum(amount), 0) as real) as net_balance
from user_ledger
where user_id = :user_id
or participant_id = :user_id;
