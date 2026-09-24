<template>
  <div class="register-layout">
    <!-- Left Area: Catalog & Products -->
    <div class="catalog-section">
      <!-- Search Bar & Barcode Scanner Simulation -->
      <div class="search-bar-container glass-panel">
        <div class="search-input-wrapper">
          <MagnifyingGlassIcon class="w-5 h-5 text-slate-400 search-icon" />
          <input 
            type="text" 
            class="search-input" 
            v-model="searchQuery" 
            placeholder="Cari nama produk atau ketik / scan barcode SKU..." 
            @keyup.enter="onBarcodeSubmit"
          />
          <button v-if="searchQuery" class="btn-clear-search" @click="searchQuery = ''">
            <XMarkIcon class="w-4 h-4 text-slate-400" />
          </button>
        </div>

        <button class="btn-scan" title="Buka Scanner Barcode" @click="isScannerModalOpen = true">
          <CameraIcon class="w-4 h-4" />
          <span>Scan</span>
        </button>
      </div>

      <!-- Categories Pills -->
      <div class="categories-container">
        <button 
          class="cat-pill" 
          :class="{ active: selectedCategoryId === null }"
          @click="selectedCategoryId = null"
        >
          Semua Produk
        </button>
        <button 
          v-for="cat in categories" 
          :key="cat.id" 
          class="cat-pill" 
          :class="{ active: selectedCategoryId === cat.id }"
          @click="selectedCategoryId = cat.id"
        >
          {{ cat.name }}
        </button>
      </div>

      <!-- Products Grid -->
      <div v-if="isLoading" class="loading-state">
        <div class="spinner"></div>
        <p>Memuat data produk...</p>
      </div>

      <div v-else-if="filteredProducts.length === 0" class="empty-products glass-panel">
        <MagnifyingGlassIcon class="w-10 h-10 text-slate-400" />
        <h3>Produk tidak ditemukan</h3>
        <p>Coba gunakan kata kunci pencarian atau kategori lain</p>
      </div>

      <div v-else class="product-grid">
        <ProductCard 
          v-for="prod in filteredProducts" 
          :key="prod.id" 
          :product="prod"
          @add-to-cart="addToCart"
        />
      </div>
    </div>

    <!-- Right Area: Cart Drawer -->
    <div class="cart-section">
      <CartDrawer 
        :cart="cart"
        :discount="discount"
        :tax-percentage="taxPercentage"
        :is-open-mobile="isMobileCartOpen"
        @update-qty="updateCartQty"
        @remove-item="removeCartItem"
        @clear-cart="clearCart"
        @update-discount="discount = $event"
        @open-payment="isPaymentModalOpen = true"
        @toggle-mobile="isMobileCartOpen = !isMobileCartOpen"
      />
    </div>

    <!-- Floating Mobile Cart Trigger -->
    <div v-if="cart.length > 0" class="mobile-cart-float" @click="isMobileCartOpen = true">
      <div class="float-left">
        <ShoppingCartIcon class="w-5 h-5" />
        <span class="float-count">{{ totalCartItems }} Item</span>
      </div>
      <div class="float-right">
        <span class="float-total">Rp {{ formatPrice(grandTotal) }}</span>
        <ArrowUpIcon class="w-4 h-4" />
      </div>
    </div>

    <!-- Modals -->
    <BarcodeScannerModal
      v-if="isScannerModalOpen"
      :products="products"
      @close="isScannerModalOpen = false"
      @scan-success="addToCart"
    />

    <PaymentModal 
      v-if="isPaymentModalOpen"
      :grand-total="grandTotal"
      :is-submitting="isSubmittingOrder"
      @close="isPaymentModalOpen = false"
      @submit-order="handleCheckout"
    />

    <ReceiptModal 
      v-if="isReceiptModalOpen && lastCompletedOrder"
      :order="lastCompletedOrder"
      @close="onReceiptClose"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import ProductCard from '../components/ProductCard.vue';
