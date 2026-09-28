<template>
  <div class="flex flex-col gap-5">
    <!-- Top Header Controls -->
    <div class="flex justify-between items-center p-5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <div class="flex flex-col gap-1">
        <h2 class="text-xl font-bold text-slate-900 dark:text-slate-100">Kelola Produk & Sub-Kategori</h2>
        <p class="text-sm text-slate-500 dark:text-slate-400">Atur produk, artist, dan tipe produk untuk analisa laporan kasir</p>
      </div>

      <div class="flex items-center gap-3">
        <AppButton variant="secondary" class="flex items-center gap-1" @click="isCategoryModalOpen = true">
          <PlusIcon class="w-4 h-4" />
          <span>Kategori Baru</span>
        </AppButton>
        <AppButton variant="primary" class="flex items-center gap-1" @click="openAddModal">
          <PlusIcon class="w-4 h-4" />
          <span>Produk Baru</span>
        </AppButton>
      </div>
    </div>

    <!-- Search & Filter Bar -->
    <div class="flex flex-wrap items-center gap-3 p-4 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <AppInput 
        type="text" 
        class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100 flex-auto min-w-[200px]" 
        v-model="searchQuery" 
        placeholder="Cari produk, SKU, artist, atau tipe..."
      />

      <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100 flex-none min-w-[140px]" v-model="selectedCatId">
        <option :value="null">Semua Kategori</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </select>

      <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100 flex-none min-w-[140px]" v-model="selectedArtist">
        <option value="">Artist</option>
        <option v-for="a in availableArtists" :key="a" :value="a">
          {{ a }}
        </option>
      </select>

      <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100 flex-none min-w-[140px]" v-model="selectedProductType">
        <option value="">Tipe Produk</option>
        <option v-for="t in availableProductTypes" :key="t" :value="t">
          {{ t }}
        </option>
      </select>

      <AppButton variant="secondary" 
        v-if="selectedCatId !== null || selectedArtist || selectedProductType || searchQuery" 
        class="text-xs" 
        @click="resetFilters"
      >
        Reset Filter
      </AppButton>
    </div>

    <!-- Products Table / Grid -->
    <div class="overflow-x-auto p-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <table class="w-full text-left text-sm">
        <thead>
          <tr>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Gambar</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Nama Produk</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Kategori</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Artist</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Tipe Produk</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Harga Jual</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Harga Modal</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Stok</th>
            <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">SKU / Barcode</th>
            <th class="text-right">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoading">
            <td colspan="10" class="text-center">Memuat data produk...</td>
          </tr>
          <tr v-else-if="filteredProducts.length === 0">
            <td colspan="10" class="text-center">Tidak ada produk ditemukan</td>
          </tr>
          <tr v-else v-for="prod in filteredProducts" :key="prod.id">
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
              <img 
                :src="prod.image_url || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400'" 
                class="w-10 h-10 rounded-lg object-cover border border-slate-200 dark:border-slate-700"
              />
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
              <div class="font-semibold text-slate-900 dark:text-slate-100-box">
                <span class="font-semibold text-slate-900 dark:text-slate-100">{{ prod.name }}</span>
              </div>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
              <span class="cat-tag">{{ prod.category ? prod.category.name : '-' }}</span>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
              <span class="sub-tag artist-tag">{{ prod.artist || '-' }}</span>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
              <span class="sub-tag type-tag">{{ prod.product_type || '-' }}</span>
            </td>
            <td class="font-bold price-text">Rp {{ formatPrice(prod.price) }}</td>
            <td class="text-slate-400 dark:text-slate-500">Rp {{ formatPrice(prod.cost_price ?? 0) }}</td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
              <span class="badge" :class="getStockBadge(prod.stock)">
                {{ prod.stock }} unit
              </span>
            </td>
            <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100"><code>{{ prod.barcode || '-' }}</code></td>
            <td class="text-right">
              <div class="action-buttons">
                <AppButton variant="secondary" class="" title="Edit Produk" @click="openEditModal(prod)">
                  <PencilIcon class="w-4 h-4 text-indigo-600" />
                </AppButton>
                <AppButton variant="secondary" class="" title="Hapus Produk" @click="deleteProduct(prod)">
                  <TrashIcon, PencilIcon class="w-4 h-4 text-red-600" />
                </AppButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
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
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppInput from '@/components/ui/AppInput.vue';

import { ref, computed, onMounted } from 'vue';
import api from '@/utils/api';
import type { Category, Product } from '../types';
import { PencilIcon, TrashIcon, PlusIcon, XMarkIcon } from '@heroicons/vue/24/outline';

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


