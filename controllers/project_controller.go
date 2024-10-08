package controllers

import (
	"go-inspect/config"
	"go-inspect/models"
	"go-inspect/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateProject 创建项目
func CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Create(&project).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建项目失败")
		return
	}

	utils.SuccessResponse(c, "项目创建成功", project)
}

// GetProject 获取单个项目
func GetProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := config.DB.Preload("Children").First(&project, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "项目不存在")
		return
	}
	utils.SuccessResponse(c, "获取项目成功", project)
}

// UpdateProject 更新项目
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := config.DB.First(&project, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "项目不存在")
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Save(&project).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新项目失败")
		return
	}

	utils.SuccessResponse(c, "项目更新成功", project)
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	// 检查是否有子项目
	var childrenCount int64
	config.DB.Model(&models.Project{}).Where("parent_id = ?", id).Count(&childrenCount)
	if childrenCount > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "无法删除含有子项目的项目")
		return
	}

	if err := config.DB.Delete(&models.Project{}, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除项目失败")
		return
	}
	utils.SuccessResponse(c, "项目删除成功", nil)
}

// ListProjects 列出所有顶级项目
func ListProjects(c *gin.Context) {
	var projects []models.Project
	if err := config.DB.Where("parent_id IS NULL").Preload("Children").Find(&projects).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取项目列表失败")
		return
	}
	utils.SuccessResponse(c, "获取项目列表成功", projects)
}

// GetProjectTree 获取项目树
func GetProjectTree(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := config.DB.Preload("Children.Children.Children").First(&project, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "项目不存在")
		return
	}
	utils.SuccessResponse(c, "获取项目树成功", project)
}
