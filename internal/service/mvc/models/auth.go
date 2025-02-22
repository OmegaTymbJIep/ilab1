package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/omegatymbjiep/ilab1/internal/data"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/controllers/requests"
)

var ErrorEmailOrUsernameTaken = fmt.Errorf("email or username is already taken")
var ErrorUserNotFound = fmt.Errorf("user not found")
var ErrorInvalidPassword = fmt.Errorf("invalid password")

type Auth struct {
	db data.MainQ

	jwtSigningKey   jwk.Key
	jwtVerifyingKey jwk.Key
}

func NewAuth(db data.MainQ, jwtSigningKey jwk.Key) (*Auth, error) {
	jwtVerifyingKey, err := jwtSigningKey.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT public key: %w", err)
	}

	return &Auth{
		db:              db,
		jwtSigningKey:   jwtSigningKey,
		jwtVerifyingKey: jwtVerifyingKey,
	}, nil
}

func (a *Auth) Login(req *requests.Login) (string, error) {
	customers := a.db.Customers()

	customer := new(data.Customer)
	ok, err := customers.WhereEmail(req.Email).Get(customer)
	if err != nil {
		return "", fmt.Errorf("failed to get customer: %w", err)
	}

	if !ok {
		return "", ErrorUserNotFound
	}

	if bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(req.Password)) != nil {
		return "", ErrorInvalidPassword
	}

	token, err := a.newUserJWT(customer.ID)
	if err != nil {
		return "", fmt.Errorf("failed to create JWT: %w", err)
	}

	return token, nil
}

func (a *Auth) Register(req *requests.Register) (string, error) {
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
		return "", err
	}

	token, err := a.newUserJWT(customer.ID)
	if err != nil {
		return "", fmt.Errorf("failed to create JWT: %w", err)
	}

	return token, nil
}

func (a *Auth) newUserJWT(customerID uuid.UUID) (string, error) {
	token, err := jwt.NewBuilder().
		Claim("uid", customerID.String()).
		Claim("exp", time.Now().Add(time.Hour*24).Unix()).
		Claim("iat", time.Now().Unix()).
		Build()
	if err != nil {
		return "", fmt.Errorf("failed to build token: %w", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.ES256(), a.jwtSigningKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signed), nil
}
