package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
)

func (a *application) GetUserHandler(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "userID")

	rcontext := r.Context()
	user, err := a.store.Users.GetById(rcontext, userID)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, store.ErrNotFound) {
			fmt.Println("yes")
			a.notFoundResponse(w, r, err)
			return
		}

		a.internalServerError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, user)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}
}
