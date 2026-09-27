<script setup lang="ts">
defineProps({
  headers: {
    type: Array as () => string[],
    required: true
  }
})
</script>

<template>
  <div class="w-full overflow-x-auto rounded-lg border border-border-color dark:border-slate-700 shadow-sm">
    <table class="w-full text-left text-sm text-text-secondary dark:text-slate-300 border-collapse">
      <!-- Judul Kolom Tabel -->
      <thead class="bg-slate-50 dark:bg-slate-800 border-b border-border-color dark:border-slate-700 text-tx-primary dark:text-slate-100">
        <tr>
          <th v-for="(header, index) in headers" :key="index" class="px-4 py-3 font-semibold whitespace-nowrap">
            {{ header }}
          </th>
        </tr>
      </thead>
      <!-- Isi Baris Tabel -->
      <tbody class="divide-y divide-border-color dark:divide-slate-700 bg-bg-secondary dark:bg-slate-900">
        <slot />
      </tbody>
    </table>
  </div>
</template>
