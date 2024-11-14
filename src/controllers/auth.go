package controllers

import (
	"blogger/database"
	"blogger/src/authentication"
	"blogger/src/models"
	"blogger/src/repositories"
	"blogger/src/responses"
	"blogger/src/security"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Login returns a token for the requesting user.
func Login(w http.ResponseWriter, r *http.Request) {
	bodyReq, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyReq, &user); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.UsersRepository(db)
	userFromDB, err := repository.GetByEmail(user.Email)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.CompareHashes(user.Password, userFromDB.Password); err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	userToken, err := authentication.GenerateToken(uint64(userFromDB.ID))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: "Congrats, you're logged in",
		Token:   userToken,
	})
}

// ChangePassword allows the user to change their password.
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := authentication.UserIDExtraction(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	bodyReq, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	var password models.Password

	if err = json.Unmarshal(bodyReq, &password); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.AuthenticationRepository(db)
	passwordFromDB, err := repository.GetPassword(userIDFromToken)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.CompareHashes(password.Current, passwordFromDB); err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("the given 'current password' don't match with our DB"))
		return
	}

	newPass, err := security.Hash(password.New)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = repository.ChangePassword(userIDFromToken, string(newPass)); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
