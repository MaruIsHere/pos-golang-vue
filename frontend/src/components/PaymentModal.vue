<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content glass-panel payment-modal-content">
      <div class="modal-header">
        <div class="header-info">
          <h3>Pembayaran Kasir</h3>
          <p class="sub-info">Pilih metode pembayaran dan selesaikan transaksi</p>
        </div>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="modal-body">
        <!-- Total Display -->
        <div class="total-card">
          <span class="total-label">Total Yang Harus Dibayar</span>
          <h2 class="total-amount">Rp {{ formatPrice(grandTotal) }}</h2>
        </div>

        <!-- Customer Name / Member Dropdown -->
        <div class="form-group">
          <label class="form-label">Nama Pelanggan / Member</label>
          <div class="customer-input-flex">
            <select class="form-control" v-model="selectedCustomerOption" @change="onCustomerSelect">
              <option value="Umum">Umum (Non-Member)</option>
              <option v-for="c in customersList" :key="c.id" :value="c.name">
                👤 {{ c.name }} {{ c.phone ? '(' + c.phone + ')' : '' }}
              </option>
              <option value="custom">-- Ketik Nama Manual --</option>
            </select>

            <input 
              v-if="selectedCustomerOption === 'custom'" 
              type="text" 
              class="form-control" 
              v-model="customerName" 
              placeholder="Ketik Nama Pelanggan..." 
            />
          </div>
        </div>

        <!-- Payment Method Grid -->
        <div class="form-group">
          <label class="form-label">Metode Pembayaran</label>
          <div class="method-grid">
            <button 
              v-for="m in paymentMethods" 
              :key="m.id" 
              class="method-card" 
              :class="{ active: paymentMethod === m.id }"
              @click="selectMethod(m.id)"
            >
              <span class="method-icon">{{ m.icon }}</span>
              <div class="method-details">
                <span class="method-name">{{ m.name }}</span>
                <span class="method-desc">{{ m.desc }}</span>
              </div>
            </button>
          </div>
        </div>

        <!-- 1. CASH SECTION -->
        <div v-if="paymentMethod === 'cash'" class="payment-section cash-section">
          <div class="form-group">
            <label class="form-label">Nominal Uang Diterima</label>
            <div class="input-prefix-wrapper">
              <span class="input-prefix">Rp</span>
              <input 
                type="number" 
                class="form-control paid-input" 
                v-model.number="paidAmount" 
                placeholder="0" 
              />
            </div>
          </div>

          <!-- Quick Cash Buttons -->
          <div class="quick-cash-grid">
            <button class="quick-cash-btn exact-btn" @click="setExactCash">
              ⚡ Uang Pas (Rp {{ formatPrice(grandTotal) }})
            </button>
            <button 
              v-for="amount in quickCashAmounts" 
              :key="amount" 
              class="quick-cash-btn"
              @click="setPaidAmount(amount)"
            >
              Rp {{ formatPrice(amount) }}
            </button>
          </div>

          <!-- Change Amount Display -->
          <div class="change-box" :class="{ 'insufficient': changeAmount < 0 }">
            <span>{{ changeAmount >= 0 ? 'Kembalian' : 'Uang Kurang' }}</span>
            <span class="change-val">
              Rp {{ formatPrice(Math.abs(changeAmount)) }}
            </span>
          </div>
        </div>

        <!-- 2. QRIS SECTION -->
        <div v-else-if="paymentMethod === 'qris'" class="payment-section qris-section">
          <div class="qris-card glass-panel">
            <div class="qris-header">
              <div class="qris-logo-badge">
                <span class="qris-text">QRIS</span>
                <span class="qris-sub">National QR Standard</span>
              </div>
              <span class="live-pulse">🟢 Standby Scan</span>
            </div>

            <!-- Custom QRIS Image or SVG QR Graphic -->
            <div class="qr-code-wrapper">
              <img 
                v-if="effectiveQrisUrl" 
                :src="effectiveQrisUrl" 
                alt="Foto QRIS Toko" 
                class="custom-qris-img" 
              />
              <svg v-else class="qr-svg" viewBox="0 0 200 200" width="160" height="160">
                <!-- Background -->
                <rect width="200" height="200" fill="#ffffff" rx="12"/>
                <!-- QR Position Detection Patterns -->
                <!-- Top Left -->
                <rect x="15" y="15" width="45" height="45" fill="#0f172a" rx="4"/>
                <rect x="23" y="23" width="29" height="29" fill="#ffffff" rx="2"/>
                <rect x="30" y="30" width="15" height="15" fill="#6366f1" rx="2"/>
                <!-- Top Right -->
                <rect x="140" y="15" width="45" height="45" fill="#0f172a" rx="4"/>
                <rect x="148" y="23" width="29" height="29" fill="#ffffff" rx="2"/>
                <rect x="155" y="30" width="15" height="15" fill="#6366f1" rx="2"/>
                <!-- Bottom Left -->
                <rect x="15" y="140" width="45" height="45" fill="#0f172a" rx="4"/>
                <rect x="23" y="148" width="29" height="29" fill="#ffffff" rx="2"/>
                <rect x="30" y="155" width="15" height="15" fill="#6366f1" rx="2"/>
                <!-- QR Data Matrix Dots Simulation -->
                <rect x="70" y="20" width="12" height="12" fill="#0f172a"/>
                <rect x="90" y="20" width="12" height="12" fill="#6366f1"/>
                <rect x="110" y="20" width="12" height="12" fill="#0f172a"/>
                <rect x="70" y="40" width="12" height="12" fill="#0f172a"/>
                <rect x="110" y="40" width="12" height="12" fill="#6366f1"/>
                <rect x="20" y="70" width="12" height="12" fill="#6366f1"/>
                <rect x="40" y="70" width="12" height="12" fill="#0f172a"/>
                <rect x="70" y="70" width="12" height="12" fill="#0f172a"/>
                <rect x="90" y="70" width="20" height="20" fill="#6366f1" rx="4"/>
                <rect x="120" y="70" width="12" height="12" fill="#0f172a"/>
                <rect x="140" y="70" width="12" height="12" fill="#6366f1"/>
                <rect x="160" y="70" width="12" height="12" fill="#0f172a"/>
                <rect x="20" y="90" width="12" height="12" fill="#0f172a"/>
                <rect x="40" y="90" width="12" height="12" fill="#6366f1"/>
                <rect x="140" y="90" width="12" height="12" fill="#0f172a"/>
                <rect x="170" y="90" width="12" height="12" fill="#6366f1"/>
                <rect x="20" y="110" width="12" height="12" fill="#6366f1"/>
                <rect x="70" y="110" width="12" height="12" fill="#0f172a"/>
                <rect x="90" y="110" width="12" height="12" fill="#6366f1"/>
                <rect x="110" y="110" width="12" height="12" fill="#0f172a"/>
                <rect x="150" y="110" width="12" height="12" fill="#6366f1"/>
                <rect x="70" y="140" width="12" height="12" fill="#0f172a"/>
                <rect x="90" y="140" width="12" height="12" fill="#6366f1"/>
                <rect x="110" y="140" width="12" height="12" fill="#0f172a"/>
                <rect x="140" y="140" width="12" height="12" fill="#6366f1"/>
                <rect x="160" y="140" width="12" height="12" fill="#0f172a"/>
                <rect x="70" y="165" width="12" height="12" fill="#6366f1"/>
                <rect x="90" y="165" width="12" height="12" fill="#0f172a"/>
                <rect x="120" y="165" width="12" height="12" fill="#6366f1"/>
                <rect x="150" y="165" width="20" height="20" fill="#0f172a" rx="4"/>
              </svg>

              <div class="qr-scan-line"></div>
            </div>

            <div class="qris-amount-tag">
              <span>Nominal Bayar:</span>
              <strong>Rp {{ formatPrice(grandTotal) }}</strong>
            </div>

            <p class="qris-hint">
              Mendukung: <strong>BCA, Mandiri, BRI, BNI, GoPay, OVO, Dana, ShopeePay, LinkAja</strong>
            </p>
          </div>
        </div>

        <!-- 3. TRANSFER BANK SECTION -->
        <div v-else-if="paymentMethod === 'transfer'" class="payment-section transfer-section">
          <!-- Select Bank -->
          <div class="form-group">
            <label class="form-label">Pilih Bank Tujuan</label>
            <div class="bank-selector-grid">
              <button 
                v-for="(acc, code) in bankAccounts" 
                :key="code" 
                class="bank-btn" 
                :class="{ active: selectedBank === code }"
                @click="selectedBank = code"
              >
                <span class="bank-code">{{ code }}</span>
                <span class="bank-title">{{ acc.name }}</span>
              </button>
            </div>
          </div>

          <!-- Account Details Box -->
          <div class="account-card glass-panel">
            <div class="acc-row">
              <span class="acc-label">Bank:</span>
              <span class="acc-val">{{ bankAccounts[selectedBank].name }}</span>
            </div>
            <div class="acc-row highlight">
              <span class="acc-label">No. Rekening:</span>
              <div class="copy-wrapper">
                <strong class="acc-number">{{ bankAccounts[selectedBank].account }}</strong>
                <button class="btn-copy" @click="copyAccount(bankAccounts[selectedBank].account)">
                  {{ copiedState ? '✓ Tersalin' : '📋 Salin' }}
                </button>
              </div>
            </div>
            <div class="acc-row">
              <span class="acc-label">Atas Nama:</span>
              <span class="acc-val">{{ bankAccounts[selectedBank].holder }}</span>
            </div>
            <div class="acc-row">
              <span class="acc-label">Nominal Transfer:</span>
              <strong class="acc-amount">Rp {{ formatPrice(grandTotal) }}</strong>
            </div>
          </div>

          <div class="transfer-hint-box">
            <span>ℹ️ Mohon verifikasi bahwa dana sudah masuk di m-Banking sebelum menekan tombol Selesaikan.</span>
          </div>
        </div>

        <!-- 4. DEBIT / EDC CARD SECTION -->
        <div v-else-if="paymentMethod === 'debit'" class="payment-section debit-section">
          <div class="edc-card glass-panel">
            <div class="edc-top">
              <span class="edc-title">💳 MESIN EDC KARTU</span>
              <div class="card-brands">
                <span class="brand-pill visa">VISA</span>
                <span class="brand-pill master">MC</span>
                <span class="brand-pill gpn">GPN</span>
              </div>
            </div>

            <div class="edc-display">
              <span class="edc-amount-label">TOTAL CHARGE</span>
              <h3 class="edc-amount-val">Rp {{ formatPrice(grandTotal) }}</h3>
              <span class="edc-status">INSERT / SWIPE KARTU DEBIT</span>
            </div>

            <div class="form-group" style="margin-top: 1rem;">
              <label class="form-label">No. Approval / Ref EDC (Opsional)</label>
              <input 
                type="text" 
                class="form-control" 
                v-model="approvalCode" 
                placeholder="cth: REF-981240" 
              />
            </div>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <button class="btn btn-secondary" @click="$emit('close')">Batal</button>
        <button 
          class="btn btn-success btn-submit" 
          :disabled="isSubmitting || (paymentMethod === 'cash' && changeAmount < 0)"
          @click="submitPayment"
        >
          <span v-if="isSubmitting">Memproses Transaksi...</span>
          <span v-else>Selesaikan & Cetak Struk 🖨️</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';

