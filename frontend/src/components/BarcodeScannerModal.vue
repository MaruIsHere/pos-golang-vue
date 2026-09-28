<template>
  <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
    <div class="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl shadow-xl w-full max-w-xl flex flex-col overflow-hidden animate-in zoom-in-95 duration-200">
      
      <!-- Header -->
      <div class="flex justify-between items-center p-5 border-b border-slate-100 dark:border-slate-700/50">
        <h2 class="text-xl font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
          <CameraIcon class="w-6 h-6 text-indigo-600 dark:text-indigo-400" />
          Scanner Barcode Produk
        </h2>
        <button class="p-1.5 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-full transition-colors text-slate-400" @click="closeModal">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>

      <!-- Tabs -->
      <div class="flex items-center gap-2 p-5 pb-0">
        <AppButton 
          :variant="activeTab === 'camera' ? 'primary' : 'secondary'"
          class="flex-1 flex items-center justify-center gap-2"
          @click="switchTab('camera')"
        >
          <CameraIcon class="w-4 h-4" />
          <span>Kamera Scanner</span>
        </AppButton>
        <AppButton 
          :variant="activeTab === 'manual' ? 'primary' : 'secondary'"
          class="flex-1 flex items-center justify-center gap-2"
          @click="switchTab('manual')"
        >
          <QrCodeIcon class="w-4 h-4" />
          <span>Alat Tembak USB</span>
        </AppButton>
      </div>

      <!-- Camera Scanner View -->
      <div v-show="activeTab === 'camera'" class="flex flex-col gap-4 p-5">
        <div id="barcode-reader-view" class="w-full min-h-[250px] bg-slate-100 dark:bg-slate-900 rounded-xl overflow-hidden border border-slate-200 dark:border-slate-700"></div>
        
        <div v-if="cameraError" class="flex flex-col items-center justify-center gap-3 p-6 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800/50 rounded-xl text-center">
          <ExclamationTriangleIcon class="w-8 h-8 text-amber-500" />
          <p class="text-sm font-semibold text-amber-800 dark:text-amber-400">{{ cameraError }}</p>
          <AppButton variant="warning" size="sm" @click="initCameraScanner">
            Coba Kamera Lagi
          </AppButton>
        </div>

        <div class="text-sm text-slate-500 text-center font-medium bg-slate-50 dark:bg-slate-900/50 py-2 rounded-lg">
          Arahkan kamera ke kode barcode / QR code pada produk
        </div>
      </div>

      <!-- Manual / Barcode Gun Input View -->
      <div v-show="activeTab === 'manual'" class="flex flex-col gap-6 p-5">
        <div class="flex flex-col gap-2">
          <label class="text-sm font-bold text-slate-700 dark:text-slate-300">Ketik atau Scan dengan Alat USB/Bluetooth:</label>
          <div class="flex items-center gap-2">
            <AppInput 
              ref="manualInputRef"
              type="text" 
              class="flex-1"
              v-model="manualCode"
              placeholder="Contoh: 8991001"
              @keyup.enter="handleManualSubmit"
            />
            <AppButton variant="primary" @click="handleManualSubmit">
              Cari & Tambah
            </AppButton>
          </div>
        </div>

        <!-- Quick Test Barcodes Grid -->
        <div class="flex flex-col gap-3">
          <span class="text-xs font-bold text-slate-500 uppercase tracking-wider">Simulasi Produk Tersedia (Klik untuk scan):</span>
          <div class="grid grid-cols-2 gap-2 max-h-[180px] overflow-y-auto pr-1">
            <button 
              v-for="prod in productsWithBarcodes" 
              :key="prod.id"
              class="px-3 py-2 bg-slate-50 hover:bg-indigo-50 dark:bg-slate-900 dark:hover:bg-indigo-900/30 border border-slate-200 dark:border-slate-700 rounded-lg flex flex-col items-start gap-1 transition-colors text-left"
              @click="simulateScanBarcode(prod.barcode!)"
            >
              <span class="text-sm font-bold text-slate-800 dark:text-slate-200 truncate w-full">{{ prod.name }}</span>
              <span class="text-xs font-mono text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-900/40 px-1.5 rounded">{{ prod.barcode }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Last Scanned Feedback Alert -->
      <div v-if="lastScannedMessage" class="mx-5 mb-5 p-3 rounded-xl border flex items-center gap-3 animate-in fade-in slide-in-from-bottom-2" :class="lastScannedSuccess ? 'bg-emerald-50 border-emerald-200 text-emerald-800 dark:bg-emerald-900/20 dark:border-emerald-800 dark:text-emerald-400' : 'bg-red-50 border-red-200 text-red-800 dark:bg-red-900/20 dark:border-red-800 dark:text-red-400'">
        <CheckCircleIcon v-if="lastScannedSuccess" class="w-6 h-6 shrink-0" />
        <ExclamationCircleIcon v-else class="w-6 h-6 shrink-0" />
        <span class="text-sm font-bold">{{ lastScannedMessage }}</span>
      </div>

      <!-- Footer -->
      <div class="flex justify-between items-center p-5 border-t border-slate-100 dark:border-slate-700/50 bg-slate-50 dark:bg-slate-900/30">
        <label class="flex items-center gap-2 cursor-pointer group">
          <input type="checkbox" v-model="keepScanningMode" class="w-4 h-4 text-indigo-600 rounded border-slate-300 focus:ring-indigo-600 dark:border-slate-600 dark:bg-slate-800 dark:checked:bg-indigo-500" />
          <span class="text-xs font-bold text-slate-600 dark:text-slate-400 group-hover:text-slate-900 dark:group-hover:text-slate-200 transition-colors">Mode Beruntun (Tetap buka scanner)</span>
        </label>
        <AppButton variant="secondary" @click="closeModal">Selesai</AppButton>
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


