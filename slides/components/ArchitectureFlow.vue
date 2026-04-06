<script setup lang="ts">
// Interactive architecture flowchart built entirely with CSS and Vue.
</script>

<template>
  <div class="page-shell arch-flow-shell">
    <div class="arch-header">
      <p class="eyebrow">04 — Architecture</p>
      <h2 class="section-title">Un système compact, observable et défendable</h2>
    </div>

    <div class="flow-canvas">
      <div class="flow-grid">
        <!-- Ligne 1 : Observabilité -->
        <div class="flow-cell pos-obs">
          <div class="flow-node node-obs">
            <div class="node-glow cyan"></div>
            <div class="node-icons">
              <img src="/icons/prometheus.svg" class="obs-icon" style="width: 18px; height: 18px; flex-shrink: 0;" />
              <img src="/icons/grafana.svg" class="obs-icon" style="width: 18px; height: 18px; flex-shrink: 0;" />
            </div>
            <div class="node-info">
              <p class="node-name">Observabilité</p>
              <p class="node-desc">Prometheus & Grafana</p>
            </div>
          </div>
          <!-- Line connecting Prom to Gateway (flowing from Gateway UP to Prom) -->
          <div class="flow-link link-down from-obs">
            <div class="particle bg-cyan"></div>
          </div>
        </div>

        <!-- Ligne 2 : Le Pipeline Principal -->
        <div class="flow-cell pos-vue">
          <div class="flow-node">
            <img src="/icons/vue.svg" class="node-icon" style="width: 26px; height: 26px; flex-shrink: 0;" />
            <div class="node-info">
              <p class="node-name">Client</p>
              <p class="node-desc">Vue 3 / Vite</p>
            </div>
          </div>
          <div class="flow-link link-right">
            <div class="particle bg-cyan"></div>
          </div>
        </div>

        <div class="flow-cell pos-gw">
          <div class="flow-node">
            <img src="/icons/go.png" class="node-icon" style="width: 26px; height: 26px; flex-shrink: 0;" />
            <div class="node-info">
              <p class="node-name">Gateway</p>
              <p class="node-desc">WebSocket Go</p>
            </div>
          </div>
          <div class="flow-link link-right">
            <div class="particle bg-orange"></div>
          </div>
          <div class="flow-link link-down">
            <div class="particle bg-slate"></div>
          </div>
        </div>

        <div class="flow-cell pos-nats">
          <div class="flow-node">
            <img src="/icons/nats.png" class="node-icon" style="width: 26px; height: 26px; flex-shrink: 0;" />
            <div class="node-info">
              <p class="node-name">Bus NATS</p>
              <p class="node-desc">Pub/Sub mémoire</p>
            </div>
          </div>
          <div class="flow-link link-right">
            <div class="particle bg-cyan"></div>
          </div>
        </div>

        <div class="flow-cell pos-msg">
          <div class="flow-node">
            <img src="/icons/go.png" class="node-icon" style="width: 26px; height: 26px; flex-shrink: 0;" />
            <div class="node-info">
              <p class="node-name">Messages</p>
              <p class="node-desc">Worker asynchrone</p>
            </div>
          </div>
          <div class="flow-link link-down">
            <div class="particle bg-orange"></div>
          </div>
        </div>

        <!-- Ligne 3 : État & Persistance -->
        <div class="flow-cell pos-redis">
          <div class="flow-node">
            <img src="/icons/redis.svg" class="node-icon" style="width: 26px; height: 26px; flex-shrink: 0;" />
            <div class="node-info">
              <p class="node-name">Redis</p>
              <p class="node-desc">Présence volatile</p>
            </div>
          </div>
        </div>

        <div class="flow-cell pos-pg">
          <div class="flow-node">
            <img src="/icons/postgres.svg" class="node-icon" style="width: 26px; height: 26px; flex-shrink: 0;" />
            <div class="node-info">
              <p class="node-name">Postgres</p>
              <p class="node-desc">Stockage durable</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.arch-flow-shell {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 1rem;
}

.arch-header {
  display: grid;
  gap: 0.3rem;
  margin-bottom: 0.5rem;
}

.flow-canvas {
  flex-grow: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem 0;
}

