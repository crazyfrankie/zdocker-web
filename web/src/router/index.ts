import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('../views/DashboardView.vue'),
    },
    {
      path: '/containers',
      name: 'containers',
      component: () => import('../views/ContainersView.vue'),
    },
    {
      path: '/containers/:id',
      name: 'container-detail',
      component: () => import('../views/ContainerDetailView.vue'),
    },
    {
      path: '/images',
      name: 'images',
      component: () => import('../views/ImagesView.vue'),
    },
    {
      path: '/networks',
      name: 'networks',
      component: () => import('../views/NetworksView.vue'),
    },
    {
      path: '/logs',
      name: 'logs',
      component: () => import('../views/LogsView.vue'),
    },
  ],
})

export default router
