package main

import (
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("internal server error: ", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Message: "internal server error",
	})
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("bad request error: ", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusBadRequest, ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("not found  error: ", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusNotFound, ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("resource already exists error:", r.Method, r.URL.Path, err.Error())
	writeJson(w, http.StatusConflict, ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}
