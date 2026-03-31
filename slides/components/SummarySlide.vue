<script setup lang="ts">
import { objectives, presentationPlan } from "../data/content"
</script>

<template>
  <div class="page-shell summary-shell">
    <p class="eyebrow">Sommaire</p>
    <h2 class="section-title">Plan de la présentation</h2>

    <div class="agenda-list">
      <article v-for="item in presentationPlan" :key="item.index" class="agenda-item">
        <div class="agenda-number">{{ item.index }}</div>
        <div class="agenda-copy">
          <p class="agenda-title">{{ item.title }}</p>
          <p class="agenda-detail">{{ item.detail }}</p>
        </div>
      </article>
    </div>
  </div>
</template>
