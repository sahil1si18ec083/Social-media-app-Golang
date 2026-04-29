package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (a *application) AuthTokenMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-processing logic (e.g., checking headers)

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			a.unauthorizedErrorResponse(w, r, errors.New("Authorization missing from header"))
			return // VERY IMPORTANT
		}
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			a.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}
		token := parts[1]

		claims, err := a.auth.ValidateToken(token)

		if err != nil {
			a.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}
		user, err := a.store.Users.GetById(r.Context(), claims.Subject)
		if err != nil {
			a.unauthorizedErrorResponse(w, r, errors.New("invalid credentials"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (a *application) BasicAuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-processing logic (e.g., checking headers)

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			a.unauthorizedErrorResponse(w, r, errors.New("Authorization missing from header"))
			return // VERY IMPORTANT
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {

			a.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {

			a.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid basic auth encoding"))
			return
		}
		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			a.unauthorizedErrorResponse(w, r, fmt.Errorf("invalid basic auth credentials"))
			return
		}

		username := creds[0]
		password := creds[1]

		user, err := a.store.Users.GetByUsername(r.Context(), username)
		if err != nil {

			a.unauthorizedErrorResponse(w, r, errors.New("invalid credentials"))
			return
		}
		err = bcrypt.CompareHashAndPassword(user.Password.Hash, []byte(password))

		if err != nil {
			a.unauthorizedErrorResponse(w, r, errors.New("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r)

	})
}
