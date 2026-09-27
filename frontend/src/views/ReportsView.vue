<template>
  <div class="reports-page">
    <!-- Top Page Header with Export Controls -->
    <div class="page-header glass-panel">
      <div class="header-title">
        <h2>Laporan Penjualan & Analytics</h2>
        <p>Ringkasan performa bisnis, grafik visual, barang terlaris, barang kurang laku, dan ekspor dokumen</p>
      </div>

      <div class="header-actions">
        <button class="btn btn-secondary flex items-center gap-1.5" @click="fetchStats">
          <ArrowPathIcon class="w-4 h-4" />
          <span>Refresh</span>
        </button>

        <button class="btn btn-excel flex items-center gap-1.5" @click="handleExportExcel" :disabled="isLoading">
          <ArrowDownTrayIcon class="w-4 h-4" />
          <span>Export Excel</span>
        </button>

        <button class="btn btn-pdf flex items-center gap-1.5" @click="handleExportPDF" :disabled="isLoading">
          <ArrowDownTrayIcon class="w-4 h-4" />
          <span>Export PDF</span>
        </button>
      </div>
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
          <span class="stat-label">Total Transaksi Selesai</span>
          <h3 class="stat-value">{{ stats.total_orders }} Transaksi</h3>
        </div>
      </div>

      <div class="stat-card glass-panel">
        <div class="stat-icon icon-items">
          <ShoppingBagIcon class="w-6 h-6 text-amber-600" />
        </div>
        <div class="stat-info">
          <span class="stat-label">Total Item Terjual</span>
          <h3 class="stat-value">{{ stats.total_items_sold }} Pcs</h3>
        </div>
      </div>
    </div>

    <!-- Section 1: Visual Charts (Grafik Penjualan) -->
    <div class="section-title">
      <h3>📈 Grafik Penjualan & Distribusi Produk</h3>
      <p class="text-xs text-secondary">Visualisasi visual tren produk terlaris, kurang laku, dan tipe barang</p>
    </div>

    <div class="charts-grid">
      <!-- Chart 1: Top vs Slow Moving Products -->
      <div class="chart-card glass-panel">
        <div class="chart-header">
          <h4>🔥 5 Barang Paling Laku vs ⚠️ 5 Barang Kurang Laku</h4>
        </div>
        <div class="chart-body">
          <Bar v-if="topVsSlowChartData.labels.length > 0" :data="topVsSlowChartData" :options="chartOptions" />
          <div v-else class="empty-chart">Belum ada data grafik penjualan</div>
        </div>
      </div>

      <!-- Chart 2: Sales Distribution per Product Type -->
      <div class="chart-card glass-panel">
        <div class="chart-header">
          <h4>🏷️ Distribusi Omset Per Tipe Produk</h4>
        </div>
        <div class="chart-body chart-body-doughnut">
          <Doughnut v-if="typeChartData.labels.length > 0" :data="typeChartData" :options="doughnutOptions" />
          <div v-else class="empty-chart">Belum ada data tipe produk</div>
        </div>
      </div>
    </div>

    <!-- Section 2: Barang Paling Laku & Barang Kurang Laku Tables -->
    <div class="dashboard-columns">
      <!-- Barang Paling Laku (Top Sellers) -->
      <div class="dash-card glass-panel">
        <div class="dash-card-header">
          <h3 class="flex items-center gap-1.5 text-emerald-600">
            <FireIcon class="w-5 h-5 text-amber-500" /> Barang Paling Laku (Top Sellers)
          </h3>
          <span class="text-xs text-secondary">Produk dengan tingkat penjualan tertinggi</span>
        </div>
        <div class="dash-card-body">
          <div v-if="!stats.top_products || stats.top_products.length === 0" class="empty-text">
            Belum ada data penjualan produk
          </div>
          <div v-else class="top-list">
            <div v-for="(p, index) in stats.top_products" :key="index" class="top-item">
              <div class="rank-badge" :class="'rank-' + (index + 1)">#{{ index + 1 }}</div>
              <div class="item-name-box">
                <span class="top-name font-bold">{{ p.product_name }}</span>
                <span class="top-qty text-xs text-secondary">{{ p.artist }} • {{ p.product_type }}</span>
              </div>
              <div class="text-right">
                <div class="top-sales font-extrabold text-emerald-600">Rp {{ formatPrice(p.total_sales) }}</div>
                <div class="text-xs font-semibold text-slate-500">{{ p.total_qty }} pcs terjual</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Barang Kurang Laku (Slow Moving / Evaluasi Stok) -->
      <div class="dash-card glass-panel">
        <div class="dash-card-header">
          <h3 class="flex items-center gap-1.5 text-rose-600">
            <ExclamationTriangleIcon class="w-5 h-5 text-rose-500" /> Barang Kurang Laku (Slow Moving)
          </h3>
          <span class="text-xs text-secondary">Produk dengan penjualan paling rendah / perlu promo</span>
        </div>
        <div class="dash-card-body">
          <div v-if="!stats.least_products || stats.least_products.length === 0" class="empty-text">
            Belum ada data evaluasi stok produk
          </div>
          <div v-else class="top-list">
            <div v-for="(p, index) in stats.least_products" :key="index" class="least-item">
              <div class="least-badge">#{{ index + 1 }}</div>
              <div class="item-name-box">
                <span class="top-name font-bold text-slate-800">{{ p.product_name }}</span>
                <span class="top-qty text-xs text-rose-500 font-medium">Sisa Stok: {{ p.stock ?? 0 }} unit</span>
              </div>
              <div class="text-right">
                <div class="font-bold text-slate-700">{{ p.total_qty }} pcs terjual</div>
                <div class="text-xs text-slate-500">Rp {{ formatPrice(p.total_sales) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Section 3: List Seluruh Barang Laku (Full Detailed Table) -->
    <div class="dash-card glass-panel full-width-card">
      <div class="dash-card-header flex justify-between items-center">
        <div>
          <h3 class="flex items-center gap-1.5">
            <ListBulletIcon class="w-5 h-5 text-indigo-600" /> List Seluruh Barang Laku
          </h3>
          <span class="text-xs text-secondary">Rincian lengkap kinerjaper barang yang telah terjual</span>
        </div>
        <span class="badge badge-info" v-if="stats.all_sold_products">
          Total {{ stats.all_sold_products.length }} Jenis Produk
        </span>
      </div>

      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>No</th>
              <th>Nama Produk</th>
              <th>Artist</th>
              <th>Tipe Produk</th>
              <th>Harga Satuan</th>
              <th>Total Kuantitas</th>
              <th class="text-right">Total Sales</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!stats.all_sold_products || stats.all_sold_products.length === 0">
              <td colspan="7" class="text-center">Belum ada transaksi barang terjual</td>
            </tr>
            <tr v-else v-for="(p, idx) in stats.all_sold_products" :key="idx">
              <td>{{ idx + 1 }}</td>
              <td class="font-bold">{{ p.product_name }}</td>
              <td><span class="sub-tag artist-tag">{{ p.artist }}</span></td>
              <td><span class="sub-tag type-tag">{{ p.product_type }}</span></td>
              <td>Rp {{ formatPrice(p.price || 0) }}</td>
              <td><span class="qty-pill">{{ p.total_qty }} Pcs</span></td>
              <td class="text-right font-bold text-emerald-600">Rp {{ formatPrice(p.total_sales) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Section 4: Breakdown Sales Per Artist & Tipe Produk -->
    <div class="sub-analytics-grid">
      <!-- Sales by Artist -->
      <div class="dash-card glass-panel">
        <div class="dash-card-header">
          <h3 class="flex items-center gap-1.5 text-amber-600">
            <UserIcon class="w-5 h-5" /> Penjualan Per Artist
          </h3>
        </div>
        <div class="dash-card-body">
          <div v-if="!stats.sales_by_artist || stats.sales_by_artist.length === 0" class="empty-text">
            Belum ada data artist
          </div>
          <div v-else class="sub-stat-list">
            <div v-for="(item, idx) in stats.sales_by_artist" :key="idx" class="sub-stat-item">
              <div class="item-name-box">
                <span class="sub-stat-name font-bold">{{ item.name }}</span>
                <span class="sub-stat-qty">{{ item.total_qty }} pcs terjual</span>
              </div>
              <div class="sub-stat-sales">Rp {{ formatPrice(item.total_sales) }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Sales by Product Type -->
      <div class="dash-card glass-panel">
        <div class="dash-card-header">
          <h3 class="flex items-center gap-1.5 text-pink-600">
            <TagIcon class="w-5 h-5" /> Penjualan Per Tipe Produk
          </h3>
        </div>
        <div class="dash-card-body">
          <div v-if="!stats.sales_by_type || stats.sales_by_type.length === 0" class="empty-text">
            Belum ada data tipe produk
          </div>
          <div v-else class="sub-stat-list">
            <div v-for="(item, idx) in stats.sales_by_type" :key="idx" class="sub-stat-item">
              <div class="item-name-box">
                <span class="sub-stat-name font-bold">{{ item.name }}</span>
                <span class="sub-stat-qty">{{ item.total_qty }} pcs terjual</span>
              </div>
              <div class="sub-stat-sales">Rp {{ formatPrice(item.total_sales) }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { DashboardStats } from '../types';
import { exportToExcel, exportToPDF } from '../utils/exportReport';
import { 
  Chart as ChartJS, 
  Title, 
  Tooltip, 
  Legend, 
  BarElement, 
  CategoryScale, 
  LinearScale, 
  ArcElement 
} from 'chart.js';
import { Bar, Doughnut } from 'vue-chartjs';
import { 
  BanknotesIcon, 
  DocumentTextIcon, 
  ShoppingBagIcon, 
  FireIcon, 
  ArrowPathIcon,
  UserIcon,
  TagIcon,
  ArrowDownTrayIcon,
  ExclamationTriangleIcon,
  ListBulletIcon
} from '@heroicons/vue/24/outline';

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, ArcElement);

const isLoading = ref(true);
const stats = ref<DashboardStats>({
  total_revenue: 0,
  total_orders: 0,
  total_items_sold: 0,
  top_products: [],
  least_products: [],
  all_sold_products: [],
  sales_by_artist: [],
  sales_by_type: [],
  recent_orders: []
});

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

const fetchStats = async () => {
  isLoading.value = true;
  try {
    const res = await fetch('/api/reports/dashboard');
    if (res.ok) stats.value = await res.json();
  } catch (err) {
    console.error('Fetch stats error:', err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchStats);

// Chart 1 Data: Top 5 vs Slowest 5
const topVsSlowChartData = computed(() => {
  const top5 = (stats.value.top_products || []).slice(0, 5);
  const slow5 = (stats.value.least_products || []).slice(0, 5);

  const labels = [
    ...top5.map(p => `🔥 ${p.product_name}`),
    ...slow5.map(p => `⚠️ ${p.product_name}`)
  ];

  const dataValues = [
    ...top5.map(p => p.total_qty),
    ...slow5.map(p => p.total_qty)
  ];

  const backgroundColors = [
    ...top5.map(() => 'rgba(16, 185, 129, 0.75)'), // Green for top sellers
    ...slow5.map(() => 'rgba(244, 63, 94, 0.75)')   // Rose for slow movers
  ];

  return {
    labels,
    datasets: [{
      label: 'Kuantitas Terjual (Pcs)',
      backgroundColor: backgroundColors,
      borderRadius: 6,
      data: dataValues
    }]
  };
});

// Chart 2 Data: Doughnut Chart for Product Types
const typeChartData = computed(() => {
  const types = stats.value.sales_by_type || [];
  return {
    labels: types.map(t => t.name),
    datasets: [{
      backgroundColor: ['#4f46e5', '#059669', '#d97706', '#db2777', '#0284c7', '#7c3aed'],
      data: types.map(t => t.total_sales)
    }]
  };
});

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false }
  },
  scales: {
    y: { beginAtZero: true }
  }
};

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'bottom' as const }
  }
};

const handleExportExcel = () => {
  exportToExcel(stats.value, 'Keseluruhan Penjualan');
};

const handleExportPDF = () => {
  exportToPDF(stats.value, 'Keseluruhan Penjualan');
};
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

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-excel {
  background-color: #10b981;
  color: #ffffff;
  font-weight: 700;
  border: none;
}
.btn-excel:hover {
  background-color: #059669;
}

.btn-pdf {
  background-color: #e11d48;
  color: #ffffff;
  font-weight: 700;
  border: none;
}
.btn-pdf:hover {
  background-color: #be123c;
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

.icon-revenue { background: rgba(16, 185, 129, 0.15); border: 1px solid rgba(16, 185, 129, 0.3); }
.icon-orders { background: rgba(99, 102, 241, 0.15); border: 1px solid rgba(99, 102, 241, 0.3); }
.icon-items { background: rgba(245, 158, 11, 0.15); border: 1px solid rgba(245, 158, 11, 0.3); }

.stat-label { font-size: 0.8rem; color: var(--text-secondary); font-weight: 600; }
.stat-value { font-size: 1.35rem; font-weight: 800; color: var(--text-primary); margin-top: 0.15rem; }

.section-title { margin-top: 0.5rem; }
.section-title h3 { font-size: 1.05rem; font-weight: 700; color: var(--text-primary); }

.charts-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.25rem;
}

@media (min-width: 1024px) {
  .charts-grid {
    grid-template-columns: 1.4fr 1fr;
  }
}

.chart-card {
  padding: 1.25rem;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
}

.chart-header h4 {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 0.75rem;
}

.chart-body {
  height: 260px;
  position: relative;
}

.chart-body-doughnut {
  display: flex;
  justify-content: center;
  align-items: center;
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
}

.top-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.top-item, .least-item {
  display: flex;
  align-items: center;
  gap: 0.85rem;
  padding: 0.75rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
}

.least-item {
  border-left: 3px solid #f43f5e;
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
}

.least-badge {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.85rem;
  background: rgba(244, 63, 94, 0.12);
  color: #e11d48;
}

.rank-1 { background: rgba(245, 158, 11, 0.15); color: #fbbf24; border: 1px solid rgba(245, 158, 11, 0.3); }
.rank-2 { background: var(--bg-primary); color: var(--text-secondary); border: 1px solid var(--border-color); }
.rank-3 { background: rgba(249, 115, 22, 0.15); color: #fb923c; border: 1px solid rgba(249, 115, 22, 0.3); }

.item-name-box {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.full-width-card {
  width: 100%;
}

.table-container {
  overflow-x: auto;
  margin-top: 0.5rem;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.875rem;
}

.data-table th {
  padding: 0.85rem;
  font-weight: 600;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-primary);
}

.data-table td {
  padding: 0.85rem;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
}

.qty-pill {
  background: var(--bg-primary);
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
  border: 1px solid var(--border-color);
}

.sub-analytics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.25rem;
}

.sub-stat-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.sub-stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.65rem 0.85rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
}

.sub-stat-sales { font-weight: 800; color: #059669; font-size: 0.875rem; }

.sub-tag {
  display: inline-block;
  padding: 0.18rem 0.5rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
}

.artist-tag { background: rgba(245, 158, 11, 0.12); color: #d97706; border: 1px solid rgba(245, 158, 11, 0.25); }
.type-tag { background: rgba(236, 72, 153, 0.12); color: #db2777; border: 1px solid rgba(236, 72, 153, 0.25); }

.empty-text, .empty-chart {
  text-align: center;
  padding: 2rem;
  color: var(--text-muted);
  font-size: 0.85rem;
}
</style>
