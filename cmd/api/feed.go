package main

import (
	"net/http"

	repository "github.com/na1tto/go-social/internal/store"
)

// getUserFeed godoc
//
//	@Summary		Fetches a feed for a user
//	@Description	Fetches the user feed with flags for filtering
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			search	query		string	false	"Search filter for the feed"
//	@Param			sort	query		string	false	"Field for sorting data"
//	@Param			limit	query		string	false	"Field for pagination limit"
//	@Param			offset	query		string	false	"Field for pagination offset"
//	@Param			tags	query		string	false	"Field for filtering by tags"
//	@Param			since	query		string	false	"Field for filtering until a date"
//	@Param			until	query		string	false	"Field for filtering since a date"
//	@Success		200		{array}		repository.Post
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := repository.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if fq.Tags == nil {
		fq.Tags = []string{}
	}

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(40), fq)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
