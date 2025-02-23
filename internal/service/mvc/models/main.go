package models

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

type Main struct {
	db data.MainQ
}

func NewMain(db data.MainQ) *Main {
	return &Main{
		db: db,
	}
}

func (m *Main) GetAccounts(customerID uuid.UUID) ([]*data.Account, error) {
	accountIDs, err := m.db.CustomersAccounts().GetCustomersByAccount(customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts by customer: %w", err)
	}

	accounts, err := m.db.Accounts().WhereID(accountIDs...).Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}
