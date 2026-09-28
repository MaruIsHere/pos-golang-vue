<template>
  <header class="sticky top-0 z-40 w-full flex items-center justify-between px-4 py-3 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 shadow-sm transition-colors duration-200">
    
    <!-- Bagian Kiri: Logo & Info Store -->
    <div class="flex items-center gap-3">
      <div class="w-10 h-10 rounded-xl bg-indigo-50 dark:bg-indigo-900/30 flex items-center justify-center border border-indigo-100 dark:border-indigo-800/50 shadow-sm">
        <ShoppingCartIcon class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
      </div>
      <div class="flex flex-col">
        <h1 class="text-sm sm:text-base font-bold text-slate-800 dark:text-slate-100 leading-tight">
          {{ storeSetting.store_name || 'POS KASIR PRO' }}
        </h1>
        <span 
          class="inline-flex items-center gap-1.5 px-1.5 py-0.5 mt-0.5 rounded text-[10px] font-bold uppercase tracking-wider w-fit"
          :class="storeSetting.db_engine === 'mysql' ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-400 border border-amber-200 dark:border-amber-800/50' : 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800/50'"
        >
          <span class="w-1.5 h-1.5 rounded-full" :class="storeSetting.db_engine === 'mysql' ? 'bg-amber-500 animate-pulse' : 'bg-emerald-500'"></span> 
          {{ (storeSetting.db_engine || 'sqlite').toUpperCase() }}
        </span>
      </div>
    </div>

    <!-- Bagian Tengah: Desktop Navigation -->
    <nav class="hidden md:flex items-center gap-1.5">
      <RouterLink
        v-for="item in navItems"
        :key="item.id"
        :to="item.to"
        class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium text-slate-500 hover:text-slate-900 hover:bg-slate-100 dark:text-slate-400 dark:hover:text-slate-100 dark:hover:bg-slate-800 transition-all duration-200"
        active-class="bg-slate-100 dark:bg-slate-800 !text-indigo-600 dark:!text-indigo-400 shadow-sm"
      >
        <component :is="item.iconComp" class="w-4 h-4 shrink-0" />
        <span>{{ item.label }}</span>
      </RouterLink>
    </nav>

    <!-- Bagian Kanan: Aksi & Profil -->
    <div class="flex items-center gap-2 sm:gap-3">
      <!-- User Profile (Hidden di HP super kecil) -->
      <div v-if="authStore.user" class="hidden sm:flex items-center gap-2 px-3 py-1.5 bg-slate-50 dark:bg-slate-800/50 border border-slate-200 dark:border-slate-700 rounded-lg">
        <UserIcon class="w-4 h-4 text-indigo-500" />
        <div class="flex flex-col leading-tight">
          <span class="text-xs font-semibold text-slate-900 dark:text-slate-100">{{ authStore.user.username }}</span>
          <span class="text-[9px] font-extrabold uppercase tracking-wider mt-0.5 w-fit" :class="roleBadgeClass">
            {{ roleLabel }}
          </span>
        </div>
      </div>

      <!-- Jam Digital (Hidden di HP) -->
      <div class="hidden lg:flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700">
        <ClockIcon class="w-4 h-4 text-slate-400" />
        <span class="text-xs font-semibold text-slate-600 dark:text-slate-300 font-mono">{{ currentTime }}</span>
      </div>

      <!-- Tombol Tema -->
      <button 
        class="p-2 rounded-lg text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors" 
        :title="isDarkMode ? 'Mode Terang' : 'Mode Gelap'"
        @click="toggleTheme"
      >
        <SunIcon v-if="isDarkMode" class="w-5 h-5 text-amber-400" />
        <MoonIcon v-else class="w-5 h-5 text-slate-400" />
      </button>

      <!-- Tombol Logout -->
      <button
        @click="logout"
        class="px-3 py-1.5 text-xs font-bold text-red-600 dark:text-red-400 bg-white dark:bg-slate-800 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors border border-slate-200 dark:border-slate-700"
      >
        Logout
      </button>
    </div>
  </header>

  <!-- Navigasi Bawah Khusus Mobile (Muncuk jika layar < md) -->
  <nav class="md:hidden fixed bottom-0 w-full flex items-center justify-around bg-white/95 dark:bg-slate-900/95 backdrop-blur-md border-t border-slate-200 dark:border-slate-800 pb-safe pt-1 shadow-[0_-4px_10px_rgba(0,0,0,0.05)] z-50">
    <RouterLink
      v-for="item in navItems"
      :key="item.id"
      :to="item.to"
      class="flex flex-col items-center gap-1 py-2 px-1 relative w-full text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 transition-colors"
      active-class="!text-indigo-600 dark:!text-indigo-400"
    >
      <component :is="item.iconComp" class="w-6 h-6 shrink-0" />
      <span class="text-[10px] font-semibold">{{ item.label }}</span>
      <span v-if="item.id === 'register' && cartCount > 0" class="absolute top-1 right-2 sm:right-6 w-4 h-4 flex items-center justify-center bg-red-500 text-white text-[9px] font-bold rounded-full border border-white dark:border-slate-900 shadow-sm animate-pulse">
        {{ cartCount }}
      </span>
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
    return 'text-amber-700 dark:text-amber-400';
  }
  if (role === 'kepala_kasir') {
    return 'text-purple-700 dark:text-purple-400';
  }
  return 'text-blue-700 dark:text-blue-400';
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