import CartDrawer from '../components/CartDrawer.vue';
import PaymentModal from '../components/PaymentModal.vue';
import ReceiptModal from '../components/ReceiptModal.vue';
import BarcodeScannerModal from '../components/BarcodeScannerModal.vue';
import { useSettingsStore } from '../stores/settings';
import type { Category, Product, CartItem, Order, CreateOrderPayload } from '../types';
import { MagnifyingGlassIcon, CameraIcon, XMarkIcon, ShoppingCartIcon, ArrowUpIcon } from '@heroicons/vue/24/outline';

defineEmits(['refresh-products']);

const settingsStore = useSettingsStore();
const { settings: storeSetting } = storeToRefs(settingsStore);

const categories = ref<Category[]>([]);
const products = ref<Product[]>([]);
const isLoading = ref(true);

const searchQuery = ref('');
const selectedCategoryId = ref<number | null>(null);

// Cart State
const cart = ref<CartItem[]>([]);
const discount = ref(0);
const isMobileCartOpen = ref(false);

// Modal States
const isScannerModalOpen = ref(false);
const isPaymentModalOpen = ref(false);
const isReceiptModalOpen = ref(false);
const lastCompletedOrder = ref<Order | null>(null);
const isSubmittingOrder = ref(false);

const taxPercentage = computed(() => storeSetting.value.tax_percentage || 10);

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

// API Data Fetching
const fetchCategories = async () => {
  try {
    const res = await fetch('/api/categories');
    if (res.ok) categories.value = await res.json();
  } catch (err) {
    console.error('Fetch categories error:', err);
  }
};

const fetchProducts = async () => {
  isLoading.value = true;
  try {
    const res = await fetch('/api/products');
    if (res.ok) products.value = await res.json();
  } catch (err) {
    console.error('Fetch products error:', err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  fetchCategories();
  fetchProducts();
});

// Filtering
const filteredProducts = computed(() => {
  return products.value.filter(p => {
    const matchesCat = selectedCategoryId.value === null || p.category_id === selectedCategoryId.value;
    const q = searchQuery.value.toLowerCase();
    const matchesQuery = !q || p.name.toLowerCase().includes(q) || (p.barcode && p.barcode.toLowerCase().includes(q));
    return matchesCat && matchesQuery;
  });
});

// Cart Actions
const addToCart = (product: Product): void => {
  const existingIndex = cart.value.findIndex(item => item.product.id === product.id);
  if (existingIndex > -1) {
    if (cart.value[existingIndex].quantity < product.stock) {
      cart.value[existingIndex].quantity++;
    }
  } else {
    cart.value.push({ product, quantity: 1, notes: '' });
  }
};

const updateCartQty = ({ index, delta }: { index: number; delta: number }): void => {
  const item = cart.value[index];
  if (!item) return;
  const newQty = item.quantity + delta;
  if (newQty <= 0) {
    cart.value.splice(index, 1);
  } else if (newQty <= item.product.stock) {
    item.quantity = newQty;
  }
};

const removeCartItem = (index: number): void => {
  cart.value.splice(index, 1);
};

const clearCart = () => {
  cart.value = [];
  discount.value = 0;
};

const totalCartItems = computed(() => cart.value.reduce((s, i) => s + i.quantity, 0));
const subtotal = computed(() => cart.value.reduce((s, i) => s + (i.product.price * i.quantity), 0));
const taxAmount = computed(() => {
  const taxable = Math.max(0, subtotal.value - discount.value);
  return (taxable * taxPercentage.value) / 100;
});
const grandTotal = computed(() => Math.max(0, subtotal.value - discount.value) + taxAmount.value);

// Barcode Scan simulation
const onBarcodeSubmit = () => {
  if (!searchQuery.value) return;
  const match = products.value.find(p => p.barcode === searchQuery.value.trim());
  if (match) {
    addToCart(match);
    searchQuery.value = '';
  }
};

const simulateScan = () => {
  if (products.value.length > 0) {
    const randomProd = products.value[Math.floor(Math.random() * products.value.length)];
    addToCart(randomProd);
  }
};

// Checkout API handler
const handleCheckout = async ({ customer_name, payment_method, paid_amount }: { customer_name: string; payment_method: string; paid_amount: number }): Promise<void> => {
  isSubmittingOrder.value = true;
  try {
    const payload: CreateOrderPayload = {
      customer_name,
      payment_method,
      paid_amount,
      discount: discount.value,
      tax: taxAmount.value,
      items: cart.value.map(i => ({
        product_id: i.product.id,
        quantity: i.quantity,
        notes: i.notes
      }))
    };

    const res = await fetch('/api/orders', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });

    if (res.ok) {
      const completedOrder = (await res.json()) as Order;
      lastCompletedOrder.value = completedOrder;
      isPaymentModalOpen.value = false;
      isReceiptModalOpen.value = true;
      clearCart();
      fetchProducts(); // refresh stock counts
    } else {
      const errData = await res.json();
      alert('Gagal memproses transaksi: ' + (errData.error || 'Terjadi kesalahan'));
    }
  } catch (err) {
    alert('Koneksi ke server gagal: ' + (err as Error).message);
  } finally {
    isSubmittingOrder.value = false;
  }
};

