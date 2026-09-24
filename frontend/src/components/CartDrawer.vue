<template>
  <!-- Mobile Backdrop -->
  <div 
    v-if="isOpenMobile" 
    class="cart-backdrop" 
    @click="$emit('toggle-mobile')"
  ></div>

  <aside class="cart-container glass-panel" :class="{ 'mobile-open': isOpenMobile }">
    <div class="cart-header">
      <div class="header-title">
        <ShoppingCartIcon class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
        <h2>Keranjang</h2>
        <span class="item-count">{{ totalItemCount }} Item</span>
      </div>
      <button v-if="cart.length > 0" class="btn-clear" @click="handleClearCart">
        Hapus Semua
      </button>
    </div>

    <!-- Empty State -->
    <div v-if="cart.length === 0" class="empty-cart">
      <div class="empty-icon-box">
        <ShoppingBagIcon class="w-12 h-12 text-slate-300 dark:text-slate-600" />
      </div>
      <p class="empty-text">Keranjang masih kosong</p>
      <span class="empty-sub">Pilih produk di sebelah kiri untuk ditambahkan</span>
    </div>

    <!-- Cart Items List -->
    <div v-else class="cart-items-list">
      <div v-for="(item, index) in cart" :key="item.product.id" class="cart-item">
        <div class="item-info">
          <h4 class="item-name">{{ item.product.name }}</h4>
          <span class="item-price">Rp {{ formatPrice(item.product.price) }}</span>
        </div>

        <div class="item-controls">
          <div class="qty-group">
            <button class="qty-btn" @click="$emit('update-qty', { index, delta: -1 })">
              <MinusIcon class="w-3.5 h-3.5" />
            </button>
            <span class="qty-val">{{ item.quantity }}</span>
            <button class="qty-btn" @click="$emit('update-qty', { index, delta: 1 })">
              <PlusIcon class="w-3.5 h-3.5" />
            </button>
          </div>
          <span class="item-subtotal">Rp {{ formatPrice(item.product.price * item.quantity) }}</span>
          <button class="btn-remove" @click="$emit('remove-item', index)">
            <TrashIcon class="w-4 h-4 text-slate-400 hover:text-red-500" />
          </button>
        </div>
      </div>
    </div>

    <!-- Summary & Checkout Footer -->
    <div v-if="cart.length > 0" class="cart-footer">
      <!-- Voucher Code Input Section -->
      <div class="voucher-container">
        <div class="voucher-input-wrapper">
          <TagIcon class="w-4 h-4 text-slate-400 voucher-icon" />
          <input 
            type="text" 
            class="voucher-input" 
            v-model="voucherInput" 
            placeholder="Kode Voucher (DISKON10...)"
            @keyup.enter="applyVoucher"
          />
          <button class="btn-voucher-apply" @click="applyVoucher">
            Gunakan
          </button>
        </div>

        <!-- Applied Voucher Alert Badge -->
        <div v-if="appliedVoucher" class="voucher-badge success">
          <div class="voucher-info">
            <span class="voucher-title flex items-center gap-1">
              <TicketIcon class="w-4 h-4 text-emerald-600" /> {{ appliedVoucher.code }}
            </span>
            <span class="voucher-desc">Potongan Rp {{ formatPrice(appliedVoucher.discountAmount) }} ({{ appliedVoucher.description }})</span>
          </div>
          <button class="btn-voucher-remove" title="Hapus Voucher" @click="removeVoucher">
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>

        <!-- Error Message -->
        <div v-if="voucherError" class="voucher-badge error">
          <span class="flex items-center gap-1">
            <ExclamationTriangleIcon class="w-4 h-4 text-red-600" /> {{ voucherError }}
          </span>
        </div>

        <!-- Voucher Hints Pills -->
        <div v-if="!appliedVoucher" class="voucher-hints">
          <span class="hints-label">Rekomendasi Voucher:</span>
          <button 
            v-for="v in availableVouchersHint" 
            :key="v.code" 
            class="hint-pill"
            @click="useHint(v.code)"
          >
            {{ v.code }}
          </button>
        </div>
      </div>

      <div class="divider"></div>

      <div class="summary-row">
        <span>Subtotal</span>
        <span>Rp {{ formatPrice(subtotal) }}</span>
      </div>

      <div class="summary-row discount-row">
        <span>Diskon {{ appliedVoucher ? '(' + appliedVoucher.code + ')' : '(Manual)' }}</span>
        <div class="discount-input-wrapper">
          <span class="input-rp">Rp</span>
          <input 
            type="number" 
            class="discount-input" 
            :value="discount" 
            @input="onManualDiscountInput(($event.target as HTMLInputElement).value)"
            placeholder="0"
            min="0"
          />
        </div>
      </div>

      <div class="summary-row">
        <span>Pajak ({{ taxPercentage }}%)</span>
        <span>Rp {{ formatPrice(taxAmount) }}</span>
      </div>

      <div class="divider"></div>

      <div class="summary-row grand-total-row">
        <span>Total Bayar</span>
        <span class="grand-total-val">Rp {{ formatPrice(grandTotal) }}</span>
      </div>

      <button class="btn btn-primary btn-checkout flex items-center justify-center gap-2" @click="$emit('open-payment')">
        <span>Bayar Sekarang</span>
        <ArrowRightIcon class="w-4 h-4" />
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { CartItem, Voucher } from '../types';
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
    const res = await fetch('/api/vouchers');
    if (res.ok) {
      const data = await res.json();
      if (Array.isArray(data) && data.length > 0) {
        availableVouchers.value = data;
      }
    }
  } catch (err) {
    console.error('Fetch vouchers error:', err);
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

<style scoped>
.cart-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(15, 23, 42, 0.4);
  z-index: 80;
}

