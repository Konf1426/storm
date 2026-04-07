<script setup lang="ts">
// Vue component for Performance Results (Storm Day) - Unified Table Version
</script>

<template>
  <div class="page-shell">
    <div class="header-minimal">
      <p class="eyebrow">11 — Performance et résilience</p>
      <h2 class="section-title">Tenue sous charge : le système dégrade gracieusement</h2>
    </div>

    <div class="results-container">
      <!-- Unified Results Table -->
      <article class="clean-card shadow-sm table-card">
        <header class="card-head">
          <div class="icon-box blue-bg"><img src="https://api.iconify.design/carbon/chart-cluster-bar.svg" class="icon-sm"/></div>
          <h3 class="card-title">Métriques de performance par phase</h3>
        </header>
        
        <div class="table-wrapper">
          <table class="perf-table">
            <thead>
              <tr>
                <th>Métrique</th>
                <th>Warmup <span class="th-sub">(Normal)</span></th>
                <th>Spike 1 <span class="th-sub">(Charge)</span></th>
                <th>Chaos <span class="th-sub">(Pannes)</span></th>
                <th>Seuil / Cible</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td class="metric-name">
                  <div class="metric-icon orange"><img src="https://api.iconify.design/carbon/timer.svg"/></div>
                  Latence p95 <span class="unit">(ms)</span>
                </td>
                <td>6.3ms</td>
                <td>81.3ms</td>
                <td class="highlight-orange">141.5ms</td>
                <td class="target">< 200ms</td>
              </tr>
              <tr>
                <td class="metric-name">
                  <div class="metric-icon cyan"><img src="https://api.iconify.design/carbon/cloud-service-management.svg"/></div>
                  Débit HTTP <span class="unit">(req/s)</span>
                </td>
                <td>923</td>
                <td class="highlight-cyan">2 502</td>
                <td>2 008</td>
                <td class="target">Scaling OK</td>
              </tr>
              <tr>
                <td class="metric-name">
                  <div class="metric-icon purple"><img src="https://api.iconify.design/carbon/connect.svg"/></div>
                  Débit Message WS <span class="unit">(msg/s)</span>
                </td>
                <td>23.4K</td>
                <td>145K</td>
                <td class="highlight-purple">171.4K</td>
                <td class="target">Linéaire</td>
              </tr>
              <tr>
                <td class="metric-name">
                  <div class="metric-icon green"><img src="https://api.iconify.design/carbon/checkmark-outline.svg"/></div>
                  Taux d'erreur <span class="unit">(%)</span>
                </td>
                <td>0%</td>
                <td>0%</td>
                <td class="highlight-green">0%</td>
                <td class="target">0 % Error</td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <div class="table-footer">
          <p class="text-xs">Les métriques confirment que malgré l'injection de pannes (Chaos), le système maintient ses SLOs sans crash ni fuite.</p>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.page-shell { display: flex; flex-direction: column; height: 100%; }
.header-minimal { margin-bottom: 0.4rem; }

.results-container {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  flex-grow: 1;
  min-height: 0;
}

.clean-card {
  background: white; border-radius: 12px; padding: 0.8rem 1.2rem;
  display: flex; flex-direction: column; gap: 0.4rem;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.03); border: 1px solid rgba(0,0,0,0.04);
}

.table-card {
  flex-grow: 1;
  padding: 0.8rem 1.2rem;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.card-head { display: flex; align-items: center; gap: 0.6rem; margin-bottom: 0.2rem; }
.card-title { font-size: 0.9rem; font-weight: 700; color: var(--storm-ink); margin: 0; }

.icon-box { width: 28px; height: 28px; border-radius: 6px; display: flex; align-items: center; justify-content: center; }
.icon-sm { width: 16px; filter: grayscale(1) brightness(0) invert(1); }
.blue-bg { background: var(--storm-ink); }

.table-wrapper {
  flex-grow: 1;
  display: flex;
  align-items: center;
  min-height: 0;
}

.perf-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  font-size: 0.8rem;
}

.perf-table th {
  text-align: left;
  padding: 0.4rem 0.6rem;
  border-bottom: 2px solid #f1f5f9;
  color: #64748b;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.6rem;
  letter-spacing: 0.05em;
}

.th-sub {
  display: block;
  font-weight: 400;
  text-transform: none;
  font-size: 0.5rem;
  color: #94a3b8;
}

.perf-table td {
  padding: 0.5rem 0.6rem;
  border-bottom: 1px solid #f1f5f9;
  color: var(--storm-ink);
  font-weight: 500;
}

.metric-name {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 700 !important;
}

.metric-icon {
  width: 20px;
  height: 20px;
  border-radius: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.metric-icon img { width: 12px; filter: invert(1); }

.orange { background: #f97316; }
.cyan { background: #06b6d4; }
.purple { background: #8b5cf6; }
.green { background: #10b981; }

.unit {
  font-weight: 400;
  color: #94a3b8;
  font-size: 0.7rem;
  margin-left: 0.15rem;
}

.highlight-orange { color: #f97316 !important; font-weight: 800 !important; background: rgba(249, 115, 22, 0.04); }
.highlight-cyan { color: #06b6d4 !important; font-weight: 800 !important; background: rgba(6, 182, 212, 0.04); }
.highlight-purple { color: #8b5cf6 !important; font-weight: 800 !important; background: rgba(139, 92, 246, 0.04); }
.highlight-green { color: #10b981 !important; font-weight: 800 !important; background: rgba(16, 185, 129, 0.04); }

.target {
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.65rem;
  color: #64748b !important;
  background: #f8fafc;
  text-align: center;
}

.table-footer {
  margin-top: 0.4rem;
  padding-top: 0.4rem;
  border-top: 1px dashed #e2e8f0;
}

.text-xs { font-size: 0.7rem; color: #64748b; line-height: 1.2; }

/* Status Summary Bottom Bar */
.status-summary {
  background: white; border: 1px solid #e2e8f0; border-radius: 12px;
  display: flex; align-items: center; justify-content: space-around; padding: 0.5rem 1rem;
}
.summary-item { display: flex; align-items: center; gap: 0.4rem; }
.sum-icon { width: 14px; }
.summary-item span { font-size: 0.7rem; color: #475569; }
.summary-item strong { color: var(--storm-ink); }
</style>
