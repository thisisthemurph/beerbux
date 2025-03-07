package repository

import (
	"context"
	"database/sql"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository/ledger"
)

// TX is a transaction interface that can be used to begin a transaction.
// This is used to limit the scope of the database/sql package to the server package.
type TX interface {
	// Begin starts a database transaction using context.Background.
	Begin() (*sql.Tx, error)
	// BeginTx starts a database transaction with the provided context and options.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type LedgerQueriesWrapper struct {
	Transaction TX
	*ledger.Queries
}

func NewLedgerQueries(db *sql.DB) *LedgerQueriesWrapper {
	return &LedgerQueriesWrapper{
		Transaction: db,
		Queries:     ledger.New(db),
	}
}
