import api from './index'

export interface SystemInfo {
  os: string
  architecture: string
  cpus: number
  memory: string
  zdocker_root: string
}

export interface VersionInfo {
  version: string
  api_version: string
  build_date: string
}

// 获取系统信息
export const getSystemInfo = () => {
  return api.get('/system/info')
}

// 获取版本信息
export const getVersion = () => {
  return api.get('/system/version')
}
