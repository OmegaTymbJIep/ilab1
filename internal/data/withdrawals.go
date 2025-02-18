package data

import (
	"time"

	"github.com/google/uuid"
)

type Withdrawals interface {
	CRUDQ[*Withdrawal, uuid.UUID]

	WhereSender(sender uuid.UUID) Withdrawals
}

type Withdrawal struct {
	Entity[uuid.UUID] `structs:"-"`

	Amount uint      `db:"amount"      structs:"amount"`
	Sender uuid.UUID `db:"sender_fkey" structs:"sender_fkey"`
}

func (w *Withdrawal) GetType() TransactionType {
	return WithdrawalTransaction
}

func (w *Withdrawal) GetAmount() uint {
	return w.Amount
}

func (w *Withdrawal) GetCreatedAt() time.Time {
	return w.CreatedAt
}
