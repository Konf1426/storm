<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    items: string[]
    reveal?: boolean
  }>(),
  {
    reveal: false,
  },
)
</script>

<template>
  <div class="bullet-stack">
    <template v-if="props.reveal">
      <div v-for="item in props.items" :key="item" v-click class="bullet-row">
        {{ item }}
      </div>
    </template>
    <template v-else>
      <div v-for="item in props.items" :key="item" class="bullet-row">
        {{ item }}
      </div>
    </template>
  </div>
</template>
