// API 请求封装
import axios, { type AxiosInstance, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

// 创建 axios 实例
const api: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message, data } = response.data

    // 业务成功
    if (code === 0) {
      return data
    }

    // 业务错误
    ElMessage.error(message || '请求失败')
    return Promise.reject(new Error(message))
  },
  (error) => {
    // HTTP 错误
    if (error.response) {
      const { status, data } = error.response

      if (status === 401) {
        // Token 过期或无效
        const userStore = useUserStore()
        userStore.logout()
        ElMessage.error('登录已过期，请重新登录')
        window.location.href = '/login'
      } else if (status === 403) {
        ElMessage.error('权限不足')
      } else {
        ElMessage.error(data?.message || '服务器错误')
      }
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }

    return Promise.reject(error)
  }
)

// API 模块
export const authApi = {
  login: (data: { username: string; password: string }) => api.post('/auth/login', data),
  register: (data: { username: string; password: string; nickname?: string }) => api.post('/auth/register', data),
  getUserInfo: () => api.get('/auth/info'),
  getMenus: () => api.get('/auth/menus'),
}

export const userApi = {
  list: (params?: { page?: number; pageSize?: number; username?: string }) => api.get('/users', { params }),
  create: (data: any) => api.post('/users', data),
  update: (id: number, data: any) => api.put(`/users/${id}`, data),
  delete: (id: number) => api.delete(`/users/${id}`),
  resetPassword: (id: number) => api.post(`/users/${id}/reset-password`),
}

export const roleApi = {
  list: () => api.get('/roles'),
  create: (data: any) => api.post('/roles', data),
  update: (id: number, data: any) => api.put(`/roles/${id}`, data),
  delete: (id: number) => api.delete(`/roles/${id}`),
  getMenus: (id: number) => api.get(`/roles/${id}/menus`),
  assignMenus: (id: number, menuIds: number[]) => api.post(`/roles/${id}/menus`, { menuIds }),
}

export const menuApi = {
  tree: () => api.get('/menus'),
  list: () => api.get('/menus/all'),
  create: (data: any) => api.post('/menus', data),
  update: (id: number, data: any) => api.put(`/menus/${id}`, data),
  delete: (id: number) => api.delete(`/menus/${id}`),
}

export const taskApi = {
  // 管理端
  list: (params?: { page?: number; pageSize?: number; type?: string }) => api.get('/tasks', { params }),
  create: (data: any) => api.post('/tasks', data),
  update: (id: number, data: any) => api.put(`/tasks/${id}`, data),
  delete: (id: number) => api.delete(`/tasks/${id}`),
  // 用户端
  userList: () => api.get('/app/tasks'),
  complete: (id: number) => api.post(`/app/tasks/${id}/complete`),
}

export const rewardApi = {
  // 管理端
  list: (params?: { page?: number; pageSize?: number }) => api.get('/rewards', { params }),
  create: (data: any) => api.post('/rewards', data),
  update: (id: number, data: any) => api.put(`/rewards/${id}`, data),
  delete: (id: number) => api.delete(`/rewards/${id}`),
  // 用户端
  userList: () => api.get('/app/rewards'),
  purchase: (id: number) => api.post(`/app/rewards/${id}/purchase`),
}

export const announcementApi = {
  // 管理端
  list: (params?: { page?: number; pageSize?: number }) => api.get('/announcements', { params }),
  create: (data: any) => api.post('/announcements', data),
  update: (id: number, data: any) => api.put(`/announcements/${id}`, data),
  delete: (id: number) => api.delete(`/announcements/${id}`),
  // 用户端
  userList: () => api.get('/app/announcements'),
}

export const dashboardApi = {
  stats: () => api.get('/dashboard/stats'),
  // 用户端
  profile: () => api.get('/app/profile'),
  logs: (params?: { page?: number; pageSize?: number; type?: string }) => api.get('/app/logs', { params }),
}

export const themeApi = {
  get: () => api.get('/theme'),
  update: (data: any) => api.put('/theme', data),
}

export default api
