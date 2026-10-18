<template>
  <div class="flex min-h-0 flex-1 flex-col">
  <div class="flex flex-col lg:flex-row gap-5 h-full lg:max-h-[calc(100vh-100px)]">
    
    <!-- KARTU KIRI: SIDEBAR FILTER -->
    <div class="w-full lg:w-[320px] shrink-0 flex flex-col gap-5 bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700/60 rounded-[24px] shadow-sm p-5 lg:overflow-y-auto custom-scrollbar">
      <div class="flex flex-col gap-1.5">
        <h2 class="text-xl font-black text-slate-800 dark:text-slate-100 tracking-tight">Katalog Produk</h2>
        <p class="text-[0.85rem] font-medium text-slate-500 dark:text-slate-400">Pencarian & filter etalase kasir.</p>
      </div>
      <div class="h-px w-full bg-slate-100 dark:bg-slate-700/50"></div>
      
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-1 gap-4">
        <div class="flex flex-col gap-2">
          <label class="text-[0.7rem] font-bold text-slate-400 uppercase tracking-wider">Pencarian</label>
          <div class="relative w-full">
            <AppInput type="text" class="w-full pl-10 pr-3 py-2.5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200/60 dark:border-slate-700/60 rounded-[14px] text-sm font-bold text-slate-800 dark:text-slate-100 placeholder:text-slate-400 focus:ring-2 focus:ring-indigo-500/20" v-model="searchQuery" placeholder="Nama / SKU..." />
            <MagnifyingGlassIcon class="w-5 h-5 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none" />
          </div>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-[0.7rem] font-bold text-slate-400 uppercase tracking-wider flex items-center gap-1">
            <BuildingStorefrontIcon class="w-3.5 h-3.5" /> Toko / Cabang
          </label>
          <select 
            class="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200/60 dark:border-slate-700/60 rounded-[14px] text-sm font-bold text-slate-700 dark:text-slate-300 outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer"
            v-model="selectedStoreFilter"
          >
            <option :value="null">Semua Toko</option>
            <option v-for="st in storeContextStore.stores" :key="st.id" :value="st.id">
              {{ st.name }} ({{ st.code }})
            </option>
          </select>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-[0.7rem] font-bold text-slate-400 uppercase tracking-wider">Kategori</label>
          <select class="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200/60 dark:border-slate-700/60 rounded-[14px] text-sm font-bold text-slate-700 dark:text-slate-300 outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer" v-model="selectedCatId">
            <option :value="null">Semua Kategori</option>
            <option v-for="cat in rootCategories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-[0.7rem] font-bold text-slate-400 uppercase tracking-wider">Merk</label>
          <select class="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200/60 dark:border-slate-700/60 rounded-[14px] text-sm font-bold text-slate-700 dark:text-slate-300 outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer" v-model="selectedMerk">
            <option value="">Semua Merk</option>
            <option v-for="a in availableMerks" :key="a" :value="a">{{ a }}</option>
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
        <div>
          <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
            <span>Daftar Produk</span>
            <span class="text-indigo-600 bg-indigo-50 dark:bg-indigo-900/40 dark:text-indigo-400 px-2.5 py-0.5 rounded-full text-xs font-black">{{ filteredProducts.length }}</span>
          </h3>
          <p class="text-xs text-slate-500 dark:text-slate-400 mt-1 max-w-md">
            Kelola produk untuk etalase kasir dan inventaris toko Anda.
          </p>
        </div>
        <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-2 w-full sm:w-auto">
          <AppButton variant="secondary" class="w-full sm:w-auto min-h-11 flex items-center justify-center gap-2 rounded-xl border-slate-200/60 font-bold text-slate-700 dark:text-slate-200" @click="openImportCopyModal">
            <ArrowDownTrayIcon class="w-4 h-4 text-slate-600 dark:text-slate-300" />
            <span>Impor & Salin</span>
          </AppButton>
          <AppButton variant="secondary" class="w-full sm:w-auto min-h-11 flex items-center justify-center gap-2 rounded-xl border-slate-200/60 font-bold text-slate-700 dark:text-slate-200" @click="openCategoryModal">
            <PlusIcon class="w-4 h-4" />
            <span>Master Data</span>
          </AppButton>
          <AppButton variant="primary" class="w-full sm:w-auto min-h-11 flex items-center justify-center gap-2 rounded-xl font-bold shadow-indigo-600/20 shadow-lg active:scale-95" @click="openAddModal">
            <PlusIcon class="w-4 h-4 text-indigo-50" />
            <span>Tambah Produk</span>
          </AppButton>
        </div>
      </div>

      <div class="flex-1 lg:overflow-y-auto p-4 sm:p-5 lg:p-6 custom-scrollbar bg-slate-50/50 dark:bg-slate-900/20">
        <!-- Loading Skeleton -->
        <div v-if="isLoading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 sm:gap-5">
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
        <div v-else-if="productsError" class="p-6 rounded-xl border border-rose-200 bg-rose-50 text-rose-700 dark:border-rose-900 dark:bg-rose-950/30 dark:text-rose-300" role="alert">
          <p class="font-semibold">Katalog gagal dimuat</p>
          <p class="mt-1 text-sm">{{ productsError }}</p>
          <AppButton variant="secondary" class="mt-4" @click="fetchProducts">Coba Lagi</AppButton>
        </div>
        <div v-else-if="filteredProducts.length === 0" class="flex flex-col items-center justify-center py-24 px-6 text-center bg-white/60 dark:bg-slate-800/40 border-2 border-dashed border-slate-200 dark:border-slate-700 rounded-[24px]">
          <div class="p-4 bg-slate-100 dark:bg-slate-800 rounded-full mb-4">
            <MagnifyingGlassIcon class="w-10 h-10 text-slate-400" />
          </div>
          <h3 class="text-lg font-bold text-slate-800 dark:text-slate-200">Tidak ada produk</h3>
          <p class="text-sm text-slate-500 mt-1 max-w-sm">
            Belum ada produk yang sesuai dengan filter yang dipilih.
          </p>
          <AppButton 
            variant="primary" 
            class="mt-4 flex items-center gap-2 shadow-lg shadow-indigo-600/20 font-bold" 
            @click="openAddModal"
          >
            <PlusIcon class="w-4 h-4" />
            <span>Tambah Produk Baru</span>
          </AppButton>
        </div>
        <div v-else class="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-3 sm:gap-5">
          <div v-for="prod in filteredProducts" :key="prod.id" class="group flex flex-col bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700/60 rounded-[20px] overflow-hidden shadow-sm hover:-translate-y-1 hover:shadow-lg hover:border-indigo-400/50 transition-all h-full">
            <div class="relative w-full h-[120px] sm:h-[150px] shrink-0 p-2 pb-0">
              <div class="w-full h-full overflow-hidden rounded-[14px] bg-slate-100 dark:bg-slate-900 relative">
                <img :src="prod.image_url || 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400'" :alt="prod.name" class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110" />
              </div>
              <AppBadge class="absolute top-3 right-3 text-[0.65rem] font-bold shadow-sm backdrop-blur-md" :variant="prod.stock > 0 ? 'info' : 'danger'">
                {{ formatQuantity(prod.stock) }} {{ unitLabel(prod.unit) }}
              </AppBadge>
              <span v-if="prod.outlet" class="absolute top-3 left-3 text-[0.7rem] font-extrabold px-2.5 py-1 rounded-full bg-slate-900/80 backdrop-blur-md text-white shadow-md flex items-center gap-1">
                <BuildingStorefrontIcon class="w-3 h-3" />
                <span>{{ prod.outlet.name }}</span>
              </span>
            </div>
            <div class="flex flex-col px-3 sm:px-4 pt-3 pb-4 flex-1">
              <div class="flex items-center gap-1.5 mb-1.5 shrink-0">
                <span class="text-[10px] font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider line-clamp-1">{{ categoryLabelById(prod.category_id) }}</span>
              </div>
              <h3 class="text-sm font-bold text-slate-800 dark:text-slate-100 leading-tight line-clamp-2 h-[2.5rem] shrink-0">{{ prod.name }}</h3>
              <div class="flex flex-col gap-0.5 mt-2 sm:mt-3 mb-3 sm:mb-4">
                <span class="text-[0.65rem] sm:text-[0.7rem] font-semibold text-slate-400 line-through">M: Rp {{ formatPrice(prod.cost_price ?? 0) }}</span>
                <span class="text-[0.85rem] sm:text-[0.95rem] font-extrabold text-indigo-600 dark:text-indigo-400">J: Rp {{ formatPrice(prod.price) }}</span>
              </div>
              <div class="flex items-center gap-1.5 mt-auto pt-3 border-t border-slate-100 dark:border-slate-700/60">
                <button class="flex-1 py-1.5 flex justify-center items-center gap-1 text-[11px] font-bold text-indigo-600 bg-indigo-50 hover:bg-indigo-600 hover:text-white dark:bg-indigo-900/30 dark:text-indigo-400 dark:hover:bg-indigo-600 rounded-[10px] transition-colors" @click="openEditModal(prod)">
                  <PencilIcon class="w-3.5 h-3.5" /> Edit
                </button>
                <button class="py-1.5 px-2 flex justify-center items-center gap-1 text-[11px] font-bold text-red-600 bg-red-50 hover:bg-red-600 hover:text-white dark:bg-red-900/30 dark:text-red-400 dark:hover:bg-red-600 rounded-[10px] transition-colors" @click="deleteProduct(prod)" title="Hapus">
                  <TrashIcon class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

  <!-- Product Form Modal (Add / Edit) -->
    <div v-if="isProductModalOpen" class="modal-overlay" @click.self="closeProductModal">
      <div class="modal-content backdrop-blur-md bg-white/90 dark:bg-slate-900/90 modal-lg">
        <div class="modal-header">
          <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">{{ editingId ? 'Edit Produk & Sub-Kategori' : 'Tambah Produk Baru' }}</h3>
          <AppButton variant="primary" size="icon" aria-label="Tutup form produk" @click="closeProductModal">
            <XMarkIcon class="w-5 h-5" />
          </AppButton>
        </div>

        <form @submit.prevent="saveProduct" class="modal-body">
          <p v-if="productError" class="rounded-lg border border-rose-200 bg-rose-50 px-3 py-2 text-sm font-medium text-rose-700 dark:border-rose-900 dark:bg-rose-950/30 dark:text-rose-300" role="alert">{{ productError }}</p>
          <div class="form-group" v-if="storeContextStore.stores.length > 0">
            <label class="form-label">Toko / Cabang Penjual *</label>
            <select v-model="form.outlet_id" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" required>
              <option v-for="st in storeContextStore.stores" :key="st.id" :value="st.id">{{ st.name }} ({{ st.code }})</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Nama Produk *</label>
            <AppInput type="text" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="form.name" required placeholder="Contoh: Keyring Chibi Character" />
          </div>

          <div class="form-group">
            <label class="form-label">Kategori Utama *</label>
            <select v-model="form.category_id" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" required>
              <option value="" disabled>{{ categories.length ? 'Pilih Kategori' : 'Buat kategori terlebih dahulu' }}</option>
              <option v-for="cat in rootCategories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
            <button v-if="categories.length === 0" type="button" class="self-start text-sm font-semibold text-indigo-600 hover:text-indigo-700 dark:text-indigo-400" @click="openCategoryModal">Tambah kategori</button>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Merk Produk</label>
              <select v-model="form.merk" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100">
                <option value="">Pilih Merk</option>
                <option v-for="merk in merks" :key="merk.id" :value="merk.name">{{ merk.name }}</option>
                <option v-if="form.merk && !merks.some(merk => merk.name === form.merk)" :value="form.merk">{{ form.merk }}</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">Tipe Produk</label>
              <select v-model="form.product_type" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100">
                <option value="">Pilih Tipe Produk</option>
                <option v-for="type in productTypes" :key="type.id" :value="type.name">{{ type.name }}</option>
                <option v-if="form.product_type && !productTypes.some(type => type.name === form.product_type)" :value="form.product_type">{{ form.product_type }}</option>
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Harga Jual (Rp / {{ unitLabel(form.unit) }}) *</label>
              <AppInput type="number" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model.number="form.price" required min="0" />
            </div>
            <div class="form-group">
              <label class="form-label">Harga Modal (Rp / {{ unitLabel(form.unit) }})</label>
              <AppInput type="number" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model.number="form.cost_price" min="0" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Stok Awal *</label>
              <input type="number" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model.number="form.stock" required min="0" :step="form.unit === 'pcs' ? 1 : 0.001" />
            </div>
            <div class="form-group">
              <label class="form-label">Satuan Stok *</label>
              <select v-model="form.unit" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" required>
                <option value="pcs">Pcs (buah)</option>
                <option value="gram">Gram (gr)</option>
                <option value="liter">Liter (L)</option>
              </select>
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Kode SKU / Barcode</label>
              <AppInput type="text" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="form.barcode" placeholder="8991001" />
            </div>
          </div>

          <div class="form-group">
            <label for="product-image" class="form-label">Gambar Produk (Opsional)</label>
            <div class="flex flex-col gap-3 rounded-xl border border-dashed border-slate-300 bg-slate-50 p-4 dark:border-slate-700 dark:bg-slate-900/50 sm:flex-row sm:items-center">
              <img
                v-if="productImagePreview"
                :src="productImagePreview"
                alt="Pratinjau gambar produk"
                class="h-20 w-20 shrink-0 rounded-xl border border-slate-200 object-cover dark:border-slate-700"
              />
              <div class="flex min-w-0 flex-1 flex-col gap-1">
                <input
                  id="product-image"
                  ref="productImageInput"
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  class="block w-full cursor-pointer text-sm text-slate-600 file:mr-3 file:rounded-lg file:border-0 file:bg-indigo-100 file:px-3 file:py-2 file:text-sm file:font-bold file:text-indigo-700 hover:file:bg-indigo-200 dark:text-slate-300 dark:file:bg-indigo-500/15 dark:file:text-indigo-300 dark:hover:file:bg-indigo-500/25"
                  :disabled="isCompressingImage || isSaving"
                  @change="onProductImageSelected"
                />
                <span class="text-xs text-slate-500 dark:text-slate-400">
                  JPG, PNG, atau WebP · Maksimal 10 MB · Otomatis dikompres menjadi maksimal 300 KB dan 1200 × 1200 px
                </span>
                <span v-if="isCompressingImage" class="text-xs font-semibold text-indigo-600 dark:text-indigo-400">Mengoptimalkan gambar...</span>
              </div>
              <AppButton
                v-if="productImagePreview"
                variant="secondary"
                size="sm"
                type="button"
                :disabled="isCompressingImage || isSaving"
                @click="clearProductImage"
              >
                Hapus Gambar
              </AppButton>
            </div>
          </div>

          <div class="modal-footer">
            <AppButton variant="secondary" type="button" @click="closeProductModal">Batal</AppButton>
            <AppButton variant="primary" type="submit" :disabled="isSaving || isCompressingImage">
              {{ isSaving ? 'Menyimpan...' : 'Simpan Produk' }}
            </AppButton>
          </div>
        </form>
      </div>
    </div>

    <!-- Category Modal -->
    <div v-if="isCategoryModalOpen" class="modal-overlay" @click.self="isCategoryModalOpen = false">
      <div class="modal-content modal-lg backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
        <div class="modal-header">
          <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Master Data Produk</h3>
          <AppButton variant="secondary" type="button" @click="isCategoryModalOpen = false">
            <XMarkIcon class="w-5 h-5 text-slate-500" />
          </AppButton>
        </div>
        <div class="px-5 pt-4">
          <div class="flex flex-wrap gap-2 border-b border-slate-200 dark:border-slate-700" role="tablist" aria-label="Jenis master data">
            <button v-for="tab in catalogTabs" :key="tab.id" type="button" role="tab" :aria-selected="activeCatalogTab === tab.id" class="border-b-2 px-4 py-2 text-sm font-semibold transition-colors" :class="activeCatalogTab === tab.id ? 'border-indigo-600 text-indigo-600 dark:text-indigo-400' : 'border-transparent text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-100'" @click="activeCatalogTab = tab.id">
              {{ tab.label }}
            </button>
          </div>
        </div>

        <form v-if="activeCatalogTab === 'categories'" @submit.prevent="saveCategory" class="modal-body">
          <div class="form-group">
            <label class="form-label">Nama Kategori</label>
            <AppInput type="text" v-model="categoryForm.name" required maxlength="100" placeholder="Contoh: Makanan Utama" />
          </div>
          <p v-if="categoryError" class="text-sm font-medium text-rose-600 dark:text-rose-400" role="alert">{{ categoryError }}</p>
          <div class="modal-footer">
            <AppButton variant="primary" type="submit" :disabled="isSavingCategory">
              {{ isSavingCategory ? 'Menyimpan...' : 'Tambah Kategori' }}
            </AppButton>
          </div>
        </form>
        <div v-if="activeCatalogTab === 'categories'" class="px-5 pb-5">
          <h4 class="mb-3 text-sm font-bold text-slate-700 dark:text-slate-200">Kategori Terdaftar ({{ rootCategories.length }})</h4>
          <p v-if="rootCategories.length === 0" class="rounded-lg bg-slate-50 px-3 py-4 text-sm text-slate-500 dark:bg-slate-800 dark:text-slate-400">Belum ada kategori.</p>
          <ul v-else class="divide-y divide-slate-200 rounded-lg border border-slate-200 dark:divide-slate-700 dark:border-slate-700">
            <li v-for="cat in rootCategories" :key="cat.id" class="flex flex-wrap items-center justify-between gap-3 px-3 py-2.5">
              <form v-if="editingCategoryId === cat.id" class="flex min-w-[240px] flex-1 items-center gap-2" @submit.prevent="updateCategory(cat.id)">
                <AppInput v-model="categoryForm.name" class="min-w-0 flex-1" type="text" required maxlength="100" aria-label="Nama kategori" />
                <AppButton variant="primary" type="submit" size="sm" :disabled="isSavingCategory">Simpan</AppButton>
                <AppButton variant="secondary" type="button" size="sm" @click="editingCategoryId = null">Batal</AppButton>
              </form>
              <template v-else>
                <span class="min-w-0 flex-1 truncate text-sm font-semibold text-slate-800 dark:text-slate-100">{{ cat.name }}</span>
                <div class="flex shrink-0 items-center gap-2">
                  <AppButton variant="secondary" size="sm" title="Ubah kategori" @click="startCategoryEdit(cat)">
                    <PencilIcon class="h-4 w-4" />
                    <span class="sr-only">Ubah kategori</span>
                  </AppButton>
                  <AppButton variant="danger" size="sm" title="Hapus kategori" @click="deleteCategory(cat)">
                    <TrashIcon class="h-4 w-4" />
                    <span class="sr-only">Hapus kategori</span>
                  </AppButton>
                </div>
              </template>
            </li>
          </ul>
        </div>

        <form v-if="activeCatalogTab === 'merks'" @submit.prevent="saveMerk" class="modal-body">
          <div class="form-group">
            <label class="form-label">Nama Merk</label>
            <AppInput v-model="merkForm.name" required maxlength="100" placeholder="Contoh: Indofood, Unilever, Nestle, dll" />
          </div>
          <p v-if="categoryError" class="text-sm font-medium text-rose-600 dark:text-rose-400" role="alert">{{ categoryError }}</p>
          <div class="modal-footer"><AppButton variant="primary" type="submit" :disabled="isSavingCategory">Tambah Merk</AppButton></div>
        </form>
        <div v-if="activeCatalogTab === 'merks'" class="px-5 pb-5">
          <h4 class="mb-3 text-sm font-bold text-slate-700 dark:text-slate-200">Merk Terdaftar ({{ merks.length }})</h4>
          <p v-if="merks.length === 0" class="rounded-lg bg-slate-50 px-3 py-4 text-sm text-slate-500 dark:bg-slate-800 dark:text-slate-400">Belum ada merk.</p>
          <ul v-else class="divide-y divide-slate-200 rounded-lg border border-slate-200 dark:divide-slate-700 dark:border-slate-700">
            <li v-for="merk in merks" :key="merk.id" class="flex flex-wrap items-center justify-between gap-3 px-3 py-2.5">
              <form v-if="editingMerkId === merk.id" class="flex min-w-[240px] flex-1 items-center gap-2" @submit.prevent="updateMerk(merk.id)">
                <AppInput v-model="editingMerkName" class="min-w-0 flex-1" required maxlength="100" aria-label="Nama merk" />
                <AppButton variant="primary" type="submit" size="sm" :disabled="isSavingCategory">Simpan</AppButton>
                <AppButton variant="secondary" type="button" size="sm" @click="editingMerkId = null">Batal</AppButton>
              </form>
              <template v-else>
                <span class="min-w-0 flex-1 truncate text-sm font-semibold text-slate-800 dark:text-slate-100">{{ merk.name }}</span>
                <div class="flex gap-2">
                  <AppButton variant="secondary" size="sm" title="Ubah merk" @click="startMerkEdit(merk)"><PencilIcon class="h-4 w-4" /></AppButton>
                  <AppButton variant="danger" size="sm" title="Hapus merk" @click="deleteMerk(merk)"><TrashIcon class="h-4 w-4" /></AppButton>
                </div>
              </template>
            </li>
          </ul>
        </div>

        <form v-if="activeCatalogTab === 'productTypes'" @submit.prevent="saveProductType" class="modal-body">
          <div class="form-group">
            <label class="form-label">Nama Tipe Produk</label>
            <AppInput v-model="productTypeForm.name" required maxlength="100" placeholder="Contoh: Merchandise" />
          </div>
          <p v-if="categoryError" class="text-sm font-medium text-rose-600 dark:text-rose-400" role="alert">{{ categoryError }}</p>
          <div class="modal-footer"><AppButton variant="primary" type="submit" :disabled="isSavingCategory">Tambah Tipe Produk</AppButton></div>
        </form>
        <div v-if="activeCatalogTab === 'productTypes'" class="px-5 pb-5">
          <h4 class="mb-3 text-sm font-bold text-slate-700 dark:text-slate-200">Tipe Produk Terdaftar ({{ productTypes.length }})</h4>
          <p v-if="productTypes.length === 0" class="rounded-lg bg-slate-50 px-3 py-4 text-sm text-slate-500 dark:bg-slate-800 dark:text-slate-400">Belum ada tipe produk.</p>
          <ul v-else class="divide-y divide-slate-200 rounded-lg border border-slate-200 dark:divide-slate-700 dark:border-slate-700">
            <li v-for="productType in productTypes" :key="productType.id" class="flex flex-wrap items-center justify-between gap-3 px-3 py-2.5">
              <form v-if="editingProductTypeId === productType.id" class="flex min-w-[240px] flex-1 items-center gap-2" @submit.prevent="updateProductType(productType.id)">
                <AppInput v-model="editingProductTypeName" class="min-w-0 flex-1" required maxlength="100" aria-label="Nama tipe produk" />
                <AppButton variant="primary" type="submit" size="sm" :disabled="isSavingCategory">Simpan</AppButton>
                <AppButton variant="secondary" type="button" size="sm" @click="editingProductTypeId = null">Batal</AppButton>
              </form>
              <template v-else>
                <span class="min-w-0 flex-1 truncate text-sm font-semibold text-slate-800 dark:text-slate-100">{{ productType.name }}</span>
                <div class="flex gap-2">
                  <AppButton variant="secondary" size="sm" title="Ubah tipe produk" @click="startProductTypeEdit(productType)"><PencilIcon class="h-4 w-4" /></AppButton>
                  <AppButton variant="danger" size="sm" title="Hapus tipe produk" @click="deleteProductType(productType)"><TrashIcon class="h-4 w-4" /></AppButton>
                </div>
              </template>
            </li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Modal Impor & Salin Produk (Gabungan Opsi A & B) -->
    <div v-if="isImportModalOpen" class="modal-overlay" @click.self="isImportModalOpen = false">
      <div class="modal-content modal-lg backdrop-blur-md bg-white/95 dark:bg-slate-900/95 max-w-4xl w-full">
        <div class="modal-header">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-indigo-50 dark:bg-indigo-950/50 flex items-center justify-center text-indigo-600 dark:text-indigo-400">
              <ArrowDownTrayIcon class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Impor & Salin Produk</h3>
              <p class="text-xs text-slate-500 dark:text-slate-400">Pilih salin dari cabang lain atau unggah spreadsheet Excel/CSV.</p>
            </div>
          </div>
          <AppButton variant="secondary" type="button" @click="isImportModalOpen = false">
            <XMarkIcon class="w-5 h-5 text-slate-500" />
          </AppButton>
        </div>

        <!-- Navigation Tabs -->
        <div class="px-6 pt-3 border-b border-slate-200 dark:border-slate-800 flex gap-4">
          <button
            type="button"
            class="pb-3 text-sm font-bold flex items-center gap-2 border-b-2 transition-all cursor-pointer"
            :class="activeImportTab === 'copy' ? 'border-indigo-600 text-indigo-600 dark:text-indigo-400' : 'border-transparent text-slate-500 hover:text-slate-700 dark:text-slate-400'"
            @click="activeImportTab = 'copy'"
          >
            <DocumentDuplicateIcon class="w-4 h-4" />
            <span>Salin dari Cabang Lain</span>
          </button>
          <button
            type="button"
            class="pb-3 text-sm font-bold flex items-center gap-2 border-b-2 transition-all cursor-pointer"
            :class="activeImportTab === 'excel' ? 'border-indigo-600 text-indigo-600 dark:text-indigo-400' : 'border-transparent text-slate-500 hover:text-slate-700 dark:text-slate-400'"
            @click="activeImportTab = 'excel'"
          >
            <TableCellsIcon class="w-4 h-4" />
            <span>Impor Excel / CSV</span>
          </button>
        </div>

        <!-- TAB 1: SALIN DARI CABANG LAIN (BRANCH CLONE) -->
        <div v-if="activeImportTab === 'copy'" class="p-6 flex flex-col gap-5">
          <div class="p-4 rounded-xl bg-indigo-50/60 dark:bg-indigo-950/30 border border-indigo-100 dark:border-indigo-900/50 flex items-start gap-3">
            <BuildingStorefrontIcon class="w-5 h-5 text-indigo-600 dark:text-indigo-400 shrink-0 mt-0.5" />
            <div class="text-xs leading-relaxed text-indigo-900 dark:text-indigo-200">
              <span class="font-bold">Salin Independen Tanpa Duplikasi:</span> Produk yang disalin akan menjadi item mandiri di toko tujuan dengan stok & harga yang dapat disesuaikan. Jika nama produk sudah ada di toko tujuan, sistem otomatis melewatinya agar database tetap rapi.
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Toko Asal (Sumber Data)</label>
              <select
                v-model="copySourceStoreId"
                class="w-full px-3.5 py-2.5 bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-bold text-slate-800 dark:text-slate-200 focus:ring-2 focus:ring-indigo-500/20"
              >
                <option value="" disabled>-- Pilih Toko Asal --</option>
                <option
                  v-for="st in storeContextStore.stores"
                  :key="st.id"
                  :value="st.id"
                  :disabled="st.id === copyTargetStoreId"
                >
                  {{ st.name }} ({{ st.code }})
                </option>
              </select>
            </div>

            <div class="flex flex-col gap-1.5">
              <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Toko Tujuan (Penerima Salinan)</label>
              <select
                v-model="copyTargetStoreId"
                class="w-full px-3.5 py-2.5 bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-bold text-slate-800 dark:text-slate-200 focus:ring-2 focus:ring-indigo-500/20"
              >
                <option value="" disabled>-- Pilih Toko Tujuan --</option>
                <option
                  v-for="st in storeContextStore.stores"
                  :key="st.id"
                  :value="st.id"
                  :disabled="st.id === copySourceStoreId"
                >
                  {{ st.name }} ({{ st.code }})
                </option>
              </select>
            </div>
          </div>

          <!-- Product Picker List -->
          <div v-if="copySourceStoreId" class="flex flex-col gap-3">
            <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2">
              <div class="flex items-center gap-2">
                <input
                  id="select-all-copy"
                  type="checkbox"
                  :checked="isAllSourceProductsSelected"
                  class="rounded text-indigo-600 focus:ring-indigo-500 h-4 w-4"
                  @change="toggleSelectAllSourceProducts"
                />
                <label for="select-all-copy" class="text-xs font-bold text-slate-700 dark:text-slate-300 cursor-pointer">
                  Pilih Semua ({{ filteredSourceProducts.length }} produk)
                </label>
              </div>

              <div class="relative w-full sm:w-64">
                <AppInput
                  v-model="copySearchQuery"
                  type="text"
                  placeholder="Cari produk toko asal..."
                  class="w-full pl-8 pr-3 py-1.5 text-xs bg-slate-50 dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700"
                />
                <MagnifyingGlassIcon class="w-3.5 h-3.5 text-slate-400 absolute left-2.5 top-1/2 -translate-y-1/2 pointer-events-none" />
              </div>
            </div>

            <!-- List of products -->
            <div class="border border-slate-200 dark:border-slate-700/80 rounded-xl overflow-hidden max-h-60 overflow-y-auto divide-y divide-slate-100 dark:divide-slate-800">
              <div
                v-if="filteredSourceProducts.length === 0"
                class="p-6 text-center text-xs text-slate-400"
              >
                Tidak ada produk di toko asal yang cocok.
              </div>
              <div
                v-for="prod in filteredSourceProducts"
                :key="prod.id"
                class="p-3 flex items-center justify-between gap-3 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors"
              >
                <div class="flex items-center gap-3 min-w-0">
                  <input
                    type="checkbox"
                    :value="prod.id"
                    v-model="selectedCopyProductIds"
                    class="rounded text-indigo-600 focus:ring-indigo-500 h-4 w-4 shrink-0"
                  />
                  <div class="min-w-0 flex flex-col">
                    <span class="text-xs font-bold text-slate-800 dark:text-slate-100 truncate">{{ prod.name }}</span>
                    <span class="text-[0.7rem] text-slate-400 truncate">
                      Kategori: {{ prod.category?.name || 'Umum' }}
                      <template v-if="prod.merk"> · Merk: {{ prod.merk }}</template>
                    </span>
                  </div>
                </div>
                <div class="text-right shrink-0">
                  <span class="text-xs font-bold text-indigo-600 dark:text-indigo-400">Rp {{ formatPrice(prod.price) }}</span>
                  <div class="text-[0.7rem] text-slate-400">Stok: {{ prod.stock }} {{ prod.unit }}</div>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="p-8 text-center text-xs text-slate-400 bg-slate-50 dark:bg-slate-800/40 rounded-xl border border-dashed border-slate-200 dark:border-slate-700">
            Silakan pilih Toko Asal terlebih dahulu untuk melihat daftar produk.
          </div>

          <div class="modal-footer pt-3 border-t border-slate-100 dark:border-slate-800 flex justify-end gap-2">
            <AppButton variant="secondary" type="button" @click="isImportModalOpen = false">Batal</AppButton>
            <AppButton
              variant="primary"
              type="button"
              :disabled="isCopying || selectedCopyProductIds.length === 0 || !copySourceStoreId || !copyTargetStoreId"
              @click="submitCopyProducts"
            >
              {{ isCopying ? 'Menyalin...' : `Salin ${selectedCopyProductIds.length} Produk Terpilih` }}
            </AppButton>
          </div>
        </div>

        <!-- TAB 2: IMPOR DARI FILE EXCEL / CSV -->
        <div v-if="activeImportTab === 'excel'" class="p-6 flex flex-col gap-5">
          <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 p-4 rounded-xl bg-emerald-50/70 dark:bg-emerald-950/30 border border-emerald-100 dark:border-emerald-900/50">
            <div class="flex items-center gap-3">
              <TableCellsIcon class="w-5 h-5 text-emerald-600 dark:text-emerald-400 shrink-0" />
              <div>
                <h4 class="text-xs font-bold text-emerald-900 dark:text-emerald-200">Gunakan Format Template Resmi</h4>
                <p class="text-[0.75rem] text-emerald-700 dark:text-emerald-300">Unduh format spreadsheet siap pakai dengan header yang telah disesuaikan.</p>
              </div>
            </div>
            <AppButton
              variant="secondary"
              size="sm"
              class="shrink-0 flex items-center gap-2 border-emerald-300 dark:border-emerald-700 text-emerald-700 dark:text-emerald-300 hover:bg-emerald-100 dark:hover:bg-emerald-900/40"
              @click="downloadExcelTemplate"
            >
              <ArrowDownTrayIcon class="w-3.5 h-3.5" />
              <span>Unduh Template (.xlsx)</span>
            </AppButton>
          </div>

          <div class="flex flex-col gap-1.5">
            <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Toko Tujuan Impor</label>
            <select
              v-model="importTargetStoreId"
              class="w-full px-3.5 py-2.5 bg-slate-50 dark:bg-slate-800/80 border border-slate-200 dark:border-slate-700 rounded-xl text-sm font-bold text-slate-800 dark:text-slate-200 focus:ring-2 focus:ring-indigo-500/20"
            >
              <option value="" disabled>-- Pilih Toko Tujuan --</option>
              <option v-for="st in storeContextStore.stores" :key="st.id" :value="st.id">
                {{ st.name }} ({{ st.code }})
              </option>
            </select>
          </div>

          <!-- Upload Dropzone -->
          <div
            v-if="parsedImportItems.length === 0"
            class="flex flex-col items-center justify-center p-8 border-2 border-dashed border-slate-300 dark:border-slate-700 rounded-2xl bg-slate-50/50 dark:bg-slate-800/30 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors cursor-pointer text-center"
            @dragover.prevent
            @drop.prevent="handleFileDrop"
            @click="triggerFileInput"
          >
            <input
              ref="excelFileInput"
              type="file"
              accept=".xlsx, .xls, .csv"
              class="hidden"
              @change="handleFileUpload"
            />
            <div class="w-12 h-12 rounded-full bg-indigo-50 dark:bg-indigo-950 flex items-center justify-center text-indigo-600 dark:text-indigo-400 mb-3">
              <DocumentArrowUpIcon class="w-6 h-6" />
            </div>
            <p class="text-sm font-bold text-slate-800 dark:text-slate-200">
              Klik atau geser file Excel / CSV ke sini
            </p>
            <p class="text-xs text-slate-400 mt-1">
              Mendukung format .xlsx, .xls, dan .csv (Maksimal 5 MB)
            </p>
            <span v-if="isParsingFile" class="mt-3 text-xs font-semibold text-indigo-600 animate-pulse">
              Sedang membaca dan memvalidasi file...
            </span>
          </div>

          <!-- Error Alert if any -->
          <div v-if="importFileError" class="p-3 rounded-xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-900 text-rose-700 dark:text-rose-300 text-xs">
            {{ importFileError }}
          </div>

          <!-- Preview Parsed Products Table -->
          <div v-if="parsedImportItems.length > 0" class="flex flex-col gap-3">
            <div class="flex items-center justify-between gap-2">
              <div class="flex items-center gap-2">
                <span class="text-xs font-bold text-slate-700 dark:text-slate-200">File: {{ importFileName }}</span>
                <span class="px-2 py-0.5 rounded-full text-[0.7rem] font-bold bg-indigo-50 text-indigo-600 dark:bg-indigo-950 dark:text-indigo-300">
                  {{ validImportItemsCount }} valid dari {{ parsedImportItems.length }} baris
                </span>
              </div>
              <AppButton variant="secondary" size="sm" @click="resetImportState">
                Ganti File
              </AppButton>
            </div>

            <div class="border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden max-h-64 overflow-y-auto">
              <table class="w-full text-left text-xs">
                <thead class="bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 font-bold sticky top-0">
                  <tr>
                    <th class="p-2.5">Status</th>
                    <th class="p-2.5">Nama Produk</th>
                    <th class="p-2.5">Kategori</th>
                    <th class="p-2.5">Merk</th>
                    <th class="p-2.5 text-right">Harga Jual</th>
                    <th class="p-2.5 text-right">Stok</th>
                    <th class="p-2.5">Satuan</th>
                    <th class="p-2.5">Barcode</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                  <tr
                    v-for="(item, idx) in parsedImportItems"
                    :key="idx"
                    :class="item.isValid ? 'hover:bg-slate-50 dark:hover:bg-slate-800/40' : 'bg-rose-50/50 dark:bg-rose-950/20'"
                  >
                    <td class="p-2.5">
                      <span
                        v-if="item.isValid"
                        class="inline-flex items-center text-[0.65rem] font-bold px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300"
                      >
                        Siap
                      </span>
                      <span
                        v-else
                        :title="item.error"
                        class="inline-flex items-center text-[0.65rem] font-bold px-1.5 py-0.5 rounded bg-rose-100 text-rose-700 dark:bg-rose-950 dark:text-rose-300"
                      >
                        Error
                      </span>
                    </td>
                    <td class="p-2.5 font-bold text-slate-800 dark:text-slate-100">{{ item.name || '(Kosong)' }}</td>
                    <td class="p-2.5 text-slate-500">{{ item.category }}</td>
                    <td class="p-2.5 text-slate-500">{{ item.merk || '-' }}</td>
                    <td class="p-2.5 text-right font-medium">Rp {{ formatPrice(item.price) }}</td>
                    <td class="p-2.5 text-right font-medium">{{ item.stock }}</td>
                    <td class="p-2.5 text-slate-500">{{ item.unit }}</td>
                    <td class="p-2.5 text-slate-400 font-mono text-[0.7rem]">{{ item.barcode || '-' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <div class="modal-footer pt-3 border-t border-slate-100 dark:border-slate-800 flex justify-end gap-2">
            <AppButton variant="secondary" type="button" @click="isImportModalOpen = false">Batal</AppButton>
            <AppButton
              variant="primary"
              type="button"
              :disabled="isImporting || validImportItemsCount === 0 || !importTargetStoreId"
              @click="submitBulkImport"
            >
              {{ isImporting ? 'Mengimpor...' : `Impor ${validImportItemsCount} Produk ke Toko` }}
            </AppButton>
          </div>
        </div>
      </div>
    </div>


  </div>
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppInput from '@/components/ui/AppInput.vue';
import AppBadge from '@/components/ui/AppBadge.vue';

import { ref, computed, onMounted } from 'vue';
import api from '@/utils/api';
import type { Category, Merk, Product, ProductType } from '../types';
import {
  MagnifyingGlassIcon,
  PencilIcon,
  TrashIcon,
  PlusIcon,
  XMarkIcon,
  BuildingStorefrontIcon,
  ArrowDownTrayIcon,
  DocumentDuplicateIcon,
  TableCellsIcon,
  DocumentArrowUpIcon
} from '@heroicons/vue/24/outline';
import * as XLSX from 'xlsx-js-style/dist/xlsx.bundle.js';
import { showAppAlert, showAppConfirm } from '@/composables/useAppDialog';
import { useStoreContextStore } from '@/stores/storeContext';

const storeContextStore = useStoreContextStore();
const selectedStoreFilter = ref<string | null>(null);

const products = ref<Product[]>([]);
const categories = ref<Category[]>([]);
const merks = ref<Merk[]>([]);
const productTypes = ref<ProductType[]>([]);
const isLoading = ref(true);

const searchQuery = ref('');
const selectedCatId = ref<string | null>(null);
const selectedMerk = ref<string>('');
const selectedProductType = ref<string>('');

const isProductModalOpen = ref(false);
const isCategoryModalOpen = ref(false);
const activeCatalogTab = ref<'categories' | 'merks' | 'productTypes'>('categories');
const isSaving = ref(false);
const isCompressingImage = ref(false);
const isSavingCategory = ref(false);
const editingId = ref<string | null>(null);
const editingCategoryId = ref<string | null>(null);
const editingMerkId = ref<string | null>(null);
const editingMerkName = ref('');
const editingProductTypeId = ref<string | null>(null);
const editingProductTypeName = ref('');
const productsError = ref('');
const categoryError = ref('');
const productError = ref('');
const selectedProductImage = ref<File | null>(null);
const productImagePreview = ref('');
const productImageInput = ref<HTMLInputElement | null>(null);

const form = ref({
  name: '',
  outlet_id: '' as string,
  category_id: '' as string,
  merk: '',
  product_type: '',
  price: 10000,
  cost_price: 5000,
  stock: 0,
  unit: 'pcs' as Product['unit'],
  barcode: '',
  image_url: ''
});

const categoryForm = ref({ name: '' });
const merkForm = ref({ name: '' });
const productTypeForm = ref({ name: '' });
const catalogTabs = [
  { id: 'categories' as const, label: 'Kategori' },
  { id: 'merks' as const, label: 'Merk' },
  { id: 'productTypes' as const, label: 'Tipe Produk' }
];
const rootCategories = computed(() => categories.value.filter(category => !category.parent_id));
const categoryLabelById = (categoryId: string | null): string => {
  const category = categories.value.find(item => item.id === categoryId);
  return category?.name ?? 'Umum';
};
const matchesCategoryBranch = (productCategoryId: string | null, selectedCategoryId: string | null): boolean => {
  let category = categories.value.find(item => item.id === productCategoryId);
  while (category) {
    if (category.id === selectedCategoryId) return true;
    category = category.parent_id
      ? categories.value.find(item => item.id === category?.parent_id)
      : undefined;
  }
  return false;
};

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);
const formatQuantity = (val: number): string => new Intl.NumberFormat('id-ID', { maximumFractionDigits: 3 }).format(val || 0);
const unitLabel = (unit: Product['unit']): string => unit === 'gram' ? 'gr' : unit === 'liter' ? 'L' : 'pcs';
const maxProductImageBytes = 300 * 1024;
const maxProductImageSourceBytes = 10 * 1024 * 1024;
const maxProductImageDimension = 1200;

const availableMerks = computed(() => {
  return merks.value.map(merk => merk.name);
});

const availableProductTypes = computed(() => {
  return productTypes.value.map(productType => productType.name);
});

const normalizeCatalog = <T extends Merk | ProductType>(data: unknown, label: string): T[] => {
  if (!Array.isArray(data)) {
    throw new Error(`API ${label} belum aktif. Restart backend dengan versi terbaru.`);
  }

  const seen = new Set<string>();
  return data.reduce<T[]>((cleaned, item) => {
    // id bisa berupa UUID string atau number (backward compat)
    if (!item || (typeof item.id !== 'string' && typeof item.id !== 'number') || typeof item.name !== 'string') return cleaned;
    const name = item.name.trim();
    const key = name.toLocaleLowerCase();
    if (!name || seen.has(key)) return cleaned;
    seen.add(key);
    cleaned.push({ ...item, name } as T);
    return cleaned;
  }, []);
};

const fetchProducts = async () => {
  isLoading.value = true;
  productsError.value = '';
  try {
    const res = await api.get('/products');
    products.value = res.data;
  } catch (err: any) {
    productsError.value = err.response?.data?.error || err.message || 'Gagal mengambil produk.';
  } finally {
    isLoading.value = false;
  }
};

const fetchCategories = async () => {
  try {
    const res = await api.get('/categories');
    categories.value = res.data;
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal mengambil kategori.';
  }
};

const fetchCatalogs = async () => {
  try {
    const [merkResponse, typeResponse] = await Promise.all([
      api.get('/merks'),
      api.get('/product-types')
    ]);
    merks.value = normalizeCatalog<Merk>(merkResponse.data, 'Merk');
    productTypes.value = normalizeCatalog<ProductType>(typeResponse.data, 'Tipe Produk');
  } catch (err: any) {
    merks.value = [];
    productTypes.value = [];
    categoryError.value = err.response?.data?.error || err.message || 'Gagal mengambil master merk dan tipe produk.';
  }
};

onMounted(() => {
  storeContextStore.fetchStores();
  fetchProducts();
  fetchCategories();
  fetchCatalogs();
});

const resetFilters = () => {
  searchQuery.value = '';
  selectedCatId.value = null;
  selectedMerk.value = '';
  selectedProductType.value = '';
  selectedStoreFilter.value = null;
};

const filteredProducts = computed(() => {
  return products.value.filter(p => {
    if (selectedStoreFilter.value && p.outlet_id && p.outlet_id !== selectedStoreFilter.value) {
      return false;
    }

    const matchCat = selectedCatId.value === null || matchesCategoryBranch(p.category_id, selectedCatId.value);
    const matchMerk = !selectedMerk.value || p.merk === selectedMerk.value || p.artist === selectedMerk.value;
    const matchType = !selectedProductType.value || p.product_type === selectedProductType.value;

    const q = searchQuery.value.toLowerCase();
    const matchQ = !q || 
      p.name.toLowerCase().includes(q) || 
      (p.barcode && p.barcode.toLowerCase().includes(q)) ||
      (p.merk && p.merk.toLowerCase().includes(q)) ||
      (p.artist && p.artist.toLowerCase().includes(q)) ||
      (p.product_type && p.product_type.toLowerCase().includes(q));

    return matchCat && matchMerk && matchType && matchQ;
  });
});

const getStockBadge = (stock: number): string => {
  if (stock <= 0) return 'badge';
  if (stock <= 10) return 'badge-warning';
  return 'badge';
};

const openCategoryModal = (): void => {
  categoryError.value = '';
  activeCatalogTab.value = 'categories';
  categoryForm.value = { name: '' };
  merkForm.value = { name: '' };
  productTypeForm.value = { name: '' };
  editingCategoryId.value = null;
  editingMerkId.value = null;
  editingProductTypeId.value = null;
  isCategoryModalOpen.value = true;
};

const openAddModal = () => {
  clearProductImage();
  editingId.value = null;
  form.value = {
    name: '',
    outlet_id: selectedStoreFilter.value || storeContextStore.activeStoreId || (storeContextStore.stores[0]?.id ?? ''),
    category_id: rootCategories.value[0]?.id ?? '',
    merk: '',
    product_type: '',
    price: 10000,
    cost_price: 5000,
    stock: 20,
    unit: 'pcs',
    barcode: '',
    image_url: ''
  };
  isProductModalOpen.value = true;
};

const openEditModal = (prod: Product): void => {
  clearProductImage();
  editingId.value = prod.id;
  form.value = {
    name: prod.name,
    outlet_id: prod.outlet_id ?? (storeContextStore.activeStoreId || (storeContextStore.stores[0]?.id ?? '')),
    category_id: prod.category_id ?? "",
    merk: prod.merk ?? prod.artist ?? '',
    product_type: prod.product_type ?? '',
    price: prod.price,
    cost_price: prod.cost_price ?? 0,
    stock: prod.stock,
    unit: prod.unit ?? 'pcs',
    barcode: prod.barcode ?? '',
    image_url: prod.image_url ?? ''
  };
  productImagePreview.value = form.value.image_url;
  isProductModalOpen.value = true;
};

const createJpegBlob = (canvas: HTMLCanvasElement, quality: number): Promise<Blob> =>
  new Promise((resolve, reject) => {
    canvas.toBlob(blob => {
      if (blob) resolve(blob);
      else reject(new Error('Browser gagal mengompres gambar.'));
    }, 'image/jpeg', quality);
  });

const compressProductImage = async (file: File): Promise<File> => {
  const bitmap = await createImageBitmap(file);
  try {
    if (bitmap.width * bitmap.height > 20_000_000) {
      throw new Error('Resolusi gambar terlalu besar. Pilih gambar maksimal 20 megapiksel.');
    }

    const scale = Math.min(1, maxProductImageDimension / Math.max(bitmap.width, bitmap.height));
    let width = Math.max(1, Math.round(bitmap.width * scale));
    let height = Math.max(1, Math.round(bitmap.height * scale));
    const canvas = document.createElement('canvas');

    while (Math.max(width, height) >= 64) {
      canvas.width = width;
      canvas.height = height;
      const context = canvas.getContext('2d');
      if (!context) throw new Error('Browser tidak mendukung pemrosesan gambar.');

      context.fillStyle = '#ffffff';
      context.fillRect(0, 0, width, height);
      context.drawImage(bitmap, 0, 0, width, height);

      for (let quality = 0.85; quality >= 0.45; quality -= 0.1) {
        const blob = await createJpegBlob(canvas, quality);
        if (blob.size <= maxProductImageBytes) {
          const name = file.name.replace(/\.[^.]+$/, '') || 'gambar-produk';
          return new File([blob], `${name}.jpg`, { type: 'image/jpeg' });
        }
      }

      width = Math.floor(width * 0.8);
      height = Math.floor(height * 0.8);
    }
  } finally {
    bitmap.close();
  }
  throw new Error('Gambar tidak dapat dikompres hingga 300 KB. Pilih gambar lain.');
};

const onProductImageSelected = async (event: Event): Promise<void> => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  productError.value = '';
  if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
    productError.value = 'Format gambar harus JPG, PNG, atau WebP.';
    input.value = '';
    return;
  }
  if (file.size > maxProductImageSourceBytes) {
    productError.value = 'Ukuran gambar asli maksimal 10 MB.';
    input.value = '';
    return;
  }

  isCompressingImage.value = true;
  try {
    const compressedImage = await compressProductImage(file);
    selectedProductImage.value = compressedImage;
    setProductImagePreview(URL.createObjectURL(compressedImage));
  } catch (err) {
    productError.value = err instanceof Error ? err.message : 'Gagal memproses gambar.';
    input.value = '';
  } finally {
    isCompressingImage.value = false;
  }
};

