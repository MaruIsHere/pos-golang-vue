<template>
  <div class="settings-page">
    <div class="page-header glass-panel">
      <div class="header-title">
        <h2>⚙️ Pengaturan Toko & Basis Data</h2>
        <p>Konfigurasi profil usaha dan switch engine database (SQLite / MySQL)</p>
      </div>
    </div>

    <div class="settings-grid">
      <!-- Database Engine Switcher Panel -->
      <div class="settings-card glass-panel highlight-card">
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

      <!-- Store Profile Settings -->
      <div class="settings-card glass-panel">
        <div class="card-header">
          <h3>🏪 Profil Toko & Struk</h3>
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
            <textarea class="form-control" rows="3" v-model="storeForm.receipt_footer"></textarea>
          </div>

          <button type="submit" class="btn btn-success" :disabled="isSavingStore">
            {{ isSavingStore ? 'Memproses...' : 'Simpan Pengaturan Toko' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const props = defineProps({
  storeSetting: { type: Object, default: () => ({}) }
});

const emit = defineEmits(['refresh-settings']);

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
  receipt_footer: ''
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
        receipt_footer: data.receipt_footer
      };
    }
  } catch (err) {
    console.error(err);
  }
};

onMounted(loadSettings);

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
    dbMessage.value = 'Terjadi kesalahan koneksi server: ' + err.message;
    dbMessageType.value = 'error';
  } finally {
    isSwitchingDb.value = false;
  }
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
      alert('Pengaturan profil toko berhasil disimpan!');
      emit('refresh-settings');
    } else {
      alert('Gagal menyimpan pengaturan toko');
    }
  } catch (err) {
    alert('Koneksi error: ' + err.message);
  } finally {
    isSavingStore.value = false;
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
