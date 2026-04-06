<script setup lang="ts">
import { presentationPlan } from "../data/content"

const icons = [
  "carbon:flow-data",           // Architecture
  "carbon:settings-check",      // Choix techniques
  "carbon:notification",        // Produit
  "carbon:play-filled-alt",     // Démo live
  "carbon:security",            // Sécurité
  "carbon:chart-line-smooth",   // Performance
  "carbon:cloud-service-management", // Cloud
  "carbon:presentation-file"    // Conclusion
]
</script>

<template>
  <div class="page-shell summary-shell">
    <div class="header-minimal">
      <p class="eyebrow">02 — Sommaire</p>
      <h2 class="section-title">Structure de la Soutenance</h2>
    </div>

    <div class="summary-layout-new">
      <div class="agenda-col">
        <div class="agenda-timeline">
          <article v-for="(item, idx) in presentationPlan" :key="item.index" class="agenda-card">
            <div class="agenda-icon-box">
              <img :src="`https://api.iconify.design/${icons[idx]}.svg?color=%230d9fb8`" class="ag-icon"/>
            </div>
            <div class="agenda-content">
              <div class="ag-index">{{ item.index }}</div>
              <div class="ag-text">
                <h3 class="ag-title">{{ item.title }}</h3>
              </div>
            </div>
          </article>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.summary-shell {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
}

.header-minimal {
  text-align: left;
  margin-bottom: 1.5rem;
}

.summary-layout-new {
  display: flex;
  justify-content: flex-start;
  align-items: flex-start;
}

.agenda-col {
  width: 100%;
  max-width: 960px;
}

.agenda-timeline {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.agenda-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  background: white;
  padding: 0.78rem 1rem;
  border-radius: 12px;
  border: 1px solid #f1f5f9;
  box-shadow: 0 2px 8px rgba(0,0,0,0.02);
  transition: transform 0.2s ease, border-color 0.2s ease;
}

.agenda-card:hover {
  border-color: var(--storm-cyan);
  transform: translateX(4px);
}

.agenda-icon-box {
  width: 32px;
  height: 32px;
  background: var(--storm-cyan-soft);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ag-icon { width: 18px; }

.agenda-content {
  display: flex;
  align-items: center;
  gap: 1.2rem;
  flex-grow: 1;
}

.ag-index {
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.8rem;
  font-weight: 800;
  color: #94a3b8;
  width: 25px;
}

.ag-text {
  display: flex;
  flex-direction: column;
}

.ag-title {
  font-size: 0.96rem;
  font-weight: 700;
  color: var(--storm-ink);
  margin: 0;
}
</style>
