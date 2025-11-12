package main

import (
	"net/http"

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
