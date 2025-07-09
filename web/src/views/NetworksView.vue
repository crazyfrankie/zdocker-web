<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  getNetworks, 
  createNetwork, 
  removeNetwork,
  type NetworkInfo,
  type CreateNetworkRequest
} from '@/api/networks'
import { 
  Plus, 
  Refresh, 
  Delete,
  Connection
} from '@element-plus/icons-vue'

const loading = ref(true)
const networks = ref<NetworkInfo[]>([])
const showCreateDialog = ref(false)
const createForm = ref<CreateNetworkRequest>({
  name: '',
  driver: 'bridge',
  subnet: ''
})

onMounted(() => {
  loadNetworks()
})

const loadNetworks = async () => {
  try {
    loading.value = true
    const response = await getNetworks()
    networks.value = response.data || []
  } catch (error) {
    console.error('获取网络列表失败:', error)
    ElMessage.error('获取网络列表失败')
  } finally {
    loading.value = false
  }
}

const showCreateForm = () => {
  showCreateDialog.value = true
  resetForm()
}

const resetForm = () => {
  createForm.value = {
    name: '',
    driver: 'bridge',
    subnet: ''
  }
}

const handleCreate = async () => {
  try {
    if (!createForm.value.name.trim()) {
      ElMessage.warning('请填写网络名称')
      return
    }
    
    await createNetwork(createForm.value)
    ElMessage.success('网络创建成功')
    showCreateDialog.value = false
    loadNetworks()
  } catch (error) {
    ElMessage.error('创建网络失败')
  }
}

const handleRemove = async (network: NetworkInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除网络 ${network.name} 吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await removeNetwork(network.name)
    ElMessage.success('网络删除成功')
    loadNetworks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除网络失败')
    }
  }
}

const getDriverColor = (driver: string) => {
  switch (driver) {
    case 'bridge':
      return 'primary'
    case 'host':
      return 'success'
    case 'none':
      return 'info'
    default:
      return 'warning'
  }
}
</script>

<template>
  <div class="networks-view">
    <div class="page-header">
      <div class="header-left">
        <h2>网络管理</h2>
        <p class="subtitle">管理容器网络配置</p>
      </div>
      <div class="header-actions">
        <el-button 
          type="primary" 
          :icon="Plus"
          @click="showCreateForm"
        >
          创建网络
        </el-button>
        <el-button 
          :icon="Refresh"
          @click="loadNetworks"
        >
          刷新
        </el-button>
      </div>
    </div>

    <el-card class="networks-card">
      <el-table 
        :data="networks" 
        v-loading="loading"
        style="width: 100%"
        empty-text="暂无网络"
      >
        <el-table-column prop="name" label="网络名称" min-width="150">
          <template #default="{ row }">
            <div class="network-name">
              <el-icon class="network-icon"><Connection /></el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="driver" label="驱动类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getDriverColor(row.driver)" size="small">
              {{ row.driver }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="subnet" label="子网" min-width="150">
          <template #default="{ row }">
            <span class="subnet">{{ row.subnet || '-' }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="状态" width="100">
          <template #default>
            <el-tag type="success" size="small">活跃</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.name !== 'bridge'"
              type="danger"
              size="small"
              :icon="Delete"
              @click="handleRemove(row)"
            >
              删除
            </el-button>
            <el-tooltip 
              v-else
              content="默认网络不能删除"
              placement="top"
            >
              <el-button
                type="danger"
                size="small"
                :icon="Delete"
                disabled
              >
                删除
              </el-button>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建网络对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      title="创建网络"
      width="500px"
      :before-close="() => { showCreateDialog = false }"
    >
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="网络名称" required>
          <el-input 
            v-model="createForm.name" 
            placeholder="输入网络名称"
          />
        </el-form-item>
        
        <el-form-item label="驱动类型">
          <el-select v-model="createForm.driver" style="width: 100%">
            <el-option label="Bridge" value="bridge" />
            <el-option label="Host" value="host" />
            <el-option label="None" value="none" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="子网">
          <el-input 
            v-model="createForm.subnet" 
            placeholder="例如: 172.18.0.0/16"
          />
          <div class="form-help">
            可选，如果不指定将自动分配
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

    <!-- 网络说明 -->
    <el-card class="help-card">
      <template #header>
        <span>网络类型说明</span>
      </template>
      <div class="help-content">
        <div class="help-item">
          <h4>Bridge 网络</h4>
          <p>默认网络类型，容器通过虚拟网桥进行通信，提供网络隔离。</p>
        </div>
        <div class="help-item">
          <h4>Host 网络</h4>
          <p>容器直接使用主机网络栈，性能最好但没有网络隔离。</p>
        </div>
        <div class="help-item">
          <h4>None 网络</h4>
          <p>容器没有网络接口，完全隔离的网络环境。</p>
        </div>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.networks-view {
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

.networks-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
}

.network-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.network-icon {
  color: #409eff;
}

.subnet {
  font-family: monospace;
  color: #606266;
}

.form-help {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.help-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.help-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.help-item h4 {
  color: #303133;
  margin-bottom: 8px;
  font-size: 16px;
}

.help-item p {
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
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
  
  .help-content {
    grid-template-columns: 1fr;
  }
}
</style>
