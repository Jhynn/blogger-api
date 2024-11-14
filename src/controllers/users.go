package controllers

import (
	"blogger/database"
	"blogger/src/authentication"
	"blogger/src/models"
	"blogger/src/repositories"
	"blogger/src/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Pong it's the response for /ping - meaning the API is online.
func Pong(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data := map[string]interface{}{
		"message": "pong",
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// IndexUser returns a listing of users from the database.
func IndexUser(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)

		return
	}
	defer db.Close()

	repository := repositories.UsersRepository(db)
	params := r.URL.Query()
	users, err := repository.Listing(params)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	page, per_page := repositories.PageAndPerPageValues(params)

	responses.JSON(w, http.StatusOK, struct {
		Data    interface{} `json:"data"`
		Page    uint64      `json:"page"`
		PerPage uint64      `json:"per_page"`
	}{
		Data:    users,
		Page:    page,
		PerPage: per_page,
	})
}

// StoreUser stores a new user in the database.
func StoreUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	requestBody, err := io.ReadAll(r.Body)

	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)

		return
	}

	var user models.User

	if err = json.Unmarshal(requestBody, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)

		return
	}

	if err = user.Prepare(models.STEP_CREATION); err != nil {
		responses.Error(w, http.StatusBadRequest, err)

		return
	}

	db, err := database.Connect()

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)

		return
	}

	repository := repositories.UsersRepository(db)
	userID, err := repository.Create(user)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)

		return
	}

	responses.JSON(w, http.StatusCreated, struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}{
		Data:    user,
		Message: fmt.Sprintf("User %d was succesfully created!", userID),
	})
}

// ShowUser retrieves the user in the database, considering its ID (from path param's URI).
func ShowUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)

	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)

		return
	}
	defer db.Close()

	repository := repositories.UsersRepository(db)
	user, err := repository.Get(userID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if user.ID == 0 {
		responses.Error(w, http.StatusBadRequest, errors.New("invalid ID, there is no user with such ID"))
		return
	}

	responses.JSON(w, http.StatusOK, struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}{
		Data:    user,
		Message: fmt.Sprintf("User %d was succesfully found!", userID),
	})
}

// UpdateUser updates the authenticated user in the database.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)

	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIDFromToken, err := authentication.UserIDExtraction(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userIDFromToken != userID {
		responses.Error(w, http.StatusForbidden, errors.New("you cannot update other user"))
		return
	}

	bodyReq, err := io.ReadAll(r.Body)

	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyReq, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare(models.STEP_UPDATE); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.UsersRepository(db)
	if err = repository.Update(userID, user); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, nil)
}

// DeleteUser deletes the authenticated user only, from the database.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIDFromToken, err := authentication.UserIDExtraction(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userIDFromToken != userID {
		responses.Error(w, http.StatusForbidden, errors.New("you cannot delete other user"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)

		return
	}
	defer db.Close()

	repository := repositories.UsersRepository(db)
	if err = repository.Delete(userID); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// Me returns the current authenticated user.
func Me(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := authentication.UserIDExtraction(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)

		return
	}
	defer db.Close()

	repository := repositories.UsersRepository(db)
	user, err := repository.Get(userIDFromToken)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if user.ID == 0 {
		responses.Error(w, http.StatusBadRequest, errors.New("invalid ID, there is no user with such ID"))
		return
	}

	responses.JSON(w, http.StatusOK, struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}{
		Data:    user,
		Message: "That's yourself!",
	})
}

// FollowUser the authenticared user will follow the given user (from path).
func FollowUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIDFromToken, err := authentication.UserIDExtraction(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userIDFromToken == userID {
		responses.Error(w, http.StatusForbidden, errors.New("you cannot follow yourself"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)

		return
	}
	defer db.Close()

	repository := repositories.UsersRepository(db)
	if err = repository.Follow(userID, userIDFromToken); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, nil)
}

// UnfollowUser the authenticared user will unfollow the given user (from path).
func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIDFromToken, err := authentication.UserIDExtraction(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userIDFromToken == userID {
		responses.Error(w, http.StatusForbidden, errors.New("you can't unfollow yourself, you are doomed to have your own company"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)

		return
	}
	defer db.Close()

	repository := repositories.UsersRepository(db)
	if err = repository.Unfollow(userID, userIDFromToken); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// GetFollowers retrieves all the followers of the given user.
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)

		return
	}
	defer db.Close()

	repository := repositories.UsersRepository(db)
	followers, err := repository.GetFollowers(userID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}{
		Data:    followers,
		Message: "That's your followers!",
	})
}

// GetFollowing retrieves all the users which the given user is following.
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)

		return
	}
	defer db.Close()

	repository := repositories.UsersRepository(db)
	following, err := repository.GetFollowing(userID)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}{
		Data:    following,
		Message: "That's your following!",
	})
}
