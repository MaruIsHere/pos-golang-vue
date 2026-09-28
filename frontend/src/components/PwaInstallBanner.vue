<template>
  <div v-if="showBanner" class="pwa-banner-wrapper">
    <div class="fixed bottom-4 left-4 right-4 md:left-auto md:right-4 md:w-96 p-5 flex flex-col gap-4 bg-white/95 dark:bg-slate-900/95 backdrop-blur-md border border-slate-200/60 dark:border-slate-700/60 rounded-[24px] shadow-xl z-50">
      <div class="pwa-info">
        <div class="pwa-icon-badge">
          <DevicePhoneMobileIcon class="w-5 h-5 text-indigo-600" />
        </div>
        <div class="pwa-text">
          <h4>Install Aplikasi POS Kasir</h4>
          <p v-if="isIos">
            Gunakan di iPhone/iPad: Ketuk tombol <strong>Share</strong> lalu pilih <strong>"Tambah ke Layar Utama"</strong>.
          </p>
          <p v-else>
            Dapatkan pengalaman terbaik sebagai aplikasi native di Android, Windows, Mac, atau Tablet Anda.
          </p>
        </div>
      </div>

      <div class="pwa-actions">
        <button v-if="!isIos && deferredPrompt" class="px-4 py-2 bg-blue-600 text-white rounded-lg font-medium flex items-center gap-1" @click="installPwa">
          <ArrowDownTrayIcon class="w-3.5 h-3.5" />
          <span>Install</span>
        </button>
        <button class="flex-1 py-2 text-sm font-bold text-slate-500 bg-slate-100 dark:bg-slate-800 dark:text-slate-400 rounded-xl hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors" title="Nanti saja" @click="dismissBanner">
          <XMarkIcon class="w-4 h-4 text-slate-400 hover:text-red-500" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { DevicePhoneMobileIcon, ArrowDownTrayIcon, XMarkIcon } from '@heroicons/vue/24/outline';

// Minimal typing for the PWA install prompt (not in TS DOM lib).
interface BeforeInstallPromptEvent extends Event {
  prompt: () => Promise<void>;
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>;
}

interface IosNavigator extends Navigator {
  standalone?: boolean;
}

const showBanner = ref(false);
const deferredPrompt = ref<BeforeInstallPromptEvent | null>(null);
const isIos = ref(false);

const checkIfIos = (): boolean => {
  const userAgent = window.navigator.userAgent.toLowerCase();
  return /iphone|ipad|ipod/.test(userAgent);
};

const checkIfStandalone = (): boolean => {
  return window.matchMedia('(display-mode: standalone)').matches || (window.navigator as IosNavigator).standalone === true;
};

onMounted(() => {
  // If already running in standalone PWA app mode, don't show prompt banner
  if (checkIfStandalone()) {
    showBanner.value = false;
    return;
  }

  isIos.value = checkIfIos();

  // Handle Chrome / Edge / Android install prompt
  window.addEventListener('beforeinstallprompt', (e) => {
    e.preventDefault();
    deferredPrompt.value = e as BeforeInstallPromptEvent;
    showBanner.value = true;
  });

  // Show banner on iOS if not standalone
  if (isIos.value && !checkIfStandalone()) {
    showBanner.value = true;
  }
});

const installPwa = async (): Promise<void> => {
  if (!deferredPrompt.value) return;

  deferredPrompt.value.prompt();
  const choiceResult = await deferredPrompt.value.userChoice;

  if (choiceResult.outcome === 'accepted') {
    console.log('User accepted the PWA install prompt');
    showBanner.value = false;
  }
  deferredPrompt.value = null;
};

const dismissBanner = (): void => {
  showBanner.value = false;
};
</script>


