import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 10000,
})

// 响应拦截器
api.interceptors.response.use(
  response => response.data,
  error => {
    console.error('API请求失败:', error)
    return Promise.reject(error)
  }
)

export default api
