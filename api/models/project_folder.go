package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProjectFolder struct {
	ID        uint      `gorm:"primary_key;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	ProjectID uint      `gorm:"index" json:"project_id"`
	ParentID  *uint     `gorm:"index" json:"parent_id"`
	AuthorID  uint      `gorm:"index;not null" json:"author_id"`
	UpdateBy  uint      `gorm:"index;default:null" json:"update_by"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (pf *ProjectFolder) Validate(action string) error {
	action = strings.ToLower(action)

	switch action {
	case "create":
		if pf.Name == "" {
			return errors.New("project folder name is required")
		}
		if pf.ProjectID == 0 {
			return errors.New("project ID is required")
		}
	case "update":
		if pf.ID == 0 {
			return errors.New("project folder ID is required for update")
		}
	default:
		return errors.New("invalid action specified")
	}

	return nil
}

func (pf *ProjectFolder) SaveProjectFolder(db *gorm.DB) (*ProjectFolder, error) {
	err := db.Debug().Create(&pf).Error
	if err != nil {
		return &ProjectFolder{}, err
	}
	return pf, nil
}

func (pf *ProjectFolder) UpdateProjectFolder(db *gorm.DB) (*ProjectFolder, error) {
	err := db.Debug().Model(&ProjectFolder{}).
		Where("id = ?", pf.ID).
		Updates(map[string]interface{}{
			"name":       pf.Name,
			"parent_id":  pf.ParentID,
			"update_by":  pf.UpdateBy,
			"updated_at": time.Now()}).Error
	if err != nil {
		return &ProjectFolder{}, err
	}
	return pf, nil
}

func (pf *ProjectFolder) DeleteProjectFolder(db *gorm.DB) (int64, error) {
	var deleteFolder func(folderID uint) (int64, error)

	deleteFolder = func(folderID uint) (int64, error) {
		var folders []ProjectFolder

		// Find all subfolders of the current folder
		err := db.Debug().Where("parent_id = ?", folderID).Find(&folders).Error
		if err != nil {
			return 0, err
		}

		// Initialize a counter for total deleted items
		totalDeleted := int64(0)

		// Recursively delete each subfolder
		for _, folder := range folders {
			count, err := deleteFolder(folder.ID)
			if err != nil {
				return 0, err
			}
			totalDeleted += count
		}

		// Delete APIs associated with the current folder
		result := db.Debug().Where("folder_id = ?", folderID).Delete(&ProjectApi{})
		if result.Error != nil {
			return 0, result.Error
		}
		totalDeleted += result.RowsAffected

		// Delete the current folder
		result = db.Debug().Where("id = ?", folderID).Delete(&ProjectFolder{})
		if result.Error != nil {
			return 0, result.Error
		}
		totalDeleted += result.RowsAffected

		return totalDeleted, nil
	}

	totalDeleted, err := deleteFolder(pf.ID)
	if err != nil {
		return 0, err
	}

	return totalDeleted, nil
}

func (pf *ProjectFolder) GetProjectFolders(db *gorm.DB) (*ProjectFolder, error) {
	err := db.Debug().Model(&ProjectFolder{}).Where("project_id = ?", pf.ProjectID).Find(&pf).Error

	if err != nil {
		return &ProjectFolder{}, err
	}

	return pf, nil
}
