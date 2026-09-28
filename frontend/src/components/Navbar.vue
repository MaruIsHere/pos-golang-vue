<template>
  <header class="glass-nav header-container">
    <div class="header-left">
      <div class="brand-logo">
        <div class="logo-icon">
          <ShoppingCartIcon class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
        </div>
        <div class="brand-info">
          <h1 class="store-name">{{ storeSetting.store_name || 'POS KASIR PRO' }}</h1>
          <span class="db-badge" :class="storeSetting.db_engine === 'mysql' ? 'db-mysql' : 'db-sqlite'">
            <span class="dot"></span> {{ (storeSetting.db_engine || 'sqlite').toUpperCase() }}
          </span>
        </div>
      </div>
    </div>

    <!-- Desktop & Tablet Navigation -->
    <nav class="desktop-nav">
      <RouterLink
        v-for="item in navItems"
        :key="item.id"
        :to="item.to"
        class="nav-btn"
        active-class="active"
      >
        <component :is="item.iconComp" class="w-4 h-4 nav-icon-svg" />
        <span class="nav-text">{{ item.label }}</span>
      </RouterLink>
    </nav>

    <div class="header-right">
      <!-- User Profile & Role Badge -->
      <div v-if="authStore.user" class="hidden sm:flex items-center gap-2 px-3 py-1 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg">
        <UserIcon class="w-4 h-4 text-accent-primary" />
        <div class="flex flex-col text-left leading-tight">
          <span class="text-xs font-semibold text-slate-900 dark:text-slate-100 leading-tight">{{ authStore.user.username }}</span>
          <span class="text-[9px] font-extrabold uppercase tracking-wider px-1.5 py-0.5 rounded text-center mt-0.5" :class="roleBadgeClass">
            {{ roleLabel }}
          </span>
        </div>
      </div>

      <!-- Theme Toggle Button -->
      <button 
        class="theme-toggle-btn" 
        :title="isDarkMode ? 'Mode Terang' : 'Mode Gelap'"
        @click="toggleTheme"
      >
        <SunIcon v-if="isDarkMode" class="w-4 h-4 text-amber-400" />
        <MoonIcon v-else class="w-4 h-4 text-slate-600" />
      </button>

      <div class="clock-display">
        <ClockIcon class="w-4 h-4 inline-block mr-1 text-slate-400" />
        <span class="time-text">{{ currentTime }}</span>
      </div>

      <button
        @click="logout"
        class="ml-1 px-3 py-1.5 text-xs font-semibold text-accent-danger bg-white dark:bg-slate-800 hover:bg-red-500/10 rounded-lg transition-colors border border-slate-200 dark:border-slate-700"
      >
        Logout
      </button>
    </div>
  </header>

  <!-- Mobile Bottom Navigation Bar -->
  <nav class="mobile-bottom-nav">
    <RouterLink
      v-for="item in navItems"
      :key="item.id"
      :to="item.to"
      class="mobile-nav-btn"
      active-class="active"
    >
      <component :is="item.iconComp" class="w-5 h-5 nav-icon-svg" />
      <span class="nav-text">{{ item.label }}</span>
      <span v-if="item.id === 'register' && cartCount > 0" class="cart-badge">{{ cartCount }}</span>
    </RouterLink>
  </nav>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useSettingsStore } from '../stores/settings';
import { useAuthStore } from '../stores/auth';
import { useTheme } from '../composables/useTheme';
import {
  ShoppingCartIcon,
  BuildingStorefrontIcon,
  CubeIcon,
  ArchiveBoxIcon,
  UsersIcon,
  DocumentTextIcon,
  ChartBarIcon,
  Cog6ToothIcon,
  ClockIcon,
  SunIcon,
  MoonIcon,
  UserIcon
} from '@heroicons/vue/24/outline';

defineProps({
  cartCount: { type: Number, default: 0 }
});

const settingsStore = useSettingsStore();
const { settings: storeSetting } = storeToRefs(settingsStore);
const authStore = useAuthStore();
const { isDarkMode, toggleTheme } = useTheme();

const logout = () => {
  authStore.logout();
};

const allNavItems = [
  { id: 'register', to: '/', label: 'Kasir', iconComp: BuildingStorefrontIcon, roles: ['kasir', 'kepala_kasir', 'owner', 'admin'] },
  { id: 'products', to: '/products', label: 'Produk', iconComp: CubeIcon, roles: ['kepala_kasir', 'owner', 'admin'] },
  { id: 'inventory', to: '/inventory', label: 'Inventoris', iconComp: ArchiveBoxIcon, roles: ['kepala_kasir', 'owner', 'admin'] },
  { id: 'customers', to: '/customers', label: 'Pelanggan', iconComp: UsersIcon, roles: ['kepala_kasir', 'owner', 'admin'] },
  { id: 'orders', to: '/orders', label: 'Riwayat', iconComp: DocumentTextIcon, roles: ['kasir', 'kepala_kasir', 'owner', 'admin'] },
  { id: 'reports', to: '/reports', label: 'Laporan', iconComp: ChartBarIcon, roles: ['kepala_kasir', 'owner', 'admin'] },
  { id: 'settings', to: '/settings', label: 'Pengaturan', iconComp: Cog6ToothIcon, roles: ['owner', 'admin'] }
];

const navItems = computed(() => {
  const currentRole = authStore.userRole || 'kasir';
  return allNavItems.filter(item => item.roles.includes(currentRole));
});

const roleLabel = computed(() => {
  const role = authStore.userRole;
  if (role === 'owner') return 'Owner';
  if (role === 'admin') return 'Admin';
  if (role === 'kepala_kasir') return 'Kepala Kasir';
  return 'Kasir';
});

const roleBadgeClass = computed(() => {
  const role = authStore.userRole;
  if (role === 'owner' || role === 'admin') {
    return 'bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-300 border border-amber-300/50';
  }
  if (role === 'kepala_kasir') {
    return 'bg-purple-100 text-purple-800 dark:bg-purple-900/40 dark:text-purple-300 border border-purple-300/50';
  }
  return 'bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-300 border border-blue-300/50';
});

const currentTime = ref('');

const updateTime = (): void => {
  const now = new Date();
  currentTime.value = now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
};

let timer: ReturnType<typeof setInterval> | null = null;
onMounted(() => {
  updateTime();
  timer = setInterval(updateTime, 1000);
});

onUnmounted(() => {
  if (timer) clearInterval(timer);
});
</script>


