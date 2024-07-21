package database

import (
	"api-center/api/models"
	"log"

	"gorm.io/gorm"
)

func MigrageTable(db *gorm.DB) (err error) {
	err = db.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatal("failed to migrate table: ", err)
	}
	return nil
}
