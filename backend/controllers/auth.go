// Package controllers 认证控制器
package controllers

import (
	"life-rpg/database"
	"life-rpg/middleware"
	"life-rpg/models"
	"life-rpg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthController 认证控制器
type AuthController struct{}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string          `json:"token"`
	UserInfo *models.SysUser `json:"userInfo"`
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param data body LoginRequest true "登录信息"
// @Success 200 {object} utils.Response
// @Router /api/auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	// 查询用户
	var user models.SysUser
	if err := database.DB.Preload("Role").Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Fail(c, "用户名或密码错误")
		return
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Fail(c, "用户名或密码错误")
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		utils.Fail(c, "账号已被禁用")
		return
	}

	// 生成Token
	roleKey := ""
	if user.Role != nil {
		roleKey = user.Role.Key
	}
	token, err := middleware.GenerateToken(user.ID, user.Username, user.RoleID, roleKey)
	if err != nil {
		utils.Fail(c, "Token生成失败")
		return
	}

	utils.Success(c, LoginResponse{
		Token:    token,
		UserInfo: &user,
	})
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

// Register 用户注册
// @Summary 用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param data body RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response
// @Router /api/auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	// 检查用户名是否存在
	var count int64
	database.DB.Model(&models.SysUser{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		utils.Fail(c, "用户名已存在")
		return
	}

	// 加密密码
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// 获取普通用户角色
	var userRole models.SysRole
	database.DB.Where("key = ?", "user").First(&userRole)

	// 创建用户
	user := models.SysUser{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		RoleID:   userRole.ID,
		Gold:     0,
		Exp:      0,
		Level:    1,
		Status:   1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.Fail(c, "注册失败")
		return
	}

	utils.SuccessWithMessage(c, "注册成功", nil)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/auth/info [get]
func (ac *AuthController) GetUserInfo(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var user models.SysUser
	if err := database.DB.Preload("Role").First(&user, userID).Error; err != nil {
		utils.Fail(c, "用户不存在")
		return
	}

	utils.Success(c, &user)
}

// GetUserMenus 获取当前用户菜单
// @Summary 获取当前用户菜单
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/auth/menus [get]
func (ac *AuthController) GetUserMenus(c *gin.Context) {
	roleKey := middleware.GetCurrentRoleKey(c)
	roleID, _ := c.Get("roleId")

	var menus []models.SysMenu

	// 管理员获取所有菜单
	if roleKey == "admin" {
		database.DB.Where("status = 1").Order("sort").Find(&menus)
	} else {
		// 普通用户获取角色关联的菜单
		database.DB.
			Joins("JOIN role_menu ON role_menu.menu_id = sys_menu.id").
			Where("role_menu.role_id = ? AND sys_menu.status = 1", roleID).
			Order("sort").
			Find(&menus)
	}

	// 构建菜单树
	menuTree := buildMenuTree(menus, 0)
	utils.Success(c, menuTree)
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []models.SysMenu, parentID uint) []models.SysMenu {
	var tree []models.SysMenu
	for _, menu := range menus {
		if menu.ParentID == parentID {
			menu.Children = buildMenuTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}
	return tree
}