const setProductImagePreview = (url: string): void => {
  if (productImagePreview.value.startsWith('blob:')) URL.revokeObjectURL(productImagePreview.value);
  productImagePreview.value = url;
};

const clearProductImage = (): void => {
  selectedProductImage.value = null;
  form.value.image_url = '';
  if (productImageInput.value) productImageInput.value.value = '';
  setProductImagePreview('');
};

const closeProductModal = (): void => {
  isProductModalOpen.value = false;
  clearProductImage();
};

const saveProduct = async () => {
  isSaving.value = true;
  productError.value = '';
  try {
    const url = editingId.value ? `/products/${editingId.value}` : '/products';
    const productPayload: any = { 
      ...form.value,
      outlet_id: form.value.outlet_id || storeContextStore.activeStoreId || (storeContextStore.stores[0]?.id ?? null)
    };
    if (selectedProductImage.value) {
      const imageForm = new FormData();
      imageForm.append('image', selectedProductImage.value);
      const { data } = await api.post<{ image_url: string }>('/products/images', imageForm, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      productPayload.image_url = data.image_url;
    }
    
    if (editingId.value) {
      await api.put(url, productPayload);
    } else {
      await api.post(url, productPayload);
    }
    
    closeProductModal();
    fetchProducts();
  } catch (err: any) {
    productError.value = err.response?.data?.error || err.message || 'Gagal menyimpan produk.';
  } finally {
    isSaving.value = false;
  }
};

const deleteProduct = async (prod: Product): Promise<void> => {
  if (await showAppConfirm(`Hapus produk "${prod.name}"?`, {
    title: 'Hapus Produk?',
    confirmLabel: 'Ya, Hapus',
    tone: 'danger'
  })) {
    try {
      await api.delete(`/products/${prod.id}`);
      fetchProducts();
    } catch (err: any) {
      const errMsg = err.response?.data?.error || err.message || 'Error occurred';
      await showAppAlert('Gagal menghapus produk: ' + errMsg, 'error');
    }
  }
};

const saveCategory = async () => {
  categoryError.value = '';
  if (!categoryForm.value.name.trim()) return;
  isSavingCategory.value = true;
  try {
    const { data: createdCategory } = await api.post('/categories', { name: categoryForm.value.name.trim() });
    await fetchCategories();
    if (categoryError.value) return;
    form.value.category_id = createdCategory.id;
    categoryForm.value = { name: '' };
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal menyimpan kategori.';
  } finally {
    isSavingCategory.value = false;
  }
};

const startCategoryEdit = (category: Category): void => {
  categoryError.value = '';
  editingCategoryId.value = category.id;
  categoryForm.value = { name: category.name };
};

const updateCategory = async (id: string): Promise<void> => {
  categoryError.value = '';
  isSavingCategory.value = true;
  try {
    await api.put(`/categories/${id}`, { name: categoryForm.value.name.trim() });
    editingCategoryId.value = null;
    await fetchCategories();
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal memperbarui kategori.';
  } finally {
    isSavingCategory.value = false;
  }
};

const deleteCategory = async (category: Category): Promise<void> => {
  if (!await showAppConfirm(`Hapus kategori "${category.name}"?`, {
    title: 'Hapus Kategori?',
    confirmLabel: 'Ya, Hapus',
    tone: 'danger'
  })) return;
  categoryError.value = '';
  try {
    await api.delete(`/categories/${category.id}`);
    if (selectedCatId.value !== null && matchesCategoryBranch(selectedCatId.value, category.id)) selectedCatId.value = null;
    await fetchCategories();
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal menghapus kategori.';
  }
};

const saveMerk = async (): Promise<void> => {
  categoryError.value = '';
  isSavingCategory.value = true;
  try {
    await api.post('/merks', { name: merkForm.value.name.trim() });
    merkForm.value = { name: '' };
    await fetchCatalogs();
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal menyimpan merk.';
  } finally {
    isSavingCategory.value = false;
  }
};

const startMerkEdit = (merk: Merk): void => {
  categoryError.value = '';
  editingMerkId.value = merk.id;
  editingMerkName.value = merk.name;
};

const updateMerk = async (id: string): Promise<void> => {
  categoryError.value = '';
  isSavingCategory.value = true;
  try {
    await api.put(`/merks/${id}`, { name: editingMerkName.value.trim() });
    editingMerkId.value = null;
    await Promise.all([fetchCatalogs(), fetchProducts()]);
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal memperbarui merk.';
  } finally {
    isSavingCategory.value = false;
  }
};

const deleteMerk = async (merk: Merk): Promise<void> => {
  if (!await showAppConfirm(`Hapus merk "${merk.name}"?`, {
    title: 'Hapus Merk?',
    confirmLabel: 'Ya, Hapus',
    tone: 'danger'
  })) return;
  categoryError.value = '';
  try {
    await api.delete(`/merks/${merk.id}`);
    await fetchCatalogs();
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal menghapus merk.';
  }
};

const saveProductType = async (): Promise<void> => {
  categoryError.value = '';
  isSavingCategory.value = true;
  try {
    await api.post('/product-types', { name: productTypeForm.value.name.trim() });
    productTypeForm.value = { name: '' };
    await fetchCatalogs();
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal menyimpan tipe produk.';
  } finally {
    isSavingCategory.value = false;
  }
};

const startProductTypeEdit = (productType: ProductType): void => {
  categoryError.value = '';
  editingProductTypeId.value = productType.id;
  editingProductTypeName.value = productType.name;
};

const updateProductType = async (id: string): Promise<void> => {
  categoryError.value = '';
  isSavingCategory.value = true;
  try {
    await api.put(`/product-types/${id}`, { name: editingProductTypeName.value.trim() });
    editingProductTypeId.value = null;
    await Promise.all([fetchCatalogs(), fetchProducts()]);
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal memperbarui tipe produk.';
  } finally {
    isSavingCategory.value = false;
  }
};

const deleteProductType = async (productType: ProductType): Promise<void> => {
  if (!await showAppConfirm(`Hapus tipe produk "${productType.name}"?`, {
    title: 'Hapus Tipe Produk?',
    confirmLabel: 'Ya, Hapus',
    tone: 'danger'
  })) return;
  categoryError.value = '';
  try {
    await api.delete(`/product-types/${productType.id}`);
    await fetchCatalogs();
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal menghapus tipe produk.';
  }
};

// -------------------------------------------------------------
// LOGIKA MODAL IMPOR & SALIN PRODUK (GABUNGAN OPSI A & B)
// -------------------------------------------------------------
const isImportModalOpen = ref(false);
const activeImportTab = ref<'copy' | 'excel'>('copy');

// TAB 1: Salin Produk Antar-Toko (Branch Clone)
const copySourceStoreId = ref<string>('');
const copyTargetStoreId = ref<string>('');
const copySearchQuery = ref<string>('');
const selectedCopyProductIds = ref<string[]>([]);
const isCopying = ref(false);

const sourceStoreProducts = computed(() => {
  if (!copySourceStoreId.value) return [];
  return products.value.filter(p => p.outlet_id === copySourceStoreId.value);
});

const filteredSourceProducts = computed(() => {
  const query = copySearchQuery.value.trim().toLowerCase();
  if (!query) return sourceStoreProducts.value;
  return sourceStoreProducts.value.filter(p =>
    p.name.toLowerCase().includes(query) ||
    (p.merk && p.merk.toLowerCase().includes(query)) ||
    (p.barcode && p.barcode.toLowerCase().includes(query))
  );
});

const isAllSourceProductsSelected = computed(() => {
  if (filteredSourceProducts.value.length === 0) return false;
  return filteredSourceProducts.value.every(p => selectedCopyProductIds.value.includes(p.id));
});

const toggleSelectAllSourceProducts = () => {
  if (isAllSourceProductsSelected.value) {
    const idsToRemove = new Set(filteredSourceProducts.value.map(p => p.id));
    selectedCopyProductIds.value = selectedCopyProductIds.value.filter(id => !idsToRemove.has(id));
  } else {
    const existing = new Set(selectedCopyProductIds.value);
    filteredSourceProducts.value.forEach(p => existing.add(p.id));
    selectedCopyProductIds.value = Array.from(existing);
  }
};

const submitCopyProducts = async () => {
  if (!copySourceStoreId.value || !copyTargetStoreId.value) {
    await showAppAlert('Pilih toko asal dan toko tujuan terlebih dahulu.', 'warning');
    return;
  }
  if (copySourceStoreId.value === copyTargetStoreId.value) {
    await showAppAlert('Toko asal dan toko tujuan tidak boleh sama.', 'warning');
    return;
  }
  if (selectedCopyProductIds.value.length === 0) {
    await showAppAlert('Pilih minimal satu produk untuk disalin.', 'warning');
    return;
  }

  isCopying.value = true;
  try {
    const payload = {
      source_outlet_id: copySourceStoreId.value,
      target_outlet_id: copyTargetStoreId.value,
      product_ids: selectedCopyProductIds.value
    };

    const res = await api.post('/stores/copy-products', payload);
    const data = res.data;
    let msg = `Berhasil menyalin ${data.copied_count || 0} produk.`;
    if (data.skipped_count > 0) {
      msg += ` (${data.skipped_count} produk dilewati karena sudah ada di toko tujuan)`;
    }
    await showAppAlert(msg, 'success', 'Salin Berhasil');
    isImportModalOpen.value = false;
    resetImportState();
    await fetchProducts();
  } catch (err: any) {
    await showAppAlert(err.response?.data?.error || err.message || 'Gagal menyalin produk', 'error');
  } finally {
    isCopying.value = false;
  }
};

// TAB 2: Impor dari File Excel / CSV
const excelFileInput = ref<HTMLInputElement | null>(null);
const importTargetStoreId = ref<string>('');
const importFileName = ref<string>('');
const importFileError = ref<string>('');
const isParsingFile = ref(false);
const isImporting = ref(false);
const parsedImportItems = ref<Array<{
  rowNum: number;
  name: string;
  category: string;
  merk: string;
  product_type: string;
  price: number;
  cost_price: number;
  stock: number;
  unit: 'pcs' | 'gram' | 'liter';
  barcode: string;
  isValid: boolean;
  error?: string;
}>>([]);

const validImportItemsCount = computed(() => {
  return parsedImportItems.value.filter(item => item.isValid).length;
});

const triggerFileInput = () => {
  excelFileInput.value?.click();
};

const resetImportState = () => {
  importFileName.value = '';
  importFileError.value = '';
  parsedImportItems.value = [];
  selectedCopyProductIds.value = [];
  if (excelFileInput.value) {
    excelFileInput.value.value = '';
  }
};

const openImportCopyModal = () => {
  const currentTarget = selectedStoreFilter.value || storeContextStore.activeStoreId || (storeContextStore.stores[0]?.id) || '';
  importTargetStoreId.value = currentTarget;
  copyTargetStoreId.value = currentTarget;

  const otherStore = storeContextStore.stores.find(s => s.id !== currentTarget);
  copySourceStoreId.value = otherStore ? otherStore.id : '';

  selectedCopyProductIds.value = [];
  copySearchQuery.value = '';
  resetImportState();
  isImportModalOpen.value = true;
};

const downloadExcelTemplate = () => {
  const sampleData = [
    {
      'Nama Produk': 'Nasi Goreng Spesial',
      'Kategori': 'Makanan',
      'Merk': 'Dapur Nusantara',
      'Tipe Produk': 'FnB',
      'Harga Jual': 25000,
      'Harga Modal': 15000,
      'Stok': 50,
      'Satuan': 'pcs',
      'Barcode': '899123456001'
    },
    {
      'Nama Produk': 'Es Kopi Susu Aren',
      'Kategori': 'Kopi',
      'Merk': 'Kopi Lokal',
      'Tipe Produk': 'FnB',
      'Harga Jual': 18000,
      'Harga Modal': 9000,
      'Stok': 100,
      'Satuan': 'pcs',
      'Barcode': '899123456002'
    },
    {
      'Nama Produk': 'Kopi Bubuk Robusta',
      'Kategori': 'Kopi',
      'Merk': 'Kopi Lokal',
      'Tipe Produk': 'Retail',
      'Harga Jual': 65000,
      'Harga Modal': 40000,
      'Stok': 1500,
      'Satuan': 'gram',
      'Barcode': '899123456003'
    }
  ];

  const ws = XLSX.utils.json_to_sheet(sampleData, {
    header: ['Nama Produk', 'Kategori', 'Merk', 'Tipe Produk', 'Harga Jual', 'Harga Modal', 'Stok', 'Satuan', 'Barcode']
  });
  ws['!cols'] = [
    { wch: 30 },
    { wch: 18 },
    { wch: 20 },
    { wch: 15 },
    { wch: 14 },
    { wch: 14 },
    { wch: 12 },
    { wch: 10 },
    { wch: 18 }
  ];

  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, 'Template Produk');
  XLSX.writeFile(wb, 'Template_Impor_Produk_POS.xlsx');
};

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;
  parseExcelFile(file);
};

