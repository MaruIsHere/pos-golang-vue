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
    <div class="w-full min-w-0 overflow-x-auto bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700/60 rounded-[24px] shadow-sm">
      <table class="w-full min-w-[1200px] table-fixed text-left text-sm">
        <colgroup>
          <col class="w-[160px]" />
          <col class="w-[200px]" />
          <col class="w-[150px]" />
          <col class="w-[160px]" />
          <col class="w-[140px]" />
          <col class="w-[170px]" />
          <col class="w-[110px]" />
          <col class="w-[180px]" />
        </colgroup>
        <thead>
          <tr>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">No. Invoice</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Waktu & Tanggal</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Kasir</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Pelanggan</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Metode Bayar</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Total Bayar</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Status</th>
            <th class="p-3 text-right font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoading" class="border-b border-slate-200 dark:border-slate-700">
            <td colspan="8" class="px-3 py-8 text-center text-slate-500 dark:text-slate-400">Memuat riwayat transaksi...</td>
          </tr>
          <tr v-else-if="orders.length === 0" class="border-b border-slate-200 dark:border-slate-700">
            <td colspan="8" class="px-3 py-8 text-center text-slate-500 dark:text-slate-400">Belum ada transaksi.</td>
          </tr>
          <tr v-else v-for="order in orders" :key="order.id">
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-800 dark:text-slate-100">
              <code class="px-2 py-1 bg-slate-100 dark:bg-slate-900 rounded-md text-xs font-mono font-bold text-slate-600 dark:text-slate-400">{{ order.invoice_no }}</code>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-800 dark:text-slate-100">{{ formatDate(order.created_at) }}</td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle">
              <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-indigo-50 text-indigo-700 dark:bg-indigo-950/60 dark:text-indigo-300 border border-indigo-200/60 dark:border-indigo-800/60">
                <UserIcon class="w-3.5 h-3.5" />
                {{ order.cashier_name || 'Kasir' }}
              </span>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-800 dark:text-slate-100">{{ order.customer_name || "Umum" }}</td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-800 dark:text-slate-100">
              <div class="flex flex-col gap-1 items-start">
                <AppBadge variant="neutral">
                  {{ (order.payment_method || "cash").toUpperCase() }}
                </AppBadge>
                <button 
                  v-if="order.payment_proof" 
                  @click="viewProof(order)" 
                  class="inline-flex items-center gap-1 text-[11px] font-bold text-indigo-600 dark:text-indigo-400 bg-indigo-50 hover:bg-indigo-100 dark:bg-indigo-950/60 px-2 py-0.5 rounded-md border border-indigo-200 dark:border-indigo-800 cursor-pointer transition-colors"
                  title="Lihat Bukti Foto Pembayaran"
                >
                  <CameraIcon class="w-3.5 h-3.5 text-indigo-500" />
                  <span>Bukti Foto</span>
                </button>
              </div>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle font-bold text-indigo-600 dark:text-indigo-400 whitespace-nowrap">Rp {{ formatPrice(order.grand_total) }}</td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-800 dark:text-slate-100">
              <span
                class="px-2 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider"
                :class="order.status === 'refunded' ? 'bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-400' : 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400'"
              >
                {{ order.status === "refunded" ? "Diretur" : "Selesai" }}
              </span>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-right">
              <div class="flex items-center justify-end gap-2 whitespace-nowrap">
                <AppButton v-if="order.payment_proof" @click="viewProof(order)" variant="secondary" size="sm" title="Lihat Foto Bukti Bayar">
                  <CameraIcon class="w-3.5 h-3.5 text-indigo-500" />
                  <span>Bukti</span>
                </AppButton>
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

    <!-- Payment Proof Image Modal Lightbox -->
    <div v-if="selectedProofOrder" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm" @click.self="selectedProofOrder = null">
      <div class="w-full max-w-lg bg-white dark:bg-slate-900 rounded-2xl shadow-2xl flex flex-col overflow-hidden border border-slate-200 dark:border-slate-700">
        <div class="flex items-center justify-between p-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/50">
          <div>
            <h3 class="text-base font-bold text-slate-900 dark:text-slate-100 flex items-center gap-2">
              <CameraIcon class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
              <span>Bukti Pembayaran Transaksi</span>
            </h3>
            <p class="text-xs text-slate-500 dark:text-slate-400">Invoice: {{ selectedProofOrder.invoice_no }} ({{ (selectedProofOrder.payment_method || '').toUpperCase() }})</p>
          </div>
          <button class="p-1.5 rounded-lg text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors" @click="selectedProofOrder = null">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        
        <div class="p-5 flex flex-col items-center gap-4 max-h-[75vh] overflow-y-auto bg-slate-900">
          <img :src="selectedProofOrder.payment_proof" alt="Foto Bukti Pembayaran" class="max-w-full max-h-[60vh] object-contain rounded-xl shadow-lg border border-slate-700" />
        </div>

        <div class="p-4 border-t border-slate-200 dark:border-slate-800 flex justify-between items-center bg-white dark:bg-slate-900">
          <span class="text-xs text-slate-500 font-semibold">Nominal Bayar: <strong class="text-indigo-600 dark:text-indigo-400 font-extrabold text-sm">Rp {{ formatPrice(selectedProofOrder.grand_total) }}</strong></span>
          <div class="flex gap-2">
            <a :href="selectedProofOrder.payment_proof" target="_blank" class="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl text-xs font-bold transition-colors inline-flex items-center gap-1 shadow-sm">
              <ArrowDownTrayIcon class="w-3.5 h-3.5" /> Buka Foto Asli
            </a>
            <AppButton variant="secondary" size="sm" @click="selectedProofOrder = null">Tutup</AppButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppInput from '@/components/ui/AppInput.vue';

import { ref, onMounted } from "vue";
import api from '@/utils/api';
import ReceiptModal from "../components/ReceiptModal.vue";
import type { Order } from "../types";
import { ArrowPathIcon, PrinterIcon, ArrowUturnLeftIcon, UserIcon, CameraIcon, XMarkIcon, ArrowDownTrayIcon } from "@heroicons/vue/24/outline";
import AppBadge from "@/components/ui/AppBadge.vue";
import { showAppAlert, showAppConfirm } from '@/composables/useAppDialog';

const orders = ref<Order[]>([]);
const isLoading = ref(true);
const selectedOrder = ref<Order | null>(null);
const selectedProofOrder = ref<Order | null>(null);

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

const viewProof = (order: Order): void => {
  selectedProofOrder.value = order;
};

const refundOrder = async (order: Order): Promise<void> => {
  if (
    !await showAppConfirm(
      `Apakah Anda yakin ingin melakukan RETUR pada Invoice #${order.invoice_no}?\n\nStok barang akan dipulihkan secara otomatis.`,
      {
        title: 'Konfirmasi Retur',
        confirmLabel: 'Ya, Proses Retur',
        tone: 'warning'
      }
    )
  ) {
    return;
  }

  try {
    await api.post(`/orders/${order.id}/refund`);
    await showAppAlert(`Transaksi #${order.invoice_no} berhasil diretur! Stok produk telah dipulihkan.`, 'success');
    fetchOrders();
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    await showAppAlert("Koneksi error: " + errMsg, 'error');
  }
};
</script>

