package controllers

import (
	"api-center/api/auth"
	"api-center/api/models"
	"api-center/api/responses"
	"api-center/utils/formaterror"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, user, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Successfully logged in",
		Data: struct {
			Id          uint   `json:"id"`
			Name        string `json:"name"`
			Email       string `json:"email"`
			AccessToken string `json:"access_token"`
		}{
			Id:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			AccessToken: token,
		},
	})
}

func (server *Server) SignIn(email, password string) (string, models.User, error) {
	var err error
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", user, err
	}

	err = models.VerifyPassword(user.Password, password)
	if err != nil || err == bcrypt.ErrMismatchedHashAndPassword {
		return "", user, err
	}

	token, err := auth.CreateToken(uint32(user.ID))
	return token, user, err
}
