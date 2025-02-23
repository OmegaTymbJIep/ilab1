package postgres

import (
	"testing"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/assets"
	"github.com/omegatymbjiep/ilab1/internal/data"
)

var _ = assets.Migrations

const (
	testDBURL = "postgres://ilab1:ilab1@localhost:5432/ilab1test?sslmode=disable"
)

func newTestMainQ(t *testing.T) data.MainQ {
	db, err := pgdb.Open(pgdb.Opts{
		URL: testDBURL,
	})
	assert.NoError(t, err, "failed to open test db (%s)", testDBURL)

	var migrations = &migrate.EmbedFileSystemMigrationSource{
		FileSystem: assets.Migrations,
		Root:       "migrations",
	}

	_, err = migrate.Exec(db.RawDB(), "postgres", migrations, migrate.Down)
	assert.NoError(t, err, "failed to apply migrations down")

	appliedN, err := migrate.Exec(db.RawDB(), "postgres", migrations, migrate.Up)
	assert.NoError(t, err, "failed to apply migrations up")
	assert.Greater(t, appliedN, 0, "no migrations applied")

	return NewMainQ(db)
}

func TestAccountsCRUD(t *testing.T) {
	db := newTestMainQ(t)
	accounts := db.Accounts()

	// Create an account
	account := &data.Account{
		Name:    "account1",
		Balance: 1000,
	}

	err := accounts.Insert(account)
	require.NoError(t, err)

	// Retrieve the account
	fetched := new(data.Account)
	ok, err := accounts.WhereID(account.ID).Get(fetched)
	require.NoError(t, err)
	require.True(t, ok, "inserted account not found")
	assert.Equal(t, account.Balance, fetched.Balance)

	// Update the account balance
	account.Balance = 2000
	err = accounts.Update(account)
	require.NoError(t, err)

	// Verify the update
	updated := new(data.Account)
	ok, err = accounts.WhereID(account.ID).Get(updated)
	require.NoError(t, err)
	require.True(t, ok, "updated account not found")
	assert.Equal(t, 2000, updated.Balance)

	// Delete the account
	err = accounts.Delete(account.ID)
	require.NoError(t, err)

	// Verify deletion
	deleted := new(data.Account)
	ok, err = accounts.WhereID(account.ID).Get(deleted)
	require.NoError(t, err)
	assert.False(t, ok, "deleted account still exists")
}

func TestCustomersCRUD(t *testing.T) {
	db := newTestMainQ(t)
	customers := db.Customers()

	// Create a customer
	customer := &data.Customer{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashed_password",
	}
	err := customers.Insert(customer)
	require.NoError(t, err)

	// Retrieve the customer by ID
	fetched := new(data.Customer)
	ok, err := customers.WhereID(customer.ID).Get(fetched)
	require.NoError(t, err)
	require.True(t, ok, "inserted customer not found")
	assert.Equal(t, customer.Email, fetched.Email)

	// Update the customer email
	customer.Email = "updated@example.com"
	err = customers.Update(customer)
	require.NoError(t, err)

	// Verify the update
	updated := new(data.Customer)
	ok, err = customers.WhereID(customer.ID).Get(updated)
	require.NoError(t, err)
	require.True(t, ok, "updated customer not found")
	assert.Equal(t, "updated@example.com", updated.Email)

	// Delete the customer
	err = customers.Delete(customer.ID)
	require.NoError(t, err)

	// Verify deletion
	deleted := new(data.Customer)
	ok, err = customers.WhereID(customer.ID).Get(deleted)
	require.NoError(t, err)
	assert.False(t, ok, "deleted customer still exists")
}

func TestWithdrawalsCRUD(t *testing.T) {
	db := newTestMainQ(t)
	withdrawals := db.Withdrawals()
	accounts := db.Accounts()

	// Create an account
	account := &data.Account{
		Name:    "account1",
		Balance: 1000,
	}

	err := accounts.Insert(account)
	require.NoError(t, err)

	// Create a withdrawal
	withdrawal := &data.Withdrawal{
		Sender: account.ID,
		Amount: 500,
	}
	err = withdrawals.Insert(withdrawal)
	require.NoError(t, err)

	// Retrieve the withdrawal by ID
	fetched := new(data.Withdrawal)
	ok, err := withdrawals.WhereSender(withdrawal.Sender).Get(fetched)
	require.NoError(t, err)
	require.True(t, ok, "inserted withdrawal not found")
	assert.Equal(t, withdrawal.Amount, fetched.Amount)

	// Update the withdrawal amount
	withdrawal.Amount = 1000
	err = withdrawals.Update(withdrawal)
	require.NoError(t, err)

	// Verify the update
	updated := new(data.Withdrawal)
	ok, err = withdrawals.WhereSender(withdrawal.Sender).Get(updated)
	require.NoError(t, err)
	require.True(t, ok, "updated withdrawal not found")
	assert.Equal(t, uint(1000), updated.Amount)

	// Delete the withdrawal
	err = withdrawals.Delete(withdrawal.ID)
	require.NoError(t, err)

	// Verify deletion
	deleted := new(data.Withdrawal)
	ok, err = withdrawals.WhereSender(withdrawal.Sender).Get(deleted)
	require.NoError(t, err)
	assert.False(t, ok, "deleted withdrawal still exists")
}

