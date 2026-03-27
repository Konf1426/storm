<script setup lang="ts">
import { computed } from "vue"

interface MetricBarsItem {
  label: string
  value: number
  display: string
  meta?: string
  tone?: "cyan" | "orange" | "slate"
}

const props = defineProps<{
  items: MetricBarsItem[]
  max?: number
}>()

const normalized = computed(() => {
  const ceiling =
    props.max ?? Math.max(...props.items.map(item => item.value), 1)

  return props.items.map(item => ({
    ...item,
    tone: item.tone ?? "cyan",
    width: `${Math.max(10, (item.value / ceiling) * 100)}%`,
  }))
})
</script>

<template>
  <div class="metric-bars">
    <article v-for="item in normalized" :key="item.label" class="metric-row">
      <div class="metric-headline">
        <span>{{ item.label }}</span>
        <strong>{{ item.display }}</strong>
      </div>
      <div class="metric-track">
        <div class="metric-fill" :class="`tone-${item.tone}`" :style="{ width: item.width }" />
      </div>
      <p v-if="item.meta" class="metric-meta">{{ item.meta }}</p>
    </article>
  </div>
</template>
