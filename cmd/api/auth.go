package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
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
	var userPayload = r.Context().Value(BodyKey).(RegisterUserPayload)

	hashedPassword, err := HashPassword(userPayload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user := store.User{
		UserName: userPayload.Username,
		Email:    userPayload.Email,
		Password: string(hashedPassword),
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err = app.store.Users.CreateAndInvite(
		r.Context(),
		&user,
		hashToken,
		app.config.mail.exp,
	)

	if err != nil {
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
