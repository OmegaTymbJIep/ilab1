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

const customerIDJWTKey = "uid"

var ErrorEmailOrUsernameTaken = fmt.Errorf("email or username is already taken")
var ErrorUserNotFound = fmt.Errorf("user not found")
var ErrorInvalidPassword = fmt.Errorf("invalid password")

type JWTWithEat struct {
	Token      string
	Expiration time.Time
}

type Auth struct {
	db data.MainQ

	jwtSigningKey   jwk.Key
	jwtVerifyingKey jwk.Key
	jwtExpiry       time.Duration
}

func NewAuth(db data.MainQ, jwtSigningKey jwk.Key, jwtExpiry time.Duration) (*Auth, error) {
	jwtVerifyingKey, err := jwtSigningKey.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT public key: %w", err)
	}

	return &Auth{
		db:              db,
		jwtSigningKey:   jwtSigningKey,
		jwtVerifyingKey: jwtVerifyingKey,
		jwtExpiry:       jwtExpiry,
	}, nil
}

func (a *Auth) Login(req *requests.Login) (*JWTWithEat, error) {
	customers := a.db.Customers()

	customer := new(data.Customer)
	ok, err := customers.WhereEmail(req.Email).Get(customer)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	if !ok {
		return nil, ErrorUserNotFound
	}

	if bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(req.Password)) != nil {
		return nil, ErrorInvalidPassword
	}

	token, err := a.newCustomerJWT(customer.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT: %w", err)
	}

	return token, nil
}

func (a *Auth) Register(req *requests.Register) (*JWTWithEat, error) {
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

	token, err := a.newCustomerJWT(customer.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT: %w", err)
	}

	return token, nil
}

func (a *Auth) newCustomerJWT(customerID uuid.UUID) (*JWTWithEat, error) {
	eat := time.Now().Add(a.jwtExpiry)

	token, err := jwt.NewBuilder().
		Claim(customerIDJWTKey, customerID.String()).
		Claim(jwt.ExpirationKey, eat.Unix()).
		Claim(jwt.IssuedAtKey, time.Now().Unix()).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build token: %w", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.ES256(), a.jwtSigningKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &JWTWithEat{
		Token:      string(signed),
		Expiration: eat,
	}, nil
}

func (a *Auth) VerifyJWT(token string) (uuid.UUID, error) {
	parsedToken, err := jwt.Parse(
		[]byte(token),
		jwt.WithValidate(true),
		jwt.WithKey(jwa.ES256(), a.jwtVerifyingKey),
		jwt.WithClock(jwt.ClockFunc(time.Now)),
		jwt.WithMaxDelta(a.jwtExpiry, jwt.ExpirationKey, jwt.IssuedAtKey),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse/validate token: %w", err)
	}

	var customerIDRaw string
	if err = parsedToken.Get(customerIDJWTKey, &customerIDRaw); err != nil {
		return uuid.Nil, fmt.Errorf("failed to get uid: %w", err)
	}

	customerID, err := uuid.Parse(customerIDRaw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse uid: %w", err)
	}

	return customerID, nil
}
