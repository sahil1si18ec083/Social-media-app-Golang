package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
		int64val, err := strconv.ParseInt(claims.Subject, 10, 64)

		user, err := a.cacheStorage.Users.Get(r.Context(), int64val)
		if err != nil {
			// fmt.Println("aa ja")
			a.internalServerError(w, r, err)
			return
		}
		// fmt.Println("yoo")
		// fmt.Println(user)
		var storeduser *store.User
		if user != nil {
			// Cache hit
			fmt.Println("hit")
			storeduser = user
		} else {
			// Cache miss
			fmt.Println("miss")
			storeduser, err = a.store.Users.GetById(r.Context(), claims.Subject)
			if err != nil {
				a.unauthorizedErrorResponse(w, r, errors.New("invalid credentials"))
				return
			}

			err = a.cacheStorage.Users.Set(r.Context(), storeduser)
			fmt.Println(err, "hhh")
		}
		// fmt.Println(storeduser, "4jj")
		ctx := r.Context()
		ctx = context.WithValue(ctx, userContextKey, storeduser)
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
		// fmt.Println(user)s
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

func (a *application) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-processing logic (e.g., checking headers)
		user, ok := r.Context().Value(userContextKey).(*store.User)
		// fmt.Println(user)s
		if !ok || user == nil {
			a.unauthorizedErrorResponse(w, r, errors.New("user missing from context"))
			return
		}
		userId := user.ID
		key := strconv.FormatInt(userId, 10)
		allowcheck, remaining, resetAfter := a.ratelimiter.Allow(key)
		w.Header().Set("RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(a.config.ratelimiter.RequestsPerTimeFrame))

		if !allowcheck {
			w.Header().Set("Retry-After", strconv.Itoa(int(resetAfter.Seconds())))
			http.Error(w, "Rate limit exceeded. Please slow down.", http.StatusTooManyRequests)
			return
		} else {
			next.ServeHTTP(w, r)
		}

	})
}
