package models

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"

	"github.com/omegatymbjiep/ilab1/internal/data"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/requests"
)

var ErrorInsufficientFunds = errors.New("insufficient funds")
var ErrorRecipientNotFound = errors.New("recipient account not found")
var ErrorATMSignatureNotUnique = errors.New("ATM signature not unique")
var ErrorInvalidATMSignature = errors.New("invalid ATM signature")

type Transactions struct {
	db           data.MainQ
	atmPublicKey *ecdsa.PublicKey
}

func NewTransactions(db data.MainQ, atmPublicKey *ecdsa.PublicKey) *Transactions {
	return &Transactions{
		db:           db,
		atmPublicKey: atmPublicKey,
	}
}

func (m *Transactions) DepositFunds(customerID uuid.UUID, req *requests.Deposit) (int, error) {
	ok, err := m.db.CustomersAccounts().HasAccount(customerID, req.AccountID)
	if err != nil {
		return 0, fmt.Errorf("failed to check account existence: %w", err)
	}
	if !ok {
		return 0, ErrorAccountNotFound
	}

	ok, err = m.verifySignature(&SignedTransaction{
		AccountID: req.AccountID.String(),
		Amount:    req.Amount,
		Signature: req.ATMSignature,
	})
	if !ok || err != nil {
		return 0, ErrorInvalidATMSignature
	}

	var newBalance int
	err = m.db.Transaction(func() error {
		account := new(data.Account)
		ok, err = m.db.Accounts().WhereID(req.AccountID).Get(account)
		if err != nil {
			return fmt.Errorf("failed to get account: %w", err)
		}
		if !ok {
			return ErrorAccountNotFound
		}

		transaction := &data.Transaction{
			Type:         data.DepositTransaction,
			Amount:       req.Amount,
			Recipient:    req.AccountID,
			ATMSignature: req.ATMSignature,
		}

		if err = m.db.Transactions().Insert(transaction); err != nil {
			if isATMNotUniqueError(err) {
				return ErrorATMSignatureNotUnique
			}

			return fmt.Errorf("failed to create transaction: %w", err)
		}

		account.Balance += int(req.Amount)

		if err = m.db.Accounts().Update(account); err != nil {
			return fmt.Errorf("failed to update account balance: %w", err)
		}

		newBalance = account.Balance
		return nil
	})

	return newBalance, err
}

func (m *Transactions) WithdrawFunds(customerID uuid.UUID, req *requests.Withdrawal) (int, error) {
	ok, err := m.db.CustomersAccounts().HasAccount(customerID, req.AccountID)
	if err != nil {
		return 0, fmt.Errorf("failed to check account existence: %w", err)
	}
	if !ok {
		return 0, ErrorAccountNotFound
	}

	var newBalance int

	err = m.db.Transaction(func() error {
		account := new(data.Account)
		ok, err = m.db.Accounts().WhereID(req.AccountID).Get(account)
		if err != nil {
			return fmt.Errorf("failed to get account: %w", err)
		}
		if !ok {
			return ErrorAccountNotFound
		}

		if account.Balance < int(req.Amount) {
			return ErrorInsufficientFunds
		}

		transaction := &data.Transaction{
			Type:   data.WithdrawalTransaction,
			Amount: req.Amount,
			Sender: req.AccountID,
		}

		if err = m.db.Transactions().Insert(transaction); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		account.Balance -= int(req.Amount)

		if err = m.db.Accounts().Update(account); err != nil {
			return fmt.Errorf("failed to update account balance: %w", err)
		}

		newBalance = account.Balance
		return nil
	})

	return newBalance, err
}

func (m *Transactions) TransferFunds(customerID uuid.UUID, req *requests.Transfer) (int, error) {
	ok, err := m.db.CustomersAccounts().HasAccount(customerID, req.SenderID)
	if err != nil {
		return 0, fmt.Errorf("failed to check account existence: %w", err)
	}
	if !ok {
		return 0, ErrorAccountNotFound
	}

	var senderBalance int

	err = m.db.Transaction(func() error {
		sender := new(data.Account)
		ok, err = m.db.Accounts().WhereID(req.SenderID).Get(sender)
		if err != nil {
			return fmt.Errorf("failed to get sender account: %w", err)
		}
		if !ok {
			return ErrorAccountNotFound
		}

		recipient := new(data.Account)
		ok, err = m.db.Accounts().WhereID(req.RecipientID).Get(recipient)
		if err != nil {
			return fmt.Errorf("failed to get recipient account: %w", err)
		}
		if !ok {
			return ErrorRecipientNotFound
		}

		if sender.Balance < int(req.Amount) {
			return ErrorInsufficientFunds
		}

		transaction := &data.Transaction{
			Type:      data.TransferTransaction,
			Amount:    req.Amount,
			Sender:    req.SenderID,
			Recipient: req.RecipientID,
		}

		if err = m.db.Transactions().Insert(transaction); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		sender.Balance -= int(req.Amount)
		recipient.Balance += int(req.Amount)

		if err = m.db.Accounts().Update(sender); err != nil {
			return fmt.Errorf("failed to update sender balance: %w", err)
		}

		if err = m.db.Accounts().Update(recipient); err != nil {
			return fmt.Errorf("failed to update recipient balance: %w", err)
		}

		senderBalance = sender.Balance
		return nil
	})

	return senderBalance, err
}

type SignedTransaction struct {
	AccountID string `json:"account_id"`
	Amount    uint   `json:"amount"`
	Signature string `json:"-"`
}

type ECDSASignature struct {
	R, S *big.Int
}

// verifySignature verifies the signature of a signed transaction
func (m *Transactions) verifySignature(signedTx *SignedTransaction) (bool, error) {
	// Marshal the transaction to JSON
	txJSON, err := json.Marshal(signedTx)
	if err != nil {
		return false, fmt.Errorf("error marshaling transaction: %v", err)
	}

	// Calculate SHA-256 hash of the JSON
	hash := sha256.Sum256(txJSON)

	// Decode the signature from Base64
	signatureBytes, err := base64.StdEncoding.DecodeString(signedTx.Signature)
	if err != nil {
		return false, fmt.Errorf("error decoding signature: %v", err)
	}

	// Unmarshal the signature
	var signature ECDSASignature
	if _, err := asn1.Unmarshal(signatureBytes, &signature); err != nil {
		return false, fmt.Errorf("error unmarshaling signature: %v", err)
	}

	// Verify the signature
	return ecdsa.Verify(m.atmPublicKey, hash[:], signature.R, signature.S), nil
}

func isATMNotUniqueError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
