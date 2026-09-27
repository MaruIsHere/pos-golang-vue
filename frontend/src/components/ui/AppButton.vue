<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'

export type ButtonVariant = 'primary' | 'secondary' | 'danger' | 'ghost' | 'outline'
export type ButtonSize = 'sm' | 'md' | 'lg' | 'icon'
export type ButtonType = 'button' | 'submit' | 'reset'

interface Props {
  variant?: ButtonVariant
  size?: ButtonSize
  isLoading?: boolean
  type?: ButtonType
  disabled?: boolean
  class?: any // Menerima custom class dari luar
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  isLoading: false,
  type: 'button',
  disabled: false
})

const variantStyles: Record<ButtonVariant, string> = {
  primary: 'bg-accent-primary text-white hover:bg-accent-primary-hover focus:ring-accent-primary',
  secondary: 'bg-bg-secondary text-text-primary border border-border-color hover:bg-bg-card-hover focus:ring-text-secondary',
  danger: 'bg-accent-danger text-white hover:bg-red-700 focus:ring-accent-danger',
  ghost: 'bg-transparent text-text-primary hover:bg-bg-card-hover dark:hover:bg-slate-800',
  // Varian khusus meniru .theme-toggle-btn Bapak
  outline: 'bg-bg-primary text-text-primary border border-border-color hover:bg-bg-card-hover focus:ring-text-secondary'
}

const sizeStyles: Record<ButtonSize, string> = {
  sm: 'h-8 px-3 text-xs rounded',
  md: 'h-10 px-4 py-2 text-sm rounded-md',
  lg: 'h-12 px-8 text-base rounded-lg',
  // Ukuran khusus icon 36x36 px (h-9 = 36px, w-9 = 36px)
  icon: 'h-9 w-9 p-0 flex items-center justify-center rounded-[10px]'
}

// Gabungkan base class, variant, size, dan class tambahan dari luar
const buttonClasses = computed(() =>
  cn(
    // 1. Base style
    'inline-flex items-center justify-center font-medium transition-all duration-150 ease-in-out focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none [-webkit-tap-highlight-color:transparent]',
    // 2. Styles dari props
    variantStyles[props.variant],
    sizeStyles[props.size],
    // 3. Class luar yang diketik di parent component
    props.class
  )
)
</script>

<template>
  <button :type="type" :disabled="disabled || isLoading" :class="buttonClasses">
    <svg v-if="isLoading" class="animate-spin" :class="size === 'icon' ? 'absolute h-4 w-4' : '-ml-1 mr-2 h-4 w-4'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    <slot />
  </button>
</template>
