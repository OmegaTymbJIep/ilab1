package postgres

import (
	"testing"

	"github.com/google/uuid"
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

func TestTransactionsCRUD(t *testing.T) {
	db := newTestMainQ(t)
	transactions := db.Transactions()
	accounts := db.Accounts()

	// Create sender and recipient accounts.
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

	// Insert a valid transfer transaction (both sender and recipient must be non-nil).
	txn := &data.Transaction{
		Type:      data.TransferTransaction,
		Amount:    300,
		Sender:    sender.ID,
		Recipient: recipient.ID,
	}
	err = transactions.Insert(txn)
	require.NoError(t, err)

	// Retrieve the transaction using a filter (by sender).
	fetched := new(data.Transaction)
	ok, err := transactions.WhereSender(txn.Sender).Get(fetched)
	require.NoError(t, err)
	require.True(t, ok, "inserted transaction not found")
	assert.Equal(t, txn.Amount, fetched.Amount)
	assert.Equal(t, txn.Type, fetched.Type)

	// Update the transaction amount.
	txn.Amount = 350
	err = transactions.Update(txn)
	require.NoError(t, err)

	// Verify the update.
	updated := new(data.Transaction)
	ok, err = transactions.WhereSender(txn.Sender).Get(updated)
	require.NoError(t, err)
	require.True(t, ok, "updated transaction not found")
	assert.Equal(t, uint(350), updated.Amount)

	// Delete the transaction.
	err = transactions.Delete(txn.ID)
	require.NoError(t, err)

	// Verify deletion.
	deleted := new(data.Transaction)
	ok, err = transactions.WhereSender(txn.Sender).Get(deleted)
	require.NoError(t, err)
	assert.False(t, ok, "deleted transaction still exists")
}

func TestInvalidDeposit(t *testing.T) {
	db := newTestMainQ(t)
	transactions := db.Transactions()
	accounts := db.Accounts()

	// Create a valid recipient account.
	recipient := &data.Account{
		Name:    "recipient_invalid_deposit",
		Balance: 500,
	}
	err := accounts.Insert(recipient)
	require.NoError(t, err)

	// Try to create a deposit with a non-nil sender.
	txnDeposit := &data.Transaction{
		Type:         data.DepositTransaction,
		Amount:       200,
		Sender:       uuid.New(), // Invalid: sender should be nil for deposits.
		Recipient:    recipient.ID,
		ATMSignature: "invalid_deposit",
	}
	err = transactions.Insert(txnDeposit)
	require.Error(t, err, "expected error when inserting deposit with non-null sender")
}

func TestInvalidWithdrawal(t *testing.T) {
	db := newTestMainQ(t)
	transactions := db.Transactions()
	accounts := db.Accounts()

	// Create a valid sender account.
	sender := &data.Account{
		Name:    "sender_invalid_withdrawal",
		Balance: 1000,
	}
	err := accounts.Insert(sender)
	require.NoError(t, err)

	// Try to create a withdrawal with a non-nil recipient.
	txnWithdrawal := &data.Transaction{
		Type:         data.WithdrawalTransaction,
		Amount:       100,
		Sender:       sender.ID,
		Recipient:    uuid.New(), // Invalid: recipient should be nil for withdrawals.
		ATMSignature: "invalid_withdrawal",
	}
	err = transactions.Insert(txnWithdrawal)
	require.Error(t, err, "expected error when inserting withdrawal with non-null recipient")
}

func TestInvalidTransfer(t *testing.T) {
	db := newTestMainQ(t)
	transactions := db.Transactions()
	accounts := db.Accounts()

	// Create valid accounts for sender and recipient.
	sender := &data.Account{
		Name:    "sender_invalid_transfer",
		Balance: 1000,
	}
	recipient := &data.Account{
		Name:    "recipient_invalid_transfer",
		Balance: 500,
	}
	err := accounts.Insert(sender)
	require.NoError(t, err)
	err = accounts.Insert(recipient)
	require.NoError(t, err)

	// Try to create a transfer with a missing recipient (recipient nil).
	txnTransferMissingRecipient := &data.Transaction{
		Type:      data.TransferTransaction,
		Amount:    300,
		Sender:    sender.ID,
		Recipient: uuid.Nil, // Invalid: transfer must have a non-nil recipient.
	}
	err = transactions.Insert(txnTransferMissingRecipient)
	require.Error(t, err, "expected error when inserting transfer with missing recipient")

	// Also try a transfer with a missing sender.
	txnTransferMissingSender := &data.Transaction{
		Type:      data.TransferTransaction,
		Amount:    300,
		Sender:    uuid.Nil, // Invalid: transfer must have a non-nil sender.
		Recipient: recipient.ID,
	}
	err = transactions.Insert(txnTransferMissingSender)
	require.Error(t, err, "expected error when inserting transfer with missing sender")
}

func TestTransactionsFilters(t *testing.T) {
	db := newTestMainQ(t)
	transactions := db.Transactions()
	accounts := db.Accounts()

	// Create two accounts.
	account1 := &data.Account{
		Name:    "account1",
		Balance: 1000,
	}
	account2 := &data.Account{
		Name:    "account2",
		Balance: 500,
	}
	err := accounts.Insert(account1)
	require.NoError(t, err)
	err = accounts.Insert(account2)
	require.NoError(t, err)

	// Insert a valid deposit: for a deposit, sender must be nil (uuid.Nil) and recipient non-nil.
	txnDeposit := &data.Transaction{
		Type:   data.DepositTransaction,
		Amount: 200,
		// Sender is omitted (zero value, uuid.Nil) to represent SQL NULL.
		Recipient:    account1.ID,
		ATMSignature: "deposit_signature",
	}
	err = transactions.Insert(txnDeposit)
	require.NoError(t, err)

	// Insert a valid withdrawal: for a withdrawal, recipient must be nil (uuid.Nil) and sender non-nil.
	txnWithdrawal := &data.Transaction{
		Type:   data.WithdrawalTransaction,
		Amount: 100,
		Sender: account1.ID,
		// Recipient is omitted (zero value) to represent SQL NULL.
	}
	err = transactions.Insert(txnWithdrawal)
	require.NoError(t, err)

	// Insert a valid transfer: both sender and recipient must be non-nil.
	txnTransfer := &data.Transaction{
		Type:         data.TransferTransaction,
		Amount:       300,
		Sender:       account1.ID,
		Recipient:    account2.ID,
		ATMSignature: "transfer_signature",
	}
	err = transactions.Insert(txnTransfer)
	require.NoError(t, err)

	// Test filtering by transaction type.
	depositTxns, err := transactions.WhereType(data.DepositTransaction).Select()
	require.NoError(t, err)
	assert.NotEmpty(t, depositTxns, "expected at least one deposit transaction")

	// Test filtering by sender: account1 should appear in both withdrawal and transfer.
	senderTxns, err := transactions.WhereSender(account1.ID).Select()
	require.NoError(t, err)
	// Deposit should not be returned since its sender is nil.
	assert.GreaterOrEqual(t, len(senderTxns), 2, "expected at least two transactions with account1 as sender")

	// Test filtering by recipient: account2 should appear in the transfer transaction.
	recipientTxns, err := transactions.WhereRecipient(account2.ID).Select()
	require.NoError(t, err)
	assert.Len(t, recipientTxns, 1, "expected exactly one transaction with account2 as recipient")

	// Test filtering by account: WhereAccount should return transactions where account1 is either sender or recipient.
	accountTxns, err := transactions.WhereAccount(account1.ID).Select()
	require.NoError(t, err)
	// account1 is recipient in the deposit and sender in the withdrawal and transfer.
	assert.GreaterOrEqual(t, len(accountTxns), 3, "expected at least three transactions involving account1")
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
