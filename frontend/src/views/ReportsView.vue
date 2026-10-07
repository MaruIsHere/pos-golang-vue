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
    <AppCard id="all-sold-products-section" class="col-span-1 lg:col-span-2 mt-6">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
        <div>
          <h3 class="flex items-center gap-1.5 font-bold text-slate-800 dark:text-slate-100">
            <ListBulletIcon class="w-5 h-5 text-indigo-600" /> List Seluruh Barang Laku
          </h3>
          <span class="text-xs text-slate-500">Rincian lengkap kinerja per barang yang telah terjual</span>
        </div>

        <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
          <!-- Search input -->
          <div class="relative flex-1 md:w-56">
            <MagnifyingGlassIcon class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
            <input 
              type="text" 
              v-model="mainTableFilter.search" 
              placeholder="Cari produk..." 
              class="w-full pl-9 pr-3 py-1.5 text-xs font-medium bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500/20 text-slate-800 dark:text-slate-100"
            />
          </div>

          <!-- Active Filter Badge & Reset -->
          <div v-if="mainTableFilter.artist || mainTableFilter.productType" class="flex items-center gap-2">
            <span class="px-2.5 py-1 bg-indigo-50 text-indigo-700 dark:bg-indigo-950/50 dark:text-indigo-400 rounded-lg text-xs font-bold border border-indigo-200 dark:border-indigo-800/40 flex items-center gap-1">
              <FunnelIcon class="w-3.5 h-3.5" />
              <span>{{ mainTableFilter.artist ? `Merk: ${mainTableFilter.artist}` : `Tipe: ${mainTableFilter.productType}` }}</span>
            </span>
            <button class="px-2 py-1 text-xs font-bold text-rose-600 hover:text-rose-700 bg-rose-50 hover:bg-rose-100 dark:bg-rose-900/30 rounded-lg transition-colors" @click="resetMainTableFilter">
              Reset
            </button>
          </div>

          <span class="px-2.5 py-1 bg-indigo-100 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-400 rounded-lg text-[10px] uppercase tracking-wider font-bold border border-indigo-200 dark:border-indigo-800/50 shrink-0">
            Total {{ filteredMainTableProducts.length }} Jenis Produk
          </span>
        </div>
      </div>

      <div class="overflow-x-auto w-full mt-2">
        <table class="w-full text-left text-sm whitespace-nowrap">
          <thead>
            <tr class="text-xs font-bold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700/60">
              <th class="pb-3 pt-1 px-3">No</th>
              <th class="pb-3 pt-1 px-3">Nama Produk</th>
              <th class="pb-3 pt-1 px-3">Merk</th>
              <th class="pb-3 pt-1 px-3">Tipe Produk</th>
              <th class="pb-3 pt-1 px-3">Harga Satuan</th>
              <th class="pb-3 pt-1 px-3">Total Kuantitas</th>
              <th class="pb-3 pt-1 px-3 text-right">Total Sales</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-700/50">
            <tr v-if="filteredMainTableProducts.length === 0">
              <td colspan="7" class="text-center py-6 text-slate-400 italic text-xs">Belum ada transaksi barang terjual yang cocok</td>
            </tr>
            <tr v-else v-for="(p, idx) in filteredMainTableProducts" :key="idx" class="hover:bg-slate-50/80 dark:hover:bg-slate-700/30 transition-colors">
              <td class="py-3 px-3 text-slate-400 font-bold text-xs">{{ idx + 1 }}</td>
              <td class="py-3 px-3 font-bold text-slate-800 dark:text-slate-100">{{ p.product_name }}</td>
              <td class="py-3 px-3"><span class="px-2.5 py-1 bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400 rounded-md text-[10px] font-bold border border-purple-200 dark:border-purple-800/50">{{ p.artist || 'Umum' }}</span></td>
              <td class="py-3 px-3"><span class="px-2.5 py-1 bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400 rounded-md text-[10px] font-bold border border-amber-200 dark:border-amber-800/50">{{ p.product_type || 'Umum' }}</span></td>
              <td class="py-3 px-3 text-slate-700 dark:text-slate-300">Rp {{ formatPrice(p.price || 0) }}</td>
              <td class="py-3 px-3"><span class="px-2.5 py-1 bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300 rounded-full text-xs font-bold border border-slate-200 dark:border-slate-700">{{ formatQuantity(p.total_qty) }} {{ unitLabel(p.unit) }}</span></td>
              <td class="py-3 px-3 text-right font-extrabold text-emerald-600 dark:text-emerald-400">Rp {{ formatPrice(p.total_sales) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </AppCard>

    <!-- Section 4: Breakdown Sales Per Merk & Tipe Produk -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5 mt-5">
      <!-- Sales by Artist / Merk -->
      <AppCard>
        <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <div>
            <h3 class="flex items-center gap-1.5 text-amber-600 font-bold">
              <UserIcon class="w-5 h-5 text-amber-500" /> Penjualan Per Merk
            </h3>
            <span class="text-[11px] text-slate-400 font-medium">Klik item untuk melihat rincian detail barang</span>
          </div>
          <span class="px-2 py-0.5 bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400 text-[10px] font-bold rounded-md border border-amber-200 dark:border-amber-800/40">
            {{ stats.sales_by_artist?.length || 0 }} Merk
          </span>
        </div>
        <div class="flex flex-col">
          <div v-if="!stats.sales_by_artist || stats.sales_by_artist.length === 0" class="text-sm text-slate-400 italic p-6 text-center">
            Belum ada data merk
          </div>
          <div v-else class="flex flex-col gap-3">
            <div 
              v-for="(item, idx) in stats.sales_by_artist" 
              :key="idx" 
              class="flex justify-between items-center p-3.5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 hover:border-amber-400 dark:hover:border-amber-500 rounded-xl shadow-xs hover:shadow-md transition-all cursor-pointer group"
              @click="openDetailModal('artist', item.name)"
            >
              <div class="flex flex-col flex-1 min-w-0 pr-2">
                <div class="flex items-center gap-2">
                  <span class="font-bold text-slate-800 dark:text-slate-100 text-sm group-hover:text-amber-600 dark:group-hover:text-amber-400 transition-colors">{{ item.name }}</span>
                  <span class="px-2 py-0.5 text-[10px] font-extrabold bg-amber-50 dark:bg-amber-950/40 text-amber-700 dark:text-amber-400 rounded-md border border-amber-200 dark:border-amber-800/50">Detail</span>
                </div>
                <span class="text-xs text-slate-500 mt-1">{{ formatQuantity(item.total_qty) }} kuantitas terjual (satuan campuran)</span>
              </div>
              <div class="flex items-center gap-3 shrink-0">
                <div class="font-extrabold text-slate-800 dark:text-slate-100 text-sm">Rp {{ formatPrice(item.total_sales) }}</div>
                <ChevronRightIcon class="w-4 h-4 text-slate-400 group-hover:text-amber-500 group-hover:translate-x-0.5 transition-all" />
              </div>
            </div>
          </div>
        </div>
      </AppCard>

      <!-- Sales by Product Type -->
      <AppCard>
        <div class="flex items-center justify-between border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <div>
            <h3 class="flex items-center gap-1.5 text-pink-600 font-bold">
              <TagIcon class="w-5 h-5 text-pink-500" /> Penjualan Per Tipe Produk
            </h3>
            <span class="text-[11px] text-slate-400 font-medium">Klik item untuk melihat rincian detail barang</span>
          </div>
          <span class="px-2 py-0.5 bg-pink-50 text-pink-700 dark:bg-pink-950/40 dark:text-pink-400 text-[10px] font-bold rounded-md border border-pink-200 dark:border-pink-800/40">
            {{ stats.sales_by_type?.length || 0 }} Tipe
          </span>
        </div>
        <div class="flex flex-col">
          <div v-if="!stats.sales_by_type || stats.sales_by_type.length === 0" class="text-sm text-slate-400 italic p-6 text-center">
            Belum ada data tipe produk
          </div>
          <div v-else class="flex flex-col gap-3">
            <div 
              v-for="(item, idx) in stats.sales_by_type" 
              :key="idx" 
              class="flex justify-between items-center p-3.5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 hover:border-pink-400 dark:hover:border-pink-500 rounded-xl shadow-xs hover:shadow-md transition-all cursor-pointer group"
              @click="openDetailModal('product_type', item.name)"
            >
              <div class="flex flex-col flex-1 min-w-0 pr-2">
                <div class="flex items-center gap-2">
                  <span class="font-bold text-slate-800 dark:text-slate-100 text-sm group-hover:text-pink-600 dark:group-hover:text-pink-400 transition-colors">{{ item.name }}</span>
                  <span class="px-2 py-0.5 text-[10px] font-extrabold bg-pink-50 dark:bg-pink-950/40 text-pink-700 dark:text-pink-400 rounded-md border border-pink-200 dark:border-pink-800/50">Detail</span>
                </div>
                <span class="text-xs text-slate-500 mt-1">{{ formatQuantity(item.total_qty) }} kuantitas terjual (satuan campuran)</span>
              </div>
              <div class="flex items-center gap-3 shrink-0">
                <div class="font-extrabold text-slate-800 dark:text-slate-100 text-sm">Rp {{ formatPrice(item.total_sales) }}</div>
                <ChevronRightIcon class="w-4 h-4 text-slate-400 group-hover:text-pink-500 group-hover:translate-x-0.5 transition-all" />
              </div>
            </div>
          </div>
        </div>
      </AppCard>
    </div>

    <!-- Modal Popup Detail Penjualan Per Merk / Tipe -->
    <div v-if="activeDetailModal.isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6 bg-slate-900/60 backdrop-blur-xs overflow-y-auto animate-fade-in" @click.self="closeDetailModal">
      <div class="relative w-full max-w-4xl bg-white dark:bg-slate-800 rounded-3xl shadow-2xl border border-slate-200 dark:border-slate-700 overflow-hidden my-8">
        
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-5 sm:p-6 border-b border-slate-100 dark:border-slate-700/80 bg-slate-50/50 dark:bg-slate-800/80">
          <div class="flex items-center gap-3">
            <div class="p-3 rounded-2xl" :class="activeDetailModal.type === 'artist' ? 'bg-amber-100 text-amber-600 dark:bg-amber-900/40 dark:text-amber-400' : 'bg-pink-100 text-pink-600 dark:bg-pink-900/40 dark:text-pink-400'">
              <UserIcon v-if="activeDetailModal.type === 'artist'" class="w-6 h-6" />
              <TagIcon v-else class="w-6 h-6" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span class="text-xs font-bold uppercase tracking-wider text-slate-400">
                  {{ activeDetailModal.type === 'artist' ? 'Rincian Penjualan Per Merk' : 'Rincian Penjualan Per Tipe Produk' }}
                </span>
              </div>
              <h3 class="text-xl font-black text-slate-800 dark:text-slate-100 flex items-center gap-2">
                <span>{{ activeDetailModal.name }}</span>
              </h3>
            </div>
          </div>
          <button class="p-2 rounded-full text-slate-400 hover:text-slate-600 hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors" @click="closeDetailModal">
            <XMarkIcon class="w-6 h-6" />
          </button>
        </div>

        <!-- Modal Stats Summary Cards -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 p-4 sm:p-6 bg-slate-100/60 dark:bg-slate-900/40 border-b border-slate-100 dark:border-slate-700/60">
          <div class="p-3.5 bg-white dark:bg-slate-800 rounded-2xl border border-slate-200/60 dark:border-slate-700/60 shadow-2xs">
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Total Sales (Omset)</span>
            <span class="text-lg font-black text-emerald-600 dark:text-emerald-400">Rp {{ formatPrice(filteredDetailStats.totalSales) }}</span>
          </div>
          <div class="p-3.5 bg-white dark:bg-slate-800 rounded-2xl border border-slate-200/60 dark:border-slate-700/60 shadow-2xs">
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Total Qty Terjual</span>
            <span class="text-lg font-black text-indigo-600 dark:text-indigo-400">{{ formatQuantity(filteredDetailStats.totalQty) }} Item</span>
          </div>
          <div class="p-3.5 bg-white dark:bg-slate-800 rounded-2xl border border-slate-200/60 dark:border-slate-700/60 shadow-2xs">
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block">Varian Jenis Produk</span>
            <span class="text-lg font-black text-slate-800 dark:text-slate-100">{{ filteredDetailStats.products.length }} Jenis</span>
          </div>
        </div>

        <!-- Modal Table Content -->
        <div class="p-4 sm:p-6">
          <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mb-4">
            <div>
              <h4 class="font-bold text-slate-800 dark:text-slate-100 text-sm flex items-center gap-1.5">
                <ListBulletIcon class="w-4 h-4 text-indigo-500" /> List Barang Laku - {{ activeDetailModal.name }}
              </h4>
              <p class="text-xs text-slate-400">Rincian lengkap kinerja per barang yang telah terjual</p>
            </div>
            
            <div class="flex items-center gap-2 w-full sm:w-auto">
              <div class="relative flex-1 sm:w-56">
                <MagnifyingGlassIcon class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
                <input 
                  type="text" 
                  v-model="detailSearchQuery" 
                  placeholder="Cari nama barang..." 
                  class="w-full pl-9 pr-3 py-1.5 text-xs font-medium bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500/20 text-slate-800 dark:text-slate-100"
                />
              </div>
              <span class="px-2.5 py-1 bg-indigo-100 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-400 rounded-lg text-[10px] uppercase tracking-wider font-bold shrink-0">
                TOTAL {{ filteredDetailProducts.length }} JENIS PRODUK
              </span>
            </div>
          </div>

          <!-- Table Container -->
          <div class="overflow-x-auto rounded-2xl border border-slate-200/80 dark:border-slate-700/80 max-h-[350px] overflow-y-auto">
            <table class="w-full text-left text-sm whitespace-nowrap">
              <thead class="sticky top-0 bg-slate-50 dark:bg-slate-800 text-slate-500 dark:text-slate-400 text-xs font-bold border-b border-slate-200/80 dark:border-slate-700/80 z-10">
                <tr>
                  <th class="py-3 px-4">No</th>
                  <th class="py-3 px-4">Nama Produk</th>
                  <th class="py-3 px-4">Merk</th>
                  <th class="py-3 px-4">Tipe Produk</th>
                  <th class="py-3 px-4">Harga Satuan</th>
                  <th class="py-3 px-4">Total Kuantitas</th>
                  <th class="py-3 px-4 text-right">Total Sales</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100 dark:divide-slate-700/50">
                <tr v-if="filteredDetailProducts.length === 0">
                  <td colspan="7" class="py-8 text-center text-slate-400 italic text-xs">
                    Tidak ada produk terjual yang cocok
                  </td>
                </tr>
                <tr v-else v-for="(p, idx) in filteredDetailProducts" :key="idx" class="hover:bg-slate-50/80 dark:hover:bg-slate-700/30 transition-colors">
                  <td class="py-3.5 px-4 font-bold text-slate-400 text-xs">{{ idx + 1 }}</td>
                  <td class="py-3.5 px-4 font-bold text-slate-800 dark:text-slate-100">{{ p.product_name }}</td>
                  <td class="py-3.5 px-4">
                    <span class="px-2.5 py-1 bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400 rounded-md text-[10px] font-bold border border-purple-200 dark:border-purple-800/50">
                      {{ p.artist || 'Umum' }}
                    </span>
                  </td>
                  <td class="py-3.5 px-4">
                    <span class="px-2.5 py-1 bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400 rounded-md text-[10px] font-bold border border-amber-200 dark:border-amber-800/50">
                      {{ p.product_type || 'Umum' }}
                    </span>
                  </td>
                  <td class="py-3.5 px-4 text-slate-700 dark:text-slate-300 font-medium">Rp {{ formatPrice(p.price || 0) }}</td>
                  <td class="py-3.5 px-4">
                    <span class="px-2.5 py-1 bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300 rounded-full text-xs font-bold border border-slate-200 dark:border-slate-700">
                      {{ formatQuantity(p.total_qty) }} {{ unitLabel(p.unit) }}
                    </span>
                  </td>
                  <td class="py-3.5 px-4 text-right font-extrabold text-emerald-600 dark:text-emerald-400">Rp {{ formatPrice(p.total_sales) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="flex items-center justify-between p-4 sm:p-5 border-t border-slate-100 dark:border-slate-700/80 bg-slate-50/50 dark:bg-slate-800/80">
          <button class="px-4 py-2 text-xs font-bold text-indigo-600 dark:text-indigo-400 hover:bg-indigo-50 dark:hover:bg-indigo-950/50 rounded-xl transition-colors flex items-center gap-1.5" @click="applyFilterToMainTable">
            <FunnelIcon class="w-4 h-4" />
            <span>Filter di Tabel Utama Page</span>
          </button>

          <button class="px-5 py-2 text-xs font-bold text-white bg-slate-800 hover:bg-slate-900 dark:bg-slate-700 dark:hover:bg-slate-600 rounded-xl transition-colors shadow-sm" @click="closeDetailModal">
            Tutup
          </button>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppCard from '@/components/ui/AppCard.vue';
import AppInput from '@/components/ui/AppInput.vue';
import { ref, computed, onMounted } from 'vue';
import api from '@/utils/api';
import type { DashboardStats, ProductSalesStat } from '../types';
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
  ListBulletIcon,
  ChevronRightIcon,
  FunnelIcon,
  MagnifyingGlassIcon,
  XMarkIcon
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

// Modal detail state
interface DetailModalState {
  isOpen: boolean;
  type: 'artist' | 'product_type';
  name: string;
}

const activeDetailModal = ref<DetailModalState>({
  isOpen: false,
  type: 'artist',
  name: ''
});

const detailSearchQuery = ref('');

// Main table filter state
const mainTableFilter = ref<{
  search: string;
  artist: string | null;
  productType: string | null;
}>({
  search: '',
  artist: null,
  productType: null
});

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);
const formatQuantity = (val: number): string => new Intl.NumberFormat('id-ID', { maximumFractionDigits: 3 }).format(val || 0);
const unitLabel = (unit?: string): string => unit === 'gram' ? 'gr' : unit === 'liter' ? 'L' : 'pcs';

// Matching helpers for brand / type (handles 'Umum' and 'Lainnya / Umum')
const isMatchingArtist = (artistInProduct?: string, filterArtistName?: string): boolean => {
  if (!filterArtistName) return true;
  const pArtist = (artistInProduct || '').trim().toLowerCase();
  const fArtist = filterArtistName.trim().toLowerCase();
  
  const isDefaultP = !pArtist || pArtist === 'umum' || pArtist === 'lainnya / umum' || pArtist === '-';
  const isDefaultF = !fArtist || fArtist === 'umum' || fArtist === 'lainnya / umum' || fArtist === '-';

  if (isDefaultF) return isDefaultP;
  return pArtist === fArtist;
};

const isMatchingProductType = (typeInProduct?: string, filterTypeName?: string): boolean => {
  if (!filterTypeName) return true;
  const pType = (typeInProduct || '').trim().toLowerCase();
  const fType = filterTypeName.trim().toLowerCase();

  const isDefaultP = !pType || pType === 'umum' || pType === 'lainnya / umum' || pType === '-';
  const isDefaultF = !fType || fType === 'umum' || fType === 'lainnya / umum' || fType === '-';

  if (isDefaultF) return isDefaultP;
  return pType === fType;
};

const openDetailModal = (type: 'artist' | 'product_type', name: string) => {
  activeDetailModal.value = {
    isOpen: true,
    type,
    name
  };
  detailSearchQuery.value = '';
};

const closeDetailModal = () => {
  activeDetailModal.value.isOpen = false;
};

const applyFilterToMainTable = () => {
  if (activeDetailModal.value.type === 'artist') {
    mainTableFilter.value.artist = activeDetailModal.value.name;
    mainTableFilter.value.productType = null;
  } else {
    mainTableFilter.value.productType = activeDetailModal.value.name;
    mainTableFilter.value.artist = null;
  }
  closeDetailModal();
  
  // Smooth scroll to main table section
  const element = document.getElementById('all-sold-products-section');
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' });
  }
};

const resetMainTableFilter = () => {
  mainTableFilter.value.search = '';
  mainTableFilter.value.artist = null;
  mainTableFilter.value.productType = null;
};

// Filtered products for modal
const rawDetailProducts = computed<ProductSalesStat[]>(() => {
  if (!activeDetailModal.value.isOpen) return [];
  const allProds = stats.value.all_sold_products || [];
  if (activeDetailModal.value.type === 'artist') {
    return allProds.filter(p => isMatchingArtist(p.artist, activeDetailModal.value.name));
  } else {
    return allProds.filter(p => isMatchingProductType(p.product_type, activeDetailModal.value.name));
  }
});

const filteredDetailProducts = computed<ProductSalesStat[]>(() => {
  const q = detailSearchQuery.value.trim().toLowerCase();
  if (!q) return rawDetailProducts.value;
  return rawDetailProducts.value.filter(p => 
    p.product_name.toLowerCase().includes(q) ||
    (p.artist && p.artist.toLowerCase().includes(q)) ||
    (p.product_type && p.product_type.toLowerCase().includes(q))
  );
});

const filteredDetailStats = computed(() => {
  const prods = rawDetailProducts.value;
  const totalSales = prods.reduce((sum, p) => sum + (p.total_sales || 0), 0);
  const totalQty = prods.reduce((sum, p) => sum + (p.total_qty || 0), 0);
  return {
    products: prods,
    totalSales,
    totalQty
  };
});

// Filtered products for Section 3 main table
const filteredMainTableProducts = computed<ProductSalesStat[]>(() => {
  const allProds = stats.value.all_sold_products || [];
  return allProds.filter(p => {
    const matchesSearch = !mainTableFilter.value.search || 
      p.product_name.toLowerCase().includes(mainTableFilter.value.search.toLowerCase());
    const matchesArtist = !mainTableFilter.value.artist || 
      isMatchingArtist(p.artist, mainTableFilter.value.artist);
    const matchesType = !mainTableFilter.value.productType || 
      isMatchingProductType(p.product_type, mainTableFilter.value.productType);

    return matchesSearch && matchesArtist && matchesType;
  });
});

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

