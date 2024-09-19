package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProjectApi struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	FolderID  uint      `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"folder_id"`
	AuthorID  uint      `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"author_id"`
	Method    string    `gorm:"type:varchar(255);not null" json:"method"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (pa *ProjectApi) Validate(action string) error {
	action = strings.ToLower(action)
	switch action {
	case "create":
		if pa.Name == "" {
			return errors.New("project api name is required")
		}
		if pa.FolderID == 0 {
			return errors.New("project folder ID is required")
		}
		if pa.AuthorID == 0 {
			return errors.New("project author ID is required")
		}
		if pa.Method == "" {
			return errors.New("project api method is required")
		}
	case "update":
		if pa.ID == 0 {
			return errors.New("project api ID is required for update")
		}
	default:
		return errors.New("invalid action specified")
	}

	return nil
}

func (pa *ProjectApi) SaveProjectApi(db *gorm.DB) (*ProjectApi, error) {
	err := db.Debug().Create(&pa).Error
	if err != nil {
		return &ProjectApi{}, err
	}
	return pa, nil
}

func (pa *ProjectApi) UpdateProjectApi(db *gorm.DB) (*ProjectApi, error) {
	err := db.Debug().Model(&ProjectApi{}).Where("id = ?", pa.ID).
		Updates(map[string]interface{}{
			"name":       pa.Name,
			"folder_id":  pa.FolderID,
			"method":     pa.Method,
			"updated_at": time.Now()}).
		Error
	if err != nil {
		return &ProjectApi{}, err
	}
	return pa, nil
}

func (pa *ProjectApi) DeleteProjectApi(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&ProjectApi{}).Where("id = ?", pa.ID).
		Take(&ProjectApi{}).
		Delete(&ProjectApi{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
