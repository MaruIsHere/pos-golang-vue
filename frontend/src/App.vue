<template>
  <div class="h-screen w-screen flex flex-col bg-slate-50 dark:bg-slate-900 overflow-hidden">
    <Navbar v-if="authStore.isAuthenticated" />

    <PwaInstallBanner />

    <main class="flex flex-1 overflow-hidden relative md:pl-20 md:pt-16 pt-12 pb-16 md:pb-0">
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
import { useAuthStore } from "./stores/auth";
import { useTheme } from "./composables/useTheme";

const settingsStore = useSettingsStore();
const authStore = useAuthStore();
const { initTheme } = useTheme();

onMounted(() => {
  initTheme();
  settingsStore.fetchSettings();
});
</script>

<style scoped>
/* App layout container styled in styles.css */
</style>
