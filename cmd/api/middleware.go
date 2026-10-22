package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read the authorization header
		authHeader := r.Header.Get("Authorization")
		fmt.Printf("AUTH HEADER = %q\n", authHeader)
		// validate the header
		if authHeader == "" {
			app.unAuthorizedBasicError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ") // Bearer <token>

		// validate the header
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unAuthorizedBasicError(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)

		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)

		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, userID)

		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
