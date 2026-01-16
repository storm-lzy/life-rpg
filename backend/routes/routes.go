// Package routes 路由配置
package routes

import (
	"life-rpg/controllers"
	"life-rpg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置路由
func SetupRoutes(r *gin.Engine) {
	// 初始化控制器
	authCtrl := &controllers.AuthController{}
	userCtrl := &controllers.UserController{}
	roleCtrl := &controllers.RoleController{}
	menuCtrl := &controllers.MenuController{}
	taskCtrl := &controllers.TaskController{}
	rewardCtrl := &controllers.RewardController{}
	announcementCtrl := &controllers.AnnouncementController{}
	dashboardCtrl := &controllers.DashboardController{}

	// API 路由组
	api := r.Group("/api")
	{
		// ===== 公开接口 =====
		auth := api.Group("/auth")
		{
			auth.POST("/login", authCtrl.Login)
			auth.POST("/register", authCtrl.Register)
		}

		// 获取H5主题配置(公开)
		api.GET("/theme", dashboardCtrl.ThemeConfig)

		// ===== 需要认证的接口 =====
		authenticated := api.Group("")
		authenticated.Use(middleware.JWTAuth())
		{
			// 通用认证接口
			authenticated.GET("/auth/info", authCtrl.GetUserInfo)
			authenticated.GET("/auth/menus", authCtrl.GetUserMenus)

			// ===== 管理端接口 (仅管理员) =====
			admin := authenticated.Group("")
			admin.Use(middleware.AdminRequired())
			{
				// 仪表盘
				admin.GET("/dashboard/stats", dashboardCtrl.Stats)

				// 用户管理
				admin.GET("/users", userCtrl.List)
				admin.POST("/users", userCtrl.Create)
				admin.PUT("/users/:id", userCtrl.Update)
				admin.DELETE("/users/:id", userCtrl.Delete)
				admin.POST("/users/:id/reset-password", userCtrl.ResetPassword)

				// 角色管理
				admin.GET("/roles", roleCtrl.List)
				admin.POST("/roles", roleCtrl.Create)
				admin.PUT("/roles/:id", roleCtrl.Update)
				admin.DELETE("/roles/:id", roleCtrl.Delete)
				admin.GET("/roles/:id/menus", roleCtrl.GetRoleMenus)
				admin.POST("/roles/:id/menus", roleCtrl.AssignMenus)

				// 菜单管理
				admin.GET("/menus", menuCtrl.List)
				admin.GET("/menus/all", menuCtrl.ListAll)
				admin.POST("/menus", menuCtrl.Create)
				admin.PUT("/menus/:id", menuCtrl.Update)
				admin.DELETE("/menus/:id", menuCtrl.Delete)

				// 任务管理
				admin.GET("/tasks", taskCtrl.List)
				admin.POST("/tasks", taskCtrl.Create)
				admin.PUT("/tasks/:id", taskCtrl.Update)
				admin.DELETE("/tasks/:id", taskCtrl.Delete)

				// 奖励管理
				admin.GET("/rewards", rewardCtrl.List)
				admin.POST("/rewards", rewardCtrl.Create)
				admin.PUT("/rewards/:id", rewardCtrl.Update)
				admin.DELETE("/rewards/:id", rewardCtrl.Delete)

				// 公告管理
				admin.GET("/announcements", announcementCtrl.List)
				admin.POST("/announcements", announcementCtrl.Create)
				admin.PUT("/announcements/:id", announcementCtrl.Update)
				admin.DELETE("/announcements/:id", announcementCtrl.Delete)

				// 主题配置
				admin.PUT("/theme", dashboardCtrl.UpdateThemeConfig)
			}

			// ===== 用户端接口 (普通用户) =====
			app := authenticated.Group("/app")
			{
				// 用户资料
				app.GET("/profile", dashboardCtrl.UserProfile)
				app.GET("/logs", dashboardCtrl.UserLogs)

				// 任务
				app.GET("/tasks", taskCtrl.UserTaskList)
				app.POST("/tasks/:id/complete", taskCtrl.CompleteTask)

				// 奖励
				app.GET("/rewards", rewardCtrl.UserRewardList)
				app.POST("/rewards/:id/purchase", rewardCtrl.Purchase)

				// 公告
				app.GET("/announcements", announcementCtrl.UserAnnouncementList)
			}
		}
	}
}
