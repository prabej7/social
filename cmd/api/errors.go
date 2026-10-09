package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("internal server error: ", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("bad request error: ", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("not found  error: ", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusNotFound, "not found error")
}
