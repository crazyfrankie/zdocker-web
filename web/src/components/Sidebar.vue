<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  Box, // 替换 Container 为 Box
  Picture, 
  Connection, 
  Monitor, 
  Setting,
  DocumentCopy,
  Histogram
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const menuItems = ref([
  {
    title: '概览',
    icon: Monitor,
    path: '/',
    name: 'dashboard'
  },
  {
    title: '容器',
    icon: Box, // 替换 Container 为 Box
    path: '/containers',
    name: 'containers'
  },
  {
    title: '镜像',
    icon: Picture,
    path: '/images',
    name: 'images'
  },
  {
    title: '网络',
    icon: Connection,
    path: '/networks',
    name: 'networks'
  },
  {
    title: '日志',
    icon: DocumentCopy,
    path: '/logs',
    name: 'logs'
  }
])

const handleMenuClick = (item: any) => {
  router.push(item.path)
}

const isActive = (path: string) => {
  return route.path === path || (path !== '/' && route.path.startsWith(path))
}
</script>

<template>
  <div class="sidebar">
    <nav class="sidebar-nav">
      <div 
        v-for="item in menuItems" 
        :key="item.name"
        class="nav-item"
        :class="{ active: isActive(item.path) }"
        @click="handleMenuClick(item)"
      >
        <el-icon class="nav-icon">
          <component :is="item.icon" />
        </el-icon>
        <span class="nav-text">{{ item.title }}</span>
      </div>
    </nav>
  </div>
</template>

<style scoped>
.sidebar {
  width: 240px;
  background: #fff;
  border-right: 1px solid #e8e8e8;
  height: 100%;
  box-shadow: 2px 0 6px rgba(0,0,0,0.1);
}

.sidebar-nav {
  padding: 20px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.3s;
  color: #666;
  margin: 2px 12px;
  border-radius: 8px;
}

.nav-item:hover {
  background-color: #f5f5f5;
  color: #2196F3;
}

.nav-item.active {
  background-color: #e3f2fd;
  color: #2196F3;
  font-weight: 500;
}

.nav-icon {
  font-size: 18px;
  margin-right: 12px;
}

.nav-text {
  font-size: 14px;
}
</style>
