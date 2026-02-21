package main

import (
	"net/http"
)

func (a *application) GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	userId := "1"
	rcontext := r.Context()

	feed, err := a.store.Posts.GetUserFeed(rcontext, userId)
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
