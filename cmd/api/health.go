package main

import (
	"net/http"
)

// healthCheckHandler godoc
//
//	@Summary		Show application health
//	@Description	Returns the current application status, environment, and version.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/v1/healthcheck [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
