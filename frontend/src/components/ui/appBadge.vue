<script setup lang="ts">
import { computed } from "vue";
import { cn } from "../../utils/cn";

export type BadgeVariant =
  "success" | "danger" | "warning" | "info" | "neutral" | "purple";

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
  success: "border text-[#047857] border-",
  danger: "",
  warning: "",
  info: "",
  neutral: "",
  purple: "",
};
</script>
