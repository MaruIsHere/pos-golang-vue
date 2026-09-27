<script setup lang="ts">
defineProps({
  title: {
    type: String,
    default: ''
  },
  bodyClass: {
    type: String,
    default: ''
  }
})
</script>

<template>
  <div class="bg-bg-card/95 backdrop-blur-md border border-border-color rounded-[14px] shadow-sm transition-colors duration-200 ease-in-out dark:bg-slate-800/90 dark:border-slate-700">
    
    <!-- Bagian Header (Otomatis muncul jika ada Title atau dimasukkan ke slot #header) -->
    <div v-if="$slots.header || title" class="px-5 py-4 border-b border-border-color dark:border-slate-700">
      <slot name="header">
        <h3 class="font-bold text-lg text-tx-primary dark:text-slate-100 flex items-center gap-2">
          {{ title }}
        </h3>
      </slot>
    </div>

    <!-- Bagian Utama (Body) -->
    <div :class="['px-5 py-4', bodyClass]">
      <slot />
    </div>

    <!-- Bagian Bawah (Otomatis muncul jika slot #footer dipanggil) -->
    <div v-if="$slots.footer" class="px-5 py-4 border-t border-border-color bg-slate-50/50 rounded-b-[14px] dark:bg-slate-900/50 dark:border-slate-700">
      <slot name="footer" />
    </div>
  </div>
</template>
