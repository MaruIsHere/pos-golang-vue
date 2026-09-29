<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-[24px] shadow-sm">
      <div class="flex flex-col gap-1">
        <h2 class="text-2xl font-black text-slate-800 dark:text-slate-100 tracking-tight">Pengaturan Sistem</h2>
        <p class="text-sm font-medium text-slate-500 dark:text-slate-400 mt-1">Konfigurasi profil toko, gambar QRIS pembayaran, tema mode terang/gelap, voucher diskon, dan manajemen staf.</p>
      </div>
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
      <!-- 0. Theme Selection Card -->
      <div class="col-span-1 xl:col-span-2 flex flex-col p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-[24px] shadow-sm">
        <div class="flex justify-between items-center border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <h3 class="flex items-center gap-2 text-lg font-bold text-slate-800 dark:text-slate-100">
            <SunIcon class="w-5 h-5 text-amber-500" />
            <span>Tema Tampilan (Mode Terang / Mode Gelap)</span>
          </h3>
          <span class="px-3 py-1 text-xs font-black rounded-full bg-indigo-100 text-indigo-700 dark:bg-indigo-900/50 dark:text-indigo-400">
            {{ isDarkMode ? '🌙 Dark Mode' : '☀️ Light Mode' }}
          </span>
        </div>

        <div class="flex flex-col gap-4 mt-2">
          <p class="text-sm text-slate-500 dark:text-slate-400 mb-4">Pilih tema warna tampilan antarmuka aplikasi kasir.</p>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div 
              class="flex flex-col gap-3 p-5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-700 rounded-xl cursor-pointer hover:border-indigo-400 hover:shadow-md transition-all relative overflow-hidden" 
              :class="!isDarkMode ? 'border-indigo-500 ring-4 ring-indigo-500/10 shadow-md bg-white dark:bg-slate-800 scale-[1.02]' : 'opacity-60 hover:opacity-100'"
              @click="setDark(false)"
            >
              <div class="w-7 h-7 text-indigo-600 dark:text-indigo-400">
                <SunIcon class="w-7 h-7 text-amber-500" />
              </div>
              <div class="flex flex-col mt-1">
                <h4 class="text-base font-bold text-slate-800 dark:text-slate-100">Mode Terang (Light Mode)</h4>
                <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed mt-1">Tampilan serba putih yang bersih, minimalis, dan terang.</p>
              </div>
            </div>

            <div 
              class="flex flex-col gap-3 p-5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-700 rounded-xl cursor-pointer hover:border-indigo-400 hover:shadow-md transition-all relative overflow-hidden" 
              :class="isDarkMode ? 'border-indigo-500 ring-4 ring-indigo-500/10 shadow-md bg-white dark:bg-slate-800 scale-[1.02]' : 'opacity-60 hover:opacity-100'"
              @click="setDark(true)"
            >
              <div class="w-7 h-7 text-indigo-600 dark:text-indigo-400">
                <MoonIcon class="w-7 h-7 text-indigo-400" />
              </div>
              <div class="flex flex-col mt-1">
                <h4 class="text-base font-bold text-slate-800 dark:text-slate-100">Mode Gelap (Dark Mode)</h4>
                <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed mt-1">Tampilan gelap modern yang nyaman di mata untuk kondisi pencahayaan redup.</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="canEditProfile" class="col-span-1 xl:col-span-2 flex flex-col p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl shadow-sm">
        <div class="border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <h3 class="text-lg font-bold text-slate-800 dark:text-slate-100">Profil Akun</h3>
        </div>

        <form class="grid grid-cols-1 md:grid-cols-[auto_1fr] gap-6" @submit.prevent="saveOwnProfile">
          <div class="flex flex-col items-center gap-3">
            <img v-if="profileForm.profile_photo" :src="profileForm.profile_photo" alt="Foto profil" class="w-24 h-24 rounded-full object-cover border border-slate-200 dark:border-slate-700" />
            <div v-else class="w-24 h-24 rounded-full bg-slate-100 dark:bg-slate-700 flex items-center justify-center text-slate-400">
              <UserGroupIcon class="w-10 h-10" />
            </div>
            <label class="cursor-pointer inline-flex items-center gap-2 px-3 py-2 text-sm font-semibold text-indigo-700 bg-indigo-50 hover:bg-indigo-100 dark:text-indigo-300 dark:bg-indigo-900/30 rounded-lg">
              <ArrowUpTrayIcon class="w-4 h-4" />
              Pilih Foto
              <input type="file" accept="image/*" class="sr-only" @change="onProfilePhotoSelected" />
            </label>
            <button v-if="profileForm.profile_photo" type="button" class="text-xs font-semibold text-rose-600 hover:text-rose-700" @click="profileForm.profile_photo = ''">
              Hapus foto
            </button>
            <span class="text-xs text-slate-500">Format gambar, maksimal 2MB</span>
          </div>

          <div class="flex flex-col gap-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Nama tampilan *</label>
                <AppInput v-model="profileForm.name" type="text" maxlength="100" required />
              </div>
              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Username *</label>
                <AppInput v-model="profileForm.username" type="text" minlength="3" maxlength="100" required />
              </div>
              <div class="flex flex-col gap-2 md:col-span-2">
                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Sandi baru</label>
                <AppInput v-model="profileForm.password" type="password" minlength="6" autocomplete="new-password" placeholder="Kosongkan jika tidak ingin mengganti" />
              </div>
            </div>

            <p v-if="profileError" class="text-sm font-semibold text-rose-600">{{ profileError }}</p>
            <p v-if="profileSuccess" class="text-sm font-semibold text-emerald-600">{{ profileSuccess }}</p>
            <div>
              <AppButton variant="primary" type="submit" :disabled="isSavingProfile">
                {{ isSavingProfile ? 'Menyimpan...' : 'Simpan Profil' }}
              </AppButton>
            </div>
          </div>
        </form>
      </div>

      <!-- 1. Store Profile & QRIS Settings -->
      <div class="flex flex-col p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-[24px] shadow-sm">
        <div class="flex justify-between items-center border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <h3 class="flex items-center gap-2 text-lg font-bold text-slate-800 dark:text-slate-100">
            <BuildingStorefrontIcon class="w-5 h-5 text-indigo-600" />
            <span>Profil Toko & Foto QRIS Pembayaran</span>
          </h3>
        </div>

        <form @submit.prevent="saveStoreSettings" class="flex flex-col gap-4 mt-2">
          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Nama Toko / Usaha *</label>
            <AppInput type="text"  v-model="storeForm.store_name" required />
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Alamat Lengkap *</label>
            <AppInput type="text"  v-model="storeForm.address" required />
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">No. Telepon / WhatsApp *</label>
            <AppInput type="text"  v-model="storeForm.phone" required />
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Pajak % (PPN / Service Charge)</label>
            <AppInput type="number"  v-model.number="storeForm.tax_percentage" min="0" max="100" />
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Potongan Diskon Pelanggan Terdaftar / Member (%)</label>
            <AppInput type="number"  v-model.number="storeForm.member_discount_percentage" min="0" max="100" step="0.5" />
            <span class="input-hint">Potongan ini otomatis dihitung saat kasir memilih/menginput nama pelanggan terdaftar saat checkout.</span>
          </div>

          <div class="flex flex-col gap-2">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Pesan Footer Struk</label>
            <AppInput  rows="2" v-model="storeForm.receipt_footer"></AppInput>
          </div>

          <!-- QRIS Image Section -->
          <div class="flex flex-col gap-2 qris-upload-section">
            <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Foto QRIS Pembayaran Toko</label>
            <p class="input-hint">Upload foto/gambar QRIS resmi toko Anda agar muncul di layar Kasir saat pelanggan memilih metode QRIS.</p>
            
            <div class="qris-preview-box">
              <div v-if="storeForm.qris_image_url" class="qris-img-container">
                <img :src="storeForm.qris_image_url" alt="QRIS Toko" class="qris-preview-img" />
                <AppButton variant="danger"  type="button" class="remove-qris" @click="storeForm.qris_image_url = ''">
                  Hapus Foto QRIS
                </AppButton>
              </div>

              <div v-else class="qris-empty-placeholder">
                <p>Belum ada foto QRIS yang di-upload</p>
              </div>
            </div>

            <div class="qris-upload-actions">
              <label class="cursor-pointer flex flex-1 sm:flex-none items-center justify-center gap-2 px-4 py-2.5 bg-indigo-50 text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400 font-bold rounded-xl hover:bg-indigo-600 hover:text-white transition-colors">
                <ArrowUpTrayIcon class="w-4 h-4" />
                <span>Upload Gambar QRIS</span>
                <input type="file" accept="image/*" @change="onQrisFileSelected" style="display: none;" />
              </label>
              
              <AppInput 
                type="text" 
                 
                v-model="storeForm.qris_image_url" 
                placeholder="Atau paste URL Gambar QRIS..." 
              />
            </div>
          </div>

          <AppButton variant="success"  type="submit" class="success -block" :disabled="isSavingStore">
            {{ isSavingStore ? 'Memproses...' : 'Simpan Pengaturan Toko & QRIS' }}
          </AppButton>
        </form>
      </div>

      <!-- 2. Manage Voucher Codes -->
      <div class="flex flex-col p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-[24px] shadow-sm">
        <div class="flex justify-between items-center border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <h3 class="flex items-center gap-2 text-lg font-bold text-slate-800 dark:text-slate-100">
            <TicketIcon class="w-5 h-5 text-indigo-600" />
            <span>Manajemen Kode Voucher Diskon</span>
          </h3>
        </div>

        <div class="flex flex-col gap-4 mt-2">
          <p class="text-sm text-slate-500 dark:text-slate-400 mb-4">
            Tambah dan kelola kode voucher diskon yang bisa digunakan oleh kasir pada halaman transaksi.
          </p>

          <!-- Add Voucher Form -->
          <form @submit.prevent="createVoucher" class="">
            <h4>Tambah Kode Voucher Baru</h4>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Kode Voucher *</label>
                <AppInput 
                  type="text" 
                   
                  v-model="newVoucher.code" 
                  placeholder="cth: DISKON30" 
                  required 
                />
              </div>

              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Tipe Diskon</label>
                <select class="w-full px-4 py-2 border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 rounded-lg outline-none focus:ring-2 focus:ring-blue-500" v-model="newVoucher.type">
                  <option value="percent">Persen (%)</option>
                  <option value="flat">Potongan Nominal (Rp)</option>
                </select>
              </div>

              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Nilai Potongan *</label>
                <AppInput 
                  type="number" 
                   
                  v-model.number="newVoucher.value" 
                  :placeholder="newVoucher.type === 'percent' ? 'cth: 15 (artinya 15%)' : 'cth: 20000 (artinya Rp 20.000)'" 
                  min="1" 
                  required 
                />
              </div>

              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Deskripsi / Keterangan</label>
                <AppInput 
                  type="text" 
                   
                  v-model="newVoucher.description" 
                  placeholder="cth: Promo Tanggal Kembar" 
                />
              </div>
            </div>

            <AppButton variant="primary"  type="submit" class="primary" :disabled="isCreatingVoucher">
              {{ isCreatingVoucher ? 'Menambahkan...' : 'Simpan Voucher Baru' }}
            </AppButton>
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
                      <AppButton variant="primary"  class="icon -delete-voucher" title="Hapus Voucher" @click="deleteVoucher(v.id ?? 0, v.code)">
                        <TrashIcon class="w-3.5 h-3.5 text-red-600 inline-block mr-1" /> Hapus
                      </AppButton>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- 3. User Management & Registration (Owner & Admin Only) -->
      <div class="col-span-1 xl:col-span-2 flex flex-col p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-[24px] shadow-sm">
        <div class="flex justify-between items-center border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <h3 class="flex items-center gap-2 text-lg font-bold text-slate-800 dark:text-slate-100">
            <UserGroupIcon class="w-5 h-5 text-indigo-600" />
            <span>Manajemen Pengguna & Registrasi Akun Staf</span>
          </h3>
          <span class="active-db-badge sqlite">
            Khusus Owner & Admin
          </span>
        </div>

        <div class="flex flex-col gap-4 mt-2">
          <p class="text-sm text-slate-500 dark:text-slate-400 mb-4">
            Tambah akun pengguna baru (Kasir, Kepala Kasir, Admin, atau Owner) dan kelola daftar staf kasir yang memiliki akses ke sistem POS.
          </p>

          <!-- Register User Form -->
          <form @submit.prevent="createUser" class="mb-6">
            <h4>Tambah Pengguna / Registrasi Staf Baru</h4>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Username *</label>
                <AppInput 
                  type="text" 
                   
                  v-model="newUser.username" 
                  placeholder="Masukkan username" 
                  required 
                />
              </div>

              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Password *</label>
                <AppInput 
                  type="password" 
                   
                  v-model="newUser.password" 
                  placeholder="Minimal 6 karakter" 
                  required 
                />
              </div>

              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Role / Peran Akses *</label>
                <select class="w-full px-4 py-2 border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 rounded-lg outline-none focus:ring-2 focus:ring-blue-500" v-model="newUser.role">
                  <option value="owner">Owner</option>
                  <option value="administrator">Administrator</option>
                  <option value="admin">Admin</option>
                  <option value="kepala_kasir">Kepala Kasir</option>
                  <option value="kasir">Kasir</option>
                </select>
              </div>
            </div>

            <div v-if="userSuccessMsg" class="text-xs text-emerald-600 font-bold mb-3">
              {{ userSuccessMsg }}
            </div>
            <div v-if="userErrorMsg" class="text-xs text-rose-600 font-bold mb-3">
              {{ userErrorMsg }}
            </div>

            <AppButton variant="primary"  type="submit" class="primary" :disabled="isCreatingUser">
              {{ isCreatingUser ? 'Menambahkan...' : 'Tambah Pengguna Baru' }}
            </AppButton>
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
              <AppButton variant="secondary"  type="button" class="secondary" @click="loadUsers">Coba Lagi</AppButton>
            </div>

            <div v-else-if="users.length === 0" class="empty-vouchers">
              <p>Belum ada pengguna terdaftar.</p>
            </div>

            <div v-else class="vouchers-table-wrapper">
              <table class="vouchers-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Nama / Username</th>
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
                      <strong class="text-slate-800 dark:text-slate-100 font-bold">{{ u.name || u.username }}</strong>
                      <span v-if="u.name && u.name !== u.username" class="block text-xs text-slate-500">{{ u.username }}</span>
                    </td>
                    <td>
                      <span class="badge" :class="getUserRoleBadgeClass(u.role)">
                        {{ formatUserRoleLabel(u.role) }}
                      </span>
                    </td>
                    <td class="text-xs text-slate-500">{{ formatDate(u.created_at) }}</td>
                    <td class="text-right">
                      <div class="user-actions">
                        <AppButton v-if="u.id !== Number(authStore.user?.id)" variant="primary"
                          type="button"
                          class="icon -edit-user"
                          :disabled="isLoadingStaffProfile || isSavingStaffProfile"
                          @click="openStaffProfileForm(u)"
                        >
                          <PencilSquareIcon class="w-3.5 h-3.5 inline-block mr-1" /> Edit Profil
                        </AppButton>
                        <AppButton variant="primary" 
                          type="button"
                          class="icon -change-password"
                          :disabled="isChangingStaffPassword"
                          @click="openStaffPasswordForm(u)"
                        >
                          Ganti Sandi
                        </AppButton>
                        <AppButton variant="primary" 
                          type="button"
                          class="icon -delete-voucher"
                          title="Hapus Pengguna"
                          @click="deleteUserAccount(u.id, u.username)"
                        >
                          <TrashIcon class="w-3.5 h-3.5 text-red-600 inline-block mr-1" /> Hapus
                        </AppButton>
                      </div>
                    </td>
                  </tr>
                    <tr v-if="staffProfileTarget?.id === u.id">
                      <td colspan="5">
                        <form class="staff-password-form" @submit.prevent="saveStaffProfile">
                          <div class="staff-password-heading">
                            <strong>Edit profil {{ u.username }}</strong>
                            <span class="text-xs text-slate-500">Nama, username, dan foto profil</span>
                          </div>
                          <div v-if="isLoadingStaffProfile" class="vouchers-loading">Memuat profil pengguna...</div>
                          <div v-else class="grid grid-cols-1 md:grid-cols-[auto_1fr] gap-4">
                            <div class="flex flex-col items-center gap-2">
                              <img v-if="staffProfileForm.profile_photo" :src="staffProfileForm.profile_photo" alt="Foto profil pengguna" class="w-20 h-20 rounded-full object-cover border border-slate-200 dark:border-slate-700" />
                              <div v-else class="w-20 h-20 rounded-full bg-slate-100 dark:bg-slate-700 flex items-center justify-center text-slate-400">
                                <UserGroupIcon class="w-8 h-8" />
                              </div>
                              <label class="cursor-pointer text-xs font-semibold text-indigo-700 dark:text-indigo-300">
                                Pilih foto
                                <input type="file" accept="image/*" class="sr-only" @change="onStaffProfilePhotoSelected" />
                              </label>
                              <button v-if="staffProfileForm.profile_photo" type="button" class="text-xs font-semibold text-rose-600" @click="staffProfileForm.profile_photo = ''">
                                Hapus foto
                              </button>
                              <span class="text-[11px] text-slate-500">Maksimal 2MB</span>
                            </div>
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                              <div class="flex flex-col gap-2">
                                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Nama tampilan *</label>
                                <AppInput v-model="staffProfileForm.name" type="text" maxlength="100" required />
                              </div>
                              <div class="flex flex-col gap-2">
                                <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Username *</label>
                                <AppInput v-model="staffProfileForm.username" type="text" minlength="3" maxlength="100" required />
                              </div>
                            </div>
                          </div>
                          <p v-if="staffProfileError" class="staff-password-error">{{ staffProfileError }}</p>
                          <p v-if="staffProfileSuccess" class="staff-password-success">{{ staffProfileSuccess }}</p>
                          <div class="user-actions">
                            <AppButton variant="secondary" type="button" class="secondary" :disabled="isSavingStaffProfile" @click="cancelStaffProfileForm">
                              Batal
                            </AppButton>
                            <AppButton variant="primary" type="submit" class="primary" :disabled="isLoadingStaffProfile || isSavingStaffProfile">
                              {{ isSavingStaffProfile ? 'Menyimpan...' : 'Simpan Profil' }}
                            </AppButton>
                          </div>
                        </form>
                      </td>
                    </tr>
                    <tr v-if="staffPasswordTarget?.id === u.id">
                      <td colspan="5">
                        <form class="staff-password-form" @submit.prevent="changeStaffPassword">
                          <div class="staff-password-heading">
                            <strong>Ganti sandi untuk {{ staffPasswordTarget?.username }}</strong>
                            <span class="text-xs text-slate-500">Minimal 6 karakter</span>
                          </div>
                          <div class="staff-password-fields">
                            <div class="flex flex-col gap-2">
                              <label class="text-sm font-semibold text-slate-700 dark:text-slate-300" :for="`new-password-${u.id}`">Sandi baru</label>
                              <AppInput
                                :id="`new-password-${u.id}`"
                                v-model="newStaffPassword"
                                
                                type="password"
                                autocomplete="new-password"
                                minlength="6"
                                required
                              />
                            </div>
                            <div class="flex flex-col gap-2">
                              <label class="text-sm font-semibold text-slate-700 dark:text-slate-300" :for="`confirm-password-${u.id}`">Konfirmasi sandi baru</label>
                              <AppInput
                                :id="`confirm-password-${u.id}`"
                                v-model="confirmStaffPassword"
                                
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
                            <AppButton variant="secondary"  type="button" class="secondary" :disabled="isChangingStaffPassword" @click="cancelStaffPasswordChange">
                              Batal
                            </AppButton>
                            <AppButton variant="primary"  type="submit" class="primary" :disabled="isChangingStaffPassword">
                              {{ isChangingStaffPassword ? 'Menyimpan...' : 'Simpan Sandi Baru' }}
                            </AppButton>
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
      <div class="flex flex-col p-6 bg-indigo-50/50 dark:bg-indigo-900/10 border border-indigo-100 dark:border-indigo-800 rounded-[24px] shadow-sm mt-6">
        <div class="flex justify-between items-center border-b border-slate-100 dark:border-slate-700 pb-4 mb-4">
          <h3 class="flex items-center gap-2 text-lg font-bold text-slate-800 dark:text-slate-100">
            <CircleStackIcon class="w-5 h-5 text-indigo-600" />
            <span>Engine Basis Data (Database Switcher)</span>
          </h3>
          <span class="px-3 py-1 text-xs font-bold rounded-full" :class="dbEngine === 'mysql' ? 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400' : 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'">
            Aktif: {{ dbEngine.toUpperCase() }}
          </span>
        </div>

        <div class="flex flex-col gap-4 mt-2">
          <p class="text-sm text-slate-500 dark:text-slate-400 mb-4">
            Pilih engine basis data yang ingin digunakan oleh backend Golang.
          </p>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div 
              class="flex flex-col gap-3 p-5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-700 rounded-xl cursor-pointer hover:border-indigo-400 hover:shadow-md transition-all relative overflow-hidden" 
              :class="{ selected: selectedEngine === 'sqlite' }"
              @click="selectedEngine = 'sqlite'"
            >
              <div class="w-7 h-7 text-indigo-600 dark:text-indigo-400">
                <FolderIcon class="w-7 h-7 text-indigo-600" />
              </div>
              <div class="flex flex-col mt-1">
                <h4>SQLite (Embedded File)</h4>
                <p>Offline-first, tanpa butuh server MySQL. Data disimpan di file <code>pos.db</code>.</p>
              </div>
            </div>

            <div 
              class="flex flex-col gap-3 p-5 bg-slate-50 dark:bg-slate-900/50 border border-slate-200 dark:border-slate-700 rounded-xl cursor-pointer hover:border-indigo-400 hover:shadow-md transition-all relative overflow-hidden" 
              :class="{ selected: selectedEngine === 'mysql' }"
              @click="selectedEngine = 'mysql'"
            >
              <div class="w-7 h-7 text-indigo-600 dark:text-indigo-400">
                <ServerIcon class="w-7 h-7 text-indigo-600" />
              </div>
              <div class="flex flex-col mt-1">
                <h4>MySQL Server</h4>
                <p>Terpusat, cocok untuk multi-kasir di jaringan lokal/server cloud.</p>
              </div>
            </div>
          </div>

          <!-- MySQL DSN Config Form -->
          <div v-if="selectedEngine === 'mysql'" class="mysql-config-box">
            <div class="flex flex-col gap-2">
              <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">MySQL Connection String (DSN)</label>
              <AppInput 
                type="text" 
                 
                v-model="mysqlDsn" 
                placeholder="root:password@tcp(127.0.0.1:3306)/pos_db?charset=utf8mb4&parseTime=True&loc=Local" 
              />
              <span class="input-hint">Format GORM MySQL: user:pass@tcp(host:port)/dbname?params</span>
            </div>
          </div>

          <div v-if="dbMessage" class="db-feedback-msg" :class="dbMessageType">
            {{ dbMessage }}
          </div>

          <AppButton variant="primary"  
            class="primary -switch-db" 
            :disabled="isSwitchingDb"
            @click="switchDatabaseEngine"
          >
            <span v-if="isSwitchingDb">Mengubah Engine & Migrasi Data...</span>
            <span v-else>Terapkan & Switch Database</span>
          </AppButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppInput from '@/components/ui/AppInput.vue';
