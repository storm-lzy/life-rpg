// Package controllers 任务控制器
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

// TaskController 任务控制器
type TaskController struct{}

// List 任务列表 (管理端)
func (tc *TaskController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	taskType := c.Query("type")

	var tasks []models.Task
	var total int64

	query := database.DB.Model(&models.Task{})
	if taskType != "" {
		query = query.Where("type = ?", taskType)
	}

	query.Count(&total)
	query.Order("sort, id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks)

	utils.PageSuccess(c, tasks, total, page, pageSize)
}

// Create 创建任务
func (tc *TaskController) Create(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	if err := database.DB.Create(&task).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", task)
}

// Update 更新任务
func (tc *TaskController) Update(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		utils.Fail(c, "任务不存在")
		return
	}

	var updateData models.Task
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	database.DB.Model(&task).Updates(updateData)
	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除任务
func (tc *TaskController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Task{}, id).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

// ===== 用户端接口 =====

// UserTaskList 用户任务列表 (H5端)
func (tc *TaskController) UserTaskList(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	// 获取所有激活的任务
	var tasks []models.Task
	database.DB.Where("is_active = ?", true).Order("sort").Find(&tasks)

	// 获取用户今日已完成的任务
	today := time.Now().Format("2006-01-02")
	var completedTaskIDs []uint
	database.DB.Model(&models.UserTask{}).
		Where("user_id = ? AND DATE(completed_at) = ?", userID, today).
		Pluck("task_id", &completedTaskIDs)

	// 构建返回结构
	type TaskWithStatus struct {
		models.Task
		Completed bool `json:"completed"`
	}

	var result []TaskWithStatus
	for _, task := range tasks {
		completed := false
		for _, id := range completedTaskIDs {
			if task.ID == id {
				completed = true
				break
			}
		}
		// 一次性任务检查是否已完成过
		if task.Type == "once" && !completed {
			var count int64
			database.DB.Model(&models.UserTask{}).
				Where("user_id = ? AND task_id = ?", userID, task.ID).
				Count(&count)
			completed = count > 0
		}
		result = append(result, TaskWithStatus{Task: task, Completed: completed})
	}

	utils.Success(c, result)
}

// CompleteTask 完成任务
func (tc *TaskController) CompleteTask(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	taskID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// 获取任务信息
	var task models.Task
	if err := database.DB.First(&task, taskID).Error; err != nil {
		utils.Fail(c, "任务不存在")
		return
	}

	if !task.IsActive {
		utils.Fail(c, "任务已下架")
		return
	}

	// 检查任务是否已完成
	today := time.Now().Format("2006-01-02")
	var count int64

	if task.Type == "daily" {
		// 每日任务：检查今天是否已完成
		database.DB.Model(&models.UserTask{}).
			Where("user_id = ? AND task_id = ? AND DATE(completed_at) = ?", userID, taskID, today).
			Count(&count)
	} else {
		// 一次性任务：检查是否已完成过
		database.DB.Model(&models.UserTask{}).
			Where("user_id = ? AND task_id = ?", userID, taskID).
			Count(&count)
	}

	if count > 0 {
		utils.Fail(c, "任务已完成")
		return
	}

	// 开始事务
	tx := database.DB.Begin()

	// 记录完成
	userTask := models.UserTask{
		UserID:      userID,
		TaskID:      uint(taskID),
		CompletedAt: time.Now(),
	}
	if err := tx.Create(&userTask).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, "完成任务失败")
		return
	}

	// 更新用户金币和经验
	var user models.SysUser
	tx.First(&user, userID)

	newGold := user.Gold + task.GoldReward
	newExp := user.Exp + task.ExpReward
	newLevel := calculateLevel(newExp)

	tx.Model(&user).Updates(map[string]interface{}{
		"gold":  newGold,
		"exp":   newExp,
		"level": newLevel,
	})

	// 记录金币流水
	if task.GoldReward > 0 {
		goldLog := models.UserLog{
			UserID:      userID,
			Type:        "gold_in",
			Amount:      task.GoldReward,
			Balance:     newGold,
			Description: "完成任务: " + task.Title,
			RefType:     "task",
			RefID:       uint(taskID),
		}
		tx.Create(&goldLog)
	}

	// 记录经验流水
	if task.ExpReward > 0 {
		expLog := models.UserLog{
			UserID:      userID,
			Type:        "exp_in",
			Amount:      task.ExpReward,
			Balance:     newExp,
			Description: "完成任务: " + task.Title,
			RefType:     "task",
			RefID:       uint(taskID),
		}
		tx.Create(&expLog)
	}

	tx.Commit()

	// 返回奖励信息
	utils.Success(c, gin.H{
		"goldReward": task.GoldReward,
		"expReward":  task.ExpReward,
		"newGold":    newGold,
		"newExp":     newExp,
		"newLevel":   newLevel,
		"levelUp":    newLevel > user.Level,
	})
}

// calculateLevel 根据经验计算等级
func calculateLevel(exp int) int {
	// 等级公式: 每升一级需要的经验 = 等级 * 100
	// Level 1: 0-99, Level 2: 100-299, Level 3: 300-599, ...
	level := 1
	requiredExp := 0
	for {
		requiredExp += level * 100
		if exp < requiredExp {
			return level
		}
		level++
		if level > 100 { // 最高100级
			return 100
		}
	}
}
