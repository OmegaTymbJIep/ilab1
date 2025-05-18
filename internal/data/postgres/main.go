package postgres

import (
	"database/sql"

	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

type mainQ struct {
	db *pgdb.DB
}

func NewMainQ(db *pgdb.DB) data.MainQ {
	return &mainQ{
		db: db,
	}
}

func (q *mainQ) Customers() data.Customers {
	return NewCustomersQ(q.db)
}

func (q *mainQ) Accounts() data.Accounts {
	return NewAccountsQ(q.db)
}

func (q *mainQ) CustomersAccounts() data.CustomersAccounts {
	return NewCustomersAccountsQ(q.db)
}

func (q *mainQ) Transactions() data.Transactions {
	return NewTransactionsQ(q.db)
}

func (q *mainQ) Transaction(fn func() error) error {
	return q.db.Transaction(fn)
}

func (q *mainQ) AuditLogs() data.AuditLogs {
	return NewAuditLogsQ(q.db)
}

func (q *mainQ) IsolatedTransaction(isolationLevel sql.IsolationLevel, fn func() error) error {
	return q.db.TransactionWithOptions(&sql.TxOptions{Isolation: isolationLevel}, fn)
}
