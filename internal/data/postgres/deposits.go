package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

const (
	depositsTableName = "deposits"

	recepientColumnName = "recepient_fkey"
)

type depositsQ struct {
	*crudQ[*data.Deposit, uuid.UUID]
}

func NewDepositsQ(db *pgdb.DB) data.Deposits {
	return &depositsQ{
		newCRUDQ[*data.Deposit, uuid.UUID](db, depositsTableName),
	}
}

func (q *depositsQ) WhereRecepient(recepient uuid.UUID) data.Deposits {
	q.sel = q.sel.Where(sq.Eq{recepientColumnName: recepient})
	return q
}
