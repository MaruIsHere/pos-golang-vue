<template>
  <div class="flex flex-col gap-5">
    <div class="flex justify-between items-center p-5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <div class="flex flex-col gap-1">
        <h2 class="text-xl font-bold text-slate-900 dark:text-slate-100">Manajemen Inventoris & Mutasi Stok</h2>
        <p class="text-sm text-slate-500 dark:text-slate-400">Kelola Penerimaan Barang (Restock) dan Pengeluaran Barang (Barang Rusak / Hilang / Expired)</p>
      </div>
    </div>

    <!-- Sub Navigation Tabs -->
    <div class="flex gap-2 overflow-x-auto shrink-0">
      <AppButton variant="secondary" 
        class="" 
        :class="{ active: activeTab === 'receive' }"
        @click="activeTab = 'receive'"
      >
        <ArrowDownTrayIcon class="w-4 h-4 inline-block mr-1.5" />
        <span>Transaksi Penerimaan (Stock In)</span>
      </AppButton>

      <AppButton variant="secondary" 
        class="" 
        :class="{ active: activeTab === 'issue' }"
        @click="activeTab = 'issue'"
      >
        <ArrowUpTrayIcon class="w-4 h-4 inline-block mr-1.5" />
        <span>Transaksi Pengeluaran (Stock Out)</span>
      </AppButton>

      <AppButton variant="secondary" 
        class="" 
        :class="{ active: activeTab === 'history' }"
        @click="activeTab = 'history'"
      >
        <ClockIcon class="w-4 h-4 inline-block mr-1.5" />
        <span>Riwayat Mutasi Stok</span>
      </AppButton>
    </div>

    <!-- Tab 1: Penerimaan Barang (Stock In / Receive) -->
    <div v-if="activeTab === 'receive'" class="flex flex-col gap-4 p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <div class="border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-900 dark:text-slate-100">
        <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Penerimaan Barang dari Supplier / Pabrik</h3>
      </div>

      <form @submit.prevent="submitStockMovement('in')" class="card-body">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="form-group">
            <label class="form-label">Pilih Produk</label>
            <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model.number="receiveForm.product_id" required>
              <option value="" disabled>-- Pilih Produk --</option>
              <option v-for="p in products" :key="p.id" :value="p.id">
                {{ p.name }}{{ p.artist ? ' (' + p.artist + ')' : '' }} (Stok Saat Ini: {{ p.stock }})
              </option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label">Kuantitas Masuk (+)</label>
            <AppInput 
              type="number" 
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" 
              v-model.number="receiveForm.quantity" 
              min="1" 
              placeholder="cth: 50" 
              required 
            />
          </div>

          <div class="form-group">
            <label class="form-label">Alasan / Sumber</label>
            <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="receiveForm.reason">
              <option value="pembelian_supplier">Pembelian dari Supplier</option>
              <option value="produksi_sendiri">Hasil Produksi Mandiri</option>
              <option value="penyesuaian_opname">Penyesuaian Stok Opname (+)</option>
            </select>
          </div>

          <div class="form-group col-span-1 sm:col-span-2">
            <label class="form-label">Catatan / No. Surat Jalan / Supplier</label>
            <AppInput 
              type="text" 
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" 
              v-model="receiveForm.notes" 
              placeholder="cth: Restock PT Maju Bersama - SJ #9921" 
            />
          </div>
        </div>

        <AppButton variant="success" type="submit" class="" :disabled="isSubmitting">
          {{ isSubmitting ? 'Memproses...' : 'Simpan Penerimaan Barang & Tambah Stok' }}
        </AppButton>
      </form>
    </div>

    <!-- Tab 2: Pengeluaran Barang (Stock Out / Issue) -->
    <div v-else-if="activeTab === 'issue'" class="flex flex-col gap-4 p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <div class="border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-900 dark:text-slate-100">
        <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Pengeluaran Barang Non-Penjualan (Rusak / Hilang / Promosi)</h3>
      </div>

      <form @submit.prevent="submitStockMovement('out')" class="card-body">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="form-group">
            <label class="form-label">Pilih Produk</label>
            <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model.number="issueForm.product_id" required>
              <option value="" disabled>-- Pilih Produk --</option>
              <option v-for="p in products" :key="p.id" :value="p.id">
                {{ p.name }}{{ p.artist ? ' (' + p.artist + ')' : '' }} (Stok Tersedia: {{ p.stock }})
              </option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label">Kuantitas Keluar (-)</label>
            <AppInput 
              type="number" 
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" 
              v-model.number="issueForm.quantity" 
              min="1" 
              placeholder="cth: 5" 
              required 
            />
          </div>

          <div class="form-group">
            <label class="form-label">Alasan Pengeluaran</label>
            <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="issueForm.reason">
              <option value="barang_rusak">Barang Rusak / Cacat</option>
              <option value="barang_hilang">Barang Hilang / Selisih Stok</option>
              <option value="expired">Kadaluarsa (Expired)</option>
              <option value="promosi">Kebutuhan Promosi / Sampling</option>
              <option value="pemakaian_sendiri">Pemakaian Internal Toko</option>
            </select>
          </div>

          <div class="form-group col-span-1 sm:col-span-2">
            <label class="form-label">Catatan Pengeluaran</label>
            <AppInput 
              type="text" 
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" 
              v-model="issueForm.notes" 
              placeholder="cth: Kemasan pecah saat pemindahan di gudang" 
            />
          </div>
        </div>

        <AppButton variant="danger" type="submit" class="" :disabled="isSubmitting">
          {{ isSubmitting ? 'Memproses...' : 'Simpan Pengeluaran Barang & Potong Stok' }}
        </AppButton>
      </form>
    </div>

    <!-- Tab 3: Riwayat Mutasi Stok -->
    <div v-else class="flex flex-col gap-4 p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <div class="border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-900 dark:text-slate-100">
        <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Log Riwayat Mutasi Masuk & Keluar Stok</h3>
      </div>

      <div class="card-body">
        <div v-if="isLoadingMovements" class="p-12 text-center text-slate-500">
          <div class="w-8 h-8 border-4 border-slate-200 dark:border-slate-700 border-t-indigo-600 rounded-full animate-spin mx-auto mb-2"></div>
          <p class="text-sm text-slate-500 dark:text-slate-400">Memuat data riwayat mutasi stok...</p>
        </div>

        <div v-else-if="movements.length === 0" class="p-12 text-center text-slate-500">
          <p class="text-sm text-slate-500 dark:text-slate-400">Belum ada riwayat mutasi stok barang.</p>
        </div>

        <div v-else class="overflow-x-auto p-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl">
          <table class="w-full text-left text-sm">
            <thead>
              <tr>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Waktu & Tanggal</th>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Jenis Mutasi</th>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Produk</th>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Jumlah</th>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Alasan</th>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Catatan</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="m in movements" :key="m.id">
                <td class="text-xs text-slate-500">{{ formatDate(m.created_at) }}</td>
                <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
                  <span class="px-2.5 py-0.5 rounded-full text-xs font-bold" :class="m.type">
                    {{ m.type === 'in' ? 'Masuk (In)' : 'Keluar (Out)' }}
                  </span>
                </td>
                <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
                  <strong class="product-name">{{ m.product ? m.product.name : 'Produk ID ' + m.product_id }}</strong>
                </td>
                <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
                  <span class="font-bold text-sm" :class="m.type">
                    {{ m.type === 'in' ? '+' : '-' }}{{ m.quantity }} Unit
                  </span>
                </td>
                <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
                  <span class="px-2 py-0.5 rounded text-xs bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400">{{ formatReason(m.reason ?? '') }}</span>
                </td>
                <td class="text-xs text-slate-500">{{ m.notes || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppInput from '@/components/ui/AppInput.vue';

import { ref, onMounted } from 'vue';
import api from '@/utils/api';
import type { Product, StockMovement } from '../types';
import { ArrowDownTrayIcon, ArrowUpTrayIcon, ClockIcon } from '@heroicons/vue/24/outline';

const activeTab = ref('receive');
const products = ref<Product[]>([]);
const movements = ref<StockMovement[]>([]);
const isLoadingMovements = ref(false);
const isSubmitting = ref(false);

interface StockForm {
  product_id: string | number;
  quantity: number;
  reason: string;
  notes: string;
}

const receiveForm = ref<StockForm>({
  product_id: '',
  quantity: 1,
  reason: 'pembelian_supplier',
  notes: ''
});

const issueForm = ref<StockForm>({
  product_id: '',
  quantity: 1,
  reason: 'barang_rusak',
  notes: ''
});

const formatDate = (str?: string): string => {
  if (!str) return '-';
  return new Date(str).toLocaleString('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  });
};

const formatReason = (reason: string): string => {
  const map: Record<string, string> = {
    pembelian_supplier: 'Pembelian Supplier',
    produksi_sendiri: 'Produksi Mandiri',
    penyesuaian_opname: 'Penyesuaian Opname',
    barang_rusak: 'Barang Rusak',
    barang_hilang: 'Barang Hilang',
    expired: 'Kadaluarsa (Expired)',
    promosi: 'Promosi / Sampling',
    pemakaian_sendiri: 'Pemakaian Internal',
    retur_penjualan: 'Retur Penjualan'
  };
  return map[reason] || reason;
};

const fetchProducts = async () => {
  try {
    const res = await api.get('/products');
    products.value = res.data;
  } catch (err: any) {
    console.error(err.response?.data?.error || err.message || 'Error occurred');
  }
};

const fetchMovements = async () => {
  isLoadingMovements.value = true;
  try {
    const res = await api.get('/stock-movements');
    movements.value = res.data;
  } catch (err: any) {
    console.error(err.response?.data?.error || err.message || 'Error occurred');
  } finally {
    isLoadingMovements.value = false;
  }
};

onMounted(() => {
  fetchProducts();
  fetchMovements();
});

const submitStockMovement = async (type: string): Promise<void> => {
  const form = type === 'in' ? receiveForm.value : issueForm.value;
  if (!form.product_id || form.quantity <= 0) return;

  isSubmitting.value = true;
  try {
    const payload = {
      product_id: form.product_id,
      type,
      quantity: form.quantity,
      reason: form.reason,
      notes: form.notes
    };

    await api.post('/stock-movements', payload);
    alert(`Berhasil menyimpan transaksi ${type === 'in' ? 'penerimaan' : 'pengeluaran'} barang!`);
    if (type === 'in') {
      receiveForm.value = { product_id: '', quantity: 1, reason: 'pembelian_supplier', notes: '' };
    } else {
      issueForm.value = { product_id: '', quantity: 1, reason: 'barang_rusak', notes: '' };
    }
    fetchProducts();
    fetchMovements();
    activeTab.value = 'history';
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    alert('Koneksi error: ' + errMsg);
  } finally {
    isSubmitting.value = false;
  }
};
</script>


