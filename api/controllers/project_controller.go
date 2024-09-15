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
)

func (server *Server) CreateProject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	project := models.Project{}
	err = json.Unmarshal(body, &project)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = project.Validate("create")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	project.AuthorID = uint(tokenID)

	// Save the project and associate members
	projectCreated, err := project.SaveProject(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, projectCreated.ID))
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusCreated,
		Message: "Project successfully created",
		Data:    nil,
	})
}

func (server *Server) UpdateProjectMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	project := models.Project{}
	err = json.Unmarshal(body, &project)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Validate the request
	err = project.Validate("update_member")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = server.DB.Debug().Preload("Members").Where("id = ?", pid).Take(&project).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("project not found"))
		return
	}

	// Extract the token ID from the request
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// Check if the authenticated user is the author of the project
	if tokenID != project.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("you are not authorized to update this members of project"))
		return
	}

	projectMemberUpdated, err := project.UpdateProjectMembers(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, projectMemberUpdated.ID))
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Project members updated successfully",
		Data:    nil,
	})
}

func (server *Server) DeleteProjectMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	project := models.Project{}
	if err = json.Unmarshal(body, &project); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Validate member IDs
	if err = project.Validate("update_member"); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Find the project and preload its members
	if err = server.DB.Debug().Preload("Members").Where("id = ?", pid).Take(&project).Error; err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("project not found"))
		return
	}

	// Extract the token ID from the request
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// Check if the authenticated user is the author of the project
	if tokenID != project.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("you are not authorized to delete members from this project"))
		return
	}

	// Remove members from the project
	for _, memberID := range project.MemberIDs {
		err = server.DB.Debug().Model(&project).Association("Members").Delete(&models.User{ID: memberID})
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Project members deleted successfully",
	})
}

// DeleteProject deletes a project by its ID
func (server *Server) DeleteProject(w http.ResponseWriter, r *http.Request) {
	// Get the project ID from the URL parameters
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("invalid project ID"))
		return
	}

	// Extract the token ID from the request
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// Check if the token ID matches the author of the project
	project := models.Project{}
	err = server.DB.Debug().Model(models.Project{}).Where("id = ?", pid).Take(&project).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("project not found"))
		return
	}

	err = project.Validate("delete")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Check if the authenticated user is the author of the project
	if tokenID != project.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("you are not authorized to delete this project"))
		return
	}

	// Delete the project
	_, err = project.DeleteProject(server.DB, uint32(pid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Project successfully deleted",
	})
}

// GetProjects retrieves all projects where the user is the author or a member
func (server *Server) GetProjects(w http.ResponseWriter, r *http.Request) {
	var projects []models.Project

	// Extract the token ID from the request
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// Query to fetch projects where the user is either the author or a member
	err = server.DB.Debug().Model(&models.Project{}).
		Preload("Members").
		Where("author_id = ? OR id IN (SELECT project_id FROM project_users WHERE user_id = ?)", tokenID, tokenID).
		Order("created_at DESC").
		Find(&projects).Error

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	projectsRespnse := helper.TransformProjects(projects)

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Projects retrieved successfully",
		Data:    projectsRespnse,
	})
}

// GetProjectByID retrieves a project by its ID if the user is the author or a member
func (server *Server) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	// Extract the project ID from the URL parameters
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("invalid project ID"))
		return
	}

	// Extract the token ID from the request
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	var project models.Project

	// Query to fetch the project by ID with validation
	err = server.DB.Debug().Model(&models.Project{}).
		Preload("Members").
		Where("id = ? AND (author_id = ? OR id IN (SELECT project_id FROM project_users WHERE user_id = ?))", pid, tokenID, tokenID).
		First(&project).Error

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("project not found or you are not authorized to view it"))
		return
	}

	projectResponse := helper.TransformProject(project)

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Project retrieved successfully",
		Data:    projectResponse,
	})
}

func (server *Server) GenerateInviteToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("invalid project ID"))
		return
	}

	// Generate the token using the project ID
	token, err := models.GenProjectToken(uint32(pid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Return the generated token in the response
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Token generated successfully",
		Data:    token,
	})
}

func (server *Server) ValidateInvitationToken(w http.ResponseWriter, r *http.Request) {
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

	project := models.Project{}
	// invite user to project
	_, err = project.InviteProjectByToken(tokenValue, uint(uid), server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Token validated successfully",
	})
}
