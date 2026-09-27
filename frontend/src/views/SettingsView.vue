<template>
  <div class="settings-page">
    <div class="page-header glass-panel">
      <div class="header-title">
        <h2>Pengaturan Toko, QRIS, Tema & Voucher</h2>
        <p>Konfigurasi profil toko, gambar QRIS pembayaran, tema mode terang/gelap, voucher diskon, dan database engine</p>
      </div>
    </div>

    <div class="settings-grid">
      <!-- 0. Theme Selection Card -->
      <div class="settings-card glass-panel full-width-card">
        <div class="card-header">
          <h3 class="flex items-center gap-1.5">
            <SunIcon class="w-5 h-5 text-amber-500" />
            <span>Tema Tampilan (Mode Terang / Mode Gelap)</span>
          </h3>
          <span class="active-db-badge" :class="isDarkMode ? 'mysql' : 'sqlite'">
            {{ isDarkMode ? 'Mode Gelap (Dark)' : 'Mode Terang (Light)' }}
          </span>
        </div>

        <div class="card-body">
          <p class="section-desc">Pilih tema warna tampilan antarmuka aplikasi kasir.</p>

          <div class="db-options-grid">
            <div 
              class="db-option-card" 
              :class="{ selected: !isDarkMode }"
              @click="setDark(false)"
            >
              <div class="db-icon">
                <SunIcon class="w-7 h-7 text-amber-500" />
              </div>
              <div class="db-info">
                <h4>Mode Terang (Light Mode)</h4>
                <p>Tampilan serba putih yang bersih, minimalis, dan terang.</p>
              </div>
            </div>

            <div 
              class="db-option-card" 
              :class="{ selected: isDarkMode }"
              @click="setDark(true)"
            >
              <div class="db-icon">
                <MoonIcon class="w-7 h-7 text-indigo-400" />
              </div>
              <div class="db-info">
                <h4>Mode Gelap (Dark Mode)</h4>
                <p>Tampilan gelap modern yang nyaman di mata untuk kondisi pencahayaan redup.</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 1. Store Profile & QRIS Settings -->
      <div class="settings-card glass-panel">
        <div class="card-header">
          <h3 class="flex items-center gap-1.5">
            <BuildingStorefrontIcon class="w-5 h-5 text-indigo-600" />
            <span>Profil Toko & Foto QRIS Pembayaran</span>
          </h3>
        </div>

        <form @submit.prevent="saveStoreSettings" class="card-body">
          <div class="form-group">
            <label class="form-label">Nama Toko / Usaha *</label>
            <input type="text" class="form-control" v-model="storeForm.store_name" required />
          </div>

          <div class="form-group">
            <label class="form-label">Alamat Lengkap *</label>
            <input type="text" class="form-control" v-model="storeForm.address" required />
          </div>

          <div class="form-group">
            <label class="form-label">No. Telepon / WhatsApp *</label>
            <input type="text" class="form-control" v-model="storeForm.phone" required />
          </div>

          <div class="form-group">
            <label class="form-label">Pajak % (PPN / Service Charge)</label>
            <input type="number" class="form-control" v-model.number="storeForm.tax_percentage" min="0" max="100" />
          </div>

          <div class="form-group">
            <label class="form-label">Potongan Diskon Pelanggan Terdaftar / Member (%)</label>
            <input type="number" class="form-control" v-model.number="storeForm.member_discount_percentage" min="0" max="100" step="0.5" />
            <span class="input-hint">Potongan ini otomatis dihitung saat kasir memilih/menginput nama pelanggan terdaftar saat checkout.</span>
          </div>

          <div class="form-group">
            <label class="form-label">Pesan Footer Struk</label>
            <textarea class="form-control" rows="2" v-model="storeForm.receipt_footer"></textarea>
          </div>

          <!-- QRIS Image Section -->
          <div class="form-group qris-upload-section">
            <label class="form-label">Foto QRIS Pembayaran Toko</label>
            <p class="input-hint">Upload foto/gambar QRIS resmi toko Anda agar muncul di layar Kasir saat pelanggan memilih metode QRIS.</p>
            
            <div class="qris-preview-box">
              <div v-if="storeForm.qris_image_url" class="qris-img-container">
                <img :src="storeForm.qris_image_url" alt="QRIS Toko" class="qris-preview-img" />
                <button type="button" class="btn-remove-qris" @click="storeForm.qris_image_url = ''">
                  Hapus Foto QRIS
                </button>
              </div>

              <div v-else class="qris-empty-placeholder">
                <p>Belum ada foto QRIS yang di-upload</p>
              </div>
            </div>

            <div class="qris-upload-actions">
              <label class="btn btn-secondary btn-upload flex items-center justify-center gap-1">
                <ArrowUpTrayIcon class="w-4 h-4" />
                <span>Upload Gambar QRIS</span>
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
            {{ isSavingStore ? 'Memproses...' : 'Simpan Pengaturan Toko & QRIS' }}
          </button>
        </form>
      </div>

      <!-- 2. Manage Voucher Codes -->
      <div class="settings-card glass-panel">
        <div class="card-header">
          <h3 class="flex items-center gap-1.5">
            <TicketIcon class="w-5 h-5 text-indigo-600" />
            <span>Manajemen Kode Voucher Diskon</span>
          </h3>
        </div>

        <div class="card-body">
          <p class="section-desc">
            Tambah dan kelola kode voucher diskon yang bisa digunakan oleh kasir pada halaman transaksi.
          </p>

          <!-- Add Voucher Form -->
          <form @submit.prevent="createVoucher" class="voucher-form-box">
            <h4>Tambah Kode Voucher Baru</h4>

            <div class="voucher-form-grid">
              <div class="form-group">
                <label class="form-label">Kode Voucher *</label>
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
                <label class="form-label">Nilai Potongan *</label>
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
              {{ isCreatingVoucher ? 'Menambahkan...' : 'Simpan Voucher Baru' }}
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
                    <th class="text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="v in vouchers" :key="v.id">
                    <td>
                      <span class="voucher-code-badge">{{ v.code }}</span>
                    </td>
                    <td>
                      <strong class="voucher-value">
                        {{ v.type === 'percent' ? v.value + '%' : 'Rp ' + formatPrice(v.value) }}
                      </strong>
                    </td>
                    <td class="voucher-desc-col">{{ v.description || '-' }}</td>
                    <td class="text-right">
                      <button class="btn-icon btn-delete-voucher" title="Hapus Voucher" @click="deleteVoucher(v.id ?? 0, v.code)">
                        <TrashIcon class="w-3.5 h-3.5 text-red-600 inline-block mr-1" /> Hapus
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- 3. User Management & Registration (Owner & Admin Only) -->
      <div class="settings-card glass-panel full-width-card">
        <div class="card-header">
          <h3 class="flex items-center gap-1.5">
            <UserGroupIcon class="w-5 h-5 text-indigo-600" />
            <span>Manajemen Pengguna & Registrasi Akun Staf</span>
          </h3>
          <span class="active-db-badge sqlite">
            Khusus Owner & Admin
          </span>
        </div>

        <div class="card-body">
          <p class="section-desc">
            Tambah akun pengguna baru (Kasir, Kepala Kasir, Admin, atau Owner) dan kelola daftar staf kasir yang memiliki akses ke sistem POS.
          </p>

          <!-- Register User Form -->
          <form @submit.prevent="createUser" class="voucher-form-box mb-6">
            <h4>Tambah Pengguna / Registrasi Staf Baru</h4>

            <div class="voucher-form-grid">
              <div class="form-group">
                <label class="form-label">Username *</label>
                <input 
                  type="text" 
                  class="form-control" 
                  v-model="newUser.username" 
                  placeholder="Masukkan username" 
                  required 
                />
              </div>

              <div class="form-group">
                <label class="form-label">Password *</label>
                <input 
                  type="password" 
                  class="form-control" 
                  v-model="newUser.password" 
                  placeholder="Minimal 6 karakter" 
                  required 
                />
              </div>

              <div class="form-group">
                <label class="form-label">Role / Peran Akses *</label>
                <select class="form-control" v-model="newUser.role">
                  <option value="kasir">Kasir (Transaksi & Struk)</option>
                  <option value="kepala_kasir">Kepala Kasir (Kelola Produk, Stok, Pelanggan, Diskon & Laporan)</option>
                  <option value="admin">Admin (Full Akses Sistem)</option>
                  <option value="owner">Owner (Full Akses & Pengaturan Toko)</option>
                </select>
              </div>
            </div>

            <div v-if="userSuccessMsg" class="text-xs text-emerald-600 font-bold mb-3">
              {{ userSuccessMsg }}
            </div>
            <div v-if="userErrorMsg" class="text-xs text-rose-600 font-bold mb-3">
              {{ userErrorMsg }}
            </div>

            <button type="submit" class="btn btn-primary" :disabled="isCreatingUser">
              {{ isCreatingUser ? 'Menambahkan...' : 'Tambah Pengguna Baru' }}
            </button>
          </form>

          <!-- Users List Table -->
          <div class="vouchers-list-container">
            <h4>
              Daftar Pengguna Terdaftar
              <span v-if="!isLoadingUsers && !usersLoadError">({{ users.length }})</span>
            </h4>

            <div v-if="isLoadingUsers" class="vouchers-loading">
              <span>Memuat data pengguna...</span>
            </div>

            <div v-else-if="usersLoadError" class="empty-vouchers">
              <p>{{ usersLoadError }}</p>
              <button type="button" class="btn btn-secondary" @click="loadUsers">Coba Lagi</button>
            </div>

            <div v-else-if="users.length === 0" class="empty-vouchers">
              <p>Belum ada pengguna terdaftar.</p>
            </div>

            <div v-else class="vouchers-table-wrapper">
              <table class="vouchers-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Username</th>
                    <th>Role / Akses</th>
                    <th>Tanggal Dibuat</th>
                    <th class="text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  <template v-for="u in users" :key="u.id">
                    <tr>
                    <td>#{{ u.id }}</td>
                    <td>
                      <strong class="text-slate-800 dark:text-slate-100 font-bold">{{ u.username }}</strong>
                    </td>
                    <td>
                      <span class="badge" :class="getUserRoleBadgeClass(u.role)">
                        {{ formatUserRoleLabel(u.role) }}
                      </span>
                    </td>
                    <td class="text-xs text-slate-500">{{ formatDate(u.created_at) }}</td>
                    <td class="text-right">
                      <div class="user-actions">
                        <button
                          type="button"
                          class="btn-icon btn-change-password"
                          :disabled="isChangingStaffPassword"
                          @click="openStaffPasswordForm(u)"
                        >
                          Ganti Sandi
                        </button>
                        <button
                          type="button"
                          class="btn-icon btn-delete-voucher"
                          title="Hapus Pengguna"
                          @click="deleteUserAccount(u.id, u.username)"
                        >
                          <TrashIcon class="w-3.5 h-3.5 text-red-600 inline-block mr-1" /> Hapus
                        </button>
                      </div>
                    </td>
                  </tr>
                    <tr v-if="staffPasswordTarget?.id === u.id">
                      <td colspan="5">
                        <form class="staff-password-form" @submit.prevent="changeStaffPassword">
                          <div class="staff-password-heading">
                            <strong>Ganti sandi untuk {{ staffPasswordTarget.username }}</strong>
                            <span class="text-xs text-slate-500">Minimal 6 karakter</span>
                          </div>
                          <div class="staff-password-fields">
                            <div class="form-group">
                              <label class="form-label" :for="`new-password-${u.id}`">Sandi baru</label>
                              <input
                                :id="`new-password-${u.id}`"
                                v-model="newStaffPassword"
                                class="form-control"
                                type="password"
                                autocomplete="new-password"
                                minlength="6"
                                required
                              />
                            </div>
                            <div class="form-group">
                              <label class="form-label" :for="`confirm-password-${u.id}`">Konfirmasi sandi baru</label>
                              <input
                                :id="`confirm-password-${u.id}`"
                                v-model="confirmStaffPassword"
                                class="form-control"
                                type="password"
                                autocomplete="new-password"
                                minlength="6"
                                required
                              />
                            </div>
                          </div>
                          <p v-if="staffPasswordSuccessMsg" class="staff-password-success">{{ staffPasswordSuccessMsg }}</p>
                          <p v-if="staffPasswordErrorMsg" class="staff-password-error">{{ staffPasswordErrorMsg }}</p>
                          <div class="user-actions">
                            <button type="button" class="btn btn-secondary" :disabled="isChangingStaffPassword" @click="cancelStaffPasswordChange">
                              Batal
                            </button>
                            <button type="submit" class="btn btn-primary" :disabled="isChangingStaffPassword">
                              {{ isChangingStaffPassword ? 'Menyimpan...' : 'Simpan Sandi Baru' }}
                            </button>
                          </div>
                        </form>
                      </td>
                    </tr>
                  </template>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- 4. Database Engine Switcher Panel -->
      <div class="settings-card glass-panel highlight-card full-width-card">
        <div class="card-header">
          <h3 class="flex items-center gap-1.5">
            <CircleStackIcon class="w-5 h-5 text-indigo-600" />
            <span>Engine Basis Data (Database Switcher)</span>
          </h3>
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
              <div class="db-icon">
                <FolderIcon class="w-7 h-7 text-indigo-600" />
              </div>
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
              <div class="db-icon">
                <ServerIcon class="w-7 h-7 text-indigo-600" />
              </div>
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
            <span v-else>Terapkan & Switch Database</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '@/utils/api';
