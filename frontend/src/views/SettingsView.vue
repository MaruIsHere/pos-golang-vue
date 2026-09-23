<template>
  <div class="settings-page">
    <div class="page-header glass-panel">
      <div class="header-title">
        <h2>⚙️ Pengaturan Toko, QRIS & Voucher</h2>
        <p>Konfigurasi profil toko, gambar QRIS pembayaran, kode voucher diskon, dan database engine</p>
      </div>
    </div>

    <div class="settings-grid">
      <!-- 1. Store Profile & QRIS Settings -->
      <div class="settings-card glass-panel">
        <div class="card-header">
          <h3>🏪 Profil Toko & Foto QRIS Pembayaran</h3>
        </div>

        <form @submit.prevent="saveStoreSettings" class="card-body">
          <div class="form-group">
            <label class="form-label">Nama Toko / Usaha</label>
            <input type="text" class="form-control" v-model="storeForm.store_name" required />
          </div>

          <div class="form-group">
            <label class="form-label">Alamat Lengkap</label>
            <input type="text" class="form-control" v-model="storeForm.address" required />
          </div>

          <div class="form-group">
            <label class="form-label">No. Telepon / WhatsApp</label>
            <input type="text" class="form-control" v-model="storeForm.phone" required />
          </div>

          <div class="form-group">
            <label class="form-label">Pajak % (PPN / Service Charge)</label>
            <input type="number" class="form-control" v-model.number="storeForm.tax_percentage" min="0" max="100" />
          </div>

          <div class="form-group">
            <label class="form-label">Pesan Footer Struk</label>
            <textarea class="form-control" rows="2" v-model="storeForm.receipt_footer"></textarea>
          </div>

          <!-- QRIS Image Section -->
          <div class="form-group qris-upload-section">
            <label class="form-label">📷 Foto QRIS Pembayaran Toko</label>
            <p class="input-hint">Upload foto/gambar QRIS resmi toko Anda agar muncul di layar Kasir saat pelanggan memilih metode QRIS.</p>
            
            <div class="qris-preview-box">
              <div v-if="storeForm.qris_image_url" class="qris-img-container">
                <img :src="storeForm.qris_image_url" alt="QRIS Toko" class="qris-preview-img" />
                <button type="button" class="btn-remove-qris" @click="storeForm.qris_image_url = ''">
                  ✕ Hapus Foto QRIS
                </button>
              </div>

              <div v-else class="qris-empty-placeholder">
                <span class="placeholder-icon">📲</span>
                <p>Belum ada foto QRIS yang di-upload</p>
              </div>
            </div>

            <div class="qris-upload-actions">
              <label class="btn btn-secondary btn-upload">
                <span>📁 Upload Gambar QRIS</span>
                <input type="file" accept="image/*" @change="onQrisFileSelected" style="display: none;" />
              </label>
              
              <input 
                type="text" 
                class="form-control code-font" 
                v-model="storeForm.qris_image_url" 
                placeholder="Atau paste URL Gambar QRIS..." 
              />
            </div>
          </div>

          <button type="submit" class="btn btn-success btn-block" :disabled="isSavingStore">
            {{ isSavingStore ? 'Memproses...' : '💾 Simpan Pengaturan Toko & QRIS' }}
          </button>
        </form>
      </div>

      <!-- 2. Manage Voucher Codes -->
      <div class="settings-card glass-panel">
        <div class="card-header">
          <h3>🎟️ Manajemen Kode Voucher Diskon</h3>
        </div>

        <div class="card-body">
          <p class="section-desc">
            Tambah dan kelola kode voucher diskon yang bisa digunakan oleh kasir pada halaman transaksi.
          </p>

          <!-- Add Voucher Form -->
          <form @submit.prevent="createVoucher" class="voucher-form-box">
            <h4>+ Tambah Kode Voucher Baru</h4>

            <div class="voucher-form-grid">
              <div class="form-group">
                <label class="form-label">Kode Voucher</label>
                <input 
                  type="text" 
                  class="form-control code-font" 
                  v-model="newVoucher.code" 
                  placeholder="cth: DISKON30" 
                  required 
                />
              </div>

              <div class="form-group">
                <label class="form-label">Tipe Diskon</label>
                <select class="form-control" v-model="newVoucher.type">
                  <option value="percent">Persen (%)</option>
                  <option value="flat">Potongan Nominal (Rp)</option>
                </select>
              </div>

              <div class="form-group">
                <label class="form-label">Nilai Potongan</label>
                <input 
                  type="number" 
                  class="form-control" 
                  v-model.number="newVoucher.value" 
                  :placeholder="newVoucher.type === 'percent' ? 'cth: 15 (artinya 15%)' : 'cth: 20000 (artinya Rp 20.000)'" 
                  min="1" 
                  required 
                />
              </div>

              <div class="form-group">
                <label class="form-label">Deskripsi / Keterangan</label>
                <input 
                  type="text" 
                  class="form-control" 
                  v-model="newVoucher.description" 
                  placeholder="cth: Promo Tanggal Kembar" 
                />
              </div>
            </div>

            <button type="submit" class="btn btn-primary" :disabled="isCreatingVoucher">
              {{ isCreatingVoucher ? 'Menambahkan...' : '+ Simpan Voucher Baru' }}
            </button>
          </form>

          <!-- Vouchers Table List -->
          <div class="vouchers-list-container">
            <h4>Daftar Voucher Aktif ({{ vouchers.length }})</h4>

            <div v-if="isLoadingVouchers" class="vouchers-loading">
              <span>Memuat data voucher...</span>
            </div>

            <div v-else-if="vouchers.length === 0" class="empty-vouchers">
              <p>Belum ada kode voucher yang dibuat.</p>
            </div>

            <div v-else class="vouchers-table-wrapper">
              <table class="vouchers-table">
                <thead>
                  <tr>
                    <th>Kode</th>
                    <th>Tipe / Nilai</th>
                    <th>Deskripsi</th>
                    <th>Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="v in vouchers" :key="v.id">
                    <td>
                      <span class="voucher-code-badge">🎟️ {{ v.code }}</span>
                    </td>
                    <td>
                      <strong class="voucher-value">
                        {{ v.type === 'percent' ? v.value + '%' : 'Rp ' + formatPrice(v.value) }}
                      </strong>
                    </td>
                    <td class="voucher-desc-col">{{ v.description || '-' }}</td>
                    <td>
                      <button class="btn-delete-voucher" title="Hapus Voucher" @click="deleteVoucher(v.id ?? 0, v.code)">
                        🗑️ Hapus
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- 3. Database Engine Switcher Panel -->
      <div class="settings-card glass-panel highlight-card full-width-card">
        <div class="card-header">
          <h3>🗄️ Engine Basis Data (Database Switcher)</h3>
          <span class="active-db-badge" :class="dbEngine === 'mysql' ? 'mysql' : 'sqlite'">
            Aktif: {{ dbEngine.toUpperCase() }}
          </span>
        </div>

        <div class="card-body">
          <p class="section-desc">
            Pilih engine basis data yang ingin digunakan oleh backend Golang.
          </p>

          <div class="db-options-grid">
            <div 
              class="db-option-card" 
              :class="{ selected: selectedEngine === 'sqlite' }"
              @click="selectedEngine = 'sqlite'"
            >
              <div class="db-icon">📂</div>
              <div class="db-info">
                <h4>SQLite (Embedded File)</h4>
                <p>Offline-first, tanpa butuh server MySQL. Data disimpan di file <code>pos.db</code>.</p>
              </div>
            </div>

            <div 
              class="db-option-card" 
              :class="{ selected: selectedEngine === 'mysql' }"
              @click="selectedEngine = 'mysql'"
            >
              <div class="db-icon">🐬</div>
              <div class="db-info">
                <h4>MySQL Server</h4>
                <p>Terpusat, cocok untuk multi-kasir di jaringan lokal/server cloud.</p>
              </div>
            </div>
          </div>

          <!-- MySQL DSN Config Form -->
          <div v-if="selectedEngine === 'mysql'" class="mysql-config-box">
            <div class="form-group">
              <label class="form-label">MySQL Connection String (DSN)</label>
              <input 
                type="text" 
                class="form-control code-font" 
                v-model="mysqlDsn" 
                placeholder="root:password@tcp(127.0.0.1:3306)/pos_db?charset=utf8mb4&parseTime=True&loc=Local" 
              />
              <span class="input-hint">Format GORM MySQL: user:pass@tcp(host:port)/dbname?params</span>
            </div>
          </div>

          <div v-if="dbMessage" class="db-feedback-msg" :class="dbMessageType">
            {{ dbMessage }}
          </div>

          <button 
            class="btn btn-primary btn-switch-db" 
            :disabled="isSwitchingDb"
            @click="switchDatabaseEngine"
          >
            <span v-if="isSwitchingDb">Mengubah Engine & Migrasi Data...</span>
            <span v-else>⚡ Terapkan & Switch Database</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { Voucher } from '../types';

