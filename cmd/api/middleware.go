package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
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

func (a *application) checkPostOwnership(rolename string, next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-processing logic (e.g., checking headers)
		user, ok := r.Context().Value(userContextKey).(*store.User)
		fmt.Println(user)
		if !ok || user == nil {
			a.unauthorizedErrorResponse(w, r, errors.New("user missing from context"))
			return
		}
		post, ok := r.Context().Value(postContextKey).(*store.Post)
		if !ok || post == nil {
			a.notFoundResponse(w, r, errors.New("post missing from context"))
			return
		}
		if user.ID == post.UserID {
			next.ServeHTTP(w, r)
			return
		}
		role, err := a.store.Roles.GetByRolename(r.Context(), rolename)
		if err != nil {
			a.internalServerError(w, r, err)
			return
		}
		if user.Role >= role.Level {
			next.ServeHTTP(w, r)
			return
		} else {
			a.forbiddenResponse(w, r, errors.New(""))
		}

	})

}
