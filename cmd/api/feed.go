package main

import (
	"net/http"

	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
)

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