const props = defineProps({
  grandTotal: { type: Number, required: true },
  isSubmitting: { type: Boolean, default: false },
  storeSetting: { type: Object, default: () => ({}) }
});

const emit = defineEmits(['close', 'submit-order']);

const formatPrice = (val) => new Intl.NumberFormat('id-ID').format(val || 0);

const customerName = ref('Umum');
const selectedCustomerOption = ref('Umum');
const customersList = ref([]);

const fetchCustomers = async () => {
  try {
    const res = await fetch('/api/customers');
    if (res.ok) customersList.value = await res.json();
  } catch (err) {
    console.error('Fetch customers error:', err);
  }
};

const fetchedQrisUrl = ref('');

const fetchLatestSettings = async () => {
  try {
    const res = await fetch('/api/settings');
    if (res.ok) {
      const data = await res.json();
      if (data.qris_image_url) {
        fetchedQrisUrl.value = data.qris_image_url;
      }
    }
  } catch (err) {
    console.error('Fetch settings error in PaymentModal:', err);
  }
};

onMounted(() => {
  fetchCustomers();
  fetchLatestSettings();
});

const effectiveQrisUrl = computed(() => {
  return props.storeSetting?.qris_image_url || fetchedQrisUrl.value || '';
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
  { id: 'cash', name: 'Tunai', icon: '💵', desc: 'Uang Tunai' },
  { id: 'qris', name: 'QRIS', icon: '📲', desc: 'Scan QR All Pay' },
  { id: 'transfer', name: 'Transfer Bank', icon: '🏦', desc: 'BCA/Mandiri/BRI/BNI' },
  { id: 'debit', name: 'Kartu Debit', icon: '💳', desc: 'Mesin EDC' }
];

const quickCashAmounts = [10000, 20000, 50000, 100000, 200000];

// Bank Transfer Accounts
const selectedBank = ref('BCA');
const bankAccounts = {
  BCA: { name: 'Bank BCA', account: '8830192831', holder: 'KASIR POS STORE' },
  Mandiri: { name: 'Bank Mandiri', account: '137001928302', holder: 'KASIR POS STORE' },
  BRI: { name: 'Bank BRI', account: '0123010293014', holder: 'KASIR POS STORE' },
  BNI: { name: 'Bank BNI', account: '0918239102', holder: 'KASIR POS STORE' }
};

const copiedState = ref(false);
const copyAccount = (accountNumber) => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(accountNumber);
    copiedState.value = true;
    setTimeout(() => { copiedState.value = false; }, 2000);
  }
};

