import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/Login.vue'
import { useDataStore } from '@/stores/dataStore'


const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      alias: ['/index.html'],
      name: 'home',
      component: HomeView
    },
    {
      path: '/collectors',
      name: 'collectors',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('@/views/CollectorView.vue')
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    },
    {
      path: '/inventory',
      name: 'inventory',
      component: () => import('@/views/InventoryView.vue')
    },
    {
    path: '/login',
      component: Login,
      beforeEnter: (to, from, next) => {
        const dataStore = useDataStore()
        if (dataStore.isLoggedIn) {
          next('/') // redirect if already logged in
        } else {
          next()
        }
      }
    }
  ]
})

export default router
