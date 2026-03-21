package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	repository "github.com/na1tto/go-social/internal/store"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile via its ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"User ID"
//	@Success		200	{object}	repository.User	"success"
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type FollowUser struct {
	UserId int64 `json:"user_id"`
}

// FollowUser godoc
//
//	@Summary		Follows a User
//	@Description	Follows a user by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int		true	"User ID"
//	@Success		200		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromContext(r)

	// TODO: Reverse back to auth user id later
	var payload FollowUser
	if err := readJson(w, r, &payload); err != nil {
		switch err {
		case repository.ErrConflict:
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followedUser.ID, payload.UserId); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UnollowUser godoc
//
//	@Summary		Unfollows a User
//	@Description	Unollows a user by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int		true	"User ID"
//	@Success		200		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := getUserFromContext(r)

	// TODO: Reverse back to auth user id later
	var payload FollowUser
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, unfollowedUser.ID, payload.UserId); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// same middleware logic that we implemented in the posts handlers
func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			switch err {
			case repository.ErrNotFound:
				app.badRequestResponse(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *repository.User {
	user, _ := r.Context().Value(userCtx).(*repository.User)
	return user
}
