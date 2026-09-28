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

      <!-- Sub-categories Dropdown Bar (Artist & Tipe Produk) -->
      <div class="subcategories-bar">
        <div class="sub-filter-item">
          <select class="sub-filter-select" v-model="selectedArtist">
            <option value="">Artist</option>
            <option v-for="a in availableArtists" :key="a" :value="a">{{ a }}</option>
          </select>
        </div>

        <div class="sub-filter-item">
          <select class="sub-filter-select" v-model="selectedProductType">
            <option value="">Tipe Produk</option>
            <option v-for="t in availableProductTypes" :key="t" :value="t">{{ t }}</option>
          </select>
        </div>

        <button 
          v-if="selectedArtist || selectedProductType" 
          class="btn-reset-sub" 
          @click="resetSubFilters"
          title="Reset Sub-Filter"
        >
          Reset
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
import api from '@/utils/api';
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
const selectedArtist = ref<string>('');
const selectedProductType = ref<string>('');

const availableArtists = computed(() => {
  const set = new Set<string>();
  products.value.forEach(p => { if (p.artist) set.add(p.artist); });
  return Array.from(set).sort();
});

const availableProductTypes = computed(() => {
  const set = new Set<string>();
  products.value.forEach(p => { if (p.product_type) set.add(p.product_type); });
  return Array.from(set).sort();
});

const resetSubFilters = () => {
  selectedArtist.value = '';
  selectedProductType.value = '';
};

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
    const res = await api.get('/categories');
    categories.value = res.data;
  } catch (err: any) {
    console.error('Fetch categories error:', err.response?.data?.error || err.message || 'Error occurred');
  }
};

const fetchProducts = async () => {
  isLoading.value = true;
  try {
    const res = await api.get('/products');
    products.value = res.data;
  } catch (err: any) {
    console.error('Fetch products error:', err.response?.data?.error || err.message || 'Error occurred');
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
    const matchesArt = !selectedArtist.value || p.artist === selectedArtist.value;
    const matchesType = !selectedProductType.value || p.product_type === selectedProductType.value;

    const q = searchQuery.value.toLowerCase();
    const matchesQuery = !q || 
      p.name.toLowerCase().includes(q) || 
      (p.barcode && p.barcode.toLowerCase().includes(q)) ||
      (p.artist && p.artist.toLowerCase().includes(q)) ||
      (p.product_type && p.product_type.toLowerCase().includes(q));

    return matchesCat && matchesArt && matchesType && matchesQuery;
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

    const res = await api.post('/orders', payload);
    const completedOrder = res.data as Order;
    lastCompletedOrder.value = completedOrder;
    isPaymentModalOpen.value = false;
    isReceiptModalOpen.value = true;
    clearCart();
    fetchProducts(); // refresh stock counts
  } catch (err: any) {
    const errorMsg = err.response?.data?.error || err.message || 'Error occurred';
    alert('Gagal memproses transaksi: ' + errorMsg);
  } finally {
    isSubmittingOrder.value = false;
  }
};

const onReceiptClose = () => {
  isReceiptModalOpen.value = false;
  lastCompletedOrder.value = null;
};
</script>


