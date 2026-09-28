<template>
  <!-- Mobile Backdrop -->
  <div 
    v-if="isOpenMobile" 
    class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[80] transition-opacity" 
    @click="$emit('toggle-mobile')"
  ></div>

  <aside 
    class="flex flex-col h-full w-full rounded-[24px] lg:rounded-none lg:border-none border border-slate-200 dark:border-slate-700 overflow-hidden bg-white dark:bg-slate-800 max-md:fixed max-md:bottom-[80px] max-md:left-4 max-md:right-4 max-md:w-auto max-md:h-[75vh] max-md:z-[90] max-md:shadow-2xl max-md:translate-y-[120%] max-md:transition-transform max-md:duration-300 max-md:ease-out" 
    :class="{ 'max-md:translate-y-0': isOpenMobile }"
  >
    <!-- Header -->
    <div class="flex items-center justify-between p-5 border-b border-slate-100 dark:border-slate-700/50 bg-white dark:bg-slate-900 shrink-0 shadow-sm z-10">
      <div class="flex items-center gap-3">
        <div class="p-2 bg-indigo-100 dark:bg-indigo-900/40 rounded-xl text-indigo-600 dark:text-indigo-400">
          <ShoppingCartIcon class="w-5 h-5" />
        </div>
        <div class="flex flex-col">
          <h2 class="text-base font-extrabold text-slate-800 dark:text-slate-100 leading-tight">Pesanan Aktif</h2>
          <span class="text-xs font-semibold text-slate-500">{{ totalItemCount }} Produk Terpilih</span>
        </div>
      </div>
      <button v-if="cart.length > 0" class="text-xs font-bold text-red-500 hover:text-red-700 bg-red-50 hover:bg-red-100 dark:bg-red-900/20 dark:hover:bg-red-900/40 px-3 py-2 rounded-xl transition-colors" @click="handleClearCart">
        Kosongkan
      </button>
    </div>

    <!-- Empty State -->
    <div v-if="cart.length === 0" class="flex-1 flex flex-col items-center justify-center p-8 text-center bg-slate-50/50 dark:bg-slate-900/10">
      <div class="p-5 bg-slate-100 dark:bg-slate-800 rounded-full mb-4">
        <ShoppingBagIcon class="w-12 h-12 text-slate-300 dark:text-slate-600" />
      </div>
      <p class="font-bold text-slate-700 dark:text-slate-300 mb-1">Keranjang masih kosong</p>
      <span class="text-xs font-medium text-slate-400 max-w-[200px]">Silakan pilih produk dari katalog di sebelah kiri.</span>
    </div>

    <!-- Cart Items & Voucher List (SCROLLABLE AREA) -->
    <div v-else class="flex-1 overflow-y-auto p-4 flex flex-col gap-4 custom-scrollbar bg-slate-50/50 dark:bg-slate-900/10">
      
      <!-- Items List -->
      <div class="flex flex-col gap-3">
        <div v-for="(item, index) in cart" :key="item.product.id" class="bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700/60 shadow-sm rounded-[20px] p-3 flex flex-col gap-3 relative group transition-all hover:border-indigo-300 dark:hover:border-indigo-700">
          

          <div class="flex flex-col pr-8">
            <h4 class="text-[0.9rem] font-bold text-slate-800 dark:text-slate-100 leading-snug line-clamp-2">{{ item.product.name }}</h4>
            <span class="text-[0.7rem] font-bold text-slate-400 mt-1 uppercase tracking-wider">@ Rp {{ formatPrice(item.product.price) }}</span>
          </div>

          <div class="flex items-center justify-between mt-1 pt-3 border-t border-slate-100 dark:border-slate-700/50">
            <div class="flex items-center bg-slate-100 dark:bg-slate-900 rounded-xl p-1 border border-slate-200/60 dark:border-slate-700/60">
              <button class="w-8 h-8 rounded-[10px] bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 font-bold flex items-center justify-center hover:bg-slate-50 shadow-sm active:scale-95 transition-all" @click="$emit('update-qty', { index, delta: -1 })">
                <MinusIcon class="w-4 h-4" />
              </button>
              <span class="text-sm font-extrabold min-w-[32px] text-center text-slate-800 dark:text-slate-200">{{ item.quantity }}</span>
              <button class="w-8 h-8 rounded-[10px] bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300 font-bold flex items-center justify-center hover:bg-slate-50 shadow-sm active:scale-95 transition-all" @click="$emit('update-qty', { index, delta: 1 })">
                <PlusIcon class="w-4 h-4" />
              </button>
            </div>
            <div class="flex items-center gap-3">
            <span class="text-[0.95rem] font-black text-indigo-600 dark:text-indigo-400">Rp {{ formatPrice(item.product.price * item.quantity) }}</span>
            <button class="w-8 h-8 flex items-center justify-center bg-red-50 text-red-500 hover:bg-red-100 hover:text-red-600 dark:bg-red-900/30 dark:text-red-400 dark:hover:bg-red-900/50 rounded-[10px] transition-colors active:scale-95 shadow-sm border border-red-100 dark:border-red-800/30 shrink-0" @click="$emit('remove-item', index)">
              <TrashIcon class="w-4 h-4" />
            </button>
          </div>
          </div>
        </div>
      </div>

      <!-- Voucher Section Moved INSIDE Scrollable Area -->
      <div class="bg-white dark:bg-slate-800 border border-dashed border-slate-300 dark:border-slate-700 rounded-[20px] p-4 mt-2">
        <span class="text-[0.7rem] font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2 block">Kode Voucher / Diskon</span>
        <div class="relative flex items-center">
          <TicketIcon class="w-5 h-5 text-slate-400 absolute left-3" />
          <input 
            type="text" 
            v-model="voucherCode" 
            placeholder="Masukkan kode voucher" 
            class="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl py-2.5 pl-10 pr-3 text-sm font-bold text-slate-800 dark:text-slate-100 placeholder:text-slate-400 placeholder:font-normal focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all uppercase"
          />
        </div>
        
        <div v-if="availableVouchersHint.length > 0" class="flex flex-wrap items-center gap-2 mt-3">
          <button 
            v-for="v in availableVouchersHint" 
            :key="v.code" 
            class="px-3 py-1.5 rounded-lg text-xs font-bold border border-indigo-100 dark:border-indigo-900/50 text-indigo-600 dark:text-indigo-400 bg-indigo-50/50 dark:bg-indigo-900/20 hover:bg-indigo-600 hover:text-white transition-colors"
            @click="useHint(v.code)"
          >
            {{ v.code }}
          </button>
        </div>
      </div>
      
      <!-- Spacer at bottom of scroll -->
      <div class="h-4"></div>
    </div>

    <!-- Super Compact Bottom Checkout Area (Sticky) -->
    <div v-if="cart.length > 0" class="border-t border-slate-100 dark:border-slate-700/50 bg-white dark:bg-slate-800 p-4 shrink-0 shadow-[0_-4px_20px_rgba(0,0,0,0.03)] z-10">
      
      <!-- Compact Breakdown -->
      <div class="flex justify-between items-center text-[0.75rem] font-bold text-slate-500 dark:text-slate-400 mb-3 px-1">
        <span>Sub: Rp {{ formatPrice(subtotal) }}</span>
        <span>Pjk: Rp {{ formatPrice(taxAmount) }}</span>
        <span v-if="discount > 0" class="text-rose-500">Disc: -Rp {{ formatPrice(discount) }}</span>
      </div>

      <!-- Giant Checkout Button with Price Inside -->
      <button 
        class="w-full py-4 px-5 rounded-[18px] shadow-[0_8px_20px_-6px_rgba(79,70,229,0.4)] flex items-center justify-between transition-all bg-indigo-600 text-white hover:bg-indigo-700 hover:scale-[1.02] active:scale-[0.98]"
        @click="$emit('open-payment')"
      >
        <div class="flex items-center gap-2">
          <span class="text-sm font-bold uppercase tracking-wider text-indigo-100">Proses Bayar</span>
          <ArrowRightIcon class="w-4 h-4 text-indigo-200" />
        </div>
        <span class="text-xl font-black tracking-wide">Rp {{ formatPrice(grandTotal) }}</span>
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import api from '@/utils/api';
import type { CartItem, Voucher } from '../types';
import AppButton from './ui/AppButton.vue';
import AppInput from './ui/AppInput.vue';
import AppBadge from './ui/AppBadge.vue';
import {
  ShoppingCartIcon,
  ShoppingBagIcon,
  TagIcon,
  TicketIcon,
  ExclamationTriangleIcon,
  ArrowRightIcon,
  XMarkIcon,
  PlusIcon,
  MinusIcon,
  TrashIcon
} from '@heroicons/vue/24/outline';

