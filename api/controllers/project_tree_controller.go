package controllers

import (
	"api-center/api/auth"
	"api-center/api/models"
	"api-center/api/responses"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) GetProjectTree(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// check project

	project := models.Project{}
	err = project.FindProjectByID(s.DB, uint(pid), uint(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// get project tree
	projectTree := models.ProjectTree{}
	tree, err := projectTree.GetProjectTree(s.DB, uint(project.ID))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "Project tree successfully retrieved",
		Data: struct {
			ProjectName string                `json:"project_name"`
			Tree        []*models.ProjectTree `json:"project_tree"`
		}{
			ProjectName: project.Name,
			Tree:        tree,
		},
	})
}
