import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import DockerView from '../views/Docker.vue'
import CloudflareView from '../views/Cloudflare.vue'
import ApplicationView from '../views/Application.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Docker',
    component: DockerView
  },
  {
    path: '/cloudflare',
    name: 'Cloudflare',
    component: CloudflareView
  },
  {
    path: '/cloudflare/:id',
    name: 'CloudflareApplication',
    component: ApplicationView
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router