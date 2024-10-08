package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string  `gorm:"uniqueIndex;type:varchar(100);not null" json:"username"`
	Password  string  `json:"-" gorm:"column:password;not null"`
	Email     string  `gorm:"uniqueIndex;type:varchar(100)" json:"email"`
	ProjectID *uint   `json:"project_id"`
	Project   Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}
