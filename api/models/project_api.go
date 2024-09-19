package models

import "time"

type ProjectApi struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	FolderID  uint      `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"folder_id"`
	AuthorID  uint      `gorm:"index;not null;constraint:OnDelete:CASCADE" json:"author_id"`
	Method    string    `gorm:"type:varchar(255);not null" json:"method"`
	Path      string    `gorm:"type:varchar(255);not null" json:"path"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
