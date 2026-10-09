package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("internal server error: ", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusInternalServerError, "internal server error")
}
