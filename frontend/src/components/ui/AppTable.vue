<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'

interface Props {
  headers: string[]
  class?: any
}

const props = defineProps<Props>()

const wrapperClasses = computed(() => cn(
  'w-full overflow-x-auto rounded-lg border border-slate-200 dark:border-slate-700 shadow-sm',
  props.class
))
</script>

<template>
  <div :class="wrapperClasses">
    <table class="data-table text-sm">
      <thead>
        <tr>
          <th v-for="(header, index) in headers" :key="index" class="px-4 py-3 font-semibold whitespace-nowrap">
            {{ header }}
          </th>
        </tr>
      </thead>
      <tbody class="bg-white dark:bg-slate-800">
        <slot />
      </tbody>
    </table>
  </div>
</template>
