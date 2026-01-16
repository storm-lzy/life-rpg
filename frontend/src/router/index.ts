// 路由配置
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

// 路由配置
const routes: RouteRecordRaw[] = [
    // 登录页
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue'),
        meta: { requiresAuth: false }
    },

    // PC 管理端
    {
        path: '/admin',
        component: () => import('@/layouts/AdminLayout.vue'),
        meta: { requiresAuth: true, requiresAdmin: true },
        children: [
            {
                path: '',
                redirect: '/admin/dashboard'
            },
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/views/admin/Dashboard.vue'),
                meta: { title: '仪表盘' }
            },
            // 系统管理
            {
                path: 'system/user',
                name: 'SystemUser',
                component: () => import('@/views/admin/system/User.vue'),
                meta: { title: '用户管理' }
            },
            {
                path: 'system/role',
                name: 'SystemRole',
                component: () => import('@/views/admin/system/Role.vue'),
                meta: { title: '角色管理' }
            },
            {
                path: 'system/menu',
                name: 'SystemMenu',
                component: () => import('@/views/admin/system/Menu.vue'),
                meta: { title: '菜单管理' }
            },
            // 游戏配置
            {
                path: 'game/task',
                name: 'GameTask',
                component: () => import('@/views/admin/game/Task.vue'),
                meta: { title: '任务管理' }
            },
            {
                path: 'game/reward',
                name: 'GameReward',
                component: () => import('@/views/admin/game/Reward.vue'),
                meta: { title: '奖励管理' }
            },
            // 公告管理
            {
                path: 'announcement',
                name: 'Announcement',
                component: () => import('@/views/admin/Announcement.vue'),
                meta: { title: '公告管理' }
            },
            // 主题配置
            {
                path: 'theme',
                name: 'Theme',
                component: () => import('@/views/admin/Theme.vue'),
                meta: { title: 'H5主题配置' }
            }
        ]
    },

    // H5 移动端
    {
        path: '/app',
        component: () => import('@/layouts/MobileLayout.vue'),
        meta: { requiresAuth: true },
        children: [
            {
                path: '',
                redirect: '/app/home'
            },
            {
                path: 'home',
                name: 'AppHome',
                component: () => import('@/views/app/Home.vue'),
                meta: { title: '首页' }
            },
            {
                path: 'quest',
                name: 'AppQuest',
                component: () => import('@/views/app/Quest.vue'),
                meta: { title: '任务大厅' }
            },
            {
                path: 'shop',
                name: 'AppShop',
                component: () => import('@/views/app/Shop.vue'),
                meta: { title: '商店' }
            },
            {
                path: 'profile',
                name: 'AppProfile',
                component: () => import('@/views/app/Profile.vue'),
                meta: { title: '我的' }
            }
        ]
    },

    // 根路径重定向
    {
        path: '/',
        redirect: '/login'
    },

    // 404
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/NotFound.vue')
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 路由守卫
router.beforeEach(async (to, _from, next) => {
    const userStore = useUserStore()

    // 设置页面标题
    document.title = (to.meta.title as string) ? `${to.meta.title} - Life RPG` : 'Life RPG'

    // 不需要认证的页面
    if (to.meta.requiresAuth === false) {
        // 已登录用户访问登录页，根据角色跳转
        if (to.path === '/login' && userStore.isLoggedIn) {
            if (userStore.isAdmin) {
                next('/admin')
            } else {
                next('/app')
            }
            return
        }
        next()
        return
    }

    // 需要认证的页面
    if (!userStore.isLoggedIn) {
        next('/login')
        return
    }

    // 获取用户信息
    if (!userStore.userInfo) {
        try {
            await userStore.fetchUserInfo()
            await userStore.fetchMenus()
        } catch {
            next('/login')
            return
        }
    }

    // 检查管理员权限
    if (to.meta.requiresAdmin && !userStore.isAdmin) {
        next('/app')
        return
    }

    next()
})

export default router
