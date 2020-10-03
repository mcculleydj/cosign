import Vue from 'vue'
import VueRouter from 'vue-router'
import Members from '@/views/Members'

Vue.use(VueRouter)

export const routes = [
  {
    path: '/',
    name: 'members',
    component: Members,
    icon: 'mdi-account-multiple',
    text: 'Cosponsors',
  },
  {
    path: '/bills',
    name: 'bills',
    component: () => import('@/views/Legislation'),
    icon: 'mdi-script-text-outline',
    text: 'Legislation',
  },
  {
    path: '/visualizations',
    name: 'visualizations',
    component: () => import('@/views/Visualizations'),
    icon: 'mdi-chart-arc',
    text: 'Visualizations',
  },
  {
    path: '/sandbox',
    name: 'sandbox',
    component: () => import('@/views/Sandbox'),
    icon: 'mdi-test-tube',
    text: 'Sandbox',
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
})

export default router
