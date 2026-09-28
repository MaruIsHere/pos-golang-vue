<template>
  <div class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-[100] flex items-center justify-center p-4 transition-opacity">
    <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-[28px] shadow-2xl w-full max-w-[500px] flex flex-col overflow-hidden animate-in zoom-in-95 duration-200">
      
      <!-- Header -->
      <div class="flex justify-between items-center px-6 py-5 border-b border-slate-100 dark:border-slate-800/60 bg-white dark:bg-slate-900 z-10">
        <h2 class="text-lg font-black text-slate-800 dark:text-slate-100 flex items-center gap-2">
          <QrCodeIcon class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
          Scanner Produk
        </h2>
        <button class="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-full transition-colors text-slate-400" @click="closeModal">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>

      <!-- Segmented Control (iOS Style Tabs) -->
      <div class="px-6 pt-5 pb-2 bg-slate-50/50 dark:bg-slate-900/30">
        <div class="relative flex items-center bg-slate-200/60 dark:bg-slate-800 p-1.5 rounded-[18px] w-full mx-auto">
          <!-- Sliding Indicator -->
          <div 
            class="absolute inset-y-1.5 left-1.5 w-[calc(50%-6px)] bg-white dark:bg-slate-700 shadow-sm rounded-2xl transition-transform duration-300 ease-out" 
            :class="activeTab === 'manual' ? 'translate-x-full' : 'translate-x-0'"
          ></div>
          
          <!-- Mode Kamera -->
          <button 
            class="relative flex-1 flex justify-center items-center gap-2 py-2.5 text-sm font-bold z-10 transition-colors duration-300 rounded-2xl outline-none" 
            :class="activeTab === 'camera' ? 'text-indigo-600 dark:text-indigo-300' : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'" 
            @click="switchTab('camera')"
          >
            <CameraIcon class="w-4 h-4" /> Kamera
          </button>
          
          <!-- Mode Alat Tembak -->
          <button 
            class="relative flex-1 flex justify-center items-center gap-2 py-2.5 text-sm font-bold z-10 transition-colors duration-300 rounded-2xl outline-none" 
            :class="activeTab === 'manual' ? 'text-indigo-600 dark:text-indigo-300' : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'" 
            @click="switchTab('manual')"
          >
            <QrCodeIcon class="w-4 h-4" /> Alat USB
          </button>
        </div>
      </div>

      <!-- Camera Scanner View -->
      <div v-show="activeTab === 'camera'" class="flex flex-col gap-4 px-6 pb-6 pt-3 bg-slate-50/50 dark:bg-slate-900/30">
        <div id="barcode-reader-view" class="w-full aspect-[4/3] bg-slate-900 rounded-[20px] overflow-hidden border-2 border-slate-200 dark:border-slate-700 shadow-inner relative flex items-center justify-center">
           <!-- Placeholder while camera loads -->
           <div class="absolute flex flex-col items-center justify-center text-slate-500 z-0">
              <CameraIcon class="w-10 h-10 mb-2 opacity-50 animate-pulse" />
           </div>
        </div>
        
        <div v-if="cameraError" class="flex flex-col items-center justify-center gap-3 p-5 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800/50 rounded-[18px] text-center shadow-sm">
          <ExclamationTriangleIcon class="w-7 h-7 text-amber-500" />
          <p class="text-xs font-bold text-amber-800 dark:text-amber-400 leading-relaxed">{{ cameraError }}</p>
          <button class="px-4 py-2 bg-amber-500 hover:bg-amber-600 text-white font-bold rounded-xl text-xs transition-colors shadow-sm" @click="initCameraScanner">
            Coba Kamera Lagi
          </button>
        </div>

        <div class="text-xs text-slate-500 text-center font-semibold bg-white dark:bg-slate-800 py-2.5 rounded-xl border border-slate-100 dark:border-slate-700">
          Arahkan kamera ke area kode barcode produk
        </div>
      </div>

      <!-- Manual / Barcode Gun Input View -->
      <div v-show="activeTab === 'manual'" class="flex flex-col gap-5 px-6 pb-6 pt-3 bg-slate-50/50 dark:bg-slate-900/30">
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">Ketik / Tembak Kode Barcode</label>
          <div class="flex items-center gap-2">
            <input 
              ref="manualInputRef"
              type="text" 
              class="flex-1 px-4 py-3 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-[16px] text-sm font-bold text-slate-800 dark:text-slate-100 placeholder:text-slate-400 focus:outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 shadow-sm"
              v-model="manualCode"
              placeholder="Contoh: 8991001"
              @keyup.enter="handleManualSubmit"
            />
            <button class="px-5 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-[16px] shadow-sm shadow-indigo-600/20 transition-all active:scale-95" @click="handleManualSubmit">
              Cari
            </button>
          </div>
        </div>

        <!-- Quick Test Barcodes Grid -->
        <div class="flex flex-col gap-2.5 mt-2">
          <span class="text-[0.65rem] font-bold text-slate-400 uppercase tracking-widest">Simulasi Scan Produk:</span>
          <div class="grid grid-cols-2 gap-2 max-h-[160px] overflow-y-auto pr-1 custom-scrollbar">
            <button 
              v-for="prod in productsWithBarcodes" 
              :key="prod.id"
              class="p-3 bg-white hover:bg-indigo-50 dark:bg-slate-800 dark:hover:bg-slate-700 border border-slate-200 dark:border-slate-700 rounded-2xl flex flex-col items-start gap-1 transition-colors text-left shadow-sm"
              @click="simulateScanBarcode(prod.barcode!)"
            >
              <span class="text-xs font-bold text-slate-700 dark:text-slate-200 line-clamp-1 w-full">{{ prod.name }}</span>
              <span class="text-[10px] font-mono font-bold text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-900/40 px-1.5 py-0.5 rounded">{{ prod.barcode }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Last Scanned Feedback Alert -->
      <div v-if="lastScannedMessage" class="mx-6 mb-5 p-3.5 rounded-[16px] border flex items-center gap-3 shadow-sm animate-in fade-in slide-in-from-bottom-2" :class="lastScannedSuccess ? 'bg-emerald-50 border-emerald-200 text-emerald-800 dark:bg-emerald-900/20 dark:border-emerald-800 dark:text-emerald-400' : 'bg-red-50 border-red-200 text-red-800 dark:bg-red-900/20 dark:border-red-800 dark:text-red-400'">
        <CheckCircleIcon v-if="lastScannedSuccess" class="w-6 h-6 shrink-0" />
        <ExclamationCircleIcon v-else class="w-6 h-6 shrink-0" />
        <span class="text-sm font-bold leading-tight">{{ lastScannedMessage }}</span>
      </div>

      <!-- Footer -->
      <div class="flex justify-between items-center px-6 py-5 border-t border-slate-100 dark:border-slate-800/60 bg-white dark:bg-slate-900 z-10">
        <label class="flex items-center gap-2.5 cursor-pointer group">
          <input type="checkbox" v-model="keepScanningMode" class="w-4 h-4 text-indigo-600 rounded-[4px] border-slate-300 focus:ring-indigo-600 dark:border-slate-600 dark:bg-slate-800 transition-colors" />
          <span class="text-xs font-bold text-slate-500 dark:text-slate-400 group-hover:text-slate-800 dark:group-hover:text-slate-200 transition-colors">Mode Beruntun</span>
        </label>
        <button class="px-6 py-2.5 bg-slate-900 hover:bg-black dark:bg-white dark:hover:bg-slate-200 text-white dark:text-slate-900 font-bold rounded-xl shadow-md transition-transform active:scale-95 text-sm" @click="closeModal">
          Selesai
        </button>
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


