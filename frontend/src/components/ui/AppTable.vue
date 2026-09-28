<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'

interface Props {
  headers: string[]
  class?: any
}

const props = defineProps<Props>()

const wrapperClasses = computed(() => cn(
  'w-full overflow-x-auto rounded-lg border border-slate-200 dark:border-slate-700 dark:border-slate-700 shadow-sm',
  props.class
))
</script>

<template>
  <div :class="wrapperClasses">
    <table class="w-full text-left text-sm text-slate-500 dark:text-slate-400 dark:text-slate-300 border-collapse">
      <thead class="bg-slate-50 dark:bg-slate-800 border-b border-slate-200 dark:border-slate-700 dark:border-slate-700 text-tx-primary dark:text-slate-100">
        <tr>
          <th v-for="(header, index) in headers" :key="index" class="px-4 py-3 font-semibold whitespace-nowrap">
            {{ header }}
          </th>
        </tr>
      </thead>
      <tbody class="divide-y divide-border-color dark:divide-slate-700 bg-white dark:bg-slate-800 dark:bg-slate-900">
        <slot />
      </tbody>
    </table>
  </div>
</template>
