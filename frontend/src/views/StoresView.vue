<template>
  <div class="flex flex-col gap-6">
    <!-- Top Header -->
    <AppCard body-class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2">
          <BuildingStorefrontIcon class="w-6 h-6 text-indigo-600 dark:text-indigo-400" />
          <h2 class="text-xl font-bold text-slate-800 dark:text-slate-100">Registrasi & Kelola Toko Kasir</h2>
        </div>
        <p class="text-sm text-slate-500 mt-1">Daftarkan toko baru untuk menggunakan sistem kasir, atur data toko, dan pisahkan kasir & katalog produk per toko</p>
      </div>

      <div class="flex items-center gap-3">
        <AppButton variant="primary" class="flex items-center gap-2" @click="openCreateModal">
          <PlusIcon class="w-4 h-4" />
          <span>Registrasi Toko Baru</span>
        </AppButton>
      </div>
    </AppCard>

    <!-- Store Stats Summary -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-5">
      <AppCard body-class="flex items-center gap-4">
        <div class="p-3.5 rounded-2xl bg-indigo-100 dark:bg-indigo-900/40 text-indigo-600 dark:text-indigo-400">
          <BuildingStorefrontIcon class="w-6 h-6" />
        </div>
        <div>
          <span class="text-xs font-bold text-slate-400 uppercase tracking-wider block">Total Toko Terdaftar</span>
          <h3 class="text-xl font-black text-slate-800 dark:text-slate-100">{{ stores.length }} Toko</h3>
        </div>
      </AppCard>

      <AppCard body-class="flex items-center gap-4">
        <div class="p-3.5 rounded-2xl bg-emerald-100 dark:bg-emerald-900/40 text-emerald-600 dark:text-emerald-400">
          <UsersIcon class="w-6 h-6" />
        </div>
        <div>
          <span class="text-xs font-bold text-slate-400 uppercase tracking-wider block">Total Kasir Aktif</span>
          <h3 class="text-xl font-black text-slate-800 dark:text-slate-100">{{ totalCashiers }} Kasir</h3>
        </div>
      </AppCard>

      <AppCard body-class="flex items-center gap-4">
        <div class="p-3.5 rounded-2xl bg-amber-100 dark:bg-amber-900/40 text-amber-600 dark:text-amber-400">
          <CubeIcon class="w-6 h-6" />
        </div>
        <div>
          <span class="text-xs font-bold text-slate-400 uppercase tracking-wider block">Varian Produk Kasir</span>
          <h3 class="text-xl font-black text-slate-800 dark:text-slate-100">{{ totalProducts }} Produk</h3>
        </div>
      </AppCard>
    </div>

    <!-- Store List Grid -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center p-12 bg-white dark:bg-slate-800 rounded-3xl border border-slate-200 dark:border-slate-700">
      <div class="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
      <span class="text-sm font-semibold text-slate-500 mt-3">Memuat data toko...</span>
    </div>

    <div v-else-if="stores.length === 0" class="flex flex-col items-center justify-center p-12 bg-white dark:bg-slate-800 rounded-3xl border border-dashed border-slate-300 dark:border-slate-700 text-center">
      <BuildingStorefrontIcon class="w-12 h-12 text-slate-400 mb-3" />
      <h3 class="text-base font-bold text-slate-800 dark:text-slate-200">Belum ada toko yang terdaftar</h3>
      <p class="text-xs text-slate-500 mt-1 max-w-md">Daftarkan toko pertama Anda agar kasir dapat terhubung ke toko dan mengelola katalog produk toko masing-masing.</p>
      <AppButton variant="primary" class="mt-4" @click="openCreateModal">Registrasi Toko Pertama</AppButton>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <AppCard 
        v-for="store in stores" 
        :key="store.id"
        class="relative overflow-hidden group hover:border-indigo-400 dark:hover:border-indigo-500 transition-all duration-200"
      >
        <div class="flex items-start justify-between border-b border-slate-100 dark:border-slate-700/60 pb-4 mb-4">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 rounded-2xl bg-indigo-50 dark:bg-indigo-900/40 text-indigo-600 dark:text-indigo-400 flex items-center justify-center font-bold text-lg border border-indigo-100 dark:border-indigo-800/40">
              <BuildingStorefrontIcon class="w-6 h-6" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span class="px-2 py-0.5 text-[10px] font-extrabold bg-indigo-100 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-400 rounded-md border border-indigo-200 dark:border-indigo-800/50 uppercase">
                  {{ store.code }}
                </span>
                <span v-if="storeContextStore.activeStoreId === store.id" class="px-2 py-0.5 text-[10px] font-extrabold bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400 rounded-md border border-emerald-200 dark:border-emerald-800/50">
                  Toko Aktif Saat Ini
                </span>
              </div>
              <h3 class="text-base font-extrabold text-slate-800 dark:text-slate-100 mt-1">{{ store.name }}</h3>
            </div>
          </div>
        </div>

        <div class="flex flex-col gap-2 text-xs text-slate-600 dark:text-slate-300 mb-5">
          <div class="flex items-center gap-2">
            <MapPinIcon class="w-4 h-4 text-slate-400 shrink-0" />
            <span class="truncate">{{ store.address || 'Alamat belum diisi' }}</span>
          </div>
          <div class="flex items-center gap-2">
            <PhoneIcon class="w-4 h-4 text-slate-400 shrink-0" />
            <span>{{ store.phone || 'No. Telp belum diisi' }}</span>
          </div>
        </div>

        <!-- Store Metrics Pill Badges -->
        <div class="grid grid-cols-2 gap-2 p-3 bg-slate-50 dark:bg-slate-900/50 rounded-2xl border border-slate-100 dark:border-slate-700/50 mb-5 text-center">
          <div>
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">Kasir Terhubung</span>
            <span class="text-sm font-extrabold text-indigo-600 dark:text-indigo-400">{{ store.total_cashiers || 0 }} Kasir</span>
          </div>
          <div>
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">Produk Kasir</span>
            <span class="text-sm font-extrabold text-emerald-600 dark:text-emerald-400">{{ store.total_products || 0 }} Produk</span>
          </div>
        </div>

        <!-- Store Actions -->
        <div class="flex items-center justify-between gap-2 pt-2 border-t border-slate-100 dark:border-slate-700/60">
          <button 
            v-if="storeContextStore.activeStoreId !== store.id"
            class="px-3 py-1.5 text-xs font-bold text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-950/50 hover:bg-indigo-100 rounded-xl transition-colors flex items-center gap-1"
            @click="selectActiveStore(store.id)"
          >
            <span>Pilih Toko Ini</span>
          </button>
          <span v-else class="text-xs font-bold text-emerald-600 dark:text-emerald-400 flex items-center gap-1">
            <CheckCircleIcon class="w-4 h-4" /> Selected
          </span>

          <div class="flex items-center gap-1">
            <button class="p-2 text-slate-500 hover:text-slate-800 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-xl transition-colors" title="Edit Toko" @click="openEditModal(store)">
              <PencilSquareIcon class="w-4 h-4" />
            </button>
            <button class="p-2 text-indigo-600 hover:bg-indigo-50 dark:hover:bg-indigo-950/50 rounded-xl transition-colors" title="Tugaskan Kasir" @click="openAssignUserModal(store)">
              <UserPlusIcon class="w-4 h-4" />
            </button>
          </div>
        </div>
      </AppCard>
    </div>

    <!-- Modal Form: Registrasi / Edit Toko Baru -->
    <div v-if="isModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-fade-in" @click.self="isModalOpen = false">
      <div class="relative w-full max-w-lg bg-white dark:bg-slate-800 rounded-3xl shadow-2xl border border-slate-200 dark:border-slate-700 overflow-hidden">
        <div class="flex items-center justify-between p-6 border-b border-slate-100 dark:border-slate-700">
          <div class="flex items-center gap-2.5">
            <BuildingStorefrontIcon class="w-6 h-6 text-indigo-600 dark:text-indigo-400" />
            <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100">
              {{ editingStoreId ? 'Edit Data Toko' : 'Registrasi Toko Baru' }}
            </h3>
          </div>
          <button class="text-slate-400 hover:text-slate-600 p-1 rounded-full" @click="isModalOpen = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveStore" class="p-6 flex flex-col gap-4">
          <div>
            <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1.5">Nama Toko *</label>
            <AppInput v-model="form.name" placeholder="Contoh: Toko Cabang Merdeka" required />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1.5">Kode Toko (Opsional)</label>
            <AppInput v-model="form.code" placeholder="Contoh: STORE-002 (Otomatis jika kosong)" />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1.5">Alamat Toko</label>
            <textarea 
              v-model="form.address" 
              rows="3" 
              placeholder="Jl. Merdeka Utama No. 88, Jakarta"
              class="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-2xl text-xs font-medium outline-none focus:ring-2 focus:ring-indigo-500/20 text-slate-800 dark:text-slate-100"
            ></textarea>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1.5">No. Telepon Toko</label>
            <AppInput v-model="form.phone" placeholder="0812-3456-7890" />
          </div>

          <div v-if="formError" class="p-3 bg-red-50 text-red-600 dark:bg-red-900/30 text-xs font-bold rounded-xl border border-red-200 dark:border-red-800">
            {{ formError }}
          </div>

          <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-700">
            <button type="button" class="px-5 py-2.5 text-xs font-bold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-xl transition-colors" @click="isModalOpen = false">
              Batal
            </button>
            <AppButton type="submit" variant="primary" :disabled="isSubmitting">
              {{ isSubmitting ? 'Menyimpan...' : (editingStoreId ? 'Perbarui Toko' : 'Registrasi Toko') }}
            </AppButton>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal Form: Penugasan Kasir ke Toko -->
    <div v-if="isAssignModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-fade-in" @click.self="isAssignModalOpen = false">
      <div class="relative w-full max-w-md bg-white dark:bg-slate-800 rounded-3xl shadow-2xl border border-slate-200 dark:border-slate-700 overflow-hidden">
        <div class="flex items-center justify-between p-6 border-b border-slate-100 dark:border-slate-700">
          <div class="flex items-center gap-2">
            <UserPlusIcon class="w-6 h-6 text-indigo-600 dark:text-indigo-400" />
            <h3 class="text-base font-bold text-slate-800 dark:text-slate-100">
              Tugaskan Kasir ke {{ targetStoreForAssign?.name }}
            </h3>
          </div>
          <button class="text-slate-400 hover:text-slate-600 p-1 rounded-full" @click="isAssignModalOpen = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveAssignUser" class="p-6 flex flex-col gap-4">
          <div>
            <label class="block text-xs font-bold text-slate-700 dark:text-slate-300 mb-1.5">Pilih Pengguna / Kasir *</label>
            <select 
              v-model="assignForm.user_id" 
              class="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-2xl text-xs font-bold text-slate-800 dark:text-slate-100 outline-none focus:ring-2 focus:ring-indigo-500/20"
              required
            >
              <option :value="null" disabled>-- Pilih Pengguna / Kasir --</option>
              <option v-for="user in usersList" :key="user.id" :value="user.id">
                {{ user.name || user.username }} ({{ user.role.toUpperCase() }}) {{ user.store_id === targetStoreForAssign?.id ? '[Sudah di toko ini]' : '' }}
              </option>
            </select>
          </div>

          <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-700">
            <button type="button" class="px-5 py-2.5 text-xs font-bold text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-xl transition-colors" @click="isAssignModalOpen = false">
              Batal
            </button>
            <AppButton type="submit" variant="primary" :disabled="isSubmitting || !assignForm.user_id">
              Tugaskan ke Toko Ini
            </AppButton>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppCard from '@/components/ui/AppCard.vue';