const emit = defineEmits(['refresh-settings']);

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

const dbEngine = ref('sqlite');
const selectedEngine = ref('sqlite');
const mysqlDsn = ref('root:@tcp(127.0.0.1:3306)/pos_db?charset=utf8mb4&parseTime=True&loc=Local');

const isSwitchingDb = ref(false);
const dbMessage = ref('');
const dbMessageType = ref('success');

const isSavingStore = ref(false);

const storeForm = ref({
  store_name: '',
  address: '',
  phone: '',
  tax_percentage: 10,
  receipt_footer: '',
  qris_image_url: ''
});

// Vouchers State
const vouchers = ref<Voucher[]>([]);
const isLoadingVouchers = ref(false);
const isCreatingVoucher = ref(false);
const newVoucher = ref({
  code: '',
  type: 'percent',
  value: 10,
  description: ''
});

const loadSettings = async () => {
  try {
    const res = await fetch('/api/settings');
    if (res.ok) {
      const data = await res.json();
      dbEngine.value = data.db_engine || 'sqlite';
      selectedEngine.value = dbEngine.value;
      if (data.mysql_dsn) mysqlDsn.value = data.mysql_dsn;

      storeForm.value = {
        store_name: data.store_name,
        address: data.address,
        phone: data.phone,
        tax_percentage: data.tax_percentage,
        receipt_footer: data.receipt_footer,
        qris_image_url: data.qris_image_url || ''
      };
    }
  } catch (err) {
    console.error(err);
  }
};

