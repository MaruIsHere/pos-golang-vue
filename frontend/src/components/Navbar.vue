<template>
  <!-- 1. NAVIGATION RAIL (KIRI) - Khusus Desktop/Tablet -->
  <nav class="hidden md:flex fixed left-0 top-0 h-full w-20 flex-col items-center bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 z-50 py-4 shadow-[4px_0_24px_rgba(0,0,0,0.02)]">
    <!-- Logo POS -->
    <div class="w-12 h-12 mb-6 rounded-2xl bg-indigo-50 dark:bg-indigo-900/30 flex items-center justify-center border border-indigo-100 dark:border-indigo-800/50 shadow-sm">
      <ShoppingCartIcon class="w-6 h-6 text-indigo-600 dark:text-indigo-400" />
    </div>

    <!-- Menu Utama -->
    <div class="flex flex-col gap-3 w-full px-2 flex-1 overflow-y-auto no-scrollbar items-center">
      <RouterLink
        v-for="item in navItems"
        :key="item.id"
        :to="item.to"
        class="group flex flex-col items-center justify-center w-14 h-14 rounded-2xl text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-slate-100 hover:bg-slate-100 dark:hover:bg-slate-800 transition-all duration-200 relative"
        active-class="bg-indigo-100 dark:bg-indigo-900/40 !text-indigo-700 dark:!text-indigo-400 shadow-inner"
        :title="item.label"
      >
        <component :is="item.iconComp" class="w-6 h-6 shrink-0 transition-transform group-hover:scale-110" />
        <span class="text-[9px] font-bold mt-1 tracking-tight">{{ item.label }}</span>
        <!-- Badge Keranjang (Khusus menu Kasir) -->
        <span v-if="item.id === 'register' && cartCount > 0" class="absolute top-0 right-0 w-4 h-4 flex items-center justify-center bg-red-500 text-white text-[9px] font-bold rounded-full border-2 border-white dark:border-slate-900 shadow-sm animate-pulse">
          {{ cartCount }}
        </span>
      </RouterLink>
    </div>

    <!-- Tombol Tema di Bawah Rel -->
    <button 
      class="mt-auto w-12 h-12 rounded-full flex items-center justify-center text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors" 
      :title="isDarkMode ? 'Mode Terang' : 'Mode Gelap'"
      @click="toggleTheme"
    >
      <SunIcon v-if="isDarkMode" class="w-5 h-5 text-amber-400" />
      <MoonIcon v-else class="w-5 h-5 text-slate-400" />
    </button>
  </nav>

  <!-- 2. TOP APP BAR (ATAS) - Khusus Desktop/Tablet -->
  <header class="hidden md:flex fixed top-0 left-20 right-0 h-16 items-center justify-between px-6 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 z-40">
    <!-- Info Toko & Database -->
    <div class="flex items-center gap-3">
      <!-- Active Outlet Selector -->
      <div v-if="storeContextStore.stores.length > 0" class="relative">
        <select 
          class="appearance-none pl-8 pr-7 py-1.5 bg-indigo-50 dark:bg-indigo-950/60 border border-indigo-200 dark:border-indigo-800 rounded-xl text-xs font-bold text-indigo-700 dark:text-indigo-300 outline-none cursor-pointer hover:bg-indigo-100 transition-colors"
          :value="storeContextStore.activeStoreId"
          @change="onStoreChange"
        >
          <option v-for="s in storeContextStore.stores" :key="s.id" :value="s.id">
            🏪 {{ s.name }} ({{ s.code }})
          </option>
        </select>
        <BuildingStorefrontIcon class="w-4 h-4 text-indigo-500 absolute left-2.5 top-2 pointer-events-none" />
      </div>

      <h1 v-else class="text-base font-extrabold text-slate-800 dark:text-slate-100 tracking-tight">
        {{ storeSetting.store_name || 'POS KASIR PRO' }}
      </h1>

      <div class="h-4 w-px bg-slate-300 dark:bg-slate-700"></div>
      <span 
        class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider"
        :class="storeSetting.db_engine === 'mysql' ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-400 border border-amber-200 dark:border-amber-800/50' : 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800/50'"
      >
        <span class="w-1.5 h-1.5 rounded-full" :class="storeSetting.db_engine === 'mysql' ? 'bg-amber-500 animate-pulse' : 'bg-emerald-500'"></span> 
        {{ (storeSetting.db_engine || 'sqlite').toUpperCase() }}
      </span>
    </div>

    <!-- Profil & Jam & Logout -->
    <div class="flex items-center gap-4">
      <div class="flex items-center gap-2">
        <ClockIcon class="w-4 h-4 text-slate-400" />
        <span class="text-xs font-bold text-slate-600 dark:text-slate-300 font-mono tracking-wider">{{ currentTime }}</span>
      </div>
      
      <div class="h-4 w-px bg-slate-300 dark:bg-slate-700"></div>

      <div v-if="authStore.user" class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center border border-slate-200 dark:border-slate-700">
          <img v-if="authStore.user.profile_photo" :src="authStore.user.profile_photo" alt="Foto profil" class="w-full h-full rounded-full object-cover" />
          <UserIcon v-else class="w-4 h-4 text-indigo-500" />
        </div>
        <div class="flex flex-col leading-none">
          <span class="text-sm font-bold text-slate-900 dark:text-slate-100">{{ authStore.user.name || authStore.user.username }}</span>
          <span class="text-[9px] font-extrabold uppercase mt-0.5" :class="roleBadgeClass">
            {{ roleLabel }}
          </span>
        </div>
      </div>

      <button
        @click="logout"
        class="ml-2 px-3 py-1.5 text-xs font-bold text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20 hover:bg-red-100 dark:hover:bg-red-900/40 rounded-lg transition-colors border border-red-100 dark:border-red-900/50"
      >
        Keluar
      </button>
    </div>
  </header>

  <!-- 3. BOTTOM NAVIGATION (BAWAH) - Khusus Mobile/HP -->
  <!-- 3. BOTTOM NAVIGATION (BAWAH) - Khusus Mobile/HP -->
  <nav class="md:hidden fixed bottom-0 left-0 right-0 w-full flex justify-around bg-white/95 dark:bg-slate-900/95 backdrop-blur-md border-t border-slate-200 dark:border-slate-800 pb-safe pt-1 shadow-[0_-4px_10px_rgba(0,0,0,0.05)] z-50">
    <RouterLink
      v-for="item in primaryNavItems"
      :key="item.id"
      :to="item.to"
      class="flex flex-col items-center justify-center gap-1 py-2 px-2 flex-1 relative text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 transition-colors"
      active-class="!text-indigo-600 dark:!text-indigo-400"
    >
      <component :is="item.iconComp" class="w-6 h-6 shrink-0" :class="{'text-indigo-600 dark:text-indigo-400': $route.path === item.to || ($route.path.startsWith(item.to) && item.to !== '/')}" />
      <span class="text-[10px] font-bold whitespace-nowrap" :class="{'text-indigo-600 dark:text-indigo-400': $route.path === item.to || ($route.path.startsWith(item.to) && item.to !== '/')}">{{ item.label }}</span>
      <span v-if="item.id === 'register' && cartCount > 0" class="absolute top-1 right-[calc(50%-1.25rem)] w-4 h-4 flex items-center justify-center bg-red-500 text-white text-[9px] font-bold rounded-full border border-white dark:border-slate-900 shadow-sm animate-pulse">
        {{ cartCount }}
      </span>
    </RouterLink>

    <!-- Tombol Lainnya -->
    <button 
      v-if="secondaryNavItems.length > 0"
      @click="isMoreMenuOpen = true"
      class="flex flex-col items-center justify-center gap-1 py-2 px-2 flex-1 relative text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 transition-colors"
    >
      <EllipsisHorizontalIcon class="w-6 h-6 shrink-0" />
      <span class="text-[10px] font-bold whitespace-nowrap">Lainnya</span>
    </button>
  </nav>

  <!-- Mobile Top Bar Khusus HP (Untuk Jam & DB Status & Theme) -->
  <header class="md:hidden fixed top-0 left-0 right-0 h-12 flex items-center justify-between px-4 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 z-40">
    <div class="flex items-center gap-2">
      <ShoppingCartIcon class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
      <span class="text-xs font-bold text-slate-800 dark:text-slate-100">{{ storeSetting.store_name || 'POS' }}</span>
    </div>
    <div class="flex items-center gap-3">
      <button @click="toggleTheme" class="text-slate-500">
        <SunIcon v-if="isDarkMode" class="w-5 h-5 text-amber-400" />
        <MoonIcon v-else class="w-5 h-5 text-slate-400" />
      </button>
      <button @click="logout" class="text-xs font-bold text-red-600 dark:text-red-400">Logout</button>
    </div>
  </header>

  <!-- Mobile More Menu Drawer -->
  <Teleport to="body">
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="opacity-0 translate-y-full"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 translate-y-full"
  >
    <div v-if="isMoreMenuOpen" class="md:hidden fixed inset-0 z-[200] flex flex-col justify-end">
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm" @click="isMoreMenuOpen = false"></div>
      
      <!-- Drawer Content -->
      <div class="relative bg-white dark:bg-slate-900 rounded-t-3xl shadow-2xl border-t border-slate-200 dark:border-slate-800 pb-safe">
        <div class="flex items-center justify-between p-4 border-b border-slate-100 dark:border-slate-800">
          <h3 class="text-base font-bold text-slate-900 dark:text-slate-100">Menu Lainnya</h3>
          <button @click="isMoreMenuOpen = false" class="p-2 rounded-full bg-slate-100 dark:bg-slate-800 text-slate-500 hover:text-slate-700 dark:hover:text-slate-300">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        
        <div class="grid grid-cols-4 gap-4 p-6">
          <RouterLink
            v-for="item in secondaryNavItems"
            :key="item.id"
            :to="item.to"
            @click="isMoreMenuOpen = false"
            class="flex flex-col items-center justify-center gap-2 text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-slate-100 transition-colors"
            active-class="!text-indigo-600 dark:!text-indigo-400"
          >
            <div class="w-12 h-12 rounded-2xl bg-slate-50 dark:bg-slate-800 flex items-center justify-center border border-slate-100 dark:border-slate-700 shadow-sm" :class="{'!bg-indigo-50 dark:!bg-indigo-900/30 !border-indigo-100 dark:!border-indigo-800/50': $route.path.startsWith(item.to) && item.to !== '/'}">
              <component :is="item.iconComp" class="w-6 h-6 shrink-0" />
            </div>
            <span class="text-[10px] font-bold text-center leading-tight">{{ item.label }}</span>
          </RouterLink>
        </div>
      </div>
    </div>
  </Transition>
  </Teleport>
