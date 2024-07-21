package seeds

import (
	"api-center/api/models"
	"log"

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
}
