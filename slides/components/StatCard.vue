<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    label: string
    value: string
    note?: string
    tone?: "cyan" | "orange" | "slate"
  }>(),
  {
    note: "",
    tone: "cyan",
  },
)
</script>

<template>
  <article class="stat-card" :class="`stat-${props.tone}`">
    <p class="stat-eyebrow">{{ props.label }}</p>
    <p class="stat-value">{{ props.value }}</p>
    <p v-if="props.note" class="stat-note">{{ props.note }}</p>
  </article>
</template>
