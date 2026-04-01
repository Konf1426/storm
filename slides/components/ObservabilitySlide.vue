<script setup lang="ts">
// Observability & SLOs detailed slide
</script>

<template>
  <div class="page-shell">
    <div class="header-minimal">
      <p class="eyebrow">07 — Performance et résilience</p>
      <h2 class="section-title">Observabilité et SLOs pour piloter la résilience</h2>
    </div>

    <div class="obs-grid">
      <!-- SLOs -->
      <article class="clean-card shadow-sm slo-card">
        <header class="card-head">
          <div class="icon-box orange-bg"><img src="https://api.iconify.design/carbon/target.svg" class="icon-sm"/></div>
          <h3 class="card-title">Objectifs de Niveau de Service (SLOs)</h3>
        </header>
        <div class="card-body">
          <p class="text-sm mb-3">
            <strong>SLO</strong> <i>(Service Level Objective)</i> : Promesses mesurables de fiabilité du système.
          </p>
          <ul class="slo-list">
             <li>
               <div class="slo-stat text-green">≥ 99.5%</div>
               <div class="slo-desc">Disponibilité globale (Uptime) requise</div>
             </li>
             <li>
               <div class="slo-stat text-orange">< 200<span class="unit">ms</span></div>
               <div class="slo-desc">Latence de publication des messages<br/><span class="def-text"><strong>p95</strong> = 95% des requêtes sont plus rapides que ce temps</span></div>
             </li>
             <li>
               <div class="slo-stat text-orange">< 200<span class="unit">ms</span></div>
               <div class="slo-desc">Temps d'établissement WebSocket<br/><span class="def-text"><strong>p95</strong> = Garantie de connexion quasi-instantanée</span></div>
             </li>
             <li>
               <div class="slo-stat text-cyan">< 1%</div>
               <div class="slo-desc">Taux d'erreur HTTP global acceptable</div>
             </li>
          </ul>
        </div>
      </article>

      <div class="right-col">
        <!-- Stack -->
        <article class="clean-card shadow-sm">
          <header class="card-head">
            <div class="icon-box cyan-bg"><img src="https://api.iconify.design/carbon/chart-line-smooth.svg" class="icon-sm"/></div>
            <h3 class="card-title">Stack Télémétrique</h3>
          </header>
          <div class="card-body gap-snug">
            <div class="stack-item">
              <img src="https://api.iconify.design/logos/prometheus.svg" class="logo-xs"/> 
              <span><strong>Prometheus</strong> : Scrape des endpoints <code>/metrics</code> (Pull)</span>
            </div>
            <div class="stack-item">
              <img src="https://api.iconify.design/logos/grafana.svg" class="logo-xs"/> 
              <span><strong>Grafana</strong> : Visualisation temps réel & alertes SLOs</span>
            </div>
            <div class="stack-item">
              <img src="https://api.iconify.design/carbon/activity.svg?color=%2364748b" class="logo-xs"/> 
              <span><strong>Probes Liveness</strong> : <code>/healthz</code> et <code>/ping-nats</code><br/><span class="def-text">Sondes vérifiant en continu que le serveur n'est pas "planté"</span></span>
            </div>
          </div>
        </article>

        <!-- Metrics -->
        <article class="clean-card shadow-sm dark-card">
          <header class="card-head">
            <div class="icon-box slate-bg"><img src="https://api.iconify.design/carbon/code.svg" class="icon-sm"/></div>
            <h3 class="card-title text-white">Métriques Métier Go Exposées</h3>
          </header>
          <div class="card-body">
             <div class="code-metrics">
               <span class="metric-line"><span class="m-type">gauge</span> storm_active_websockets</span>
               <span class="metric-line"><span class="m-type">counter</span> storm_auth_rate_limited_total</span>
               <span class="metric-line"><span class="m-type">histogram</span> storm_nats_publish_duration_seconds</span>
               <span class="metric-line"><span class="m-type">gauge</span> storm_save_queue_length</span>
             </div>
          </div>
        </article>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Base Layout shared across new slides */
