package data

import (
	"database/sql"
	"time"
)

type MainQ interface {
	Customers() Customers
	Accounts() Accounts
	CustomersAccounts() CustomersAccounts
	Transactions() Transactions
	AuditLogs() AuditLogs

	Transaction(func() error) error
	IsolatedTransaction(sql.IsolationLevel, func() error) error
}

type CRUDQ[T IEntity[ID], ID comparable] interface {
	Insert(entity T) error
	Get(T) (bool, error)
	Select() ([]T, error)
	Update(entity T) error
	Delete(id ID) error
	Count() (uint64, error)
}

type IEntity[ID comparable] interface {
	GetID() *ID
	GetCreatedAt() time.Time
}

type Entity[I comparable] struct {
	ID        I         `db:"id"         structs:"id"`
	CreatedAt time.Time `db:"created_at" structs:"created_at"`
}

func (e *Entity[I]) GetID() *I {
	return &e.ID
}

func (e *Entity[I]) GetCreatedAt() time.Time {
	return e.CreatedAt
}
