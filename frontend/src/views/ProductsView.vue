<template>
  <div class="flex flex-col lg:flex-row gap-5 h-full lg:max-h-[calc(100vh-100px)]">
    
    <!-- KARTU KIRI: SIDEBAR FILTER -->
    <div class="w-full lg:w-[320px] shrink-0 flex flex-col gap-5 bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700/60 rounded-[24px] shadow-sm p-5 lg:overflow-y-auto custom-scrollbar">
      <div class="flex flex-col gap-1.5">
        <h2 class="text-xl font-black text-slate-800 dark:text-slate-100 tracking-tight">Katalog Produk</h2>
        <p class="text-[0.85rem] font-medium text-slate-500 dark:text-slate-400">Pencarian & filter etalase kasir.</p>
      </div>
      <div class="h-px w-full bg-slate-100 dark:bg-slate-700/50"></div>
      
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-2">
          <label class="text-[0.7rem] font-bold text-slate-400 uppercase tracking-wider">Pencarian</label>
          <div class="relative w-full">
            <AppInput type="text" class="w-full pl-10 pr-3 py-2.5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200/60 dark:border-slate-700/60 rounded-[14px] text-sm font-bold text-slate-800 dark:text-slate-100 placeholder:text-slate-400 focus:ring-2 focus:ring-indigo-500/20" v-model="searchQuery" placeholder="Nama / SKU..." />
            <MagnifyingGlassIcon class="w-5 h-5 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" />
          </div>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-[0.7rem] font-bold text-slate-400 uppercase tracking-wider">Kategori</label>
          <select class="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200/60 dark:border-slate-700/60 rounded-[14px] text-sm font-bold text-slate-700 dark:text-slate-300 outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer" v-model="selectedCatId">
            <option :value="null">Semua Kategori</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-[0.7rem] font-bold text-slate-400 uppercase tracking-wider">Artist</label>
          <select class="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200/60 dark:border-slate-700/60 rounded-[14px] text-sm font-bold text-slate-700 dark:text-slate-300 outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer" v-model="selectedArtist">
            <option value="">Semua Artist</option>
            <option v-for="a in availableArtists" :key="a" :value="a">{{ a }}</option>
          </select>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-[0.7rem] font-bold text-slate-400 uppercase tracking-wider">Tipe Produk</label>
          <select class="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200/60 dark:border-slate-700/60 rounded-[14px] text-sm font-bold text-slate-700 dark:text-slate-300 outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer" v-model="selectedProductType">
            <option value="">Semua Tipe</option>
            <option v-for="t in availableProductTypes" :key="t" :value="t">{{ t }}</option>
          </select>
        </div>
      </div>
    </div>

    <!-- KARTU KANAN: KONTEN PRODUK (SCROLLABLE) -->
    <div class="flex-1 flex flex-col bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700/60 rounded-[24px] shadow-sm overflow-hidden min-w-0">
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 p-5 border-b border-slate-100 dark:border-slate-700/50 shrink-0 bg-white dark:bg-slate-800 z-10">
        <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
          Daftar Produk 
          <span class="text-indigo-600 bg-indigo-50 dark:bg-indigo-900/40 dark:text-indigo-400 px-2.5 py-0.5 rounded-full text-xs font-black">{{ filteredProducts.length }}</span>
        </h3>
        <div class="flex items-center gap-2 w-full sm:w-auto">
          <AppButton variant="secondary" class="flex-1 sm:flex-none flex items-center justify-center gap-2 rounded-xl border-slate-200/60 font-bold text-slate-600 dark:text-slate-300" @click="isCategoryModalOpen = true">
            <PlusIcon class="w-4 h-4" />
            <span>Kategori</span>
          </AppButton>
          <AppButton variant="primary" class="flex-1 sm:flex-none flex items-center justify-center gap-2 rounded-xl font-bold shadow-indigo-600/20 shadow-lg active:scale-95" @click="openAddModal">
            <PlusIcon class="w-4 h-4 text-indigo-50" />
            <span>Produk Baru</span>
          </AppButton>
        </div>
      </div>

      <div class="flex-1 overflow-y-auto p-5 lg:p-6 custom-scrollbar bg-slate-50/50 dark:bg-slate-900/20">
        <div v-if="isLoading" class="flex flex-col items-center justify-center py-24 gap-4 bg-white/40 dark:bg-slate-800/40 rounded-[24px]">
          <div class="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
          <p class="text-sm font-semibold text-slate-500">Memuat katalog...</p>
        </div>
        <div v-else-if="filteredProducts.length === 0" class="flex flex-col items-center justify-center py-24 px-6 text-center bg-white/60 dark:bg-slate-800/40 border-2 border-dashed border-slate-200 dark:border-slate-700 rounded-[24px]">
          <div class="p-4 bg-slate-100 dark:bg-slate-800 rounded-full mb-4">
            <MagnifyingGlassIcon class="w-10 h-10 text-slate-400" />
          </div>
          <h3 class="text-lg font-bold text-slate-800 dark:text-slate-200">Tidak ada produk</h3>
          <p class="text-sm text-slate-500 mt-1 max-w-sm">Produk yang dicari tidak ditemukan.</p>
        </div>
        <div v-else class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-5">
          <div v-for="prod in filteredProducts" :key="prod.id" class="group flex flex-col bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700/60 rounded-[20px] overflow-hidden shadow-sm hover:-translate-y-1 hover:shadow-lg hover:border-indigo-400/50 transition-all h-full">
            <div class="relative w-full h-[150px] shrink-0 p-2 pb-0">
              <div class="w-full h-full overflow-hidden rounded-[14px] bg-slate-100 dark:bg-slate-900 relative">
                <img :src="prod.image_url || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400'" :alt="prod.name" class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110" />
              </div>
              <AppBadge class="absolute top-3 right-3 text-[0.65rem] font-bold shadow-sm backdrop-blur-md" :variant="prod.stock > 0 ? 'info' : 'danger'">
                {{ prod.stock }} unit
              </AppBadge>
            </div>
            <div class="flex flex-col px-4 pt-3 pb-4 flex-1">
              <div class="flex items-center gap-1.5 mb-1.5 shrink-0">
                <span class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider line-clamp-1">{{ prod.category ? prod.category.name : 'Umum' }}</span>
              </div>
              <h3 class="text-sm font-bold text-slate-800 dark:text-slate-100 leading-tight line-clamp-2 h-[2.5rem] shrink-0">{{ prod.name }}</h3>
              <div class="flex flex-col gap-0.5 mt-3 mb-4">
                <span class="text-[0.7rem] font-semibold text-slate-400 line-through">M: Rp {{ formatPrice(prod.cost_price ?? 0) }}</span>
                <span class="text-[0.95rem] font-extrabold text-indigo-600 dark:text-indigo-400">J: Rp {{ formatPrice(prod.price) }}</span>
              </div>
              <div class="flex items-center gap-2 mt-auto pt-3 border-t border-slate-100 dark:border-slate-700/60">
                <button class="flex-1 py-1.5 flex justify-center items-center gap-1.5 text-xs font-bold text-indigo-600 bg-indigo-50 hover:bg-indigo-600 hover:text-white dark:bg-indigo-900/30 dark:text-indigo-400 dark:hover:bg-indigo-600 rounded-[10px] transition-colors" @click="openEditModal(prod)">
                  <PencilIcon class="w-3.5 h-3.5" /> Edit
                </button>
                <button class="flex-1 py-1.5 flex justify-center items-center gap-1.5 text-xs font-bold text-red-600 bg-red-50 hover:bg-red-600 hover:text-white dark:bg-red-900/30 dark:text-red-400 dark:hover:bg-red-600 rounded-[10px] transition-colors" @click="deleteProduct(prod)">
                  <TrashIcon class="w-3.5 h-3.5" /> Hapus
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Product Form Modal (Add / Edit) -->
    <div v-if="isProductModalOpen" class="modal-overlay" @click.self="isProductModalOpen = false">
      <div class="modal-content backdrop-blur-md bg-white/90 dark:bg-slate-900/90 modal-lg">
        <div class="modal-header">
          <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">{{ editingId ? 'Edit Produk & Sub-Kategori' : 'Tambah Produk Baru' }}</h3>
          <AppButton variant="primary" class="" @click="isProductModalOpen = false">
            <XMarkIcon class="w-5 h-5 text-slate-500" />
          </AppButton>
        </div>

        <form @submit.prevent="saveProduct" class="modal-body">
          <div class="form-group">
            <label class="form-label">Nama Produk *</label>
            <AppInput type="text" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="form.name" required placeholder="Contoh: Keyring Chibi Character" />
          </div>

          <div class="form-group">
            <label class="form-label">Kategori Utama *</label>
            <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="form.category_id" required>
              <option value="" disabled>Pilih Kategori</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Artist</label>
              <AppInput type="text" list="artist-list" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="form.artist" placeholder="Contoh: Nama Artist / Kreator" />
              <datalist id="artist-list">
                <option v-for="a in availableArtists" :key="a" :value="a" />
              </datalist>
            </div>
            <div class="form-group">
              <label class="form-label">Tipe Produk</label>
              <AppInput type="text" list="type-list" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="form.product_type" placeholder="Contoh: Art Print, Merchandise, Apparel" />
              <datalist id="type-list">
                <option v-for="t in availableProductTypes" :key="t" :value="t" />
              </datalist>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Harga Jual (Rp) *</label>
              <AppInput type="number" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model.number="form.price" required min="0" />
            </div>
            <div class="form-group">
              <label class="form-label">Harga Modal (Rp)</label>
              <AppInput type="number" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model.number="form.cost_price" min="0" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Stok Awal *</label>
              <AppInput type="number" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model.number="form.stock" required min="0" />
            </div>
            <div class="form-group">
              <label class="form-label">Kode SKU / Barcode</label>
              <AppInput type="text" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="form.barcode" placeholder="8991001" />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">URL Gambar (Opsional)</label>
            <AppInput type="url" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="form.image_url" placeholder="https://..." />
          </div>

          <div class="modal-footer">
            <AppButton variant="secondary" type="button" class="" @click="isProductModalOpen = false">Batal</AppButton>
            <AppButton variant="primary" type="submit" class="" :disabled="isSaving">
              {{ isSaving ? 'Menyimpan...' : 'Simpan Produk' }}
            </AppButton>
          </div>
        </form>
      </div>
    </div>

    <!-- Category Modal -->
    <div v-if="isCategoryModalOpen" class="modal-overlay" @click.self="isCategoryModalOpen = false">
      <div class="modal-content backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
        <div class="modal-header">
          <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Tambah Kategori Utama</h3>
          <AppButton variant="primary" class="" @click="isCategoryModalOpen = false">
            <XMarkIcon class="w-5 h-5 text-slate-500" />
          </AppButton>
        </div>
        <form @submit.prevent="saveCategory" class="modal-body">
          <div class="form-group">
            <label class="form-label">Nama Kategori</label>
            <AppInput type="text" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="catForm.name" required placeholder="Contoh: Merchandise" />
          </div>
          <div class="modal-footer">
            <AppButton variant="secondary" type="button" class="" @click="isCategoryModalOpen = false">Batal</AppButton>
            <AppButton variant="primary" type="submit" class="">Simpan Kategori</AppButton>
          </div>
        </form>
      </div>
    </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppInput from '@/components/ui/AppInput.vue';

