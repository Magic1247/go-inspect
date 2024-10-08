package models

import "gorm.io/gorm"

type InspectionItem struct {
	gorm.Model
	Title           string            `gorm:"type:varchar(255);not null" json:"title"`
	Details         string            `gorm:"type:text" json:"details"`
	ExecutionMethod string            `gorm:"type:varchar(100);not null" json:"execution_method"`
	Points          []InspectionPoint `gorm:"many2many:point_items;" json:"points"`
}
