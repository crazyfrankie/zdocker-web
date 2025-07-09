<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  getContainers, 
  startContainer, 
  stopContainer, 
  removeContainer,
  createContainer,
  type Container,
  type CreateContainerRequest
} from '@/api/containers'
import { 
  Plus, 
  Refresh, 
  VideoPlay, 
  VideoPause, 
  Delete,
  View,
  Setting
} from '@element-plus/icons-vue'

const loading = ref(true)
const containers = ref<Container[]>([])
const showCreateDialog = ref(false)
const createForm = ref<CreateContainerRequest>({
  image: '',
  command: '',
  name: '',
  detach: true,
  tty: false,
  volume: '',
  memory: '',
  cpu_share: '',
  cpu_set: '',
  network: '',
  environment: {},
  port_mapping: []
})
const envInput = ref('')
const portInput = ref('')

onMounted(() => {
  loadContainers()
})

const loadContainers = async () => {
  try {
    loading.value = true
    const response = await getContainers()
    containers.value = response.data || []
  } catch (error) {
    console.error('获取容器列表失败:', error)
    ElMessage.error('获取容器列表失败')
  } finally {
    loading.value = false
  }
}

const handleStart = async (container: Container) => {
  try {
    await startContainer(container.id)
    ElMessage.success(`容器 ${container.name || container.id} 启动成功`)
    loadContainers()
  } catch (error) {
    ElMessage.error('启动容器失败')
  }
}

const handleStop = async (container: Container) => {
  try {
    await stopContainer(container.id)
    ElMessage.success(`容器 ${container.name || container.id} 停止成功`)
    loadContainers()
  } catch (error) {
    ElMessage.error('停止容器失败')
  }
}

const handleRemove = async (container: Container) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除容器 ${container.name || container.id} 吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await removeContainer(container.id)
    ElMessage.success('容器删除成功')
    loadContainers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除容器失败')
    }
  }
}