</template>


<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useSettingsStore } from '../stores/settings';
import { useAuthStore } from '../stores/auth';
import { useStoreContextStore } from '../stores/storeContext';
import { useTheme } from '../composables/useTheme';
import api from '../utils/api';
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
  UserIcon,
  EllipsisHorizontalIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline';

defineProps({
  cartCount: { type: Number, default: 0 }
});

const settingsStore = useSettingsStore();
const { settings: storeSetting } = storeToRefs(settingsStore);
const authStore = useAuthStore();
const storeContextStore = useStoreContextStore();
const { isDarkMode, toggleTheme } = useTheme();

const logout = () => {
  authStore.logout();
};

const onStoreChange = (e: Event) => {
  const target = e.target as HTMLSelectElement;
  if (target && target.value) {
    storeContextStore.setActiveStore(target.value);
    window.location.reload();
  }
};

const syncProfile = async (): Promise<void> => {
  if (!['owner', 'admin'].includes(String(authStore.userRole || '').toLowerCase())) return;
  try {
    const { data } = await api.get('/profile');
    authStore.updateUser(data.user);
  } catch {
    // Keep the locally cached profile available when the API is temporarily unreachable.
  }
};

const allNavItems = [
  { id: 'register', to: '/', label: 'Kasir', iconComp: ShoppingCartIcon, roles: ['kasir', 'kepala_kasir', 'owner', 'admin'] },
  { id: 'products', to: '/products', label: 'Produk', iconComp: CubeIcon, roles: ['kepala_kasir', 'owner', 'admin'] },
  { id: 'stores', to: '/stores', label: 'Toko', iconComp: BuildingStorefrontIcon, roles: ['owner', 'admin'] },
  { id: 'inventory', to: '/inventory', label: 'Inventori', iconComp: ArchiveBoxIcon, roles: ['kepala_kasir', 'owner', 'admin'] },
  { id: 'customers', to: '/customers', label: 'Pelanggan', iconComp: UsersIcon, roles: ['kepala_kasir', 'owner', 'admin'] },
  { id: 'orders', to: '/orders', label: 'Riwayat', iconComp: DocumentTextIcon, roles: ['kasir', 'kepala_kasir', 'owner', 'admin'] },
  { id: 'reports', to: '/reports', label: 'Laporan', iconComp: ChartBarIcon, roles: ['kepala_kasir', 'owner', 'admin'] },
  { id: 'settings', to: '/settings', label: 'Pengaturan', iconComp: Cog6ToothIcon, roles: ['owner', 'admin'] }
];


const isMoreMenuOpen = ref(false);

const primaryNavItems = computed(() => {
  if (navItems.value.length <= 4) return navItems.value;
  return navItems.value.slice(0, 3);
});

const secondaryNavItems = computed(() => {
  if (navItems.value.length <= 4) return [];
  return navItems.value.slice(3);
});

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
  storeContextStore.fetchStores();
  void syncProfile();
  updateTime();
  timer = setInterval(updateTime, 1000);
});

onUnmounted(() => {
  if (timer) clearInterval(timer);
});
</script>
