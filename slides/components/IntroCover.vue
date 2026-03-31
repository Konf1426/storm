<script setup lang="ts">
import { deckMeta, objectives } from "../data/content"
</script>

<template>
  <div class="page-shell intro-shell">
    <div class="intro-centered">
      <div class="intro-eyebrow-line">
        <span class="eyebrow">MT4 — Soutenance de Projet</span>
        <div class="team-pills">
          <span v-for="member in deckMeta.teamMembers" :key="member" class="team-pill">{{ member }}</span>
        </div>
      </div>

      <div class="hero-block">
        <h1 class="hero-title-large">STORM</h1>
        <p class="hero-subtitle-large">{{ deckMeta.subtitle }}</p>
        <div class="hero-divider"></div>
        <p class="hero-positioning">{{ deckMeta.positioning }}</p>
      </div>

      <div class="targets-grid">
        <article v-for="obj in objectives" :key="obj.label" class="clean-card target-card">
          <p class="target-label">{{ obj.label }}</p>
          <p class="target-value" :class="`text-${obj.tone || 'cyan'}`">{{ obj.value }}</p>
          <p class="target-note">{{ obj.note }}</p>
        </article>
      </div>

      <div class="tech-stack-footer">
        <div class="tech-item"><img src="https://api.iconify.design/logos/go.svg" class="tech-icon"/></div>
        <div class="tech-item"><img src="https://api.iconify.design/logos/nats-icon.svg" class="tech-icon"/></div>
        <div class="tech-item"><img src="https://api.iconify.design/logos/postgresql.svg" class="tech-icon"/></div>
        <div class="tech-item"><img src="https://api.iconify.design/logos/redis.svg" class="tech-icon"/></div>
        <div class="tech-item"><img src="https://api.iconify.design/logos/prometheus.svg" class="tech-icon"/></div>
        <div class="tech-item"><img src="https://api.iconify.design/logos/grafana.svg" class="tech-icon"/></div>
        <div class="tech-item"><img src="https://api.iconify.design/logos/docker-icon.svg" class="tech-icon"/></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.intro-shell {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.intro-centered {
  max-width: 960px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 2.5rem;
}

.intro-eyebrow-line {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
  align-items: center;
}

.team-pills {
  display: flex;
  gap: 0.5rem;
}
.team-pill {
  font-size: 0.65rem;
  font-weight: 700;
  color: #64748b;
  background: white;
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  border: 1px solid #e2e8f0;
}

.hero-block {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.hero-title-large {
  font-size: 6rem;
  font-weight: 900;
  letter-spacing: -0.06em;
  line-height: 0.8;
  color: var(--storm-ink);
  background: linear-gradient(135deg, #16222b 0%, #3f5560 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.hero-subtitle-large {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--storm-cyan);
  text-transform: uppercase;
  letter-spacing: 0.2em;
  margin-top: 0.5rem;
}

.hero-divider {
  width: 60px;
  height: 4px;
  background: var(--storm-orange);
  margin: 1.5rem 0;
  border-radius: 2px;
}

.hero-positioning {
  font-size: 0.95rem;
  color: #475569;
  max-width: 700px;
  line-height: 1.5;
  font-weight: 500;
}

.targets-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
  width: 100%;
}

.target-card {
  background: white;
  padding: 1.2rem;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  box-shadow: 0 4px 20px rgba(0,0,0,0.03);
  border: 1px solid rgba(0,0,0,0.04);
}

.target-label {
  font-size: 0.6rem;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #94a3b8;
}

.target-value {
  font-size: 1.6rem;
  font-weight: 800;
  letter-spacing: -0.04em;
  margin: 0.2rem 0;
}

.target-note {
  font-size: 0.65rem;
  color: #64748b;
  line-height: 1.3;
}

.text-cyan { color: var(--storm-cyan); }
.text-orange { color: var(--storm-orange); }
.text-slate { color: var(--storm-slate); }

.tech-stack-footer {
  display: flex;
  gap: 1.5rem;
  align-items: center;
  opacity: 0.6;
  margin-top: 1rem;
}

.tech-icon {
  height: 20px;
  filter: grayscale(1);
  transition: filter 0.3s ease, opacity 0.3s ease;
}

.tech-item:hover .tech-icon {
  filter: grayscale(0);
  opacity: 1;
}
</style>

