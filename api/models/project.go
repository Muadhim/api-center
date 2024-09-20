package models

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Project struct {
	ID        uint      `gorm:"primary_key;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	AuthorID  uint      `gorm:"index;not null" json:"author_id"`
	Members   []User    `gorm:"many2many:project_users" json:"members,omitempty"`
	MemberIDs []uint    `gorm:"-" json:"member_ids"` // Use this field to capture the member IDs from the request\
	GroupID   uint      `gorm:"index" json:"group_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (p *Project) Validate(action string) error {
	action = strings.ToLower(action)

	switch action {
	case "create":
		if p.Name == "" {
			return errors.New("project name is required")
		}
	case "update_member":
		if len(p.MemberIDs) == 0 {
			return errors.New("at least one member must be provided for update")
		}
	case "delete":
		if p.ID == 0 {
			return errors.New("project ID is required for deletion")
		}
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

func (p *Project) UpdateProjectMembers(db *gorm.DB) (*Project, error) {
	// Fetch the users that correspond to the provided member IDs
	if len(p.MemberIDs) > 0 {
		user := User{}
		users, err := user.FindUsersByIDs(db, p.MemberIDs)
		if err != nil {
			return &Project{}, err
		}
		p.Members = users
	}

	err := db.Debug().Save(&p).Error
	if err != nil {
		return &Project{}, err
	}

	return p, nil
}

func GenProjectToken(pid uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["project_id"] = pid
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func (p *Project) InviteProjectByToken(t string, uid uint, db *gorm.DB) (*Project, error) {
	// Parse the token
	token, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Parse project_id from token
		pid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["project_id"]), 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid project_id in token: %v", err)
		}

		// Set project ID
		p.ID = uint(pid)

		// Add the user to MemberIDs
		p.MemberIDs = append(p.MemberIDs, uid)

		// Validate update action
		if err := p.Validate("update_member"); err != nil {
			return nil, err
		}

		// Update project members in the database
		updatedProject, err := p.UpdateProjectMembers(db)
		if err != nil {
			return nil, err
		}

		return updatedProject, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// FindProjectByID retrieves a project by its ID if the user is the author or a member
func (p *Project) FindProjectByID(db *gorm.DB, pid uint, uid uint) ( error) {
	err := db.Debug().Preload("Members").Where("id = ? AND (author_id = ? OR id IN (SELECT project_id FROM project_users WHERE user_id = ?))", pid, uid, uid).First(&p).Error
	if err != nil {
		return err
	}
	return  nil
}
