// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: ledger-queries.sql

package ledger

import (
	"context"
)

const calculateUserTotal = `-- name: CalculateUserTotal :one
select
    user_id,
    cast(coalesce(sum(case when amount > 0 then amount else 0 end), 0) as real) as credit,
    cast(coalesce(sum(case when amount < 0 then amount * -1 else 0 end), 0) as real) as debit,
    cast(coalesce(sum(amount), 0) as real) as net
from ledger
where user_id = ?
group by user_id
`

type CalculateUserTotalRow struct {
	UserID string
	Credit float64
	Debit  float64
	Net    float64
}

func (q *Queries) CalculateUserTotal(ctx context.Context, userID string) (CalculateUserTotalRow, error) {
	row := q.db.QueryRowContext(ctx, calculateUserTotal, userID)
	var i CalculateUserTotalRow
	err := row.Scan(
		&i.UserID,
		&i.Credit,
		&i.Debit,
		&i.Net,
	)
	return i, err
}

const insertLedger = `-- name: InsertLedger :exec
insert into ledger (id, transaction_id, session_id, user_id, amount)
values (?, ?, ?, ?, ?)
`

type InsertLedgerParams struct {
	ID            string
	TransactionID string
	SessionID     string
	UserID        string
	Amount        float64
}

func (q *Queries) InsertLedger(ctx context.Context, arg InsertLedgerParams) error {
	_, err := q.db.ExecContext(ctx, insertLedger,
		arg.ID,
		arg.TransactionID,
		arg.SessionID,
		arg.UserID,
		arg.Amount,
	)
	return err
}
