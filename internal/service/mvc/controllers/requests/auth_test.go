package requests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegister_ValidInput(t *testing.T) {
	form := url.Values{
		"email":         {"test@example.com"},
		"username":      {"validUser"},
		"password_hash": {"$2a$12$yUm4iidTxeB/pVObFGnp4OnbfABLRS8Ga01cY3rt5j.oOSjC5aeZu"},
		"first_name":    {"John"},
		"last_name":     {"Doe"},
	}
	r, _ := http.NewRequest("POST", "/register", nil)
	r.PostForm = form

	reg, err := NewRegister(r)
	assert.NoError(t, err)
	assert.NotNil(t, reg)
	assert.Equal(t, "test@example.com", reg.Email)
	assert.Equal(t, "validUser", reg.Username)
}

func TestNewRegister_InvalidEmail(t *testing.T) {
	form := url.Values{
		"email":         {"invalid-email"},
		"username":      {"validUser"},
		"password_hash": {"$2a$12$yUm4iidTxeB/pVObFGnp4OnbfABLRS8Ga01cY3rt5j.oOSjC5aeZu"},
	}
	r, _ := http.NewRequest("POST", "/register", nil)
	r.PostForm = form

	_, err := NewRegister(r)
	assert.Error(t, err)
}

func TestNewRegister_ShortUsername(t *testing.T) {
	form := url.Values{
		"email":         {"test@example.com"},
		"username":      {"ab"}, // Too short
		"password_hash": {"$2a$12$yUm4iidTxeB/pVObFGnp4OnbfABLRS8Ga01cY3rt5j.oOSjC5aeZu"},
	}
	r, _ := http.NewRequest("POST", "/register", nil)
	r.PostForm = form

	_, err := NewRegister(r)
	assert.Error(t, err)
}

func TestNewRegister_InvalidBcryptHash(t *testing.T) {
	form := url.Values{
		"email":         {"test@example.com"},
		"username":      {"validUser"},
		"password_hash": {"invalidhash"}, // Not a bcrypt hash
	}
	r, _ := http.NewRequest("POST", "/register", nil)
	r.PostForm = form

	_, err := NewRegister(r)
	assert.Error(t, err)
}

func TestNewRegister_ValidOptionalNames(t *testing.T) {
	form := url.Values{
		"email":         {"test@example.com"},
		"username":      {"validUser"},
		"password_hash": {"$2a$12$yUm4iidTxeB/pVObFGnp4OnbfABLRS8Ga01cY3rt5j.oOSjC5aeZu"},
	}
	r, _ := http.NewRequest("POST", "/register", nil)
	r.PostForm = form

	reg, err := NewRegister(r)
	assert.NoError(t, err)
	assert.NotNil(t, reg)
	assert.Empty(t, reg.FirstName)
	assert.Empty(t, reg.LastName)
}
