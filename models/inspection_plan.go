package models

import (
	"time"

	"gorm.io/gorm"
)

type TriggerType string

const (
	TriggerTypeMonthly TriggerType = "monthly"
	TriggerTypeWeekly  TriggerType = "weekly"
	TriggerTypeManual  TriggerType = "manual"
)

type InspectionPlan struct {
	gorm.Model
	Name            string          `gorm:"type:varchar(100);not null" json:"name"`
	ProjectID       uint            `gorm:"not null" json:"project_id"`
	Project         Project         `gorm:"foreignKey:ProjectID" json:"project"`
	RouteID         uint            `gorm:"not null" json:"route_id"`
	Route           InspectionRoute `gorm:"foreignKey:RouteID" json:"route"`
	TriggerType     TriggerType     `gorm:"type:varchar(20);not null" json:"trigger_type"`
	TriggerDay      int             `gorm:"type:int" json:"trigger_day"` // 1-31 for monthly, 0-6 for weekly (0 = Sunday)
	AssignerID      uint            `gorm:"not null" json:"assigner_id"`
	Assigner        User            `gorm:"foreignKey:AssignerID" json:"assigner"`
	Assignees       []User          `gorm:"many2many:plan_assignees;" json:"assignees"`
	LastTriggeredAt *time.Time      `json:"last_triggered_at"`
}
