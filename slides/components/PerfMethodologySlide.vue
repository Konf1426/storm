<script setup lang="ts">
// Vue component for Performance Methodology
</script>

<template>
  <div class="page-shell">
    <div class="header-minimal">
      <p class="eyebrow">08 — Performance et résilience</p>
      <h2 class="section-title">Méthodologie : charger avec k6, profiler avec pprof, puis analyser</h2>
    </div>

    <div class="perf-grid">
      <!-- k6 -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box purple-bg"><img src="https://api.iconify.design/simple-icons/k6.svg?color=%237d64ff" class="icon-sm color-keep"/></div>
          <h3 class="card-title">Injection de charge (Grafana k6)</h3>
        </header>
        <div class="card-body">
          <p class="text-sm">Génération de trafic combiné HTTP + WebSocket avec des paliers progressifs :</p>
          <div class="k6-timeline">
             <div class="step-box border-cyan">
                <span class="step-time">2 min</span>
                <span class="step-name text-cyan">Warm-up</span>
                <span class="step-vu">50 VUs</span>
             </div>
             <div class="arr">➔</div>
             <div class="step-box border-slate">
                <span class="step-time">3 min</span>
                <span class="step-name text-slate">Spike 1</span>
                <span class="step-vu">150 VUs</span>
             </div>
             <div class="arr">➔</div>
             <div class="step-box border-orange bg-orange-light">
                <span class="step-time">3 min</span>
                <span class="step-name text-orange">Spike 2 + Chaos</span>
                <span class="step-vu">200 VUs</span>
             </div>
          </div>
          <div class="k6-visual-graph">
             <div class="graph-bar" style="height: 20%"></div>
             <div class="graph-bar" style="height: 40%"></div>
             <div class="graph-bar" style="height: 80%"></div>
             <div class="graph-bar" style="height: 60%"></div>
             <div class="graph-bar" style="height: 100%"></div>
             <div class="graph-label">Load Curve (VUs)</div>
          </div>
          <p class="def-text mt-2 text-center text-orange">
            ⚠️ <strong>Chaos Engineering :</strong> Arrêt inopiné de nœuds NATS/Postgres pendant le Spike final.
          </p>
        </div>
      </article>

      <!-- pprof -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box cyan-white-bg"><img src="https://api.iconify.design/simple-icons/go.svg?color=%2300add8" class="icon-sm color-keep lg-icon"/></div>
          <h3 class="card-title">Diagnostic & Profilage (pprof)</h3>
        </header>
        <div class="card-body">
          <p class="text-sm mb-3">Captures CPU et analyse de la mémoire (Heap) en pleine condition de charge pour repérer les <i>Hotspots</i> (goulets) :</p>
          
          <ul class="hotspot-list">
            <li>
               <img src="https://api.iconify.design/carbon/analytics.svg?color=%230ea5e9" class="hs-icon"/>
               <div class="hs-text"><strong>Parsing JWT & Crypto</strong> : Validation des signatures HMAC très coûteuses en CPU.</div>
            </li>
            <li>
               <img src="https://api.iconify.design/carbon/network-4.svg?color=%230ea5e9" class="hs-icon"/>
               <div class="hs-text"><strong>Syscalls Réseau</strong> : Saturation lors du broadcast NATS vers des milliers de sockets.</div>
            </li>
            <li>
               <img src="https://api.iconify.design/carbon/json.svg?color=%230ea5e9" class="hs-icon"/>
               <div class="hs-text"><strong>Sérialisation JSON</strong> : Coût d'allocation mémoire via <code>encoding/json</code> standard.</div>
            </li>
          </ul>
        </div>
      </article>

      <!-- Disclaimer -->
      <div class="disclaimer-bar">
         <div class="icon-box slate-bg flex-shrink-0"><img src="https://api.iconify.design/carbon/laptop.svg" class="icon-sm"/></div>
         <div class="disc-text">
            <strong>Environnement de test local (Docker Compose) :</strong> 
            Les métriques qui suivent (diapo suivante) proviennent toutes de stress-tests isolés sur un environnement local unifié. Elles valident la résilience structurelle de notre architecture et l'absence de fuites, sans prétendre aux performances brutes d'un cluster Cloud horizontal.
         </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Base Layout shared across new slides */