.flow-grid {
  display: grid;
  grid-template-columns: 1fr 1.25fr 1.25fr 1fr;
  grid-template-rows: repeat(3, auto);
  gap: 1.8rem 2rem;
  width: 100%;
  max-width: 860px;
}

/* Placement */
.pos-obs { grid-column: 2; grid-row: 1; }
.pos-vue { grid-column: 1; grid-row: 2; }
.pos-gw { grid-column: 2; grid-row: 2; }
.pos-nats { grid-column: 3; grid-row: 2; }
.pos-msg { grid-column: 4; grid-row: 2; }
.pos-redis { grid-column: 2; grid-row: 3; }
.pos-pg { grid-column: 4; grid-row: 3; }

.flow-cell {
  position: relative;
  display: flex;
  justify-content: center;
  min-width: 0;
}

.flow-node {
  position: relative;
  width: 100%;
  background: var(--storm-surface);
  border: 1px solid rgba(22, 34, 43, 0.08);
  border-radius: 16px;
  box-shadow: var(--storm-shadow);
  backdrop-filter: blur(12px);
  padding: 0.6rem 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.65rem;
  z-index: 2;
}

.node-obs {
  justify-content: center;
  border-color: rgba(13, 159, 184, 0.2);
}

.node-glow {
  position: absolute;
  inset: 0;
  border-radius: inherit;
  opacity: 0.15;
  z-index: -1;
}

.node-glow.cyan { background: var(--storm-cyan); }

.node-icons {
  display: flex;
  gap: 0.3rem;
}

.obs-icon {
  width: 18px;
  height: 18px;
  object-fit: contain;
}

.node-icon {
  width: 26px;
  height: 26px;
  object-fit: contain;
  flex-shrink: 0;
}

.node-info {
  display: grid;
  gap: 0.1rem;
}

.node-name {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--storm-ink);
  line-height: 1.1;
}

.node-desc {
  font-size: 0.7rem;
  color: var(--storm-muted);
  line-height: 1.2;
}

/* Links & Animations */
.flow-link {
  position: absolute;
  background: rgba(22, 34, 43, 0.12);
  z-index: 1;
}

/* Rightward link */
.link-right {
  top: 50%;
  right: -2rem; /* Matches grid gap */
  width: 2rem;
  height: 2px;
  transform: translateY(-50%);
}

.link-right .particle {
  position: absolute;
  top: -2px; left: 0;
  width: 6px; height: 6px;
  border-radius: 50%;
  animation: flowAnimX 1.5s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

@keyframes flowAnimX {
  0% { transform: translateX(0); opacity: 0; }
  15% { opacity: 1; }
  85% { opacity: 1; }
  100% { transform: translateX(2rem); opacity: 0; }
}

/* Downward link */
.link-down {
  bottom: -1.8rem;
  left: 50%;
  width: 2px;
  height: 1.8rem;
  transform: translateX(-50%);
}

/* The particle going down */
.link-down .particle {
  position: absolute;
  top: 0; left: -2px;
  width: 6px; height: 6px;
  border-radius: 50%;
  animation: flowAnimY 1.5s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

@keyframes flowAnimY {
  0% { transform: translateY(0); opacity: 0; }
  15% { opacity: 1; }
  85% { opacity: 1; }
  100% { transform: translateY(1.8rem); opacity: 0; }
}

/* Reverse flow from Gateway to Obs */
.from-obs.link-down .particle {
  animation: flowAnimYReverse 1.5s cubic-bezier(0.4, 0, 0.2, 1) infinite;
  top: auto;
  bottom: 0;
}

@keyframes flowAnimYReverse {
  0% { transform: translateY(0); opacity: 0; }
  15% { opacity: 1; }
  85% { opacity: 1; }
  100% { transform: translateY(-1.8rem); opacity: 0; }
}

/* BG glows */
.bg-cyan { background: var(--storm-cyan); box-shadow: 0 0 6px rgba(13, 159, 184, 0.6); }
.bg-orange { background: var(--storm-orange); box-shadow: 0 0 6px rgba(239, 138, 41, 0.6); }
.bg-slate { background: var(--storm-slate); box-shadow: 0 0 6px rgba(63, 85, 96, 0.6); }
</style>
