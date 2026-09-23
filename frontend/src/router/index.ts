import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';

// Lazy-loaded views — 1 route per tab (pengganti currentTab di App.vue lama).
const routes: RouteRecordRaw[] = [
  { path: '/', name: 'register', component: () => import('../views/RegisterView.vue') },
  { path: '/products', name: 'products', component: () => import('../views/ProductsView.vue') },
  { path: '/inventory', name: 'inventory', component: () => import('../views/InventoryView.vue') },
  { path: '/customers', name: 'customers', component: () => import('../views/CustomersView.vue') },
  { path: '/orders', name: 'orders', component: () => import('../views/OrdersView.vue') },
  { path: '/reports', name: 'reports', component: () => import('../views/ReportsView.vue') },
  { path: '/settings', name: 'settings', component: () => import('../views/SettingsView.vue') },
  // TODO(auth): tambah { path: '/login', ... } + beforeEach guard (lihat docs/TODO.md item 4).
  { path: '/:pathMatch(.*)*', redirect: '/' },
];

const router = createRouter({
  // Backend Gin sudah serve ../frontend/dist/index.html via NoRoute,
  // jadi history mode aman untuk production single-binary.
  history: createWebHistory(),
  routes,
});

export default router;