import type { Voucher } from '../types';
import { useTheme } from '../composables/useTheme';
import { 
  BuildingStorefrontIcon, 
  TicketIcon, 
  CircleStackIcon, 
  FolderIcon, 
  ServerIcon, 
  TrashIcon, 
  ArrowUpTrayIcon,
  SunIcon,
  MoonIcon,
  UserGroupIcon
} from '@heroicons/vue/24/outline';

const emit = defineEmits(['refresh-settings']);

const { isDarkMode, setDark } = useTheme();

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
  member_discount_percentage: 5,
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

// Users State
const users = ref<any[]>([]);
const isLoadingUsers = ref(false);
const usersLoadError = ref('');
const isCreatingUser = ref(false);
const userSuccessMsg = ref('');
const userErrorMsg = ref('');
const newUser = ref({
  username: '',
  password: '',
  role: 'kasir'
});
const staffPasswordTarget = ref<{ id: number; username: string } | null>(null);
const newStaffPassword = ref('');
const confirmStaffPassword = ref('');
const isChangingStaffPassword = ref(false);
const staffPasswordSuccessMsg = ref('');
const staffPasswordErrorMsg = ref('');

const loadSettings = async () => {
  try {
    const res = await api.get('/settings');
    const data = res.data;
    dbEngine.value = data.db_engine || 'sqlite';
    selectedEngine.value = dbEngine.value;
    if (data.mysql_dsn) mysqlDsn.value = data.mysql_dsn;

    storeForm.value = {
      store_name: data.store_name,
      address: data.address,
      phone: data.phone,
      tax_percentage: data.tax_percentage,
      member_discount_percentage: data.member_discount_percentage ?? 5,
      receipt_footer: data.receipt_footer,
      qris_image_url: data.qris_image_url || ''
    };
  } catch (err: any) {
    console.error(err.response?.data?.error || err.message || 'Error occurred');
  }
};

