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
import { ref, onMounted, onUnmounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useSettingsStore } from '../stores/settings';
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
  MoonIcon
} from '@heroicons/vue/24/outline';

defineProps({
  cartCount: { type: Number, default: 0 }
});

const settingsStore = useSettingsStore();
const { settings: storeSetting } = storeToRefs(settingsStore);
const { isDarkMode, toggleTheme } = useTheme();

const navItems = [
  { id: 'register', to: '/', label: 'Kasir', iconComp: BuildingStorefrontIcon },
  { id: 'products', to: '/products', label: 'Produk', iconComp: CubeIcon },
  { id: 'inventory', to: '/inventory', label: 'Inventoris', iconComp: ArchiveBoxIcon },
  { id: 'customers', to: '/customers', label: 'Pelanggan', iconComp: UsersIcon },
  { id: 'orders', to: '/orders', label: 'Riwayat', iconComp: DocumentTextIcon },
  { id: 'reports', to: '/reports', label: 'Laporan', iconComp: ChartBarIcon },
  { id: 'settings', to: '/settings', label: 'Pengaturan', iconComp: Cog6ToothIcon }
];

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

<style scoped>
.header-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1.5rem;
  z-index: 50;
  height: 64px;
  flex-shrink: 0;
  background: var(--bg-glass);
  border-bottom: 1px solid var(--border-color);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.logo-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: #eef2ff;
  border: 1px solid #c7d2fe;
  display: flex;
  align-items: center;
  justify-content: center;
}
html.dark .logo-icon {
  background: #1e1b4b;
  border-color: #3730a3;
}

.store-name {
  font-size: 1.05rem;
  font-weight: 800;
  letter-spacing: -0.02em;
  color: var(--text-primary);
  line-height: 1.2;
}

.db-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.12rem 0.5rem;
  border-radius: 999px;
  font-size: 0.65rem;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.db-sqlite {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}
html.dark .db-sqlite {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
  border-color: rgba(59, 130, 246, 0.3);
}

.db-mysql {
  background: #fffbeb;
  color: #b45309;
  border: 1px solid #fde68a;
}
html.dark .db-mysql {
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
  border-color: rgba(245, 158, 11, 0.3);
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

/* Desktop Nav */
.desktop-nav {
  display: none;
  gap: 0.25rem;
  background: var(--bg-primary);
  padding: 0.25rem;
  border-radius: 12px;
  border: 1px solid var(--border-color);
}

@media (min-width: 768px) {
  .desktop-nav {
    display: flex;
  }
}

.nav-btn {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.45rem 0.85rem;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-family: var(--font-family);
  font-size: 0.825rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  text-decoration: none;
}

.nav-btn:hover {
  color: var(--text-primary);
  background: var(--bg-card);
}

.nav-btn.active {
  background: var(--bg-card);
  color: var(--accent-primary);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.theme-toggle-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
}
.theme-toggle-btn:hover {
  background: var(--bg-card-hover);
}

.clock-display {
  padding: 0.35rem 0.75rem;
  background: var(--bg-primary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
  font-weight: 600;
  font-size: 0.8rem;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
}

/* Mobile Bottom Nav */
.mobile-bottom-nav {
  display: flex;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: var(--bg-glass);
  backdrop-filter: blur(16px);
  border-top: 1px solid var(--border-color);
  z-index: 100;
  align-items: center;
  justify-content: space-around;
  box-shadow: 0 -4px 10px rgba(0, 0, 0, 0.05);
}

@media (min-width: 768px) {
  .mobile-bottom-nav {
    display: none;
  }
}

.mobile-nav-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.15rem;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-family: var(--font-family);
  font-size: 0.65rem;
  font-weight: 600;
  cursor: pointer;
  flex: 1;
  height: 100%;
  position: relative;
  text-decoration: none;
}

.mobile-nav-btn.active {
  color: var(--accent-primary);
}

.cart-badge {
  position: absolute;
  top: 6px;
  right: calc(50% - 18px);
  background: var(--accent-danger);
  color: #ffffff;
  font-size: 0.65rem;
  font-weight: 800;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