const onReceiptClose = () => {
  isReceiptModalOpen.value = false;
  lastCompletedOrder.value = null;
};
</script>

<style scoped>
.register-layout {
  display: flex;
  height: 100%;
  gap: 1.25rem;
  position: relative;
  overflow: hidden;
}

.catalog-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  overflow-y: auto;
  padding-bottom: 80px; /* spacing for mobile float nav */
}

@media (min-width: 768px) {
  .catalog-section {
    padding-bottom: 1rem;
  }
}

.cart-section {
  width: 380px;
  height: 100%;
  display: none;
}

@media (min-width: 768px) {
  .cart-section {
    display: block;
  }
}

/* Search Bar */
.search-bar-container {
  display: flex;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  align-items: center;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  flex-shrink: 0;
}

.search-input-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 0.5rem 0.85rem;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-family: var(--font-family);
  font-size: 0.9rem;
  outline: none;
}

.btn-clear-search {
  background: transparent;
  border: none;
  cursor: pointer;
}

.btn-scan {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.65rem 1rem;
  background: rgba(99, 102, 241, 0.12);
  border: 1px solid rgba(99, 102, 241, 0.3);
  color: var(--accent-primary);
  border-radius: 8px;
  font-weight: 700;
  font-size: 0.85rem;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.15s ease;
}
.btn-scan:hover {
  background: rgba(99, 102, 241, 0.2);
}

/* Categories Pills */
.categories-container {
  display: flex;
  gap: 0.5rem;
  overflow-x: auto;
  padding: 0.25rem 0.25rem 0.5rem 0.25rem;
  flex-shrink: 0;
  min-height: 48px;
  align-items: center;
}

.categories-container::-webkit-scrollbar {
  display: none;
}

.cat-pill {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 1rem;
  min-height: 36px;
  border-radius: 999px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  font-family: var(--font-family);
  font-size: 0.85rem;
  font-weight: 600;
  line-height: 1;
  white-space: nowrap;
  cursor: pointer;
  transition: all 0.15s ease;
}

.cat-pill:hover {
  background: var(--bg-card-hover);
  color: var(--text-primary);
}

.cat-pill.active {
  background: var(--accent-primary);
  color: #ffffff;
  border-color: var(--accent-primary);
  box-shadow: 0 2px 6px rgba(79, 70, 229, 0.3);
}

.loading-state, .empty-products {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  text-align: center;
  gap: 0.75rem;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  color: var(--text-secondary);
}

.spinner {
  width: 36px;
  height: 36px;
  border: 3px solid var(--border-color);
  border-top-color: var(--accent-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Floating Mobile Cart Trigger */
.mobile-cart-float {
  position: fixed;
  bottom: 74px;
  left: 1rem;
  right: 1rem;
  background: var(--accent-primary);
  color: #ffffff;
  padding: 0.85rem 1.25rem;
  border-radius: 999px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.4);
  z-index: 70;
  cursor: pointer;
}

@media (min-width: 768px) {
  .mobile-cart-float {
    display: none;
  }
}

.float-left, .float-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.float-count {
  font-size: 0.85rem;
  font-weight: 700;
  background: rgba(255, 255, 255, 0.2);
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
}

.float-total {
  font-size: 1.1rem;
  font-weight: 800;
}
</style>
