package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID        uint      `gorm:"primary_key;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	AuthorID  uint      `gorm:"index;not null" json:"author_id"`
	Projects  []Project `gorm:"foreignkey:GroupID" json:"projects,omitempty"`
	Members   []User    `gorm:"many2many:group_users" json:"members,omitempty"`
	MemberIDs []uint    `gorm:"-" json:"member_ids"` // Use this field to capture the member IDs from the request\
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (g *Group) Validate(action string) error {
	action = strings.ToLower(action)
	switch action {
	case "create":
		if g.Name == "" {
			return errors.New("group name is required")
		}
	case "update_member":
		if len(g.MemberIDs) == 0 {
			return errors.New("at least one member must be provided for update")
		}
	case "delete":
		if g.ID == 0 {
			return errors.New("group ID is required for deletion")
		}
	default:
		return errors.New("invalid action specified")
	}
	return nil
}

// SaveGroup creates a new group and associates members if provided
func (g *Group) SaveGroup(db *gorm.DB) (*Group, error) {
	// Create the group with the associated members
	err := db.Debug().Create(&g).Error
	if err != nil {
		return &Group{}, err
	}
	return g, nil
}

func (g *Group) FindGroupByID(db *gorm.DB, gid uint) (*Group, error) {
	err := db.Debug().Model(&Group{}).Where("id = ?", gid).Take(&g).Error
	if err != nil {
		return &Group{}, err
	}
	if g.ID != gid {
		return &Group{}, errors.New("group not found")
	}
	return g, err
}
