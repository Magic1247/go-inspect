package controllers

import (
	"go-inspect/config"
	"go-inspect/models"
	"go-inspect/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateInspectionItem 创建巡检项
func CreateInspectionItem(c *gin.Context) {
	var item models.InspectionItem
	if err := c.ShouldBindJSON(&item); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Create(&item).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建巡检项失败")
		return
	}

	utils.SuccessResponse(c, "巡检项创建成功", item)
}

// GetInspectionItem 获取单个巡检项
func GetInspectionItem(c *gin.Context) {
	id := c.Param("id")
	var item models.InspectionItem
	if err := config.DB.First(&item, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检项不存在")
		return
	}
	utils.SuccessResponse(c, "获取巡检项成功", item)
}

// UpdateInspectionItem 更新巡检项
func UpdateInspectionItem(c *gin.Context) {
	id := c.Param("id")
	var item models.InspectionItem
	if err := config.DB.First(&item, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检项不存在")
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Save(&item).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新巡检项失败")
		return
	}

	utils.SuccessResponse(c, "巡检项更新成功", item)
}

// DeleteInspectionItem 删除巡检项
func DeleteInspectionItem(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.InspectionItem{}, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除巡检项失败")
		return
	}
	utils.SuccessResponse(c, "巡检项删除成功", nil)
}

// ListInspectionItems 列出某个点位的所有巡检项
func ListInspectionItems(c *gin.Context) {
	pointID, err := strconv.Atoi(c.Query("point_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的点位ID")
		return
	}

	var items []models.InspectionItem
	if err := config.DB.Where("point_id = ?", pointID).Find(&items).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取巡检项列表失败")
		return
	}

	utils.SuccessResponse(c, "获取巡检项列表成功", items)
}

// AddItemToPoint 将巡检项添加到巡检点
func AddItemToPoint(c *gin.Context) {
	itemID := c.Param("id")
	var pointIDs []uint
	if err := c.ShouldBindJSON(&pointIDs); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var item models.InspectionItem
	if err := config.DB.First(&item, itemID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检项不存在")
		return
	}

	var points []models.InspectionPoint
	if err := config.DB.Find(&points, pointIDs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "一个或多个巡检点不存在")
		return
	}

	if err := config.DB.Model(&item).Association("Points").Append(points); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "添加巡检项到巡检点失败")
		return
	}

	utils.SuccessResponse(c, "巡检项成功添加到巡检点", nil)
}

// RemoveItemFromPoint 从巡检点移除巡检项
func RemoveItemFromPoint(c *gin.Context) {
	itemID := c.Param("id")
	pointID := c.Param("pointId")

	var item models.InspectionItem
	if err := config.DB.First(&item, itemID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检项不存在")
		return
	}

	var point models.InspectionPoint
	if err := config.DB.First(&point, pointID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检点不存在")
		return
	}

	if err := config.DB.Model(&item).Association("Points").Delete(&point); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "从巡检点移除巡检项失败")
		return
	}

	utils.SuccessResponse(c, "巡检项成功从巡检点移除", nil)
}
