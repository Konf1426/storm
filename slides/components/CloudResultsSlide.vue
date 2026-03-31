<script setup lang="ts">
// Cloud Results & Impact Slide - Premium Dashboard Design
import { azureLoadTest, budgetScenarios } from "../data/content"
const finalRun = azureLoadTest.stages[2] // 10 000 VUs
const steadyState = budgetScenarios.find(s => s.name === "Preuve Azure AKS")
</script>

<template>
  <div class="page-shell">
    <div class="header-minimal">
      <p class="eyebrow">12 — Résultats Cloud : Étape 2</p>
      <h2 class="section-title">Le Saut de Performance : 10 000 VUs sur Azure</h2>
    </div>

    <div class="results-dashboard">
      <!-- Main Metric: Success Rate -->
      <div class="main-metric shadow-xl count-up">
        <div class="metric-glow"></div>
        <div class="metric-label">Taux de Succès HTTP</div>
        <div class="metric-val">~100<span class="unit">%</span></div>
        <div class="metric-sub">Aucune erreur applicative sous charge</div>
      </div>

      <!-- Grid of secondary metrics -->
      <div class="metrics-bento">
        <article class="mini-stat-card shadow-sm fade-in delay-1">
          <div class="stat-label">Latence Message</div>
          <div class="stat-val">~85<span class="unit">ms</span></div>
          <p class="stat-desc">Publication NATS asynchrone</p>
        </article>

        <article class="mini-stat-card shadow-sm fade-in delay-2">
          <div class="stat-label">Débit WebSocket</div>
          <div class="stat-val">210<span class="unit">M</span></div>
          <p class="stat-desc">Messages cumulés sur le run</p>
        </article>

        <article class="mini-stat-card shadow-sm fade-in delay-3">
          <div class="stat-label">Login Médian</div>
          <div class="stat-val">144<span class="unit">ms</span></div>
          <p class="stat-desc">Validation JWT & Refresh Token</p>
        </article>
      </div>

      <!-- Budget & Scaling context -->
      <div class="budget-panel shadow-sm fade-in delay-4">
        <h4 class="panel-title">Dimensionnement & Coûts</h4>
        <div class="budget-grid">
           <div class="budget-item">
              <span class="bj-label">Cluster AKS</span>
              <span class="bj-val">26 vCPUs</span>
           </div>
           <div class="budget-item">
              <span class="bj-label">Gateway</span>
              <span class="bj-val">30 Réplicas</span>
           </div>
           <div class="budget-item">
              <span class="bj-label">PostgreSQL</span>
              <span class="bj-val">16 Cores</span>
           </div>
           <div class="budget-item highlight">
              <span class="bj-label">Coût du Run</span>
              <span class="bj-val">3,33 $/h</span>
           </div>
        </div>
      </div>
    </div>

    <div class="final-insight text-center">
       Le véritable saut de performance vient du <strong>découplage</strong> : la publication immédiate dans NATS libère le Gateway de l'attente PostgreSQL.
    </div>
  </div>
</template>

<style scoped>
.page-shell { display: flex; flex-direction: column; height: 100%; gap: 1rem; }
.header-minimal { margin-bottom: 0.2rem; }

.results-dashboard {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 1.2rem;
  flex-grow: 1;
}

/* Main Dashboard Card */
.main-metric {
  grid-row: span 1;
  background: #0f172a; color: white; border-radius: 20px;
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  position: relative; overflow: hidden; padding: 1.5rem;
  border: 1px solid rgba(255,255,255,0.1);
}

.metric-glow {
  position: absolute; top: -50%; left: -50%; width: 200%; height: 200%;
  background: radial-gradient(circle at center, rgba(6, 182, 212, 0.1) 0%, transparent 60%);
  pointer-events: none;
}

.metric-label { font-size: 0.8rem; font-weight: 700; color: var(--storm-cyan); text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 0.5rem; }
.metric-val { font-size: 4rem; font-weight: 850; letter-spacing: -0.04em; line-height: 1; }
.metric-val .unit { font-size: 1.5rem; color: rgba(255,255,255,0.4); margin-left: 0.2rem; }
.metric-sub { font-size: 0.75rem; color: #94a3b8; font-weight: 500; margin-top: 0.8rem; }

/* Secondary metrics grid */
.metrics-bento {
  display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;
}

.mini-stat-card {
  background: white; border-radius: 16px; padding: 1rem;
  display: flex; flex-direction: column; justify-content: center;
  border: 1px solid #f1f5f9;
}

.stat-label { font-size: 0.65rem; font-weight: 700; color: #64748b; text-transform: uppercase; margin-bottom: 0.4rem; }
.stat-val { font-size: 1.4rem; font-weight: 800; color: var(--storm-ink); letter-spacing: -0.02em; }
.stat-val .unit { font-size: 0.7rem; color: #94a3b8; margin-left: 0.1rem; }
.stat-desc { font-size: 0.6rem; color: #94a3b8; margin-top: 0.2rem; }

/* Budget Panel */
.budget-panel {
  grid-column: span 2;
  background: white; border-radius: 20px; padding: 1.5rem;
  display: flex; flex-direction: column; gap: 1rem;
  border: 1px solid #f1f5f9;
}

.panel-title { font-size: 0.75rem; font-weight: 800; color: #1e293b; text-transform: uppercase; letter-spacing: 0.05em; margin: 0; }

.budget-grid {
  display: grid; grid-template-columns: repeat(4, 1fr); gap: 1.5rem;
}

.budget-item { display: flex; flex-direction: column; gap: 0.2rem; border-left: 1px solid #f1f5f9; padding-left: 1rem; }
.bj-label { font-size: 0.65rem; color: #64748b; }
.bj-val { font-size: 0.95rem; font-weight: 800; color: var(--storm-ink); }
.budget-item.highlight .bj-val { color: #f97316; }

.final-insight { font-size: 0.95rem; color: #475569; font-style: italic; opacity: 0.8; margin-top: 0.5rem; line-height: 1.3; }

/* Animations */
.fade-in { opacity: 0; animation: fadeIn 0.6s ease forwards; }
@keyframes fadeIn { to { opacity: 1; } }
.count-up { animation: countUp 1s ease-out; }
@keyframes countUp { from { filter: blur(5px); opacity: 0; transform: scale(0.95); } to { filter: blur(0); opacity: 1; transform: scale(1); } }
.delay-1 { animation-delay: 0.2s; }
.delay-2 { animation-delay: 0.4s; }
.delay-3 { animation-delay: 0.6s; }
.delay-4 { animation-delay: 0.8s; }
</style>
