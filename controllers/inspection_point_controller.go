package controllers

import (
	"go-inspect/config"
	"go-inspect/models"
	"go-inspect/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateInspectionPoint 创建巡检点
func CreateInspectionPoint(c *gin.Context) {
	var point models.InspectionPoint
	if err := c.ShouldBindJSON(&point); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Create(&point).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建巡检点失败")
		return
	}

	utils.SuccessResponse(c, "巡检点创建成功", point)
}

// GetInspectionPoint 获取单个巡检点
func GetInspectionPoint(c *gin.Context) {
	id := c.Param("id")
	var point models.InspectionPoint
	if err := config.DB.Preload("Items").First(&point, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检点不存在")
		return
	}
	utils.SuccessResponse(c, "获取巡检点成功", point)
}

// UpdateInspectionPoint 更新巡检点
func UpdateInspectionPoint(c *gin.Context) {
	id := c.Param("id")
	var point models.InspectionPoint
	if err := config.DB.First(&point, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检点不存在")
		return
	}

	if err := c.ShouldBindJSON(&point); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Save(&point).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新巡检点失败")
		return
	}

	utils.SuccessResponse(c, "巡检点更新成功", point)
}

// DeleteInspectionPoint 删除巡检点
func DeleteInspectionPoint(c *gin.Context) {
	id := c.Param("id")
	// 直接删除巡检点，不删除关联的巡检项
	if err := config.DB.Delete(&models.InspectionPoint{}, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除巡检点失败")
		return
	}
	utils.SuccessResponse(c, "巡检点删除成功", nil)
}

// ListInspectionPoints 列出所有巡检点
func ListInspectionPoints(c *gin.Context) {
	var points []models.InspectionPoint
	if err := config.DB.Preload("Items").Find(&points).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取巡检点列表失败")
		return
	}
	utils.SuccessResponse(c, "获取巡检点列表成功", points)
}
