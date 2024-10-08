package models

import "gorm.io/gorm"

type InspectionRoute struct {
	gorm.Model
	Name        string            `gorm:"type:varchar(100);not null" json:"name"`
	Description string            `gorm:"type:text" json:"description"`
	Points      []InspectionPoint `gorm:"many2many:route_points;" json:"points"`
	ProjectID   uint              `gorm:"not null" json:"project_id"`
	Project     Project           `gorm:"foreignkey:ProjectID" json:"project"`
}
