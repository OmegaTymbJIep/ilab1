package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

const (
	customersTableName = "customers"

	emailColumnName    = "email"
	usernameColumnName = "username"
)

type customersQ struct {
	*crudQ[*data.Customer, uuid.UUID]
}

func NewCustomersQ(db *pgdb.DB) data.Customers {
	return &customersQ{
		newCRUDQ[*data.Customer, uuid.UUID](db, customersTableName),
	}
}

func (q *customersQ) WhereID(id ...uuid.UUID) data.Customers {
	q.sel = q.sel.Where(sq.Eq{idColumnName: id})
	return q
}

func (q *customersQ) WhereEmail(email string) data.Customers {
	q.sel = q.sel.Where(sq.Eq{emailColumnName: email})
	return q
}

func (q *customersQ) WhereUsername(username string) data.Customers {
	q.sel = q.sel.Where(sq.Eq{usernameColumnName: username})
	return q
}

// IsUnique checks that it doesn't exist a customer with the same email or username.
func (q *customersQ) IsUnique(email, username string) (bool, error) {
	var count int

	err := q.db.Get(&count, sq.Select("COUNT(*)").
		From(customersTableName).
		Where(sq.Or{
			sq.Eq{emailColumnName: email},
			sq.Eq{usernameColumnName: username},
		}),
	)

	return count == 0, err
}
