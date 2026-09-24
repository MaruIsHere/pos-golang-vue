<template>
  <div class="modal-overlay" @click.self="closeModal">
    <div class="modal-content glass-panel scanner-modal">
      <div class="modal-header">
        <div class="header-title-box">
          <CameraIcon class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
          <h3>Scan Barcode Produk</h3>
        </div>
        <button class="btn-close" @click="closeModal">
          <XMarkIcon class="w-5 h-5 text-slate-500" />
        </button>
      </div>

      <div class="scanner-body">
        <!-- Tabs for Scan Method -->
        <div class="scan-tabs">
          <button 
            class="scan-tab-btn" 
            :class="{ active: activeTab === 'camera' }"
            @click="switchTab('camera')"
          >
            <CameraIcon class="w-4 h-4" />
            <span>Kamera HP / WebCam</span>
          </button>
          <button 
            class="scan-tab-btn" 
            :class="{ active: activeTab === 'manual' }"
            @click="switchTab('manual')"
          >
            <QrCodeIcon class="w-4 h-4" />
            <span>Barcode Gun / Manual</span>
          </button>
        </div>

        <!-- Camera Scanner View -->
        <div v-show="activeTab === 'camera'" class="camera-container">
          <div id="barcode-reader-view" class="reader-box"></div>
          
          <div v-if="cameraError" class="camera-error-box">
            <ExclamationTriangleIcon class="w-8 h-8 text-amber-500" />
            <p>{{ cameraError }}</p>
            <button class="btn btn-secondary btn-sm" @click="initCameraScanner">
              Coba Kamera Lagi
            </button>
          </div>

          <div class="scanner-hint">
            Arahkan kamera ke kode barcode / QR code pada produk
          </div>
        </div>

        <!-- Manual / Barcode Gun Input View -->
        <div v-show="activeTab === 'manual'" class="manual-container">
          <div class="form-group">
            <label class="form-label">Ketik atau Scan dengan Barcode Gun (USB/Bluetooth):</label>
            <div class="input-with-button">
              <input 
                ref="manualInputRef"
                type="text" 
                class="form-control barcode-input" 
                v-model="manualCode"
                placeholder="Contoh: 8991001"
                @keyup.enter="handleManualSubmit"
              />
              <button class="btn btn-primary" @click="handleManualSubmit">
                Cari & Tambah
              </button>
            </div>
          </div>

          <!-- Quick Test Barcodes Grid -->
          <div class="quick-test-section">
            <span class="quick-test-label">Simulasi Scan Kode Barcode Produk Tersedia:</span>
            <div class="barcode-pills">
              <button 
                v-for="prod in productsWithBarcodes" 
                :key="prod.id"
                class="barcode-pill-btn"
                @click="simulateScanBarcode(prod.barcode!)"
              >
                <span class="pill-name">{{ prod.name }}</span>
                <span class="pill-code">{{ prod.barcode }}</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Last Scanned Feedback Alert -->
        <div v-if="lastScannedMessage" class="scan-alert" :class="lastScannedSuccess ? 'success' : 'error'">
          <CheckCircleIcon v-if="lastScannedSuccess" class="w-5 h-5 text-emerald-500 flex-shrink-0" />
          <ExclamationCircleIcon v-else class="w-5 h-5 text-red-500 flex-shrink-0" />
          <span>{{ lastScannedMessage }}</span>
        </div>
      </div>

      <div class="modal-footer">
        <div class="continuous-mode">
          <label class="toggle-label flex items-center gap-2 cursor-pointer">
            <input type="checkbox" v-model="keepScanningMode" class="toggle-checkbox" />
            <span class="text-xs font-semibold text-slate-600 dark:text-slate-300">Mode Scan Beruntun (Tetap buka scanner)</span>
          </label>
        </div>
        <button class="btn btn-secondary" @click="closeModal">Selesai</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { Html5Qrcode, Html5QrcodeSupportedFormats } from 'html5-qrcode';
import type { Product } from '../types';
import { 
  CameraIcon, 
  QrCodeIcon, 
  XMarkIcon, 
  ExclamationTriangleIcon, 
  CheckCircleIcon, 
  ExclamationCircleIcon 
} from '@heroicons/vue/24/outline';

const props = defineProps({
  products: { type: Array as () => Product[], default: () => [] }
});

const emit = defineEmits(['close', 'scan-success']);

const activeTab = ref<'camera' | 'manual'>('camera');
const manualCode = ref('');
const manualInputRef = ref<HTMLInputElement | null>(null);
const cameraError = ref('');
const lastScannedMessage = ref('');
const lastScannedSuccess = ref(true);
const keepScanningMode = ref(true);

let html5QrcodeScanner: Html5Qrcode | null = null;
let scanCooldown = false;

const productsWithBarcodes = computed(() => {
  return props.products.filter(p => p.barcode && p.barcode.trim() !== '').slice(0, 8);
});

const playBeepSound = () => {
  try {
    const audioCtx = new (window.AudioContext || (window as any).webkitAudioContext)();
    const osc = audioCtx.createOscillator();
    const gain = audioCtx.createGain();
    osc.type = 'sine';
    osc.frequency.value = 1800;
    gain.gain.setValueAtTime(0.15, audioCtx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.00001, audioCtx.currentTime + 0.12);
    osc.connect(gain);
    gain.connect(audioCtx.destination);
    osc.start();
    osc.stop(audioCtx.currentTime + 0.12);
  } catch (e) {
    // Audio autoplay restrict fallback
  }
};

