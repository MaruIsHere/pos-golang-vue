<template>
  <div class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 p-4 print:static print:inset-auto print:bg-transparent print:p-0 print:flex-none print:items-start print:justify-start" @click.self="$emit('close')">
    <div :class="['w-full rounded-xl shadow-xl bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 flex flex-col max-h-[90vh] print:shadow-none print:bg-transparent print:text-black print:max-h-none print:max-w-full transition-all duration-200', paperSize === '58mm' ? 'max-w-[360px]' : 'max-w-[440px]']">
      <div class="flex items-center justify-between p-3.5 border-b border-slate-200 dark:border-slate-700 print:hidden">
        <div class="flex items-center gap-2">
          <h3 class="text-base font-bold">Struk</h3>
          <!-- Ukuran Kertas Selector -->
          <div class="flex items-center bg-slate-100 dark:bg-slate-700/70 p-1 rounded-lg text-xs font-medium">
            <button 
              type="button" 
              class="px-2 py-0.5 rounded transition-all"
              :class="paperSize === '58mm' ? 'bg-indigo-600 text-white font-semibold shadow-xs' : 'text-slate-600 dark:text-slate-300 hover:text-slate-900'"
              @click="paperSize = '58mm'"
            >
              58mm
            </button>
            <button 
              type="button" 
              class="px-2 py-0.5 rounded transition-all"
              :class="paperSize === '80mm' ? 'bg-indigo-600 text-white font-semibold shadow-xs' : 'text-slate-600 dark:text-slate-300 hover:text-slate-900'"
              @click="paperSize = '80mm'"
            >
              80mm
            </button>
          </div>
        </div>
        <button class="bg-transparent border-none cursor-pointer p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors" @click="$emit('close')">
          <XMarkIcon class="w-5 h-5 text-slate-500" />
        </button>
      </div>

      <!-- Thermal Receipt Preview -->
      <div class="flex-1 overflow-y-auto p-4 print:p-0 print:overflow-visible print:absolute print:left-0 print:top-0 print:w-full">
        <div 
          id="receipt-print-area"
          :class="[
            'bg-white text-slate-900 font-mono p-5 text-[12.5px] leading-relaxed shadow-sm border border-slate-200 rounded-lg mx-auto print:shadow-none print:border-none print:p-0 print:mx-0 print:w-full',
            paperSize === '58mm' ? 'paper-58mm max-w-[320px]' : 'paper-80mm max-w-[400px]'
          ]"
        >
          <div class="text-center mb-2">
            <h2 class="text-base font-black uppercase mb-1 leading-snug">{{ storeSetting.store_name || 'KASIR POS PRO' }}</h2>
            <p class="text-[11.5px] text-slate-600 print:text-black leading-snug">{{ storeSetting.address || 'Jl. Utama No. 88, Jakarta' }}</p>
            <p class="text-[11.5px] text-slate-600 print:text-black">Telp: {{ storeSetting.phone || '0812-3456-7890' }}</p>
          </div>

          <div class="text-center text-slate-400 print:text-black my-1.5 text-[11px] overflow-hidden whitespace-nowrap">------------------------------------------</div>

          <div class="flex flex-col gap-0.5">
            <div class="flex justify-between text-[11.5px]">
              <span>No. Struk:</span>
              <span class="font-bold">{{ order.invoice_no }}</span>
            </div>
            <div class="flex justify-between text-[11.5px]">
              <span>Tanggal:</span>
              <span>{{ formatDate(order.created_at) }}</span>
            </div>
            <div class="flex justify-between text-[11.5px]">
              <span>Kasir:</span>
              <span>{{ order.cashier_name || 'Kasir 1' }}</span>
            </div>
            <div class="flex justify-between text-[11.5px]">
              <span>Pelanggan:</span>
              <span>{{ order.customer_name || 'Umum' }}</span>
            </div>
            <div v-if="order.table_number" class="flex justify-between text-[11.5px]">
              <span>Meja:</span>
              <span>{{ order.table_number }}</span>
            </div>
          </div>

          <div class="text-center text-slate-400 print:text-black my-1.5 text-[11px] overflow-hidden whitespace-nowrap">------------------------------------------</div>

          <!-- Items Table -->
          <div class="flex flex-col gap-1.5">
            <div v-for="item in order.order_items" :key="item.id">
              <div class="font-bold text-[12px] leading-tight">{{ item.product_name }}</div>
              <div class="flex justify-between text-[11.5px] text-slate-700 print:text-black">
                <span>{{ formatQuantity(item.quantity) }} {{ unitLabel(item.unit) }} x {{ formatPrice(item.product_price) }}</span>
                <span class="font-bold">Rp {{ formatPrice(item.subtotal) }}</span>
              </div>
            </div>
          </div>

          <div class="text-center text-slate-400 print:text-black my-1.5 text-[11px] overflow-hidden whitespace-nowrap">------------------------------------------</div>

          <!-- Totals -->
          <div class="flex flex-col gap-0.5">
            <div class="flex justify-between text-[11.5px]">
              <span>Subtotal:</span>
              <span>Rp {{ formatPrice(order.total_amount) }}</span>
            </div>
            <div v-if="order.discount > 0" class="flex justify-between text-[11.5px]">
              <span>Diskon:</span>
              <span>- Rp {{ formatPrice(order.discount) }}</span>
            </div>
            <div class="flex justify-between text-[11.5px]">
              <span>Pajak:</span>
              <span>Rp {{ formatPrice(order.tax) }}</span>
            </div>
            <div class="flex justify-between text-[14px] font-black my-1 py-1 border-y border-dashed border-slate-900 print:border-black">
              <span>TOTAL:</span>
              <span>Rp {{ formatPrice(order.grand_total) }}</span>
            </div>
            <div class="flex justify-between text-[11.5px]">
              <span>Bayar ({{ (order.payment_method || 'cash').toUpperCase() }}):</span>
              <span>Rp {{ formatPrice(order.paid_amount) }}</span>
            </div>
            <div class="flex justify-between text-[11.5px]">
              <span>Kembali:</span>
              <span>Rp {{ formatPrice(order.change_amount) }}</span>
            </div>
            <div v-if="order.payment_proof" class="flex justify-between text-[11.5px] font-bold text-indigo-700 print:text-black mt-1 pt-1 border-t border-dashed border-slate-300 print:border-slate-400">
              <span>Bukti Pembayaran:</span>
              <span>Terlampir (Foto)</span>
            </div>
          </div>

          <!-- Proof Preview in Thermal Box (Print Hidden) -->
          <div v-if="order.payment_proof" class="mt-3 pt-2 border-t border-slate-200 text-center print:hidden">
            <span class="text-[10px] font-bold uppercase text-slate-500 block mb-1">Foto Bukti Transfer/QRIS:</span>
            <img :src="order.payment_proof" alt="Foto Bukti Pembayaran" class="w-full max-h-36 object-contain rounded border border-slate-200 mx-auto" />
          </div>

          <div class="text-center text-slate-400 print:text-black my-1.5 text-[11px] overflow-hidden whitespace-nowrap">------------------------------------------</div>

          <div class="text-center mt-2">
            <p class="text-[11.5px] font-semibold whitespace-pre-line leading-tight">{{ storeSetting.receipt_footer || 'Terima kasih telah berbelanja!' }}</p>
            <p class="text-[10px] text-slate-500 print:text-black mt-1.5">Powered by POS Golang + Vue.js</p>
          </div>
        </div>
      </div>

      <div class="p-3.5 border-t border-slate-200 dark:border-slate-700 flex items-center justify-between gap-3 print:hidden">
        <span class="text-xs text-slate-500 dark:text-slate-400">Kertas: <strong>{{ paperSize }}</strong></span>
        <div class="flex gap-2">
          <AppButton variant="secondary" @click="$emit('close')">Tutup</AppButton>
          <AppButton variant="primary" class="flex items-center gap-1.5" @click="printReceipt">
            <PrinterIcon class="w-4 h-4" />
            <span>Cetak Struk</span>
          </AppButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { storeToRefs } from 'pinia';
