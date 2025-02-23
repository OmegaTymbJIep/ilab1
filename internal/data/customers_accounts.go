package data

import (
	"github.com/google/uuid"
)

type CustomersAccounts interface {
	AddCustomersToAccount(accountID uuid.UUID, customersID ...uuid.UUID) error
	AddAccountsToCustomer(customerID uuid.UUID, accountsID ...uuid.UUID) error

	RemoveCustomersFromAccount(accountID uuid.UUID, customersID ...uuid.UUID) error
	RemoveAccountsFromCustomer(customerID uuid.UUID, accountsID ...uuid.UUID) error

	GetCustomersByAccount(accountID uuid.UUID) ([]uuid.UUID, error)
	GetAccountsByCustomer(customerID uuid.UUID) ([]uuid.UUID, error)

	HasAccount(accountID uuid.UUID, customerID uuid.UUID) (bool, error)
}
