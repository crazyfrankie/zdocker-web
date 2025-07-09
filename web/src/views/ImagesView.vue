<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Picture, 
  Refresh, 
  Delete,
  Download
} from '@element-plus/icons-vue'

const loading = ref(true)
const images = ref<any[]>([])

onMounted(() => {
  loadImages()
})

const loadImages = async () => {
  try {
    loading.value = true
    // 由于zdocker没有专门的镜像管理，这里模拟一些常见的镜像
    images.value = [
      {
        id: 'ubuntu:latest',
        name: 'ubuntu',
        tag: 'latest',
        size: '72.8 MB',
        created: '2024-01-15T10:30:00Z'
      },
      {
        id: 'alpine:latest',
        name: 'alpine',
        tag: 'latest',
        size: '5.6 MB',
        created: '2024-01-10T08:20:00Z'
      },
      {
        id: 'busybox:latest',
        name: 'busybox',
        tag: 'latest',
        size: '1.2 MB',
        created: '2024-01-05T14:15:00Z'
      }
    ]
  } catch (error) {
    console.error('获取镜像列表失败:', error)
    ElMessage.error('获取镜像列表失败')
  } finally {
    loading.value = false
  }
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString()
}

const handlePull = () => {
  ElMessage.info('镜像拉取功能开发中...')
}

const handleRemove = (image: any) => {
  ElMessage.info(`删除镜像 ${image.name}:${image.tag} 功能开发中...`)
}
</script>

<template>
  <div class="images-view">
    <div class="page-header">
      <div class="header-left">
        <h2>镜像管理</h2>
        <p class="subtitle">管理容器镜像</p>
      </div>
      <div class="header-actions">
        <el-button 
          type="primary" 
          :icon="Download"
          @click="handlePull"
        >
          拉取镜像
        </el-button>
        <el-button 
          :icon="Refresh"
          @click="loadImages"
        >
          刷新
        </el-button>
      </div>
    </div>

    <el-card class="images-card">
      <el-table 
        :data="images" 
        v-loading="loading"
        style="width: 100%"
        empty-text="暂无镜像"
      >
        <el-table-column label="镜像" min-width="200">
          <template #default="{ row }">
            <div class="image-info">
              <el-icon class="image-icon"><Picture /></el-icon>
              <div class="image-details">
                <div class="image-name">{{ row.name }}:{{ row.tag }}</div>
                <div class="image-id">{{ row.id }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="size" label="大小" width="120" />
        
        <el-table-column prop="created" label="创建时间" width="160">
          <template #default="{ row }">
            <span>{{ formatTime(row.created) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="状态" width="100">
          <template #default>
            <el-tag type="success" size="small">可用</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              @click="handleRemove(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 镜像说明 -->
    <el-card class="help-card">
      <template #header>
        <span>镜像说明</span>
      </template>
      <div class="help-content">
        <el-alert
          title="提示"
          type="info"
          :closable="false"
          show-icon
        >
          <p>ZDocker 是一个简化的容器运行时实现，专注于学习容器技术原理。</p>
          <p>镜像管理功能目前处于开发阶段，上述镜像列表为示例数据。</p>
          <p>实际使用时，可以通过 zdocker run 命令直接使用系统中的文件系统作为容器根目录。</p>
        </el-alert>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.images-view {
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

.images-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
}

.image-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.image-icon {
  font-size: 24px;
  color: #409eff;
}

.image-details {
  flex: 1;
}

.image-name {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
}

.image-id {
  font-size: 12px;
  color: #909399;
  font-family: monospace;
  margin-top: 2px;
}

.help-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.help-content p {
  margin-bottom: 8px;
  line-height: 1.5;
}

.help-content p:last-child {
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
}
</style>