const processScannedCode = (code: string) => {
  if (scanCooldown) return;
  scanCooldown = true;
  setTimeout(() => { scanCooldown = false; }, 1200);

  const cleanCode = code.trim().toLowerCase();
  const matched = props.products.find(p => p.barcode && p.barcode.trim().toLowerCase() === cleanCode);

  if (matched) {
    playBeepSound();
    lastScannedSuccess.value = true;
    lastScannedMessage.value = `Berhasil! "${matched.name}" ditambahkan ke keranjang.`;
    emit('scan-success', matched);

    if (!keepScanningMode.value) {
      setTimeout(() => {
        closeModal();
      }, 500);
    }
  } else {
    lastScannedSuccess.value = false;
    lastScannedMessage.value = `Kode '${code}' tidak ditemukan pada daftar produk.`;
  }
};

const handleManualSubmit = () => {
  if (!manualCode.value) return;
  processScannedCode(manualCode.value);
  manualCode.value = '';
};

const simulateScanBarcode = (barcode: string) => {
  processScannedCode(barcode);
};

const switchTab = (tab: 'camera' | 'manual') => {
  activeTab.value = tab;
  if (tab === 'camera') {
    nextTick(() => {
      initCameraScanner();
    });
  } else {
    stopCameraScanner();
    nextTick(() => {
      manualInputRef.value?.focus();
    });
  }
};

const initCameraScanner = async () => {
  cameraError.value = '';
  stopCameraScanner();

  try {
    const html5Qr = new Html5Qrcode("barcode-reader-view", {
      formatsToSupport: [
        Html5QrcodeSupportedFormats.EAN_13,
        Html5QrcodeSupportedFormats.EAN_8,
        Html5QrcodeSupportedFormats.CODE_128,
        Html5QrcodeSupportedFormats.CODE_39,
        Html5QrcodeSupportedFormats.UPC_A,
        Html5QrcodeSupportedFormats.UPC_E,
        Html5QrcodeSupportedFormats.QR_CODE
      ],
      verbose: false
    });
    html5QrcodeScanner = html5Qr;

    await html5Qr.start(
      { facingMode: "environment" },
      {
        fps: 10,
        qrbox: { width: 250, height: 160 }
      },
      (decodedText) => {
        processScannedCode(decodedText);
      },
      () => {
        // Scanning in progress...
      }
    );
  } catch (err: any) {
    console.error("Camera init error:", err);
    cameraError.value = "Tidak dapat mengakses kamera perangkat. Pastikan izin kamera telah diberikan atau gunakan mode 'Barcode Gun / Manual'.";
  }
};

const stopCameraScanner = async () => {
  if (html5QrcodeScanner) {
    try {
      if (html5QrcodeScanner.isScanning) {
        await html5QrcodeScanner.stop();
      }
    } catch (e) {
      console.warn("Scanner stop error:", e);
    } finally {
      html5QrcodeScanner = null;
    }
  }
};

const closeModal = async () => {
  await stopCameraScanner();
  emit('close');
};

onMounted(() => {
  initCameraScanner();
});

onUnmounted(() => {
  stopCameraScanner();
});
</script>

<style scoped>
.scanner-modal {
  max-width: 540px;
  width: 95%;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.header-title-box {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.header-title-box h3 {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-primary);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-primary);
}

.btn-close {
  background: transparent;
  border: none;
  cursor: pointer;
}

.scanner-body {
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

/* Tabs */
.scan-tabs {
  display: flex;
  gap: 0.5rem;
  background: var(--bg-primary);
  padding: 0.25rem;
  border-radius: 10px;
  border: 1px solid var(--border-color);
}

.scan-tab-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-family: var(--font-family);
  font-size: 0.825rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
}

.scan-tab-btn.active {
  background: var(--bg-card);
  color: var(--accent-primary);
  box-shadow: var(--shadow-sm);
}

/* Camera Reader Box */
.camera-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
}

.reader-box {
  width: 100%;
  max-width: 440px;
  min-height: 220px;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--border-color);
  background: #000000;
  position: relative;
}

.camera-error-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
  text-align: center;
  gap: 0.5rem;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  color: var(--text-secondary);
  font-size: 0.85rem;
}

.scanner-hint {
  font-size: 0.78rem;
  color: var(--text-muted);
  text-align: center;
}

/* Manual Section */
.manual-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.input-with-button {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.25rem;
}

.barcode-input {
  flex: 1;
  font-size: 0.95rem;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.quick-test-section {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  background: var(--bg-primary);
  padding: 0.85rem;
  border-radius: 10px;
  border: 1px solid var(--border-color);
}

.quick-test-label {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-secondary);
}

.barcode-pills {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
}

.barcode-pill-btn {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 0.35rem 0.65rem;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.barcode-pill-btn:hover {
  border-color: var(--accent-primary);
  background: var(--bg-card-hover);
}

.pill-name {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--text-primary);
}

.pill-code {
  font-size: 0.68rem;
  color: var(--accent-primary);
  font-family: monospace;
}

/* Alert Notification */
.scan-alert {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.65rem 0.85rem;
  border-radius: 8px;
  font-size: 0.825rem;
  font-weight: 600;
  animation: fadeIn 0.15s ease;
}

.scan-alert.success {
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid rgba(16, 185, 129, 0.3);
  color: #34d399;
}

.scan-alert.error {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.85rem 1.25rem;
  border-top: 1px solid var(--border-color);
  background: var(--bg-primary);
}

.toggle-checkbox {
  accent-color: var(--accent-primary);
  width: 16px;
  height: 16px;
}
</style>
