<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content glass-panel receipt-modal-wrapper">
      <div class="modal-header no-print">
        <h3>Struk Belanja</h3>
        <button class="btn-close" @click="$emit('close')">
          <XMarkIcon class="w-5 h-5 text-slate-500" />
        </button>
      </div>

      <!-- Thermal Receipt Preview -->
      <div class="receipt-paper" id="receipt-print-area">
        <div class="receipt-header">
          <h2 class="receipt-store-name">{{ storeSetting.store_name || 'KASIR POS PRO' }}</h2>
          <p class="receipt-address">{{ storeSetting.address || 'Jl. Utama No. 88, Jakarta' }}</p>
          <p class="receipt-phone">Telp: {{ storeSetting.phone || '0812-3456-7890' }}</p>
        </div>

        <div class="receipt-divider">--------------------------------</div>

        <div class="receipt-info-rows">
          <div class="r-row">
            <span>No. Struk:</span>
            <span class="bold">{{ order.invoice_no }}</span>
          </div>
          <div class="r-row">
            <span>Tanggal:</span>
            <span>{{ formatDate(order.created_at) }}</span>
          </div>
          <div class="r-row">
            <span>Kasir:</span>
            <span>{{ order.cashier_name || 'Kasir 1' }}</span>
          </div>
          <div class="r-row">
            <span>Pelanggan:</span>
            <span>{{ order.customer_name || 'Umum' }}</span>
          </div>
        </div>

        <div class="receipt-divider">--------------------------------</div>

        <!-- Items Table -->
        <div class="receipt-items">
          <div v-for="item in order.order_items" :key="item.id" class="r-item">
            <div class="r-item-title">{{ item.product_name }}</div>
            <div class="r-item-details">
              <span>{{ item.quantity }} x {{ formatPrice(item.product_price) }}</span>
              <span class="r-item-subtotal">Rp {{ formatPrice(item.subtotal) }}</span>
            </div>
          </div>
        </div>

        <div class="receipt-divider">--------------------------------</div>

        <!-- Totals -->
        <div class="receipt-totals">
          <div class="r-total-row">
            <span>Subtotal:</span>
            <span>Rp {{ formatPrice(order.total_amount) }}</span>
          </div>
          <div v-if="order.discount > 0" class="r-total-row">
            <span>Diskon:</span>
            <span>- Rp {{ formatPrice(order.discount) }}</span>
          </div>
          <div class="r-total-row">
            <span>Pajak:</span>
            <span>Rp {{ formatPrice(order.tax) }}</span>
          </div>
          <div class="r-total-row grand">
            <span>TOTAL:</span>
            <span>Rp {{ formatPrice(order.grand_total) }}</span>
          </div>
          <div class="r-total-row">
            <span>Bayar ({{ (order.payment_method || 'cash').toUpperCase() }}):</span>
            <span>Rp {{ formatPrice(order.paid_amount) }}</span>
          </div>
          <div class="r-total-row">
            <span>Kembali:</span>
            <span>Rp {{ formatPrice(order.change_amount) }}</span>
          </div>
        </div>

        <div class="receipt-divider">--------------------------------</div>

        <div class="receipt-footer">
          <p class="footer-msg">{{ storeSetting.receipt_footer || 'Terima kasih telah berbelanja!' }}</p>
          <p class="powered-by">Powered by POS Golang + Vue.js</p>
        </div>
      </div>

      <div class="modal-footer no-print">
        <button class="btn btn-secondary" @click="$emit('close')">Tutup</button>
        <button class="btn btn-primary flex items-center gap-1.5" @click="printReceipt">
          <PrinterIcon class="w-4 h-4" />
          <span>Cetak Struk</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useSettingsStore } from '../stores/settings';
import type { Order } from '../types';
import { PrinterIcon, XMarkIcon } from '@heroicons/vue/24/outline';

defineProps({
  order: { type: Object as () => Order, required: true }
});

defineEmits(['close']);

const settingsStore = useSettingsStore();
const { settings: storeSetting } = storeToRefs(settingsStore);

const formatPrice = (val: number): string => new Intl.NumberFormat('id-ID').format(val || 0);

const formatDate = (dateStr?: string): string => {
  if (!dateStr) return new Date().toLocaleString('id-ID');
  return new Date(dateStr).toLocaleString('id-ID');
};

const printReceipt = () => {
  window.print();
};
</script>

<style scoped>
.receipt-modal-wrapper {
  max-width: 420px;
  background: #ffffff;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid #e2e8f0;
}

.modal-header h3 {
  font-size: 1.1rem;
  font-weight: 700;
  color: #0f172a;
}

.btn-close {
  background: transparent;
  border: none;
  cursor: pointer;
}

/* Thermal Paper Look */
.receipt-paper {
  background: #ffffff;
  color: #111827;
  font-family: 'Courier New', Courier, monospace;
  padding: 1.5rem 1.25rem;
  font-size: 0.82rem;
  line-height: 1.4;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  margin: 1rem 1.25rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
}

.receipt-header {
  text-align: center;
  margin-bottom: 0.5rem;
}

.receipt-store-name {
  font-size: 1.1rem;
  font-weight: 900;
  text-transform: uppercase;
  margin-bottom: 0.2rem;
}

.receipt-address, .receipt-phone {
  font-size: 0.75rem;
  color: #4b5563;
}

.receipt-divider {
  text-align: center;
  color: #9ca3af;
  margin: 0.4rem 0;
  font-size: 0.75rem;
}

.receipt-info-rows .r-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
}

.bold {
  font-weight: 700;
}

.receipt-items {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.r-item-title {
  font-weight: 700;
}

.r-item-details {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: #374151;
}

.r-item-subtotal {
  font-weight: 700;
}

.receipt-totals {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.r-total-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.78rem;
}

.r-total-row.grand {
  font-size: 0.95rem;
  font-weight: 900;
  margin: 0.25rem 0;
  border-top: 1px dashed #111;
  border-bottom: 1px dashed #111;
  padding: 0.2rem 0;
}

.receipt-footer {
  text-align: center;
  margin-top: 0.5rem;
}

.footer-msg {
  font-size: 0.75rem;
  font-weight: 600;
  white-space: pre-line;
}

.powered-by {
  font-size: 0.65rem;
  color: #6b7280;
  margin-top: 0.5rem;
}

.modal-footer {
  padding: 1rem 1.25rem;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

@media print {
  body * {
    visibility: hidden;
  }
  #receipt-print-area, #receipt-print-area * {
    visibility: visible;
  }
  #receipt-print-area {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    margin: 0;
    box-shadow: none;
    border: none;
  }
  .no-print {
    display: none !important;
  }
}
</style>