import { ref, computed, onMounted } from 'vue';
import api from '@/utils/api';
import type { Category, Product } from '../types';
import { PencilSquareIcon, TrashIcon, PlusIcon, XMarkIcon } from '@heroicons/vue/24/outline';

const products = ref<Product[]>([]);
const categories = ref<Category[]>([]);
const isLoading = ref(true);

const searchQuery = ref('');
const selectedCatId = ref<number | null>(null);
const selectedArtist = ref<string>('');
const selectedProductType = ref<string>('');

const isProductModalOpen = ref(false);
const isCategoryModalOpen = ref(false);
const isSaving = ref(false);
const editingId = ref<number | null>(null);

const form = ref({
  name: '',
  category_id: '' as string | number,
  artist: '',
  product_type: '',
  price: 0,
  cost_price: 0,
  stock: 0,
  barcode: '',
  image_url: ''
});

const catForm = ref({ name: '' });

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

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

const fetchCategories = async () => {
  try {
    const res = await api.get('/categories');
    categories.value = res.data;
  } catch (err: any) {
    console.error('Fetch categories error:', err.response?.data?.error || err.message || 'Error occurred');
  }
};

onMounted(() => {
  fetchProducts();
  fetchCategories();
});

const resetFilters = () => {
  searchQuery.value = '';
  selectedCatId.value = null;
  selectedArtist.value = '';
  selectedProductType.value = '';
};

