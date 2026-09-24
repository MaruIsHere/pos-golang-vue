<template>
  <div class="app-container">
    <Navbar />

    <PwaInstallBanner />

    <main class="main-content">
      <div class="view-container">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import Navbar from './components/Navbar.vue';
import PwaInstallBanner from './components/PwaInstallBanner.vue';
import { useSettingsStore } from './stores/settings';
import { useTheme } from './composables/useTheme';

const settingsStore = useSettingsStore();
const { initTheme } = useTheme();

onMounted(() => {
  initTheme();
  settingsStore.fetchSettings();
});
</script>

<style scoped>
/* App layout container styled in styles.css */
</style>
