package postgres

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

const (
	transactionsTableName = "transactions"

	typeColumnName      = "type"
	senderColumnName    = "sender_fkey"
	recipientColumnName = "recipient_fkey"
)

type transactionsQ struct {
	*crudQ[*data.Transaction, uuid.UUID]
}

func NewTransactionsQ(db *pgdb.DB) data.Transactions {
	return &transactionsQ{
		newCRUDQ[*data.Transaction, uuid.UUID](db, transactionsTableName),
	}
}

func (q *transactionsQ) Insert(transaction *data.Transaction) error {
	entry := structs.Map(transaction)

	if transaction.Sender == uuid.Nil {
		entry[senderColumnName] = uuid.NullUUID{}
	} else {
		entry[senderColumnName] = uuid.NullUUID{UUID: transaction.Sender, Valid: true}
	}

	if transaction.Recipient == uuid.Nil {
		entry[recipientColumnName] = uuid.NullUUID{}
	} else {
		entry[recipientColumnName] = uuid.NullUUID{UUID: transaction.Recipient, Valid: true}
	}

	return q.db.Get(transaction.GetID(),
		sq.Insert(transactionsTableName).
			SetMap(entry).
			Suffix(fmt.Sprintf("RETURNING %s ", idColumnName)),
	)
}

func (q *transactionsQ) WhereType(t data.TransactionType) data.Transactions {
	q.sel = q.sel.Where(sq.Eq{typeColumnName: t})
	return q
}

func (q *transactionsQ) WhereSender(sender uuid.UUID) data.Transactions {
	q.sel = q.sel.Where(sq.Eq{senderColumnName: sender})
	return q
}

func (q *transactionsQ) WhereRecipient(recipient uuid.UUID) data.Transactions {
	q.sel = q.sel.Where(sq.Eq{recipientColumnName: recipient})
	return q
}

func (q *transactionsQ) WhereAccount(account uuid.UUID) data.Transactions {
	q.sel = q.sel.Where(sq.Or{
		sq.Eq{senderColumnName: account},
		sq.Eq{recipientColumnName: account},
	})
	return q
}

func (q *transactionsQ) Limit(limit uint64) data.Transactions {
	q.sel = q.sel.Limit(limit)
	return q
}

func (q *transactionsQ) Offset(offset uint64) data.Transactions {
	q.sel = q.sel.Offset(offset)
	return q
}

func (q *transactionsQ) OrderBy(sort ...string) data.Transactions {
	q.sel = q.sel.OrderBy(sort...)

	return q
}
