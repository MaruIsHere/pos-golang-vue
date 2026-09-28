<template>
  <div class="flex flex-col gap-5">
    <!-- Top Page Header with Export Controls -->
    <AppCard body-class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-5">
      <div>
        <h2 class="text-xl font-bold text-slate-800 dark:text-slate-100">Laporan Penjualan & Analytics</h2>
        <p class="text-sm text-slate-500">Ringkasan performa bisnis, grafik visual, barang terlaris, barang kurang laku, dan ekspor dokumen</p>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <AppButton variant="secondary"  class="secondary flex items-center gap-1.5" @click="fetchStats">
          <ArrowPathIcon class="w-4 h-4" />
          <span>Refresh</span>
        </AppButton>

        <AppButton variant="success"  class="excel flex items-center gap-1.5" @click="handleExportExcel" :disabled="isLoading">
          <ArrowDownTrayIcon class="w-4 h-4" />
          <span>Export Excel</span>
        </AppButton>

        <AppButton variant="danger"  class="pdf flex items-center gap-1.5" @click="handleExportPDF" :disabled="isLoading">
          <ArrowDownTrayIcon class="w-4 h-4" />
          <span>Export PDF</span>
        </AppButton>
      </div>
    </AppCard>

    <!-- Stat Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
      <AppCard body-class="flex items-center gap-4">
        <div class="p-3 rounded-xl bg-emerald-100 dark:bg-emerald-900/30">
          <BanknotesIcon class="w-6 h-6 text-emerald-600" />
        </div>
        <div>
          <span class="text-sm font-semibold text-slate-500 dark:text-slate-400">Total Omset Penjualan</span>
          <h3 class="text-xl font-bold text-slate-800 dark:text-slate-100">Rp {{ formatPrice(stats.total_revenue) }}</h3>
        </div>
      </AppCard>

      <AppCard body-class="flex items-center gap-4">
        <div class="p-3 rounded-xl bg-indigo-100 dark:bg-indigo-900/30">
          <DocumentTextIcon class="w-6 h-6 text-indigo-600" />
        </div>
        <div>
          <span class="text-sm font-semibold text-slate-500 dark:text-slate-400">Total Transaksi Selesai</span>
          <h3 class="text-xl font-bold text-slate-800 dark:text-slate-100">{{ stats.total_orders }} Transaksi</h3>
        </div>
      </AppCard>

      <AppCard body-class="flex items-center gap-4">
        <div class="p-3 rounded-xl bg-amber-100 dark:bg-amber-900/30">
          <ShoppingBagIcon class="w-6 h-6 text-amber-600" />
        </div>
        <div>
          <span class="text-sm font-semibold text-slate-500 dark:text-slate-400">Total Item Terjual</span>
          <h3 class="text-xl font-bold text-slate-800 dark:text-slate-100">{{ stats.total_items_sold }} Pcs</h3>
        </div>
      </AppCard>
    </div>

    <!-- Section 1: Visual Charts (Grafik Penjualan) -->
    <div class="section-title">
      <h3>Grafik Penjualan & Distribusi Produk</h3>
      <p class="text-xs text-secondary">Visualisasi visual tren produk terlaris, kurang laku, dan tipe barang</p>
    </div>

    <div class="charts-grid">
      <!-- Chart 1: Top vs Slow Moving Products -->
      <AppCard>
        <div class="chart-header">
          <h4>5 Barang Paling Laku vs 5 Barang Kurang Laku</h4>
        </div>
        <div class="chart-body">
          <Bar v-if="topVsSlowChartData.labels.length > 0" :data="topVsSlowChartData" :options="chartOptions" />
          <div v-else class="empty-chart">Belum ada data grafik penjualan</div>
        </div>
      </AppCard>

      <!-- Chart 2: Sales Distribution per Product Type -->
      <AppCard>
        <div class="chart-header">
          <h4>Distribusi Omset Per Tipe Produk</h4>
        </div>
        <div class="chart-body chart-body-doughnut">
          <Doughnut v-if="typeChartData.labels.length > 0" :data="typeChartData" :options="doughnutOptions" />
          <div v-else class="empty-chart">Belum ada data tipe produk</div>
        </div>
      </AppCard>
    </div>

    <!-- Section 2: Barang Paling Laku & Barang Kurang Laku Tables -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- Barang Paling Laku (Top Sellers) -->
      <AppCard>
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
      </AppCard>

      <!-- Barang Kurang Laku (Slow Moving / Evaluasi Stok) -->
      <AppCard>
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
      </AppCard>
    </div>

    <!-- Section 3: List Seluruh Barang Laku (Full Detailed Table) -->
    <AppCard class="full-width-card">
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
    </AppCard>

    <!-- Section 4: Breakdown Sales Per Artist & Tipe Produk -->
    <div class="sub-analytics-grid">
      <!-- Sales by Artist -->
      <AppCard>
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
      </AppCard>

      <!-- Sales by Product Type -->
      <AppCard>
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
      </AppCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppCard from '@/components/ui/AppCard.vue';
import AppInput from '@/components/ui/AppInput.vue';
import { ref, computed, onMounted } from 'vue';
import api from '@/utils/api';
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
    const res = await api.get('/reports/dashboard');
    stats.value = res.data;
  } catch (err: any) {
    console.error('Fetch stats error:', err.response?.data?.error || err.message || 'Error occurred');
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
    ...top5.map(p => p.product_name),
    ...slow5.map(p => p.product_name)
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

const handleExportExcel = async () => {
  await exportToExcel(stats.value, 'Keseluruhan Penjualan');
};

const handleExportPDF = () => {
  exportToPDF(stats.value, 'Keseluruhan Penjualan');
};
</script>


