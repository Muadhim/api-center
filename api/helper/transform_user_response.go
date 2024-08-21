package helper

import (
	"api-center/api/models"
	"api-center/api/responses"
)

func TransformUser(user models.User) responses.User {
	return responses.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func TransformUsers(users []models.User) []responses.User {
	if users == nil {
		return []responses.User{}
	}

	var userResponses []responses.User
	for _, user := range users {
		userResponses = append(userResponses, TransformUser(user))
	}
	return userResponses
}
