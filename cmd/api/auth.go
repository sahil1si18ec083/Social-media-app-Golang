package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/mailer"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}
type WelcomeEmailData struct {
	Username      string
	ActivationURL string
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
		fmt.Println(err)
		a.internalServerError(w, r, err)
		return
	}
	plainToken := uuid.New().String()
	fmt.Println(plainToken, "              HHHH")
	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	expiry_time := a.config.auth.token.exp

	err = a.store.Users.CreateAndInvite(r.Context(), &user, hashToken, expiry_time)
	if err != nil {
		fmt.Println(err)
		a.internalServerError(w, r, err)
		return
	}
	fmt.Println(user)
	// send the invitation mail

	activationURL := fmt.Sprintf("%s/confirm/%s", a.config.frontendURL, plainToken)
	fmt.Println(activationURL)
	vars := WelcomeEmailData{Username: user.Username, ActivationURL: activationURL}
	isProdEnv := a.config.env == "production"

	status, err := a.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		// do something
		// saga pattern
		fmt.Println(err)
	}
	fmt.Println(status)
	err = writeJSON(w, http.StatusOK, user)
	if err != nil {

		a.internalServerError(w, r, err)
		return
	}

}
