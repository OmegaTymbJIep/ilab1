package postgres

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

const (
	customersAccountsTableName = "customers_accounts"

	accountFkeyColumnName  = "account_fkey"
	customerFkeyColumnName = "customer_fkey"
)

type customersAccountsQ struct {
	db *pgdb.DB
}

func NewCustomersAccountsQ(db *pgdb.DB) data.CustomersAccounts {
	return &customersAccountsQ{
		db: db,
	}
}

func (q *customersAccountsQ) addCustomerToAccount(accountID, customerID uuid.UUID) error {
	return q.db.Exec(
		sq.Insert(customersAccountsTableName).
			Columns(accountFkeyColumnName, customerFkeyColumnName).
			Values(accountID, customerID),
	)
}

func (q *customersAccountsQ) AddCustomersToAccount(accountID uuid.UUID, customersID ...uuid.UUID) error {
	if len(customersID) == 0 {
		return nil
	}

	if len(customersID) == 1 {
		return q.addCustomerToAccount(accountID, customersID[0])
	}

	db := q.db.Clone()
	return db.Transaction(func() error {
		for _, customerID := range customersID {
			if err := db.Exec(
				sq.Insert(customersAccountsTableName).
					Columns(accountFkeyColumnName, customerFkeyColumnName).
					Values(accountID, customerID),
			); err != nil {
				return err
			}
		}

		return nil
	})
}

func (q *customersAccountsQ) AddAccountsToCustomer(customerID uuid.UUID, accountsID ...uuid.UUID) error {
	if len(accountsID) == 0 {
		return nil
	}

	if len(accountsID) == 1 {
		return q.addCustomerToAccount(accountsID[0], customerID)
	}

	db := q.db.Clone()
	return db.Transaction(func() error {
		for _, accountID := range accountsID {
			if err := db.Exec(
				sq.Insert(customersAccountsTableName).
					Columns(accountFkeyColumnName, customerFkeyColumnName).
					Values(accountID, customerID),
			); err != nil {
				return err
			}
		}

		return nil
	})
}

func (q *customersAccountsQ) removeCustomerFromAccount(accountID, customerID uuid.UUID) error {
	if err := q.db.Exec(
		sq.Delete(customersAccountsTableName).
			Where(sq.Eq{
				accountFkeyColumnName:  accountID,
				customerFkeyColumnName: customerID,
			}),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	}

	return nil
}

func (q *customersAccountsQ) RemoveCustomersFromAccount(accountID uuid.UUID, customersID ...uuid.UUID) error {
	if len(customersID) == 0 {
		return nil
	}

	if len(customersID) == 1 {
		return q.removeCustomerFromAccount(accountID, customersID[0])
	}

	db := q.db.Clone()
	return db.Transaction(func() error {
		for _, customerID := range customersID {
			if err := db.Exec(
				sq.Delete(customersAccountsTableName).
					Where(sq.Eq{
						accountFkeyColumnName:  accountID,
						customerFkeyColumnName: customerID,
					}),
			); err != nil {
				return err
			}
		}

		return nil
	})
}

func (q *customersAccountsQ) RemoveAccountsFromCustomer(customerID uuid.UUID, accountsID ...uuid.UUID) error {
	if len(accountsID) == 0 {
		return nil
	}

	if len(accountsID) == 1 {
		return q.removeCustomerFromAccount(accountsID[0], customerID)
	}

	db := q.db.Clone()
	return db.Transaction(func() error {
		for _, accountID := range accountsID {
			if err := db.Exec(
				sq.Delete(customersAccountsTableName).
					Where(sq.Eq{
						accountFkeyColumnName:  accountID,
						customerFkeyColumnName: customerID,
					}),
			); err != nil {
				return err
			}
		}

		return nil
	})
}

func (q *customersAccountsQ) GetCustomersByAccount(accountID uuid.UUID) ([]uuid.UUID, error) {
	var result []uuid.UUID

	if err := q.db.Select(&result,
		sq.Select(customerFkeyColumnName).
			From(customersAccountsTableName).
			Where(sq.Eq{accountFkeyColumnName: accountID}),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, nil
		}

		return nil, err
	}

	return result, nil
}

func (q *customersAccountsQ) GetAccountsByCustomer(customerID uuid.UUID) ([]uuid.UUID, error) {
	var result []uuid.UUID

	if err := q.db.Select(&result,
		sq.Select(accountFkeyColumnName).
			From(customersAccountsTableName).
			Where(sq.Eq{customerFkeyColumnName: customerID}),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, nil
		}

		return nil, err
	}

	return result, nil
}
