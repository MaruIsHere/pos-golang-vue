<template>
  <div class="customers-page">
    <div class="page-header glass-panel">
      <div class="header-title">
        <h2>👥 Master Data Pelanggan & Member</h2>
        <p>Kelola profil data pelanggan, nomor telepon/WhatsApp, dan poin member toko</p>
      </div>
    </div>

    <div class="customers-grid">
      <!-- Add / Edit Customer Form -->
      <div class="customers-card glass-panel">
        <div class="card-header">
          <h3>{{ editingId ? '✏️ Edit Data Pelanggan' : '+ Tambah Pelanggan Baru' }}</h3>
        </div>

        <form @submit.prevent="saveCustomer" class="card-body">
          <div class="form-group">
            <label class="form-label">Nama Pelanggan</label>
            <input 
              type="text" 
              class="form-control" 
              v-model="custForm.name" 
              placeholder="cth: Budi Santoso" 
              required 
            />
          </div>

          <div class="form-group">
            <label class="form-label">No. Telepon / WhatsApp</label>
            <input 
              type="text" 
              class="form-control" 
              v-model="custForm.phone" 
              placeholder="cth: 081234567890" 
            />
          </div>

          <div class="form-group">
            <label class="form-label">Email</label>
            <input 
              type="email" 
              class="form-control" 
              v-model="custForm.email" 
              placeholder="cth: budi@gmail.com" 
            />
          </div>

          <div class="form-group">
            <label class="form-label">Alamat Lengkap</label>
            <textarea 
              class="form-control" 
              rows="2" 
              v-model="custForm.address" 
              placeholder="Alamat rumah / kantor pelanggan..."
            ></textarea>
          </div>

          <div class="form-group">
            <label class="form-label">Poin Loyalty</label>
            <input 
              type="number" 
              class="form-control" 
              v-model.number="custForm.points" 
              min="0" 
            />
          </div>

          <div class="form-actions">
            <button v-if="editingId" type="button" class="btn btn-secondary" @click="resetForm">
              Batal Edit
            </button>

            <button type="submit" class="btn btn-success" :disabled="isSaving">
              {{ isSaving ? 'Memproses...' : (editingId ? 'Simpan Perubahan' : '+ Tambah Pelanggan') }}
            </button>
          </div>
        </form>
      </div>

      <!-- Customers Table List -->
      <div class="customers-card glass-panel main-list-card">
        <div class="card-header">
          <h3>📋 Daftar Master Pelanggan ({{ filteredCustomers.length }})</h3>
          
          <div class="search-box">
            <span class="search-icon">🔍</span>
            <input 
              type="text" 
              class="search-input" 
              v-model="searchQuery" 
              placeholder="Cari nama / HP..." 
            />
          </div>
        </div>

        <div class="card-body">
          <div v-if="isLoading" class="loading-box">
            <div class="spinner"></div>
            <p>Memuat data pelanggan...</p>
          </div>

          <div v-else-if="filteredCustomers.length === 0" class="empty-box">
            <span class="empty-icon">👥</span>
            <p>Belum ada data pelanggan yang tersimpan.</p>
          </div>

          <div v-else class="table-wrapper">
            <table class="cust-table">
              <thead>
                <tr>
                  <th>Nama</th>
                  <th>No. Telepon</th>
                  <th>Email</th>
                  <th>Alamat</th>
                  <th>Poin</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in filteredCustomers" :key="c.id">
                  <td>
                    <strong class="cust-name">👤 {{ c.name }}</strong>
                  </td>
                  <td>{{ c.phone || '-' }}</td>
                  <td>{{ c.email || '-' }}</td>
                  <td class="addr-col">{{ c.address || '-' }}</td>
                  <td>
                    <span class="points-badge">⭐ {{ c.points }} Poin</span>
                  </td>
                  <td>
                    <div class="action-btns">
                      <button class="btn-edit" title="Edit" @click="editCustomer(c)">
                        ✏️
                      </button>
                      <button class="btn-delete" title="Hapus" @click="deleteCustomer(c.id, c.name)">
                        🗑️
                      </button>
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
import { ref, computed, onMounted } from 'vue';
import type { Customer } from '../types';