interface AppliedVoucher {
  code: string;
  discountAmount: number;
  description: string;
}

const props = defineProps({
  cart: { type: Array as () => CartItem[], required: true },
  discount: { type: Number, default: 0 },
  taxPercentage: { type: Number, default: 10 },
  isOpenMobile: { type: Boolean, default: false }
});

const emit = defineEmits(['update-qty', 'remove-item', 'clear-cart', 'update-discount', 'open-payment', 'toggle-mobile']);

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

const totalItemCount = computed(() => {
  return props.cart.reduce((sum, item) => sum + item.quantity, 0);
});

const subtotal = computed(() => {
  return props.cart.reduce((sum, item) => sum + (item.product.price * item.quantity), 0);
});

const taxAmount = computed(() => {
  const taxable = Math.max(0, subtotal.value - props.discount);
  return (taxable * props.taxPercentage) / 100;
});

const grandTotal = computed(() => {
  return Math.max(0, subtotal.value - props.discount) + taxAmount.value;
});

// Voucher Logic
const voucherInput = ref('');
const appliedVoucher = ref<AppliedVoucher | null>(null);
const voucherError = ref('');

const availableVouchers = ref<Voucher[]>([
  { code: 'DISKON10', type: 'percent', value: 10, description: 'Diskon 10%' },
  { code: 'DISKON20', type: 'percent', value: 20, description: 'Diskon 20%' },
  { code: 'HEMAT10K', type: 'flat', value: 10000, description: 'Potongan Rp 10.000' },
  { code: 'HEMAT50K', type: 'flat', value: 50000, description: 'Potongan Rp 50.000' },
  { code: 'POSHEMAT', type: 'percent', value: 15, description: 'Diskon POS 15%' }
]);

