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
	Members   []User    `gorm:"many2many:project_users" json:"members,omitempty"`
	MemberIDs []uint    `gorm:"-" json:"member_ids"` // Use this field to capture the member IDs from the request
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Project) Validate(action string) error {
	action = strings.ToLower(action)

	switch action {
	case "create":
		// Ensure Name is provided
		if p.Name == "" {
			return errors.New("project name is required")
		}
		// Optionally validate Members (if necessary)

	case "create_project":
		// Ensure Name is provided
		if p.Name == "" {
			return errors.New("project name is required")
		}
		return nil

	case "update_member":
		// Validate that Members are provided for update
		if len(p.MemberIDs) == 0 {
			return errors.New("at least one member must be provided for update")
		}
		return nil

	case "delete":
		// Validate that Project ID is provided for delete
		if p.ID == 0 {
			return errors.New("project ID is required for deletion")
		}
		return nil

	default:
		return errors.New("invalid action specified")
	}
	return nil
}

// SaveProject creates a new project and associates members if provided
func (p *Project) SaveProject(db *gorm.DB) (*Project, error) {
	// Fetch the users that correspond to the provided member IDs
	if len(p.MemberIDs) > 0 {
		user := User{}
		users, err := user.FindUsersByIDs(db, p.MemberIDs)
		if err != nil {
			return &Project{}, err
		}
		p.Members = users
	}

	// Create the project with the associated members
	err := db.Debug().Create(&p).Error
	if err != nil {
		return &Project{}, err
	}
	return p, nil
}

func (p *Project) DeleteProject(db *gorm.DB, pid uint32) (int64, error) {
	db = db.Debug().Model(&Project{}).Where("id = ?", pid).Take(&Project{}).Delete(&Project{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
