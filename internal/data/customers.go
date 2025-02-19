package data

import (
	"time"

	"github.com/google/uuid"
)

type Customers interface {
	CRUDQ[*Customer, uuid.UUID]

	WhereID(id ...uuid.UUID) Customers
	WhereEmail(email string) Customers
	WhereUsername(username string) Customers
}

type Customer struct {
	Entity[uuid.UUID] `structs:"-"`

	Email        string    `db:"email"         structs:"email"`
	Username     string    `db:"username"      structs:"username"`
	PasswordHash string    `db:"password_hash" structs:"password_hash"`
	FirstName    *string   `db:"first_name"    structs:"first_name"`
	LastName     *string   `db:"last_name"     structs:"last_name"`
	UpdatedAt    time.Time `db:"updated_at"    structs:"-"`

	account []Accounts
}

func (c *Customer) GetID() *uuid.UUID {
	return &c.ID
}
