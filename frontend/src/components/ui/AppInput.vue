<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'

interface Props {
  modelValue?: string | number
  label?: string
  type?: string
  placeholder?: string
  error?: string
  required?: boolean
  class?: any // Untuk custom styling via cn()
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  label: '',
  type: 'text',
  placeholder: '',
  error: '',
  required: false
})

defineEmits(['update:modelValue'])

const inputClasses = computed(() => cn(
  'w-full px-3.5 py-2.5 text-[0.9rem] rounded-[10px] transition-all duration-150 ease-in-out focus:outline-none focus:ring-[3px]',
  'bg-white dark:bg-slate-800 text-tx-primary border',
  'dark:bg-slate-900 dark:border-slate-700 dark:text-slate-50',
  props.error 
    ? 'border-accent-danger focus:border-accent-danger focus:ring-accent-danger/20' 
    : 'border-slate-200 dark:border-slate-700 focus:border-accent-primary focus:ring-accent-primary/12',
  props.class
))
</script>

<template>
  <div class="flex flex-col gap-1.5 w-full">
    <label v-if="label" class="text-sm font-semibold text-slate-500 dark:text-slate-400 dark:text-slate-300">
      {{ label }} <span v-if="required" class="text-accent-danger">*</span>
    </label>
    <input
      :type="type"
      :placeholder="placeholder"
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      :required="required"
      :class="inputClasses"
    />
    <span v-if="error" class="text-xs text-accent-danger font-medium">{{ error }}</span>
  </div>
</template>
