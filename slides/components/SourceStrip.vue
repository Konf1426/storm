<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    items: string[]
    title?: string
  }>(),
  {
    title: "Sources internes",
  },
)
</script>

<template>
  <div class="source-strip">
    <span class="source-title">{{ props.title }}</span>
    <span class="source-item" v-for="item in props.items" :key="item">
      {{ item }}
    </span>
  </div>
</template>
