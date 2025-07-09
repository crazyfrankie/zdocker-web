import api from './index'

export interface Container {
  id: string
  name: string
  image: string
  command: string
  status: string
  created_time: string
  pid: string
  volume: string
  port_mapping: string
}

export interface CreateContainerRequest {
  image: string
  command: string
  name?: string
  detach?: boolean
  tty?: boolean
  volume?: string
  memory?: string
  cpu_share?: string
  cpu_set?: string
  network?: string
  environment?: Record<string, string>
  port_mapping?: string[]
}

export interface ExecRequest {
  command: string[]
}

// 获取容器列表
export const getContainers = () => {
  return api.get('/containers')
}

// 获取单个容器
export const getContainer = (id: string) => {
  return api.get(`/containers/${id}`)
}

// 创建容器
export const createContainer = (data: CreateContainerRequest) => {
  return api.post('/containers', data)
}

// 启动容器
export const startContainer = (id: string) => {
  return api.post(`/containers/${id}/start`)
}

// 停止容器
export const stopContainer = (name: string) => {
  return api.post(`/containers/stop/${name}`)
}

// 删除容器
export const removeContainer = (name: string) => {
  return api.delete(`/containers/${name}`)
}

// 获取容器日志
export const getContainerLogs = (name: string) => {
  return api.get(`/containers/logs/${name}`)
}

// 在容器中执行命令
export const execContainer = (id: string, data: ExecRequest) => {
  return api.post(`/containers/${id}/exec`, data)
}
