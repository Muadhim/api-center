package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProjectApi struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	Name            string    `gorm:"type:varchar(255);not null" json:"name"`
	FolderID        uint      `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"folder_id"`
	AuthorID        uint      `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"author_id"`
	UpdateBy        uint      `gorm:"index;default:null" json:"update_by"`
	Method          string    `gorm:"type:varchar(255);not null" json:"method"`
	Path            string    `gorm:"type:varchar(255)" json:"path"`
	Header          string    `gorm:"type:text" json:"header"`
	Body            string    `gorm:"type:text" json:"body"`
	ExampleRequest  string    `gorm:"type:text" json:"example_request"`
	ExampleResponse string    `gorm:"type:text" json:"example_response"`
	ProjectID       uint      `gorm:"-" json:"project_id"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
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
		if pa.Method == "" {
			return errors.New("project api method is required")
		}
	case "update":
		if pa.ID == 0 {
			return errors.New("project api ID is required for update")
		}
		if pa.FolderID == 0 {
			return errors.New("project folder ID is required")
		}
		if pa.Method == "" {
			return errors.New("project api method is required")
		}
		if pa.ProjectID == 0 {
			return errors.New("project ID is required")
		}
	case "update-detail":
		if pa.ID == 0 {
			return errors.New("project api ID is required for update")
		}
		if pa.FolderID == 0 {
			return errors.New("project folder ID is required")
		}
		if pa.Path == "" {
			return errors.New("project api path is required")
		}
		if pa.Method == "" {
			return errors.New("project api method is required")
		}
		if pa.Header == "" {
			return errors.New("project api header is required")
		}
		if pa.Body == "" {
			return errors.New("project api body is required")
		}
		if pa.ExampleRequest == "" {
			return errors.New("project api example request is required")
		}
		if pa.ProjectID == 0 {
			return errors.New("project ID is required")
		}
		if pa.ExampleResponse == "" {
			return errors.New("project api example response is required")
		}
	default:
		return errors.New("invalid action specified")
	}

	return nil
}

func (pa *ProjectApi) SaveProjectApi(db *gorm.DB) (*ProjectApi, error) {
	err := pa.checkUserInProject(db, pa.UpdateBy)
	if err != nil {
		return &ProjectApi{}, err
	}

	err = db.Debug().Create(&pa).Error
	if err != nil {
		return &ProjectApi{}, err
	}
	return pa, nil
}

func (pa *ProjectApi) UpdateProjectApi(db *gorm.DB) (*ProjectApi, error) {
	err := pa.checkUserInProject(db, pa.UpdateBy)
	if err != nil {
		return &ProjectApi{}, err
	}

	err = db.Debug().Model(&ProjectApi{}).Where("id = ?", pa.ID).
		Updates(map[string]interface{}{
			"name":       pa.Name,
			"folder_id":  pa.FolderID,
			"method":     pa.Method,
			"update_by":  pa.UpdateBy,
			"updated_at": time.Now()}).
		Error
	if err != nil {
		return &ProjectApi{}, err
	}
	return pa, nil
}

func (pa *ProjectApi) UpdateDetailApi(db *gorm.DB, uid uint) (*ProjectApi, error) {

	err := pa.checkUserInProject(db, uid)
	if err != nil {
		return &ProjectApi{}, err
	}

	err = db.Debug().Model(&ProjectApi{}).Where("id = ?", pa.ID).
		Updates(map[string]interface{}{
			"method":           pa.Method,
			"path":             pa.Path,
			"header":           pa.Header,
			"body":             pa.Body,
			"example_request":  pa.ExampleRequest,
			"example_response": pa.ExampleResponse,
			"update_by":        pa.UpdateBy,
			"updated_at":       time.Now()}).
		Error

	if err != nil {
		return &ProjectApi{}, err
	}
	return pa, nil
}

func (pa *ProjectApi) DeleteProjectApi(db *gorm.DB) (int64, error) {
	err := pa.checkUserInProject(db, pa.UpdateBy)
	if err != nil {
		return 0, err
	}

	db = db.Debug().Model(&ProjectApi{}).Where("id = ?", pa.ID).
		Take(&ProjectApi{}).
		Delete(&ProjectApi{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}

func (pa *ProjectApi) GetProjectApiDetail(db *gorm.DB, uid uint) (*ProjectApi, error) {
	err := pa.checkUserInProject(db, uid)
	if err != nil {
		return &ProjectApi{}, err
	}

	err = db.Debug().Model(&ProjectApi{}).Where("id = ?", pa.ID).First(&pa).Error
	if err != nil {
		return &ProjectApi{}, err
	}

	return pa, nil
}

func (pa *ProjectApi) checkUserInProject(db *gorm.DB, uid uint) error {
	p := Project{}
	err := db.Debug().Model(&Project{}).Where("id = ?", pa.ProjectID).
		Preload("Members").
		Take(&p).Error
	if err != nil {
		return err
	}

	if p.AuthorID != uid {
		isMember := false
		for _, m := range p.Members {
			if m.ID == uid {
				isMember = true
				break
			}
		}
		if !isMember {
			return errors.New("only project author or members can update")
		}
	}
	return nil
}
