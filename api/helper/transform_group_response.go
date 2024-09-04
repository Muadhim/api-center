package helper

import (
	"api-center/api/models"
	"api-center/api/responses"
	"time"
)

func TransformGroup(group models.Group) responses.Group {
	memberIDs := make([]uint, len(group.Members))
	for i, member := range group.Members {
		memberIDs[i] = member.ID
	}

	members := TransformUsers(group.Members)

	return responses.Group{
		ID:        group.ID,
		Name:      group.Name,
		AuthorID:  group.AuthorID,
		Members:   members,
		CreatedAt: group.CreatedAt.Format(time.RFC3339),
		UpdatedAt: group.UpdatedAt.Format(time.RFC3339),
	}
}

func TransformGroups(groups []models.Group) []responses.Group {
	var groupsResp []responses.Group
	for _, group := range groups {
		groupsResp = append(groupsResp, TransformGroup(group))
	}
	return groupsResp
}
