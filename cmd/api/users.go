package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
)

const userContextKey = contextKey("user")

type FollowRequestPayload struct {
	UserId int64 `json:"user_id"`
}
type FollowResponse struct {
	Message string `json:"message"`
}

func (a *application) GetUserHandler(w http.ResponseWriter, r *http.Request) {

	rcontext := r.Context()
	user, ok := GetUserFromContext(rcontext)
	if !ok {
		a.internalServerError(w, r, errors.New("post missing from context"))
		return // VERY IMPORTANT
	}

	err := writeJSON(w, http.StatusOK, user)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}
}

func (a *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		user, err := a.store.Users.GetById(r.Context(), userID)
		if err != nil {

			if errors.Is(err, store.ErrNotFound) {
				a.notFoundResponse(w, r, err)
				return
			}

			a.internalServerError(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func GetUserFromContext(ctx context.Context) (*store.User, bool) {

	user, ok := ctx.Value(userContextKey).(*store.User)
	if user == nil {
		return nil, false
	}
	if !ok {
		return nil, false
	}
	return user, true

}

func (a *application) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	rcontext := r.Context()
	user, ok := GetUserFromContext(rcontext)
	if !ok {
		a.internalServerError(w, r, errors.New("post missing from context"))
		return // VERY IMPORTANT
	}
	var payload FollowRequestPayload

	err := readJSON(w, r, &payload)
	if err != nil {

		if !errors.Is(err, io.EOF) {
			a.badRequestResponse(w, r, err)
			return
		}

	}
	err = Validate.Struct(&payload)
	if err != nil {

		a.badRequestResponse(w, r, err)
		return
	}

	follower_id := user.ID
	userId := payload.UserId

	err = a.store.Follower.Follow(rcontext, strconv.FormatInt(follower_id, 10), strconv.FormatInt(userId, 10))
	if err != nil {
		if errors.Is(err, store.ErrAlreadyFollowing) || errors.Is(err, store.ErrSelfFollow) {
			a.conflictResponse(w, r, err)
			return
		}

		if errors.Is(err, store.ErrNotFound) {
			a.notFoundResponse(w, r, err)
			return
		}

		a.internalServerError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, FollowResponse{Message: "Followed Successfully"})
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}

}

func (a *application) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	rcontext := r.Context()
	user, ok := GetUserFromContext(rcontext)
	if !ok {
		a.internalServerError(w, r, errors.New("post missing from context"))
		return
	}
	var payload FollowRequestPayload

	err := readJSON(w, r, &payload)
	if err != nil {

		if !errors.Is(err, io.EOF) {
			a.badRequestResponse(w, r, err)
			return
		}

	}
	err = Validate.Struct(&payload)
	if err != nil {

		a.badRequestResponse(w, r, err)
		return
	}

	follower_id := user.ID
	userId := payload.UserId

	err = a.store.Follower.UnFollow(rcontext, strconv.FormatInt(follower_id, 10), strconv.FormatInt(userId, 10))
	if err != nil {

		if errors.Is(err, store.ErrFollowNotFound) {
			a.notFoundResponse(w, r, err)
			return
		}

		a.internalServerError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, FollowResponse{Message: "UnFollowed Successfully"})
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}
}
