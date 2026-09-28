<script setup lang="ts">
import { computed } from "vue";
import { cn } from "../../utils/cn";

export type BadgeVariant = "success" | "danger" | "warning" | "info" | "neutral" | "purple";

export type BadgeSize = "sm" | "md";

interface Props {
  variant?: BadgeVariant;
  size?: BadgeSize;
  dot?: boolean;
  class?: string;
}

const props = withDefaults(defineProps<Props>(), {
  variant: "neutral",
  size: "md",
  dot: false,
});

const variantStyles: Record<BadgeVariant, string> = {
  success:
    "bg-emerald-50  text-emerald-700 border border-emerald-200/60 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20",
  danger:
    "bg-red-50 text-red-700 border border-red-200/60 dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/20",
  warning:
    "bg-amber-50 text-amber-700 border border-amber-200/60 dark:bg-amber-500/10 dark:text-amber-400 dark:border-amber-500/20",
  info: "bg-blue-50 text-blue-700 border border-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:border-blue-500/20",
  neutral:
    "bg-slate-50 text-slate-700 border border-slate-200/60 dark:bg-slate-500/10 dark:text-slate-400 dark:border-slate-500/20",
  purple:
    "bg-purple-50 text-purple-700 border border-purple-200/60 dark:bg-purple-500/10 dark:text-purple-400 dark:border-purple-500/20",
};

const dotStyles: Record<BadgeVariant, string> = {
  success: "bg-emerald-500",
  danger: "bg-red-500",
  warning: "bg-amber-500",
  info: "bg-blue-500",
  neutral: "bg-slate-500",
  purple: "bg-pulple-500",
};

const sizeStyles: Record<BadgeSize, string> = {
  sm: "px-2 py-0.5 text-xs",
  md: "px-2.5 py-1 text-xs",
};

const badgeClasses = computed(() =>
  cn(
    "inline-flex items-center py-2 px-1 rounded-full text-xs font-semibold",
    variantStyles[props.variant],
    sizeStyles[props.size],
    props.class,
  ),
);
</script>

<template>
  <span :class="badgeClasses">
    <span
      v-if="dot"
      class="h-1.5 w-1.5 rounded-full shrink-0"
      :class="dotStyles[variant]"
      aria-hidden="true"
    />
    <slot />
  </span>
</template>
