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

func (g *Group) DeleteGroup(db *gorm.DB, gid uint) (int64, error) {
	// Delete the group members
	err := db.Debug().Model(&Group{}).Where("id = ?", gid).Association("Members").Clear()
	if err != nil {
		return 0, err
	}

	// Delete the group
	db = db.Debug().Model(&Group{}).Where("id = ?", gid).Take(&Group{}).Delete(&Group{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (g *Group) UpdateGroupMembers(db *gorm.DB) (*Group, error) {
	if len(g.MemberIDs) > 0 {
		user := User{}
		users, err := user.FindUsersByIDs(db, g.MemberIDs)
		if err != nil {
			return &Group{}, err
		}
		g.Members = users
	}

	err := db.Debug().Save(&g).Error
	if err != nil {
		return &Group{}, err
	}
	return g, nil
}
func GenGroupToken(pid uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["group_id"] = pid
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func (g *Group) InviteGroupByToken(t string, uid uint, db *gorm.DB) (*Group, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return &Group{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	gid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["group_id"]), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid group_id in token: %v", err)
	}

	g.ID = uint(gid)
	g.MemberIDs = append(g.MemberIDs, uid)
	
	if err := g.Validate("update_member"); err != nil {
		return nil, err
	}

	updateGroup, err := g.UpdateGroupMembers(db)
	if err != nil {
		return nil, err
	}
	return updateGroup, nil
	}

	return &Group{}, errors.New("invalid token")
}