const loadVouchers = async () => {
  isLoadingVouchers.value = true;
  try {
    const res = await api.get('/vouchers');
    vouchers.value = res.data;
  } catch (err: any) {
    console.error('Fetch vouchers error:', err.response?.data?.error || err.message || 'Error occurred');
  } finally {
    isLoadingVouchers.value = false;
  }
};

const loadUsers = async () => {
  isLoadingUsers.value = true;
  usersLoadError.value = '';
  const controller = new AbortController();
  const timeoutId = window.setTimeout(() => controller.abort(), 10000);
  try {
    const res = await api.get('/users', { signal: controller.signal });
    users.value = res.data || [];
  } catch (err: any) {
    usersLoadError.value = err.name === 'AbortError'
      ? 'Permintaan daftar pengguna terlalu lama. Periksa koneksi backend, lalu coba lagi.'
      : err.response?.data?.error || err.message || 'Gagal memuat daftar pengguna.';
  } finally {
    window.clearTimeout(timeoutId);
    isLoadingUsers.value = false;
  }
};

const createUser = async () => {
  userSuccessMsg.value = '';
  userErrorMsg.value = '';
  isCreatingUser.value = true;
  try {
    const res = await api.post('/users', newUser.value);
    userSuccessMsg.value = res.data.message || 'Pengguna baru berhasil ditambahkan!';
    newUser.value = { username: '', password: '', role: 'kasir' };
    await loadUsers();
  } catch (err: any) {
    userErrorMsg.value = err.response?.data?.error || err.message || 'Gagal menambahkan pengguna';
  } finally {
    isCreatingUser.value = false;
  }
};