import { useSettingsStore } from '../stores/settings';
import type { Order } from '../types';
import { PrinterIcon, XMarkIcon } from '@heroicons/vue/24/outline';
import AppButton from './ui/AppButton.vue';

defineProps({
  order: { type: Object as () => Order, required: true }
});

defineEmits(['close']);

const paperSize = ref<'58mm' | '80mm'>('58mm');

const settingsStore = useSettingsStore();
const { settings: storeSetting } = storeToRefs(settingsStore);

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);
const formatQuantity = (val: number): string => new Intl.NumberFormat('id-ID', { maximumFractionDigits: 3 }).format(val || 0);
const unitLabel = (unit: string): string => unit === 'gram' ? 'gr' : unit === 'liter' ? 'L' : 'pcs';

const formatDate = (dateStr?: string): string => {
  if (!dateStr) return new Date().toLocaleString('id-ID');
  return new Date(dateStr).toLocaleString('id-ID');
};

const printReceipt = () => {
  // Dynamically set @page size for thermal receipt printing (58mm / 80mm roll)
  let styleEl = document.getElementById('dynamic-print-style');
  if (!styleEl) {
    styleEl = document.createElement('style');
    styleEl.id = 'dynamic-print-style';
    document.head.appendChild(styleEl);
  }
  styleEl.innerHTML = `@page { size: ${paperSize.value} auto; margin: 0mm; }`;

  window.print();
};
</script>
