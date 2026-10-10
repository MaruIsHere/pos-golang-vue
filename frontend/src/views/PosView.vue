<template>
  <div class="h-[calc(100vh-8rem)] md:h-[calc(100vh-6.5rem)] w-full flex flex-col lg:grid lg:grid-cols-[1fr_350px] xl:grid-cols-[1fr_400px] gap-6 relative">
    
    <!-- Left Area: Catalog & Products -->
    <div class="flex flex-col gap-4 h-full min-h-0 relative">
      
      <!-- Outlet POS Context Banner -->
      <div v-if="storeContextStore.activeStore" class="flex items-center justify-between px-4 py-2 bg-indigo-50/80 dark:bg-indigo-900/30 border border-indigo-100 dark:border-indigo-800/40 rounded-2xl shrink-0">
        <div class="flex items-center gap-2">
          <BuildingStorefrontIcon class="w-4 h-4 text-indigo-600 dark:text-indigo-400" />
          <span class="text-xs font-bold text-slate-700 dark:text-slate-200">
            Etalase Kasir Toko: <strong class="text-indigo-600 dark:text-indigo-400 font-extrabold">{{ storeContextStore.activeStore.name }}</strong> ({{ storeContextStore.activeStore.code }})
          </span>
        </div>
        <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-indigo-100 text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300">
          {{ filteredProducts.length }} Produk Kasir
        </span>
      </div>

      <!-- Apple-style Search Bar & Scan -->
      <div class="flex items-center gap-3 shrink-0 relative z-10">
        <div class="relative flex-1 flex items-center bg-white dark:bg-slate-800 rounded-full px-5 py-3 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.05)] border border-slate-100 dark:border-slate-700/50 focus-within:ring-4 focus-within:ring-indigo-500/10 focus-within:border-indigo-300 dark:focus-within:border-indigo-600 transition-all duration-300">
          <MagnifyingGlassIcon class="w-5 h-5 text-slate-400 shrink-0" />
          <input 
            type="text" 
            class="w-full bg-transparent border-none outline-none text-[0.95rem] font-medium text-slate-800 dark:text-slate-100 ml-3 placeholder:text-slate-400" 
            v-model="searchQuery" 
            placeholder="Cari nama produk, merk, atau scan barcode..." 
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

      <!-- Cashier Filter Section (Segmented Pills with Labels & Separators) -->
      <div class="flex flex-col gap-2 shrink-0 bg-white/60 dark:bg-slate-800/50 p-2.5 sm:p-3 rounded-2xl border border-slate-200/60 dark:border-slate-700/60 shadow-xs">
        
        <!-- Row 1: Kategori Filter -->
        <div class="flex items-center gap-2.5 overflow-x-auto no-scrollbar snap-x py-0.5">
          <div class="flex items-center gap-1.5 shrink-0 px-2.5 py-1 text-xs font-bold text-slate-700 dark:text-slate-200 bg-slate-100 dark:bg-slate-700/60 rounded-lg">
            <TagIcon class="w-3.5 h-3.5 text-indigo-500" />
            <span>Kategori</span>
          </div>
          <div class="h-4 w-px bg-slate-200 dark:bg-slate-700 shrink-0"></div>
          
          <button 
            class="px-4 py-1.5 rounded-full text-xs font-bold whitespace-nowrap transition-all duration-200 shadow-xs shrink-0 snap-start" 
            :class="selectedCategoryId === null ? 'bg-indigo-600 text-white shadow-indigo-600/20' : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700'"
            @click="selectCategory(null)"
          >
            Semua Produk
          </button>
          <button 
            v-for="cat in rootCategories" 
            :key="cat.id" 
            class="px-4 py-1.5 rounded-full text-xs font-bold whitespace-nowrap transition-all duration-200 shadow-xs shrink-0 snap-start" 
            :class="selectedCategoryId === cat.id ? 'bg-indigo-600 text-white shadow-indigo-600/20' : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700'"
            @click="selectCategory(cat.id)"
          >
            {{ categoryLabel(cat) }}
          </button>
        </div>

        <!-- Row 2: Merk Filter (Show if availableArtists has items) -->
        <div v-if="availableArtists.length > 0" class="flex items-center gap-2.5 overflow-x-auto no-scrollbar snap-x py-0.5 border-t border-slate-100 dark:border-slate-700/50 pt-2">
          <div class="flex items-center gap-1.5 shrink-0 px-2.5 py-1 text-xs font-bold text-slate-700 dark:text-slate-200 bg-slate-100 dark:bg-slate-700/60 rounded-lg">
            <BookmarkIcon class="w-3.5 h-3.5 text-emerald-500" />
            <span>Merk</span>
          </div>
          <div class="h-4 w-px bg-slate-200 dark:bg-slate-700 shrink-0"></div>

          <button 
            class="px-4 py-1.5 rounded-full text-xs font-bold whitespace-nowrap transition-all duration-200 shadow-xs shrink-0 snap-start" 
            :class="selectedArtist === null ? 'bg-emerald-600 text-white shadow-emerald-600/20' : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700'"
            @click="selectedArtist = null"
          >
            Semua Merk
          </button>
          <button 
            v-for="artist in availableArtists" 
            :key="artist.id" 
            class="px-4 py-1.5 rounded-full text-xs font-bold whitespace-nowrap transition-all duration-200 shadow-xs shrink-0 snap-start" 
            :class="selectedArtist === artist.name ? 'bg-emerald-600 text-white shadow-emerald-600/20' : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700'"
            @click="selectedArtist = artist.name"
          >
            {{ artist.name }}
          </button>
        </div>

        <!-- Row 3: Tipe Filter (Show if availableProductTypes has items) -->
        <div v-if="availableProductTypes.length > 0" class="flex items-center gap-2.5 overflow-x-auto no-scrollbar snap-x py-0.5 border-t border-slate-100 dark:border-slate-700/50 pt-2">
          <div class="flex items-center gap-1.5 shrink-0 px-2.5 py-1 text-xs font-bold text-slate-700 dark:text-slate-200 bg-slate-100 dark:bg-slate-700/60 rounded-lg">
            <Squares2X2Icon class="w-3.5 h-3.5 text-amber-500" />
            <span>Tipe</span>
          </div>
          <div class="h-4 w-px bg-slate-200 dark:bg-slate-700 shrink-0"></div>

          <button 
            class="px-4 py-1.5 rounded-full text-xs font-bold whitespace-nowrap transition-all duration-200 shadow-xs shrink-0 snap-start" 
            :class="selectedProductType === null ? 'bg-amber-600 text-white shadow-amber-600/20' : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700'"
            @click="selectedProductType = null"
          >
            Semua Tipe
          </button>
          <button 
            v-for="type in availableProductTypes" 
            :key="type.id" 
            class="px-4 py-1.5 rounded-full text-xs font-bold whitespace-nowrap transition-all duration-200 shadow-xs shrink-0 snap-start" 
            :class="selectedProductType === type.name ? 'bg-amber-600 text-white shadow-amber-600/20' : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 hover:bg-slate-50 dark:hover:bg-slate-700'"
            @click="selectedProductType = type.name"
          >
            {{ type.name }}
          </button>
        </div>

        <!-- Filter Reset Action Row -->
        <div v-if="selectedCategoryId !== null || selectedArtist !== null || selectedProductType !== null" class="flex items-center justify-between border-t border-slate-100 dark:border-slate-700/50 pt-2 px-1">
          <span class="text-[11px] font-medium text-slate-400">Filter aktif diterapkan</span>
          <button class="px-2.5 py-1 text-[11px] font-bold text-red-600 hover:text-red-700 dark:text-red-400 bg-red-50 hover:bg-red-100 dark:bg-red-900/30 rounded-lg transition-colors flex items-center gap-1" @click="resetAllFilters">
            <XMarkIcon class="w-3.5 h-3.5" />
            <span>Reset Filter</span>
          </button>
        </div>

      </div>
      <div v-if="catalogError" class="flex items-center justify-between gap-3 rounded-lg border border-rose-200 bg-rose-50 px-3 py-2 text-xs text-rose-700 dark:border-rose-900 dark:bg-rose-950/30 dark:text-rose-300" role="alert">
        <span>{{ catalogError }}</span>
        <button type="button" class="shrink-0 font-bold underline" @click="fetchCatalogs">Coba lagi</button>
      </div>

      <!-- Loading Skeleton -->
      <div v-if="isLoading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 2xl:grid-cols-5 gap-4 sm:gap-5 overflow-hidden content-start pt-2">
        <div v-for="i in 10" :key="i" class="flex flex-col rounded-[24px] bg-white/50 dark:bg-slate-800/50 border border-slate-200/50 dark:border-slate-700/50 shadow-sm h-[250px] sm:h-[260px] animate-pulse">
          <div class="relative w-full h-[140px] sm:h-[150px] p-2 shrink-0">
             <div class="w-full h-full rounded-[18px] bg-slate-200/70 dark:bg-slate-700/70"></div>
          </div>
          <div class="flex flex-col px-4 pb-4 pt-2 flex-1">
             <div class="h-2.5 w-1/3 bg-slate-200/70 dark:bg-slate-700/70 rounded mb-3"></div>
             <div class="h-3.5 w-3/4 bg-slate-200/70 dark:bg-slate-700/70 rounded mb-2"></div>
             <div class="h-3.5 w-1/2 bg-slate-200/70 dark:bg-slate-700/70 rounded mb-3"></div>
             <div class="h-3.5 w-2/3 bg-indigo-100 dark:bg-indigo-900/30 rounded mt-auto"></div>
          </div>
        </div>
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
      @set-qty="setCartQty"
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
import { ref, computed, onMounted, watch } from 'vue';
import api from '@/utils/api';
import { storeToRefs } from 'pinia';
import ProductCard from '../components/ProductCard.vue';
import CartDrawer from '../components/CartDrawer.vue';
import PaymentModal from '../components/PaymentModal.vue';
import ReceiptModal from '../components/ReceiptModal.vue';
import BarcodeScannerModal from '../components/BarcodeScannerModal.vue';
import { useSettingsStore } from '../stores/settings';
import { useAuthStore } from '../stores/auth';
import { useStoreContextStore } from '../stores/storeContext';
import type { Artist, Category, Product, ProductType, CartItem, Order, CreateOrderPayload } from '../types';
import { MagnifyingGlassIcon, CameraIcon, XMarkIcon, ShoppingCartIcon, ArrowUpIcon, TagIcon, BookmarkIcon, Squares2X2Icon, BuildingStorefrontIcon } from '@heroicons/vue/24/outline';
import { showAppAlert } from '@/composables/useAppDialog';

