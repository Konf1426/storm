<script setup lang="ts">
import { computed } from "vue"
import { useNav } from "@slidev/client"

const { currentPage, total } = useNav()

const label = computed(() => {
  const current = String(currentPage.value).padStart(2, "0")
  const count = String(total.value).padStart(2, "0")
  return `${current} / ${count}`
})
</script>

<template>
  <div class="deck-page-number" aria-hidden="true">
    {{ label }}
  </div>
</template>
