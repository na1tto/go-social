package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	repository "github.com/na1tto/go-social/internal/store"
)

// these are the values the we will accept from the user
type createPostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload
	if err := readJson(w, r, &payload); err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &repository.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		//TODO: change after auth
		UserId: 1,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJson(w, http.StatusCreated, post); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			writeJsonError(w, http.StatusNotFound, err.Error())
		default:
			writeJsonError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	
	if err := writeJson(w, http.StatusOK, post); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
