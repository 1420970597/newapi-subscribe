import axios from 'axios'
import { useAuthStore } from '../store/auth'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout()
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error)
  }
)

// 认证
export const authApi = {
  login: (data: { username: string; password: string }) => api.post('/auth/login', data),
  register: (data: { username: string; password: string; email?: string }) => api.post('/auth/register', data),
  loginNewAPI: (data: { username: string; password: string }) => api.post('/auth/login/newapi', data),
  me: () => api.get('/auth/me'),
}

// 套餐
export const planApi = {
  list: () => api.get('/plans'),
  get: (id: number) => api.get(`/plans/${id}`),
  getModels: (id: number) => api.get(`/plans/${id}/models`),
}

// 订阅
export const subscriptionApi = {
  current: () => api.get('/subscriptions/current'),
  purchase: (data: any) => api.post('/subscriptions/purchase', data),
  renew: (data: { period_days: number }) => api.post('/subscriptions/renew', data),
  usage: (params?: any) => api.get('/subscriptions/usage', { params }),
  usageDetail: (params?: any) => api.get('/subscriptions/usage/detail', { params }),
}

// 订单
export const orderApi = {
  list: (params?: any) => api.get('/orders', { params }),
  get: (id: number) => api.get(`/orders/${id}`),
  pay: (data: { order_id: number; payment_method: string }) => api.post('/orders/pay', data),
}

// 用户
export const userApi = {
  updateProfile: (data: any) => api.put('/user/profile', data),
  bindNewAPI: (data: { username: string; password: string }) => api.post('/user/bind-newapi', data),
  updateEmailSettings: (data: any) => api.put('/user/email-settings', data),
}

// 管理员
export const adminApi = {
  // 用户
  getUsers: (params?: any) => api.get('/admin/users', { params }),
  getUser: (id: number) => api.get(`/admin/users/${id}`),
  getUserUsage: (id: number, params?: any) => api.get(`/admin/users/${id}/usage`, { params }),
  updateUser: (id: number, data: any) => api.put(`/admin/users/${id}`, data),

  // 订阅
  getSubscriptions: (params?: any) => api.get('/admin/subscriptions', { params }),

  // 订单
  getOrders: (params?: any) => api.get('/admin/orders', { params }),

  // 套餐
  createPlan: (data: any) => api.post('/admin/plans', data),
  updatePlan: (id: number, data: any) => api.put(`/admin/plans/${id}`, data),
  deletePlan: (id: number) => api.delete(`/admin/plans/${id}`),

  // 设置
  getSettings: () => api.get('/admin/settings'),
  updateSettings: (data: any) => api.put('/admin/settings', data),

  // 同步
  triggerSync: () => api.post('/admin/sync/trigger'),

  // new-api
  getNewAPIGroups: () => api.get('/admin/newapi/groups'),
}

export default api
