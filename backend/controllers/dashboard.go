// Package controllers 仪表盘统计控制器
package controllers

import (
	"strconv"
	"time"

	"life-rpg/database"
	"life-rpg/middleware"
	"life-rpg/models"
	"life-rpg/utils"

	"github.com/gin-gonic/gin"
)

// DashboardController 仪表盘控制器
type DashboardController struct{}

// Stats 获取统计数据 (管理端)
func (dc *DashboardController) Stats(c *gin.Context) {
	// 用户总数
	var userCount int64
	database.DB.Model(&models.SysUser{}).Count(&userCount)

	// 今日产生金币
	today := time.Now().Format("2006-01-02")
	var todayGold int64
	database.DB.Model(&models.UserLog{}).
		Where("type = ? AND DATE(created_at) = ?", "gold_in", today).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&todayGold)

	// 今日完成任务数
	var todayTasks int64
	database.DB.Model(&models.UserTask{}).
		Where("DATE(completed_at) = ?", today).
		Count(&todayTasks)

	// 活跃任务数
	var activeTaskCount int64
	database.DB.Model(&models.Task{}).Where("is_active = ?", true).Count(&activeTaskCount)

	// 活跃奖励数
	var activeRewardCount int64
	database.DB.Model(&models.Reward{}).Where("is_active = ?", true).Count(&activeRewardCount)

	// 最近7天每日金币产出
	var dailyGoldStats []struct {
		Date string `json:"date"`
		Gold int    `json:"gold"`
	}
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var gold int
		database.DB.Model(&models.UserLog{}).
			Where("type = ? AND DATE(created_at) = ?", "gold_in", date).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&gold)
		dailyGoldStats = append(dailyGoldStats, struct {
			Date string `json:"date"`
			Gold int    `json:"gold"`
		}{Date: date, Gold: gold})
	}

	// 最近7天每日完成任务数
	var dailyTaskStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		database.DB.Model(&models.UserTask{}).
			Where("DATE(completed_at) = ?", date).
			Count(&count)
		dailyTaskStats = append(dailyTaskStats, struct {
			Date  string `json:"date"`
			Count int64  `json:"count"`
		}{Date: date, Count: count})
	}

	utils.Success(c, gin.H{
		"userCount":         userCount,
		"todayGold":         todayGold,
		"todayTasks":        todayTasks,
		"activeTaskCount":   activeTaskCount,
		"activeRewardCount": activeRewardCount,
		"dailyGoldStats":    dailyGoldStats,
		"dailyTaskStats":    dailyTaskStats,
	})
}

// ===== 用户端接口 =====

// UserProfile 用户个人资料 (H5端)
func (dc *DashboardController) UserProfile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var user models.SysUser
	database.DB.Preload("Role").First(&user, userID)

	// 计算升级所需经验
	nextLevelExp := user.Level * 100
	currentLevelExp := (user.Level - 1) * user.Level / 2 * 100
	expProgress := user.Exp - currentLevelExp

	utils.Success(c, gin.H{
		"user":          user,
		"nextLevelExp":  nextLevelExp,
		"expProgress":   expProgress,
		"expPercentage": float64(expProgress) / float64(nextLevelExp) * 100,
	})
}

// UserLogs 用户流水记录 (H5端)
func (dc *DashboardController) UserLogs(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	logType := c.Query("type") // gold_in, gold_out, exp_in

	var logs []models.UserLog
	var total int64

	query := database.DB.Model(&models.UserLog{}).Where("user_id = ?", userID)
	if logType != "" {
		query = query.Where("type = ?", logType)
	}

	query.Count(&total)
	query.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	utils.PageSuccess(c, logs, total, page, pageSize)
}

// ThemeConfig 获取H5主题配置
func (dc *DashboardController) ThemeConfig(c *gin.Context) {
	var theme models.ThemeConfig
	database.DB.First(&theme)

	// 如果没有配置则返回默认值
	if theme.ID == 0 {
		theme = models.ThemeConfig{
			PrimaryColor:   "#1989fa",
			SecondaryColor: "#ff976a",
			GoldColor:      "#ffd700",
			ExpColor:       "#07c160",
		}
	}

	utils.Success(c, theme)
}

// UpdateThemeConfig 更新H5主题配置 (管理端)
func (dc *DashboardController) UpdateThemeConfig(c *gin.Context) {
	var theme models.ThemeConfig
	if err := c.ShouldBindJSON(&theme); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	// 获取现有配置
	var existing models.ThemeConfig
	database.DB.First(&existing)

	if existing.ID == 0 {
		// 创建新配置
		database.DB.Create(&theme)
	} else {
		// 更新配置
		database.DB.Model(&existing).Updates(theme)
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}
