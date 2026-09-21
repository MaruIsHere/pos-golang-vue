<template>
  <header class="glass-nav header-container">
    <div class="header-left">
      <div class="brand-logo">
        <div class="logo-icon">🛒</div>
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
      <button 
        v-for="item in navItems" 
        :key="item.id" 
        class="nav-btn" 
        :class="{ active: currentTab === item.id }"
        @click="$emit('change-tab', item.id)"
      >
        <span class="nav-icon">{{ item.icon }}</span>
        <span class="nav-text">{{ item.label }}</span>
      </button>
    </nav>

    <div class="header-right">
      <div class="clock-display">
        <span class="time-text">{{ currentTime }}</span>
      </div>
    </div>
  </header>

  <!-- Mobile Bottom Navigation Bar -->
  <nav class="mobile-bottom-nav">
    <button 
      v-for="item in navItems" 
      :key="item.id" 
      class="mobile-nav-btn" 
      :class="{ active: currentTab === item.id }"
      @click="$emit('change-tab', item.id)"
    >
      <span class="nav-icon">{{ item.icon }}</span>
      <span class="nav-text">{{ item.label }}</span>
      <span v-if="item.id === 'register' && cartCount > 0" class="cart-badge">{{ cartCount }}</span>
    </button>
  </nav>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

const props = defineProps({
  currentTab: { type: String, default: 'register' },
  storeSetting: { type: Object, default: () => ({}) },
  cartCount: { type: Number, default: 0 }
});

defineEmits(['change-tab']);

const navItems = [
  { id: 'register', label: 'Kasir', icon: '🏪' },
  { id: 'products', label: 'Produk', icon: '📦' },
  { id: 'inventory', label: 'Inventoris', icon: '📥' },
  { id: 'customers', label: 'Pelanggan', icon: '👥' },
  { id: 'orders', label: 'Riwayat', icon: '📜' },
  { id: 'reports', label: 'Laporan', icon: '📊' },
  { id: 'settings', label: 'Pengaturan', icon: '⚙️' }
];

const currentTime = ref('');

const updateTime = () => {
  const now = new Date();
  currentTime.value = now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
};

let timer = null;
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
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-purple));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  box-shadow: var(--shadow-glow);
}

.store-name {
  font-size: 1.1rem;
  font-weight: 800;
  letter-spacing: -0.02em;
  color: #ffffff;
  line-height: 1.2;
}

.db-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
  font-size: 0.65rem;
  font-weight: 800;
  letter-spacing: 0.05em;
}

.db-sqlite {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.4);
}

.db-mysql {
  background: rgba(245, 158, 11, 0.2);
  color: #fbbf24;
  border: 1px solid rgba(245, 158, 11, 0.4);
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
  gap: 0.35rem;
  background: rgba(15, 23, 42, 0.6);
  padding: 0.3rem;
  border-radius: var(--radius-lg);
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
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border-radius: var(--radius-md);
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-family: var(--font-family);
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.nav-btn:hover {
  color: var(--text-primary);
  background: rgba(255, 255, 255, 0.05);
}

.nav-btn.active {
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-purple));
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.clock-display {
  padding: 0.4rem 0.8rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  font-weight: 700;
  font-size: 0.85rem;
  color: var(--text-secondary);
}

/* Mobile Bottom Nav */
.mobile-bottom-nav {
  display: flex;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border-top: 1px solid var(--border-color);
  z-index: 100;
  align-items: center;
  justify-content: space-around;
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
  font-size: 0.7rem;
  font-weight: 600;
  cursor: pointer;
  flex: 1;
  height: 100%;
  position: relative;
}

.mobile-nav-btn.active {
  color: #6366f1;
}

.mobile-nav-btn .nav-icon {
  font-size: 1.2rem;
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
