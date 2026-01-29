package main

import "net/http"

func (a *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("heath api working 123w"))

}
