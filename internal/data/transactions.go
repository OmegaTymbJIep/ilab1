package data

import (
	"time"
)

type TransactionType byte

const (
	DepositTransaction    TransactionType = 0x0
	TransferTransaction   TransactionType = 0x1
	WithdrawalTransaction TransactionType = 0x2
)

func (t TransactionType) String() string {
	switch t {
	case DepositTransaction:
		return "deposit"
	case TransferTransaction:
		return "transfer"
	case WithdrawalTransaction:
		return "withdrawal"
	default:
		return "unknown"
	}
}

type Transaction interface {
	GetType() TransactionType
	GetAmount() uint
	GetCreatedAt() time.Time
}