const openStaffPasswordForm = (user: { id: number; username: string }) => {
  staffPasswordTarget.value = user;
  newStaffPassword.value = '';
  confirmStaffPassword.value = '';
  staffPasswordSuccessMsg.value = '';
  staffPasswordErrorMsg.value = '';
};

const cancelStaffPasswordChange = () => {
  staffPasswordTarget.value = null;
  newStaffPassword.value = '';
  confirmStaffPassword.value = '';
  staffPasswordSuccessMsg.value = '';
  staffPasswordErrorMsg.value = '';
};

const changeStaffPassword = async () => {
  const target = staffPasswordTarget.value;
  if (!target) return;

  staffPasswordSuccessMsg.value = '';
  staffPasswordErrorMsg.value = '';
  if (newStaffPassword.value.length < 6) {
    staffPasswordErrorMsg.value = 'Sandi minimal 6 karakter.';
    return;
  }
  if (newStaffPassword.value !== confirmStaffPassword.value) {
    staffPasswordErrorMsg.value = 'Konfirmasi sandi tidak cocok.';
    return;
  }

  isChangingStaffPassword.value = true;
  try {
    const res = await api.put(`/users/${target.id}/password`, { password: newStaffPassword.value });
    staffPasswordSuccessMsg.value = res.data.message || 'Sandi staff berhasil diperbarui.';
    newStaffPassword.value = '';
    confirmStaffPassword.value = '';
  } catch (err: any) {
    staffPasswordErrorMsg.value = err.response?.data?.error || err.message || 'Gagal mengganti sandi staff.';
  } finally {
    isChangingStaffPassword.value = false;
  }
};

