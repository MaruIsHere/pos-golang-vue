<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" @click.self="$emit('close')">
    <div class="w-full max-w-lg backdrop-blur-md bg-white/90 dark:bg-slate-900/90 rounded-xl shadow-xl flex flex-col overflow-hidden">
      <!-- Header -->
      <div class="flex items-center justify-between p-5 border-b border-slate-200 dark:border-slate-700">
        <div>
          <h3 class="text-lg font-bold text-slate-900 dark:text-slate-100">Pembayaran Kasir</h3>
          <p class="text-sm text-slate-500 dark:text-slate-400">Pilih metode pembayaran dan selesaikan transaksi</p>
        </div>
        <button class="p-1.5 rounded-md text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-slate-900 dark:hover:text-slate-100 transition-colors" @click="$emit('close')">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>

      <!-- Body -->
      <div class="p-5 flex flex-col gap-4 max-h-[80vh] overflow-y-auto">
        <!-- Total Display -->
        <div class="bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl p-4 text-center">
          <div class="flex flex-col gap-1 items-center">
            <span class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wide">Total Yang Harus Dibayar</span>
            <div v-if="isMatchedMember" class="flex items-center gap-1">
              <AppBadge variant="success">
                <CheckBadgeIcon class="w-4 h-4 mr-1" />
                Diskon Member {{ memberDiscountPercent }}% (-Rp {{ formatPrice(memberDiscountAmount) }})
              </AppBadge>
            </div>
          </div>
          <div class="mt-2">
            <h2 class="text-3xl font-extrabold text-blue-600 dark:text-blue-400">Rp {{ formatPrice(finalGrandTotal) }}</h2>
            <span v-if="isMatchedMember" class="text-sm text-slate-500 line-through">Semula: Rp {{ formatPrice(grandTotal) }}</span>
          </div>
        </div>

        <!-- Customer Name / Member Dropdown -->
        <div class="flex flex-col gap-2">
          <label class="flex justify-between items-center text-sm font-medium text-slate-700 dark:text-slate-300">
            <span>Nama Pelanggan / Member</span>
            <AppBadge v-if="isMatchedMember" variant="success">
              <CheckBadgeIcon class="w-3.5 h-3.5 mr-1" />
              Diskon {{ memberDiscountPercent }}% Ditambahkan
            </AppBadge>
          </label>
          <div class="flex flex-col gap-2">
            <select class="w-full px-3 py-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg text-slate-900 dark:text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500" v-model="selectedCustomerOption" @change="onCustomerSelect">
              <option value="Umum">Umum (Non-Member)</option>
              <option v-for="c in customersList" :key="c.id" :value="c.name">
                {{ c.name }} {{ c.phone ? '(' + c.phone + ')' : '' }} (Diskon Member {{ memberDiscountPercent }}%)
              </option>
              <option value="custom">-- Ketik Nama Manual --</option>
            </select>

            <AppInput 
              v-if="selectedCustomerOption === 'custom'" 
              type="text" 
              v-model="customerName" 
              placeholder="Ketik Nama Pelanggan..." 
            />
          </div>
          <p v-if="isMatchedMember" class="text-xs text-emerald-600 dark:text-emerald-400 mt-1">
            Pelanggan "{{ matchedMember?.name }}" terdaftar di database! Potongan member {{ memberDiscountPercent }}% (-Rp {{ formatPrice(memberDiscountAmount) }}) otomatis diterapkan.
          </p>
        </div>

        <!-- Payment Method Grid -->
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium text-slate-700 dark:text-slate-300">Metode Pembayaran</label>
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-2.5">
            <button 
              v-for="m in paymentMethods" 
              :key="m.id" 
              class="flex flex-col items-center gap-1.5 p-3 rounded-lg border transition-all text-center" 
              :class="[
                paymentMethod === m.id 
                  ? 'bg-indigo-50 dark:bg-indigo-900/30 border-indigo-500 text-indigo-600 dark:text-indigo-400 ring-2 ring-indigo-500/20' 
                  : 'bg-slate-50 dark:bg-slate-800/50 border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800'
              ]"
              @click="selectMethod(m.id)"
            >
              <component :is="m.iconComp" class="w-5 h-5" />
              <div class="flex flex-col">
                <span class="text-xs font-bold">{{ m.name }}</span>
                <span class="text-[10px] opacity-70">{{ m.desc }}</span>
              </div>
            </button>
          </div>
        </div>

        <!-- 1. CASH SECTION -->
        <div v-if="paymentMethod === 'cash'" class="bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg p-4 flex flex-col gap-3">
          <div class="flex flex-col gap-2">
            <label class="text-sm font-medium text-slate-700 dark:text-slate-300">Nominal Uang Diterima</label>
            <div class="relative flex items-center">
              <span class="absolute left-4 font-bold text-blue-600 dark:text-blue-400 z-10">Rp</span>
              <AppInput 
                type="number" 
                class="w-full text-lg font-bold"
                style="padding-left: 2.75rem;"
                v-model.number="paidAmount" 
                placeholder="0" 
              />
            </div>
          </div>

          <!-- Quick Cash Buttons -->
          <div class="grid grid-cols-2 gap-2 mt-2">
            <AppButton variant="success" class="col-span-2 !bg-emerald-50 dark:!bg-emerald-900/20 !border-emerald-200 dark:!border-emerald-800 !text-emerald-600 dark:!text-emerald-400 hover:!bg-emerald-100 dark:hover:!bg-emerald-900/40" @click="setExactCash">
              <BoltIcon class="w-4 h-4 mr-1" /> Uang Pas (Rp {{ formatPrice(grandTotal) }})
            </AppButton>
            <AppButton 
              v-for="amount in quickCashAmounts" 
              :key="amount" 
              variant="secondary"
              @click="setPaidAmount(amount)"
            >
              Rp {{ formatPrice(amount) }}
            </AppButton>
          </div>

          <!-- Change Amount Display -->
          <div class="flex justify-between items-center p-3 rounded-lg mt-2 font-bold text-sm border"
               :class="changeAmount >= 0 ? 'bg-emerald-50 dark:bg-emerald-900/20 border-emerald-200 dark:border-emerald-800 text-emerald-600 dark:text-emerald-400' : 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800 text-red-600 dark:text-red-400'">
            <span>{{ changeAmount >= 0 ? 'Kembalian' : 'Uang Kurang' }}</span>
            <span class="text-lg">
              Rp {{ formatPrice(Math.abs(changeAmount)) }}
            </span>
          </div>
        </div>

        <!-- 2. QRIS SECTION -->
        <div v-else-if="paymentMethod === 'qris'" class="bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg p-4 flex flex-col items-center gap-3 text-center">
          <div class="w-full flex justify-between items-center">
            <div class="flex flex-col items-start">
              <span class="text-xl font-black text-red-600 tracking-wide">QRIS</span>
              <span class="text-[10px] text-slate-500">National QR Standard</span>
            </div>
            <span class="text-xs font-bold text-emerald-500 bg-emerald-50 dark:bg-emerald-900/20 px-2.5 py-1 rounded-full border border-emerald-200 dark:border-emerald-800">Standby Scan</span>
          </div>

          <!-- Custom QRIS Image or SVG QR Graphic -->
          <div class="relative p-2.5 bg-white border border-slate-200 rounded-xl shadow-sm overflow-hidden">
            <img 
              v-if="effectiveQrisUrl" 
              :src="effectiveQrisUrl" 
              alt="Foto QRIS Toko" 
              class="w-40 h-40 object-contain rounded-lg" 
            />
            <svg v-else class="w-40 h-40" viewBox="0 0 200 200">
              <!-- Background -->
              <rect width="200" height="200" fill="#ffffff" rx="12"/>
              <!-- QR Position Detection Patterns -->
              <!-- Top Left -->
              <rect x="15" y="15" width="45" height="45" fill="#0f172a" rx="4"/>
              <rect x="23" y="23" width="29" height="29" fill="#ffffff" rx="2"/>
              <rect x="30" y="30" width="15" height="15" fill="#4f46e5" rx="2"/>
              <!-- Top Right -->
              <rect x="140" y="15" width="45" height="45" fill="#0f172a" rx="4"/>
              <rect x="148" y="23" width="29" height="29" fill="#ffffff" rx="2"/>
              <rect x="155" y="30" width="15" height="15" fill="#4f46e5" rx="2"/>
              <!-- Bottom Left -->
              <rect x="15" y="140" width="45" height="45" fill="#0f172a" rx="4"/>
              <rect x="23" y="148" width="29" height="29" fill="#ffffff" rx="2"/>
              <rect x="30" y="155" width="15" height="15" fill="#4f46e5" rx="2"/>
              <!-- QR Data Matrix Dots Simulation -->
              <rect x="70" y="20" width="12" height="12" fill="#0f172a"/>
              <rect x="90" y="20" width="12" height="12" fill="#4f46e5"/>
              <rect x="110" y="20" width="12" height="12" fill="#0f172a"/>
              <rect x="70" y="40" width="12" height="12" fill="#0f172a"/>
              <rect x="110" y="40" width="12" height="12" fill="#4f46e5"/>
              <rect x="20" y="70" width="12" height="12" fill="#4f46e5"/>
              <rect x="40" y="70" width="12" height="12" fill="#0f172a"/>
              <rect x="70" y="70" width="12" height="12" fill="#0f172a"/>
              <rect x="90" y="70" width="20" height="20" fill="#4f46e5" rx="4"/>
              <rect x="120" y="70" width="12" height="12" fill="#0f172a"/>
              <rect x="140" y="70" width="12" height="12" fill="#4f46e5"/>
              <rect x="160" y="70" width="12" height="12" fill="#0f172a"/>
              <rect x="20" y="90" width="12" height="12" fill="#0f172a"/>
              <rect x="40" y="90" width="12" height="12" fill="#4f46e5"/>
              <rect x="140" y="90" width="12" height="12" fill="#0f172a"/>
              <rect x="170" y="90" width="12" height="12" fill="#4f46e5"/>
              <rect x="20" y="110" width="12" height="12" fill="#4f46e5"/>
              <rect x="70" y="110" width="12" height="12" fill="#0f172a"/>
              <rect x="90" y="110" width="12" height="12" fill="#4f46e5"/>
              <rect x="110" y="110" width="12" height="12" fill="#0f172a"/>
              <rect x="150" y="110" width="12" height="12" fill="#4f46e5"/>
              <rect x="70" y="140" width="12" height="12" fill="#0f172a"/>
              <rect x="90" y="140" width="12" height="12" fill="#4f46e5"/>
              <rect x="110" y="140" width="12" height="12" fill="#0f172a"/>
              <rect x="140" y="140" width="12" height="12" fill="#4f46e5"/>
              <rect x="160" y="140" width="12" height="12" fill="#0f172a"/>
              <rect x="70" y="165" width="12" height="12" fill="#4f46e5"/>
              <rect x="90" y="165" width="12" height="12" fill="#0f172a"/>
              <rect x="120" y="165" width="12" height="12" fill="#4f46e5"/>
              <rect x="150" y="165" width="20" height="20" fill="#0f172a" rx="4"/>
            </svg>

            <!-- Using simple keyframes for inline style -->
            <component is="style">
              @keyframes scanMove {
                0% { top: 10px; }
                50% { top: 155px; }
                100% { top: 10px; }
              }
            </component>
            <div class="absolute left-2 right-2 h-[3px] bg-gradient-to-r from-transparent via-indigo-500 to-transparent shadow-[0_0_6px_#4f46e5]" style="animation: scanMove 2s infinite ease-in-out;"></div>
          </div>

          <div class="flex items-center gap-2 text-sm bg-indigo-50 dark:bg-indigo-900/20 px-3.5 py-1.5 rounded-full text-indigo-600 dark:text-indigo-400 border border-indigo-200 dark:border-indigo-800">
            <span>Nominal Bayar:</span>
            <strong>Rp {{ formatPrice(grandTotal) }}</strong>
          </div>

          <p class="text-xs text-slate-500">
            Mendukung: <strong>BCA, Mandiri, BRI, BNI, GoPay, OVO, Dana, ShopeePay, LinkAja</strong>
          </p>
        </div>

        <!-- 3. TRANSFER BANK SECTION -->
        <div v-else-if="paymentMethod === 'transfer'" class="bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg p-4 flex flex-col gap-4">
          <!-- Select Bank -->
          <div class="flex flex-col gap-2">
            <label class="text-sm font-medium text-slate-700 dark:text-slate-300">Pilih Bank Tujuan</label>
            <div class="grid grid-cols-4 gap-2">
              <button 
                v-for="(acc, code) in bankAccounts" 
                :key="code" 
                class="flex flex-col items-center p-2 bg-slate-50 dark:bg-slate-800/50 border rounded-md transition-all text-slate-600 dark:text-slate-400" 
                :class="selectedBank === code ? 'bg-indigo-50 dark:bg-indigo-900/30 border-indigo-500 text-indigo-600 dark:text-indigo-400' : 'border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800'"
                @click="selectedBank = code"
              >
                <span class="font-bold text-sm">{{ code }}</span>
                <span class="text-[10px] text-slate-500">{{ acc.name }}</span>
              </button>
            </div>
          </div>

          <!-- Account Details Box -->
          <div class="p-3.5 flex flex-col gap-1.5 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg">
            <div class="flex justify-between items-center text-xs text-slate-600 dark:text-slate-400">
              <span>Bank:</span>
              <span>{{ bankAccounts[selectedBank].name }}</span>
            </div>
            <div class="flex justify-between items-center text-xs bg-slate-50 dark:bg-slate-900 p-2 rounded-md border border-slate-200 dark:border-slate-700">
              <span>No. Rekening:</span>
              <div class="flex items-center gap-2">
                <strong class="text-sm font-extrabold text-slate-900 dark:text-slate-100 tracking-wide">{{ bankAccounts[selectedBank].account }}</strong>
                <button class="px-2.5 py-1 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded text-[10px] font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700 cursor-pointer transition-colors" @click="copyAccount(bankAccounts[selectedBank].account)">
                  <span v-if="copiedState" class="flex items-center gap-1 text-emerald-500"><CheckIcon class="w-3.5 h-3.5" /> Tersalin</span>
                  <span v-else class="flex items-center gap-1"><ClipboardDocumentIcon class="w-3.5 h-3.5" /> Salin</span>
                </button>
              </div>
            </div>
            <div class="flex justify-between items-center text-xs text-slate-600 dark:text-slate-400">
              <span>Atas Nama:</span>
              <span>{{ bankAccounts[selectedBank].holder }}</span>
            </div>
            <div class="flex justify-between items-center text-xs text-slate-600 dark:text-slate-400">
              <span>Nominal Transfer:</span>
              <strong class="text-sm text-blue-600 dark:text-blue-400">Rp {{ formatPrice(grandTotal) }}</strong>
            </div>
          </div>

          <div class="text-xs text-amber-600 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 p-2.5 rounded-md flex items-start">
            <InformationCircleIcon class="w-4 h-4 mr-1 shrink-0" />
            <span>Mohon verifikasi bahwa dana sudah masuk di m-Banking sebelum menekan tombol Selesaikan.</span>
          </div>
        </div>

        <!-- 4. DEBIT / EDC CARD SECTION -->
        <div v-else-if="paymentMethod === 'debit'" class="bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg p-4 flex flex-col gap-3">
          <div class="flex flex-col gap-3 p-4 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg">
            <div class="flex justify-between items-center">
              <span class="text-sm font-bold text-slate-900 dark:text-slate-100">MESIN EDC KARTU</span>
              <div class="flex gap-1.5">
                <span class="text-[10px] font-extrabold px-1.5 py-0.5 rounded bg-blue-800 text-white">VISA</span>
                <span class="text-[10px] font-extrabold px-1.5 py-0.5 rounded bg-red-600 text-white">MC</span>
                <span class="text-[10px] font-extrabold px-1.5 py-0.5 rounded bg-sky-700 text-white">GPN</span>
              </div>
            </div>

            <div class="bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg p-4 text-center">
              <span class="text-xs text-slate-500 font-semibold">TOTAL CHARGE</span>
              <h3 class="text-2xl font-extrabold text-blue-600 dark:text-blue-400">Rp {{ formatPrice(grandTotal) }}</h3>
              <span class="text-xs font-bold text-blue-500 dark:text-blue-400">INSERT / SWIPE KARTU DEBIT</span>
            </div>

            <div class="flex flex-col gap-2 mt-2">
              <label class="text-sm font-medium text-slate-700 dark:text-slate-300">No. Approval / Ref EDC (Opsional)</label>
              <AppInput 
                type="text" 
                v-model="approvalCode" 
                placeholder="cth: REF-981240" 
              />
            </div>
          </div>
        </div>
      </div>

      <div class="p-4 border-t border-slate-200 dark:border-slate-700 flex justify-end gap-3 bg-white dark:bg-slate-800">
        <AppButton variant="secondary" @click="$emit('close')">Batal</AppButton>
        <AppButton 
          variant="primary" 
          class="flex-1"
          :disabled="isSubmitting || (paymentMethod === 'cash' && changeAmount < 0)"
          @click="submitPayment"
        >
          <span v-if="isSubmitting">Memproses Transaksi...</span>
          <span v-else class="flex items-center justify-center gap-1.5">
            <PrinterIcon class="w-4 h-4" />
            <span>Selesaikan & Cetak Struk</span>
          </span>
        </AppButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import api from '@/utils/api';
