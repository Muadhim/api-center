package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID        uint      `gorm:"primary_key;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	AuthorID  uint      `gorm:"index;not null" json:"author_id"`
	Members   []User    `gorm:"many2many:project_users;" json:"members"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Project) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if p.Name == "" {
			return errors.New("Required project name")
		}
		return nil
	case "update":
		return nil
	case "delete":
		return nil
	default:
		return nil
	}
}

func (p *Project) SaveProject(db *gorm.DB) (*Project, error) {
	err := db.Debug().Create(&p).Error
	if err != nil {
		return &Project{}, err
	}
	return p, err
}

func (p *Project) DeleteProject(db *gorm.DB, pid uint32) (int64, error) {
	db = db.Debug().Model(&Project{}).Where("id = ?", pid).Take(&Project{}).Delete(&Project{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