.page-shell { display: flex; flex-direction: column; height: 100%; }
.header-minimal { margin-bottom: 0.8rem; }

.perf-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr auto;
  gap: 1.2rem;
  flex-grow: 1;
  min-height: 0;
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

.card-head { display: flex; align-items: center; gap: 0.6rem; margin-bottom: 0.8rem; }
.card-title { font-size: 0.95rem; font-weight: 700; color: var(--storm-ink); margin: 0; }

.icon-box { width: 26px; height: 26px; border-radius: 6px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.flex-shrink-0 { flex-shrink: 0; }
.icon-sm { width: 14px; filter: grayscale(1) brightness(0) invert(1); }
.icon-sm.lg-icon { width: 20px; }
.icon-sm.color-keep { filter: none; }

.purple-bg { background: #f3e8ff; border: 1px solid #d8b4fe; } /* k6 is purple */
.cyan-white-bg { background: #e0f2fe; border: 1px solid #bae6fd; } /* Gopher is light blue */
.slate-bg { background: var(--storm-slate); }

.card-body { display: flex; flex-direction: column; flex-grow: 1; }

.text-sm { font-size: 0.72rem; line-height: 1.35; color: #475569; margin: 0; }
.mb-3 { margin-bottom: 0.8rem; }

/* k6 Timeline */
.k6-timeline {
  display: flex; align-items: center; justify-content: space-between;
  margin-top: auto; margin-bottom: auto; /* center vertically */
  padding: 0.8rem 0;
}
.step-box {
  background: #f8fafc; padding: 0.5rem 0.6rem; border-radius: 6px; text-align: center;
  border: 1px solid #e2e8f0; display: flex; flex-direction: column; gap: 0.2rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.01);
  width: 30%;
}
.border-cyan { border-top: 3px solid var(--storm-cyan); }
.border-slate { border-top: 3px solid #64748b; }
.border-orange { border-top: 3px solid var(--storm-orange); }
.bg-orange-light { background: #fffbeb; border-color: #fde68a; }

.step-time { font-family: monospace; font-size: 0.6rem; color: #94a3b8; font-weight: 600; }
.step-name { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; }
.step-vu { font-size: 0.65rem; color: var(--storm-ink); font-weight: 600; }

.text-cyan { color: var(--storm-cyan); }
.text-slate { color: #475569; }
.text-orange { color: var(--storm-orange); }

.arr { font-size: 0.8rem; color: #cbd5e1; }

/* Visual Graph */
.k6-visual-graph {
  display: flex; align-items: flex-end; gap: 4px; height: 40px; margin: 1rem 0;
  padding: 4px; border-bottom: 1px solid #e2e8f0; position: relative;
}
.graph-bar { flex: 1; background: var(--storm-cyan); border-radius: 2px 2px 0 0; opacity: 0.3; }
.graph-bar:last-of-type { background: var(--storm-orange); opacity: 0.6; }
.graph-label { position: absolute; right: 0; top: -10px; font-size: 0.5rem; color: #94a3b8; font-weight: 700; text-transform: uppercase; }

/* Info Text */
.def-text { font-size: 0.6rem; color: #d97706; font-weight: 500; }
.text-center { text-align: center; }

/* Hotspot List */
.hotspot-list {
  list-style: none; padding: 0; margin: 0;
  display: flex; flex-direction: column; gap: 0.5rem; justify-content: center; flex-grow:1;
}
.hotspot-list li {
  display: flex; align-items: flex-start; gap: 0.6rem;
  background: #f8fafc; border: 1px solid #f1f5f9; padding: 0.6rem; border-radius: 6px;
}
.hs-text { font-size: 0.65rem; color: #475569; line-height: 1.3; }
.hs-text strong { color: var(--storm-ink); }
.hs-icon { width: 14px; margin-top: 2px; }

/* Disclaimer Bar */
.disclaimer-bar {
  grid-column: 1 / -1;
  background: white; border: 1px solid #e2e8f0; border-radius: 8px;
  display: flex; align-items: center; gap: 1rem; padding: 0.8rem 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
}
.disc-text {
  font-size: 0.65rem; color: #64748b; line-height: 1.4;
}
.disc-text strong { color: var(--storm-ink); }
</style>
