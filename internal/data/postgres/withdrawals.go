package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

const (
	withdrawalsTableName = "withdrawals"

	senderColumnName = "sender_fkey"
)

type withdrawalsQ struct {
	*crudQ[*data.Withdrawal, uuid.UUID]
}

func NewWithdrawalsQ(db *pgdb.DB) data.Withdrawals {
	return &withdrawalsQ{
		newCRUDQ[*data.Withdrawal, uuid.UUID](db, withdrawalsTableName),
	}
}

func (q *withdrawalsQ) WhereSender(sender uuid.UUID) data.Withdrawals {
	q.sel = q.sel.Where(sq.Eq{senderColumnName: sender})
	return q
}