import AppInput from '@/components/ui/AppInput.vue';
import { ref, computed, onMounted } from 'vue';
import api from '@/utils/api';
import type { Store, UserProfile, Product } from '@/types';
import { useStoreContextStore } from '@/stores/storeContext';
import { 
  BuildingStorefrontIcon, 
  PlusIcon, 
  UsersIcon, 
  CubeIcon, 
  MapPinIcon, 
  PhoneIcon, 
  PencilSquareIcon, 
  UserPlusIcon, 
  CheckCircleIcon,
  XMarkIcon 
} from '@heroicons/vue/24/outline';

const storeContextStore = useStoreContextStore();

const stores = ref<Store[]>([]);
const allProducts = ref<Product[]>([]);
const usersList = ref<UserProfile[]>([]);
const isLoading = ref(true);
const isSubmitting = ref(false);

const isModalOpen = ref(false);
const editingStoreId = ref<number | null>(null);
const formError = ref('');
const form = ref({
  name: '',
  code: '',
  address: '',
  phone: ''
});

const isAssignModalOpen = ref(false);
const targetStoreForAssign = ref<Store | null>(null);
const assignForm = ref<{ user_id: number | null }>({
  user_id: null
});

const totalCashiers = computed(() => stores.value.reduce((sum, s) => sum + (s.total_cashiers || 0), 0));
const totalProducts = computed(() => {
  if (!allProducts.value.length) return 0;
  const storeProds = allProducts.value.filter(p => !p.is_master);
  if (storeProds.length === 0) {
    return allProducts.value.filter(p => p.is_master).length;
  }
  const uniqueKeys = new Set(
    storeProds.map(p => p.master_product_id ? `master_${p.master_product_id}` : `name_${p.name.trim().toLowerCase()}`)
  );
  return uniqueKeys.size;
});

