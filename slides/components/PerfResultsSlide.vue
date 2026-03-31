<script setup lang="ts">
// Vue component for Performance Results (Storm Day)
</script>

<template>
  <div class="page-shell">
    <div class="header-minimal">
      <p class="eyebrow">09 — Résultats du Storm Day</p>
      <h2 class="section-title">Tenue sous charge : le système dégrade gracieusement</h2>
    </div>

    <div class="results-grid">
      <!-- Latency p95 -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box orange-bg"><img src="https://api.iconify.design/carbon/timer.svg" class="icon-sm"/></div>
          <h3 class="card-title">Latence p95 (Round-trip)</h3>
        </header>
        <div class="card-body">
          <div class="metric-trend">
            <div class="trend-item">
              <span class="trend-bar" style="height: 10%"></span>
              <span class="trend-val">6.3ms</span>
              <span class="trend-label">Warmup</span>
            </div>
            <div class="trend-item">
              <span class="trend-bar" style="height: 55%"></span>
              <span class="trend-val">81.3ms</span>
              <span class="trend-label">Spike 1</span>
            </div>
            <div class="trend-item">
              <span class="trend-bar peak" style="height: 100%"></span>
              <span class="trend-val text-orange">141.5ms</span>
              <span class="trend-label">Chaos</span>
            </div>
          </div>
          <p class="text-xs mt-auto">Latence maintenue sous le seuil critique (< 200ms) malgré l'injection de pannes.</p>
        </div>
      </article>

      <!-- Debit HTTP -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box cyan-bg"><img src="https://api.iconify.design/carbon/cloud-service-management.svg" class="icon-sm"/></div>
          <h3 class="card-title">Débit Requêtes HTTP</h3>
        </header>
        <div class="card-body">
          <div class="metric-value-box">
            <span class="huge-val">2 502</span>
            <span class="unit-val">req/s</span>
          </div>
          <div class="mini-stats">
            <div class="ms-row"><span>Error Rate</span> <span class="text-green">0.00%</span></div>
            <div class="ms-row"><span>Success (2xx)</span> <span>100%</span></div>
          </div>
          <p class="text-xs mt-auto">Stabilité parfaite de la couche Auth/Gateway sous un débit soutenu.</p>
        </div>
      </article>

      <!-- Debit WebSocket -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box purple-bg"><img src="https://api.iconify.design/carbon/connect.svg" class="icon-sm"/></div>
          <h3 class="card-title">Débit Messages WebSocket</h3>
        </header>
        <div class="card-body">
          <div class="metric-value-box">
            <span class="huge-val">171.4K</span>
            <span class="unit-val">msg/s</span>
          </div>
          <div class="mini-stats">
            <div class="ms-row"><span>Broadcasting</span> <span class="text-purple">NATS Pub/Sub</span></div>
            <div class="ms-row"><span>Connect p95</span> <span>28.4ms</span></div>
          </div>
          <p class="text-xs mt-auto">Le bus NATS garantit une distribution instantanée même en pic de charge.</p>
        </div>
      </article>

      <!-- Synthèse technique bottom bar -->
      <div class="status-summary full-span shadow-sm">
        <div class="summary-item">
           <img src="https://api.iconify.design/carbon/checkmark-filled.svg?color=%2310b981" class="sum-icon"/>
           <span><strong>Résilience</strong> : Aucun crash service constaté lors des arrêts de nœuds.</span>
        </div>
        <div class="summary-item">
           <img src="https://api.iconify.design/carbon/checkmark-filled.svg?color=%2310b981" class="sum-icon"/>
           <span><strong>Stabilité</strong> : Zero memory leak identifié lors des profils pprof post-charge.</span>
        </div>
        <div class="summary-item">
           <img src="https://api.iconify.design/carbon/checkmark-filled.svg?color=%2310b981" class="sum-icon"/>
           <span><strong>Scalabilité</strong> : Goulets identifiés (JSON/RSA) prêts pour optimisation.</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-shell { display: flex; flex-direction: column; height: 100%; }
.header-minimal { margin-bottom: 1rem; }

.results-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: 1fr auto;
  gap: 1.2rem;
  flex-grow: 1;
}

.full-span { grid-column: 1 / -1; }

.clean-card {
  background: white; border-radius: 12px; padding: 1.2rem;
  display: flex; flex-direction: column; gap: 1rem;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.03); border: 1px solid rgba(0,0,0,0.04);
}

.card-head { display: flex; align-items: center; gap: 0.6rem; }
.card-title { font-size: 0.9rem; font-weight: 700; color: var(--storm-ink); margin: 0; }

.icon-box { width: 26px; height: 26px; border-radius: 6px; display: flex; align-items: center; justify-content: center; }
.icon-sm { width: 14px; filter: grayscale(1) brightness(0) invert(1); }

.orange-bg { background: #f97316; }
.cyan-bg { background: #06b6d4; }
.purple-bg { background: #8b5cf6; }

.card-body { display: flex; flex-direction: column; flex-grow: 1; gap: 1rem; justify-content: center; }

/* Latency Trend */
.metric-trend {
  display: flex; align-items: flex-end; justify-content: space-around; height: 80px; padding-bottom: 0.5rem;
  border-bottom: 1px solid #f1f5f9;
}
.trend-item { display: flex; flex-direction: column; align-items: center; gap: 0.4rem; width: 30%; }
.trend-bar { width: 100%; background: #e2e8f0; border-radius: 4px 4px 0 0; transition: height 0.3s ease; }
.trend-bar.peak { background: #fdba74; }
.trend-val { font-family: monospace; font-size: 0.65rem; font-weight: 700; color: #475569; }
.trend-label { font-size: 0.55rem; font-weight: 600; text-transform: uppercase; color: #94a3b8; }
.text-orange { color: #f97316; }

/* Huge Value Display */
.metric-value-box { text-align: center; padding: 0.5rem 0; }
.huge-val { font-size: 2.2rem; font-weight: 800; color: var(--storm-ink); letter-spacing: -0.04em; line-height: 1; }
.unit-val { font-size: 0.8rem; font-weight: 600; color: #94a3b8; margin-left: 0.3rem; }

.mini-stats { display: flex; flex-direction: column; gap: 0.3rem; }
.ms-row { display: flex; justify-content: space-between; font-size: 0.65rem; color: #64748b; font-weight: 500; border-bottom: 1px dashed #f1f5f9; padding-bottom: 0.2rem; }
.ms-row span:last-child { font-weight: 700; color: var(--storm-ink); }
.text-green { color: #10b981 !important; }
.text-purple { color: #8b5cf6 !important; }

.text-xs { font-size: 0.65rem; color: #94a3b8; line-height: 1.3; text-align: center; }
.mt-auto { margin-top: auto; }

/* Status Summary Bottom Bar */
.status-summary {
  background: white; border: 1px solid #e2e8f0; border-radius: 12px;
  display: flex; align-items: center; justify-content: space-around; padding: 0.8rem 1.2rem;
  margin-top: 0.5rem;
}
.summary-item { display: flex; align-items: center; gap: 0.6rem; }
.sum-icon { width: 14px; }
.summary-item span { font-size: 0.72rem; color: #475569; line-height: 1.3; }
.summary-item strong { color: var(--storm-ink); }
</style>
