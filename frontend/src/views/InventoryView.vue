<template>
  <div class="flex flex-col gap-5">
    <div class="flex justify-between items-center p-5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <div class="flex flex-col gap-1">
        <h2 class="text-xl font-bold text-slate-900 dark:text-slate-100">Manajemen Inventoris & Mutasi Stok</h2>
        <p class="text-sm text-slate-500 dark:text-slate-400">Kelola Penerimaan Barang (Restock) dan Pengeluaran Barang (Barang Rusak / Hilang / Expired)</p>
      </div>
    </div>

    <!-- Sub Navigation Tabs -->
    <div class="flex gap-2 overflow-x-auto shrink-0">
      <AppButton variant="secondary" 
        class="" 
        :class="{ active: activeTab === 'receive' }"
        @click="activeTab = 'receive'"
      >
        <ArrowDownTrayIcon class="w-4 h-4 inline-block mr-1.5" />
        <span>Transaksi Penerimaan (Stock In)</span>
      </AppButton>

      <AppButton variant="secondary" 
        class="" 
        :class="{ active: activeTab === 'issue' }"
        @click="activeTab = 'issue'"
      >
        <ArrowUpTrayIcon class="w-4 h-4 inline-block mr-1.5" />
        <span>Transaksi Pengeluaran (Stock Out)</span>
      </AppButton>

      <AppButton variant="secondary" 
        class="" 
        :class="{ active: activeTab === 'history' }"
        @click="activeTab = 'history'"
      >
        <ClockIcon class="w-4 h-4 inline-block mr-1.5" />
        <span>Riwayat Mutasi Stok</span>
      </AppButton>
    </div>

    <!-- Tab 1: Penerimaan Barang (Stock In / Receive) -->
    <div v-if="activeTab === 'receive'" class="flex flex-col gap-4 p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <div class="border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-900 dark:text-slate-100">
        <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Penerimaan Barang dari Supplier / Pabrik</h3>
      </div>

      <form @submit.prevent="submitStockMovement('in')" class="card-body">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="form-group">
            <label class="form-label">Pilih Produk</label>
            <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100 font-medium" v-model.number="receiveForm.product_id" @change="onProductSelect('in')" required>
              <option value="" disabled>-- Pilih Produk --</option>
              <option v-for="p in products" :key="p.id" :value="p.id">
                {{ p.name }}{{ p.artist ? ' [Merk: ' + p.artist + ']' : '' }} (Stok Saat Ini: {{ formatQuantity(p.stock) }} {{ unitLabel(p.unit) }})
              </option>
            </select>
          </div>

          <div class="form-group">
            <div class="flex justify-between items-center">
              <label class="form-label">Kuantitas Masuk (+)</label>
              <button type="button" class="text-xs font-bold text-indigo-600 dark:text-indigo-400 hover:underline flex items-center gap-1" @click="toggleCalc('in')">
                <CalculatorIcon class="w-3.5 h-3.5" />
                {{ showCalcIn ? 'Tutup Kalkulator' : 'Kalkulator Gram & Liter' }}
              </button>
            </div>
            <input
              type="number" 
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100 font-bold" 
              v-model.number="receiveForm.quantity" 
              :min="selectedUnit(receiveForm.product_id) === 'pcs' ? 1 : 0.001"
              :step="selectedUnit(receiveForm.product_id) === 'pcs' ? 1 : 0.001"
              placeholder="cth: 50" 
              required 
            />
            <span class="text-xs text-slate-500 font-semibold">Satuan: {{ unitLabel(selectedUnit(receiveForm.product_id)) }}</span>
          </div>

          <!-- Kalkulator Konversi Stok Khusus Gram & Liter (Receive) -->
          <div v-if="showCalcIn" class="col-span-1 sm:col-span-2 p-4 bg-indigo-50/70 dark:bg-indigo-950/40 border border-indigo-200 dark:border-indigo-800/60 rounded-xl flex flex-col gap-3">
            <div class="flex items-center justify-between border-b border-indigo-200/60 dark:border-indigo-800/60 pb-2">
              <span class="text-xs font-black uppercase text-indigo-800 dark:text-indigo-300 tracking-wider flex items-center gap-1.5">
                <CalculatorIcon class="w-4 h-4 text-indigo-600" />
                Perhitungan Khusus Stok (Gram / Liter / Kemasan)
              </span>
              <div class="flex gap-1.5">
                <button type="button" class="px-2.5 py-1 text-xs font-bold rounded-lg transition-colors" :class="calcModeIn === 'gram' ? 'bg-indigo-600 text-white' : 'bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300'" @click="calcModeIn = 'gram'">Gram (gr)</button>
                <button type="button" class="px-2.5 py-1 text-xs font-bold rounded-lg transition-colors" :class="calcModeIn === 'liter' ? 'bg-indigo-600 text-white' : 'bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300'" @click="calcModeIn = 'liter'">Liter (L)</button>
                <button type="button" class="px-2.5 py-1 text-xs font-bold rounded-lg transition-colors" :class="calcModeIn === 'pack' ? 'bg-indigo-600 text-white' : 'bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300'" @click="calcModeIn = 'pack'">Per Pack / Dus</button>
              </div>
            </div>

            <!-- Mode Gram -->
            <div v-if="calcModeIn === 'gram'" class="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Konversi dari Kg</label>
                <div class="flex items-center gap-1 mt-1">
                  <input type="number" step="0.001" min="0" v-model.number="calcGramKgIn" placeholder="0.5" class="w-full px-3 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
                  <span class="text-xs font-bold text-slate-500">kg</span>
                </div>
                <span class="text-[10px] text-slate-500">= {{ (calcGramKgIn || 0) * 1000 }} gr</span>
              </div>
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Jumlah Kemasan / Bal</label>
                <input type="number" min="0" v-model.number="calcGramPacksIn" placeholder="10" class="w-full px-3 py-1.5 mt-1 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
              </div>
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Isi per Kemasan</label>
                <div class="flex items-center gap-1 mt-1">
                  <input type="number" step="0.01" min="0" v-model.number="calcGramPerPackIn" placeholder="250" class="w-full px-3 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
                  <select v-model="calcGramPackUnitIn" class="text-xs font-bold px-2 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg">
                    <option value="gr">gr</option>
                    <option value="kg">kg</option>
                  </select>
                </div>
              </div>
            </div>

            <!-- Mode Liter -->
            <div v-else-if="calcModeIn === 'liter'" class="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Konversi dari ml</label>
                <div class="flex items-center gap-1 mt-1">
                  <input type="number" step="1" min="0" v-model.number="calcLiterMlIn" placeholder="750" class="w-full px-3 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
                  <span class="text-xs font-bold text-slate-500">ml</span>
                </div>
                <span class="text-[10px] text-slate-500">= {{ ((calcLiterMlIn || 0) / 1000).toFixed(3) }} Liter</span>
              </div>
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Jumlah Botol / Jerigen</label>
                <input type="number" min="0" v-model.number="calcLiterPacksIn" placeholder="4" class="w-full px-3 py-1.5 mt-1 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
              </div>
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Volume per Botol/Wadah</label>
                <div class="flex items-center gap-1 mt-1">
                  <input type="number" step="0.01" min="0" v-model.number="calcLiterPerPackIn" placeholder="5" class="w-full px-3 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
                  <select v-model="calcLiterPackUnitIn" class="text-xs font-bold px-2 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg">
                    <option value="L">Liter (L)</option>
                    <option value="ml">ml</option>
                  </select>
                </div>
              </div>
            </div>

            <!-- Mode Pack / Dus -->
            <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Jumlah Dus / Karton</label>
                <input type="number" min="0" v-model.number="calcPackBoxesIn" placeholder="5" class="w-full px-3 py-1.5 mt-1 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
              </div>
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Isi Pcs per Dus</label>
                <input type="number" min="0" v-model.number="calcPackPcsPerBoxIn" placeholder="24" class="w-full px-3 py-1.5 mt-1 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
              </div>
            </div>

            <div class="flex items-center justify-between pt-2 border-t border-indigo-200/60 dark:border-indigo-800/60">
              <span class="text-xs font-bold text-indigo-900 dark:text-indigo-200">
                Hasil Perhitungan: <strong class="text-sm font-black text-indigo-700 dark:text-indigo-300">{{ formatQuantity(calculatedResultIn) }} {{ unitLabel(selectedUnit(receiveForm.product_id)) }}</strong>
              </span>
              <button type="button" class="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-bold transition-colors shadow-sm" @click="applyCalcResult('in')">
                Gunakan Hasil ini ke Kuantitas
              </button>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Alasan / Sumber</label>
            <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="receiveForm.reason">
              <option value="pembelian_supplier">Pembelian dari Supplier</option>
              <option value="produksi_sendiri">Hasil Produksi Mandiri</option>
              <option value="penyesuaian_opname">Penyesuaian Stok Opname (+)</option>
            </select>
          </div>

          <div class="form-group col-span-1 sm:col-span-2">
            <label class="form-label">Catatan / No. Surat Jalan / Supplier</label>
            <AppInput 
              type="text" 
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" 
              v-model="receiveForm.notes" 
              placeholder="cth: Restock PT Maju Bersama - SJ #9921" 
            />
          </div>
        </div>

        <AppButton variant="primary" type="submit" class="" :disabled="isSubmitting">
          {{ isSubmitting ? 'Memproses...' : 'Simpan Penerimaan Barang & Tambah Stok' }}
        </AppButton>
      </form>
    </div>

    <!-- Tab 2: Pengeluaran Barang (Stock Out / Issue) -->
    <div v-else-if="activeTab === 'issue'" class="flex flex-col gap-4 p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <div class="border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-900 dark:text-slate-100">
        <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Pengeluaran Barang Non-Penjualan (Rusak / Hilang / Promosi)</h3>
      </div>

      <form @submit.prevent="submitStockMovement('out')" class="card-body">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="form-group">
            <label class="form-label">Pilih Produk</label>
            <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100 font-medium" v-model.number="issueForm.product_id" @change="onProductSelect('out')" required>
              <option value="" disabled>-- Pilih Produk --</option>
              <option v-for="p in products" :key="p.id" :value="p.id">
                {{ p.name }}{{ p.artist ? ' [Merk: ' + p.artist + ']' : '' }} (Stok Tersedia: {{ formatQuantity(p.stock) }} {{ unitLabel(p.unit) }})
              </option>
            </select>
          </div>

          <div class="form-group">
            <div class="flex justify-between items-center">
              <label class="form-label">Kuantitas Keluar (-)</label>
              <button type="button" class="text-xs font-bold text-indigo-600 dark:text-indigo-400 hover:underline flex items-center gap-1" @click="toggleCalc('out')">
                <CalculatorIcon class="w-3.5 h-3.5" />
                {{ showCalcOut ? 'Tutup Kalkulator' : 'Kalkulator Gram & Liter' }}
              </button>
            </div>
            <input
              type="number" 
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100 font-bold" 
              v-model.number="issueForm.quantity" 
              :min="selectedUnit(issueForm.product_id) === 'pcs' ? 1 : 0.001"
              :step="selectedUnit(issueForm.product_id) === 'pcs' ? 1 : 0.001"
              placeholder="cth: 5" 
              required 
            />
            <span class="text-xs text-slate-500 font-semibold">Satuan: {{ unitLabel(selectedUnit(issueForm.product_id)) }}</span>
          </div>

          <!-- Kalkulator Konversi Stok Khusus Gram & Liter (Issue) -->
          <div v-if="showCalcOut" class="col-span-1 sm:col-span-2 p-4 bg-indigo-50/70 dark:bg-indigo-950/40 border border-indigo-200 dark:border-indigo-800/60 rounded-xl flex flex-col gap-3">
            <div class="flex items-center justify-between border-b border-indigo-200/60 dark:border-indigo-800/60 pb-2">
              <span class="text-xs font-black uppercase text-indigo-800 dark:text-indigo-300 tracking-wider flex items-center gap-1.5">
                <CalculatorIcon class="w-4 h-4 text-indigo-600" />
                Perhitungan Khusus Stok (Gram / Liter / Kemasan)
              </span>
              <div class="flex gap-1.5">
                <button type="button" class="px-2.5 py-1 text-xs font-bold rounded-lg transition-colors" :class="calcModeOut === 'gram' ? 'bg-indigo-600 text-white' : 'bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300'" @click="calcModeOut = 'gram'">Gram (gr)</button>
                <button type="button" class="px-2.5 py-1 text-xs font-bold rounded-lg transition-colors" :class="calcModeOut === 'liter' ? 'bg-indigo-600 text-white' : 'bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300'" @click="calcModeOut = 'liter'">Liter (L)</button>
                <button type="button" class="px-2.5 py-1 text-xs font-bold rounded-lg transition-colors" :class="calcModeOut === 'pack' ? 'bg-indigo-600 text-white' : 'bg-white dark:bg-slate-800 text-slate-700 dark:text-slate-300'" @click="calcModeOut = 'pack'">Per Pack / Dus</button>
              </div>
            </div>

            <!-- Mode Gram -->
            <div v-if="calcModeOut === 'gram'" class="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Konversi dari Kg</label>
                <div class="flex items-center gap-1 mt-1">
                  <input type="number" step="0.001" min="0" v-model.number="calcGramKgOut" placeholder="0.5" class="w-full px-3 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
                  <span class="text-xs font-bold text-slate-500">kg</span>
                </div>
                <span class="text-[10px] text-slate-500">= {{ (calcGramKgOut || 0) * 1000 }} gr</span>
              </div>
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Jumlah Kemasan / Bal</label>
                <input type="number" min="0" v-model.number="calcGramPacksOut" placeholder="10" class="w-full px-3 py-1.5 mt-1 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
              </div>
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Isi per Kemasan</label>
                <div class="flex items-center gap-1 mt-1">
                  <input type="number" step="0.01" min="0" v-model.number="calcGramPerPackOut" placeholder="250" class="w-full px-3 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
                  <select v-model="calcGramPackUnitOut" class="text-xs font-bold px-2 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg">
                    <option value="gr">gr</option>
                    <option value="kg">kg</option>
                  </select>
                </div>
              </div>
            </div>

            <!-- Mode Liter -->
            <div v-else-if="calcModeOut === 'liter'" class="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Konversi dari ml</label>
                <div class="flex items-center gap-1 mt-1">
                  <input type="number" step="1" min="0" v-model.number="calcLiterMlOut" placeholder="750" class="w-full px-3 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
                  <span class="text-xs font-bold text-slate-500">ml</span>
                </div>
                <span class="text-[10px] text-slate-500">= {{ ((calcLiterMlOut || 0) / 1000).toFixed(3) }} Liter</span>
              </div>
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Jumlah Botol / Jerigen</label>
                <input type="number" min="0" v-model.number="calcLiterPacksOut" placeholder="4" class="w-full px-3 py-1.5 mt-1 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
              </div>
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Volume per Botol/Wadah</label>
                <div class="flex items-center gap-1 mt-1">
                  <input type="number" step="0.01" min="0" v-model.number="calcLiterPerPackOut" placeholder="5" class="w-full px-3 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
                  <select v-model="calcLiterPackUnitOut" class="text-xs font-bold px-2 py-1.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg">
                    <option value="L">Liter (L)</option>
                    <option value="ml">ml</option>
                  </select>
                </div>
              </div>
            </div>

            <!-- Mode Pack / Dus -->
            <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Jumlah Dus / Karton</label>
                <input type="number" min="0" v-model.number="calcPackBoxesOut" placeholder="5" class="w-full px-3 py-1.5 mt-1 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
              </div>
              <div>
                <label class="text-xs font-bold text-slate-600 dark:text-slate-300">Isi Pcs per Dus</label>
                <input type="number" min="0" v-model.number="calcPackPcsPerBoxOut" placeholder="24" class="w-full px-3 py-1.5 mt-1 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-sm font-bold" />
              </div>
            </div>

            <div class="flex items-center justify-between pt-2 border-t border-indigo-200/60 dark:border-indigo-800/60">
              <span class="text-xs font-bold text-indigo-900 dark:text-indigo-200">
                Hasil Perhitungan: <strong class="text-sm font-black text-indigo-700 dark:text-indigo-300">{{ formatQuantity(calculatedResultOut) }} {{ unitLabel(selectedUnit(issueForm.product_id)) }}</strong>
              </span>
              <button type="button" class="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-bold transition-colors shadow-sm" @click="applyCalcResult('out')">
                Gunakan Hasil ini ke Kuantitas
              </button>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Alasan Pengeluaran</label>
            <select class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" v-model="issueForm.reason">
              <option value="barang_rusak">Barang Rusak / Cacat</option>
              <option value="barang_hilang">Barang Hilang / Selisih Stok</option>
              <option value="expired">Kadaluarsa (Expired)</option>
              <option value="promosi">Kebutuhan Promosi / Sampling</option>
              <option value="pemakaian_sendiri">Pemakaian Internal Toko</option>
            </select>
          </div>

          <div class="form-group col-span-1 sm:col-span-2">
            <label class="form-label">Catatan Pengeluaran</label>
            <AppInput 
              type="text" 
              class="w-full px-3 py-2 bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100" 
              v-model="issueForm.notes" 
              placeholder="cth: Kemasan pecah saat pemindahan di gudang" 
            />
          </div>
        </div>

        <AppButton variant="danger" type="submit" class="" :disabled="isSubmitting">
          {{ isSubmitting ? 'Memproses...' : 'Simpan Pengeluaran Barang & Potong Stok' }}
        </AppButton>
      </form>
    </div>

    <!-- Tab 3: Riwayat Mutasi Stok -->
    <div v-else class="flex flex-col gap-4 p-6 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl backdrop-blur-md bg-white/90 dark:bg-slate-900/90">
      <div class="border-b border-slate-200 dark:border-slate-700 pb-3 text-lg font-bold text-slate-900 dark:text-slate-100">
        <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Log Riwayat Mutasi Masuk & Keluar Stok</h3>
      </div>

      <div class="card-body">
        <div v-if="isLoadingMovements" class="p-12 text-center text-slate-500">
          <div class="w-8 h-8 border-4 border-slate-200 dark:border-slate-700 border-t-indigo-600 rounded-full animate-spin mx-auto mb-2"></div>
          <p class="text-sm text-slate-500 dark:text-slate-400">Memuat data riwayat mutasi stok...</p>
        </div>

        <div v-else-if="movements.length === 0" class="p-12 text-center text-slate-500">
          <p class="text-sm text-slate-500 dark:text-slate-400">Belum ada riwayat mutasi stok barang.</p>
        </div>

        <div v-else class="overflow-x-auto p-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl">
          <table class="data-table min-w-[760px] text-sm">
            <thead>
              <tr>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Waktu & Tanggal</th>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Jenis Mutasi</th>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Produk</th>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Jumlah</th>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Alasan</th>
                <th class="p-3 font-semibold text-slate-500 dark:text-slate-400 border-b border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900 whitespace-nowrap">Catatan</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="m in movements" :key="m.id">
                <td class="text-xs text-slate-500">{{ formatDate(m.created_at) }}</td>
                <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
                  <span 
                    class="px-2.5 py-1 rounded-full text-[10px] uppercase tracking-wider font-extrabold" 
                    :class="m.type === 'in' ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400' : 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-400'"
                  >
                    {{ m.type === 'in' ? 'Masuk (In)' : 'Keluar (Out)' }}
                  </span>
                </td>
                <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
                  <strong class="product-name">{{ m.product ? m.product.name : 'Produk ID ' + m.product_id }}</strong>
                  <span v-if="m.product && m.product.artist" class="text-xs text-slate-400 ml-1.5">[Merk: {{ m.product.artist }}]</span>
                </td>
                <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
                  <span class="font-bold text-sm" :class="m.type">
                    {{ m.type === 'in' ? '+' : '-' }}{{ formatQuantity(m.quantity) }} {{ unitLabel(m.unit) }}
                  </span>
                </td>
                <td class="p-3 border-b border-slate-200 dark:border-slate-700 align-middle text-slate-900 dark:text-slate-100">
                  <span class="px-2 py-0.5 rounded text-xs bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400">{{ formatReason(m.reason ?? '') }}</span>
                </td>
                <td class="text-xs text-slate-500">{{ m.notes || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppButton from '@/components/ui/AppButton.vue';
import AppInput from '@/components/ui/AppInput.vue';

import { ref, computed, onMounted, watch } from 'vue';
import api from '@/utils/api';
import type { Product, StockMovement } from '../types';
import { ArrowDownTrayIcon, ArrowUpTrayIcon, ClockIcon, CalculatorIcon } from '@heroicons/vue/24/outline';
import { showAppAlert } from '@/composables/useAppDialog';
import { useStoreContextStore } from '@/stores/storeContext';

const storeContextStore = useStoreContextStore();
const activeTab = ref('receive');
const products = ref<Product[]>([]);
const movements = ref<StockMovement[]>([]);
const isLoadingMovements = ref(false);
const isSubmitting = ref(false);

interface StockForm {
  product_id: string | number;
  quantity: number;
  reason: string;
  notes: string;
}

const receiveForm = ref<StockForm>({
  product_id: '',
  quantity: 1,
  reason: 'pembelian_supplier',
  notes: ''
});

const issueForm = ref<StockForm>({
  product_id: '',
  quantity: 1,
  reason: 'barang_rusak',
  notes: ''
});

// State Kalkulator Konversi Stok (In)
const showCalcIn = ref(false);
const calcModeIn = ref<'gram' | 'liter' | 'pack'>('gram');
const calcGramKgIn = ref<number | null>(null);
const calcGramPacksIn = ref<number | null>(null);
const calcGramPerPackIn = ref<number | null>(null);
const calcGramPackUnitIn = ref<'gr' | 'kg'>('gr');
const calcLiterMlIn = ref<number | null>(null);
const calcLiterPacksIn = ref<number | null>(null);
const calcLiterPerPackIn = ref<number | null>(null);
const calcLiterPackUnitIn = ref<'L' | 'ml'>('L');
const calcPackBoxesIn = ref<number | null>(null);
const calcPackPcsPerBoxIn = ref<number | null>(null);

// State Kalkulator Konversi Stok (Out)
const showCalcOut = ref(false);
const calcModeOut = ref<'gram' | 'liter' | 'pack'>('gram');
const calcGramKgOut = ref<number | null>(null);
const calcGramPacksOut = ref<number | null>(null);
const calcGramPerPackOut = ref<number | null>(null);
const calcGramPackUnitOut = ref<'gr' | 'kg'>('gr');
const calcLiterMlOut = ref<number | null>(null);
const calcLiterPacksOut = ref<number | null>(null);
const calcLiterPerPackOut = ref<number | null>(null);
const calcLiterPackUnitOut = ref<'L' | 'ml'>('L');
const calcPackBoxesOut = ref<number | null>(null);
const calcPackPcsPerBoxOut = ref<number | null>(null);

const formatDate = (str?: string): string => {
  if (!str) return '-';
  return new Date(str).toLocaleString('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit'
  });
};
const formatQuantity = (quantity: number): string => new Intl.NumberFormat('id-ID', { maximumFractionDigits: 3 }).format(quantity || 0);
const unitLabel = (unit?: Product['unit']): string => unit === 'gram' ? 'gr' : unit === 'liter' ? 'L' : 'pcs';
const selectedUnit = (productId: string | number): Product['unit'] => products.value.find(product => product.id === Number(productId))?.unit ?? 'pcs';

const toggleCalc = (formType: 'in' | 'out') => {
  if (formType === 'in') {
    showCalcIn.value = !showCalcIn.value;
  } else {
    showCalcOut.value = !showCalcOut.value;
  }
};

const onProductSelect = (formType: 'in' | 'out') => {
  const productId = formType === 'in' ? receiveForm.value.product_id : issueForm.value.product_id;
  const unit = selectedUnit(productId);
  if (unit === 'gram') {
    if (formType === 'in') calcModeIn.value = 'gram';
    else calcModeOut.value = 'gram';
  } else if (unit === 'liter') {
    if (formType === 'in') calcModeIn.value = 'liter';
    else calcModeOut.value = 'liter';
  } else {
    if (formType === 'in') calcModeIn.value = 'pack';
    else calcModeOut.value = 'pack';
  }
};

const calculatedResultIn = computed((): number => {
  if (calcModeIn.value === 'gram') {
    if (calcGramKgIn.value) return (calcGramKgIn.value || 0) * 1000;
    if (calcGramPacksIn.value && calcGramPerPackIn.value) {
      const perPack = calcGramPackUnitIn.value === 'kg' ? (calcGramPerPackIn.value * 1000) : calcGramPerPackIn.value;
      return (calcGramPacksIn.value || 0) * perPack;
    }
  } else if (calcModeIn.value === 'liter') {
    if (calcLiterMlIn.value) return (calcLiterMlIn.value || 0) / 1000;
    if (calcLiterPacksIn.value && calcLiterPerPackIn.value) {
      const perPack = calcLiterPackUnitIn.value === 'ml' ? (calcLiterPerPackIn.value / 1000) : calcLiterPerPackIn.value;
      return (calcLiterPacksIn.value || 0) * perPack;
    }
  } else if (calcModeIn.value === 'pack') {
    if (calcPackBoxesIn.value && calcPackPcsPerBoxIn.value) {
      return (calcPackBoxesIn.value || 0) * (calcPackPcsPerBoxIn.value || 0);
    }
  }
  return 0;
});

const calculatedResultOut = computed((): number => {
  if (calcModeOut.value === 'gram') {
    if (calcGramKgOut.value) return (calcGramKgOut.value || 0) * 1000;
    if (calcGramPacksOut.value && calcGramPerPackOut.value) {
      const perPack = calcGramPackUnitOut.value === 'kg' ? (calcGramPerPackOut.value * 1000) : calcGramPerPackOut.value;
      return (calcGramPacksOut.value || 0) * perPack;
    }
  } else if (calcModeOut.value === 'liter') {
    if (calcLiterMlOut.value) return (calcLiterMlOut.value || 0) / 1000;
    if (calcLiterPacksOut.value && calcLiterPerPackOut.value) {
      const perPack = calcLiterPackUnitOut.value === 'ml' ? (calcLiterPerPackOut.value / 1000) : calcLiterPerPackOut.value;
      return (calcLiterPacksOut.value || 0) * perPack;
    }
  } else if (calcModeOut.value === 'pack') {
    if (calcPackBoxesOut.value && calcPackPcsPerBoxOut.value) {
      return (calcPackBoxesOut.value || 0) * (calcPackPcsPerBoxOut.value || 0);
    }
  }
  return 0;
});

const applyCalcResult = (formType: 'in' | 'out') => {
  if (formType === 'in') {
    if (calculatedResultIn.value > 0) {
      receiveForm.value.quantity = calculatedResultIn.value;
    }
  } else {
    if (calculatedResultOut.value > 0) {
      issueForm.value.quantity = calculatedResultOut.value;
    }
  }
};

const formatReason = (reason: string): string => {
  const map: Record<string, string> = {
    pembelian_supplier: 'Pembelian Supplier',
    produksi_sendiri: 'Produksi Mandiri',
    penyesuaian_opname: 'Penyesuaian Opname',
    barang_rusak: 'Barang Rusak',
    barang_hilang: 'Barang Hilang',
    expired: 'Kadaluarsa (Expired)',
    promosi: 'Promosi / Sampling',
    pemakaian_sendiri: 'Pemakaian Internal',
    retur_penjualan: 'Retur Penjualan'
  };
  return map[reason] || reason;
};

const fetchProducts = async () => {
  try {
    const params: Record<string, any> = { is_master: false };
    if (storeContextStore.activeStoreId) {
      params.store_id = storeContextStore.activeStoreId;
    }
    const res = await api.get('/products', { params });
    products.value = res.data;
  } catch (err: any) {
    console.error(err.response?.data?.error || err.message || 'Error occurred');
  }
};

const fetchMovements = async () => {
  isLoadingMovements.value = true;
  try {
    const res = await api.get('/stock-movements');
    movements.value = res.data;
  } catch (err: any) {
    console.error(err.response?.data?.error || err.message || 'Error occurred');
  } finally {
    isLoadingMovements.value = false;
  }
};

onMounted(() => {
  storeContextStore.fetchStores();
  fetchProducts();
  fetchMovements();
});

watch(() => storeContextStore.activeStoreId, () => {
  fetchProducts();
});

const submitStockMovement = async (type: string): Promise<void> => {
  const form = type === 'in' ? receiveForm.value : issueForm.value;
  if (!form.product_id || form.quantity <= 0) return;

  isSubmitting.value = true;
  try {
    const payload = {
      product_id: form.product_id,
      type,
      quantity: form.quantity,
      reason: form.reason,
      notes: form.notes
    };

    await api.post('/stock-movements', payload);
    await showAppAlert(`Berhasil menyimpan transaksi ${type === 'in' ? 'penerimaan' : 'pengeluaran'} barang!`, 'success');
    if (type === 'in') {
      receiveForm.value = { product_id: '', quantity: 1, reason: 'pembelian_supplier', notes: '' };
      showCalcIn.value = false;
    } else {
      issueForm.value = { product_id: '', quantity: 1, reason: 'barang_rusak', notes: '' };
      showCalcOut.value = false;
    }
    fetchProducts();
    fetchMovements();
    activeTab.value = 'history';
  } catch (err: any) {
    const errMsg = err.response?.data?.error || err.message || 'Error occurred';
    await showAppAlert('Koneksi error: ' + errMsg, 'error');
  } finally {
    isSubmitting.value = false;
  }
};
</script>