const deleteUserAccount = async (id: number, username: string) => {
  if (!confirm(`Apakah Anda yakin ingin menghapus akun pengguna "${username}"?`)) return;
  try {
    await api.delete(`/users/${id}`);
    alert('User berhasil dihapus!');
    await loadUsers();
  } catch (err: any) {
    alert('Gagal menghapus user: ' + (err.response?.data?.error || err.message));
  }
};

const formatUserRoleLabel = (role: string): string => {
  const r = (role || '').toLowerCase();
  if (r === 'owner') return 'Owner (Pemilik)';
  if (r === 'admin') return 'Admin (Full Akses)';
  if (r === 'kepala_kasir') return 'Kepala Kasir';
  return 'Kasir';
};

const getUserRoleBadgeClass = (role: string): string => {
  const r = (role || '').toLowerCase();
  if (r === 'owner' || r === 'admin') return 'badge-warning';
  if (r === 'kepala_kasir') return 'badge-success';
  return 'badge-info';
};

const formatDate = (dateStr: string): string => {
  if (!dateStr) return '-';
  try {
    return new Date(dateStr).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    });
  } catch {
    return dateStr;
  }
};

onMounted(() => {
  loadSettings();
  loadVouchers();
  loadUsers();
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
    await api.put('/settings', storeForm.value);
    alert('Pengaturan profil toko & foto QRIS berhasil disimpan!');
    emit('refresh-settings');
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    alert('Koneksi error: ' + errMsg);
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

    await api.post('/vouchers', payload);
    alert(`Kode voucher '${payload.code}' berhasil ditambahkan!`);
    newVoucher.value = { code: '', type: 'percent', value: 10, description: '' };
    loadVouchers();
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    alert('Koneksi error: ' + errMsg);
  } finally {
    isCreatingVoucher.value = false;
  }
};

const deleteVoucher = async (id: number, code: string): Promise<void> => {
  if (!confirm(`Apakah Anda yakin ingin menghapus voucher '${code}'?`)) return;

  try {
    await api.delete(`/vouchers/${id}`);
    loadVouchers();
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    alert('Koneksi error: ' + errMsg);
  }
};

