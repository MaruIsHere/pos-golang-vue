<template>
  <div class="h-screen w-screen flex flex-col bg-slate-100/70 dark:bg-slate-950 overflow-hidden">
    <Navbar v-if="authStore.isAuthenticated" />

    <PwaInstallBanner />

    <main class="flex flex-1 overflow-hidden relative" :class="{ 'md:pl-20 md:pt-16 pt-12 pb-16 md:pb-0': authStore.isAuthenticated }">
      <div class="flex-1 h-full w-full overflow-y-auto overflow-x-hidden flex flex-col relative" :class="{ 'p-3 md:p-5': authStore.isAuthenticated }">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" class="flex-1 w-full min-h-0" />
          </transition>
        </router-view>
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
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease-out;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
