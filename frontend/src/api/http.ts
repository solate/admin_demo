import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

const baseURL = import.meta.env.VITE_API_BASE_URL || '/api'

const http: AxiosInstance = axios.create({
  baseURL,
  timeout: 15000
})

http.interceptors.request.use((config: AxiosRequestConfig) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers = config.headers || {}
    ;(config.headers as any).Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response: AxiosResponse) => {
    const data = response.data
    if (typeof data === 'object' && data && 'code' in data && data.code !== 0) {
      ElMessage.error(data.message || '请求失败')
      return Promise.reject(data)
    }
    return data?.data ?? data
  },
  (error) => {
    const status = error?.response?.status
    if (status === 401) {
      ElMessage.error('未登录或登录已过期')
      localStorage.removeItem('token')
      location.href = '/login'
    } else {
      ElMessage.error(error?.response?.data?.message || error.message || '网络错误')
    }
    return Promise.reject(error)
  }
)

export default http


