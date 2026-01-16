// Package controllers 用户控制器
package controllers

import (
	"strconv"

	"life-rpg/database"
	"life-rpg/models"
	"life-rpg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserController 用户控制器
type UserController struct{}

// List 用户列表
// @Summary 获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param username query string false "用户名"
// @Success 200 {object} utils.Response
// @Router /api/users [get]
func (uc *UserController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	username := c.Query("username")

	var users []models.SysUser
	var total int64

	query := database.DB.Model(&models.SysUser{}).Preload("Role")
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)

	utils.PageSuccess(c, users, total, page, pageSize)
}

// Create 创建用户
// @Summary 创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body models.SysUser true "用户信息"
// @Success 200 {object} utils.Response
// @Router /api/users [post]
func (uc *UserController) Create(c *gin.Context) {
	var user models.SysUser
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	// 检查用户名是否存在
	var count int64
	database.DB.Model(&models.SysUser{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		utils.Fail(c, "用户名已存在")
		return
	}

	// 加密密码
	if user.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	} else {
		// 默认密码123456
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", user)
}

// Update 更新用户
// @Summary 更新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param data body models.SysUser true "用户信息"
// @Success 200 {object} utils.Response
// @Router /api/users/{id} [put]
func (uc *UserController) Update(c *gin.Context) {
	id := c.Param("id")
	var user models.SysUser
	if err := database.DB.First(&user, id).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}

	var updateData models.SysUser
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	// 如果传了新密码则加密
	if updateData.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		updateData.Password = string(hashedPassword)
	} else {
		updateData.Password = user.Password
	}

	database.DB.Model(&user).Updates(updateData)
	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除用户
// @Summary 删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response
// @Router /api/users/{id} [delete]
func (uc *UserController) Delete(c *gin.Context) {
	id := c.Param("id")

	// 防止删除管理员
	if id == "1" {
		utils.Fail(c, "不能删除超级管理员")
		return
	}

	if err := database.DB.Delete(&models.SysUser{}, id).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// ResetPassword 重置密码
// @Summary 重置用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response
// @Router /api/users/{id}/reset-password [post]
func (uc *UserController) ResetPassword(c *gin.Context) {
	id := c.Param("id")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	database.DB.Model(&models.SysUser{}).Where("id = ?", id).Update("password", string(hashedPassword))

	utils.SuccessWithMessage(c, "密码已重置为123456", nil)
}
