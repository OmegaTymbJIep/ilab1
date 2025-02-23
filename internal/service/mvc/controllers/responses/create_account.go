package responses

import (
	"time"

	"github.com/google/uuid"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

type CreateAccount struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Balance   int       `json:"balance"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func NewCreateAccount(account *data.Account) *CreateAccount {
	return &CreateAccount{
		ID:        account.ID,
		Name:      account.Name,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt.Format(time.DateTime),
		UpdatedAt: account.UpdatedAt.Format(time.DateTime),
	}
}
