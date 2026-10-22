package main

import (
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJson(w, http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Message: "internal server error",
	})
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJson(w, http.StatusBadRequest, ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found  error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJson(w, http.StatusNotFound, ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("resource already exists", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJson(w, http.StatusConflict, ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}

func (app *application) unAuthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorized", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJson(w, http.StatusUnauthorized, ErrorResponse{
		Success: false,
		Message: "unauthorized",
	})
}

func (app *application) unAuthorizedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorized", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8`)

	writeJson(w, http.StatusUnauthorized, ErrorResponse{
		Success: false,
		Message: "unauthorized",
	})
}

func (app *application) forbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("forbidden error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJson(w, http.StatusForbidden, ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}
