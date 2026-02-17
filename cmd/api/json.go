package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}
func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	res := map[string]interface{}{
		"response": data,
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(res)

}

func writeJSONError(w http.ResponseWriter, status int, err string) {
	errorResponse := map[string]string{
		"error": err,
	}
	writeJSON(w, status, errorResponse)

}

func readJSON(w http.ResponseWriter, r *http.Request, body interface{}) error {

	fmt.Println("testing")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(body)
	return err

}
