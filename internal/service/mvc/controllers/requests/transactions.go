package requests

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type Deposit struct {
	AccountID    uuid.UUID `json:"account_id" validate:"required"`
	Amount       uint      `json:"amount" validate:"required,gt=0"`
	ATMSignature string    `json:"atm_signature" validate:"required"`
}

func NewDeposit(r *http.Request) (*Deposit, error) {
	var req Deposit
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	return &req, nil
}

type Withdrawal struct {
	AccountID uuid.UUID `json:"account_id" validate:"required"`
	Amount    uint      `json:"amount" validate:"required,gt=0"`
}

func NewWithdrawal(r *http.Request) (*Withdrawal, error) {
	var req Withdrawal
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	return &req, nil
}

type Transfer struct {
	SenderID    uuid.UUID `json:"sender_id" validate:"required"`
	RecipientID uuid.UUID `json:"recipient_id" validate:"required"`
	Amount      uint      `json:"amount" validate:"required,gt=0"`
}

func NewTransfer(r *http.Request) (*Transfer, error) {
	var req Transfer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	if req.SenderID == req.RecipientID {
		return nil, errors.New("sender and recipient must be different")
	}

	return &req, nil
}
