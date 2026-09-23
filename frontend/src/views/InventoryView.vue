<template>
  <div class="inventory-page">
    <div class="page-header glass-panel">
      <div class="header-title">
        <h2>📦 Manajemen Inventoris & Mutasi Stok</h2>
        <p>Kelola Penerimaan Barang (Restock) dan Pengeluaran Barang (Barang Rusak / Hilang / Expired)</p>
      </div>
    </div>

    <!-- Sub Navigation Tabs -->
    <div class="inventory-tabs">
      <button 
        class="tab-btn" 
        :class="{ active: activeTab === 'receive' }"
        @click="activeTab = 'receive'"
      >
        <span>📥 Transaksi Penerimaan (Stock In)</span>
      </button>

      <button 
        class="tab-btn" 
        :class="{ active: activeTab === 'issue' }"
        @click="activeTab = 'issue'"
      >
        <span>📤 Transaksi Pengeluaran (Stock Out)</span>
      </button>

      <button 
        class="tab-btn" 
        :class="{ active: activeTab === 'history' }"
        @click="activeTab = 'history'"
      >
        <span>📋 Riwayat Mutasi Stok</span>
      </button>
    </div>

    <!-- Tab 1: Penerimaan Barang (Stock In / Receive) -->
    <div v-if="activeTab === 'receive'" class="inventory-card glass-panel">
      <div class="card-header">
        <h3>📥 Penerimaan Barang dari Supplier / Pabrik</h3>
      </div>

      <form @submit.prevent="submitStockMovement('in')" class="card-body">
        <div class="form-grid">
          <div class="form-group">
            <label class="form-label">Pilih Produk</label>
            <select class="form-control" v-model.number="receiveForm.product_id" required>
              <option value="" disabled>-- Pilih Produk --</option>
              <option v-for="p in products" :key="p.id" :value="p.id">
                {{ p.name }} (Stok Saat Ini: {{ p.stock }})
              </option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label">Kuantitas Masuk (+)</label>
            <input 
              type="number" 
              class="form-control" 
              v-model.number="receiveForm.quantity" 
              min="1" 
              placeholder="cth: 50" 
              required 
            />
          </div>

          <div class="form-group">
            <label class="form-label">Alasan / Sumber</label>
            <select class="form-control" v-model="receiveForm.reason">
              <option value="pembelian_supplier">Pembelian dari Supplier</option>
              <option value="produksi_sendiri">Hasil Produksi Mandiri</option>
              <option value="penyesuaian_opname">Penyesuaian Stok Opname (+)</option>
            </select>
          </div>

          <div class="form-group full-width">
            <label class="form-label">Catatan / No. Surat Jalan / Supplier</label>
            <input 
              type="text" 
              class="form-control" 
              v-model="receiveForm.notes" 
              placeholder="cth: Restock PT Maju Bersama - SJ #9921" 
            />
          </div>
        </div>

        <button type="submit" class="btn btn-success btn-submit" :disabled="isSubmitting">
          {{ isSubmitting ? 'Memproses...' : '📥 Simpan Penerimaan Barang & Tambah Stok' }}
        </button>
      </form>
    </div>

    <!-- Tab 2: Pengeluaran Barang (Stock Out / Issue) -->
    <div v-else-if="activeTab === 'issue'" class="inventory-card glass-panel">
      <div class="card-header">
        <h3>📤 Pengeluaran Barang Non-Penjualan (Rusak / Hilang / Promosi)</h3>
      </div>

      <form @submit.prevent="submitStockMovement('out')" class="card-body">
        <div class="form-grid">
          <div class="form-group">
            <label class="form-label">Pilih Produk</label>
            <select class="form-control" v-model.number="issueForm.product_id" required>
              <option value="" disabled>-- Pilih Produk --</option>
              <option v-for="p in products" :key="p.id" :value="p.id">
                {{ p.name }} (Stok Tersedia: {{ p.stock }})
              </option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label">Kuantitas Keluar (-)</label>
            <input 
              type="number" 
              class="form-control" 
              v-model.number="issueForm.quantity" 
              min="1" 
              placeholder="cth: 5" 
              required 
            />
          </div>

          <div class="form-group">
            <label class="form-label">Alasan Pengeluaran</label>
            <select class="form-control" v-model="issueForm.reason">
              <option value="barang_rusak">Barang Rusak / Cacat</option>
              <option value="barang_hilang">Barang Hilang / Selisih Stok</option>
              <option value="expired">Kadaluarsa (Expired)</option>
              <option value="promosi">Kebutuhan Promosi / Sampling</option>
              <option value="pemakaian_sendiri">Pemakaian Internal Toko</option>
            </select>
          </div>

          <div class="form-group full-width">
            <label class="form-label">Catatan Pengeluaran</label>
            <input 
              type="text" 
              class="form-control" 
              v-model="issueForm.notes" 
              placeholder="cth: Kemasan pecah saat pemindahan di gudang" 
            />
          </div>
        </div>

        <button type="submit" class="btn btn-danger btn-submit" :disabled="isSubmitting">
          {{ isSubmitting ? 'Memproses...' : '📤 Simpan Pengeluaran Barang & Potong Stok' }}
        </button>
      </form>
    </div>

    <!-- Tab 3: Riwayat Mutasi Stok -->
    <div v-else class="inventory-card glass-panel">
      <div class="card-header">
        <h3>📋 Log Riwayat Mutasi Masuk & Keluar Stok</h3>
      </div>

      <div class="card-body">
        <div v-if="isLoadingMovements" class="loading-box">
          <div class="spinner"></div>
          <p>Memuat data riwayat mutasi stok...</p>
        </div>

        <div v-else-if="movements.length === 0" class="empty-box">
          <span class="empty-icon">📂</span>
          <p>Belum ada riwayat mutasi stok barang.</p>
        </div>

        <div v-else class="table-wrapper">
          <table class="movement-table">
            <thead>
              <tr>
                <th>Waktu & Tanggal</th>
                <th>Jenis Mutasi</th>
                <th>Produk</th>
                <th>Jumlah</th>
                <th>Alasan</th>
                <th>Catatan</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="m in movements" :key="m.id">
                <td class="date-col">{{ formatDate(m.created_at) }}</td>
                <td>
                  <span class="type-badge" :class="m.type">
                    {{ m.type === 'in' ? '📥 Masuk (In)' : '📤 Keluar (Out)' }}
                  </span>
                </td>
                <td>
                  <strong class="product-name">{{ m.product ? m.product.name : 'Produk ID ' + m.product_id }}</strong>
                </td>
                <td>
                  <span class="qty-badge" :class="m.type">
                    {{ m.type === 'in' ? '+' : '-' }}{{ m.quantity }} Unit
                  </span>
                </td>
                <td>
                  <span class="reason-tag">{{ formatReason(m.reason ?? '') }}</span>
                </td>
                <td class="notes-col">{{ m.notes || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { Product, StockMovement } from '../types';

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
    const res = await fetch('/api/products');
    if (res.ok) products.value = await res.json();
  } catch (err) {
    console.error(err);
  }
};

