package data

import (
	"time"

	"github.com/google/uuid"
)

type Deposits interface {
	CRUDQ[*Deposit, uuid.UUID]

	WhereRecepient(recepient uuid.UUID) Deposits
}

type Deposit struct {
	Entity[uuid.UUID] `structs:"-"`

	Amount       uint      `db:"amount"         structs:"amount"`
	Recepient    uuid.UUID `db:"recepient_fkey" structs:"recepient_fkey"`
	ATMSignature string    `db:"atm_signature"  structs:"atm_signature"`
}

func (d *Deposit) GetType() TransactionType {
	return DepositTransaction
}

func (d *Deposit) GetAmount() uint {
	return d.Amount
}

func (d *Deposit) GetCreatedAt() time.Time {
	return d.CreatedAt
}
