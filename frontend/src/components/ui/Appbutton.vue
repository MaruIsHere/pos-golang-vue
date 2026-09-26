<template>
  <button :class="buttonClass">
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed, useAttrs } from "vue";
import { cn } from "../../utils/cn";

export type ButtonVariant = "primary" | "success" | "secondary"| "danger";
export type ButtonSize = "sm" | "md" | "lg";

interface Props {
  variant?: ButtonVariant;
  size?: ButtonSize;
  class?: any;
}

const props = withDefaults(defineProps<Props>(), {
  variant: "primary",
  size: "md",
});

const variantStyles: Record<ButtonVariant, string> = {
  primary:
    "bg-accent-primary text-white shadow-sm shadow-indigo-600/30 hover:bg-accent-primary-hover hover:shadow-md hover:shadow-indigo-600/40 active:transform-y-0.5",
  success:
    "bg-accent-secondary text-white shadow-sm shadow-indigo-500/30 hover:bg-accent-secondary-hover hover:shadow-md hover:shadow-indigo-500/40",
  secondary:
    "bg-bg-secondary text-text-primary border-border-color shadow-shadow-sm hover:bg-bg-card-hover hover:border-border-color",
  danger: "bg-accent-danger text-white hover:bg-[#b91c1c]",
};

const styleSize: Record<ButtonSize, string> = {
  sm: "px-3 py-1.5 text-xs",
  md: "px-4 py-2 text-sm",
  lg: "px-6 py-3.5 text-base",
};

const buttonClass = computed(() =>
  cn(
    // 1. Base style (default)
    "inline-flex items-center justify-center gap-2 py-2.5 px-5 font-font-family text-base font-semibold rounded-md border-solid border border-transparent cursor-pointer transition-all select-none",
    // 2.
    variantStyles[props.variant],
    styleSize[props.size],

    props.class,
  ),
);
</script>
