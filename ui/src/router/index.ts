import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ProductView from '../views/ProductView.vue'
import CoverageView from '../views/CoverageView.vue'
import TestView from '../views/TestView.vue'
import NotFound from '../views/NotFound.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    alias: ['/index.html'],
    name: 'home',
    component: HomeView
  },
  {
    path: '/product',
    alias: ['/product.html'],
    name: 'product',
    component: ProductView
  },
  {
    path: '/coverage',
    alias: ['/coverage.html'],
    name: 'coverage',
    component: CoverageView
  },
  {
    path: '/tests',
    alias: ['/tests.html'],
    name: 'tests',
    component: TestView
  },
  {
    path: '/404',
    alias: ['/404.html'],
    name: '404',
    component: NotFound
  },
  {
    path: '/:pathMatch(.*)*',
    alias: ['/404.html'],
    name: '404',
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