.cart-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

@media (max-width: 767px) {
  .cart-container {
    position: fixed;
    bottom: 60px;
    left: 0;
    right: 0;
    height: 75vh;
    z-index: 90;
    border-radius: var(--radius-xl) var(--radius-xl) 0 0;
    transform: translateY(105%);
    transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  }
  .cart-container.mobile-open {
    transform: translateY(0);
  }
}

.cart-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-primary);
}

.header-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.header-title h2 {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-primary);
}

.item-count {
  background: rgba(99, 102, 241, 0.12);
  color: var(--accent-primary);
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
  border: 1px solid rgba(99, 102, 241, 0.3);
}

.btn-clear {
  background: transparent;
  border: none;
  color: var(--accent-danger);
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
}

.empty-cart {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  text-align: center;
}

.empty-icon-box {
  margin-bottom: 0.75rem;
}

.empty-text {
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 0.25rem;
}

.empty-sub {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.cart-items-list {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.cart-item {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.item-info {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.item-name {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--text-primary);
}

.item-price {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.item-controls {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.qty-group {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: var(--bg-primary);
  border-radius: 8px;
  padding: 0.15rem;
  border: 1px solid var(--border-color);
}

.qty-btn {
  width: 26px;
  height: 26px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  background: var(--bg-card);
  color: var(--text-primary);
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}
.qty-btn:hover {
  background: var(--bg-card-hover);
}

.qty-val {
  font-size: 0.85rem;
  font-weight: 800;
  min-width: 20px;
  text-align: center;
  color: var(--text-primary);
}

.item-subtotal {
  font-size: 0.85rem;
  font-weight: 800;
  color: var(--accent-secondary);
}

.btn-remove {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0.2rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cart-footer {
  padding: 1rem 1.25rem;
  border-top: 1px solid var(--border-color);
  background: var(--bg-primary);
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

/* Voucher Code Styling */
.voucher-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.voucher-input-wrapper {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 0.25rem 0.5rem;
}

.voucher-input {
  flex: 1;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-family: var(--font-family);
  font-size: 0.8rem;
  text-transform: uppercase;
  outline: none;
}

.btn-voucher-apply {
  padding: 0.35rem 0.75rem;
  background: var(--accent-primary);
  border: none;
  border-radius: 6px;
  color: #ffffff;
  font-weight: 700;
  font-size: 0.75rem;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn-voucher-apply:hover {
  opacity: 0.9;
}

.voucher-badge {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 700;
}

.voucher-badge.success {
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid rgba(16, 185, 129, 0.3);
  color: #34d399;
}

.voucher-badge.error {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
}

.voucher-info {
  display: flex;
  flex-direction: column;
}

.voucher-title {
  font-weight: 800;
}

.voucher-desc {
  font-size: 0.7rem;
  opacity: 0.9;
}

.btn-voucher-remove {
  background: transparent;
  border: none;
  color: inherit;
  cursor: pointer;
  padding: 0.2rem;
}

.voucher-hints {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  flex-wrap: wrap;
}

.hints-label {
  font-size: 0.7rem;
  color: var(--text-muted);
}

.hint-pill {
  padding: 0.15rem 0.45rem;
  background: var(--bg-secondary);
  border: 1px dashed var(--border-color);
  border-radius: 999px;
  color: var(--accent-primary);
  font-size: 0.68rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}

.hint-pill:hover {
  background: rgba(99, 102, 241, 0.15);
  border-color: var(--accent-primary);
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.discount-input-wrapper {
  display: flex;
  align-items: center;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 0.15rem 0.4rem;
}

.input-rp {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-right: 0.25rem;
}

.discount-input {
  width: 75px;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-family: var(--font-family);
  font-size: 0.85rem;
  text-align: right;
  outline: none;
}

.divider {
  height: 1px;
  background: var(--border-color);
  margin: 0.2rem 0;
}

.grand-total-row {
  font-size: 1rem;
  font-weight: 800;
  color: var(--text-primary);
}

.grand-total-val {
  font-size: 1.25rem;
  color: var(--accent-secondary);
  font-weight: 800;
}

.btn-checkout {
  width: 100%;
  padding: 0.85rem;
  font-size: 1rem;
  margin-top: 0.35rem;
}
</style>
