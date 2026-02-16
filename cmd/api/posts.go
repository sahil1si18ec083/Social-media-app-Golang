package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (a *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// userId := 1
	fmt.Println(".....................................")
	var payload CreatePostPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		fmt.Println(err)
		a.badRequestResponse(w, r, err)
		return
	}
	err = Validate.Struct(&payload)
	if err != nil {
		fmt.Println(err)
		a.badRequestResponse(w, r, err)
		return
	}
	post := &store.Post{Title: payload.Title, Content: payload.Content,
		Tags: payload.Tags,
		// change userID after auth
		UserID: 1}
	fmt.Println(post)
	rcontext := r.Context()
	err = a.store.Posts.Create(rcontext, post)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}
	err = writeJSON(w, http.StatusCreated, post)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}
	// app.store.Posts.Create(r.Context(), &store.Post{})
}

func (a *application) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	fmt.Println(postID)
	rcontext := r.Context()

	post, err := a.store.Posts.GetById(rcontext, postID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			a.notFoundResponse(w, r, err)
			return
		}

		a.internalServerError(w, r, err)
		return
	}
	comments, err := a.store.Comments.GetByPostId(rcontext, postID)
	fmt.Println(comments)
	if err != nil {
		// for future
		fmt.Println(err)
		return
	}
	post.Comment = *comments
	err = writeJSON(w, http.StatusOK, post)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}

}
