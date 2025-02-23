package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateAccount struct {
	Name string `json:"name" validate:"omitempty,min=3,max=50"`
}

// NewCreateAccount parses HTTP request and validates input
func NewCreateAccount(r *http.Request) (*CreateAccount, error) {
	var requestBody CreateAccount

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, fmt.Errorf("failed to decode json request body: %w", err)
	}

	if err := validate.Struct(requestBody); err != nil {
		return nil, err
	}

	return &requestBody, nil
}
