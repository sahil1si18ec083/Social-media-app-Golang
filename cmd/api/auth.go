package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/mailer"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}
type LoginUserPayload struct {
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
		fmt.Println(err)

		a.badRequestResponse(w, r, err)
		return
	}
	err = Validate.Struct(&payload)
	if err != nil {
		fmt.Println(err)
		a.badRequestResponse(w, r, err)
		return
	}

	user := store.User{Username: payload.Username, Email: payload.Email, Role: 1}

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

	// send the invitation mail

	activationURL := fmt.Sprintf("%s/confirm/%s", a.config.frontendURL, plainToken)
	fmt.Println(activationURL)
	vars := WelcomeEmailData{Username: user.Username, ActivationURL: activationURL}

	isProdEnv := a.config.env == "production"

	_, err = a.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)

	if err != nil {
		fmt.Println(err)

		a.internalServerError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, "Activation URL send")
	if err != nil {
		fmt.Println(err)
		a.internalServerError(w, r, err)
		return
	}

}

func (a *application) LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload LoginUserPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		fmt.Println(err)

		a.badRequestResponse(w, r, err)
		return
	}
	err = Validate.Struct(&payload)
	if err != nil {
		fmt.Println(err)
		a.badRequestResponse(w, r, err)
		return
	}
	rcontext := r.Context()
	user, err := a.store.Users.GetByEmail(rcontext, payload.Email)

	if err != nil {
		a.unauthorizedErrorResponse(w, r, errors.New("invalid credentials"))
		return
	}
	if !user.IsActive {
		a.unauthorizedErrorResponse(w, r, errors.New("invalid credentials"))
		return
	}
	err = bcrypt.CompareHashAndPassword(user.Password.Hash, []byte(payload.Password))

	if err != nil {
		a.unauthorizedErrorResponse(w, r, errors.New("invalid credentials"))
		return
	}

	jwtToken, err := a.auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, FollowResponse{Message: jwtToken})
	if err != nil {

		a.internalServerError(w, r, err)
		return
	}

}
