<template>
  <span :class="classes">
    <slot />
  </span>
</template>

<script setup>
import { computed } from "vue"

const props = defineProps({
  variant: {
    type: String,
    default: "default",
  },
})

const base =
  "inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-semibold uppercase tracking-[0.2em]"

const variants = {
  default: "border-border bg-white/80 text-foreground",
  outline: "border-border bg-transparent text-foreground",
  success: "border-emerald-200 bg-emerald-100 text-emerald-700",
  warning: "border-amber-200 bg-amber-100 text-amber-700",
}

const classes = computed(() => `${base} ${variants[props.variant] || variants.default}`)
</script>