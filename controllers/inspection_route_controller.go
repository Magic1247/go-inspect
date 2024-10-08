package controllers

import (
	"go-inspect/config"
	"go-inspect/models"
	"go-inspect/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateInspectionRoute 创建巡检路线
func CreateInspectionRoute(c *gin.Context) {
	var route models.InspectionRoute
	if err := c.ShouldBindJSON(&route); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 检查用户是否有权限创建该项目的巡检路线
	if !utils.HasProjectAccess(c, route.ProjectID) {
		utils.ErrorResponse(c, http.StatusForbidden, "无权在该项目中创建巡检路线")
		return
	}

	if err := config.DB.Create(&route).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建巡检路线失败")
		return
	}

	utils.SuccessResponse(c, "巡检路线创建成功", route)
}

// GetInspectionRoute 获取单个巡检路线
func GetInspectionRoute(c *gin.Context) {
	id := c.Param("id")
	var route models.InspectionRoute
	if err := config.DB.Preload("Points").First(&route, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检路线不存在")
		return
	}

	// 检查用户是否有权限访问该巡检路线
	if !utils.HasProjectAccess(c, route.ProjectID) {
		utils.ErrorResponse(c, http.StatusForbidden, "无权访问该巡检路线")
		return
	}

	utils.SuccessResponse(c, "获取巡检路线成功", route)
}

// UpdateInspectionRoute 更新巡检路线
func UpdateInspectionRoute(c *gin.Context) {
	id := c.Param("id")
	var route models.InspectionRoute
	if err := config.DB.First(&route, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检路线不存在")
		return
	}

	// 检查用户是否有权限更新该巡检路线
	if !utils.HasProjectAccess(c, route.ProjectID) {
		utils.ErrorResponse(c, http.StatusForbidden, "无权更新该巡检路线")
		return
	}

	if err := c.ShouldBindJSON(&route); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Save(&route).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新巡检路线失败")
		return
	}

	utils.SuccessResponse(c, "巡检路线更新成功", route)
}

// DeleteInspectionRoute 删除巡检路线
func DeleteInspectionRoute(c *gin.Context) {
	id := c.Param("id")
	var route models.InspectionRoute
	if err := config.DB.First(&route, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "巡检路线不存在")
		return
	}

	// 检查用户是否有权限删除该巡检路线
	if !utils.HasProjectAccess(c, route.ProjectID) {
		utils.ErrorResponse(c, http.StatusForbidden, "无权删除该巡检路线")
		return
	}

	if err := config.DB.Delete(&route).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除巡检路线失败")
		return
	}
	utils.SuccessResponse(c, "巡检路线删除成功", nil)
}

// ListInspectionRoutes 列出所有巡检路线
func ListInspectionRoutes(c *gin.Context) {
	projectID := c.Query("project_id")
	var routes []models.InspectionRoute
	query := config.DB.Preload("Points")

	if projectID != "" {
		// 检查用户是否有权限访问该项目
		if !utils.HasProjectAccess(c, utils.StringToUint(projectID)) {
			utils.ErrorResponse(c, http.StatusForbidden, "无权访问该项目的巡检路线")
			return
		}
		query = query.Where("project_id = ?", projectID)
	} else {
		// 获取用户有权限访问的所有项目ID
		projectIDs := utils.GetAccessibleProjectIDs(c)
		query = query.Where("project_id IN ?", projectIDs)
	}

	if err := query.Find(&routes).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取巡检路线列表失败")
		return
	}
	utils.SuccessResponse(c, "获取巡检路线列表成功", routes)
}

// AddPointToRoute 向路线添加巡检点
func AddPointToRoute(c *gin.Context) {
	routeID := c.Param("id")
	var route models.InspectionRoute
	var point models.InspectionPoint

	if err := config.DB.First(&route, routeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "巡检路线不存在"})
		return
	}

	var input struct {
		PointID uint `json:"point_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&point, input.PointID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "巡检点不存在"})
		return
	}

	if err := config.DB.Model(&route).Association("Points").Append(&point); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加巡检点失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "巡检点已添加到路线"})
}

// RemovePointFromRoute 从路线移除巡检点
func RemovePointFromRoute(c *gin.Context) {
	routeID := c.Param("id")
	pointID := c.Param("pointId")
	var route models.InspectionRoute
	var point models.InspectionPoint

	if err := config.DB.First(&route, routeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "巡检路线不存在"})
		return
	}

	if err := config.DB.First(&point, pointID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "巡检点不存在"})
		return
	}

	if err := config.DB.Model(&route).Association("Points").Delete(&point); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移除巡检点失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "巡检点已从路线移除"})
}
