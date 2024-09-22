package controllers

import (
	"api-center/api/auth"
	"api-center/api/models"
	"api-center/api/responses"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) CreateProjectApi(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	projectApi := models.ProjectApi{}
	err = json.Unmarshal(body, &projectApi)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = projectApi.Validate("create")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	projectApi.AuthorID = uint(tokenID)

	projectApiCreated, err := projectApi.SaveProjectApi(s.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, projectApiCreated.ID))
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusCreated,
		Message: "Project api successfully created",
	})
}

func (s *Server) DeleteProjectApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId, err := strconv.ParseUint(vars["projectId"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	projectApi := models.ProjectApi{
		ProjectID: uint(projectId),
		ID:        uint(id),
	}

	_, err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	_, err = projectApi.DeleteProjectApi(s.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Project api successfully deleted",
	})
}

func (s *Server) UpdateProjectApi(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	projectApi := models.ProjectApi{}
	err = json.Unmarshal(body, &projectApi)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = projectApi.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	projectApi.UpdateBy = uint(uid)

	_, err = projectApi.UpdateProjectApi(s.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Project api successfully updated",
	})
}

func (s *Server) UpdateProjectApiDetail(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	projectApi := models.ProjectApi{}
	err = json.Unmarshal(body, &projectApi)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = projectApi.Validate("update-detail")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	projectApi.UpdateBy = uint(uid)

	_, err = projectApi.UpdateDetailApi(s.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Project api successfully updated",
	})
}

func (s *Server) GetProjectApiDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId, err := strconv.ParseUint(vars["projectId"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println("id", id)
	projectApi := models.ProjectApi{}
	projectApi.ID = uint(id)
	projectApi.ProjectID = uint(projectId)

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	result, err := projectApi.GetProjectApiDetail(s.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Project api successfully retrieved",
		Data:    result,
	})
}