const switchDatabaseEngine = async () => {
  isSwitchingDb.value = true;
  dbMessage.value = '';
  try {
    const res = await api.post('/settings/switch-db', {
      engine: selectedEngine.value,
      mysql_dsn: mysqlDsn.value
    });

    const data = res.data;
    dbEngine.value = data.db_engine;
    dbMessage.value = data.message || 'Engine database berhasil diperbarui!';
    dbMessageType.value = 'success';
    emit('refresh-settings');
  } catch (err: any) {
    dbMessage.value = 'Terjadi kesalahan koneksi server: ' + (err.response?.data?.error || err.message || 'Error occurred');
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
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.header-title h2 {
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--text-primary);
}

.header-title p {
  font-size: 0.8rem;
  color: var(--text-secondary);
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
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.highlight-card {
  border-color: var(--border-color);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.85rem;
}

.card-header h3 {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-primary);
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
  background: var(--bg-primary);
  border: 1px dashed var(--border-color);
  border-radius: 10px;
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
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.btn-remove-qris {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
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
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  padding: 1rem;
  border-radius: 10px;
  margin-bottom: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
}

.voucher-form-box h4 {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-primary);
}

.voucher-form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}

.vouchers-list-container h4 {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 0.75rem;
}

.vouchers-table-wrapper {
  overflow-x: auto;
  border: 1px solid var(--border-color);
  border-radius: 8px;
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
  background: var(--bg-primary);
  color: var(--text-secondary);
  font-weight: 600;
}

.voucher-code-badge {
  font-weight: 700;
  color: var(--accent-primary);
  background: rgba(99, 102, 241, 0.12);
  padding: 0.15rem 0.5rem;
  border-radius: 6px;
}

.voucher-value {
  color: var(--accent-secondary);
}

.voucher-desc-col {
  color: var(--text-secondary);
}

.btn-delete-voucher {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  color: #ef4444;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
}
.btn-delete-voucher:hover {
  background: rgba(239, 68, 68, 0.1);
}

.user-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 0.4rem;
}

.btn-change-password {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
}

.btn-change-password:hover {
  border-color: var(--accent-primary);
  color: var(--accent-primary);
}

.btn-change-password:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.staff-password-form {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  padding: 1rem;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
}

.staff-password-heading {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 0.5rem;
  color: var(--text-primary);
}

.staff-password-fields {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.75rem;
}

.staff-password-success,
.staff-password-error {
  margin: 0;
  font-size: 0.8rem;
  font-weight: 600;
}

.staff-password-success {
  color: #059669;
}

.staff-password-error {
  color: #e11d48;
}

@media (min-width: 640px) {
  .staff-password-fields {
    grid-template-columns: 1fr 1fr;
  }
}

.empty-vouchers, .vouchers-loading {
  padding: 1.5rem;
  text-align: center;
  color: var(--text-secondary);
  font-size: 0.85rem;
  background: var(--bg-primary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.active-db-badge {
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
}

.active-db-badge.sqlite {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}
html.dark .active-db-badge.sqlite {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
  border-color: rgba(59, 130, 246, 0.3);
}

.active-db-badge.mysql {
  background: #fffbeb;
  color: #b45309;
  border: 1px solid #fde68a;
}
html.dark .active-db-badge.mysql {
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
  border-color: rgba(245, 158, 11, 0.3);
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
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.db-option-card:hover {
  background: var(--bg-card-hover);
}

.db-option-card.selected {
  background: rgba(99, 102, 241, 0.1);
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2);
}

.db-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.db-info h4 {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-primary);
}

.db-info p {
  font-size: 0.78rem;
  color: var(--text-secondary);
  margin-top: 0.2rem;
}

.mysql-config-box {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  padding: 1rem;
  border-radius: 10px;
  margin-top: 0.75rem;
}

.code-font {
  font-family: monospace;
  font-size: 0.85rem;
}

.input-hint {
  font-size: 0.7rem;
  color: var(--text-secondary);
  margin-top: 0.2rem;
  display: block;
}

.db-feedback-msg {
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 700;
}

.db-feedback-msg.success {
  background: #ecfdf5;
  color: #047857;
  border: 1px solid #a7f3d0;
}

.db-feedback-msg.error {
  background: #fef2f2;
  color: #b91c1c;
  border: 1px solid #fecaca;
}

.btn-switch-db {
  width: 100%;
  padding: 0.85rem;
  margin-top: 0.5rem;
}

.text-right { text-align: right; }
</style>
