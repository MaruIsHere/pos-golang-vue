<template>
  <!-- Mobile Backdrop -->
  <div 
    v-if="isOpenMobile" 
    class="fixed top-0 left-0 right-0 bottom-0 bg-slate-900/40 z-[80]" 
    @click="$emit('toggle-mobile')"
  ></div>

  <aside 
    class="flex flex-col h-full w-full rounded-lg border border-slate-200 dark:border-slate-700 overflow-hidden shadow-sm backdrop-blur-md bg-white/90 dark:bg-slate-900/90 max-md:fixed max-md:bottom-[60px] max-md:left-0 max-md:right-0 max-md:h-[75vh] max-md:z-[90] max-md:rounded-t-xl max-md:rounded-b-none max-md:translate-y-[105%] max-md:transition-transform max-md:duration-300 max-md:ease-[cubic-bezier(0.16,1,0.3,1)]" 
    :class="{ 'max-md:translate-y-0': isOpenMobile }"
  >
    <div class="flex items-center justify-between p-4 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900">
      <div class="flex items-center gap-2">
        <ShoppingCartIcon class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
        <h2 class="text-[1.05rem] font-bold text-slate-900 dark:text-slate-100">Keranjang</h2>
        <AppBadge variant="primary" class="text-xs font-bold">{{ totalItemCount }} Item</AppBadge>
      </div>
      <AppButton v-if="cart.length > 0" variant="ghost" class="text-red-500 text-[0.8rem] font-semibold hover:text-red-600 p-0" @click="handleClearCart">
        Hapus Semua
      </AppButton>
    </div>

    <!-- Empty State -->
    <div v-if="cart.length === 0" class="flex-1 flex flex-col items-center justify-center p-8 text-center">
      <div class="mb-3">
        <ShoppingBagIcon class="w-12 h-12 text-slate-300 dark:text-slate-600" />
      </div>
      <p class="font-bold text-slate-900 dark:text-slate-100 mb-1">Keranjang masih kosong</p>
      <span class="text-[0.8rem] text-slate-500">Pilih produk di sebelah kiri untuk ditambahkan</span>
    </div>

    <!-- Cart Items List -->
    <div v-else class="flex-1 overflow-y-auto p-4 flex flex-col gap-3">
      <div v-for="(item, index) in cart" :key="item.product.id" class="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-[10px] p-3 flex flex-col gap-2">
        <div class="flex justify-between items-start">
          <h4 class="text-[0.85rem] font-bold text-slate-900 dark:text-slate-100">{{ item.product.name }}</h4>
          <span class="text-[0.8rem] text-slate-500">Rp {{ formatPrice(item.product.price) }}</span>
        </div>

        <div class="flex items-center justify-between">
          <div class="flex items-center gap-1.5 bg-slate-50 dark:bg-slate-900 rounded-lg p-0.5 border border-slate-200 dark:border-slate-700">
            <AppButton variant="ghost" class="w-[26px] h-[26px] rounded-md border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 font-bold p-0 flex items-center justify-center hover:bg-slate-50 dark:hover:bg-slate-700" @click="$emit('update-qty', { index, delta: -1 })">
              <MinusIcon class="w-3.5 h-3.5" />
            </AppButton>
            <span class="text-[0.85rem] font-extrabold min-w-[20px] text-center text-slate-900 dark:text-slate-100">{{ item.quantity }}</span>
            <AppButton variant="ghost" class="w-[26px] h-[26px] rounded-md border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-900 dark:text-slate-100 font-bold p-0 flex items-center justify-center hover:bg-slate-50 dark:hover:bg-slate-700" @click="$emit('update-qty', { index, delta: 1 })">
              <PlusIcon class="w-3.5 h-3.5" />
            </AppButton>
          </div>
          <span class="text-[0.85rem] font-extrabold text-indigo-600 dark:text-indigo-400">Rp {{ formatPrice(item.product.price * item.quantity) }}</span>
          <AppButton variant="ghost" class="p-1 flex items-center justify-center" @click="$emit('remove-item', index)">
            <TrashIcon class="w-4 h-4 text-slate-400 hover:text-red-500" />
          </AppButton>
        </div>
      </div>
    </div>

    <!-- Summary & Checkout Footer -->
    <div v-if="cart.length > 0" class="p-4 border-t border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 flex flex-col gap-2.5">
      <!-- Voucher Code Input Section -->
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-1.5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg p-1 px-2">
          <TagIcon class="w-4 h-4 text-slate-400" />
          <AppInput 
            type="text" 
            class="flex-1 text-[0.8rem] uppercase bg-transparent border-none outline-none text-slate-900 dark:text-slate-100" 
            v-model="voucherInput" 
            placeholder="Kode Voucher (DISKON10...)"
            @keyup.enter="applyVoucher"
          />
          <AppButton variant="primary" class="py-1.5 px-3 text-[0.75rem] font-bold rounded-md" @click="applyVoucher">
            Gunakan
          </AppButton>
        </div>

        <!-- Applied Voucher Alert Badge -->
        <AppBadge v-if="appliedVoucher" variant="success" class="flex items-center justify-between p-2 rounded-lg text-[0.75rem] font-bold bg-emerald-500/15 border border-emerald-500/30 text-emerald-500">
          <div class="flex flex-col">
            <span class="font-extrabold flex items-center gap-1">
              <TicketIcon class="w-4 h-4 text-emerald-600" /> {{ appliedVoucher.code }}
            </span>
            <span class="text-[0.7rem] opacity-90">Potongan Rp {{ formatPrice(appliedVoucher.discountAmount) }} ({{ appliedVoucher.description }})</span>
          </div>
          <AppButton variant="ghost" class="p-1" title="Hapus Voucher" @click="removeVoucher">
            <XMarkIcon class="w-4 h-4" />
          </AppButton>
        </AppBadge>

        <!-- Error Message -->
        <AppBadge v-if="voucherError" variant="danger" class="flex items-center gap-1 p-2 rounded-lg text-[0.75rem] font-bold bg-red-500/15 border border-red-500/30 text-red-500">
          <span class="flex items-center gap-1">
            <ExclamationTriangleIcon class="w-4 h-4 text-red-600" /> {{ voucherError }}
          </span>
        </AppBadge>

        <!-- Voucher Hints Pills -->
        <div v-if="!appliedVoucher" class="flex items-center gap-1.5 flex-wrap">
          <span class="text-[0.7rem] text-slate-500">Rekomendasi Voucher:</span>
          <AppButton 
            v-for="v in availableVouchersHint" 
            :key="v.code" 
            variant="outline"
            class="px-2 py-0.5 rounded-full text-[0.68rem] font-bold border-dashed text-indigo-600 border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 hover:bg-indigo-50 dark:hover:bg-indigo-900/30 hover:border-indigo-600"
            @click="useHint(v.code)"
          >
            {{ v.code }}
          </AppButton>
        </div>
      </div>

      <div class="h-px bg-slate-200 dark:bg-slate-700 my-1"></div>

      <div class="flex justify-between items-center text-[0.85rem] text-slate-600 dark:text-slate-400">
        <span>Subtotal</span>
        <span>Rp {{ formatPrice(subtotal) }}</span>
      </div>

      <div class="flex justify-between items-center text-[0.85rem] text-slate-600 dark:text-slate-400">
        <span>Diskon {{ appliedVoucher ? '(' + appliedVoucher.code + ')' : '(Manual)' }}</span>
        <div class="flex items-center bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-md px-1.5 py-0.5">
          <span class="text-[0.75rem] text-slate-500 mr-1">Rp</span>
          <AppInput 
            type="number" 
            class="w-[75px] bg-transparent border-none text-right text-[0.85rem] text-slate-900 dark:text-slate-100 outline-none" 
            :modelValue="discount" 
            @update:modelValue="val => onManualDiscountInput(String(val))"
            placeholder="0"
            min="0"
          />
        </div>
      </div>

      <div class="flex justify-between items-center text-[0.85rem] text-slate-600 dark:text-slate-400">
        <span>Pajak ({{ taxPercentage }}%)</span>
        <span>Rp {{ formatPrice(taxAmount) }}</span>
      </div>

      <div class="h-px bg-slate-200 dark:bg-slate-700 my-1"></div>

      <div class="flex justify-between items-center text-base font-extrabold text-slate-900 dark:text-slate-100">
        <span>Total Bayar</span>
        <span class="text-[1.25rem] text-indigo-600 dark:text-indigo-400 font-extrabold">Rp {{ formatPrice(grandTotal) }}</span>
      </div>

      <AppButton variant="primary" class="w-full p-3 text-base mt-1.5 flex items-center justify-center gap-2" @click="$emit('open-payment')">
        <span>Bayar Sekarang</span>
        <ArrowRightIcon class="w-4 h-4" />
      </AppButton>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import api from '@/utils/api';
import type { CartItem, Voucher } from '../types';
import AppButton from './ui/AppButton.vue';
import AppInput from './ui/AppInput.vue';
import AppBadge from './ui/appBadge.vue';
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
