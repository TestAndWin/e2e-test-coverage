import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import ProductView from '../views/ProductView.vue';
import CoverageView from '../views/CoverageView.vue';
import TestView from '../views/TestView.vue';
import LogInView from '../views/LogInView.vue';
import LogOutView from '../views/LogOutView.vue';
import NotFound from '../views/NotFound.vue';
import AdminView from '../views/AdminView.vue';
import MyAccountView from '../views/MyAccountView.vue';

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
    path: '/login',
    alias: ['/login.html'],
    name: 'login',
    component: LogInView
  },
  {
    path: '/logout',
    alias: ['/logout.html'],
    name: 'logout',
    component: LogOutView
  },
  {
    path: '/admin',
    alias: ['/admin.html'],
    name: 'admin',
    component: AdminView
  },
  {
    path: '/myaccount',
    alias: ['/myaccount.html'],
    name: 'myaccount',
    component: MyAccountView
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
