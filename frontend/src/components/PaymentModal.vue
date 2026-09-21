<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content glass-panel">
      <div class="modal-header">
        <div class="header-info">
          <h3>Pembayaran Kasir</h3>
          <p class="sub-info">Selesaikan transaksi belanja</p>
        </div>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="modal-body">
        <!-- Total Display -->
        <div class="total-card">
          <span class="total-label">Total Yang Harus Dibayar</span>
          <h2 class="total-amount">Rp {{ formatPrice(grandTotal) }}</h2>
        </div>

        <!-- Customer Name -->
        <div class="form-group">
          <label class="form-label">Nama Pelanggan (Opsional)</label>
          <input 
            type="text" 
            class="form-control" 
            v-model="customerName" 
            placeholder="Umum / Pelanggan" 
          />
        </div>

        <!-- Payment Method Tabs -->
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
              <span class="method-name">{{ m.name }}</span>
            </button>
          </div>
        </div>

        <!-- Cash Details -->
        <div v-if="paymentMethod === 'cash'" class="cash-section">
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
              Uang Pas (Rp {{ formatPrice(grandTotal) }})
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

        <!-- Non-Cash Simulation -->
        <div v-else class="non-cash-box">
          <div class="qris-simulation">
            <div class="qr-placeholder">
              <span>📱 QRIS / EDC</span>
            </div>
            <p>Silakan scan QR code atau gesek kartu pada mesin EDC.</p>
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
          <span v-if="isSubmitting">Memproses...</span>
          <span v-else>Selesaikan & Cetak Struk</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';

const props = defineProps({
  grandTotal: { type: Number, required: true },
  isSubmitting: { type: Boolean, default: false }
});

const emit = defineEmits(['close', 'submit-order']);

const formatPrice = (val) => new Intl.NumberFormat('id-ID').format(val || 0);

const customerName = ref('Umum');
const paymentMethod = ref('cash');
const paidAmount = ref(props.grandTotal);

const paymentMethods = [
  { id: 'cash', name: 'Tunai', icon: '💵' },
  { id: 'qris', name: 'QRIS', icon: '📲' },
  { id: 'debit', name: 'Kartu Debit/EDC', icon: '💳' }
];

const quickCashAmounts = [10000, 20000, 50000, 100000, 200000];

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
  emit('submit-order', {
    customer_name: customerName.value,
    payment_method: paymentMethod.value,
    paid_amount: paidAmount.value
  });
};
</script>

<style scoped>
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
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.total-card {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(139, 92, 246, 0.15));
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: var(--radius-lg);
  padding: 1.25rem;
  text-align: center;
}

.total-label {
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.total-amount {
  font-size: 2rem;
  font-weight: 800;
  color: var(--accent-secondary);
  margin-top: 0.25rem;
}

.method-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.75rem;
}

.method-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.35rem;
  padding: 0.85rem 0.5rem;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.method-card.active {
  background: rgba(99, 102, 241, 0.2);
  border-color: var(--accent-primary);
  color: #ffffff;
  box-shadow: 0 0 12px rgba(99, 102, 241, 0.3);
}

.method-icon {
  font-size: 1.4rem;
}

.method-name {
  font-size: 0.75rem;
  font-weight: 700;
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
  font-size: 1.2rem;
  font-weight: 800;
}

.quick-cash-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.5rem;
  margin-top: 0.75rem;
}

.quick-cash-btn {
  padding: 0.5rem;
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 0.8rem;
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
  padding: 1rem;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: var(--radius-md);
  margin-top: 0.75rem;
  font-weight: 700;
}

.change-box.insufficient {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.3);
  color: #f87171;
}

.change-val {
  font-size: 1.25rem;
  font-weight: 800;
}

.non-cash-box {
  padding: 1.5rem;
  text-align: center;
  background: rgba(15, 23, 42, 0.6);
  border-radius: var(--radius-md);
  border: 1px dashed var(--border-color);
}

.qr-placeholder {
  font-size: 1.5rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
  color: var(--accent-primary);
}

.modal-footer {
  padding: 1.25rem 1.5rem;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.btn-submit {
  flex: 1;
}
</style>
