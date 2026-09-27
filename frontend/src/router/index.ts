import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('../views/LoginView.vue') },
  { path: '/register', name: 'user-register', component: () => import('../views/RegisterView.vue') },
  { path: '/', name: 'pos', component: () => import('../views/PosView.vue'), meta: { requiresAuth: true } },
  { path: '/products', name: 'products', component: () => import('../views/ProductsView.vue'), meta: { requiresAuth: true } },
  { path: '/inventory', name: 'inventory', component: () => import('../views/InventoryView.vue'), meta: { requiresAuth: true } },
  { path: '/customers', name: 'customers', component: () => import('../views/CustomersView.vue'), meta: { requiresAuth: true } },
  { path: '/orders', name: 'orders', component: () => import('../views/OrdersView.vue'), meta: { requiresAuth: true } },
  { path: '/reports', name: 'reports', component: () => import('../views/ReportsView.vue'), meta: { requiresAuth: true } },
  { path: '/settings', name: 'settings', component: () => import('../views/SettingsView.vue'), meta: { requiresAuth: true } },
  { path: '/:pathMatch(.*)*', redirect: '/' },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})

export default router;