const filteredProducts = computed(() => {
  return products.value.filter(p => {
    const matchCat = selectedCatId.value === null || p.category_id === selectedCatId.value;
    const matchArt = !selectedArtist.value || p.artist === selectedArtist.value;
    const matchType = !selectedProductType.value || p.product_type === selectedProductType.value;

    const q = searchQuery.value.toLowerCase();
    const matchQ = !q || 
      p.name.toLowerCase().includes(q) || 
      (p.barcode && p.barcode.toLowerCase().includes(q)) ||
      (p.artist && p.artist.toLowerCase().includes(q)) ||
      (p.product_type && p.product_type.toLowerCase().includes(q));

    return matchCat && matchArt && matchType && matchQ;
  });
});

const getStockBadge = (stock: number): string => {
  if (stock <= 0) return 'badge';
  if (stock <= 10) return 'badge-warning';
  return 'badge';
};

const openAddModal = () => {
  editingId.value = null;
  form.value = {
    name: '',
    category_id: categories.value.length > 0 ? categories.value[0].id : '',
    artist: '',
    product_type: '',
    price: 10000,
    cost_price: 5000,
    stock: 20,
    barcode: '',
    image_url: ''
  };
  isProductModalOpen.value = true;
};

const openEditModal = (prod: Product): void => {
  editingId.value = prod.id;
  form.value = {
    name: prod.name,
    category_id: prod.category_id,
    artist: prod.artist ?? '',
    product_type: prod.product_type ?? '',
    price: prod.price,
    cost_price: prod.cost_price ?? 0,
    stock: prod.stock,
    barcode: prod.barcode ?? '',
    image_url: prod.image_url ?? ''
  };
  isProductModalOpen.value = true;
};

const saveProduct = async () => {
  isSaving.value = true;
  try {
    const url = editingId.value ? `/products/${editingId.value}` : '/products';
    
    if (editingId.value) {
      await api.put(url, form.value);
    } else {
      await api.post(url, form.value);
    }
    
    isProductModalOpen.value = false;
    fetchProducts();
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    alert('Gagal menyimpan produk: ' + errMsg);
  } finally {
    isSaving.value = false;
  }
};

const deleteProduct = async (prod: Product): Promise<void> => {
  if (confirm(`Hapus produk"${prod.name}"?`)) {
    try {
      await api.delete(`/products/${prod.id}`);
      fetchProducts();
    } catch (err: any) {
      const errMsg = err.response?.data?.error || err.message || 'Error occurred';
      alert('Gagal menghapus produk: ' + errMsg);
    }
  }
};

const saveCategory = async () => {
  if (!catForm.value.name) return;
  try {
    await api.post('/categories', catForm.value);
    catForm.value.name = '';
    isCategoryModalOpen.value = false;
    fetchCategories();
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    alert('Gagal menyimpan kategori: ' + errMsg);
  }
};
</script>


