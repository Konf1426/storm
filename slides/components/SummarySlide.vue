<script setup lang="ts">
import { presentationPlan } from "../data/content"

const icons = [
  "carbon:pedestrian-child", // Démarche
  "carbon:flow-data",       // Architecture
  "carbon:settings-check",   // Choix techniques
  "carbon:security",         // Sécurité
  "carbon:chart-line-smooth",// Performance
  "carbon:cloud-service-management", // Scalabilité
  "carbon:presentation-file" // Conclusion
]
</script>

<template>
  <div class="page-shell summary-shell">
    <div class="header-minimal">
      <p class="eyebrow">02 — Sommaire</p>
      <h2 class="section-title">Structure de la Soutenance</h2>
    </div>

    <div class="summary-layout-new">
      <!-- Agenda List -->
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
                <p class="ag-detail">{{ item.detail }}</p>
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
  justify-content: center;
}

.header-minimal {
  text-align: center;
  margin-bottom: 2rem;
}

.summary-layout-new {
  display: flex;
  justify-content: center;
  align-items: center;
}

.agenda-col {
  width: 100%;
  max-width: 900px;
}

.agenda-timeline {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.agenda-card:last-child {
  grid-column: 1 / -1;
  max-width: 440px;
  margin: 0 auto;
}

.agenda-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  background: white;
  padding: 0.6rem 1rem;
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
  gap: 0.1rem;
}

.ag-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: var(--storm-ink);
  margin: 0;
}

.ag-detail {
  font-size: 0.72rem;
  color: #64748b;
  margin: 0;
}

/* Sidebar */
.summary-sidebar {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.info-card {
  padding: 1.5rem;
  background: linear-gradient(180deg, #ffffff 0%, #f8fafc 100%);
}

.sidebar-title {
  font-size: 0.8rem;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #475569;
  margin-bottom: 1.2rem;
}

.sidebar-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.sidebar-list li {
  display: flex;
  align-items: flex-start;
  gap: 0.8rem;
}

.s-icon { width: 14px; margin-top: 3px; flex-shrink: 0; }
.sidebar-list li span {
  font-size: 0.78rem;
  font-weight: 500;
  color: #334155;
  line-height: 1.3;
}

.timing-badge {
  background: white;
  border: 1px solid #e2e8f0;
  padding: 0.8rem 1rem;
  border-radius: 999px;
  display: flex;
  align-items: center;
  gap: 0.6rem;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 600;
  color: #64748b;
}

.t-icon { width: 14px; }
</style>

