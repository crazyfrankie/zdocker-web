<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getContainers, getContainerLogs } from '@/api/containers'
import { 
  DocumentCopy, 
  Refresh, 
  Download,
  Delete
} from '@element-plus/icons-vue'

const loading = ref(true)
const containers = ref<any[]>([])
const selectedContainer = ref('')
const logs = ref('')
const logsLoading = ref(false)
const autoRefresh = ref(false)
const refreshInterval = ref<any>(null)

onMounted(() => {
  loadContainers()
})

const loadContainers = async () => {
  try {
    loading.value = true
    const response = await getContainers()
    containers.value = response.data || []
    if (containers.value.length > 0 && !selectedContainer.value) {
      selectedContainer.value = containers.value[0].id
      loadLogs()
    }
  } catch (error) {
    console.error('获取容器列表失败:', error)
  } finally {
    loading.value = false
  }
}

const loadLogs = async () => {
  if (!selectedContainer.value) return
  
  try {
    logsLoading.value = true
    const response = await getContainerLogs(selectedContainer.value)
    logs.value = response.data || '暂无日志'
  } catch (error) {
    console.error('获取日志失败:', error)
    logs.value = '获取日志失败'
  } finally {
    logsLoading.value = false
  }
}

const handleContainerChange = () => {
  loadLogs()
}

const handleAutoRefreshChange = () => {
  if (autoRefresh.value) {
    refreshInterval.value = setInterval(loadLogs, 3000)
  } else {
    if (refreshInterval.value) {
      clearInterval(refreshInterval.value)
      refreshInterval.value = null
    }
  }
}

const clearLogs = () => {
  logs.value = ''
}

const downloadLogs = () => {
  if (!logs.value) return
  
  const blob = new Blob([logs.value], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `container-${selectedContainer.value}-logs.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const getContainerName = (id: string) => {
  const container = containers.value.find(c => c.id === id)
  return container?.name || id.substring(0, 12)
}

// 组件销毁时清理定时器
import { onUnmounted } from 'vue'
onUnmounted(() => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
  }
})
</script>

<template>
  <div class="logs-view">
    <div class="page-header">
      <div class="header-left">
        <h2>容器日志</h2>
        <p class="subtitle">查看和管理容器日志输出</p>
      </div>
      <div class="header-actions">
        <el-select 
          v-model="selectedContainer" 
          placeholder="选择容器"
          style="width: 200px; margin-right: 12px;"
          @change="handleContainerChange"
          :loading="loading"
        >
          <el-option
            v-for="container in containers"
            :key="container.id"
            :label="container.name || container.id.substring(0, 12)"
            :value="container.id"
          >
            <div class="container-option">
              <span class="container-name">{{ container.name || container.id.substring(0, 12) }}</span>
              <el-tag 
                :type="container.status === 'running' ? 'success' : 'info'" 
                size="small"
              >
                {{ container.status }}
              </el-tag>
            </div>
          </el-option>
        </el-select>
        
        <el-checkbox 
          v-model="autoRefresh"
          @change="handleAutoRefreshChange"
          style="margin-right: 12px;"
        >
          自动刷新
        </el-checkbox>
        
        <el-button 
          :icon="Refresh"
          @click="loadLogs"
          :loading="logsLoading"
        >
          刷新
        </el-button>
        
        <el-button 
          :icon="Download"
          @click="downloadLogs"
          :disabled="!logs"
        >
          下载
        </el-button>
        
        <el-button 
          :icon="Delete"
          @click="clearLogs"
          :disabled="!logs"
        >
          清空
        </el-button>
      </div>
    </div>

    <el-card class="logs-card" v-if="selectedContainer">
      <template #header>
        <div class="logs-header">
          <span>
            <el-icon><DocumentCopy /></el-icon>
            {{ getContainerName(selectedContainer) }} 的日志
          </span>
          <div class="logs-info">
            <span v-if="autoRefresh" class="auto-refresh-indicator">
              <el-icon class="spinning"><Refresh /></el-icon>
              自动刷新中
            </span>
          </div>
        </div>
      </template>
      
      <div class="logs-container" v-loading="logsLoading">
        <div class="logs-content">
          <pre class="logs-text">{{ logs }}</pre>
        </div>
      </div>
    </el-card>

    <!-- 无容器提示 -->
    <el-card class="empty-card" v-else-if="!loading && containers.length === 0">
      <el-empty description="暂无容器">
        <el-button type="primary" @click="$router.push('/containers')">
          去创建容器
        </el-button>
      </el-empty>
    </el-card>

    <!-- 使用说明 -->
    <el-card class="help-card">
      <template #header>
        <span>使用说明</span>
      </template>
      <div class="help-content">
        <ul>
          <li>选择要查看日志的容器</li>
          <li>开启"自动刷新"可以实时查看最新日志</li>
          <li>点击"下载"可以将日志保存为文本文件</li>
          <li>使用"清空"可以临时清空显示的日志内容</li>
        </ul>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.logs-view {
  width: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.header-left h2 {
  font-size: 28px;
  color: #303133;
  margin-bottom: 8px;
}

.subtitle {
  color: #909399;
  font-size: 14px;
}

.header-actions {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.container-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.container-name {
  flex: 1;
  margin-right: 8px;
}

.logs-card, .empty-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logs-info {
  display: flex;
  align-items: center;
}

.auto-refresh-indicator {
  color: #67c23a;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.logs-container {
  height: 500px;
  overflow: hidden;
}

.logs-content {
  height: 100%;
  overflow-y: auto;
  background-color: #1e1e1e;
  border-radius: 4px;
  padding: 16px;
}

.logs-text {
  font-family: 'Courier New', Courier, monospace;
  font-size: 12px;
  line-height: 1.4;
  color: #d4d4d4;
  white-space: pre-wrap;
  word-wrap: break-word;
  margin: 0;
}

.help-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.help-content ul {
  margin: 0;
  padding-left: 20px;
}

.help-content li {
  margin-bottom: 8px;
  color: #606266;
  line-height: 1.5;
}

.help-content li:last-child {
  margin-bottom: 0;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-start;
  }
  
  .header-actions .el-select {
    width: 100% !important;
    margin-right: 0 !important;
    margin-bottom: 8px;
  }
}
</style>