const fetchStores = async () => {
  isLoading.value = true;
  try {
    const [storesRes, productsRes] = await Promise.all([
      api.get('/stores'),
      api.get('/products')
    ]);
    stores.value = storesRes.data;
    allProducts.value = productsRes.data;
    storeContextStore.fetchStores();
  } catch (err: any) {
    console.error('Fetch stores error:', err);
  } finally {
    isLoading.value = false;
  }
};

const fetchUsers = async () => {
  try {
    const res = await api.get('/users');
    usersList.value = res.data;
  } catch (err) {
    console.error('Fetch users error:', err);
  }
};

const openCreateModal = () => {
  editingStoreId.value = null;
  formError.value = '';
  form.value = { name: '', code: '', address: '', phone: '' };
  isModalOpen.value = true;
};

const openEditModal = (store: Store) => {
  editingStoreId.value = store.id;
  formError.value = '';
  form.value = {
    name: store.name,
    code: store.code,
    address: store.address || '',
    phone: store.phone || ''
  };
  isModalOpen.value = true;
};

const saveStore = async () => {
  formError.value = '';
  if (!form.value.name.trim()) {
    formError.value = 'Nama toko wajib diisi';
    return;
  }

  isSubmitting.value = true;
  try {
    if (editingStoreId.value) {
      await api.put(`/stores/${editingStoreId.value}`, form.value);
    } else {
      await api.post('/stores', form.value);
    }
    isModalOpen.value = false;
    await fetchStores();
  } catch (err: any) {
    formError.value = err.response?.data?.error || 'Gagal menyimpan data toko';
  } finally {
    isSubmitting.value = false;
  }
};

const openAssignUserModal = (store: Store) => {
  targetStoreForAssign.value = store;
  assignForm.value.user_id = null;
  isAssignModalOpen.value = true;
  fetchUsers();
};

const saveAssignUser = async () => {
  if (!targetStoreForAssign.value || !assignForm.value.user_id) return;
  isSubmitting.value = true;
  try {
    await api.post('/stores/assign-user', {
      user_id: assignForm.value.user_id,
      store_id: targetStoreForAssign.value.id
    });
    isAssignModalOpen.value = false;
    await fetchStores();
  } catch (err: any) {
    alert(err.response?.data?.error || 'Gagal menugaskan kasir ke toko');
  } finally {
    isSubmitting.value = false;
  }
};

const selectActiveStore = (id: number) => {
  storeContextStore.setActiveStore(id);
};

onMounted(() => {
  fetchStores();
  fetchUsers();
});
</script>