const customers = ref<Customer[]>([]);
const isLoading = ref(true);
const isSaving = ref(false);
const searchQuery = ref('');

const editingId = ref<number | null>(null);
const custForm = ref({
  name: '',
  phone: '',
  email: '',
  address: '',
  points: 0
});

const loadCustomers = async () => {
  isLoading.value = true;
  try {
    const res = await fetch('/api/customers');
    if (res.ok) customers.value = await res.json();
  } catch (err) {
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(loadCustomers);

const filteredCustomers = computed(() => {
  if (!searchQuery.value) return customers.value;
  const q = searchQuery.value.toLowerCase();
  return customers.value.filter(c => 
    c.name.toLowerCase().includes(q) || (c.phone && c.phone.includes(q))
  );
});

const resetForm = () => {
  editingId.value = null;
  custForm.value = { name: '', phone: '', email: '', address: '', points: 0 };
};

const editCustomer = (cust: Customer): void => {
  editingId.value = cust.id;
  custForm.value = {
    name: cust.name,
    phone: cust.phone || '',
    email: cust.email || '',
    address: cust.address || '',
    points: cust.points || 0
  };
};

const saveCustomer = async () => {
  if (!custForm.value.name) return;

  isSaving.value = true;
  try {
    const url = editingId.value ? `/api/customers/${editingId.value}` : '/api/customers';
    const method = editingId.value ? 'PUT' : 'POST';

    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(custForm.value)
    });

    if (res.ok) {
      alert(`Data pelanggan berhasil ${editingId.value ? 'diperbarui' : 'ditambahkan'}!`);
      resetForm();
      loadCustomers();
    } else {
      const errData = await res.json();
      alert('Gagal: ' + (errData.error || 'Terjadi kesalahan'));
    }
  } catch (err) {
    alert('Koneksi error: ' + (err as Error).message);
  } finally {
    isSaving.value = false;
  }
};

const deleteCustomer = async (id: number, name: string): Promise<void> => {
  if (!confirm(`Apakah Anda yakin ingin menghapus pelanggan '${name}'?`)) return;

  try {
    const res = await fetch(`/api/customers/${id}`, { method: 'DELETE' });
    if (res.ok) {
      loadCustomers();
    } else {
      alert('Gagal menghapus pelanggan');
    }
  } catch (err) {
    alert('Koneksi error: ' + (err as Error).message);
  }
};
</script>

<style scoped>
.customers-page {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.page-header {
  padding: 1.25rem 1.5rem;
}

.header-title h2 {
  font-size: 1.25rem;
  font-weight: 800;
}

.customers-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.25rem;
}

@media (min-width: 1024px) {
  .customers-grid {
    grid-template-columns: 360px 1fr;
  }
}

.customers-card {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.85rem;
}

.card-header h3 {
  font-size: 1.1rem;
  font-weight: 800;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 0.25rem 0.6rem;
}

.search-input {
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-size: 0.8rem;
  outline: none;
  width: 130px;
}

.form-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.table-wrapper {
  overflow-x: auto;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
}

.cust-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}

.cust-table th, .cust-table td {
  padding: 0.75rem 0.85rem;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

.cust-table th {
  background: rgba(15, 23, 42, 0.8);
  color: var(--text-secondary);
  font-weight: 700;
}

.cust-name {
  color: #ffffff;
}

.points-badge {
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
  border: 1px solid rgba(245, 158, 11, 0.3);
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 800;
}

.addr-col {
  color: var(--text-muted);
  font-size: 0.8rem;
}

.action-btns {
  display: flex;
  gap: 0.35rem;
}

.btn-edit, .btn-delete {
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 0.25rem 0.4rem;
  cursor: pointer;
  font-size: 0.8rem;
}

.btn-delete:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.4);
}

.loading-box, .empty-box {
  padding: 3rem;
  text-align: center;
  color: var(--text-muted);
}
</style>
