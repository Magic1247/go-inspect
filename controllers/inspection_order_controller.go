package controllers

import (
	"go-inspect/config"
	"go-inspect/models"
	"go-inspect/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ListInspectionOrders 列出巡检工单
func ListInspectionOrders(c *gin.Context) {
	var orders []models.InspectionOrder
	if err := config.DB.Preload("Plan").Preload("Assigner").Preload("Assignee").Find(&orders).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取巡检工单列表失败")
		return
	}
	utils.SuccessResponse(c, "获取巡检工单列表成功", orders)
}

// GetInspectionOrder 获取单个巡检工单
func GetInspectionOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.InspectionOrder
	if err := config.DB.Preload("Plan").Preload("Assigner").Preload("Assignee").First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检工单不存在")
		return
	}
	utils.SuccessResponse(c, "获取巡检工单成功", order)
}

// AssignInspectionOrder 分配巡检工单
func AssignInspectionOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.InspectionOrder
	if err := config.DB.First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检工单不存在")
		return
	}

	var input struct {
		AssigneeID uint `json:"assignee_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order.AssigneeID = &input.AssigneeID
	order.Status = models.OrderStatusAssigned
	if err := config.DB.Save(&order).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "分配巡检工单失败")
		return
	}

	utils.SuccessResponse(c, "巡检工单分配成功", order)
}

// StartInspectionOrder 开始巡检
func StartInspectionOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.InspectionOrder
	if err := config.DB.First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检工单不存在")
		return
	}

	if order.Status != models.OrderStatusAssigned {
		utils.ErrorResponse(c, http.StatusBadRequest, "巡检工单状态不正确")
		return
	}

	now := time.Now()
	order.Status = models.OrderStatusInProgress
	order.StartTime = &now
	if err := config.DB.Save(&order).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "开始巡检失败")
		return
	}

	utils.SuccessResponse(c, "巡检已开始", order)
}

// CheckInspectionPoint 确认巡检点
func CheckInspectionPoint(c *gin.Context) {
	orderID := c.Param("id")
	pointID := c.Param("pointId")

	var order models.InspectionOrder
	if err := config.DB.First(&order, orderID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检工单不存在")
		return
	}

	if order.Status != models.OrderStatusInProgress {
		utils.ErrorResponse(c, http.StatusBadRequest, "巡检工单状态不正确")
		return
	}

	var check models.InspectionPointCheck
	if err := config.DB.Where("order_id = ? AND point_id = ?", orderID, pointID).First(&check).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检点确认记录不存在")
		return
	}

	var input struct {
		Status  models.CheckStatus `json:"status" binding:"required"`
		Comment string             `json:"comment"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	check.Status = input.Status
	check.Comment = input.Comment
	if err := config.DB.Save(&check).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新巡检点确认失败")
		return
	}

	// 更新已完成的巡检点数量
	order.CompletedChecks++
	if err := config.DB.Save(&order).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新巡检工单失败")
		return
	}

	utils.SuccessResponse(c, "巡检点确认成功", check)
}

// CompleteInspectionOrder 完成巡检
func CompleteInspectionOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.InspectionOrder
	if err := config.DB.First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检工单不存在")
		return
	}

	if order.Status != models.OrderStatusInProgress {
		utils.ErrorResponse(c, http.StatusBadRequest, "巡检工单状态不正确")
		return
	}

	if order.CompletedChecks != order.TotalChecks {
		utils.ErrorResponse(c, http.StatusBadRequest, "还有未确认的巡检点")
		return
	}

	var input struct {
		InspectionData string `json:"inspection_data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	order.Status = models.OrderStatusCompleted
	order.EndTime = &now
	order.InspectionData = input.InspectionData
	if err := config.DB.Save(&order).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "完成巡检失败")
		return
	}

	utils.SuccessResponse(c, "巡检已完成", order)
}
