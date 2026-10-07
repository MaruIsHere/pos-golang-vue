<template>
  <div class="flex flex-col gap-5">
    <div class="flex justify-between items-center p-5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl shadow-sm">
      <div class="flex flex-col gap-1">
        <h2 class="text-xl font-bold text-slate-800 dark:text-slate-100">Master Data Pelanggan & Member</h2>
        <p class="text-sm text-slate-500 dark:text-slate-400">Kelola profil data pelanggan, nomor telepon/WhatsApp, dan poin member toko</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Add / Edit Customer Form -->
      <div class="col-span-1 flex flex-col p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl shadow-sm">
        <div class="border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-800 dark:text-slate-100">
          <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100">
            <template v-if="editingId">
              <span class="flex items-center gap-1.5"><PencilSquareIcon class="w-4 h-4 text-indigo-600" /> Edit Data Pelanggan</span>
            </template>
            <template v-else>
              <span class="flex items-center gap-1.5"><UserPlusIcon class="w-4 h-4 text-indigo-600" /> Tambah Pelanggan Baru</span>
            </template>
          </h3>
        </div>

        <form @submit.prevent="saveCustomer" class="flex flex-col gap-4 mt-2">
          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Nama Pelanggan *</label>
            <AppInput
              type="text"
              class="w-full"
              v-model="custForm.name"
              placeholder="cth: Budi Santoso"
              required
            />
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">No. Telepon / WhatsApp</label>
            <AppInput
              type="text"
              class="w-full"
              v-model="custForm.phone"
              placeholder="cth: 081234567890"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Email</label>
            <AppInput
              type="email"
              class="w-full"
              v-model="custForm.email"
              placeholder="cth: budi@gmail.com"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Alamat Lengkap</label>
            <textarea
              class="w-full text-slate-900 dark:text-slate-100 placeholder:text-slate-500 dark:placeholder:text-slate-400"
              rows="2"
              v-model="custForm.address"
              placeholder="Alamat rumah / kantor pelanggan..."
            ></textarea>
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Poin Loyalty</label>
            <AppInput
              type="number"
              class="w-full"
              v-model.number="custForm.points"
              min="0"
            />
          </div>

          <div class="flex gap-3 mt-4">
            <AppButton variant="secondary"
              v-if="editingId"
              type="button"
              class=""
              @click="resetForm"
            >
              Batal Edit
            </AppButton>

            <AppButton variant="primary" type="submit" class="" :disabled="isSaving">
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
      <div class="col-span-1 lg:col-span-2 flex flex-col p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl shadow-sm overflow-hidden">
        <div class="border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-800 dark:text-slate-100">
          <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100">Daftar Master Pelanggan ({{ filteredCustomers.length }})</h3>

          <div class="flex items-center gap-3 mb-5">
            <MagnifyingGlassIcon class="w-4 h-4 text-slate-400" />
            <AppInput
              type="text"
              class="flex-auto min-w-[200px]"
              v-model="searchQuery"
              placeholder="Cari nama / HP..."
            />
          </div>
        </div>

        <div class="flex flex-col gap-4 mt-2">
          <div v-if="isLoading" class="p-12 text-center text-slate-500">
            <div class="w-8 h-8 border-4 border-slate-200 dark:border-slate-700 border-t-indigo-600 rounded-full animate-spin mx-auto mb-2"></div>
            <p class="text-sm text-slate-500 dark:text-slate-400">Memuat data pelanggan...</p>
          </div>

          <div v-else-if="filteredCustomers.length === 0" class="p-12 text-center text-slate-500">
            <p class="text-sm text-slate-500 dark:text-slate-400">Belum ada data pelanggan yang tersimpan.</p>
          </div>

          <div v-else class="overflow-x-auto p-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl">
            <table class="data-table min-w-[780px] table-fixed text-sm">
              <thead>
                <tr>
                  <th class="w-[18%]">Nama</th>
                  <th class="w-[15%]">No. Telepon</th>
                  <th class="w-[21%]">Email</th>
                  <th class="w-[18%]">Alamat</th>
                  <th class="w-[13%]">Poin</th>
                  <th class="w-[15%] text-right">Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in filteredCustomers" :key="c.id">
                  <td class="font-semibold text-slate-800 dark:text-slate-100">
                    <strong class="text-slate-800 dark:text-slate-100 font-bold">{{ c.name }}</strong>
                  </td>
                  <td class="text-slate-700 dark:text-slate-200">{{ c.phone ||"-" }}</td>
                  <td class="text-slate-700 dark:text-slate-200">{{ c.email ||"-" }}</td>
                  <td class="max-w-[150px] truncate text-slate-600 dark:text-slate-300">{{ c.address ||"-" }}</td>
                  <td class="text-slate-800 dark:text-slate-100">
                    <span class="px-2.5 py-1 bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-400 rounded-full font-bold text-[10px] uppercase tracking-wider border border-amber-200 dark:border-amber-800/50">{{ c.points }} Poin</span>
                  </td>
                  <td class="text-right">
                    <div class="flex items-center justify-end gap-2">
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
import { showAppAlert, showAppConfirm } from '@/composables/useAppDialog';
import type { Customer } from"../types";
import { PencilSquareIcon, TrashIcon, UserPlusIcon, MagnifyingGlassIcon } from"@heroicons/vue/24/outline";

const customers = ref<Customer[]>([]);
const isLoading = ref(true);
const isSaving = ref(false);
const searchQuery = ref("");

const editingId = ref<string | null>(null);
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
    await showAppAlert("Koneksi error: " + errMsg, 'error');
  } finally {
    isSaving.value = false;
  }
};

const deleteCustomer = async (id: string, name: string): Promise<void> => {
  if (!await showAppConfirm(`Apakah Anda yakin ingin menghapus pelanggan '${name}'?`, {
    title: 'Hapus Pelanggan?',
    confirmLabel: 'Ya, Hapus',
    tone: 'danger'
  }))
    return;

  try {
    await api.delete(`/customers/${id}`);
    loadCustomers();
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    await showAppAlert("Koneksi error: " + errMsg, 'error');
  }
};
</script>

