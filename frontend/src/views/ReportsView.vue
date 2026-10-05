<template>
  <div class="flex flex-col gap-5">
    <!-- Top Page Header with Export Controls -->
    <AppCard body-class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-5">
      <div>
        <h2 class="text-xl font-bold text-slate-800 dark:text-slate-100">Laporan Penjualan & Analytics</h2>
        <p class="text-sm text-slate-500">Ringkasan performa bisnis, grafik visual, barang terlaris, barang kurang laku, dan ekspor dokumen</p>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <AppButton variant="primary" class="flex items-center gap-1.5" @click="fetchStats">
          <ArrowPathIcon class="w-4 h-4" />
          <span>Refresh</span>
        </AppButton>

        <AppButton variant="success"  class="flex items-center gap-1.5" @click="handleExportExcel" :disabled="isLoading">
          <ArrowDownTrayIcon class="w-4 h-4" />
          <span>Export Excel</span>
        </AppButton>

        <AppButton variant="danger"  class="flex items-center gap-1.5" @click="handleExportPDF" :disabled="isLoading">
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
          <h3 class="text-xl font-bold text-slate-800 dark:text-slate-100">{{ formatQuantity(stats.total_items_sold) }}</h3>
        </div>
      </AppCard>
    </div>

    <!-- Section 1: Visual Charts (Grafik Penjualan) -->
    <div class="flex flex-col gap-1 mt-6 mb-4">
      <h3>Grafik Penjualan & Distribusi Produk</h3>
      <p class="text-xs text-secondary">Visualisasi visual tren produk terlaris, kurang laku, dan tipe barang</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- Chart 1: Top vs Slow Moving Products -->
      <AppCard>
        <div class="font-bold text-slate-800 dark:text-slate-100 mb-4 flex items-center gap-2">
          <h4>5 Barang Paling Laku vs 5 Barang Kurang Laku</h4>
        </div>
        <div class="p-5 bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 shadow-sm min-h-[300px]">
          <Bar v-if="topVsSlowChartData.labels.length > 0" :data="topVsSlowChartData" :options="chartOptions" />
          <div v-else class="text-slate-400 text-sm italic flex h-full items-center justify-center">Belum ada data grafik penjualan</div>
        </div>
      </AppCard>

      <!-- Chart 2: Sales Distribution per Product Type -->
      <AppCard>
        <div class="font-bold text-slate-800 dark:text-slate-100 mb-4 flex items-center gap-2">
          <h4>Distribusi Omset Per Tipe Produk</h4>
        </div>
        <div class="p-5 bg-white dark:bg-slate-800 rounded-2xl border border-slate-200 dark:border-slate-700 shadow-sm flex items-center justify-center min-h-[300px]">
          <Doughnut v-if="typeChartData.labels.length > 0" :data="typeChartData" :options="doughnutOptions" />
          <div v-else class="text-slate-400 text-sm italic flex h-full items-center justify-center">Belum ada data tipe produk</div>
        </div>
      </AppCard>
    </div>

    <!-- Section 2: Barang Paling Laku & Barang Kurang Laku Tables -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- Barang Paling Laku (Top Sellers) -->
      <AppCard>
        <div class="flex flex-col gap-1 border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <h3 class="flex items-center gap-1.5 text-emerald-600">
            <FireIcon class="w-5 h-5 text-amber-500" /> Barang Paling Laku (Top Sellers)
          </h3>
          <span class="text-xs text-secondary">Produk dengan tingkat penjualan tertinggi</span>
        </div>
        <div class="flex flex-col">
          <div v-if="!stats.top_products || stats.top_products.length === 0" class="text-sm text-slate-400 italic p-6 text-center">
            Belum ada data penjualan produk
          </div>
          <div v-else class="flex flex-col gap-3">
            <div v-for="(p, index) in stats.top_products" :key="index" class="flex items-center gap-3 p-3 bg-slate-50 dark:bg-slate-900/40 rounded-xl border border-slate-200 dark:border-slate-700/50">
              <div class="w-8 h-8 flex shrink-0 items-center justify-center rounded-lg bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 shadow-sm font-black text-slate-400 text-xs"">#{{ index + 1 }}</div>
              <div class="flex flex-col flex-1 min-w-0">
                <span class="font-bold text-slate-800 dark:text-slate-100 text-sm truncate">{{ p.product_name }}</span>
                <span class="text-xs text-slate-500 truncate mt-0.5">{{ p.artist }} • {{ p.product_type }}</span>
              </div>
              <div class="text-right">
                <div class="font-extrabold text-emerald-600 dark:text-emerald-400 text-sm">Rp {{ formatPrice(p.total_sales) }}</div>
                <div class="text-xs font-semibold text-slate-500">{{ formatQuantity(p.total_qty) }} {{ unitLabel(p.unit) }} terjual</div>
              </div>
            </div>
          </div>
        </div>
      </AppCard>

      <!-- Barang Kurang Laku (Slow Moving / Evaluasi Stok) -->
      <AppCard>
        <div class="flex flex-col gap-1 border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <h3 class="flex items-center gap-1.5 text-rose-600">
            <ExclamationTriangleIcon class="w-5 h-5 text-rose-500" /> Barang Kurang Laku (Slow Moving)
          </h3>
          <span class="text-xs text-secondary">Produk dengan penjualan paling rendah / perlu promo</span>
        </div>
        <div class="flex flex-col">
          <div v-if="!stats.least_products || stats.least_products.length === 0" class="text-sm text-slate-400 italic p-6 text-center">
            Belum ada data evaluasi stok produk
          </div>
          <div v-else class="flex flex-col gap-3">
            <div v-for="(p, index) in stats.least_products" :key="index" class="flex items-center gap-3 p-3 bg-slate-50 dark:bg-slate-900/40 rounded-xl border border-slate-200 dark:border-slate-700/50">
              <div class="w-8 h-8 flex shrink-0 items-center justify-center rounded-lg bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 shadow-sm font-black text-slate-400 text-xs">#{{ index + 1 }}</div>
              <div class="flex flex-col flex-1 min-w-0">
                <span class="top-name font-bold text-slate-800">{{ p.product_name }}</span>
                <span class="top-qty text-xs text-rose-500 font-medium">Sisa Stok: {{ formatQuantity(p.stock ?? 0) }} {{ unitLabel(p.unit) }}</span>
              </div>
              <div class="text-right">
                <div class="font-bold text-slate-700">{{ formatQuantity(p.total_qty) }} {{ unitLabel(p.unit) }} terjual</div>
                <div class="text-xs text-slate-500">Rp {{ formatPrice(p.total_sales) }}</div>
              </div>
            </div>
          </div>
        </div>
      </AppCard>
    </div>

    <!-- Section 3: List Seluruh Barang Laku (Full Detailed Table) -->
    <AppCard class="col-span-1 lg:col-span-2 mt-6">
      <div class="flex justify-between items-center border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
        <div>
          <h3 class="flex items-center gap-1.5">
            <ListBulletIcon class="w-5 h-5 text-indigo-600" /> List Seluruh Barang Laku
          </h3>
          <span class="text-xs text-secondary">Rincian lengkap kinerjaper barang yang telah terjual</span>
        </div>
        <span class="px-2.5 py-1 bg-indigo-100 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-400 rounded-lg text-[10px] uppercase tracking-wider font-bold border border-indigo-200 dark:border-indigo-800/50" v-if="stats.all_sold_products">
          Total {{ stats.all_sold_products.length }} Jenis Produk
        </span>
      </div>

      <div class="overflow-x-auto w-full mt-2">
        <table class="w-full text-left text-sm whitespace-nowrap">
          <thead>
            <tr>
              <th>No</th>
              <th>Nama Produk</th>
              <th>Merk</th>
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
              <td><span class="px-2 py-1 bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400 rounded-md text-[10px] font-bold border border-purple-200 dark:border-purple-800/50">{{ p.artist }}</span></td>
              <td><span class="px-2 py-1 bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400 rounded-md text-[10px] font-bold border border-amber-200 dark:border-amber-800/50">{{ p.product_type }}</span></td>
              <td>Rp {{ formatPrice(p.price || 0) }}</td>
              <td><span class="px-2 py-1 bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300 rounded-full text-xs font-bold border border-slate-200 dark:border-slate-700">{{ formatQuantity(p.total_qty) }} {{ unitLabel(p.unit) }}</span></td>
              <td class="text-right font-bold text-emerald-600">Rp {{ formatPrice(p.total_sales) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </AppCard>

    <!-- Section 4: Breakdown Sales Per Merk & Tipe Produk -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5 mt-5">
      <!-- Sales by Artist / Merk -->
      <AppCard>
        <div class="flex flex-col gap-1 border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <h3 class="flex items-center gap-1.5 text-amber-600">
            <UserIcon class="w-5 h-5" /> Penjualan Per Merk
          </h3>
        </div>
        <div class="flex flex-col">
          <div v-if="!stats.sales_by_artist || stats.sales_by_artist.length === 0" class="text-sm text-slate-400 italic p-6 text-center">
            Belum ada data merk
          </div>
          <div v-else class="flex flex-col gap-3">
            <div v-for="(item, idx) in stats.sales_by_artist" :key="idx" class="flex justify-between items-center p-3 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl shadow-sm hover:shadow-md transition-shadow">
              <div class="flex flex-col flex-1 min-w-0">
                <span class="font-bold text-slate-800 dark:text-slate-100 text-sm">{{ item.name }}</span>
                <span class="text-xs text-slate-500 mt-0.5">{{ formatQuantity(item.total_qty) }} kuantitas terjual (satuan campuran)</span>
              </div>
              <div class="font-bold text-slate-700 dark:text-slate-300">Rp {{ formatPrice(item.total_sales) }}</div>
            </div>
          </div>
        </div>
      </AppCard>

      <!-- Sales by Product Type -->
      <AppCard>
        <div class="flex flex-col gap-1 border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <h3 class="flex items-center gap-1.5 text-pink-600">
            <TagIcon class="w-5 h-5" /> Penjualan Per Tipe Produk
          </h3>
        </div>
        <div class="flex flex-col">
          <div v-if="!stats.sales_by_type || stats.sales_by_type.length === 0" class="text-sm text-slate-400 italic p-6 text-center">
            Belum ada data tipe produk
          </div>
          <div v-else class="flex flex-col gap-3">
            <div v-for="(item, idx) in stats.sales_by_type" :key="idx" class="flex justify-between items-center p-3 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl shadow-sm hover:shadow-md transition-shadow">
              <div class="flex flex-col flex-1 min-w-0">
                <span class="font-bold text-slate-800 dark:text-slate-100 text-sm">{{ item.name }}</span>
                <span class="text-xs text-slate-500 mt-0.5">{{ formatQuantity(item.total_qty) }} kuantitas terjual (satuan campuran)</span>
              </div>
              <div class="font-bold text-slate-700 dark:text-slate-300">Rp {{ formatPrice(item.total_sales) }}</div>
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
const formatQuantity = (val: number): string => new Intl.NumberFormat('id-ID', { maximumFractionDigits: 3 }).format(val || 0);
const unitLabel = (unit?: string): string => unit === 'gram' ? 'gr' : unit === 'liter' ? 'L' : 'pcs';

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
      label: 'Kuantitas Terjual (Satuan Produk)',
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

