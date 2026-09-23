<template>
  <div class="products-page">
    <!-- Top Header Controls -->
    <div class="page-header glass-panel">
      <div class="header-title">
        <h2>📦 Kelola Produk & Kategori</h2>
        <p>Tambah, edit, dan atur stok barang kasir</p>
      </div>

      <div class="header-actions">
        <button class="btn btn-secondary" @click="isCategoryModalOpen = true">
          <span>+ Kategori Baru</span>
        </button>
        <button class="btn btn-primary" @click="openAddModal">
          <span>+ Produk Baru</span>
        </button>
      </div>
    </div>

    <!-- Search & Filter Bar -->
    <div class="filter-bar glass-panel">
      <input 
        type="text" 
        class="form-control search-input" 
        v-model="searchQuery" 
        placeholder="Cari produk berdasarkan nama atau SKU..."
      />

      <select class="form-control cat-select" v-model="selectedCatId">
        <option :value="null">Semua Kategori</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </select>
    </div>

    <!-- Products Table / Grid -->
    <div class="table-container glass-panel">
      <table class="data-table">
        <thead>
          <tr>
            <th>Gambar</th>
            <th>Nama Produk</th>
            <th>Kategori</th>
            <th>Harga Jual</th>
            <th>Harga Modal</th>
            <th>Stok</th>
            <th>SKU / Barcode</th>
            <th class="text-right">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="isLoading">
            <td colspan="8" class="text-center">Memuat data produk...</td>
          </tr>
          <tr v-else-if="filteredProducts.length === 0">
            <td colspan="8" class="text-center">Tidak ada produk ditemukan</td>
          </tr>
          <tr v-else v-for="prod in filteredProducts" :key="prod.id">
            <td>
              <img 
                :src="prod.image_url || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400'" 
                class="prod-thumb"
              />
            </td>
            <td>
              <div class="prod-name-box">
                <span class="prod-name">{{ prod.name }}</span>
              </div>
            </td>
            <td>
              <span class="cat-tag">{{ prod.category ? prod.category.name : '-' }}</span>
            </td>
            <td class="font-bold price-text">Rp {{ formatPrice(prod.price) }}</td>
            <td class="text-muted">Rp {{ formatPrice(prod.cost_price ?? 0) }}</td>
            <td>
              <span class="badge" :class="getStockBadge(prod.stock)">
                {{ prod.stock }} unit
              </span>
            </td>
            <td><code>{{ prod.barcode || '-' }}</code></td>
            <td class="text-right">
              <div class="action-buttons">
                <button class="btn-icon btn-edit" @click="openEditModal(prod)">✏️</button>
                <button class="btn-icon btn-delete" @click="deleteProduct(prod)">🗑️</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Product Form Modal (Add / Edit) -->
    <div v-if="isProductModalOpen" class="modal-overlay" @click.self="isProductModalOpen = false">
      <div class="modal-content glass-panel">
        <div class="modal-header">
          <h3>{{ editingId ? 'Edit Produk' : 'Tambah Produk Baru' }}</h3>
          <button class="btn-close" @click="isProductModalOpen = false">✕</button>
        </div>

        <form @submit.prevent="saveProduct" class="modal-body">
          <div class="form-group">
            <label class="form-label">Nama Produk *</label>
            <input type="text" class="form-control" v-model="form.name" required placeholder="Contoh: Kopi Susu Aren" />
          </div>

          <div class="form-group">
            <label class="form-label">Kategori *</label>
            <select class="form-control" v-model="form.category_id" required>
              <option value="" disabled>Pilih Kategori</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Harga Jual (Rp) *</label>
              <input type="number" class="form-control" v-model.number="form.price" required min="0" />
            </div>
            <div class="form-group">
              <label class="form-label">Harga Modal (Rp)</label>
              <input type="number" class="form-control" v-model.number="form.cost_price" min="0" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Stok Awal *</label>
              <input type="number" class="form-control" v-model.number="form.stock" required min="0" />
            </div>
            <div class="form-group">
              <label class="form-label">Kode SKU / Barcode</label>
              <input type="text" class="form-control" v-model="form.barcode" placeholder="8991001" />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">URL Gambar (Opsional)</label>
            <input type="url" class="form-control" v-model="form.image_url" placeholder="https://..." />
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="isProductModalOpen = false">Batal</button>
            <button type="submit" class="btn btn-primary" :disabled="isSaving">
              {{ isSaving ? 'Menyimpan...' : 'Simpan Produk' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Category Modal -->
    <div v-if="isCategoryModalOpen" class="modal-overlay" @click.self="isCategoryModalOpen = false">
      <div class="modal-content glass-panel">
        <div class="modal-header">
          <h3>Tambah Kategori Baru</h3>
          <button class="btn-close" @click="isCategoryModalOpen = false">✕</button>
        </div>
        <form @submit.prevent="saveCategory" class="modal-body">
          <div class="form-group">
            <label class="form-label">Nama Kategori</label>
            <input type="text" class="form-control" v-model="catForm.name" required placeholder="Contoh: Snack" />
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="isCategoryModalOpen = false">Batal</button>
            <button type="submit" class="btn btn-primary">Simpan Kategori</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { Category, Product } from '../types';

const products = ref<Product[]>([]);
const categories = ref<Category[]>([]);
const isLoading = ref(true);

const searchQuery = ref('');
const selectedCatId = ref<number | null>(null);

const isProductModalOpen = ref(false);
const isCategoryModalOpen = ref(false);
const isSaving = ref(false);
const editingId = ref<number | null>(null);

const form = ref({
  name: '',
  category_id: '' as string | number,
  price: 0,
  cost_price: 0,
  stock: 0,
  barcode: '',
  image_url: ''
});

const catForm = ref({ name: '' });

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

const fetchProducts = async () => {
  isLoading.value = true;
  try {
    const res = await fetch('/api/products');
    if (res.ok) products.value = await res.json();
  } finally {
    isLoading.value = false;
  }
};

const fetchCategories = async () => {
  try {
    const res = await fetch('/api/categories');
    if (res.ok) categories.value = await res.json();
  } catch (err) {
    console.error(err);
  }
};

onMounted(() => {
  fetchProducts();
  fetchCategories();
});

const filteredProducts = computed(() => {
  return products.value.filter(p => {
    const matchCat = selectedCatId.value === null || p.category_id === selectedCatId.value;
    const q = searchQuery.value.toLowerCase();
    const matchQ = !q || p.name.toLowerCase().includes(q) || (p.barcode && p.barcode.toLowerCase().includes(q));
    return matchCat && matchQ;
  });
});

const getStockBadge = (stock: number): string => {
  if (stock <= 0) return 'badge-danger';
  if (stock <= 10) return 'badge-warning';
  return 'badge-success';
};

const openAddModal = () => {
  editingId.value = null;
  form.value = {
    name: '',
    category_id: categories.value.length > 0 ? categories.value[0].id : '',
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
    const url = editingId.value ? `/api/products/${editingId.value}` : '/api/products';
    const method = editingId.value ? 'PUT' : 'POST';

    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    });

    if (res.ok) {
      isProductModalOpen.value = false;
      fetchProducts();
    } else {
      alert('Gagal menyimpan produk');
    }
  } finally {
    isSaving.value = false;
  }
};

