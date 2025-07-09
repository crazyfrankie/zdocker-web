<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getContainers } from '@/api/containers'
import { getNetworks } from '@/api/networks'
import { getSystemInfo } from '@/api/system'
import { Box, Plus, Monitor, Connection } from '@element-plus/icons-vue'

const loading = ref(true)
const containerCount = ref(0)
const runningContainers = ref(0)
const networkCount = ref(0)
const systemInfo = ref<any>({})
const recentContainers = ref<any[]>([])

onMounted(async () => {
  try {
    await loadData()
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
})

const loadData = async () => {
  const [containersRes, networksRes, systemRes] = await Promise.all([
    getContainers(),
    getNetworks(),
    getSystemInfo()
  ])

  const containers = containersRes.data || []
  containerCount.value = containers.length
  runningContainers.value = containers.filter((c: any) => c.status === 'running').length
  networkCount.value = networksRes.data?.length || 0
  systemInfo.value = systemRes.data || {}
  
  // 取最近创建的5个容器
  recentContainers.value = containers
    .sort((a: any, b: any) => new Date(b.created_time).getTime() - new Date(a.created_time).getTime())
    .slice(0, 5)
}

const getStatusColor = (status: string) => {
  return status === 'running' ? 'success' : 'info'
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString()
}
</script>

<template>
  <div class="dashboard" v-loading="loading">
    <div class="dashboard-header">
      <h2>概览</h2>
      <p class="subtitle">ZDocker 容器管理面板</p>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon container-icon">
            <el-icon><Box /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ containerCount }}</div>
            <div class="stat-label">总容器数</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon running-icon">
            <el-icon><Monitor /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ runningContainers }}</div>
            <div class="stat-label">运行中</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon network-icon">
            <el-icon><Connection /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ networkCount }}</div>
            <div class="stat-label">网络数量</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon cpu-icon">
            <el-icon><Monitor /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ systemInfo.cpus || 4 }}</div>
            <div class="stat-label">CPU 核心</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 系统信息 -->
    <div class="info-grid">
      <el-card class="info-card">
        <template #header>
          <span>系统信息</span>
        </template>
        <div class="system-info">
          <div class="info-item">
            <span class="info-label">操作系统:</span>
            <span class="info-value">{{ systemInfo.os || 'Linux' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">架构:</span>
            <span class="info-value">{{ systemInfo.architecture || 'amd64' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">内存:</span>
            <span class="info-value">{{ systemInfo.memory || 'Unknown' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">ZDocker 根目录:</span>
            <span class="info-value">{{ systemInfo.zdocker_root || '/var/lib/zdocker' }}</span>
          </div>
        </div>
      </el-card>

      <!-- 最近容器 -->
      <el-card class="info-card">
        <template #header>
          <span>最近容器</span>
        </template>
        <div class="recent-containers">
          <div 
            v-for="container in recentContainers" 
            :key="container.id"
            class="container-item"
          >
            <div class="container-name">
              <el-tag 
                :type="getStatusColor(container.status)" 
                size="small"
              >
                {{ container.status }}
              </el-tag>
              <span class="name">{{ container.name || container.id.substring(0, 12) }}</span>
            </div>
            <div class="container-info">
              <span class="image">{{ container.image }}</span>
              <span class="time">{{ formatTime(container.created_time) }}</span>
            </div>
          </div>
          <div v-if="recentContainers.length === 0" class="no-data">
            暂无容器
          </div>
        </div>
      </el-card>
    </div>

    <!-- 快速操作 -->
    <el-card class="quick-actions">
      <template #header>
        <span>快速操作</span>
      </template>
      <div class="actions-grid">
        <el-button 
          type="primary" 
          :icon="Plus"
          @click="$router.push('/containers')"
        >
          创建容器
        </el-button>
        <el-button 
          :icon="Box"
          @click="$router.push('/containers')"
        >
          管理容器
        </el-button>
        <el-button 
          :icon="Connection"
          @click="$router.push('/networks')"
        >
          管理网络
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.dashboard {
  width: 100%;
}

.dashboard-header {
  margin-bottom: 24px;
}

.dashboard-header h2 {
  font-size: 28px;
  color: #303133;
  margin-bottom: 8px;
}

.subtitle {
  color: #909399;
  font-size: 14px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.container-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.running-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.network-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.cpu-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 24px;
}

.info-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.system-info {
  space-y: 12px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}

.info-label {
  color: #909399;
  font-size: 14px;
}

.info-value {
  color: #303133;
  font-size: 14px;
  font-weight: 500;
}

.recent-containers {
  space-y: 12px;
}

.container-item {
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.container-item:last-child {
  border-bottom: none;
}

.container-name {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.name {
  font-weight: 500;
  color: #303133;
}

.container-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
}

.no-data {
  text-align: center;
  color: #909399;
  padding: 20px;
}

.quick-actions {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.actions-grid {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  }
}
</style>
