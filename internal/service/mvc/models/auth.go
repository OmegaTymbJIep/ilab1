package models

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/omegatymbjiep/ilab1/internal/data"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/requests"
)

var ErrorEmailOrUsernameTaken = fmt.Errorf("email or username is already taken")

type Auth struct {
	db data.MainQ
}

func NewAuth(db data.MainQ) *Auth {
	return &Auth{db: db}
}

func (a *Auth) Register(req *requests.Register) (*uuid.UUID, error) {
	customers := a.db.Customers()

	customer := &data.Customer{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: req.PasswordHash,
		FirstName:    &req.FirstName,
		LastName:     &req.LastName,
	}

	if err := a.db.Transaction(func() error {
		isUnique, err := customers.IsUnique(req.Email, req.Username)
		if err != nil {
			return fmt.Errorf("failed to check uniqueness: %w", err)
		}

		if !isUnique {
			return ErrorEmailOrUsernameTaken
		}

		if err = customers.Insert(customer); err != nil {
			return fmt.Errorf("failed to insert customer: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &customer.ID, nil
}
