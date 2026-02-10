package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)

}

func writeJSONError(w http.ResponseWriter, status int, err error) {
	errorResponse := map[string]string{
		"error": err.Error(),
	}
	writeJSON(w, status, errorResponse)

}

func readJSON(w http.ResponseWriter, r *http.Request, body interface{}) error {

	fmt.Println("testing")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(body)
	return err

}
