package data

import (
	"time"

	"github.com/google/uuid"
)

type Transfers interface {
	CRUDQ[*Transfer, uuid.UUID]

	WhereSender(sender uuid.UUID) Transfers
	WhereRecepient(recepient uuid.UUID) Transfers
}

type Transfer struct {
	Entity[uuid.UUID] `structs:"-"`

	Amount    uint      `db:"amount"         structs:"amount"`
	Sender    uuid.UUID `db:"sender_fkey"    structs:"sender_fkey"`
	Recepient uuid.UUID `db:"recepient_fkey" structs:"recepient_fkey"`
}

func (t *Transfer) GetType() TransactionType {
	return TransferTransaction
}

func (t *Transfer) GetAmount() uint {
	return t.Amount
}

func (t *Transfer) GetCreatedAt() time.Time {
	return t.CreatedAt
}
