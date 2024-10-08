package controllers

import (
	"go-inspect/config"
	"go-inspect/models"
	"go-inspect/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// 如果提供了项目ID，验证项目是否存在
	if user.ProjectID != nil {
		var project models.Project
		if err := config.DB.First(&project, *user.ProjectID).Error; err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "指定的项目不存在")
			return
		}
	}

	if user.Username == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户名不能为空")
		return
	}

	if user.Password == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "密码不能为空")
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "用户名已存在")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "密码加密失败")
		return
	}
	user.Password = string(hashedPassword)

	if err := config.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "用户注册失败")
		return
	}

	utils.SuccessResponse(c, "用户注册成功", gin.H{"id": user.ID, "username": user.Username, "email": user.Email, "project_id": user.ProjectID})
}

// Login 用户登录
func Login(c *gin.Context) {
	var loginForm struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginForm); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", loginForm.Username).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	utils.SuccessResponse(c, "登录成功", gin.H{"token": token})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userId, _ := c.Get("userId")
	var user models.User
	if err := config.DB.Preload("Project").First(&user, userId).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	utils.SuccessResponse(c, "获取用户信息成功", user)
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	userId, _ := c.Get("userId")
	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	var updateForm struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		ProjectID *uint  `json:"project_id"`
	}

	if err := c.ShouldBindJSON(&updateForm); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if updateForm.Username != "" {
		user.Username = updateForm.Username
	}
	if updateForm.Email != "" {
		user.Email = updateForm.Email
	}
	if updateForm.ProjectID != nil {
		// 验证项目是否存在
		var project models.Project
		if err := config.DB.First(&project, *updateForm.ProjectID).Error; err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "指定的项目不存在")
			return
		}
		user.ProjectID = updateForm.ProjectID
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新用户信息失败")
		return
	}

	utils.SuccessResponse(c, "用户信息更新成功", user)
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	userId, _ := c.Get("userId")
	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	var passwordForm struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&passwordForm); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordForm.OldPassword)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "旧密码错误")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordForm.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "密码加密失败")
		return
	}

	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "修改密码失败")
		return
	}

	utils.SuccessResponse(c, "密码修改成功", nil)
}
