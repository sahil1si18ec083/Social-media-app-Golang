package main

import "net/http"

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")

}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {

	writeJSONError(w, http.StatusNotFound, "not found")
}
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	writeJSONError(w, http.StatusConflict, err.Error())
}
