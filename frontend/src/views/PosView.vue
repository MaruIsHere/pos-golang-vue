<template>
  <div class="h-[calc(100vh-8rem)] md:h-[calc(100vh-6.5rem)] w-full flex flex-col lg:grid lg:grid-cols-[1fr_350px] xl:grid-cols-[1fr_400px] gap-6 relative">
    
    <!-- Left Area: Catalog & Products -->
    <div class="flex flex-col gap-5 h-full min-h-0 relative">
      
      <!-- Apple-style Search Bar & Scan -->
      <div class="flex items-center gap-3 shrink-0 relative z-10">
        <div class="relative flex-1 flex items-center bg-white dark:bg-slate-800 rounded-full px-5 py-3 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.05)] border border-slate-100 dark:border-slate-700/50 focus-within:ring-4 focus-within:ring-indigo-500/10 focus-within:border-indigo-300 dark:focus-within:border-indigo-600 transition-all duration-300">
          <MagnifyingGlassIcon class="w-5 h-5 text-slate-400 shrink-0" />
          <input 
            type="text" 
            class="w-full bg-transparent border-none outline-none text-[0.95rem] font-medium text-slate-800 dark:text-slate-100 ml-3 placeholder:text-slate-400" 
            v-model="searchQuery" 
            placeholder="Cari nama produk, artist, atau scan barcode..." 
            @keyup.enter="onBarcodeSubmit"
          />
          <button v-if="searchQuery" class="p-1.5 bg-slate-100 hover:bg-slate-200 dark:bg-slate-700 dark:hover:bg-slate-600 rounded-full transition-colors text-slate-500" @click="searchQuery = ''">
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>

        <button class="flex items-center gap-2 px-6 py-3.5 bg-slate-900 text-white dark:bg-white dark:text-slate-900 hover:scale-105 active:scale-95 rounded-full font-bold text-sm shrink-0 transition-transform shadow-[0_4px_15px_-3px_rgba(0,0,0,0.15)]" title="Buka Scanner Barcode" @click="isScannerModalOpen = true">
          <CameraIcon class="w-5 h-5" />
          <span class="hidden sm:inline">Scanner</span>
        </button>
      </div>

      <!-- M3 Segmented Categories Pills -->
      <div class="flex items-center gap-2.5 overflow-x-auto pb-2 shrink-0 no-scrollbar snap-x">
        <button 
          class="px-5 py-2 rounded-full text-[0.85rem] font-bold whitespace-nowrap transition-all duration-300 shadow-sm snap-start" 
          :class="selectedCategoryId === null ? 'bg-indigo-600 text-white border-transparent shadow-indigo-600/20' : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700'"
          @click="selectedCategoryId = null"
        >
          Semua Produk
        </button>
        <button 
          v-for="cat in categories" 
          :key="cat.id" 
          class="px-5 py-2 rounded-full text-[0.85rem] font-bold whitespace-nowrap transition-all duration-300 shadow-sm snap-start" 
          :class="selectedCategoryId === cat.id ? 'bg-indigo-600 text-white border-transparent shadow-indigo-600/20' : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700'"
          @click="selectedCategoryId = cat.id"
        >
          {{ cat.name }}
        </button>
      </div>

      <!-- Compact Sub-categories Filter -->
      <div class="flex items-center gap-3 shrink-0">
        <select class="px-4 py-2 text-xs font-bold bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-700 dark:text-slate-300 outline-none focus:ring-2 focus:ring-indigo-500/20 shadow-sm cursor-pointer transition-all" v-model="selectedArtist">
          <option value="">Semua Artist</option>
          <option v-for="a in availableArtists" :key="a" :value="a">{{ a }}</option>
        </select>
        <select class="px-4 py-2 text-xs font-bold bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-700 dark:text-slate-300 outline-none focus:ring-2 focus:ring-indigo-500/20 shadow-sm cursor-pointer transition-all" v-model="selectedProductType">
          <option value="">Semua Tipe Produk</option>
          <option v-for="t in availableProductTypes" :key="t" :value="t">{{ t }}</option>
        </select>
        <button v-if="selectedArtist || selectedProductType" class="px-3 py-2 text-xs font-bold text-red-600 bg-red-50 hover:bg-red-100 dark:bg-red-900/20 dark:hover:bg-red-900/40 rounded-xl transition-colors shadow-sm ml-auto" @click="resetSubFilters">
          Reset
        </button>
      </div>

      <!-- Products Grid (Independent Scroll) -->
      <div v-if="isLoading" class="flex flex-col items-center justify-center flex-1 min-h-[300px] gap-4 bg-white/40 dark:bg-slate-800/20 rounded-[32px]">
        <div class="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
        <p class="text-sm font-semibold text-slate-500">Memuat katalog...</p>
      </div>

      <div v-else-if="filteredProducts.length === 0" class="flex flex-col items-center justify-center flex-1 min-h-[300px] p-8 text-center bg-white/60 dark:bg-slate-800/40 border-2 border-dashed border-slate-200 dark:border-slate-700 rounded-[32px]">
        <div class="p-4 bg-slate-100 dark:bg-slate-800 rounded-full mb-4">
          <MagnifyingGlassIcon class="w-10 h-10 text-slate-400" />
        </div>
        <h3 class="text-[1.1rem] font-bold text-slate-800 dark:text-slate-200">Produk tidak ditemukan</h3>
        <p class="text-sm text-slate-500 mt-1 max-w-sm">Coba gunakan kata kunci pencarian yang lebih singkat atau pilih kategori lain.</p>
      </div>

      <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 2xl:grid-cols-5 gap-4 sm:gap-5 overflow-y-auto pb-24 lg:pb-6 pr-2 -mr-2 min-h-0 custom-scrollbar content-start">
        <ProductCard 
          v-for="prod in filteredProducts" 
          :key="prod.id" 
          :product="prod"
          @add-to-cart="addToCart"
        />
      </div>
    </div>

    <!-- Cart Drawer (Handles both Desktop and Mobile) -->
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

    <!-- Floating Mobile Cart Trigger (Khusus Layar Kecil) -->
    <div v-if="cart.length > 0" class="lg:hidden fixed bottom-[5.5rem] left-4 right-4 flex items-center justify-between p-4 px-6 bg-slate-900 hover:bg-slate-800 dark:bg-white dark:hover:bg-slate-100 text-white dark:text-slate-900 rounded-full shadow-[0_8px_30px_rgb(0,0,0,0.15)] font-bold cursor-pointer z-[60] active:scale-[0.98] transition-all" @click="isMobileCartOpen = true">
      <div class="flex items-center gap-3">
        <div class="relative">
          <ShoppingCartIcon class="w-6 h-6" />
          <span class="absolute -top-2 -right-2 w-5 h-5 flex items-center justify-center bg-indigo-600 text-white text-[10px] font-black rounded-full border-2 border-slate-900 dark:border-white">{{ totalCartItems }}</span>
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


