<template>
  <div class="orders-page">
    <div class="page-header glass-panel">
      <div class="header-title">
        <h2>📜 Riwayat Transaksi</h2>
        <p>Daftar transaksi penjualan dan cetak ulang struk</p>
      </div>

      <button class="btn btn-secondary" @click="fetchOrders">
        <span>🔄 Refresh</span>
      </button>
    </div>

    <!-- Orders Table -->
    <div class="table-container glass-panel">
      <table class="data-table">
        <thead>
          <tr>
            <th>No. Invoice</th>
            <th>Waktu & Tanggal</th>
            <th>Pelanggan</th>
            <th>Metode Bayar</th>
            <th>Total Bayar</th>
            <th>Status</th>
            <th class="text-right">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoading">
            <td colspan="7" class="text-center">Memuat riwayat transaksi...</td>
          </tr>
          <tr v-else-if="orders.length === 0">
            <td colspan="7" class="text-center">Belum ada transaksi recorded</td>
          </tr>
          <tr v-else v-for="order in orders" :key="order.id">
            <td>
              <code class="invoice-code">{{ order.invoice_no }}</code>
            </td>
            <td>{{ formatDate(order.created_at) }}</td>
            <td>{{ order.customer_name || 'Umum' }}</td>
            <td>
              <span class="pay-method-badge">
                {{ (order.payment_method || 'cash').toUpperCase() }}
              </span>
            </td>
            <td class="font-bold price-text">Rp {{ formatPrice(order.grand_total) }}</td>
            <td>
              <span class="badge badge-success">Selesai</span>
            </td>
            <td class="text-right">
              <button class="btn btn-secondary btn-sm" @click="openReceipt(order)">
                <span>🖨️ Struk</span>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Receipt Modal -->
    <ReceiptModal 
      v-if="selectedOrder"
      :order="selectedOrder"
      :store-setting="storeSetting"
      @close="selectedOrder = null"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import ReceiptModal from '../components/ReceiptModal.vue';

const props = defineProps({
  storeSetting: { type: Object, default: () => ({}) }
});

const orders = ref([]);
const isLoading = ref(true);
const selectedOrder = ref(null);

const formatPrice = (val) => new Intl.NumberFormat('id-ID').format(val || 0);

const formatDate = (dateStr) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleString('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  });
};

const fetchOrders = async () => {
  isLoading.value = true;
  try {
    const res = await fetch('/api/orders?limit=50');
    if (res.ok) orders.value = await res.json();
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchOrders);

const openReceipt = (order) => {
  selectedOrder.value = order;
};
</script>

<style scoped>
.orders-page {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
}

.header-title h2 {
  font-size: 1.25rem;
  font-weight: 800;
}

.table-container {
  overflow-x: auto;
  padding: 0.5rem;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.88rem;
}

.data-table th {
  padding: 0.85rem 1rem;
  font-weight: 700;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-color);
  background: rgba(15, 23, 42, 0.4);
}

.data-table td {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
}

.invoice-code {
  background: rgba(99, 102, 241, 0.15);
  color: #818cf8;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-weight: 700;
}

.pay-method-badge {
  background: rgba(255, 255, 255, 0.08);
  padding: 0.15rem 0.5rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 700;
}

.price-text {
  color: var(--accent-secondary);
}

.btn-sm {
  padding: 0.35rem 0.75rem;
  font-size: 0.8rem;
}

.text-right { text-align: right; }
.text-center { text-align: center; }
</style>
