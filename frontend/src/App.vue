<template>
  <div class="h-screen w-screen flex flex-col bg-bg-primary overflow-hidde">
    <Navbar />

    <PwaInstallBanner />

    <main class="flex flex-1 overflow-hidden relative">
      <div class="flex-1 overflow-y-auto p-5">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import Navbar from "./components/Navbar.vue";
import PwaInstallBanner from "./components/PwaInstallBanner.vue";
import { useSettingsStore } from "./stores/settings";
import { useTheme } from "./composables/useTheme";

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
