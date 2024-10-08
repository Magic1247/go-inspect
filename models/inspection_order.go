package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusAssigned   OrderStatus = "assigned"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusCompleted  OrderStatus = "completed"
)

type InspectionOrder struct {
	gorm.Model
	PlanID          uint           `gorm:"not null" json:"plan_id"`
	Plan            InspectionPlan `gorm:"foreignKey:PlanID" json:"plan"`
	Status          OrderStatus    `gorm:"type:varchar(20);not null" json:"status"`
	AssignerID      *uint          `json:"assigner_id"`
	Assigner        *User          `gorm:"foreignKey:AssignerID" json:"assigner"`
	AssigneeID      *uint          `json:"assignee_id"`
	Assignee        *User          `gorm:"foreignKey:AssigneeID" json:"assignee"`
	StartTime       *time.Time     `json:"start_time"`
	EndTime         *time.Time     `json:"end_time"`
	InspectionData  string         `gorm:"type:text" json:"inspection_data"`
	CompletedChecks int            `gorm:"default:0" json:"completed_checks"`
	TotalChecks     int            `gorm:"default:0" json:"total_checks"`
}
