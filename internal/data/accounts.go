package data

import (
	"time"

	"github.com/google/uuid"
)

type Accounts interface {
	CRUDQ[*Account, uuid.UUID]

	WhereID(id ...uuid.UUID) Accounts
}

type Account struct {
	Entity[uuid.UUID] `structs:"-"`

	Balance   int       `db:"balance"    structs:"balance"`
	UpdatedAt time.Time `db:"updated_at" structs:"-"`
}
