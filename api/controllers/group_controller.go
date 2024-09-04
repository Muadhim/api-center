package controllers

import (
	"api-center/api/auth"
	"api-center/api/helper"
	"api-center/api/models"
	"api-center/api/responses"
	"encoding/json"
	"fmt"
	"net/http"
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
	err = json.Unmarshal([]byte(req.Name), &group)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

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