const handleViewDetails = (container: Container) => {
  // 跳转到容器详情页面
  window.open(`/containers/${container.id}`, '_blank')
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

const showCreateForm = () => {
  showCreateDialog.value = true
  resetForm()
}

const resetForm = () => {
  createForm.value = {
    image: '',
    command: '',
    name: '',
    detach: true,
    tty: false,
    volume: '',
    memory: '',
    cpu_share: '',
    cpu_set: '',
    network: '',
    environment: {},
    port_mapping: []
  }
  envInput.value = ''
  portInput.value = ''
}

const addEnvironment = () => {
  if (!envInput.value.trim()) return
  
  const [key, value] = envInput.value.split('=')
  if (key && value) {
    createForm.value.environment![key.trim()] = value.trim()
    envInput.value = ''
  } else {
    ElMessage.warning('环境变量格式应为: KEY=VALUE')
  }
}

const removeEnvironment = (key: string) => {
  delete createForm.value.environment![key]
}

const addPortMapping = () => {
  if (!portInput.value.trim()) return
  
  if (!createForm.value.port_mapping) {
    createForm.value.port_mapping = []
  }
  
  createForm.value.port_mapping.push(portInput.value.trim())
  portInput.value = ''
}

const removePortMapping = (index: number) => {
  createForm.value.port_mapping?.splice(index, 1)
}

const handleCreate = async () => {
  try {
    if (!createForm.value.image.trim() || !createForm.value.command.trim()) {
      ElMessage.warning('请填写镜像名称和命令')
      return
    }
    
    await createContainer(createForm.value)
    ElMessage.success('容器创建成功')
    showCreateDialog.value = false
    loadContainers()
  } catch (error) {
    ElMessage.error('创建容器失败')
  }
}
</script>

<template>
  <div class="containers-view">
    <div class="page-header">
      <div class="header-left">
        <h2>容器管理</h2>
        <p class="subtitle">管理和监控你的容器</p>
      </div>
      <div class="header-actions">
        <el-button 
          type="primary" 
          :icon="Plus"
          @click="showCreateForm"
        >
          创建容器
        </el-button>
        <el-button 
          :icon="Refresh"
          @click="loadContainers"
        >
          刷新
        </el-button>
      </div>
    </div>

    <el-card class="containers-card">
      <el-table 
        :data="containers" 
        v-loading="loading"
        style="width: 100%"
        empty-text="暂无容器"
      >
        <el-table-column prop="name" label="名称" min-width="120">
          <template #default="{ row }">
            <span class="container-name">
              {{ row.name || row.id.substring(0, 12) }}
            </span>
          </template>
        </el-table-column>
        
        <el-table-column prop="image" label="镜像" min-width="120" />
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="command" label="命令" min-width="150" show-overflow-tooltip />
        
        <el-table-column prop="port_mapping" label="端口映射" width="120">
          <template #default="{ row }">
            <span>{{ row.port_mapping || '-' }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_time" label="创建时间" width="160">
          <template #default="{ row }">
            <span>{{ formatTime(row.created_time) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button
                v-if="row.status !== 'running'"
                type="success"
                size="small"
                :icon="VideoPlay"
                @click="handleStart(row)"
              >
                启动
              </el-button>
              
              <el-button
                v-if="row.status === 'running'"
                type="warning"
                size="small"
                :icon="VideoPause"
                @click="handleStop(row)"
              >
                停止
              </el-button>
              
              <el-button
                type="info"
                size="small"
                :icon="View"
                @click="handleViewDetails(row)"
              >
                详情
              </el-button>
              
              <el-button
                type="danger"
                size="small"
                :icon="Delete"
                @click="handleRemove(row)"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建容器对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      title="创建容器"
      width="600px"
      :before-close="() => { showCreateDialog = false }"
    >
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="镜像名称" required>
          <el-input 
            v-model="createForm.image" 
            placeholder="例如: ubuntu:latest"
          />
        </el-form-item>
        
        <el-form-item label="执行命令" required>
          <el-input 
            v-model="createForm.command" 
            placeholder="例如: /bin/bash"
          />
        </el-form-item>
        
        <el-form-item label="容器名称">
          <el-input 
            v-model="createForm.name" 
            placeholder="可选，不填写将自动生成"
          />
        </el-form-item>
        
        <el-form-item label="运行选项">
          <el-checkbox v-model="createForm.detach">后台运行 (-d)</el-checkbox>
          <el-checkbox v-model="createForm.tty">分配终端 (-t)</el-checkbox>
        </el-form-item>
        
        <el-form-item label="数据卷">
          <el-input 
            v-model="createForm.volume" 
            placeholder="例如: /host/path:/container/path"
          />
        </el-form-item>
        
        <el-form-item label="内存限制">
          <el-input 
            v-model="createForm.memory" 
            placeholder="例如: 512m, 1g"
          />
        </el-form-item>
        
        <el-form-item label="CPU限制">
          <el-input 
            v-model="createForm.cpu_share" 
            placeholder="CPU份额，例如: 512"
          />
        </el-form-item>
        
        <el-form-item label="网络">
          <el-input 
            v-model="createForm.network" 
            placeholder="网络名称，例如: bridge"
          />
        </el-form-item>
        
        <el-form-item label="环境变量">
          <div class="env-section">
            <div class="env-input">
              <el-input 
                v-model="envInput" 
                placeholder="KEY=VALUE"
                @keyup.enter="addEnvironment"
              />
              <el-button @click="addEnvironment">添加</el-button>
            </div>
            <div class="env-list">
              <el-tag
                v-for="(value, key) in createForm.environment"
                :key="key"
                closable
                @close="removeEnvironment(key)"
                style="margin: 4px;"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item label="端口映射">
          <div class="port-section">
            <div class="port-input">
              <el-input 
                v-model="portInput" 
                placeholder="8080:80"
                @keyup.enter="addPortMapping"
              />
              <el-button @click="addPortMapping">添加</el-button>
            </div>
            <div class="port-list">
              <el-tag
                v-for="(port, index) in createForm.port_mapping"
                :key="index"
                closable
                @close="removePortMapping(index)"
                style="margin: 4px;"
              >
                {{ port }}
              </el-tag>
            </div>
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="handleCreate">创建</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.containers-view {
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
  gap: 12px;
}

.containers-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.container-name {
  font-weight: 500;
  color: #303133;
}

.action-buttons {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.env-section, .port-section {
  width: 100%;
}

.env-input, .port-input {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.env-list, .port-list {
  min-height: 40px;
  border: 1px dashed #dcdfe6;
  border-radius: 4px;
  padding: 8px;
  background-color: #fafafa;
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
  
  .action-buttons {
    flex-direction: column;
  }
}
</style>
