import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/queue'
  },
  {
    path: '/queue',
    name: 'Queue',
    component: () => import('../views/Queue.vue'),
    meta: { showTabBar: true }
  },
  {
    path: '/intro',
    name: 'Intro',
    component: () => import('../views/Intro.vue'),
    meta: { showTabBar: true }
  },
  {
    path: '/gallery',
    name: 'Gallery',
    component: () => import('../views/Gallery.vue'),
    meta: { showTabBar: true }
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router 