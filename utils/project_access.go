package utils

import (
	"go-inspect/config"
	"go-inspect/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HasProjectAccess 检查用户是否有权限访问指定的项目
func HasProjectAccess(c *gin.Context, projectID uint) bool {
	userID, _ := c.Get("userId")
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return false
	}

	// 如果用户没有关联项目，假设他是超级管理员，可以访问所有项目
	if user.ProjectID == nil {
		return true
	}

	// 检查用户的项目ID是否匹配
	return *user.ProjectID == projectID
}

// GetAccessibleProjectIDs 获取用户有权限访问的所有项目ID
func GetAccessibleProjectIDs(c *gin.Context) []uint {
	userID, _ := c.Get("userId")
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return []uint{}
	}

	// 如果用户没有关联项目，假设他是超级管理员，可以访问所有项目
	if user.ProjectID == nil {
		var allProjectIDs []uint
		config.DB.Model(&models.Project{}).Pluck("id", &allProjectIDs)
		return allProjectIDs
	}

	// 返回用户的项目ID
	return []uint{*user.ProjectID}
}

// StringToUint 将字符串转换为uint
func StringToUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint(i)
}

// isProjectOrSubproject 检查 childID 是否是 parentID 本身或其子项目
func isProjectOrSubproject(parentID, childID uint) bool {
	if parentID == childID {
		return true
	}

	var child models.Project
	if err := config.DB.First(&child, childID).Error; err != nil {
		return false
	}

	for child.ParentID != nil {
		if *child.ParentID == parentID {
			return true
		}
		if err := config.DB.First(&child, *child.ParentID).Error; err != nil {
			return false
		}
	}

	return false
}

// getProjectAndSubprojectIDs 获取指定项目及其所有子项目的ID
func getProjectAndSubprojectIDs(projectID uint) []uint {
	var ids []uint
	ids = append(ids, projectID)

	var children []models.Project
	config.DB.Where("parent_id = ?", projectID).Find(&children)

	for _, child := range children {
		ids = append(ids, getProjectAndSubprojectIDs(child.ID)...)
	}

	return ids
}