import { storeToRefs } from 'pinia';
import { useSettingsStore } from '../stores/settings';
import type { Customer } from '../types';
import { 
  BanknotesIcon, 
  QrCodeIcon, 
  BuildingLibraryIcon, 
  CreditCardIcon, 
  PrinterIcon, 
  XMarkIcon, 
  BoltIcon,
  CheckIcon,
  ClipboardDocumentIcon,
  InformationCircleIcon,
  CheckBadgeIcon
} from '@heroicons/vue/24/outline';
import AppButton from './ui/AppButton.vue';
import AppInput from './ui/AppInput.vue';
import AppBadge from './ui/AppBadge.vue';

const props = defineProps({
  grandTotal: { type: Number, required: true },
  isSubmitting: { type: Boolean, default: false }
});

const emit = defineEmits(['close', 'submit-order']);

const settingsStore = useSettingsStore();
const { settings: storeSetting } = storeToRefs(settingsStore);

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

const customerName = ref('Umum');
const selectedCustomerOption = ref('Umum');
const customersList = ref<Customer[]>([]);

const fetchCustomers = async () => {
  try {
    const res = await api.get('/customers');
    customersList.value = res.data;
  } catch (err: any) {
    console.error('Fetch customers error:', err.response?.data?.error || err.message || 'Error occurred');
  }
};