const fetchMovements = async () => {
  isLoadingMovements.value = true;
  try {
    const res = await fetch('/api/stock-movements');
    if (res.ok) movements.value = await res.json();
  } catch (err) {
    console.error(err);
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

    const res = await fetch('/api/stock-movements', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });

    if (res.ok) {
      alert(`Berhasil menyimpan transaksi ${type === 'in' ? 'penerimaan' : 'pengeluaran'} barang!`);
      if (type === 'in') {
        receiveForm.value = { product_id: '', quantity: 1, reason: 'pembelian_supplier', notes: '' };
      } else {
        issueForm.value = { product_id: '', quantity: 1, reason: 'barang_rusak', notes: '' };
      }
      fetchProducts();
      fetchMovements();
      activeTab.value = 'history';
    } else {
      const errData = await res.json();
      alert('Gagal: ' + (errData.error || 'Terjadi kesalahan'));
    }
  } catch (err) {
    alert('Koneksi error: ' + (err as Error).message);
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<style scoped>
.inventory-page {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.page-header {
  padding: 1.25rem 1.5rem;
}

.header-title h2 {
  font-size: 1.25rem;
  font-weight: 800;
}

.inventory-tabs {
  display: flex;
  gap: 0.5rem;
  overflow-x: auto;
}

.tab-btn {
  padding: 0.75rem 1.25rem;
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-family: var(--font-family);
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s;
}

.tab-btn:hover {
  background: var(--bg-card-hover);
  color: var(--text-primary);
}

.tab-btn.active {
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-purple));
  color: #ffffff;
  border-color: transparent;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

.inventory-card {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.card-header {
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.85rem;
}

.card-header h3 {
  font-size: 1.1rem;
  font-weight: 800;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1rem;
}

@media (min-width: 640px) {
  .form-grid {
    grid-template-columns: 1fr 1fr;
  }
  .full-width {
    grid-column: span 2;
  }
}

.btn-submit {
  width: 100%;
  padding: 0.85rem;
  font-size: 0.95rem;
  margin-top: 0.75rem;
}

.table-wrapper {
  overflow-x: auto;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
}

.movement-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}

.movement-table th, .movement-table td {
  padding: 0.75rem 1rem;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

.movement-table th {
  background: rgba(15, 23, 42, 0.8);
  color: var(--text-secondary);
  font-weight: 700;
}

.type-badge {
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 800;
}

.type-badge.in {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.type-badge.out {
  background: rgba(239, 68, 68, 0.2);
  color: #f87171;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.product-name {
  color: #ffffff;
}

.qty-badge {
  font-weight: 800;
  font-size: 0.9rem;
}

.qty-badge.in { color: #34d399; }
.qty-badge.out { color: #f87171; }

.reason-tag {
  background: rgba(30, 41, 59, 0.8);
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-size: 0.78rem;
  color: #a5b4fc;
}

.date-col, .notes-col {
  color: var(--text-muted);
  font-size: 0.8rem;
}

.loading-box, .empty-box {
  padding: 3rem;
  text-align: center;
  color: var(--text-muted);
}
</style>
