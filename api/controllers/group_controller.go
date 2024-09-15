package controllers

import (
	"api-center/api/auth"
	"api-center/api/helper"
	"api-center/api/models"
	"api-center/api/responses"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) CreateGroup(w http.ResponseWriter, r *http.Request) {
	type groupRequest struct {
		Name string `json:"name"`
	}
	var req groupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	group := models.Group{}
	group.Name = req.Name

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	group.AuthorID = uint(tokenID)

	groupCreated, err := group.SaveGroup(s.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, groupCreated.ID))
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusCreated,
		Message: "Group successfully created",
		Data:    nil,
	})
}

func (s *Server) GetGroups(w http.ResponseWriter, r *http.Request) {
	var groups []models.Group

	// Extrac the token ID from the request
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	err = s.DB.Debug().Model(&models.Group{}).
		Preload("Members").
		Where("author_id = ? OR id IN (SELECT group_id FROM group_users WHERE user_id = ?)", tokenID, tokenID).
		Order("created_at DESC").
		Find(&groups).Error

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	gropsResponse := helper.TransformGroups(groups)

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "OK",
		Data:    gropsResponse,
	})
}

func (s *Server) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pid, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	group := models.Group{}
	err = s.DB.Debug().Model(&models.Group{}).Where("id = ? AND author_id = ?", pid, tokenID).Take(&group).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = group.DeleteGroup(s.DB, uint(pid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Group successfully deleted",
	})
}

func (s Server) GenInviteGroupToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	token, err := models.GenGroupToken(uint32(pid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Token generated successfully",
		Data:    token,
	})
}
func (s *Server) ValidateInvitationGroupToken(w http.ResponseWriter, r *http.Request) {
	type tokenRequest struct {
		Token string `json:"token"`
	}

	var req tokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenValue := req.Token
	uid, err := auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	group := models.Group{}
	_, err = group.InviteGroupByToken(tokenValue, uint(uid), s.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Group successfully joined",
	})
}
