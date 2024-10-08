package models

import "gorm.io/gorm"

type CheckStatus string

const (
	CheckStatusPending CheckStatus = "pending"
	CheckStatusPassed  CheckStatus = "passed"
	CheckStatusFailed  CheckStatus = "failed"
)

type InspectionPointCheck struct {
	gorm.Model
	OrderID uint            `gorm:"not null" json:"order_id"`
	Order   InspectionOrder `gorm:"foreignKey:OrderID" json:"order"`
	PointID uint            `gorm:"not null" json:"point_id"`
	Point   InspectionPoint `gorm:"foreignKey:PointID" json:"point"`
	Status  CheckStatus     `gorm:"type:varchar(20);not null" json:"status"`
	Comment string          `gorm:"type:text" json:"comment"`
}
