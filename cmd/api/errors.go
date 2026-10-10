package main

import (
	"log"
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("internal server error: ", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Message: "internal server error",
	})
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("bad request error: ", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusBadRequest, ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("not found  error: ", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusNotFound, ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}
