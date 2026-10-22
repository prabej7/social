package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prabej7/social/internal/mailer"
	"github.com/prabej7/social/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=75"`
}

func HashPassword(text string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

type UserWithToken struct {
	store.User
	Token string `json:"token"`
}

// Regiser User godoc
//
//	@Summary		Register a new user
//	@Description	Creates a user account and sends an invitation email.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"Registration details"
//	@Success		201		{object}	store.User
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var userPayload = ctx.Value(BodyKey).(RegisterUserPayload)

	hashedPassword, err := HashPassword(userPayload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	role, err := app.store.Roles.GetByName(ctx, "user")

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user := store.User{
		UserName: userPayload.Username,
		Email:    userPayload.Email,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err = app.store.Users.CreateAndInvite(
		ctx,
		&user,
		hashToken,
		app.config.mail.exp,
	)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	activationURL := fmt.Sprintf("http://%s/api/v1/user/activate/%s", r.Host, plainToken)

	vars := struct {
		UserName      string
		ActivationURL string
	}{
		UserName:      user.UserName,
		ActivationURL: activationURL,
	}

	isProdEnv := app.config.env == "production"

	if err = app.mailer.Send(
		mailer.UserWelcomeTemplate,
		user.UserName,
		user.Email,
		vars,
		!isProdEnv,
	); err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type CreateTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload := ctx.Value(BodyKey).(CreateTokenPayload)

	user, err := app.store.Users.GetByEmail(ctx, payload.Email)

	if err != nil {
		switch err {
		case store.NotFoundError:
			app.unAuthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		app.unAuthorizedError(w, r, fmt.Errorf("invalid credentials"))
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.iss,
		"aud": app.config.auth.iss,
	}
	token, err := app.authenticator.GenerateToken(claims)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, map[string]any{
		"token": token,
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
