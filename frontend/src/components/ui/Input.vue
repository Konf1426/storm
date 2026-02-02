<template>
  <input
    v-bind="$attrs"
    :class="classes"
  />
</template>

<script setup>
const classes =
  "mt-2 w-full rounded-2xl border border-border bg-white/90 px-4 py-2 text-sm text-foreground shadow-sm focus:border-primary focus:outline-none focus:ring-2 focus:ring-ring/30"
</script>