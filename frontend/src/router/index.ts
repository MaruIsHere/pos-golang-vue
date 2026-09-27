import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('../views/LoginView.vue') },
  { path: '/register', name: 'user-register', component: () => import('../views/RegisterView.vue') },
  { 
    path: '/', 
    name: 'pos', 
    component: () => import('../views/PosView.vue'), 
    meta: { requiresAuth: true, roles: ['kasir', 'kepala_kasir', 'owner', 'admin'] } 
  },
  { 
    path: '/orders', 
    name: 'orders', 
    component: () => import('../views/OrdersView.vue'), 
    meta: { requiresAuth: true, roles: ['kasir', 'kepala_kasir', 'owner', 'admin'] } 
  },
  { 
    path: '/products', 
    name: 'products', 
    component: () => import('../views/ProductsView.vue'), 
    meta: { requiresAuth: true, roles: ['kepala_kasir', 'owner', 'admin'] } 
  },
  { 
    path: '/inventory', 
    name: 'inventory', 
    component: () => import('../views/InventoryView.vue'), 
    meta: { requiresAuth: true, roles: ['kepala_kasir', 'owner', 'admin'] } 
  },
  { 
    path: '/customers', 
    name: 'customers', 
    component: () => import('../views/CustomersView.vue'), 
    meta: { requiresAuth: true, roles: ['kepala_kasir', 'owner', 'admin'] } 
  },
  { 
    path: '/reports', 
    name: 'reports', 
    component: () => import('../views/ReportsView.vue'), 
    meta: { requiresAuth: true, roles: ['kepala_kasir', 'owner', 'admin'] } 
  },
  { 
    path: '/settings', 
    name: 'settings', 
    component: () => import('../views/SettingsView.vue'), 
    meta: { requiresAuth: true, roles: ['owner', 'admin'] } 
  },
  { path: '/:pathMatch(.*)*', redirect: '/' },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();

  // 1. Check for invalid or corrupted auth state in localStorage
  if (authStore.isAuthenticated && (!authStore.user || !authStore.userRole)) {
    authStore.logout();
    if (to.path !== '/login') {
      next('/login');
    } else {
      next();
    }
    return;
  }

  // 2. Unauthenticated user trying to access protected route
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    if (to.path !== '/login') {
      next('/login');
    } else {
      next();
    }
    return;
  }

  // 3. Authenticated user trying to access guest routes (/login or /register)
  if ((to.path === '/login' || to.path === '/register') && authStore.isAuthenticated) {
    const role = authStore.userRole;
    const targetPath = (role === 'owner' || role === 'admin') ? '/reports' : '/';
    next(targetPath);
    return;
  }

  // 4. Role checking for protected routes
  if (to.meta.requiresAuth && to.meta.roles) {
    const allowedRoles = to.meta.roles as string[];
    const currentRole = authStore.userRole || '';

    if (!allowedRoles.includes(currentRole)) {
      let fallbackPath = '/';
      if (currentRole === 'owner' || currentRole === 'admin') {
        fallbackPath = '/reports';
      }
      
      if (to.path !== fallbackPath) {
        next(fallbackPath);
      } else {
        next();
      }
      return;
    }
  }

  next();
});

export default router;
