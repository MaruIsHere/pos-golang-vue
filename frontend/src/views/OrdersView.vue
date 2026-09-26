<template>
  <div class="flex flex-col gap-5 p-1">
    <div
      class="flex justify-between items-center py-5 px-6 bg-bg-card border-border-color rounded-xl"
    >
      <!-- <div class="header-title"> -->
      <div class="">
        <h2 class="text-xl font-bold text-text-primary">
          Riwayat Transaksi & Retur Penjualan
        </h2>
        <p class="text-base text-text-secondary">
          Daftar transaksi penjualan, cetak ulang struk, dan proses retur/refund
          barang
        </p>
      </div>

      <appButton variant="primary" @click="fetchOrders">
        <ArrowPathIcon class="w-4 h-4" />
        <span>Refresh</span>
      </appButton>
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
            <td colspan="7" class="text-center">
              Belum ada transaksi recorded
            </td>
          </tr>
          <tr v-else v-for="order in orders" :key="order.id">
            <td>
              <code class="invoice-code">{{ order.invoice_no }}</code>
            </td>
            <td>{{ formatDate(order.created_at) }}</td>
            <td>{{ order.customer_name || "Umum" }}</td>
            <td>
              <span class="pay-method-badge">
                {{ (order.payment_method || "cash").toUpperCase() }}
              </span>
            </td>
            <td class="font-bold price-text">
              Rp {{ formatPrice(order.grand_total) }}
            </td>
            <td>
              <span
                class="badge"
                :class="
                  order.status === 'refunded'
                    ? 'badge-refunded'
                    : 'badge-success'
                "
              >
                {{ order.status === "refunded" ? "Diretur" : "Selesai" }}
              </span>
            </td>
            <td class="text-right">
              <div class="action-flex">
                <appButton
                  @click="openReceipt(order)"
                  variant="secondary"
                  size="sm"
                >
                  <PrinterIcon class="w-3.5 h-3.5" />
                  <span>Struk</span>
                </appButton>
                <appButton
                  v-if="order.status !== 'refunded'"
                  @click="refundOrder(order)"
                  title="Retur Transaksi & Pulihkan Stok"
                  variant="danger"
                  size="sm"
                >
                  <ArrowUturnLeftIcon class="w-3.5 h-3.5" />
                  <span>Retur</span>
                </appButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Receipt Modal -->
    <ReceiptModal
      v-if="selectedOrder"
      :order="selectedOrder"
      @close="selectedOrder = null"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import ReceiptModal from "../components/ReceiptModal.vue";
import type { Order } from "../types";
import appButton from "../components/ui/Appbutton.vue";
import {
  ArrowPathIcon,
  PrinterIcon,
  ArrowUturnLeftIcon,
} from "@heroicons/vue/24/outline";

const orders = ref<Order[]>([]);
const isLoading = ref(true);
const selectedOrder = ref<Order | null>(null);

const formatPrice = (val: number): string =>
  new Intl.NumberFormat("id-ID").format(val || 0);

const formatDate = (dateStr?: string): string => {
  if (!dateStr) return "-";
  return new Date(dateStr).toLocaleString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const fetchOrders = async () => {
  isLoading.value = true;
  try {
    const res = await fetch("/api/orders?limit=50");
    if (res.ok) orders.value = await res.json();
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchOrders);

const openReceipt = (order: Order): void => {
  selectedOrder.value = order;
};

const refundOrder = async (order: Order): Promise<void> => {
  if (
    !confirm(
      `Apakah Anda yakin ingin melakukan RETUR pada Invoice #${order.invoice_no}?\n\nStok barang akan dipulihkan secara otomatis.`,
    )
  ) {
    return;
  }

  try {
    const res = await fetch(`/api/orders/${order.id}/refund`, {
      method: "POST",
    });

    if (res.ok) {
      alert(
        `Transaksi #${order.invoice_no} berhasil diretur! Stok produk telah dipulihkan.`,
      );
      fetchOrders();
    } else {
      const errData = await res.json();
      alert(
        "Gagal merefur transaksi: " + (errData.error || "Terjadi kesalahan"),
      );
    }
  } catch (err) {
    alert("Koneksi error: " + (err as Error).message);
  }
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

.table-container {
  overflow-x: auto;
  padding: 0.5rem;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.875rem;
}

.data-table th {
  padding: 0.85rem 1rem;
  font-weight: 600;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-primary);
}

.data-table td {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
  color: var(--text-primary);
}

.invoice-code {
  background: rgba(99, 102, 241, 0.12);
  color: var(--accent-primary);
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-weight: 700;
  font-size: 0.8rem;
  border: 1px solid rgba(99, 102, 241, 0.3);
}

.pay-method-badge {
  background: var(--bg-primary);
  padding: 0.15rem 0.5rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.price-text {
  color: #059669;
}

.badge {
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
}

.badge-success {
  background: rgba(16, 185, 129, 0.15);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.badge-refunded {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.action-flex {
  display: flex;
  justify-content: flex-end;
  gap: 0.4rem;
}

.btn-danger-outline {
  background: rgba(239, 68, 68, 0.12);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
}

.btn-danger-outline:hover {
  background: rgba(239, 68, 68, 0.25);
}

.btn-sm {
  padding: 0.35rem 0.65rem;
  font-size: 0.78rem;
}

.text-right {
  text-align: right;
}
.text-center {
  text-align: center;
}
</style>
