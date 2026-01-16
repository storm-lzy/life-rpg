// Package controllers 奖励控制器
package controllers

import (
	"strconv"

	"life-rpg/database"
	"life-rpg/middleware"
	"life-rpg/models"
	"life-rpg/utils"

	"github.com/gin-gonic/gin"
)

// RewardController 奖励控制器
type RewardController struct{}

// List 奖励列表 (管理端)
func (rc *RewardController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var rewards []models.Reward
	var total int64

	query := database.DB.Model(&models.Reward{})
	query.Count(&total)
	query.Order("sort, id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rewards)

	utils.PageSuccess(c, rewards, total, page, pageSize)
}

// Create 创建奖励
func (rc *RewardController) Create(c *gin.Context) {
	var reward models.Reward
	if err := c.ShouldBindJSON(&reward); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	if err := database.DB.Create(&reward).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", reward)
}

// Update 更新奖励
func (rc *RewardController) Update(c *gin.Context) {
	id := c.Param("id")
	var reward models.Reward
	if err := database.DB.First(&reward, id).Error; err != nil {
		utils.Fail(c, "奖励不存在")
		return
	}

	var updateData models.Reward
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	database.DB.Model(&reward).Updates(updateData)
	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除奖励
func (rc *RewardController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Reward{}, id).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

// ===== 用户端接口 =====

// UserRewardList 用户奖励列表 (H5端)
func (rc *RewardController) UserRewardList(c *gin.Context) {
	var rewards []models.Reward
	database.DB.Where("is_active = ?", true).Order("sort").Find(&rewards)
	utils.Success(c, rewards)
}

// Purchase 购买奖励
func (rc *RewardController) Purchase(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	rewardID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// 获取奖励信息
	var reward models.Reward
	if err := database.DB.First(&reward, rewardID).Error; err != nil {
		utils.Fail(c, "奖励不存在")
		return
	}

	if !reward.IsActive {
		utils.Fail(c, "奖励已下架")
		return
	}

	// 检查库存
	if reward.Stock == 0 {
		utils.Fail(c, "库存不足")
		return
	}

	// 获取用户信息
	var user models.SysUser
	database.DB.First(&user, userID)

	// 检查金币是否足够
	if user.Gold < reward.Cost {
		utils.Fail(c, "金币不足")
		return
	}

	// 开始事务
	tx := database.DB.Begin()

	// 扣除金币
	newGold := user.Gold - reward.Cost
	tx.Model(&user).Update("gold", newGold)

	// 减少库存
	if reward.Stock > 0 {
		tx.Model(&reward).Update("stock", reward.Stock-1)
	}

	// 记录流水
	log := models.UserLog{
		UserID:      userID,
		Type:        "gold_out",
		Amount:      reward.Cost,
		Balance:     newGold,
		Description: "兑换奖励: " + reward.Title,
		RefType:     "reward",
		RefID:       uint(rewardID),
	}
	tx.Create(&log)

	tx.Commit()

	utils.Success(c, gin.H{
		"cost":    reward.Cost,
		"newGold": newGold,
		"reward":  reward.Title,
	})
}
