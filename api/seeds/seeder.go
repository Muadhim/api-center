package seeds

import (
	"api-center/api/models"
	"log"
	"time"

	"gorm.io/gorm"
)

func Load(db *gorm.DB) {
	if db == nil {
		log.Fatal("db instance is nil")
	}

	// Drop table if exists
	err := db.Migrator().DropTable(&models.User{})
	if err != nil {
		log.Fatalf("cannot drop table %v", err)
	}

	// Auto-migrate User model
	err = db.Debug().AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("cannot migrate table %v", err)
	}

	// Seed Users
	for _, user := range Users {
		err := db.Debug().Create(&user).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

	// Drop table if exists
	err = db.Migrator().DropTable(&models.Project{})
	if err != nil {
		log.Fatalf("cannot drop table %v", err)
	}

	// Auto-migrate User model
	err = db.Debug().AutoMigrate(&models.Project{})
	if err != nil {
		log.Fatalf("cannot migrate table %v", err)
	}

	var users []models.User
	db.Find(&users)

	if len(users) < 4 {
		log.Fatalf("Not enough users in the database to seed projects. Need at least 4 users.")
		return
	}

	projects := []models.Project{
		{
			Name: "Project Alpha",
			Members: []models.User{
				users[0], // Assigning existing users
				users[1],
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name: "Project Beta",
			Members: []models.User{
				users[2], // Assigning existing users
				users[3],
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, project := range projects {
		err := db.Create(&project).Error
		if err != nil {
			log.Fatalf("Failed to seed projects: %v", err)
		}
	}

}