const selectMethod = (method) => {
  paymentMethod.value = method;
  if (method !== 'cash') {
    paidAmount.value = props.grandTotal;
  }
};

const setExactCash = () => {
  paidAmount.value = props.grandTotal;
};

const setPaidAmount = (amount) => {
  paidAmount.value = amount;
};

const changeAmount = computed(() => {
  return (paidAmount.value || 0) - props.grandTotal;
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
    paid_amount: paymentMethod.value === 'cash' ? paidAmount.value : props.grandTotal
  });
};
</script>

<style scoped>
.payment-modal-content {
  max-width: 540px;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.header-info h3 {
  font-size: 1.2rem;
  font-weight: 800;
}

.sub-info {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.btn-close {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 1.2rem;
  cursor: pointer;
}

.modal-body {
  padding: 1.25rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.total-card {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(139, 92, 246, 0.15));
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: var(--radius-lg);
  padding: 1rem;
  text-align: center;
}

.total-label {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.total-amount {
  font-size: 1.8rem;
  font-weight: 800;
  color: var(--accent-secondary);
  margin-top: 0.15rem;
}

.method-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.6rem;
}

@media (min-width: 480px) {
  .method-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

.method-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  padding: 0.75rem 0.4rem;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
  text-align: center;
}

.method-card.active {
  background: rgba(99, 102, 241, 0.25);
  border-color: var(--accent-primary);
  color: #ffffff;
  box-shadow: 0 0 12px rgba(99, 102, 241, 0.35);
}

.method-icon {
  font-size: 1.3rem;
}

.method-details {
  display: flex;
  flex-direction: column;
}

.method-name {
  font-size: 0.78rem;
  font-weight: 800;
}

.method-desc {
  font-size: 0.65rem;
  color: var(--text-muted);
}

.payment-section {
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 1rem;
}

.input-prefix-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-prefix {
  position: absolute;
  left: 1rem;
  font-weight: 800;
  color: var(--accent-secondary);
}

.paid-input {
  padding-left: 2.75rem;
  font-size: 1.15rem;
  font-weight: 800;
}

.quick-cash-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.4rem;
  margin-top: 0.6rem;
}

.quick-cash-btn {
  padding: 0.45rem;
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  transition: background 0.2s;
}

.quick-cash-btn:hover {
  background: var(--bg-card-hover);
}

.exact-btn {
  grid-column: span 2;
  background: rgba(16, 185, 129, 0.15);
  border-color: rgba(16, 185, 129, 0.4);
  color: #34d399;
}

.change-box {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: var(--radius-md);
  margin-top: 0.6rem;
  font-weight: 700;
  font-size: 0.85rem;
}

.change-box.insufficient {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.3);
  color: #f87171;
}

