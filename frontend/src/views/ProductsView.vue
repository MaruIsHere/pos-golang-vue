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
          <label class="text-[0.7rem] font-bold text-slate-400 uppercase tracking-wider">Kategori</label>
          <select class="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200/60 dark:border-slate-700/60 rounded-[14px] text-sm font-bold text-slate-700 dark:text-slate-300 outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer" v-model="selectedCatId">
            <option :value="null">Semua Kategori</option>
            <option v-for="cat in rootCategories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-[0.7rem] font-bold text-slate-400 uppercase tracking-wider">Merk</label>
          <select class="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200/60 dark:border-slate-700/60 rounded-[14px] text-sm font-bold text-slate-700 dark:text-slate-300 outline-none focus:ring-2 focus:ring-indigo-500/20 cursor-pointer" v-model="selectedArtist">
            <option value="">Semua Merk</option>
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
    <div class="flex-1 flex flex-col bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700/60 rounded-[24px] shadow-sm lg:overflow-hidden min-w-0">
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 p-5 border-b border-slate-100 dark:border-slate-700/50 shrink-0 bg-white dark:bg-slate-800 z-10">
        <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
          Daftar Produk 
          <span class="text-indigo-600 bg-indigo-50 dark:bg-indigo-900/40 dark:text-indigo-400 px-2.5 py-0.5 rounded-full text-xs font-black">{{ filteredProducts.length }}</span>
        </h3>
        <div class="flex items-center gap-2 w-full sm:w-auto">
          <AppButton variant="secondary" class="min-h-11 flex-1 sm:flex-none flex items-center justify-center gap-2 rounded-xl border-slate-200/60 font-bold text-slate-700 dark:text-slate-200" @click="openCategoryModal">
            <PlusIcon class="w-4 h-4" />
            <span>Master Data</span>
          </AppButton>
          <AppButton variant="primary" class="min-h-11 flex-1 sm:flex-none flex items-center justify-center gap-2 rounded-xl font-bold shadow-indigo-600/20 shadow-lg active:scale-95" @click="openAddModal">
            <PlusIcon class="w-4 h-4 text-indigo-50" />
            <span>Produk Baru</span>
          </AppButton>
        </div>
      </div>

      <div class="flex-1 lg:overflow-y-auto p-4 sm:p-5 lg:p-6 custom-scrollbar bg-slate-50/50 dark:bg-slate-900/20">
        <div v-if="isLoading" class="flex flex-col items-center justify-center py-24 gap-4 bg-white/40 dark:bg-slate-800/40 rounded-[24px]">
          <div class="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
          <p class="text-sm font-semibold text-slate-500">Memuat katalog...</p>
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
          <p class="text-sm text-slate-500 mt-1 max-w-sm">Produk yang dicari tidak ditemukan.</p>
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
              <div class="flex sm:flex-row flex-col items-center gap-2 mt-auto pt-3 border-t border-slate-100 dark:border-slate-700/60">
                <button class="w-full sm:flex-1 py-1.5 flex justify-center items-center gap-1.5 text-xs font-bold text-indigo-600 bg-indigo-50 hover:bg-indigo-600 hover:text-white dark:bg-indigo-900/30 dark:text-indigo-400 dark:hover:bg-indigo-600 rounded-[10px] transition-colors" @click="openEditModal(prod)">
                  <PencilIcon class="w-3.5 h-3.5" /> <span class="sm:hidden lg:inline">Edit</span>
                </button>
                <button class="w-full sm:flex-1 py-1.5 flex justify-center items-center gap-1.5 text-xs font-bold text-red-600 bg-red-50 hover:bg-red-600 hover:text-white dark:bg-red-900/30 dark:text-red-400 dark:hover:bg-red-600 rounded-[10px] transition-colors" @click="deleteProduct(prod)">
                  <TrashIcon class="w-3.5 h-3.5" /> <span class="sm:hidden lg:inline">Hapus</span>
                </button>
              </div>
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
              <select v-model="form.artist" class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100">
                <option value="">Pilih Merk</option>
                <option v-for="artist in artists" :key="artist.id" :value="artist.name">{{ artist.name }}</option>
                <option v-if="form.artist && !artists.some(artist => artist.name === form.artist)" :value="form.artist">{{ form.artist }}</option>
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

        <form v-if="activeCatalogTab === 'artists'" @submit.prevent="saveArtist" class="modal-body">
          <div class="form-group">
            <label class="form-label">Nama Merk</label>
            <AppInput v-model="artistForm.name" required maxlength="100" placeholder="Contoh: Indofood, Unilever, Nestle, dll" />
          </div>
          <p v-if="categoryError" class="text-sm font-medium text-rose-600 dark:text-rose-400" role="alert">{{ categoryError }}</p>
          <div class="modal-footer"><AppButton variant="primary" type="submit" :disabled="isSavingCategory">Tambah Merk</AppButton></div>
        </form>
        <div v-if="activeCatalogTab === 'artists'" class="px-5 pb-5">
          <h4 class="mb-3 text-sm font-bold text-slate-700 dark:text-slate-200">Merk Terdaftar ({{ artists.length }})</h4>
          <p v-if="artists.length === 0" class="rounded-lg bg-slate-50 px-3 py-4 text-sm text-slate-500 dark:bg-slate-800 dark:text-slate-400">Belum ada merk.</p>
          <ul v-else class="divide-y divide-slate-200 rounded-lg border border-slate-200 dark:divide-slate-700 dark:border-slate-700">
            <li v-for="artist in artists" :key="artist.id" class="flex flex-wrap items-center justify-between gap-3 px-3 py-2.5">
              <form v-if="editingArtistId === artist.id" class="flex min-w-[240px] flex-1 items-center gap-2" @submit.prevent="updateArtist(artist.id)">
                <AppInput v-model="editingArtistName" class="min-w-0 flex-1" required maxlength="100" aria-label="Nama merk" />
                <AppButton variant="primary" type="submit" size="sm" :disabled="isSavingCategory">Simpan</AppButton>
                <AppButton variant="secondary" type="button" size="sm" @click="editingArtistId = null">Batal</AppButton>
              </form>
              <template v-else>
                <span class="min-w-0 flex-1 truncate text-sm font-semibold text-slate-800 dark:text-slate-100">{{ artist.name }}</span>
                <div class="flex gap-2">
                  <AppButton variant="secondary" size="sm" title="Ubah merk" @click="startArtistEdit(artist)"><PencilIcon class="h-4 w-4" /></AppButton>
                  <AppButton variant="danger" size="sm" title="Hapus merk" @click="deleteArtist(artist)"><TrashIcon class="h-4 w-4" /></AppButton>
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
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppInput from '@/components/ui/AppInput.vue';
import AppBadge from '@/components/ui/AppBadge.vue';

