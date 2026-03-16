package main

import (
	"fmt"
	"net/http"

	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (a *application) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	err := readJSON(w, r, &payload)
	if err != nil {

		a.badRequestResponse(w, r, err)
		return
	}
	err = Validate.Struct(&payload)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	fmt.Println(payload)

	user := store.User{Username: payload.Username, Email: payload.Email}

	err = user.Password.SetPassword(payload.Password, &user)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}
	err = a.store.Users.Create(r.Context(), &user)
	if err != nil {
		a.internalServerError(w, r, err)
	}

}