const handleFileDrop = (event: DragEvent) => {
  const file = event.dataTransfer?.files?.[0];
  if (!file) return;
  parseExcelFile(file);
};

const parseExcelFile = (file: File) => {
  importFileError.value = '';
  importFileName.value = file.name;
  isParsingFile.value = true;

  const reader = new FileReader();
  reader.onload = (e) => {
    try {
      const data = new Uint8Array(e.target?.result as ArrayBuffer);
      const workbook = XLSX.read(data, { type: 'array' });
      const firstSheetName = workbook.SheetNames[0];
      const worksheet = workbook.Sheets[firstSheetName];
      const rawRows: Record<string, any>[] = XLSX.utils.sheet_to_json(worksheet, { defval: '' });

      if (rawRows.length === 0) {
        importFileError.value = 'File Excel kosong atau tidak memiliki data produk.';
        parsedImportItems.value = [];
        return;
      }

      parsedImportItems.value = rawRows.map((row, index) => {
        const name = String(row['Nama Produk'] || row['Nama'] || row['nama'] || row['Product Name'] || row['name'] || '').trim();
        const category = String(row['Kategori'] || row['kategori'] || row['Category'] || row['category'] || 'Umum').trim();
        const merk = String(row['Merk'] || row['merk'] || row['Brand'] || row['brand'] || row['Artist'] || '').trim();
        const productType = String(row['Tipe Produk'] || row['Tipe'] || row['tipe'] || row['Product Type'] || row['product_type'] || '').trim();
        const price = Number(row['Harga Jual'] || row['Harga'] || row['harga'] || row['Price'] || row['price'] || 0);
        const costPrice = Number(row['Harga Modal'] || row['HPP'] || row['Cost Price'] || row['cost_price'] || 0);
        const stock = Number(row['Stok'] || row['stok'] || row['Stock'] || row['stock'] || row['Qty'] || 0);
        let unit = String(row['Satuan'] || row['satuan'] || row['Unit'] || row['unit'] || 'pcs').toLowerCase().trim();
        if (unit !== 'pcs' && unit !== 'gram' && unit !== 'liter') {
          unit = 'pcs';
        }
        const barcode = String(row['Barcode'] || row['barcode'] || row['SKU'] || row['sku'] || '').trim();

        const isValid = name.length > 0 && !isNaN(price) && price >= 0 && !isNaN(stock) && stock >= 0;

        return {
          rowNum: index + 2,
          name,
          category,
          merk,
          product_type: productType,
          price: isNaN(price) ? 0 : price,
          cost_price: isNaN(costPrice) ? 0 : costPrice,
          stock: isNaN(stock) ? 0 : stock,
          unit: unit as 'pcs' | 'gram' | 'liter',
          barcode,
          isValid,
          error: !name ? 'Nama produk wajib diisi' : (isNaN(price) || price < 0) ? 'Harga jual tidak valid' : ''
        };
      });
    } catch (err: any) {
      importFileError.value = 'Gagal membaca file: ' + (err.message || 'Format tidak didukung');
    } finally {
      isParsingFile.value = false;
    }
  };
  reader.readAsArrayBuffer(file);
};