const fetchedQrisUrl = ref('');

const fetchLatestSettings = async () => {
  try {
    const res = await api.get('/settings');
    const data = res.data;
    if (data.qris_image_url) {
      fetchedQrisUrl.value = data.qris_image_url;
    }
  } catch (err: any) {
    console.error('Fetch settings error in PaymentModal:', err.response?.data?.error || err.message || 'Error occurred');
  }
};

onMounted(() => {
  fetchCustomers();
  fetchLatestSettings();
});

const effectiveQrisUrl = computed(() => {
  return storeSetting.value?.qris_image_url || fetchedQrisUrl.value || '';
});

// Member Discount Calculation
const memberDiscountPercent = computed(() => storeSetting.value?.member_discount_percentage ?? 5);

const matchedMember = computed(() => {
  if (!customerName.value || customerName.value.trim() === '' || customerName.value === 'Umum') return null;
  const q = customerName.value.trim().toLowerCase();
  return customersList.value.find(c => c.name.trim().toLowerCase() === q);
});

const isMatchedMember = computed(() => !!matchedMember.value);

const memberDiscountAmount = computed(() => {
  if (!isMatchedMember.value) return 0;
  return Math.round((props.grandTotal * memberDiscountPercent.value) / 100);
});

const finalGrandTotal = computed(() => {
  return Math.max(0, props.grandTotal - memberDiscountAmount.value);
});

