package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prabej7/social/internal/store"
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
		user, err := app.getUser(ctx, userID)

		if err != nil {
			app.unAuthorizedError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func (app *application) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.getUserFromContext(r, userContextKey)
		post := getPostContext(r)

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := app.checkRolePrecedence(r.Context(), user, role)
		if err != nil {
			app.internalServerError(w, r, err)
		}

		if !allowed {
			app.forbiddenError(w, r, fmt.Errorf("forbidden"))
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}

func (app *application) getUser(ctx context.Context, id int64) (*store.User, error) {
	var user *store.User
	isCacheEnabled := app.config.redisConfig.enabled
	if isCacheEnabled {
		user, err := app.cache.Users.Get(ctx, id)

		if err != nil {
			return nil, err
		}

		if user != nil {
			return user, nil
		}
	}

	user, err := app.store.Users.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	if isCacheEnabled {
		if err := app.cache.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
