package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	repository "github.com/na1tto/go-social/internal/store"
)

type postKey string

const postCtx postKey = "post"

// these are the values the we will accept from the user
type createPostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

// CreatePost godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post in the database
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		repository.Post	true	"Post Infos"
//	@Success		201			{object}	repository.Post	"Success"
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/ [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)

	post := &repository.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserId:  user.ID,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// GetPost godoc
//
//	@Summary		Retrieves a Post
//	@Description	Retrieves a post by ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		int				true	"Post ID"
//	@Success		201		{object}	repository.Post	"Success"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error	"Post not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postId} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	// adding the comments to the post before returning to the user
	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// DeletePost godoc
//
//	@Summary		Deletes a Post
//	@Description	Deletes a post by ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path	int	true	"Post ID"
//	@Success		204		"Success"
//	@Failure		404		{object}	error	"Post not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postId} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	err = app.store.Posts.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent) // we don't want to return any data via json here
}

// these are the data that the user can send to update their posts
type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Fields to Update"
//	@Success		200		{object}	repository.Post		"Success"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error	"Post not found"
//	@Failure		409		{object}	error	"Version conflict"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postId} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		switch {
		case errors.Is(err, repository.StatusConflict):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

type commentPayload struct {
	Content string `json:"content" validate:"required,max=300"`
}

// UpdatePost godoc
//
//	@Summary		Creates a Comment
//	@Description	Comments on a post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		int					true	"Post ID"
//	@Param			payload	body		commentPayload		true	"Comment Field"
//	@Success		201		{object}	repository.Comment	"Success"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postId} [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {

	post := getPostFromCtx(r)

	var payload commentPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &repository.Comment{
		Content: payload.Content,
		UserId:  40,
		PostId:  post.ID,
	}

	ctx := r.Context()

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// middleware for fetching a post in the database
func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postId")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err) //handling errors internally
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *repository.Post {
	post, _ := r.Context().Value(postCtx).(*repository.Post)
	return post
}