import { computed, ref, onMounted } from 'vue';
import api from '@/utils/api';
import type { Voucher } from '../types';
import { useTheme } from '../composables/useTheme';
import { useAuthStore } from '../stores/auth';
import { 
  BuildingStorefrontIcon, 
  TicketIcon, 
  CircleStackIcon, 
  FolderIcon, 
  ServerIcon, 
  TrashIcon, 
  PencilSquareIcon,
  ArrowUpTrayIcon,
  SunIcon,
  MoonIcon,
  UserGroupIcon
} from '@heroicons/vue/24/outline';

const emit = defineEmits(['refresh-settings']);

const { isDarkMode, setDark } = useTheme();
const authStore = useAuthStore();
const canEditProfile = computed(() => ['owner', 'admin'].includes(String(authStore.userRole || '').toLowerCase()));
const profileForm = ref({
  username: authStore.user?.username || '',
  name: authStore.user?.name || authStore.user?.username || '',
  profile_photo: authStore.user?.profile_photo || '',
  password: ''
});
const isSavingProfile = ref(false);
const profileError = ref('');
const profileSuccess = ref('');

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

const onProfilePhotoSelected = (event: Event): void => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  if (!file.type.startsWith('image/')) {
    profileError.value = 'File harus berupa gambar.';
    input.value = '';
    return;
  }
  if (file.size > 2 * 1024 * 1024) {
    profileError.value = 'Ukuran foto maksimal 2MB.';
    input.value = '';
    return;
  }

  const reader = new FileReader();
  reader.onload = () => {
    profileForm.value.profile_photo = String(reader.result || '');
    profileError.value = '';
  };
  reader.onerror = () => {
    profileError.value = 'Gagal membaca file foto.';
  };
  reader.readAsDataURL(file);
};