const onCustomerSelect = () => {
  if (selectedCustomerOption.value !== 'custom') {
    customerName.value = selectedCustomerOption.value;
  } else {
    customerName.value = '';
  }
};

const paymentMethod = ref('cash');
const paidAmount = ref(props.grandTotal);
const approvalCode = ref('');

const paymentMethods = [
  { id: 'cash', name: 'Tunai', iconComp: BanknotesIcon, desc: 'Uang Tunai' },
  { id: 'qris', name: 'QRIS', iconComp: QrCodeIcon, desc: 'Scan QR All Pay' },
  { id: 'transfer', name: 'Transfer Bank', iconComp: BuildingLibraryIcon, desc: 'BCA/Mandiri/BRI' },
  { id: 'debit', name: 'Kartu Debit', iconComp: CreditCardIcon, desc: 'Mesin EDC' }
];

const quickCashAmounts = [10000, 20000, 50000, 100000, 200000];

// Bank Transfer Accounts
const selectedBank = ref<string>('BCA');
const bankAccounts: Record<string, { name: string; account: string; holder: string }> = {
  BCA: { name: 'Bank BCA', account: '8830192831', holder: 'KASIR POS STORE' },
  Mandiri: { name: 'Bank Mandiri', account: '137001928302', holder: 'KASIR POS STORE' },
  BRI: { name: 'Bank BRI', account: '0123010293014', holder: 'KASIR POS STORE' },
  BNI: { name: 'Bank BNI', account: '0918239102', holder: 'KASIR POS STORE' }
};

const copiedState = ref(false);
const copyAccount = (accountNumber: string): void => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(accountNumber);
    copiedState.value = true;
    setTimeout(() => { copiedState.value = false; }, 2000);
  }
};

const selectMethod = (method: string): void => {
  paymentMethod.value = method;
  if (method !== 'cash') {
    paidAmount.value = finalGrandTotal.value;
  }
};

const setExactCash = () => {
  paidAmount.value = finalGrandTotal.value;
};

const setPaidAmount = (amount: number): void => {
  paidAmount.value = amount;
};

const changeAmount = computed(() => {
  return (paidAmount.value || 0) - finalGrandTotal.value;
});

const submitPayment = () => {
  let methodLabel = paymentMethod.value;
  if (paymentMethod.value === 'transfer') {
    methodLabel = `transfer (${selectedBank.value})`;
  } else if (paymentMethod.value === 'debit') {
    methodLabel = `debit (${approvalCode.value ? approvalCode.value : 'EDC'})`;
  }

  emit('submit-order', {
    customer_name: customerName.value,
    payment_method: methodLabel,
    paid_amount: paymentMethod.value === 'cash' ? paidAmount.value : finalGrandTotal.value
  });
};
</script>
