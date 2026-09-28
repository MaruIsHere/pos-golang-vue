<template>
  <div class="h-full w-full flex flex-col lg:grid lg:grid-cols-3 xl:grid-cols-4 gap-4 relative">
    <!-- Left Area: Catalog & Products -->
    <div class="lg:col-span-2 xl:col-span-3 flex flex-col gap-4 h-full min-h-0">
      
      <!-- Search Bar & Barcode Scanner -->
      <div class="flex items-center gap-3 p-3 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl shadow-sm shrink-0">
        <div class="relative flex-1 flex items-center bg-slate-100 dark:bg-slate-900/50 rounded-xl px-3 py-2 border border-slate-200 dark:border-slate-700/50 focus-within:border-indigo-500 focus-within:ring-2 focus-within:ring-indigo-500/20 transition-all">
          <MagnifyingGlassIcon class="w-5 h-5 text-slate-400 shrink-0" />
          <input 
            type="text" 
            class="w-full bg-transparent border-none outline-none text-sm font-semibold text-slate-800 dark:text-slate-100 ml-2 placeholder:text-slate-500 placeholder:font-normal" 
            v-model="searchQuery" 
            placeholder="Cari nama produk atau ketik / scan barcode SKU..." 
            @keyup.enter="onBarcodeSubmit"
          />
          <button v-if="searchQuery" class="p-1 hover:bg-slate-200 dark:hover:bg-slate-800 rounded-full transition-colors" @click="searchQuery = ''">
            <XMarkIcon class="w-4 h-4 text-slate-400" />
          </button>
        </div>

        <button class="flex items-center gap-2 px-4 py-2.5 bg-indigo-50 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-400 hover:bg-indigo-100 dark:hover:bg-indigo-900/60 rounded-xl font-bold text-sm shrink-0 transition-colors border border-indigo-200 dark:border-indigo-800/50 shadow-sm active:scale-95" title="Buka Scanner Barcode" @click="isScannerModalOpen = true">
          <CameraIcon class="w-5 h-5" />
          <span class="hidden sm:inline">Scan</span>
        </button>
      </div>

      <!-- Categories Pills -->
      <div class="flex items-center gap-2 overflow-x-auto pb-1 shrink-0 no-scrollbar scroll-smooth">
        <button 
          class="px-4 py-2 rounded-xl text-sm font-bold whitespace-nowrap transition-all border shadow-sm" 
          :class="selectedCategoryId === null ? 'bg-indigo-600 text-white border-indigo-600 dark:bg-indigo-500 dark:border-indigo-500' : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-700 hover:border-indigo-400 hover:text-indigo-600'"
          @click="selectedCategoryId = null"
        >
          Semua Produk
        </button>
        <button 
          v-for="cat in categories" 
          :key="cat.id" 
          class="px-4 py-2 rounded-xl text-sm font-bold whitespace-nowrap transition-all border shadow-sm" 
          :class="selectedCategoryId === cat.id ? 'bg-indigo-600 text-white border-indigo-600 dark:bg-indigo-500 dark:border-indigo-500' : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-700 hover:border-indigo-400 hover:text-indigo-600'"
          @click="selectedCategoryId = cat.id"
        >
          {{ cat.name }}
        </button>
      </div>

      <!-- Sub-categories Dropdown Bar -->
      <div class="flex flex-wrap items-center gap-2 shrink-0 pb-2 border-b border-slate-200 dark:border-slate-700/50">
        <div class="relative flex-1 sm:flex-none">
          <select class="px-3 py-1.5 text-xs font-bold bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-700 dark:text-slate-300 outline-none focus:border-indigo-500 w-full sm:w-auto shadow-sm cursor-pointer" v-model="selectedArtist">
            <option value="">Filter Artist</option>
            <option v-for="a in availableArtists" :key="a" :value="a">{{ a }}</option>
          </select>
        </div>

        <div class="relative flex-1 sm:flex-none">
          <select class="px-3 py-1.5 text-xs font-bold bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-700 dark:text-slate-300 outline-none focus:border-indigo-500 w-full sm:w-auto shadow-sm cursor-pointer" v-model="selectedProductType">
            <option value="">Tipe Produk</option>
            <option v-for="t in availableProductTypes" :key="t" :value="t">{{ t }}</option>
          </select>
        </div>

        <button 
          v-if="selectedArtist || selectedProductType" 
          class="px-3 py-1.5 text-xs font-bold text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-900/50 transition-colors shadow-sm ml-auto sm:ml-0" 
          @click="resetSubFilters"
          title="Reset Sub-Filter"
        >
          Reset Filter
        </button>
      </div>

      <!-- Products Grid -->
      <div v-if="isLoading" class="flex flex-col items-center justify-center flex-1 min-h-[300px] gap-3">
        <div class="w-8 h-8 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
        <p class="text-sm font-semibold text-slate-500">Memuat katalog produk...</p>
      </div>

      <div v-else-if="filteredProducts.length === 0" class="flex flex-col items-center justify-center flex-1 min-h-[300px] p-8 text-center bg-white/50 dark:bg-slate-800/30 border-2 border-dashed border-slate-200 dark:border-slate-700 rounded-3xl">
        <MagnifyingGlassIcon class="w-12 h-12 text-slate-300 dark:text-slate-600 mb-3" />
        <h3 class="text-lg font-bold text-slate-700 dark:text-slate-300">Produk tidak ditemukan</h3>
        <p class="text-sm text-slate-500 mt-1">Coba gunakan kata kunci pencarian atau kategori lain</p>
      </div>

      <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 2xl:grid-cols-5 gap-3 sm:gap-4 overflow-y-auto pb-24 lg:pb-4 min-h-0 pr-1 custom-scrollbar">
        <ProductCard 
          v-for="prod in filteredProducts" 
          :key="prod.id" 
          :product="prod"
          @add-to-cart="addToCart"
        />
      </div>
    </div>

    <!-- Right Area: Cart Drawer (Desktop) -->
    <div class="hidden lg:flex flex-col lg:col-span-1 h-full min-h-0 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl shadow-sm overflow-hidden">
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

    <!-- Floating Mobile Cart Trigger (Khusus Layar Kecil) -->
    <div v-if="cart.length > 0" class="lg:hidden fixed bottom-20 left-4 right-4 flex items-center justify-between p-4 bg-indigo-600 hover:bg-indigo-700 text-white rounded-2xl shadow-xl font-bold cursor-pointer z-40 active:scale-95 transition-all" @click="isMobileCartOpen = true">
      <div class="flex items-center gap-2">
        <div class="relative">
          <ShoppingCartIcon class="w-6 h-6" />
          <span class="absolute -top-2 -right-2 w-4 h-4 flex items-center justify-center bg-red-500 text-white text-[10px] font-black rounded-full border border-indigo-600">{{ totalCartItems }}</span>
        </div>
        <span class="text-sm">Keranjang</span>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-base tracking-wide">Rp {{ formatPrice(grandTotal) }}</span>
        <ArrowUpIcon class="w-4 h-4 animate-bounce" />
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