const saveOwnProfile = async (): Promise<void> => {
  profileError.value = '';
  profileSuccess.value = '';
  isSavingProfile.value = true;
  try {
    const { data } = await api.put('/profile', profileForm.value);
    authStore.updateProfile(data.token, data.user);
    profileForm.value.password = '';
    profileSuccess.value = data.message || 'Profil berhasil diperbarui.';
  } catch (err: any) {
    profileError.value = err.response?.data?.error || err.message || 'Gagal menyimpan profil.';
  } finally {
    isSavingProfile.value = false;
  }
};

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
const staffProfileTarget = ref<{ id: number; username: string } | null>(null);
const staffProfileForm = ref({ username: '', name: '', profile_photo: '' });
const isLoadingStaffProfile = ref(false);
const isSavingStaffProfile = ref(false);
const staffProfileError = ref('');
const staffProfileSuccess = ref('');

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

const openStaffProfileForm = async (user: { id: number; username: string }) => {
  cancelStaffProfileForm();
  staffProfileTarget.value = user;
  isLoadingStaffProfile.value = true;
  try {
    const { data } = await api.get(`/users/${user.id}`);
    staffProfileForm.value = {
      username: data.username || '',
      name: data.name || data.username || '',
      profile_photo: data.profile_photo || ''
    };
  } catch (err: any) {
    staffProfileError.value = err.response?.data?.error || err.message || 'Gagal memuat profil pengguna.';
  } finally {
    isLoadingStaffProfile.value = false;
  }
};

