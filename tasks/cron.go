package tasks

import (
	"go-inspect/config"
	"go-inspect/models"
	"time"

	"github.com/robfig/cron/v3"
)

var cronJob *cron.Cron

// InitCronJobs 初始化并启动定时任务
func InitCronJobs() {
	cronJob = cron.New()

	// 每天凌晨检查并触发巡检计划
	cronJob.AddFunc("0 0 * * *", triggerInspectionPlans)

	// 启动定时任务
	cronJob.Start()
}

// StopCronJobs 停止所有定时任务
func StopCronJobs() {
	if cronJob != nil {
		cronJob.Stop()
	}
}

func triggerInspectionPlans() {
	var plans []models.InspectionPlan
	config.DB.Find(&plans)

	now := time.Now()
	for _, plan := range plans {
		if shouldTriggerPlan(plan, now) {
			createInspectionOrder(plan)
			plan.LastTriggeredAt = &now
			config.DB.Save(&plan)
		}
	}
}

func createInspectionOrder(plan models.InspectionPlan) {
	order := models.InspectionOrder{
		PlanID: plan.ID,
		Status: models.OrderStatusPending,
	}

	if plan.AssignerID != 0 {
		order.AssignerID = &plan.AssignerID
		order.Status = models.OrderStatusAssigned
	}

	config.DB.Create(&order)

	// 获取关联的巡检路线
	var route models.InspectionRoute
	config.DB.Preload("Points").First(&route, plan.RouteID)

	// 创建巡检点确认记录
	for _, point := range route.Points {
		check := models.InspectionPointCheck{
			OrderID: order.ID,
			PointID: point.ID,
			Status:  models.CheckStatusPending,
		}
		config.DB.Create(&check)
	}

	// 更新工单的总巡检点数量
	order.TotalChecks = len(route.Points)
	config.DB.Save(&order)
}

func shouldTriggerPlan(plan models.InspectionPlan, now time.Time) bool {
	switch plan.TriggerType {
	case models.TriggerTypeMonthly:
		return now.Day() == plan.TriggerDay
	case models.TriggerTypeWeekly:
		return int(now.Weekday()) == plan.TriggerDay
	case models.TriggerTypeManual:
		return false // 手动触发的计划不会自动触发
	default:
		return false
	}
}
