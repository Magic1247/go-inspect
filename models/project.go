package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	ParentID    *uint     `json:"parent_id"`
	Parent      *Project  `gorm:"foreignkey:ParentID" json:"parent,omitempty"`
	Children    []Project `gorm:"foreignkey:ParentID" json:"children,omitempty"`
}
