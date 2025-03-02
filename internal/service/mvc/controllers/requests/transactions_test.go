package requests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewDeposit(t *testing.T) {
	accountID := uuid.New()
	tests := []struct {
		name    string
		body    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid input",
			body: map[string]interface{}{
				"account_id":    accountID,
				"amount":        1000,
				"atm_signature": "valid-signature",
			},
			wantErr: false,
		},
		{
			name: "missing account_id",
			body: map[string]interface{}{
				"amount":        1000,
				"atm_signature": "valid-signature",
			},
			wantErr: true,
		},
		{
			name: "invalid amount",
			body: map[string]interface{}{
				"account_id":    accountID,
				"amount":        0,
				"atm_signature": "valid-signature",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			r, _ := http.NewRequest("POST", "/deposit", bytes.NewBuffer(body))

			got, err := NewDeposit(r)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, accountID, got.AccountID)
				assert.Equal(t, uint(1000), got.Amount)
				assert.Equal(t, "valid-signature", got.ATMSignature)
			}
		})
	}
}

func TestNewWithdrawal(t *testing.T) {
	accountID := uuid.New()
	tests := []struct {
		name    string
		body    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid input",
			body: map[string]interface{}{
				"account_id": accountID,
				"amount":     500,
			},
			wantErr: false,
		},
		{
			name: "missing account_id",
			body: map[string]interface{}{
				"amount": 500,
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			body: map[string]interface{}{
				"account_id": accountID,
				"amount":     0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			r, _ := http.NewRequest("POST", "/withdraw", bytes.NewBuffer(body))

			got, err := NewWithdrawal(r)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, accountID, got.AccountID)
				assert.Equal(t, uint(500), got.Amount)
			}
		})
	}
}

func TestNewTransfer(t *testing.T) {
	senderID := uuid.New()
	recipientID := uuid.New()
	tests := []struct {
		name    string
		body    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid input",
			body: map[string]interface{}{
				"sender_id":    senderID,
				"recipient_id": recipientID,
				"amount":       1000,
			},
			wantErr: false,
		},
		{
			name: "same sender and recipient",
			body: map[string]interface{}{
				"sender_id":    senderID,
				"recipient_id": senderID,
				"amount":       1000,
			},
			wantErr: true,
		},
		{
			name: "missing sender",
			body: map[string]interface{}{
				"recipient_id": recipientID,
				"amount":       1000,
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			body: map[string]interface{}{
				"sender_id":    senderID,
				"recipient_id": recipientID,
				"amount":       0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			r, _ := http.NewRequest("POST", "/transfer", bytes.NewBuffer(body))

			got, err := NewTransfer(r)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, senderID, got.SenderID)
				assert.Equal(t, recipientID, got.RecipientID)
				assert.Equal(t, uint(1000), got.Amount)
			}
		})
	}
}
