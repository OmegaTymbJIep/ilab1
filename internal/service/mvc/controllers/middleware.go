package controllers

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Auth) VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt := r.Header.Get("Authorization")
		if jwt == "" {
			Unauthorized(w, r, fmt.Errorf("no jwt provided"))
			return
		}

		customerID, err := c.model.VerifyJWT(jwt)
		if err != nil {
			Unauthorized(w, r, fmt.Errorf("invalid jwt: %w", err))
			return
		}

		ctx := context.WithValue(r.Context(), customerIDCtxKey, customerID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