const loadVouchers = async () => {
  isLoadingVouchers.value = true;
  try {
    const res = await fetch('/api/vouchers');
    if (res.ok) {
      vouchers.value = await res.json();
    }
  } catch (err) {
    console.error('Fetch vouchers error:', err);
  } finally {
    isLoadingVouchers.value = false;
  }
};

onMounted(() => {
  loadSettings();
  loadVouchers();
});

const onQrisFileSelected = (event: Event): void => {
  const input = event.target as HTMLInputElement | null;
  const file = input?.files?.[0];
  if (!file) return;

  if (file.size > 5 * 1024 * 1024) {
    alert('Ukuran file foto terlalu besar (maksimal 5MB)');
    return;
  }

  const reader = new FileReader();
  reader.onload = (e: ProgressEvent<FileReader>) => {
    storeForm.value.qris_image_url = (e.target?.result as string) ?? '';
  };
  reader.readAsDataURL(file);
};

const saveStoreSettings = async () => {
  isSavingStore.value = true;
  try {
    const res = await fetch('/api/settings', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(storeForm.value)
    });

    if (res.ok) {
      alert('Pengaturan profil toko & foto QRIS berhasil disimpan!');
      emit('refresh-settings');
    } else {
      alert('Gagal menyimpan pengaturan toko');
    }
  } catch (err) {
    alert('Koneksi error: ' + (err as Error).message);
  } finally {
    isSavingStore.value = false;
  }
};

const createVoucher = async () => {
  if (!newVoucher.value.code) return;

  isCreatingVoucher.value = true;
  try {
    const payload = {
      code: newVoucher.value.code.trim().toUpperCase(),
      type: newVoucher.value.type,
      value: newVoucher.value.value,
      description: newVoucher.value.description
    };

    const res = await fetch('/api/vouchers', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });

    if (res.ok) {
      alert(`Kode voucher '${payload.code}' berhasil ditambahkan!`);
      newVoucher.value = { code: '', type: 'percent', value: 10, description: '' };
      loadVouchers();
    } else {
      const errData = await res.json();
      alert('Gagal menambahkan voucher: ' + (errData.error || 'Terjadi kesalahan'));
    }
  } catch (err) {
    alert('Koneksi error: ' + (err as Error).message);
  } finally {
    isCreatingVoucher.value = false;
  }
};

const deleteVoucher = async (id: number, code: string): Promise<void> => {
  if (!confirm(`Apakah Anda yakin ingin menghapus voucher '${code}'?`)) return;

  try {
    const res = await fetch(`/api/vouchers/${id}`, {
      method: 'DELETE'
    });

    if (res.ok) {
      loadVouchers();
    } else {
      alert('Gagal menghapus voucher');
    }
  } catch (err) {
    alert('Koneksi error: ' + (err as Error).message);
  }
};

const switchDatabaseEngine = async () => {
  isSwitchingDb.value = true;
  dbMessage.value = '';
  try {
    const res = await fetch('/api/settings/switch-db', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        engine: selectedEngine.value,
        mysql_dsn: mysqlDsn.value
      })
    });

    const data = await res.json();
    if (res.ok) {
      dbEngine.value = data.db_engine;
      dbMessage.value = data.message || 'Engine database berhasil diperbarui!';
      dbMessageType.value = 'success';
      emit('refresh-settings');
    } else {
      dbMessage.value = data.error || 'Gagal mengubah database';
      dbMessageType.value = 'error';
    }
  } catch (err) {
    dbMessage.value = 'Terjadi kesalahan koneksi server: ' + (err as Error).message;
    dbMessageType.value = 'error';
  } finally {
    isSwitchingDb.value = false;
  }
};
</script>

<style scoped>
.settings-page {
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

.settings-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.25rem;
}

