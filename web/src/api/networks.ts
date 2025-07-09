import api from './index'

export interface NetworkInfo {
  name: string
  driver: string
  subnet: string
}

export interface CreateNetworkRequest {
  name: string
  driver?: string
  subnet?: string
}

// 获取网络列表
export const getNetworks = () => {
  return api.get('/networks')
}

// 创建网络
export const createNetwork = (data: CreateNetworkRequest) => {
  return api.post('/networks', data)
}

// 删除网络
export const removeNetwork = (id: string) => {
  return api.delete(`/networks/${id}`)
}
