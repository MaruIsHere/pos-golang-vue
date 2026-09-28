<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 print:static print:inset-auto print:bg-transparent print:p-0 print:flex-none print:items-start print:justify-start" @click.self="$emit('close')">
    <div class="w-full max-w-[420px] rounded-xl shadow-xl bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 flex flex-col max-h-[90vh] print:shadow-none print:bg-transparent print:text-black print:max-h-none print:max-w-full">
      <div class="flex items-center justify-between p-4 border-b border-slate-200 dark:border-slate-700 print:hidden">
        <h3 class="text-lg font-bold">Struk Belanja</h3>
        <button class="bg-transparent border-none cursor-pointer p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors" @click="$emit('close')">
          <XMarkIcon class="w-5 h-5 text-slate-500" />
        </button>
      </div>

      <!-- Thermal Receipt Preview -->
      <div class="flex-1 overflow-y-auto p-4 print:p-0 print:overflow-visible print:absolute print:left-0 print:top-0 print:w-full">
        <div class="bg-white text-slate-900 font-mono p-6 text-[13px] leading-relaxed shadow-sm border border-slate-200 rounded-lg mx-auto print:shadow-none print:border-none print:p-0 print:mx-0 print:w-full" id="receipt-print-area">
          <div class="text-center mb-2">
            <h2 class="text-lg font-black uppercase mb-1">{{ storeSetting.store_name || 'KASIR POS PRO' }}</h2>
            <p class="text-[12px] text-slate-600 print:text-black">{{ storeSetting.address || 'Jl. Utama No. 88, Jakarta' }}</p>
            <p class="text-[12px] text-slate-600 print:text-black">Telp: {{ storeSetting.phone || '0812-3456-7890' }}</p>
          </div>

          <div class="text-center text-slate-400 print:text-black my-2 text-[12px]">--------------------------------</div>

          <div class="flex flex-col gap-1">
            <div class="flex justify-between text-[12px]">
              <span>No. Struk:</span>
              <span class="font-bold">{{ order.invoice_no }}</span>
            </div>
            <div class="flex justify-between text-[12px]">
              <span>Tanggal:</span>
              <span>{{ formatDate(order.created_at) }}</span>
            </div>
            <div class="flex justify-between text-[12px]">
              <span>Kasir:</span>
              <span>{{ order.cashier_name || 'Kasir 1' }}</span>
            </div>
            <div class="flex justify-between text-[12px]">
              <span>Pelanggan:</span>
              <span>{{ order.customer_name || 'Umum' }}</span>
            </div>
          </div>

          <div class="text-center text-slate-400 print:text-black my-2 text-[12px]">--------------------------------</div>

          <!-- Items Table -->
          <div class="flex flex-col gap-2">
            <div v-for="item in order.order_items" :key="item.id">
              <div class="font-bold">{{ item.product_name }}</div>
              <div class="flex justify-between text-[12px] text-slate-700 print:text-black">
                <span>{{ item.quantity }} x {{ formatPrice(item.product_price) }}</span>
                <span class="font-bold">Rp {{ formatPrice(item.subtotal) }}</span>
              </div>
            </div>
          </div>

          <div class="text-center text-slate-400 print:text-black my-2 text-[12px]">--------------------------------</div>

          <!-- Totals -->
          <div class="flex flex-col gap-1">
            <div class="flex justify-between text-[12.5px]">
              <span>Subtotal:</span>
              <span>Rp {{ formatPrice(order.total_amount) }}</span>
            </div>
            <div v-if="order.discount > 0" class="flex justify-between text-[12.5px]">
              <span>Diskon:</span>
              <span>- Rp {{ formatPrice(order.discount) }}</span>
            </div>
            <div class="flex justify-between text-[12.5px]">
              <span>Pajak:</span>
              <span>Rp {{ formatPrice(order.tax) }}</span>
            </div>
            <div class="flex justify-between text-[15px] font-black my-1 py-1 border-y border-dashed border-slate-900 print:border-black">
              <span>TOTAL:</span>
              <span>Rp {{ formatPrice(order.grand_total) }}</span>
            </div>
            <div class="flex justify-between text-[12.5px]">
              <span>Bayar ({{ (order.payment_method || 'cash').toUpperCase() }}):</span>
              <span>Rp {{ formatPrice(order.paid_amount) }}</span>
            </div>
            <div class="flex justify-between text-[12.5px]">
              <span>Kembali:</span>
              <span>Rp {{ formatPrice(order.change_amount) }}</span>
            </div>
          </div>

          <div class="text-center text-slate-400 print:text-black my-2 text-[12px]">--------------------------------</div>

          <div class="text-center mt-2">
            <p class="text-[12px] font-semibold whitespace-pre-line">{{ storeSetting.receipt_footer || 'Terima kasih telah berbelanja!' }}</p>
            <p class="text-[10.5px] text-slate-500 print:text-black mt-2">Powered by POS Golang + Vue.js</p>
          </div>
        </div>
      </div>

      <div class="p-4 border-t border-slate-200 dark:border-slate-700 flex justify-end gap-3 print:hidden">
        <AppButton variant="secondary" @click="$emit('close')">Tutup</AppButton>
        <AppButton variant="primary" class="flex items-center gap-1.5" @click="printReceipt">
          <PrinterIcon class="w-4 h-4" />
          <span>Cetak Struk</span>
        </AppButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useSettingsStore } from '../stores/settings';
import type { Order } from '../types';
import { PrinterIcon, XMarkIcon } from '@heroicons/vue/24/outline';
import AppButton from './ui/AppButton.vue';

defineProps({
  order: { type: Object as () => Order, required: true }
});

defineEmits(['close']);

const settingsStore = useSettingsStore();
const { settings: storeSetting } = storeToRefs(settingsStore);

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

const formatDate = (dateStr?: string): string => {
  if (!dateStr) return new Date().toLocaleString('id-ID');
  return new Date(dateStr).toLocaleString('id-ID');
};

const printReceipt = () => {
  window.print();
};
</script>
