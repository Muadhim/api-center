package seeds

import "api-center/api/models"

var Users = []models.User{
	{
		Name:     "ucub victor",
		Email:    "ucub@gmail.com",
		Password: "ucubPassword", // Ensure to hash password if needed
	},
	{
		Name:     "abdol",
		Email:    "abdol@gmail.com",
		Password: "abdolPassword", // Ensure to hash password if needed
	},
}