const cancelStaffProfileForm = () => {
  staffProfileTarget.value = null;
  staffProfileForm.value = { username: '', name: '', profile_photo: '' };
  staffProfileError.value = '';
  staffProfileSuccess.value = '';
};

const onStaffProfilePhotoSelected = (event: Event): void => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  if (!file.type.startsWith('image/')) {
    staffProfileError.value = 'File harus berupa gambar.';
    input.value = '';
    return;
  }
  if (file.size > 2 * 1024 * 1024) {
    staffProfileError.value = 'Ukuran foto maksimal 2MB.';
    input.value = '';
    return;
  }

  const reader = new FileReader();
  reader.onload = () => {
    staffProfileForm.value.profile_photo = String(reader.result || '');
    staffProfileError.value = '';
  };
  reader.onerror = () => {
    staffProfileError.value = 'Gagal membaca file foto.';
  };
  reader.readAsDataURL(file);
};

const saveStaffProfile = async () => {
  const target = staffProfileTarget.value;
  if (!target) return;
  staffProfileError.value = '';
  staffProfileSuccess.value = '';
  isSavingStaffProfile.value = true;
  try {
    const { data } = await api.put(`/users/${target.id}`, staffProfileForm.value);
    users.value = users.value.map(user => user.id === target.id
      ? { ...user, username: data.user.username, name: data.user.name }
      : user);
    staffProfileTarget.value = { id: data.user.id, username: data.user.username };
    staffProfileSuccess.value = data.message || 'Profil pengguna berhasil diperbarui.';
  } catch (err: any) {
    staffProfileError.value = err.response?.data?.error || err.message || 'Gagal menyimpan profil pengguna.';
  } finally {
    isSavingStaffProfile.value = false;
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


