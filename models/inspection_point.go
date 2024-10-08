package models

import "gorm.io/gorm"

type InspectionPoint struct {
	gorm.Model
	Name        string           `gorm:"type:varchar(100);not null" json:"name"`
	Description string           `gorm:"type:text" json:"description"`
	Location    string           `gorm:"type:varchar(255)" json:"location"`
	Items       []InspectionItem `gorm:"many2many:point_items;" json:"items"`
}