@media (min-width: 1024px) {
  .settings-grid {
    grid-template-columns: 1fr 1fr;
  }
  .full-width-card {
    grid-column: span 2;
  }
}

.settings-card {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.highlight-card {
  border-color: rgba(99, 102, 241, 0.4);
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

/* QRIS Section */
.qris-upload-section {
  border-top: 1px dashed var(--border-color);
  padding-top: 1rem;
  margin-top: 0.5rem;
}

.qris-preview-box {
  margin: 0.75rem 0;
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(15, 23, 42, 0.8);
  border: 1px dashed var(--border-color);
  border-radius: var(--radius-md);
  padding: 1rem;
  min-height: 160px;
}

.qris-img-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.qris-preview-img {
  max-width: 200px;
  max-height: 200px;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
}

.btn-remove-qris {
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid rgba(239, 68, 68, 0.4);
  color: #f87171;
  padding: 0.25rem 0.65rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 700;
  cursor: pointer;
}

.qris-empty-placeholder {
  text-align: center;
  color: var(--text-muted);
  font-size: 0.82rem;
}

.placeholder-icon {
  font-size: 2.2rem;
  display: block;
  margin-bottom: 0.35rem;
}

.qris-upload-actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.btn-upload {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  text-align: center;
  font-size: 0.85rem;
}

.btn-block {
  width: 100%;
  padding: 0.85rem;
  font-size: 0.95rem;
  margin-top: 0.5rem;
}

/* Voucher Box */
.voucher-form-box {
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid var(--border-color);
  padding: 1rem;
  border-radius: var(--radius-md);
  margin-bottom: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
}

.voucher-form-box h4 {
  font-size: 0.95rem;
  font-weight: 800;
  color: var(--accent-secondary);
}

.voucher-form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}

.vouchers-list-container h4 {
  font-size: 0.95rem;
  font-weight: 800;
  margin-bottom: 0.75rem;
}

.vouchers-table-wrapper {
  overflow-x: auto;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
}

.vouchers-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.82rem;
}

.vouchers-table th, .vouchers-table td {
  padding: 0.65rem 0.85rem;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

.vouchers-table th {
  background: rgba(15, 23, 42, 0.8);
  color: var(--text-secondary);
  font-weight: 700;
}

.voucher-code-badge {
  font-weight: 800;
  color: #a5b4fc;
}

.voucher-value {
  color: #34d399;
}

.voucher-desc-col {
  color: var(--text-muted);
}

.btn-delete-voucher {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 700;
  cursor: pointer;
}

.empty-vouchers, .vouchers-loading {
  padding: 1.5rem;
  text-align: center;
  color: var(--text-muted);
  font-size: 0.85rem;
  background: rgba(15, 23, 42, 0.4);
  border-radius: var(--radius-md);
}

.active-db-badge {
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 800;
}

.active-db-badge.sqlite {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.4);
}

.active-db-badge.mysql {
  background: rgba(245, 158, 11, 0.2);
  color: #fbbf24;
  border: 1px solid rgba(245, 158, 11, 0.4);
}

.section-desc {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.db-options-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.85rem;
  margin-top: 0.75rem;
}

@media (min-width: 640px) {
  .db-options-grid {
    grid-template-columns: 1fr 1fr;
  }
}

.db-option-card {
  display: flex;
  align-items: flex-start;
  gap: 0.85rem;
  padding: 1rem;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s;
}

.db-option-card:hover {
  background: var(--bg-card-hover);
}

.db-option-card.selected {
  background: rgba(99, 102, 241, 0.15);
  border-color: var(--accent-primary);
  box-shadow: 0 0 15px rgba(99, 102, 241, 0.25);
}

.db-icon {
  font-size: 1.6rem;
}

.db-info h4 {
  font-size: 0.95rem;
  font-weight: 800;
  color: var(--text-primary);
}

.db-info p {
  font-size: 0.78rem;
  color: var(--text-muted);
  margin-top: 0.2rem;
}

.mysql-config-box {
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid var(--border-color);
  padding: 1rem;
  border-radius: var(--radius-md);
  margin-top: 0.75rem;
}

.code-font {
  font-family: monospace;
  font-size: 0.85rem;
}

.input-hint {
  font-size: 0.7rem;
  color: var(--text-muted);
  margin-top: 0.2rem;
  display: block;
}

.db-feedback-msg {
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  font-size: 0.85rem;
  font-weight: 700;
}

.db-feedback-msg.success {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.4);
}

.db-feedback-msg.error {
  background: rgba(239, 68, 68, 0.2);
  color: #f87171;
  border: 1px solid rgba(239, 68, 68, 0.4);
}

.btn-switch-db {
  width: 100%;
  padding: 0.85rem;
  margin-top: 0.5rem;
}
</style>
