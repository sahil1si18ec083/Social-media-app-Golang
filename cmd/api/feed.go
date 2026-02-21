package main

import (
	"net/http"

	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
)

// GetUserFeedHandler godoc
//
//	@Summary		Get user feed
//	@Description	Returns feed posts for the authenticated user including posts from followed users
//	@Tags			users
//	@ID				get-user-feed
//	@Produce		json
//
//	@Param			limit	query		int		false	"Number of posts to return (default 20)"
//	@Param			offset	query		int		false	"Pagination offset"
//	@Param			sort	query		string	false	"Sort order: asc or desc"
//	@Param			tags	query		[]string	false	"Filter by tags (repeat parameter)"
//	@Param			search	query		string	false	"Search text in title or content"
//
//	@Success		200		{array}		store.PostWithMetadata	"Feed returned successfully"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//
//	@Router			/users/feed [get]
func (a *application) GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	userId := "1"
	rcontext := r.Context()
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
		Search: "",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	feed, err := a.store.Posts.GetUserFeed(rcontext, userId, fq)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}
	err = writeJSON(w, http.StatusOK, feed)
	if err != nil {

		a.internalServerError(w, r, err)
		return
	}

}