const fetchVouchers = async () => {
  try {
    const res = await api.get('/vouchers');
    const data = res.data;
    if (Array.isArray(data) && data.length > 0) {
      availableVouchers.value = data;
    }
  } catch (err: any) {
    console.error('Fetch vouchers error:', err.response?.data?.error || err.message || 'Error occurred');
  }
};

onMounted(fetchVouchers);

const availableVouchersHint = computed(() => {
  return availableVouchers.value.slice(0, 4);
});

const applyVoucher = async () => {
  voucherError.value = '';
  const code = voucherInput.value.trim().toUpperCase();
  if (!code) {
    voucherError.value = 'Silakan ketik kode voucher terlebih dahulu';
    return;
  }

  // Refresh latest vouchers from server
  await fetchVouchers();

  const found = availableVouchers.value.find(v => v.code.toUpperCase() === code);
  if (!found) {
    voucherError.value = `Kode voucher '${code}' tidak valid / tidak ditemukan`;
    return;
  }

  let amount = 0;
  if (found.type === 'percent') {
    amount = Math.round((subtotal.value * found.value) / 100);
  } else {
    amount = found.value;
  }

  if (amount > subtotal.value) {
    amount = subtotal.value;
  }

  appliedVoucher.value = {
    code: found.code,
    discountAmount: amount,
    description: found.description ?? ''
  };
  voucherInput.value = '';
  
  emit('update-discount', amount);
};

const removeVoucher = () => {
  appliedVoucher.value = null;
  voucherInput.value = '';
  voucherError.value = '';
  emit('update-discount', 0);
};

const useHint = (code: string): void => {
  voucherInput.value = code;
  applyVoucher();
};

const onManualDiscountInput = (val: string): void => {
  appliedVoucher.value = null;
  voucherError.value = '';
  emit('update-discount', parseFloat(val) || 0);
};

const handleClearCart = () => {
  removeVoucher();
  emit('clear-cart');
};
</script>
