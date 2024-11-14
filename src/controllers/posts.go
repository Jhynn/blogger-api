package controllers

import (
	"blogger/database"
	"blogger/src/authentication"
	"blogger/src/models"
	"blogger/src/repositories"
	"blogger/src/responses"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

// IndexPost returns a listing of posts from the database.
func IndexPost(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := authentication.UserIDExtraction(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		db.Close()
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.PostRepository(db)
	var posts []models.Post

	posts, err = repository.ListingPost(userIDFromToken)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}{
		Data:    posts,
		Message: "that's your feed",
	})
}

// StorePost stores a new user in the database.
func StorePost(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := authentication.UserIDExtraction(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	bodyReq, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post

	if err = json.Unmarshal(bodyReq, &post); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	post.AuthorID = userIDFromToken
	if err = post.Prepare(); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		db.Close()
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	respository := repositories.PostRepository(db)
	post.ID, err = respository.CreatePost(post)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}{
		Data:    post,
		Message: "post succesfully created",
	})
}

// ShowPost retrieves the post in the database, considering its ID (from path param's URI).
func ShowPost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		db.Close()
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.PostRepository(db)
	var post models.Post

	post, err = repository.GetPost(postID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}{
		Data:    post,
		Message: "post was succesfully found",
	})
}

// UpdatePost updates the authenticated user in the database.
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := authentication.UserIDExtraction(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	postID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		db.Close()
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.PostRepository(db)
	postFromDB, err := repository.GetPost(postID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if postFromDB.AuthorID != userIDFromToken {
		responses.Error(w, http.StatusForbidden, errors.New("you can't update another user's post"))
		return
	}

	bodyReq, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post
	if err = json.Unmarshal(bodyReq, &post); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = post.Prepare(); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = repository.UpdatePost(postID, post); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// DeletePost deletes the authenticated user only, from the database.
func DeletePost(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := authentication.UserIDExtraction(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	postID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		db.Close()
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.PostRepository(db)
	postFromDB, err := repository.GetPost(postID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if postFromDB.AuthorID != userIDFromToken {
		responses.Error(w, http.StatusForbidden, errors.New("you can't delete another user's post"))
		return
	}

	if err = repository.DeletePost(postID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// ShowPost retrieves the post in the database, considering its ID (from path param's URI).
func ShowUserPosts(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		db.Close()
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.PostRepository(db)
	posts, err := repository.ListingUserPosts(userID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, struct {
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message"`
	}{
		Data:    posts,
		Message: "post was succesfully found",
	})
}

// LikePost increases the like property of the post.
func LikePost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		db.Close()
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.PostRepository(db)

	if err = repository.LikePost(postID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// UnlikePost decreases the like property of the post.
func UnlikePost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		db.Close()
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.PostRepository(db)

	if err = repository.UnlikePost(postID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
