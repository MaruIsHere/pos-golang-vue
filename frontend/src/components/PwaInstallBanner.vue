<template>
  <div v-if="showBanner" class="pwa-banner-wrapper">
    <div class="pwa-banner glass-panel">
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
        <button v-if="!isIos && deferredPrompt" class="btn btn-primary btn-pwa-install flex items-center gap-1" @click="installPwa">
          <ArrowDownTrayIcon class="w-3.5 h-3.5" />
          <span>Install</span>
        </button>
        <button class="btn-pwa-dismiss" title="Nanti saja" @click="dismissBanner">
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

<style scoped>
.pwa-banner-wrapper {
  position: fixed;
  top: 74px;
  right: 1.5rem;
  z-index: 999;
  max-width: 440px;
  animation: slideDown 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

@media (max-width: 640px) {
  .pwa-banner-wrapper {
    top: auto;
    bottom: 74px;
    left: 1rem;
    right: 1rem;
    max-width: none;
  }
}

@keyframes slideDown {
  from {
    transform: translateY(-20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.pwa-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.85rem 1.15rem;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-md);
}

.pwa-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.pwa-icon-badge {
  padding: 0.4rem;
  background: rgba(99, 102, 241, 0.12);
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: 8px;
}

.pwa-text h4 {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--text-primary);
}

.pwa-text p {
  font-size: 0.72rem;
  color: var(--text-secondary);
  line-height: 1.3;
  margin-top: 0.15rem;
}

.pwa-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  white-space: nowrap;
}

.btn-pwa-install {
  padding: 0.4rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 700;
}

.btn-pwa-dismiss {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0.2rem;
}
</style>
