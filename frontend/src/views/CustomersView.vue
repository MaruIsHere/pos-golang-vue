<template>
  <div class="flex flex-col gap-5">
    <div class="flex justify-between items-center p-5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <div class="flex flex-col gap-1">
        <h2 class="text-xl font-bold text-slate-900 dark:text-slate-100">Master Data Pelanggan & Member</h2>
        <p class="text-sm text-slate-500 dark:text-slate-400">Kelola profil data pelanggan, nomor telepon/WhatsApp, dan poin member toko</p>
      </div>
    </div>

    <div class="customers-grid">
      <!-- Add / Edit Customer Form -->
      <div class="customers-card backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
        <div class="border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-900 dark:text-slate-100">
          <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">
            <template v-if="editingId">
              <span class="flex items-center gap-1.5"><PencilSquareIcon class="w-4 h-4 text-indigo-600" /> Edit Data Pelanggan</span>
            </template>
            <template v-else>
              <span class="flex items-center gap-1.5"><UserPlusIcon class="w-4 h-4 text-indigo-600" /> Tambah Pelanggan Baru</span>
            </template>
          </h3>
        </div>

        <form @submit.prevent="saveCustomer" class="card-body">
          <div class="form-group">
            <label class="form-label">Nama Pelanggan *</label>
            <AppInput
              type="text"
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100"
              v-model="custForm.name"
              placeholder="cth: Budi Santoso"
              required
            />
          </div>

          <div class="form-group">
            <label class="form-label">No. Telepon / WhatsApp</label>
            <AppInput
              type="text"
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100"
              v-model="custForm.phone"
              placeholder="cth: 081234567890"
            />
          </div>

          <div class="form-group">
            <label class="form-label">Email</label>
            <AppInput
              type="email"
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100"
              v-model="custForm.email"
              placeholder="cth: budi@gmail.com"
            />
          </div>

          <div class="form-group">
            <label class="form-label">Alamat Lengkap</label>
            <textarea
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100"
              rows="2"
              v-model="custForm.address"
              placeholder="Alamat rumah / kantor pelanggan..."
            ></textarea>
          </div>

          <div class="form-group">
            <label class="form-label">Poin Loyalty</label>
            <AppInput
              type="number"
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100"
              v-model.number="custForm.points"
              min="0"
            />
          </div>

          <div class="form-actions">
            <AppButton variant="secondary"
              v-if="editingId"
              type="button"
              class=""
              @click="resetForm"
            >
              Batal Edit
            </AppButton>

            <AppButton variant="success" type="submit" class="" :disabled="isSaving">
              {{
                isSaving
                  ?"Memproses..."
                  : editingId
                    ?"Simpan Perubahan"
                    :"Tambah Pelanggan"
              }}
            </AppButton>
          </div>
        </form>
      </div>

      <!-- Customers Table List -->
      <div class="customers-card backdrop-blur-md bg-white/90 dark:bg-slate-900/90 main-list-card">
        <div class="border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-900 dark:text-slate-100">
          <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Daftar Master Pelanggan ({{ filteredCustomers.length }})</h3>

          <div class="search-box">
            <MagnifyingGlassIcon class="w-4 h-4 text-slate-400" />
            <AppInput
              type="text"
              class="flex-auto min-w-[200px]"
              v-model="searchQuery"
              placeholder="Cari nama / HP..."
            />
          </div>
        </div>

        <div class="card-body">
          <div v-if="isLoading" class="p-12 text-center text-slate-500">
            <div class="w-8 h-8 border-4 border-slate-200 dark:border-slate-700 border-t-indigo-600 rounded-full animate-spin mx-auto mb-2"></div>
            <p class="text-sm text-slate-500 dark:text-slate-400">Memuat data pelanggan...</p>
          </div>

          <div v-else-if="filteredCustomers.length === 0" class="p-12 text-center text-slate-500">
            <p class="text-sm text-slate-500 dark:text-slate-400">Belum ada data pelanggan yang tersimpan.</p>
          </div>

          <div v-else class="overflow-x-auto p-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl">
            <table class="cust-table">
              <thead>
                <tr>
                  <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Nama</th>
                  <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">No. Telepon</th>
                  <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Email</th>
                  <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Alamat</th>
                  <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Poin</th>
                  <th class="text-right">Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in filteredCustomers" :key="c.id">
                  <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
                    <strong class="cust-name">{{ c.name }}</strong>
                  </td>
                  <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">{{ c.phone ||"-" }}</td>
                  <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">{{ c.email ||"-" }}</td>
                  <td class="addr-col">{{ c.address ||"-" }}</td>
                  <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
                    <span class="points-badge">{{ c.points }} Poin</span>
                  </td>
                  <td class="text-right">
                    <div class="action-btns">
                      <AppButton variant="secondary"
                        class=""
                        title="Edit"
                        @click="editCustomer(c)"
                      >
                        <PencilSquareIcon class="w-4 h-4 text-indigo-600" />
                      </AppButton>
                      <AppButton variant="secondary"
                        class=""
                        title="Hapus"
                        @click="deleteCustomer(c.id, c.name)"
                      >
                        <TrashIcon class="w-4 h-4 text-red-600" />
                      </AppButton>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppInput from '@/components/ui/AppInput.vue';

import { ref, computed, onMounted } from"vue";
import api from '@/utils/api';
import type { Customer } from"../types";
import { PencilSquareIcon, TrashIcon, UserPlusIcon, MagnifyingGlassIcon } from"@heroicons/vue/24/outline";

const customers = ref<Customer[]>([]);
const isLoading = ref(true);
const isSaving = ref(false);
const searchQuery = ref("");

const editingId = ref<number | null>(null);
const custForm = ref({
  name:"",
  phone:"",
  email:"",
  address:"",
  points: 0,
});

const loadCustomers = async () => {
  isLoading.value = true;
  try {
    const res = await api.get("/customers");
    customers.value = res.data;
  } catch (err: any) {
    console.error(err.response?.data?.error || err.message || 'Error occurred');
  } finally {
    isLoading.value = false;
  }
};

onMounted(loadCustomers);

const filteredCustomers = computed(() => {
  if (!searchQuery.value) return customers.value;
  const q = searchQuery.value.toLowerCase();
  return customers.value.filter(
    (c) => c.name.toLowerCase().includes(q) || (c.phone && c.phone.includes(q)),
  );
});

const resetForm = () => {
  editingId.value = null;
  custForm.value = { name:"", phone:"", email:"", address:"", points: 0 };
};

const editCustomer = (cust: Customer): void => {
  editingId.value = cust.id;
  custForm.value = {
    name: cust.name,
    phone: cust.phone ||"",
    email: cust.email ||"",
    address: cust.address ||"",
    points: cust.points || 0,
  };
};

const saveCustomer = async () => {
  if (!custForm.value.name) return;

  isSaving.value = true;
  try {
    const url = editingId.value
      ? `/customers/${editingId.value}`
      :"/customers";
    
    if (editingId.value) {
      await api.put(url, custForm.value);
    } else {
      await api.post(url, custForm.value);
    }

    resetForm();
    loadCustomers();
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    alert("Koneksi error:" + errMsg);
  } finally {
    isSaving.value = false;
  }
};

const deleteCustomer = async (id: number, name: string): Promise<void> => {
  if (!confirm(`Apakah Anda yakin ingin menghapus pelanggan '${name}'?`))
    return;

  try {
    await api.delete(`/customers/${id}`);
    loadCustomers();
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    alert("Koneksi error:" + errMsg);
  }
};
</script>


