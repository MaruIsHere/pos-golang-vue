<template>
  <div class="reports-page">
    <div class="page-header glass-panel">
      <div class="header-title">
        <h2>Laporan Penjualan & Dashboard</h2>
        <p>Ringkasan performa penjualan dan statistik produk terlaris</p>
      </div>

      <button class="btn btn-secondary flex items-center gap-1" @click="fetchStats">
        <ArrowPathIcon class="w-4 h-4" />
        <span>Refresh</span>
      </button>
    </div>

    <!-- Stat Cards Grid -->
    <div class="stats-grid">
      <div class="stat-card glass-panel">
        <div class="stat-icon icon-revenue">
          <BanknotesIcon class="w-6 h-6 text-emerald-600" />
        </div>
        <div class="stat-info">
          <span class="stat-label">Total Omset Penjualan</span>
          <h3 class="stat-value">Rp {{ formatPrice(stats.total_revenue) }}</h3>
        </div>
      </div>

      <div class="stat-card glass-panel">
        <div class="stat-icon icon-orders">
          <DocumentTextIcon class="w-6 h-6 text-indigo-600" />
        </div>
        <div class="stat-info">
          <span class="stat-label">Total Transaksi</span>
          <h3 class="stat-value">{{ stats.total_orders }} Transaksi</h3>
        </div>
      </div>

      <div class="stat-card glass-panel">
        <div class="stat-icon icon-items">
          <ShoppingBagIcon class="w-6 h-6 text-amber-600" />
        </div>
        <div class="stat-info">
          <span class="stat-label">Item Terjual</span>
          <h3 class="stat-value">{{ stats.total_items_sold }} Pcs</h3>
        </div>
      </div>
    </div>

    <!-- Dashboard Content Columns -->
    <div class="dashboard-columns">
      <!-- Top Products Card -->
      <div class="dash-card glass-panel">
        <div class="dash-card-header">
          <h3 class="flex items-center gap-1.5">
            <FireIcon class="w-5 h-5 text-amber-500" /> 5 Produk Terlaris
          </h3>
        </div>
        <div class="dash-card-body">
          <div v-if="!stats.top_products || stats.top_products.length === 0" class="empty-text">
            Belum ada data penjualan
          </div>
          <div v-else class="top-list">
            <div v-for="(p, index) in stats.top_products" :key="index" class="top-item">
              <div class="rank-badge" :class="'rank-' + (index + 1)">#{{ index + 1 }}</div>
              <div class="item-name-box">
                <span class="top-name">{{ p.product_name }}</span>
                <span class="top-qty">{{ p.total_qty }} pcs terjual</span>
              </div>
              <div class="top-sales">Rp {{ formatPrice(p.total_sales) }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Orders Card -->
      <div class="dash-card glass-panel">
        <div class="dash-card-header">
          <h3 class="flex items-center gap-1.5">
            <ClockIcon class="w-5 h-5 text-indigo-500" /> Transaksi Terakhir
          </h3>
        </div>
        <div class="dash-card-body">
          <div v-if="!stats.recent_orders || stats.recent_orders.length === 0" class="empty-text">
            Belum ada transaksi
          </div>
          <div v-else class="recent-list">
            <div v-for="ro in stats.recent_orders" :key="ro.id" class="recent-item">
              <div class="recent-left">
                <code class="inv-code">{{ ro.invoice_no }}</code>
                <span class="recent-cust">{{ ro.customer_name || 'Umum' }} • {{ (ro.payment_method || 'cash').toUpperCase() }}</span>
              </div>
              <div class="recent-right">
                <span class="recent-total">Rp {{ formatPrice(ro.grand_total) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { DashboardStats } from '../types';
import { BanknotesIcon, DocumentTextIcon, ShoppingBagIcon, FireIcon, ClockIcon, ArrowPathIcon } from '@heroicons/vue/24/outline';

const stats = ref<DashboardStats>({
  total_revenue: 0,
  total_orders: 0,
  total_items_sold: 0,
  top_products: [],
  recent_orders: []
});

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

const fetchStats = async () => {
  try {
    const res = await fetch('/api/reports/dashboard');
    if (res.ok) stats.value = await res.json();
  } catch (err) {
    console.error('Fetch stats error:', err);
  }
};

onMounted(fetchStats);
</script>

<style scoped>
.reports-page {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.header-title h2 {
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--text-primary);
}

.header-title p {
  font-size: 0.8rem;
  color: var(--text-secondary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 1.25rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  padding: 1.25rem 1.5rem;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-revenue {
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.icon-orders {
  background: rgba(99, 102, 241, 0.15);
  border: 1px solid rgba(99, 102, 241, 0.3);
}

.icon-items {
  background: rgba(245, 158, 11, 0.15);
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.stat-label {
  font-size: 0.8rem;
  color: var(--text-secondary);
  font-weight: 600;
}

.stat-value {
  font-size: 1.35rem;
  font-weight: 800;
  color: var(--text-primary);
  margin-top: 0.15rem;
}

.dashboard-columns {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.25rem;
}

@media (min-width: 1024px) {
  .dashboard-columns {
    grid-template-columns: 1fr 1fr;
  }
}

.dash-card {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.dash-card-header {
  padding-bottom: 0.85rem;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 1rem;
}

.dash-card-header h3 {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-primary);
}

.top-list, .recent-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.top-item {
  display: flex;
  align-items: center;
  gap: 0.85rem;
  padding: 0.75rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
}

.rank-badge {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.85rem;
  background: var(--bg-primary);
  color: var(--text-secondary);
}

.rank-1 { background: rgba(245, 158, 11, 0.15); color: #fbbf24; border: 1px solid rgba(245, 158, 11, 0.3); }
.rank-2 { background: var(--bg-primary); color: var(--text-secondary); border: 1px solid var(--border-color); }
.rank-3 { background: rgba(249, 115, 22, 0.15); color: #fb923c; border: 1px solid rgba(249, 115, 22, 0.3); }

.item-name-box {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.top-name {
  font-weight: 700;
  font-size: 0.875rem;
  color: var(--text-primary);
}

.top-qty {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.top-sales {
  font-weight: 800;
  color: #059669;
}

.recent-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
}

.recent-left {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.inv-code {
  font-weight: 700;
  color: var(--accent-primary);
  font-size: 0.8rem;
}

.recent-cust {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.recent-total {
  font-weight: 800;
  color: #059669;
}

.empty-text {
  text-align: center;
  padding: 2rem;
  color: var(--text-muted);
}
</style>
