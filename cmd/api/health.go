package main

import (
	"fmt"
	"net/http"
)

func (a *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": "version",
	}
	fmt.Println("hey maa")

	err := writeJSON(w, http.StatusOK, data)
	if err != nil {
		a.internalServerError(w, r, err)
	}
	// a.store.Posts.Create(r.Context())

}
