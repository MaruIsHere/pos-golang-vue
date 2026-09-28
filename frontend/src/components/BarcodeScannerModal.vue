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
            <button class="px-3 py-1.5 bg-slate-200 dark:bg-slate-700 text-slate-900 dark:text-slate-100 rounded-md text-sm font-medium" @click="initCameraScanner">
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
              <button class="px-4 py-2 bg-blue-600 text-white rounded-lg font-medium" @click="handleManualSubmit">
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
          <CheckCircleIcon v-if="lastScannedSuccess" class="w-5 h-5 text-emerald-500 shrink-0" />
          <ExclamationCircleIcon v-else class="w-5 h-5 text-red-500 shrink-0" />
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
        <button class="px-4 py-2 bg-slate-200 dark:bg-slate-700 text-slate-900 dark:text-slate-100 rounded-lg font-medium" @click="closeModal">Selesai</button>
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


