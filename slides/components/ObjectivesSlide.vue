<script setup lang="ts">
const phases = [
  {
    number: "01",
    title: "Architecture",
    detail: "Infrastructure as Code, CI/CD, contrat d'API",
    accent: "var(--storm-cyan)",
    bg: "var(--storm-cyan-soft)",
  },
  {
    number: "02",
    title: "Développement",
    detail: "Services Go, couverture > 80 %, scans de sécurité",
    accent: "var(--storm-orange)",
    bg: "var(--storm-orange-soft)",
  },
  {
    number: "03",
    title: "Performance",
    detail: "Profilage pprof, tests de charge k6, optimisation",
    accent: "var(--storm-cyan)",
    bg: "var(--storm-cyan-soft)",
  },
  {
    number: "04",
    title: "Résilience",
    detail: "Chaos engineering, Storm Day, post-mortems",
    accent: "var(--storm-orange)",
    bg: "var(--storm-orange-soft)",
  },
]
</script>

<template>
  <div class="page-shell obj-shell">
    <div class="obj-header">
      <p class="eyebrow">03 — Objectifs du projet</p>
      <h2 class="section-title">Notre Cap : De l'Architecture à la Résilience</h2>
    </div>

    <div class="obj-timeline">
      <article
        v-for="phase in phases"
        :key="phase.number"
        class="obj-phase"
      >
        <div
          class="obj-phase-number"
          :style="{ background: phase.bg, color: phase.accent }"
        >
          {{ phase.number }}
        </div>
        <div
          class="obj-phase-bar"
          :style="{ background: phase.accent }"
        />
        <div class="obj-phase-content">
          <p class="obj-phase-title">{{ phase.title }}</p>
          <p class="obj-phase-detail">{{ phase.detail }}</p>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.obj-shell {
  height: 100%;
  align-content: center;
  gap: 0.9rem;
}

.obj-header {
  display: grid;
  gap: 0.3rem;
}

.obj-timeline {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.7rem;
}

.obj-phase {
  display: grid;
  grid-template-columns: 44px 3px 1fr;
  gap: 0.7rem;
  align-items: start;
  border: 1px solid rgba(22, 34, 43, 0.08);
  border-radius: 20px;
  background: var(--storm-surface);
  box-shadow: var(--storm-shadow);
  backdrop-filter: blur(16px);
  padding: 0.82rem 0.9rem;
}

.obj-phase-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 14px;
  font-family: "IBM Plex Mono", monospace;
  font-size: 1rem;
  font-weight: 800;
  letter-spacing: -0.02em;
}

.obj-phase-bar {
  width: 3px;
  height: 100%;
  min-height: 36px;
  border-radius: 999px;
  opacity: 0.4;
  margin-top: 4px;
}

.obj-phase-content {
  display: grid;
  gap: 0.2rem;
}

.obj-phase-title {
  font-size: 0.92rem;
  font-weight: 700;
  color: var(--storm-ink);
}

.obj-phase-detail {
  font-size: 0.78rem;
  line-height: 1.3;
  color: var(--storm-muted);
}

@media (max-width: 720px) {
  .obj-timeline {
    grid-template-columns: 1fr;
  }
}
</style>