.change-val {
  font-size: 1.15rem;
  font-weight: 800;
}

/* QRIS Section */
.qris-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 1.25rem;
  text-align: center;
  background: rgba(15, 23, 42, 0.8);
  border-radius: var(--radius-md);
}

.qris-header {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.qris-logo-badge {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.qris-text {
  font-size: 1.2rem;
  font-weight: 900;
  color: #ef4444;
  letter-spacing: 0.05em;
}

.qris-sub {
  font-size: 0.65rem;
  color: var(--text-muted);
}

.live-pulse {
  font-size: 0.72rem;
  font-weight: 700;
  color: #34d399;
  background: rgba(16, 185, 129, 0.15);
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.qr-code-wrapper {
  position: relative;
  padding: 0.6rem;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.3);
}

.qr-scan-line {
  position: absolute;
  top: 8px;
  left: 8px;
  right: 8px;
  height: 3px;
  background: linear-gradient(90deg, transparent, #ef4444, transparent);
  box-shadow: 0 0 8px #ef4444;
  animation: scanMove 2s infinite ease-in-out;
}

@keyframes scanMove {
  0% { top: 10px; }
  50% { top: 155px; }
  100% { top: 10px; }
}

.qris-amount-tag {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  font-size: 0.9rem;
  background: rgba(99, 102, 241, 0.15);
  padding: 0.35rem 0.85rem;
  border-radius: 999px;
  color: #a5b4fc;
}

.qris-hint {
  font-size: 0.75rem;
  color: var(--text-muted);
}

/* Transfer Section */
.bank-selector-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.5rem;
}

.bank-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0.5rem;
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.bank-btn.active {
  background: rgba(99, 102, 241, 0.25);
  border-color: var(--accent-primary);
  color: #ffffff;
}

.bank-code {
  font-weight: 800;
  font-size: 0.85rem;
}

.bank-title {
  font-size: 0.65rem;
  color: var(--text-muted);
}

.account-card {
  margin-top: 0.75rem;
  padding: 0.85rem;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  background: rgba(15, 23, 42, 0.8);
  border-radius: var(--radius-md);
}

.acc-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  color: var(--text-secondary);
}

.acc-row.highlight {
  background: rgba(99, 102, 241, 0.1);
  padding: 0.4rem 0.6rem;
  border-radius: 6px;
}

.copy-wrapper {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.acc-number {
  font-size: 0.95rem;
  font-weight: 800;
  color: #ffffff;
  letter-spacing: 0.05em;
}

.btn-copy {
  padding: 0.2rem 0.5rem;
  background: rgba(99, 102, 241, 0.2);
  border: 1px solid rgba(99, 102, 241, 0.4);
  border-radius: 4px;
  color: #a5b4fc;
  font-size: 0.7rem;
  font-weight: 700;
  cursor: pointer;
}

.acc-amount {
  font-size: 0.95rem;
  color: var(--accent-secondary);
}

.transfer-hint-box {
  margin-top: 0.6rem;
  font-size: 0.72rem;
  color: #fbbf24;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.2);
  padding: 0.5rem 0.75rem;
  border-radius: var(--radius-sm);
}

/* Debit Section */
.edc-card {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(15, 23, 42, 0.8);
  border-radius: var(--radius-md);
}

.edc-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.edc-title {
  font-size: 0.85rem;
  font-weight: 800;
  color: var(--text-primary);
}

.card-brands {
  display: flex;
  gap: 0.35rem;
}

.brand-pill {
  font-size: 0.65rem;
  font-weight: 900;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
}

.brand-pill.visa { background: #1a1f71; color: #fff; }
.brand-pill.master { background: #eb001b; color: #fff; }
.brand-pill.gpn { background: #00529c; color: #fff; }

.edc-display {
  background: #020617;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 1rem;
  text-align: center;
}

.edc-amount-label {
  font-size: 0.7rem;
  color: var(--text-muted);
}

.edc-amount-val {
  font-size: 1.5rem;
  font-weight: 800;
  color: #38bdf8;
}

.edc-status {
  font-size: 0.72rem;
  font-weight: 700;
  color: #34d399;
}

.modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.btn-submit {
  flex: 1;
}

.custom-qris-img {
  width: 160px;
  height: 160px;
  object-fit: contain;
  border-radius: 8px;
}
</style>