const submitBulkImport = async () => {
  if (!importTargetStoreId.value) {
    await showAppAlert('Pilih toko tujuan terlebih dahulu.', 'warning');
    return;
  }
  const validItems = parsedImportItems.value.filter(i => i.isValid);
  if (validItems.length === 0) {
    await showAppAlert('Tidak ada data produk valid untuk diimpor.', 'warning');
    return;
  }

  isImporting.value = true;
  try {
    const payload = {
      outlet_id: importTargetStoreId.value,
      products: validItems.map(item => ({
        name: item.name,
        category: item.category,
        merk: item.merk,
        product_type: item.product_type,
        price: item.price,
        cost_price: item.cost_price,
        stock: item.stock,
        unit: item.unit,
        barcode: item.barcode
      }))
    };

    const res = await api.post('/products/bulk-import', payload);
    const data = res.data;
    let msg = `Berhasil mengimpor ${data.imported_count || 0} produk.`;
    if (data.skipped_count > 0) {
      msg += ` (${data.skipped_count} produk dilewati karena sudah terdaftar)`;
    }
    await showAppAlert(msg, 'success', 'Impor Sukses');
    isImportModalOpen.value = false;
    resetImportState();
    await Promise.all([fetchProducts(), fetchCatalogs(), fetchCategories()]);
  } catch (err: any) {
    await showAppAlert(err.response?.data?.error || err.message || 'Gagal mengimpor produk', 'error');
  } finally {
    isImporting.value = false;
  }
};
</script>
