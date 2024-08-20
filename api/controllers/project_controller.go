package controllers

import (
	"api-center/api/auth"
	"api-center/api/models"
	"api-center/api/responses"
	"api-center/utils/formaterror"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	uid, err := strconv.ParseUint(fmt.Sprint(tokenID), 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("Error parsing token ID"))
		return
	}

	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	project.AuthorID = uint(uid)

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
		Data:    projectCreated,
	})
}
