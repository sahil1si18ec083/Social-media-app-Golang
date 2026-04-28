package main

import (
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
		if len(parts) != 2 || parts[0] == "Beaer" {
			a.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}
		token := parts[1]
		fmt.Println(token)

		err := a.auth.ValidateToken(token)
		fmt.Println(err)
		if err != nil {
			a.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (a *application) BasicAuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-processing logic (e.g., checking headers)
		fmt.Println("yah yah")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			a.unauthorizedErrorResponse(w, r, errors.New("Authorization missing from header"))
			return // VERY IMPORTANT
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			fmt.Println("a")
			a.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			fmt.Println(err)
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
		fmt.Println(username, password)
		fmt.Println("bye")

		user, err := a.store.Users.GetByUsername(r.Context(), username)
		if err != nil {
			fmt.Println(err)
			a.unauthorizedErrorResponse(w, r, errors.New("invalid credentials"))
			return
		}
		err = bcrypt.CompareHashAndPassword(user.Password.Hash, []byte(password))
		fmt.Println(err)
		if err != nil {
			a.unauthorizedErrorResponse(w, r, errors.New("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r)

	})
}
