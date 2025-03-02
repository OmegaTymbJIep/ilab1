package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/omegatymbjiep/ilab1/internal/data"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/requests"
)

var ErrorNonZeroBalance = errors.New("account with non-zero balance")
var ErrorAccountNotFound = errors.New("account not found")

type Accounts struct {
	db data.MainQ
}

func NewMain(db data.MainQ) *Accounts {
	return &Accounts{
		db: db,
	}
}

func (m *Accounts) CreateAccount(customerID uuid.UUID, req *requests.CreateAccount) (*data.Account, error) {
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

func (m *Accounts) GetAccountList(customerID uuid.UUID) ([]*data.Account, error) {
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

func (m *Accounts) GetAccount(customerID, accountID uuid.UUID) (*data.Account, error) {
	ok, err := m.db.CustomersAccounts().HasAccount(customerID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to check account existence: %w", err)
	}
	if !ok {
		return nil, ErrorAccountNotFound
	}

	account := new(data.Account)
	if ok, err := m.db.Accounts().WhereID(accountID).Get(account); err != nil || !ok {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

func (m *Accounts) GetAccountTransactions(customerID, accountID uuid.UUID) ([]*data.Transaction, error) {
	ok, err := m.db.CustomersAccounts().HasAccount(customerID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to check account existence: %w", err)
	}
	if !ok {
		return nil, ErrorAccountNotFound
	}

	transactions, err := m.db.Transactions().WhereAccount(accountID).Select()
	if err != nil {
		return nil, fmt.Errorf("failed to get transers: %w", err)
	}

	return transactions, nil
}

func (m *Accounts) DeleteAccount(customerID, accountID uuid.UUID) error {
	account, err := m.GetAccount(customerID, accountID)
	if err != nil {
		return err
	}

	if account.Balance != 0 {
		return ErrorNonZeroBalance
	}

	err = m.db.Transaction(func() error {
		err = m.db.CustomersAccounts().RemoveAccountsFromCustomer(customerID, accountID)
		if err != nil {
			return fmt.Errorf("failed to remove customer association: %w", err)
		}

		if err = m.db.Accounts().Delete(accountID); err != nil {
			return fmt.Errorf("failed to delete account: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
