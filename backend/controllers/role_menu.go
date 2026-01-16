// Package controllers 角色和菜单控制器
package controllers

import (
	"strconv"

	"life-rpg/database"
	"life-rpg/models"
	"life-rpg/utils"

	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct{}

// List 角色列表
func (rc *RoleController) List(c *gin.Context) {
	var roles []models.SysRole
	database.DB.Order("sort").Find(&roles)
	utils.Success(c, roles)
}

// Create 创建角色
func (rc *RoleController) Create(c *gin.Context) {
	var role models.SysRole
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	if err := database.DB.Create(&role).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", role)
}

// Update 更新角色
func (rc *RoleController) Update(c *gin.Context) {
	id := c.Param("id")
	var role models.SysRole
	if err := database.DB.First(&role, id).Error; err != nil {
		utils.Fail(c, "角色不存在")
		return
	}

	var updateData models.SysRole
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	database.DB.Model(&role).Updates(updateData)
	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除角色
func (rc *RoleController) Delete(c *gin.Context) {
	id := c.Param("id")

	// 防止删除内置角色
	if id == "1" || id == "2" {
		utils.Fail(c, "不能删除内置角色")
		return
	}

	if err := database.DB.Delete(&models.SysRole{}, id).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

// AssignMenusRequest 分配菜单请求
type AssignMenusRequest struct {
	MenuIDs []uint `json:"menuIds"`
}

// AssignMenus 为角色分配菜单
func (rc *RoleController) AssignMenus(c *gin.Context) {
	id := c.Param("id")
	roleID, _ := strconv.ParseUint(id, 10, 32)

	var req AssignMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	// 删除原有的角色菜单关联
	database.DB.Where("role_id = ?", roleID).Delete(&models.RoleMenu{})

	// 创建新的关联
	for _, menuID := range req.MenuIDs {
		database.DB.Create(&models.RoleMenu{RoleID: uint(roleID), MenuID: menuID})
	}

	utils.SuccessWithMessage(c, "分配成功", nil)
}

// GetRoleMenus 获取角色菜单ID列表
func (rc *RoleController) GetRoleMenus(c *gin.Context) {
	id := c.Param("id")

	var roleMenus []models.RoleMenu
	database.DB.Where("role_id = ?", id).Find(&roleMenus)

	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	utils.Success(c, menuIDs)
}

// MenuController 菜单控制器
type MenuController struct{}

// List 菜单树列表
func (mc *MenuController) List(c *gin.Context) {
	var menus []models.SysMenu
	database.DB.Order("sort").Find(&menus)

	// 构建树形结构
	tree := buildMenuTree(menus, 0)
	utils.Success(c, tree)
}

// ListAll 所有菜单列表(扁平)
func (mc *MenuController) ListAll(c *gin.Context) {
	var menus []models.SysMenu
	database.DB.Order("sort").Find(&menus)
	utils.Success(c, menus)
}

// Create 创建菜单
func (mc *MenuController) Create(c *gin.Context) {
	var menu models.SysMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	if err := database.DB.Create(&menu).Error; err != nil {
		utils.Fail(c, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", menu)
}

// Update 更新菜单
func (mc *MenuController) Update(c *gin.Context) {
	id := c.Param("id")
	var menu models.SysMenu
	if err := database.DB.First(&menu, id).Error; err != nil {
		utils.Fail(c, "菜单不存在")
		return
	}

	var updateData models.SysMenu
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	database.DB.Model(&menu).Updates(updateData)
	utils.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除菜单
func (mc *MenuController) Delete(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有子菜单
	var count int64
	database.DB.Model(&models.SysMenu{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		utils.Fail(c, "存在子菜单，无法删除")
		return
	}

	if err := database.DB.Delete(&models.SysMenu{}, id).Error; err != nil {
		utils.Fail(c, "删除失败")
		return
	}

	// 删除关联
	database.DB.Where("menu_id = ?", id).Delete(&models.RoleMenu{})

	utils.SuccessWithMessage(c, "删除成功", nil)
}
