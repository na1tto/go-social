package main

import (
	"net/http"

	repository "github.com/na1tto/go-social/internal/store"
)

type RegistredUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// registerUserHandler godoc
//
// @Summary 							Registers a user
// @Description 					Registers a user
// @Tags									authentication
// @Accept 								json
// @Produce								json
// @Param	payload body		RegisterUserPayload	true	"User credentials"
// @Success	201	{object}	repository.User	"User
// @Failure 400	{object} 	error
// @Failure 500 {object} 	error
// @Router 								/authentication/user [post]
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

	// create the user in the db and store an invitation token to validate the invite email later
	if err := app.store.Users.CreateAndInvite(ctx, user, "uuidv4"); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// mail

	if err := writeJson(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