const deleteProduct = async (prod: Product): Promise<void> => {
  if (confirm(`Hapus produk "${prod.name}"?`)) {
    const res = await fetch(`/api/products/${prod.id}`, { method: 'DELETE' });
    if (res.ok) fetchProducts();
  }
};

const saveCategory = async () => {
  if (!catForm.value.name) return;
  const res = await fetch('/api/categories', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(catForm.value)
  });
  if (res.ok) {
    catForm.value.name = '';
    isCategoryModalOpen.value = false;
    fetchCategories();
  }
};
</script>

<style scoped>
.products-page {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
}

.header-title h2 {
  font-size: 1.25rem;
  font-weight: 800;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
}

.filter-bar {
  display: flex;
  gap: 1rem;
  padding: 0.85rem 1.25rem;
}

.search-input {
  flex: 1;
}

.cat-select {
  width: 220px;
}

.table-container {
  overflow-x: auto;
  padding: 0.5rem;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.88rem;
}

.data-table th {
  padding: 0.85rem 1rem;
  font-weight: 700;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-color);
  background: rgba(15, 23, 42, 0.4);
}

.data-table td {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
}

.prod-thumb {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-sm);
  object-fit: cover;
}

.prod-name {
  font-weight: 700;
  color: var(--text-primary);
}

.cat-tag {
  background: rgba(255, 255, 255, 0.08);
  padding: 0.15rem 0.5rem;
  border-radius: 6px;
  font-size: 0.75rem;
}

.price-text {
  color: var(--accent-secondary);
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 0.35rem;
}

.btn-icon {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-color);
  padding: 0.35rem 0.5rem;
  border-radius: 6px;
  cursor: pointer;

}
.btn-icon:hover {
  background: rgba(255, 255, 255, 0.15);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.text-right { text-align: right; }
.text-center { text-align: center; }
</style>
