<template>
  <button
    :type="type"
    :class="classes"
    :disabled="disabled"
  >
    <slot />
  </button>
</template>

<script setup>
import { computed } from "vue"

const props = defineProps({
  variant: {
    type: String,
    default: "primary",
  },
  type: {
    type: String,
    default: "button",
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const base =
  "inline-flex items-center justify-center rounded-full px-5 py-2 text-sm font-medium transition focus:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:opacity-60 disabled:pointer-events-none"

const variants = {
  primary: "bg-primary text-primary-foreground shadow-glow hover:shadow-lg",
  secondary: "bg-muted text-foreground hover:bg-muted/80",
  outline: "border border-border bg-white/70 text-foreground hover:bg-white",
  ghost: "text-foreground hover:bg-muted/70",
}

const classes = computed(() => `${base} ${variants[props.variant] || variants.primary}`)
</script>