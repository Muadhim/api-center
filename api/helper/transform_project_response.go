package helper

import (
	"api-center/api/models"
	"api-center/api/responses"
	"time"
)

func TransformProject(project models.Project) responses.Project {
	memberIDs := make([]uint, len(project.Members))
	for i, member := range project.Members {
		memberIDs[i] = member.ID
	}

	members := TransformUsers(project.Members)

	return responses.Project{
		ID:        project.ID,
		Name:      project.Name,
		AuthorID:  project.AuthorID,
		Members:   members,
		CreatedAt: project.CreatedAt.Format(time.RFC3339),
		UpdatedAt: project.UpdatedAt.Format(time.RFC3339),
	}
}

func TransformProjects(projects []models.Project) []responses.Project {
	var projectsResp []responses.Project
	for _, project := range projects {
		projectsResp = append(projectsResp, TransformProject(project))
	}
	return projectsResp
}
