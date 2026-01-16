// 用户状态管理
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api'

// 用户信息接口
export interface UserInfo {
    id: number
    username: string
    nickname: string
    avatar: string
    roleId: number
    role?: {
        id: number
        name: string
        key: string
    }
    gold: number
    exp: number
    level: number
    status: number
}

// 菜单项接口
export interface MenuItem {
    id: number
    parentId: number
    name: string
    path: string
    component: string
    icon: string
    sort: number
    type: number
    permission: string
    visible: number
    children?: MenuItem[]
}

export const useUserStore = defineStore('user', () => {
    // 状态
    const token = ref<string>(localStorage.getItem('token') || '')
    const userInfo = ref<UserInfo | null>(null)
    const menus = ref<MenuItem[]>([])

    // 计算属性
    const isLoggedIn = computed(() => !!token.value)
    const isAdmin = computed(() => userInfo.value?.role?.key === 'admin')
    const username = computed(() => userInfo.value?.username || '')
    const gold = computed(() => userInfo.value?.gold || 0)
    const exp = computed(() => userInfo.value?.exp || 0)
    const level = computed(() => userInfo.value?.level || 1)

    // 登录
    async function login(username: string, password: string) {
        const data: any = await authApi.login({ username, password })
        token.value = data.token
        userInfo.value = data.userInfo
        localStorage.setItem('token', data.token)
        return data
    }

    // 获取用户信息
    async function fetchUserInfo() {
        if (!token.value) return
        try {
            const data: any = await authApi.getUserInfo()
            userInfo.value = data
            return data
        } catch (error) {
            logout()
            throw error
        }
    }

    // 获取菜单
    async function fetchMenus() {
        if (!token.value) return
        try {
            const data: any = await authApi.getMenus()
            menus.value = data || []
            return data
        } catch {
            menus.value = []
        }
    }

    // 更新用户金币和经验
    function updateUserStats(newGold: number, newExp: number, newLevel: number) {
        if (userInfo.value) {
            userInfo.value.gold = newGold
            userInfo.value.exp = newExp
            userInfo.value.level = newLevel
        }
    }

    // 登出
    function logout() {
        token.value = ''
        userInfo.value = null
        menus.value = []
        localStorage.removeItem('token')
    }

    return {
        token,
        userInfo,
        menus,
        isLoggedIn,
        isAdmin,
        username,
        gold,
        exp,
        level,
        login,
        fetchUserInfo,
        fetchMenus,
        updateUserStats,
        logout,
    }
})
