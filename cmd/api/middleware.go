package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read the authorization header
		authHeader := r.Header.Get("Authorization")

		// validate the header
		if authHeader == "" {
			app.unAuthorizedBasicError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		// split the header in Basic <base64>
		parts := strings.Split(authHeader, " ")

		// validate the header
		if len(parts) != 2 || parts[0] != "Basic" {
			app.unAuthorizedBasicError(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		// decode the base64 string
		decoded, err := base64.StdEncoding.DecodeString(parts[1])

		if err != nil {
			app.unAuthorizedBasicError(w, r, err)
			return
		}

		// validate the credentials [username, password]
		creds := strings.SplitN(string(decoded), ":", 2)

		username := app.config.auth.basic.username
		password := app.config.auth.basic.password
		if len(creds) != 2 || creds[0] != username || creds[0] != password {
			app.unAuthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r)

	})
}
