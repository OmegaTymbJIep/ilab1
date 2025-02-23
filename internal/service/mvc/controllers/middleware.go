package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Auth) VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customerID, err := c.verifyJWT(r)
		if err != nil {
			Unauthorized(w, r, fmt.Errorf("invalid jwt: %w", err))
			return
		}

		ctx := context.WithValue(r.Context(), customerIDCtxKey, customerID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *Auth) verifyJWT(r *http.Request) (uuid.UUID, error) {
	jwtCookie, err := r.Cookie(JWTCookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return uuid.Nil, fmt.Errorf("no jwt provided")
	}

	customerID, err := c.model.VerifyJWT(jwtCookie.Value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid jwt: %w", err)
	}

	return customerID, nil
}