defineEmits(['refresh-products']);

const settingsStore = useSettingsStore();
const { settings: storeSetting } = storeToRefs(settingsStore);
const storeContextStore = useStoreContextStore();

const categories = ref<Category[]>([]);
const artists = ref<Artist[]>([]);
const productTypes = ref<ProductType[]>([]);
const products = ref<Product[]>([]);
const isLoading = ref(true);
const catalogError = ref('');

const searchQuery = ref('');
const selectedCategoryId = ref<string | null>(null);
const selectedArtist = ref<string | null>(null);
const selectedProductType = ref<string | null>(null);

const rootCategories = computed(() => categories.value.filter(category => !category.parent_id));
const availableArtists = computed(() => artists.value);
const availableProductTypes = computed(() => productTypes.value);

const categoryLabel = (category: Category): string => {
  const names = [category.name];
  const visited = new Set([category.id]);
  let parentId = category.parent_id ?? null;
  while (parentId) {
    const parent = categories.value.find(item => item.id === parentId);
    if (!parent || visited.has(parent.id)) break;
    names.unshift(parent.name);
    visited.add(parent.id);
    parentId = parent.parent_id ?? null;
  }
  return names.join(' / ');
};

const resetSubFilters = () => {
  selectedArtist.value = null;
  selectedProductType.value = null;
};

