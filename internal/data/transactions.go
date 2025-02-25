package data

import (
	"github.com/google/uuid"
)

type TransactionType int

const (
	DepositTransaction TransactionType = iota
	WithdrawalTransaction
	TransferTransaction
)

func (t TransactionType) String() string {
	switch t {
	case DepositTransaction:
		return "deposit"
	case WithdrawalTransaction:
		return "withdrawal"
	case TransferTransaction:
		return "transfer"
	default:
		return "unknown"
	}
}

type Transactions interface {
	CRUDQ[*Transaction, uuid.UUID]

	WhereType(TransactionType) Transactions
	WhereSender(sender uuid.UUID) Transactions
	WhereRecipient(recipient uuid.UUID) Transactions
	WhereAccount(account uuid.UUID) Transactions

	Limit(limit uint64) Transactions
	Offset(offset uint64) Transactions
	OrderBy(orderBy ...string) Transactions
}

type Transaction struct {
	Entity[uuid.UUID] `structs:"-"`

	Type         TransactionType `db:"type"           structs:"type"`
	Amount       uint            `db:"amount"         structs:"amount"`
	Sender       uuid.UUID       `db:"sender_fkey"    structs:"sender_fkey"`
	Recipient    uuid.UUID       `db:"recipient_fkey" structs:"recipient_fkey"`
	ATMSignature string          `db:"atm_signature"  structs:"atm_signature"`
}
