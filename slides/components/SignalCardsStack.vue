<script setup lang="ts">
interface SignalCardItem {
  title: string
  detail: string
}

const props = defineProps<{
  items: SignalCardItem[]
}>()
</script>

<template>
  <div class="bullet-stack">
    <article v-for="item in props.items" :key="item.title" class="bullet-card">
      <strong>{{ item.title }}</strong>
      <span>{{ item.detail }}</span>
    </article>
  </div>
</template>
