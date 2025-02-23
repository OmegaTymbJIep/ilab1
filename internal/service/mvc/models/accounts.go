package models

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/omegatymbjiep/ilab1/internal/data"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/requests"
)

type Main struct {
	db data.MainQ
}

func NewMain(db data.MainQ) *Main {
	return &Main{
		db: db,
	}
}

func (m *Main) CreateAccount(customerID uuid.UUID, req *requests.CreateAccount) (*data.Account, error) {
	account := &data.Account{
		Name:    req.Name,
		Balance: 0,
	}

	err := m.db.Transaction(func() error {
		if err := m.db.Accounts().Insert(account); err != nil {
			return fmt.Errorf("failed to insert account: %w", err)
		}

		err := m.db.CustomersAccounts().AddAccountsToCustomer(customerID, account.ID)
		if err != nil {
			return fmt.Errorf("failed to add account to customer: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if ok, err := m.db.Accounts().WhereID(account.ID).Get(account); err != nil || !ok {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

func (m *Main) GetAccounts(customerID uuid.UUID) ([]*data.Account, error) {
	accountIDs, err := m.db.CustomersAccounts().GetAccountsByCustomer(customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	accounts, err := m.db.Accounts().WhereID(accountIDs...).Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}
