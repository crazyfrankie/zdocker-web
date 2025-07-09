<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  getContainer, 
  getContainerLogs, 
  execContainer,
  startContainer,
  stopContainer,
  type Container 
} from '@/api/containers'
import { 
  VideoPlay, 
  VideoPause, 
  Refresh, 
  DocumentCopy,
  Monitor,
  Setting
} from '@element-plus/icons-vue'

const route = useRoute()
const containerId = route.params.id as string

const loading = ref(true)
const container = ref<Container | null>(null)
const logs = ref('')
const logsLoading = ref(false)
const activeTab = ref('info')
const execCommand = ref('')
const execOutput = ref('')
const execLoading = ref(false)

onMounted(() => {
  loadContainerDetails()
  if (activeTab.value === 'logs') {
    loadLogs()
  }
})

const loadContainerDetails = async () => {
  try {
    loading.value = true
    const response = await getContainer(containerId)
    container.value = response.data
  } catch (error) {
    console.error('获取容器详情失败:', error)
    ElMessage.error('获取容器详情失败')
  } finally {
    loading.value = false
  }
}

const loadLogs = async () => {
  try {
    logsLoading.value = true
    const response = await getContainerLogs(containerId)
    logs.value = response.data || '暂无日志'
  } catch (error) {
    console.error('获取容器日志失败:', error)
    logs.value = '获取日志失败'
  } finally {
    logsLoading.value = false
  }
}

const handleStart = async () => {
  try {
    await startContainer(containerId)
    ElMessage.success('容器启动成功')
    loadContainerDetails()
  } catch (error) {
    ElMessage.error('启动容器失败')
  }
}

const handleStop = async () => {
  try {
    await stopContainer(containerId)
    ElMessage.success('容器停止成功')
    loadContainerDetails()
  } catch (error) {
    ElMessage.error('停止容器失败')
  }
}

const executeCommand = async () => {
  if (!execCommand.value.trim()) {
    ElMessage.warning('请输入要执行的命令')
    return
  }
  
  try {
    execLoading.value = true
    const response = await execContainer(containerId, {
      command: execCommand.value.split(' ')
    })
    execOutput.value += `$ ${execCommand.value}\n${response.data.output}\n\n`
    execCommand.value = ''
  } catch (error) {
    execOutput.value += `$ ${execCommand.value}\n执行失败\n\n`
    ElMessage.error('执行命令失败')
  } finally {
    execLoading.value = false
  }
}

const handleTabChange = (tabName: string) => {
  if (tabName === 'logs' && !logs.value) {
    loadLogs()
  }
}

const getStatusType = (status: string) => {
  switch (status) {
    case 'running':
      return 'success'
    case 'stopped':
      return 'info'
    case 'error':
      return 'danger'
    default:
      return 'warning'
  }
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString()
}
</script>

<template>
  <div class="container-detail" v-loading="loading">
    <div class="page-header" v-if="container">
      <div class="header-left">
        <h2>{{ container.name || container.id.substring(0, 12) }}</h2>
        <div class="header-info">
          <el-tag :type="getStatusType(container.status)" size="large">
            {{ container.status }}
          </el-tag>
          <span class="container-id">ID: {{ container.id }}</span>
        </div>
      </div>
      <div class="header-actions">
        <el-button
          v-if="container.status !== 'running'"
          type="success"
          :icon="VideoPlay"
          @click="handleStart"
        >
          启动
        </el-button>
        <el-button
          v-if="container.status === 'running'"
          type="warning"
          :icon="VideoPause"
          @click="handleStop"
        >
          停止
        </el-button>
        <el-button :icon="Refresh" @click="loadContainerDetails">
          刷新
        </el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange" v-if="container">
      <!-- 基本信息 -->
      <el-tab-pane label="基本信息" name="info">
        <el-card class="info-card">
          <div class="info-grid">
            <div class="info-item">
              <span class="label">容器名称:</span>
              <span class="value">{{ container.name || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">容器ID:</span>
              <span class="value">{{ container.id }}</span>
            </div>
            <div class="info-item">
              <span class="label">镜像:</span>
              <span class="value">{{ container.image }}</span>
            </div>
            <div class="info-item">
              <span class="label">状态:</span>
              <el-tag :type="getStatusType(container.status)">
                {{ container.status }}
              </el-tag>
            </div>
            <div class="info-item">
              <span class="label">进程ID:</span>
              <span class="value">{{ container.pid || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">创建时间:</span>
              <span class="value">{{ formatTime(container.created_time) }}</span>
            </div>
            <div class="info-item">
              <span class="label">启动命令:</span>
              <span class="value">{{ container.command }}</span>
            </div>
            <div class="info-item">
              <span class="label">数据卷:</span>
              <span class="value">{{ container.volume || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">端口映射:</span>
              <span class="value">{{ container.port_mapping || '-' }}</span>
            </div>
          </div>
        </el-card>
      </el-tab-pane>

      <!-- 日志 -->
      <el-tab-pane label="日志" name="logs">
        <el-card class="logs-card">
          <template #header>
            <div class="logs-header">
              <span>容器日志</span>
              <el-button 
                size="small" 
                :icon="Refresh" 
                @click="loadLogs"
                :loading="logsLoading"
              >
                刷新日志
              </el-button>
            </div>
          </template>
          <div class="logs-content" v-loading="logsLoading">
            <pre class="logs-text">{{ logs }}</pre>
          </div>
        </el-card>
      </el-tab-pane>

      <!-- 控制台 -->
      <el-tab-pane label="控制台" name="console" v-if="container.status === 'running'">
        <el-card class="console-card">
          <template #header>
            <span>执行命令</span>
          </template>
          <div class="console-content">
            <div class="command-input">
              <el-input
                v-model="execCommand"
                placeholder="输入要执行的命令，例如: ls -la"
                @keyup.enter="executeCommand"
                :disabled="execLoading"
              >
                <template #append>
                  <el-button 
                    @click="executeCommand"
                    :loading="execLoading"
                  >
                    执行
                  </el-button>
                </template>
              </el-input>
            </div>
            <div class="command-output">
              <pre class="output-text">{{ execOutput || '输入命令并按回车执行...' }}</pre>
            </div>
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.container-detail {
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

.header-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.container-id {
  color: #909399;
  font-size: 14px;
  font-family: monospace;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.info-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.label {
  font-weight: 500;
  color: #606266;
  min-width: 100px;
}

.value {
  color: #303133;
  word-break: break-all;
  text-align: right;
  flex: 1;
  margin-left: 16px;
}

.logs-card, .console-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logs-content {
  height: 400px;
  overflow-y: auto;
  background-color: #f8f9fa;
  border-radius: 4px;
  padding: 16px;
}

.logs-text {
  font-family: 'Courier New', Courier, monospace;
  font-size: 12px;
  line-height: 1.4;
  color: #333;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.console-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.command-output {
  height: 300px;
  overflow-y: auto;
  background-color: #1e1e1e;
  border-radius: 4px;
  padding: 16px;
}

.output-text {
  font-family: 'Courier New', Courier, monospace;
  font-size: 12px;
  line-height: 1.4;
  color: #d4d4d4;
  white-space: pre-wrap;
  word-wrap: break-word;
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
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .info-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .value {
    text-align: left;
    margin-left: 0;
  }
}
</style>
