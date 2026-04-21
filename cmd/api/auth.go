package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/na1tto/go-social/internal/mailer"
	repository "github.com/na1tto/go-social/internal/store"
)

type RegistredUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*repository.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegistredUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken			"User
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegistredUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// pass the payload data to the user model to be stored
	user := &repository.User{
		UserName: payload.Username,
		Email:    payload.Email,
	}
	// hash the user password using bcrypt in the storage layer setting it to the model
	user.Password.Set(payload.Password)

	ctx := r.Context()

	// the plain text token will be displayed only for the user via email later
	plainToken := uuid.New().String()

	// store the hashed token in the db (testing a technique)
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	// create the user in the db and store an hashed invitation token to validate the invite email later
	if err := app.store.Users.CreateAndInvite(ctx, user, hashedToken, app.config.mail.exp); err != nil {
		switch err {
		case repository.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case repository.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	isProdEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.UserName,
		ActivationURL: activationURL,
	}

	// mail

	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.UserName, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// rollback user creation if email fails (SAGA patter)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent", "status code", status)

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload					body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200						{string}	string
//	@Failure		400						{object}	error
//	@Failure		401						{object}	error
//	@Failure		500						{object}	error
//	@Router			/authentication/token 	[post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// parse payload credentials
	var payload CreateUserTokenPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// fetch the user (check if the user exists) from the payload
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			app.unauthorizedErrorResponse(w, r, err) // do not return 404 to prevent enumaration attacks
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// send it to the client

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
