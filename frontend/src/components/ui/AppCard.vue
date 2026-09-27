<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'

interface Props {
  title?: string
  bodyClass?: any
  class?: any
}

const props = withDefaults(defineProps<Props>(), {
  title: ''
})

const cardClasses = computed(() => cn(
  'bg-bg-card/95 backdrop-blur-md border border-border-color rounded-[14px] shadow-sm transition-colors duration-200 ease-in-out dark:bg-slate-800/90 dark:border-slate-700',
  props.class
))

const bodyClasses = computed(() => cn('px-5 py-4', props.bodyClass))
</script>

<template>
  <div :class="cardClasses">
    <div v-if="$slots.header || title" class="px-5 py-4 border-b border-border-color dark:border-slate-700">
      <slot name="header">
        <h3 class="font-bold text-lg text-tx-primary dark:text-slate-100 flex items-center gap-2">
          {{ title }}
        </h3>
      </slot>
    </div>
    <div :class="bodyClasses">
      <slot />
    </div>
    <div v-if="$slots.footer" class="px-5 py-4 border-t border-border-color bg-slate-50/50 rounded-b-[14px] dark:bg-slate-900/50 dark:border-slate-700">
      <slot name="footer" />
    </div>
  </div>
</template>
