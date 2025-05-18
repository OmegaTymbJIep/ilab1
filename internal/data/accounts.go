package data

import (
	"time"

	"github.com/google/uuid"
)

type Accounts interface {
	CRUDQ[*Account, uuid.UUID]

	IsDeleted(bool) Accounts
	WhereID(id ...uuid.UUID) Accounts
	// LDelete - Logical Delete - marks the account as deleted.
	LDelete(id uuid.UUID) error
}

type Account struct {
	Entity[uuid.UUID] `structs:"-"`

	Name      string    `db:"name"       structs:"name"`
	Balance   int       `db:"balance"    structs:"balance"`
	IsDeleted bool      `db:"is_deleted" structs:"is_deleted"`
	UpdatedAt time.Time `db:"updated_at" structs:"-"`
}