const resetAllFilters = () => {
  selectedCategoryId.value = null;
  selectedArtist.value = null;
  selectedProductType.value = null;
};

const selectCategory = (categoryId: string | null): void => {
  selectedCategoryId.value = categoryId;
  resetSubFilters();
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

const fetchCatalogs = async () => {
  catalogError.value = '';
  try {
    const [artistResponse, typeResponse] = await Promise.all([
      api.get('/artists'),
      api.get('/product-types')
    ]);
    if (!Array.isArray(artistResponse.data) || !Array.isArray(typeResponse.data)) {
      throw new Error('Server API belum menyediakan master Artist dan Tipe Produk. Mulai ulang backend dengan versi terbaru.');
    }
    artists.value = artistResponse.data.filter((item: Artist) => item && item.id && item.name?.trim());
    productTypes.value = typeResponse.data.filter((item: ProductType) => item && item.id && item.name?.trim());
  } catch (err: any) {
    artists.value = [];
    productTypes.value = [];
    catalogError.value = err.response?.data?.error || err.message || 'Gagal mengambil master Artist dan Tipe Produk.';
  }
};

const fetchProducts = async () => {
  isLoading.value = true;
  try {
    const params: Record<string, any> = { is_master: false };
    if (storeContextStore.activeStoreId) {
      params.outlet_id = storeContextStore.activeStoreId;
    }
    const res = await api.get('/products', { params });
    products.value = res.data;
  } catch (err: any) {
    console.error('Fetch products error:', err.response?.data?.error || err.message || 'Error occurred');
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  storeContextStore.fetchStores();
  fetchCategories();
  fetchCatalogs();
  fetchProducts();
});

watch(() => storeContextStore.activeStoreId, () => {
  cart.value = [];
  fetchProducts();
});

// Filtering
const filteredProducts = computed(() => {
  return products.value.filter(p => {
    // Filter out master products from POS (POS only displays store POS products)
    if (p.is_master) return false;
    if (storeContextStore.activeStoreId && p.outlet_id && p.outlet_id !== storeContextStore.activeStoreId) {
      return false;
    }

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
      cart.value[existingIndex].quantity = Math.min(cart.value[existingIndex].quantity + 1, product.stock);
    }
  } else {
    cart.value.push({ product, quantity: Math.min(1, product.stock), notes: '' });
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

const setCartQty = ({ index, quantity }: { index: number; quantity: number }): void => {
  const item = cart.value[index];
  if (!item || (item.product.unit !== 'gram' && item.product.unit !== 'liter') || !Number.isFinite(quantity)) return;
  if (quantity <= 0) {
    cart.value.splice(index, 1);
  } else if (quantity <= item.product.stock) {
    item.quantity = Math.round(quantity * 1000) / 1000;
  }
};

const removeCartItem = (index: number): void => {
  cart.value.splice(index, 1);
};

const clearCart = () => {
  cart.value = [];
  discount.value = 0;
};

const totalCartItems = computed(() => new Intl.NumberFormat('id-ID', { maximumFractionDigits: 3 }).format(cart.value.reduce((s, i) => s + i.quantity, 0)));
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

const authStore = useAuthStore();

// Checkout API handler
const handleCheckout = async ({ customer_name, table_number, payment_method, paid_amount, payment_proof }: { customer_name: string; table_number?: string; payment_method: string; paid_amount: number; payment_proof?: string }): Promise<void> => {
  isSubmittingOrder.value = true;
  try {
    const activeUser = authStore.user;
    const cashierName = activeUser?.name || activeUser?.Name || activeUser?.username || activeUser?.Username || 'Kasir';

    const payload: CreateOrderPayload = {
      outlet_id: storeContextStore.activeStoreId || null,
      customer_name,
      table_number,
      cashier_name: cashierName,
      payment_method,
      payment_proof: payment_proof || '',
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
    isMobileCartOpen.value = false;
    isReceiptModalOpen.value = true;
    clearCart();
    fetchProducts(); // refresh stock counts
  } catch (err: any) {
    const errorMsg = err.response?.data?.error || err.message || 'Error occurred';
    await showAppAlert('Gagal memproses transaksi: ' + errorMsg, 'error');
  } finally {
    isSubmittingOrder.value = false;
  }
};

const onReceiptClose = () => {
  isReceiptModalOpen.value = false;
  lastCompletedOrder.value = null;
};
</script>
