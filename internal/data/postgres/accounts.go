package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

const (
	accountsTableName = "accounts"
)

type accountsQ struct {
	*crudQ[*data.Account, uuid.UUID]
}

func NewAccountsQ(db *pgdb.DB) data.Accounts {
	return &accountsQ{
		newCRUDQ[*data.Account, uuid.UUID](db, accountsTableName),
	}
}

func (q *accountsQ) WhereID(id ...uuid.UUID) data.Accounts {
	q.sel = q.sel.Where(sq.Eq{idColumnName: id})
	return q
}
