package requests

import (
	"net/http"
)

type Register struct {
	Email        string `validate:"required,email"`
	Username     string `validate:"required,min=3,max=50"`
	PasswordHash string `validate:"required,bcrypt"`
	FirstName    string `validate:"omitempty,min=2,max=50"`
	LastName     string `validate:"omitempty,min=2,max=50"`
}

// NewRegister parses HTTP request and validates input
func NewRegister(r *http.Request) (*Register, error) {
	req := Register{
		Email:        r.FormValue("email"),
		Username:     r.FormValue("username"),
		PasswordHash: r.FormValue("password_hash"),
		FirstName:    r.FormValue("first_name"),
		LastName:     r.FormValue("last_name"),
	}

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	return &req, nil
}

type Login struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,hexadecimal,len=64"`
}

// NewLogin parses HTTP request and validates input
func NewLogin(r *http.Request) (*Login, error) {
	req := Login{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	return &req, nil
}
