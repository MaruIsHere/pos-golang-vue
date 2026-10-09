<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="opacity-0 translate-y-10"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 translate-y-10"
  >
    <div v-if="showBanner" class="fixed bottom-20 left-4 right-4 md:bottom-6 md:left-auto md:right-6 md:w-[420px] z-[200]">
      <div class="flex flex-col gap-4 bg-white/95 dark:bg-slate-900/95 backdrop-blur-xl border border-slate-200/80 dark:border-slate-700/80 rounded-[28px] p-5 shadow-[0_20px_60px_-15px_rgba(0,0,0,0.1)] dark:shadow-[0_20px_60px_-15px_rgba(0,0,0,0.5)]">
        
        <div class="flex items-start gap-4">
          <div class="w-12 h-12 shrink-0 bg-gradient-to-br from-indigo-500 to-indigo-700 rounded-2xl flex items-center justify-center shadow-lg shadow-indigo-500/30">
            <DevicePhoneMobileIcon class="w-6 h-6 text-white" />
          </div>
          
          <div class="flex flex-col pt-0.5">
            <h4 class="text-base font-bold text-slate-900 dark:text-slate-100 tracking-tight leading-tight">
              Instal Aplikasi POS Kasir
            </h4>
            <p v-if="isIos" class="text-xs text-slate-500 dark:text-slate-400 mt-1 leading-relaxed">
              Tekan tombol <strong class="text-slate-700 dark:text-slate-300">Share</strong> di Safari, lalu pilih <strong class="text-slate-700 dark:text-slate-300">"Tambah ke Layar Utama"</strong>.
            </p>
            <p v-else class="text-xs text-slate-500 dark:text-slate-400 mt-1 leading-relaxed">
              Pasang sebagai aplikasi *native* agar lebih cepat, stabil, dan bisa dipakai secara *offline*.
            </p>
          </div>
        </div>

        <div class="flex items-center gap-3 mt-1">
          <button 
            class="flex-1 py-2.5 px-4 text-xs font-bold text-slate-600 dark:text-slate-300 bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 rounded-xl transition-colors active:scale-[0.98]"
            @click="dismissBanner"
          >
            Nanti Saja
          </button>
          
          <button 
            v-if="!isIos && deferredPrompt"
            class="flex-1 py-2.5 px-4 text-xs font-bold text-white bg-indigo-600 hover:bg-indigo-700 rounded-xl flex items-center justify-center gap-1.5 shadow-md shadow-indigo-600/20 transition-colors active:scale-[0.98]"
            @click="installPwa"
          >
            <ArrowDownTrayIcon class="w-4 h-4" />
            <span>Instal Sekarang</span>
          </button>
        </div>
      </div>
    </div>
  </Transition>
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


