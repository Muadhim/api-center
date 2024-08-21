package seeds

import "api-center/api/models"

var Users = []models.User{
	{
		Name:     "ucub",
		Email:    "ucub@api.com",
		Password: "$2a$10$t.Qtf4t3AJDnayTqM3e90O7OrA58eDs71JUA6Atm5P8nUw.Yq9C7S", // Ensure to hash password if needed
	},
}
