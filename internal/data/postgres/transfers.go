package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

const (
	transfersTableName = "transfers"
)

type transfersQ struct {
	*crudQ[*data.Transfer, uuid.UUID]
}

func NewTransfersQ(db *pgdb.DB) data.Transfers {
	return &transfersQ{
		newCRUDQ[*data.Transfer, uuid.UUID](db, transfersTableName),
	}
}

func (q *transfersQ) WhereSender(sender uuid.UUID) data.Transfers {
	q.sel = q.sel.Where(sq.Eq{senderColumnName: sender})
	return q
}

func (q *transfersQ) WhereRecepient(recepient uuid.UUID) data.Transfers {
	q.sel = q.sel.Where(sq.Eq{recepientColumnName: recepient})
	return q
}
