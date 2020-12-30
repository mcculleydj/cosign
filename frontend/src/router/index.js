import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '@/views/Home'

Vue.use(VueRouter)

export const routes = [
  {
    path: '/',
    name: 'home',
    component: Home,
    icon: 'mdi-home',
    text: 'Home',
  },
  {
    path: '/cosponsors',
    name: 'cosponsors',
    component: () => import('@/views/Cosponsors'),
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
  // TODO: 404 page
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
})

export default router
