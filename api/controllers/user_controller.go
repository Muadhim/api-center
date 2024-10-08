package controllers

import (
	"api-center/api/auth"
	"api-center/api/helper"
	"api-center/api/models"
	"api-center/api/responses"
	"api-center/utils/formaterror"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// check if user is exist by email
	userExist, err := user.FindUserByEmail(server.DB, user.Email)

	if userExist != nil && err != gorm.ErrRecordNotFound {
		responses.ERROR(w, http.StatusConflict,
			fmt.Errorf("user with email %s already exist", user.Email))
		return
	}

	userCreated, err := user.SaveUser(server.DB)

	if err != nil {
		formatedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusCreated,
		Message: "User successfully created",
		Data:    nil,
	})
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	users, err := user.FindAllUser(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	usersResponse := helper.TransformUsers(*users)
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Success Get All Users",
		Data:    usersResponse,
	})
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userResponse := helper.TransformUser(*userGotten)
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Success Get User",
		Data:    userResponse,
	})
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadGateway, err)
		return
	}

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

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != uint(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	userResponse := helper.TransformUser(*updatedUser)

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Success Update User",
		Data:    userResponse,
	})
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := models.User{}
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	if tokenID != 0 && tokenID != uint(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusNoContent,
		Message: "Success Delete User",
		Data:    nil,
	})
}

func (server *Server) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Otp         string `json:"otp"`
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := models.User{}
	user.Email = request.Email
	user.Password = request.NewPassword
	if err := user.Validate("forgot_password"); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	otp := models.OtpStore{}
	otp.Otp = request.Otp
	otp.Email = request.Email

	if otp.Otp == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("otp is required"))
		return
	}

	if otp.Email == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("email is required"))
		return
	}
	err := otp.ValidateOtp(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = user.ChangePassword(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if err := otp.DeleteOtp(server.DB); err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Password changed successfully",
	})
}
