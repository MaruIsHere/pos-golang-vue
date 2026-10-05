<template>
  <div class="flex flex-col gap-5 p-1">
    <div
      class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 py-5 px-6 bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700/60 rounded-[24px] shadow-sm"
    >
      <!-- <div class="flex flex-col gap-1"> -->
      <div class="">
        <h2 class="text-xl font-bold text-slate-800 dark:text-slate-100">Riwayat Transaksi & Retur Penjualan</h2>
        <p class="text-base text-slate-500 dark:text-slate-400">
          Daftar transaksi penjualan, cetak ulang struk, dan proses retur/refund barang
        </p>
      </div>

      <AppButton variant="primary" @click="fetchOrders">
        <ArrowPathIcon class="w-4 h-4" />
        <span>Refresh</span>
      </AppButton>
    </div>

    <!-- Orders Table -->
    <div class="overflow-x-auto bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700/60 rounded-[24px] shadow-sm">
      <table class="w-full text-left text-sm">
        <thead>
          <tr>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">No. Invoice</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Waktu & Tanggal</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Pelanggan</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Metode Bayar</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Total Bayar</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Status</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 text-right whitespace-nowrap">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoading">
            <td colspan="7" class="text-center p-4 text-slate-500">Memuat riwayat transaksi...</td>
          </tr>
          <tr v-else-if="orders.length === 0">
            <td colspan="7" class="text-center p-4 text-slate-500">Belum ada transaksi recorded</td>
          </tr>
          <tr v-else v-for="order in orders" :key="order.id">
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-800 dark:text-slate-100">
              <code class="px-2 py-1 bg-slate-100 dark:bg-slate-900 rounded-md text-xs font-mono font-bold text-slate-600 dark:text-slate-400">{{ order.invoice_no }}</code>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-800 dark:text-slate-100">{{ formatDate(order.created_at) }}</td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-800 dark:text-slate-100">{{ order.customer_name ||"Umum" }}</td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-800 dark:text-slate-100">
              <AppBadge variant="neutral">
                {{ (order.payment_method ||"cash").toUpperCase() }}
              </AppBadge>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle font-bold text-indigo-600 dark:text-indigo-400">Rp {{ formatPrice(order.grand_total) }}</td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-800 dark:text-slate-100">
              <span
                class="px-2 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider"
                :class="order.status === 'refunded' ? 'bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-400' : 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400'"
              >
                {{ order.status ==="refunded" ?"Diretur" :"Selesai" }}
              </span>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-right">
              <div class="flex items-center justify-end gap-2">
                <AppButton @click="openReceipt(order)" variant="secondary" size="sm">
                  <PrinterIcon class="w-3.5 h-3.5" />
                  <span>Struk</span>
                </AppButton>
                <AppButton
                  v-if="order.status !== 'refunded'"
                  @click="refundOrder(order)"
                  title="Retur Transaksi & Pulihkan Stok"
                  variant="danger"
                  size="sm"
                >
                  <ArrowUturnLeftIcon class="w-3.5 h-3.5" />
                  <span>Retur</span>
                </AppButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Receipt Modal -->
    <ReceiptModal v-if="selectedOrder" :order="selectedOrder" @close="selectedOrder = null" />
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppInput from '@/components/ui/AppInput.vue';

import { ref, onMounted } from"vue";
import api from '@/utils/api';
import ReceiptModal from"../components/ReceiptModal.vue";
import type { Order } from"../types";
import appButton from"../components/ui/AppButton.vue";
import { ArrowPathIcon, PrinterIcon, ArrowUturnLeftIcon } from"@heroicons/vue/24/outline";
import AppBadge from"@/components/ui/AppBadge.vue";

const orders = ref<Order[]>([]);
const isLoading = ref(true);
const selectedOrder = ref<Order | null>(null);

const formatPrice = (val: number): string => new Intl.NumberFormat("id-ID").format(val || 0);

const formatDate = (dateStr?: string): string => {
  if (!dateStr) return"-";
  return new Date(dateStr).toLocaleString("id-ID", {
    day:"2-digit",
    month:"short",
    year:"numeric",
    hour:"2-digit",
    minute:"2-digit",
  });
};

const fetchOrders = async () => {
  isLoading.value = true;
  try {
    const res = await api.get("/orders?limit=50");
    orders.value = res.data;
  } catch (err: any) {
    console.error(err.response?.data?.error || err.message || 'Error occurred');
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
    await api.post(`/orders/${order.id}/refund`);
    alert(`Transaksi #${order.invoice_no} berhasil diretur! Stok produk telah dipulihkan.`);
    fetchOrders();
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    alert("Koneksi error:" + errMsg);
  }
};
</script>


