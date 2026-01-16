// Package controllers 公告控制器
package controllers

import (
	"strconv"

	"life-rpg/database"
	"life-rpg/models"
	"life-rpg/utils"

	"github.com/gin-gonic/gin"
)

// AnnouncementController 公告控制器
type AnnouncementController struct{}

// List 公告列表 (管理端)
func (ac *AnnouncementController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var announcements []models.Announcement
	var total int64

	query := database.DB.Model(&models.Announcement{})
	query.Count(&total)
	query.Order("sort, id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&announcements)

	utils.PageSuccess(c, announcements, total, page, pageSize)
}

// Create 创建公告
func (ac *AnnouncementController) Create(c *gin.Context) {
	var announcement models.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	if err := database.DB.Create(&announcement).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", announcement)
}

// Update 更新公告
func (ac *AnnouncementController) Update(c *gin.Context) {
	id := c.Param("id")
	var announcement models.Announcement
	if err := database.DB.First(&announcement, id).Error; err != nil {
		utils.Fail(c, "公告不存在")
		return
	}

	var updateData models.Announcement
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	database.DB.Model(&announcement).Updates(updateData)
	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除公告
func (ac *AnnouncementController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Announcement{}, id).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

// ===== 用户端接口 =====

// UserAnnouncementList 用户公告列表 (H5端)
func (ac *AnnouncementController) UserAnnouncementList(c *gin.Context) {
	var announcements []models.Announcement
	database.DB.Where("is_active = ?", true).Order("sort, created_at desc").Limit(10).Find(&announcements)
	utils.Success(c, announcements)
}
