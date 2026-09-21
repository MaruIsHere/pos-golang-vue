<template>
  <div v-if="showBanner" class="pwa-banner-wrapper">
    <div class="pwa-banner glass-panel">
      <div class="pwa-info">
        <div class="pwa-icon-badge">📲</div>
        <div class="pwa-text">
          <h4>Install Aplikasi POS Kasir</h4>
          <p v-if="isIos">
            Gunakan di iPhone/iPad: Ketuk tombol <strong>Share 📤</strong> lalu pilih <strong>"Tambah ke Layar Utama"</strong>.
          </p>
          <p v-else>
            Dapatkan pengalaman terbaik sebagai aplikasi native di Android, Windows, Mac, atau Tablet Anda.
          </p>
        </div>
      </div>

      <div class="pwa-actions">
        <button v-if="!isIos && deferredPrompt" class="btn btn-primary btn-pwa-install" @click="installPwa">
          ⚡ Install Sekarang
        </button>
        <button class="btn-pwa-dismiss" title="Nanti saja" @click="dismissBanner">
          ✕
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const showBanner = ref(false);
const deferredPrompt = ref(null);
const isIos = ref(false);

const checkIfIos = () => {
  const userAgent = window.navigator.userAgent.toLowerCase();
  return /iphone|ipad|ipod/.test(userAgent);
};

const checkIfStandalone = () => {
  return window.matchMedia('(display-mode: standalone)').matches || window.navigator.standalone === true;
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
    deferredPrompt.value = e;
    showBanner.value = true;
  });

  // Show banner on iOS if not standalone
  if (isIos.value && !checkIfStandalone()) {
    showBanner.value = true;
  }
});

const installPwa = async () => {
  if (!deferredPrompt.value) return;

  deferredPrompt.value.prompt();
  const choiceResult = await deferredPrompt.value.userChoice;

  if (choiceResult.outcome === 'accepted') {
    console.log('User accepted the PWA install prompt');
    showBanner.value = false;
  }
  deferredPrompt.value = null;
};

const dismissBanner = () => {
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
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.95), rgba(15, 23, 42, 0.95));
  border: 1px solid rgba(99, 102, 241, 0.5);
  border-radius: var(--radius-lg);
  box-shadow: 0 10px 30px rgba(0,0,0,0.5);
}

.pwa-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.pwa-icon-badge {
  font-size: 1.5rem;
  padding: 0.4rem;
  background: rgba(99, 102, 241, 0.2);
  border-radius: 10px;
}

.pwa-text h4 {
  font-size: 0.88rem;
  font-weight: 800;
  color: #ffffff;
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
  padding: 0.45rem 0.85rem;
  font-size: 0.78rem;
  font-weight: 800;
}

.btn-pwa-dismiss {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 1rem;
  cursor: pointer;
  padding: 0.2rem;
}
.btn-pwa-dismiss:hover {
  color: var(--accent-danger);
}
</style>