.page-shell { display: flex; flex-direction: column; height: 100%; }
.header-minimal { margin-bottom: 0.8rem; }

.obs-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.2rem;
  flex-grow: 1;
  min-height: 0;
}

.right-col {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
}

/* Cards */
.clean-card {
  background: white;
  border-radius: 12px;
  padding: 1.2rem 1.4rem;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0,0,0,0.04);
}

.dark-card {
  background: #0f172a;
  border-color: #1e293b;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.2);
}

.card-head { display: flex; align-items: center; gap: 0.6rem; margin-bottom: 0.6rem; }
.card-title { font-size: 1rem; font-weight: 700; color: var(--storm-ink); margin: 0; }
.text-white { color: white; }

.icon-box { width: 26px; height: 26px; border-radius: 6px; display: flex; align-items: center; justify-content: center; }
.icon-sm { width: 14px; filter: grayscale(1) brightness(0) invert(1); }

.orange-bg { background: #f97316; box-shadow: 0 4px 12px rgba(249, 115, 22, 0.2); }
.cyan-bg { background: #06b6d4; box-shadow: 0 4px 12px rgba(6, 182, 212, 0.2); }
.slate-bg { background: #334155; box-shadow: 0 4px 12px rgba(51, 65, 85, 0.3); }

.card-body { display: flex; flex-direction: column; flex-grow: 1; }
.gap-snug { gap: 0.6rem; justify-content: center; }

.text-sm { font-size: 0.72rem; line-height: 1.35; color: #475569; margin: 0; }
.mb-3 { margin-bottom: 0.8rem; }

/* SLO List */
.slo-card {
  background: linear-gradient(180deg, #ffffff 0%, #f8fafc 100%);
}
.slo-list {
  list-style: none; padding: 0; margin: 0;
  display: flex; flex-direction: column; gap: 0.6rem; flex-grow: 1; justify-content: center;
}
.slo-list li {
  display: flex; align-items: center; gap: 1rem;
  background: white; padding: 0.6rem 1rem; border-radius: 8px; border: 1px solid #e2e8f0;
  box-shadow: 0 2px 4px rgba(0,0,0,0.01);
}
.slo-stat {
  font-size: 1.3rem; font-weight: 800; font-family: monospace; letter-spacing: -0.03em;
  width: 90px; text-align: right;
}
.unit { font-size: 0.8rem; font-weight: 600; opacity: 0.8; }
.slo-desc { font-size: 0.75rem; font-weight: 600; color: #334155; line-height: 1.2; }

.text-green { color: #10b981; }
.text-orange { color: #f59e0b; }
.text-cyan { color: #0ea5e9; }

.def-text { display: block; font-size: 0.6rem; color: #94a3b8; font-weight: 500; margin-top: 0.2rem; }

/* Stack Items */
.stack-item {
  display: flex; align-items: center; gap: 0.6rem;
  background: #f8fafc; padding: 0.5rem 0.8rem; border-radius: 6px; border: 1px solid #f1f5f9;
}
.stack-item span { font-size: 0.7rem; color: #334155; line-height: 1.3; }
.stack-item strong { color: var(--storm-ink); }
.stack-item code { font-size: 0.6rem; background: white; padding: 0.1rem 0.3rem; border-radius: 3px; border: 1px solid #e2e8f0; color: #0369a1; }
.logo-xs { width: 16px; height: 16px; flex-shrink: 0; }

/* Code Metrics Dark */
.code-metrics {
  background: #020617; padding: 0.8rem; border-radius: 6px;
  display: flex; flex-direction: column; gap: 0.4rem;
  font-family: "IBM Plex Mono", monospace; font-size: 0.62rem; color: #94a3b8;
  border: 1px solid #1e293b; flex-grow: 1; justify-content: center;
}
.metric-line { display: flex; gap: 0.6rem; align-items: center; }
.m-type { color: #8b5cf6; font-weight: 700; width: 50px; }
</style>