func TestDepositsCRUD(t *testing.T) {
	db := newTestMainQ(t)
	deposits := db.Deposits()
	accounts := db.Accounts()

	// Create an account for deposit
	account := &data.Account{
		Name:    "account1",
		Balance: 0,
	}
	err := accounts.Insert(account)
	require.NoError(t, err)

	// Create a deposit
	deposit := &data.Deposit{
		Recepient:    account.ID,
		ATMSignature: "test_signature",
		Amount:       500,
	}
	err = deposits.Insert(deposit)
	require.NoError(t, err)

	// Retrieve the deposit
	fetched := new(data.Deposit)
	ok, err := deposits.WhereRecepient(deposit.Recepient).Get(fetched)
	require.NoError(t, err)
	require.True(t, ok, "inserted deposit not found")
	assert.Equal(t, deposit.Amount, fetched.Amount)

	// Delete the deposit
	err = deposits.Delete(deposit.ID)
	require.NoError(t, err)

	// Verify deletion
	deleted := new(data.Deposit)
	ok, err = deposits.WhereRecepient(deposit.Recepient).Get(deleted)
	require.NoError(t, err)
	assert.False(t, ok, "deleted deposit still exists")
}

func TestTransfersCRUD(t *testing.T) {
	db := newTestMainQ(t)
	transfers := db.Transfers()
	accounts := db.Accounts()

	// Create sender and recipient accounts
	sender := &data.Account{
		Name:    "sender",
		Balance: 1000,
	}
	recipient := &data.Account{
		Name:    "recipient",
		Balance: 500,
	}

	err := accounts.Insert(sender)
	require.NoError(t, err)
	err = accounts.Insert(recipient)
	require.NoError(t, err)

	// Create a transfer
	transfer := &data.Transfer{
		Sender:    sender.ID,
		Recepient: recipient.ID,
		Amount:    200,
	}
	err = transfers.Insert(transfer)
	require.NoError(t, err)

	// Retrieve the transfer
	fetched := new(data.Transfer)
	ok, err := transfers.WhereSender(transfer.Sender).Get(fetched)
	require.NoError(t, err)
	require.True(t, ok, "inserted transfer not found")
	assert.Equal(t, transfer.Amount, fetched.Amount)

	// Delete the transfer
	err = transfers.Delete(transfer.ID)
	require.NoError(t, err)

	// Verify deletion
	deleted := new(data.Transfer)
	ok, err = transfers.WhereSender(transfer.Sender).Get(deleted)
	require.NoError(t, err)
	assert.False(t, ok, "deleted transfer still exists")
}

func TestCustomersAccountsCRUD(t *testing.T) {
	db := newTestMainQ(t)
	customers := db.Customers()
	accounts := db.Accounts()
	customersAccounts := db.CustomersAccounts()

	// Create a customer
	customer := &data.Customer{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashed_password",
	}
	err := customers.Insert(customer)
	require.NoError(t, err)

	// Create an account
	account := &data.Account{Name: "account1", Balance: 1000}
	err = accounts.Insert(account)
	require.NoError(t, err)

	// Link customer to account
	err = customersAccounts.AddCustomersToAccount(account.ID, customer.ID)
	require.NoError(t, err)

	// Retrieve linked customers by account
	retrievedCustomers, err := customersAccounts.GetCustomersByAccount(account.ID)
	require.NoError(t, err)
	assert.Contains(t, retrievedCustomers, customer.ID)

	// Retrieve linked accounts by customer
	retrievedAccounts, err := customersAccounts.GetAccountsByCustomer(customer.ID)
	require.NoError(t, err)
	assert.Contains(t, retrievedAccounts, account.ID)

	// Remove customer from account
	err = customersAccounts.RemoveCustomersFromAccount(account.ID, customer.ID)
	require.NoError(t, err)

	// Verify customer removal
	retrievedCustomers, err = customersAccounts.GetCustomersByAccount(account.ID)
	require.NoError(t, err)
	assert.NotContains(t, retrievedCustomers, customer.ID)
}

func TestCustomerDeletionRestriction(t *testing.T) {
	db := newTestMainQ(t)
	customers := db.Customers()
	accounts := db.Accounts()
	customersAccounts := db.CustomersAccounts()

	// Create a customer
	customer := &data.Customer{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: "hashed_password",
	}
	err := customers.Insert(customer)
	require.NoError(t, err)

	// Create an account
	account := &data.Account{Name: "account1", Balance: 1000}
	err = accounts.Insert(account)
	require.NoError(t, err)

	// Link customer to account
	err = customersAccounts.AddCustomersToAccount(account.ID, customer.ID)
	require.NoError(t, err)

	// Attempt to delete customer (should fail due to ON DELETE RESTRICT)
	err = customers.Delete(customer.ID)
	require.Error(t, err, "customer deletion should be restricted")
}

func TestTransactionCommit(t *testing.T) {
	db := newTestMainQ(t)
	accounts := db.Accounts()

	// Create an account inside a transaction
	err := db.Transaction(func() error {
		account := &data.Account{Name: "account1", Balance: 500}
		err := accounts.Insert(account)
		if err != nil {
			return err
		}
		return nil
	})

	require.NoError(t, err, "transaction should commit successfully")

	// Verify the account was created
	allAccounts, err := accounts.Select()
	require.NoError(t, err)
	assert.NotEmpty(t, allAccounts, "account should exist after commit")
}

func TestTransactionRollback(t *testing.T) {
	db := newTestMainQ(t)
	accounts := db.Accounts()

	initialAccounts, err := accounts.Select()
	require.NoError(t, err)

	err = db.Transaction(func() error {
		account := &data.Account{Name: "account1", Balance: 1000}
		err = accounts.Insert(account)
		if err != nil {
			return err
		}
		return assert.AnError // Force rollback
	})

	require.Error(t, err, "transaction should be rolled back")

	// Verify that no new account was created
	afterRollbackAccounts, err := accounts.Select()
	require.NoError(t, err)
	assert.Equal(t, len(initialAccounts), len(afterRollbackAccounts), "rollback should not persist changes")
}