import { ref, computed, onMounted } from 'vue';
import api from '@/utils/api';
import type { Artist, Category, Product, ProductType } from '../types';
import { MagnifyingGlassIcon, PencilIcon, TrashIcon, PlusIcon, XMarkIcon } from '@heroicons/vue/24/outline';
import { showAppAlert, showAppConfirm } from '@/composables/useAppDialog';

const products = ref<Product[]>([]);
const categories = ref<Category[]>([]);
const artists = ref<Artist[]>([]);
const productTypes = ref<ProductType[]>([]);
const isLoading = ref(true);

const searchQuery = ref('');
const selectedCatId = ref<number | null>(null);
const selectedArtist = ref<string>('');
const selectedProductType = ref<string>('');

const isProductModalOpen = ref(false);
const isCategoryModalOpen = ref(false);
const activeCatalogTab = ref<'categories' | 'artists' | 'productTypes'>('categories');
const isSaving = ref(false);
const isCompressingImage = ref(false);
const isSavingCategory = ref(false);
const editingId = ref<number | null>(null);
const editingCategoryId = ref<number | null>(null);
const editingArtistId = ref<number | null>(null);
const editingArtistName = ref('');
const editingProductTypeId = ref<number | null>(null);
const editingProductTypeName = ref('');
const productsError = ref('');
const categoryError = ref('');
const productError = ref('');
const selectedProductImage = ref<File | null>(null);
const productImagePreview = ref('');
const productImageInput = ref<HTMLInputElement | null>(null);

const form = ref({
  name: '',
  category_id: '' as string | number,
  artist: '',
  product_type: '',
  price: 0,
  cost_price: 0,
  stock: 0,
  unit: 'pcs' as Product['unit'],
  barcode: '',
  image_url: ''
});

