package seeds

import "api-center/api/models"

var Users = []models.User{
	{
		Name:     "ucub",
		Email:    "ucub@gmail.com",
		Password: "ucubPassword", // Ensure to hash password if needed
	},
}
