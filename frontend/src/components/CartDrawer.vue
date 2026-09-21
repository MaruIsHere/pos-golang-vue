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
        <span class="cart-icon">🛒</span>
        <h2>Keranjang</h2>
        <span class="item-count">{{ totalItemCount }} Item</span>
      </div>
      <button v-if="cart.length > 0" class="btn-clear" @click="$emit('clear-cart')">
        Hapus Semua
      </button>
    </div>

    <!-- Empty State -->
    <div v-if="cart.length === 0" class="empty-cart">
      <div class="empty-icon">🛍️</div>
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
            <button class="qty-btn" @click="$emit('update-qty', { index, delta: -1 })">-</button>
            <span class="qty-val">{{ item.quantity }}</span>
            <button class="qty-btn" @click="$emit('update-qty', { index, delta: 1 })">+</button>
          </div>
          <span class="item-subtotal">Rp {{ formatPrice(item.product.price * item.quantity) }}</span>
          <button class="btn-remove" @click="$emit('remove-item', index)">✕</button>
        </div>
      </div>
    </div>

    <!-- Summary & Checkout Footer -->
    <div v-if="cart.length > 0" class="cart-footer">
      <div class="summary-row">
        <span>Subtotal</span>
        <span>Rp {{ formatPrice(subtotal) }}</span>
      </div>

      <div class="summary-row discount-row">
        <span>Diskon (Rp)</span>
        <input 
          type="number" 
          class="discount-input" 
          :value="discount" 
          @input="$emit('update-discount', parseFloat($event.target.value) || 0)"
          placeholder="0"
          min="0"
        />
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

      <button class="btn btn-primary btn-checkout" @click="$emit('open-payment')">
        <span>Bayar Sekarang</span>
        <span>➡️</span>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  cart: { type: Array, required: true },
  discount: { type: Number, default: 0 },
  taxPercentage: { type: Number, default: 10 },
  isOpenMobile: { type: Boolean, default: false }
});

defineEmits(['update-qty', 'remove-item', 'clear-cart', 'update-discount', 'open-payment', 'toggle-mobile']);

const formatPrice = (val) => new Intl.NumberFormat('id-ID').format(val || 0);

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
</script>

<style scoped>
.cart-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.7);
  z-index: 80;
}

.cart-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

@media (max-width: 767px) {
  .cart-container {
    position: fixed;
    bottom: 64px;
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
  background: rgba(15, 23, 42, 0.4);
}

.header-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.header-title h2 {
  font-size: 1.1rem;
  font-weight: 800;
}

.item-count {
  background: rgba(99, 102, 241, 0.2);
  color: #818cf8;
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
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

.empty-icon {
  font-size: 3rem;
  margin-bottom: 0.75rem;
  opacity: 0.6;
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
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
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
  color: var(--text-secondary);
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
  background: rgba(30, 41, 59, 0.8);
  border-radius: 8px;
  padding: 0.15rem;
  border: 1px solid var(--border-color);
}

.qty-btn {
  width: 26px;
  height: 26px;
  border-radius: 6px;
  border: none;
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  font-weight: 700;
  cursor: pointer;
}

.qty-val {
  font-size: 0.85rem;
  font-weight: 800;
  min-width: 20px;
  text-align: center;
}

.item-subtotal {
  font-size: 0.85rem;
  font-weight: 800;
  color: var(--accent-secondary);
}

.btn-remove {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 0.9rem;
  cursor: pointer;
  padding: 0.2rem;
}
.btn-remove:hover {
  color: var(--accent-danger);
}

.cart-footer {
  padding: 1.25rem;
  border-top: 1px solid var(--border-color);
  background: rgba(15, 23, 42, 0.6);
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.discount-input {
  width: 90px;
  padding: 0.25rem 0.5rem;
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: #fff;
  font-family: var(--font-family);
  font-size: 0.85rem;
  text-align: right;
}

.divider {
  height: 1px;
  background: var(--border-color);
  margin: 0.25rem 0;
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
  margin-top: 0.5rem;
}
</style>
