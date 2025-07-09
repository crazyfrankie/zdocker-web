<script setup lang="ts">
import { RouterView } from 'vue-router'
import Sidebar from './components/Sidebar.vue'
import Header from './components/Header.vue'
</script>

<template>
  <div class="app-container">
    <Header />
    <div class="main-content">
      <Sidebar />
      <div class="content-area">
        <RouterView />
      </div>
    </div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f5f5;
}

.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.content-area {
  flex: 1;
  width: 100%;
  padding: 20px;
  overflow-y: auto;
  background-color: #f5f5f5;
}

/* 确保所有页面组件都有统一的宽度 */
.content-area > * {
  width: 100%;
  max-width: none;
}

/* 页面容器统一样式 */
.page-container {
  width: 100%;
  margin: 0 auto;
}

/* 页面标题统一样式 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  width: 100%;
}

.page-header h2 {
  font-size: 28px;
  color: #303133;
  margin-bottom: 8px;
}

.page-header .subtitle {
  color: #909399;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

/* 卡片统一样式 */
.page-card {
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
  width: 100%;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-start;
  }
  
  .content-area {
    padding: 12px;
  }
}
</style>
