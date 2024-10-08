package controllers

import (
	"fmt"
	"go-inspect/config"
	"go-inspect/models"
	"go-inspect/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateInspectionPlan 创建巡检计划
func CreateInspectionPlan(c *gin.Context) {
	var plan models.InspectionPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 检查用户是否有权限创建该项目的巡检计划
	if !utils.HasProjectAccess(c, plan.ProjectID) {
		utils.ErrorResponse(c, http.StatusForbidden, "无权在该项目中创建巡检计划")
		return
	}

	// 验证触发类型和触发日
	if err := validateTrigger(plan.TriggerType, plan.TriggerDay); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Create(&plan).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建巡检计划失败")
		return
	}

	utils.SuccessResponse(c, "巡检计划创建成功", plan)
}

// GetInspectionPlan 获取单个巡检计划
func GetInspectionPlan(c *gin.Context) {
	id := c.Param("id")
	var plan models.InspectionPlan
	if err := config.DB.Preload("Project").Preload("Route").Preload("Assigner").Preload("Assignees").First(&plan, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检计划不存在")
		return
	}

	// 检查用户是否有权限访问该巡检计划
	if !utils.HasProjectAccess(c, plan.ProjectID) {
		utils.ErrorResponse(c, http.StatusForbidden, "无权访问该巡检计划")
		return
	}

	utils.SuccessResponse(c, "获取巡检计划成功", plan)
}

// UpdateInspectionPlan 更新巡检计划
func UpdateInspectionPlan(c *gin.Context) {
	id := c.Param("id")
	var plan models.InspectionPlan
	if err := config.DB.First(&plan, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检计划不存在")
		return
	}

	// 检查用户是否有权限更新该巡检计划
	if !utils.HasProjectAccess(c, plan.ProjectID) {
		utils.ErrorResponse(c, http.StatusForbidden, "无权更新该巡检计划")
		return
	}

	if err := c.ShouldBindJSON(&plan); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 验证触发类型和触发日
	if err := validateTrigger(plan.TriggerType, plan.TriggerDay); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Save(&plan).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新巡检计划失败")
		return
	}

	utils.SuccessResponse(c, "巡检计划更新成功", plan)
}

// DeleteInspectionPlan 删除巡检计划
func DeleteInspectionPlan(c *gin.Context) {
	id := c.Param("id")
	var plan models.InspectionPlan
	if err := config.DB.First(&plan, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检计划不存在")
		return
	}

	// 检查用户是否有权限删除该巡检计划
	if !utils.HasProjectAccess(c, plan.ProjectID) {
		utils.ErrorResponse(c, http.StatusForbidden, "无权删除该巡检计划")
		return
	}

	if err := config.DB.Delete(&plan).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除巡检计划失败")
		return
	}
	utils.SuccessResponse(c, "巡检计划删除成功", nil)
}

// ListInspectionPlans 列出所有巡检计划
func ListInspectionPlans(c *gin.Context) {
	projectID := c.Query("project_id")
	var plans []models.InspectionPlan
	query := config.DB.Preload("Project").Preload("Route").Preload("Assigner").Preload("Assignees")

	if projectID != "" {
		// 检查用户是否有权限访问该项目
		if !utils.HasProjectAccess(c, utils.StringToUint(projectID)) {
			utils.ErrorResponse(c, http.StatusForbidden, "无权访问该项目的巡检计划")
			return
		}
		query = query.Where("project_id = ?", projectID)
	} else {
		// 获取用户有权限访问的所有项目ID
		projectIDs := utils.GetAccessibleProjectIDs(c)
		query = query.Where("project_id IN ?", projectIDs)
	}

	if err := query.Find(&plans).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取巡检计划列表失败")
		return
	}
	utils.SuccessResponse(c, "获取巡检计划列表成功", plans)
}

// TriggerInspectionPlan 手动触发巡检计划
func TriggerInspectionPlan(c *gin.Context) {
	id := c.Param("id")
	var plan models.InspectionPlan
	if err := config.DB.First(&plan, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检计划不存在")
		return
	}

	// 检查用户是否有权限触发该巡检计划
	if !utils.HasProjectAccess(c, plan.ProjectID) {
		utils.ErrorResponse(c, http.StatusForbidden, "无权触发该巡检计划")
		return
	}

	// 在这里添加触发巡检计划的逻辑
	// 例如：创建一个新的巡检任务，更新最后触发时间等

	now := time.Now()
	plan.LastTriggeredAt = &now
	if err := config.DB.Save(&plan).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新巡检计划失败")
		return
	}

	utils.SuccessResponse(c, "巡检计划已手动触发", plan)
}

func validateTrigger(triggerType models.TriggerType, triggerDay int) error {
	switch triggerType {
	case models.TriggerTypeMonthly:
		if triggerDay < 1 || triggerDay > 31 {
			return fmt.Errorf("月触发日必须在1到31之间")
		}
	case models.TriggerTypeWeekly:
		if triggerDay < 0 || triggerDay > 6 {
			return fmt.Errorf("周触发日必须在0到6之间（0表示周日）")
		}
	case models.TriggerTypeManual:
		// 手动触发不需要验证triggerDay
	default:
		return fmt.Errorf("无效的触发类型")
	}
	return nil
}
