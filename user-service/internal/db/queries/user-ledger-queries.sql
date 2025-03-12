-- name: InsertUserLedger :exec
insert into user_ledger (user_id, participant_id, amount, type)
values (:user_id, :participant_id, :amount, case when :amount < 0 then 'debit' else 'credit' end);
