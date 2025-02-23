package models

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/stretchr/testify/require"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

// Example of a mocked MainQ. If you only care about JWT logic, a nil or no-op mock is fine.
type mockDB struct{}

func (m *mockDB) Customers() data.Customers                 { return nil }
func (m *mockDB) Accounts() data.Accounts                   { return nil }
func (m *mockDB) CustomersAccounts() data.CustomersAccounts { return nil }
func (m *mockDB) Withdrawals() data.Withdrawals             { return nil }
func (m *mockDB) Deposits() data.Deposits                   { return nil }
func (m *mockDB) Transfers() data.Transfers                 { return nil }
func (m *mockDB) Transaction(fn func() error) error         { return fn() }
func (m *mockDB) IsolatedTransaction(_ sql.IsolationLevel, fn func() error) error {
	return fn()
}

func TestNewCustomerJWT(t *testing.T) {
	// Generate an ephemeral ECDSA key pair for testing
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed to generate test ECDSA key")

	privJWK, err := jwk.Import(ecdsaKey)
	require.NoError(t, err, "failed to parse private key into jwk")

	// Instantiate the Auth model with a short expiry to test
	auth, err := NewAuth(&mockDB{}, privJWK, 2*time.Minute)
	require.NoError(t, err, "failed to instantiate Auth")

	// Use any random UUID for the customer
	custID := uuid.New()

	token, err := auth.newCustomerJWT(custID)
	require.NoError(t, err, "failed to create JWT")

	require.NotEmpty(t, token, "JWT should not be empty")
}

func TestVerifyJWT(t *testing.T) {
	// Generate an ephemeral ECDSA key pair for testing
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "failed to generate test ECDSA key")

	privJWK, err := jwk.Import(ecdsaKey)
	require.NoError(t, err, "failed to parse private key into jwk")

	// Set a short expiry so we can test expiry scenarios too
	auth, err := NewAuth(&mockDB{}, privJWK, 2*time.Second)
	require.NoError(t, err, "failed to instantiate Auth")

	// Generate a valid token
	custID := uuid.New()
	validToken, err := auth.newCustomerJWT(custID)
	require.NoError(t, err, "failed to create valid JWT")

	t.Run("valid token", func(t *testing.T) {
		parsedID, err := auth.VerifyJWT(validToken)
		require.NoError(t, err, "valid token should verify without error")
		require.Equal(t, custID, parsedID, "parsed ID should match the original customer ID")
	})

	t.Run("malformed token", func(t *testing.T) {
		_, err := auth.VerifyJWT("this-is-not-a-valid-token")
		require.Error(t, err, "should fail to parse a malformed token")
	})

	t.Run("signature changed / tampered token", func(t *testing.T) {
		// Just remove the last character from the token or otherwise alter it
		tampered := validToken[:len(validToken)-1] + "x"
		_, err := auth.VerifyJWT(tampered)
		require.Error(t, err, "should fail with a tampered signature")
	})

	t.Run("expired token", func(t *testing.T) {
		// Wait for the token to expire
		time.Sleep(3 * time.Second)
		_, err := auth.VerifyJWT(validToken)
		require.Error(t, err, "should fail when token is expired")
		require.True(t, errors.Is(err, jwt.TokenExpiredError()), "expired token error expected")
	})
}
