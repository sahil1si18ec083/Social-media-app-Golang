package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"
)

type contextKey string

const postContextKey = contextKey("post")

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}
type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}
type DeletePostResponse struct {
	Message string `json:"message"`
}

func (a *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreatePostPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	err = Validate.Struct(&payload)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	post := &store.Post{Title: payload.Title, Content: payload.Content,
		Tags: payload.Tags,
		// change userID after auth
		UserID: 1}

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

	rcontext := r.Context()

	post, ok := getPostFromContext(rcontext)
	if !ok {
		a.internalServerError(w, r, errors.New("post missing from context"))
		return // VERY IMPORTANT
	}
	comments, err := a.store.Comments.GetByPostId(rcontext, postID)
	if err != nil {
		a.internalServerError(w, r, err)
		return

	}
	post.Comment = *comments
	err = writeJSON(w, http.StatusOK, post)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}

}

func (a *application) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")

	rcontext := r.Context()

	err := a.store.Posts.Delete(rcontext, postID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			a.notFoundResponse(w, r, err)
			return
		}

		a.internalServerError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, DeletePostResponse{Message: "Post Deleted Successfully"})
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}

}

func (a *application) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {

	// PATCH CALL FOR UPDATING A POST

	postID := chi.URLParam(r, "postID")
	rcontext := r.Context()

	post, ok := getPostFromContext(rcontext)
	if !ok {
		a.internalServerError(w, r, errors.New("post missing from context"))
		return // VERY IMPORTANT
	}
	var payload UpdatePostPayload

	err := readJSON(w, r, &payload)
	if err != nil {

		if !errors.Is(err, io.EOF) {
			a.badRequestResponse(w, r, err)
			return
		}

	}
	err = Validate.Struct(&payload)
	if err != nil {

		a.badRequestResponse(w, r, err)
		return
	}
	if payload.Title == nil && payload.Content == nil {

		err = writeJSON(w, http.StatusNoContent, post)
		if err != nil {
			a.internalServerError(w, r, err)
			return
		}
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content

	}
	if payload.Title != nil {
		post.Title = *payload.Title

	}

	err = a.store.Posts.Update(rcontext, post, postID)
	if err != nil {
		fmt.Println(err)
		a.internalServerError(w, r, err)
		return
	}
	err = writeJSON(w, http.StatusOK, post)
	if err != nil {
		fmt.Println(err, "2")
		a.internalServerError(w, r, err)
		return
	}
	fmt.Print("yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy")
}

func (a *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")
		post, err := a.store.Posts.GetById(r.Context(), postID)
		if err != nil {
			fmt.Println(err)
			if errors.Is(err, store.ErrNotFound) {
				a.notFoundResponse(w, r, err)
				return
			}

			a.internalServerError(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), postContextKey, post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func getPostFromContext(ctx context.Context) (*store.Post, bool) {

	post, ok := ctx.Value(postContextKey).(*store.Post)
	if post == nil {
		return nil, false
	}
	if !ok {
		return nil, false
	}
	return post, true

}
