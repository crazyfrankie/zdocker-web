<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Bell, User, Setting } from '@element-plus/icons-vue'
import { getSystemInfo, getVersion } from '@/api/system'

const systemInfo = ref<any>({})
const version = ref('')

onMounted(async () => {
  try {
    const [sysInfo, versionInfo] = await Promise.all([
      getSystemInfo(),
      getVersion()
    ])
    systemInfo.value = sysInfo.data
    version.value = versionInfo.data.version
  } catch (error) {
    console.error('获取系统信息失败:', error)
  }
})
</script>

<template>
  <div class="header">
    <div class="header-left">
      <h1 class="title">
        <span class="logo">🐳</span>
        ZDocker Desktop
      </h1>
      <span class="version">{{ version || 'v1.0.0' }}</span>
    </div>
    
    <div class="header-right">
      <div class="system-info">
        <span>{{ systemInfo.os || 'Linux' }} | {{ systemInfo.cpus || 4 }} CPUs</span>
      </div>
      
      <el-dropdown>
        <span class="user-menu">
          <el-icon><User /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item>设置</el-dropdown-item>
            <el-dropdown-item>关于</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<style scoped>
.header {
  height: 60px;
  background: linear-gradient(135deg, #2196F3 0%, #1976D2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: bold;
}

.logo {
  font-size: 24px;
}

.version {
  background: rgba(255,255,255,0.2);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.system-info {
  font-size: 14px;
  opacity: 0.9;
}

.user-menu {
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  transition: background-color 0.3s;
}

.user-menu:hover {
  background-color: rgba(255,255,255,0.1);
}
</style>