const categoryForm = ref({ name: '' });
const artistForm = ref({ name: '' });
const productTypeForm = ref({ name: '' });
const catalogTabs = [
  { id: 'categories' as const, label: 'Kategori' },
  { id: 'artists' as const, label: 'Merk' },
  { id: 'productTypes' as const, label: 'Tipe Produk' }
];
const rootCategories = computed(() => categories.value.filter(category => !category.parent_id));
const categoryLabelById = (categoryId: number): string => {
  const category = categories.value.find(item => item.id === categoryId);
  return category?.name ?? 'Umum';
};
const matchesCategoryBranch = (productCategoryId: number, selectedCategoryId: number): boolean => {
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

const availableArtists = computed(() => {
  return artists.value.map(artist => artist.name);
});

const availableProductTypes = computed(() => {
  return productTypes.value.map(productType => productType.name);
});

const normalizeCatalog = <T extends Artist | ProductType>(data: unknown, label: string): T[] => {
  if (!Array.isArray(data)) {
    throw new Error(`API ${label} belum aktif. Restart backend dengan versi terbaru.`);
  }

  const seen = new Set<string>();
  return data.reduce<T[]>((cleaned, item) => {
    if (!item || typeof item.id !== 'number' || typeof item.name !== 'string') return cleaned;
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
    const [artistResponse, typeResponse] = await Promise.all([
      api.get('/artists'),
      api.get('/product-types')
    ]);
    artists.value = normalizeCatalog<Artist>(artistResponse.data, 'Artist');
    productTypes.value = normalizeCatalog<ProductType>(typeResponse.data, 'Tipe Produk');
  } catch (err: any) {
    artists.value = [];
    productTypes.value = [];
    categoryError.value = err.response?.data?.error || err.message || 'Gagal mengambil master artist dan tipe produk.';
  }
};

onMounted(() => {
  fetchProducts();
  fetchCategories();
  fetchCatalogs();
});

const resetFilters = () => {
  searchQuery.value = '';
  selectedCatId.value = null;
  selectedArtist.value = '';
  selectedProductType.value = '';
};

const filteredProducts = computed(() => {
  return products.value.filter(p => {
    const matchCat = selectedCatId.value === null || matchesCategoryBranch(p.category_id, selectedCatId.value);
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

const openCategoryModal = (): void => {
  categoryError.value = '';
  activeCatalogTab.value = 'categories';
  categoryForm.value = { name: '' };
  artistForm.value = { name: '' };
  productTypeForm.value = { name: '' };
  editingCategoryId.value = null;
  editingArtistId.value = null;
  editingProductTypeId.value = null;
  isCategoryModalOpen.value = true;
};

const openAddModal = () => {
  clearProductImage();
  editingId.value = null;
  form.value = {
    name: '',
    category_id: rootCategories.value[0]?.id ?? '',
    artist: '',
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
    category_id: prod.category_id,
    artist: prod.artist ?? '',
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
    const productPayload = { ...form.value };
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

const updateCategory = async (id: number): Promise<void> => {
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

const saveArtist = async (): Promise<void> => {
  categoryError.value = '';
  isSavingCategory.value = true;
  try {
    await api.post('/artists', { name: artistForm.value.name.trim() });
    artistForm.value = { name: '' };
    await fetchCatalogs();
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal menyimpan artist.';
  } finally {
    isSavingCategory.value = false;
  }
};

const startArtistEdit = (artist: Artist): void => {
  categoryError.value = '';
  editingArtistId.value = artist.id;
  editingArtistName.value = artist.name;
};

const updateArtist = async (id: number): Promise<void> => {
  categoryError.value = '';
  isSavingCategory.value = true;
  try {
    await api.put(`/artists/${id}`, { name: editingArtistName.value.trim() });
    editingArtistId.value = null;
    await Promise.all([fetchCatalogs(), fetchProducts()]);
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal memperbarui artist.';
  } finally {
    isSavingCategory.value = false;
  }
};

const deleteArtist = async (artist: Artist): Promise<void> => {
  if (!await showAppConfirm(`Hapus merk "${artist.name}"?`, {
    title: 'Hapus Merk?',
    confirmLabel: 'Ya, Hapus',
    tone: 'danger'
  })) return;
  categoryError.value = '';
  try {
    await api.delete(`/artists/${artist.id}`);
    await fetchCatalogs();
  } catch (err: any) {
    categoryError.value = err.response?.data?.error || err.message || 'Gagal menghapus artist.';
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

const updateProductType = async (id: number): Promise<void> => {
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
</script>
