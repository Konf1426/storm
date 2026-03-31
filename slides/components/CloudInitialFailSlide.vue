<script setup lang="ts">
// Cloud Step 1: Initial Failure - Strict Bento DA (like Slide 5)
</script>

<template>
  <div class="page-shell bento-shell">
    <div class="bento-header">
      <p class="eyebrow">11 — Azure Cloud : Étape 1</p>
      <h2 class="section-title">Le Test de Départ : 1 000 VUs</h2>
    </div>

    <div class="bento-grid grid-fail">
      
      <!-- Goulot d'étranglement (Carte large, gauche) -->
      <article class="bento-card bento-bottleneck">
        <div class="bento-card-header">
          <div class="bento-icon orange-glow"><img src="https://api.iconify.design/carbon/firewall.svg" style="width:16px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Effondrement & Rate Limiting</h3>
        </div>
        <p class="bento-desc">
          Le jour 1, le système n'est pas dimensionné pour le volume d'un testeur k6. Un garde-fou Azure bloque instantanément l'assaut HTTP.
        </p>
        
        <div class="rl-visual">
          <div class="rl-label-row">
            <span class="rl-small">Attaquant (k6)</span>
            <span class="rl-small">Gateway AGIC</span>
            <span class="rl-small">Cluster AKS</span>
          </div>
          
          <div class="rl-flow">
            <!-- Paquets -->
            <div class="rl-packets">
               <div class="rl-packet"></div>
               <div class="rl-packet"></div>
               <div class="rl-packet"></div>
               <div class="rl-packet"></div>
               <div class="rl-packet warning"></div>
               <div class="rl-packet danger"></div>
            </div>
            
            <!-- Mur -->
            <div class="rl-wall">
               <span>Rate Limit</span>
               <code class="rl-limit-code">5 req / min / IP</code>
            </div>

            <!-- AKS (Vide) -->
            <div class="rl-aks">
               <div class="aks-box">Node B2s</div>
            </div>
          </div>
        </div>

        <div class="bento-tag-row">
          <span class="bento-tag-label">Infra de base</span>
          <span class="bento-tag-value">1 Nœud B2s, 2 Gateways, PG 4 vCores</span>
        </div>
      </article>

      <!-- Taux de Succès (Haut Droite) -->
      <article class="bento-card bento-success-rate">
        <div class="bento-card-header">
          <div class="bento-icon orange-glow"><img src="https://api.iconify.design/carbon/warning.svg" style="width:14px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Taux de Succès</h3>
        </div>
        <div class="kpi-visual">
          <div class="kpi-value c-orange">0,01<span class="kpi-unit">%</span></div>
          <span class="kpi-subc">Requêtes massivement rejetées (429)</span>
        </div>
      </article>

      <!-- Latences & Coûts (Bas Droite) -->
      <article class="bento-card bento-fail-metrics">
        <div class="bento-card-header">
          <div class="bento-icon slate-glow"><img src="https://api.iconify.design/carbon/timer.svg" style="width:14px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Latences Engorgées</h3>
        </div>
        <div class="metrics-visual">
          <div class="mv-row">
            <span class="mv-label">Login API</span>
            <span class="mv-val">> 4.0 s</span>
          </div>
          <div class="mv-row border-b">
            <span class="mv-label">Msg Exchange</span>
            <span class="mv-val">~2.5 s</span>
          </div>
          <div class="mv-row mv-budget mt-1">
            <span class="mv-label c-cyan">Budget / H</span>
            <span class="mv-val code-font">0,15 $</span>
          </div>
        </div>
      </article>

    </div>
  </div>
</template>

<style scoped>
.bento-shell { display: flex; flex-direction: column; height: 100%; gap: 0.5rem; }
.bento-header { display: grid; gap: 0.1rem; margin-bottom: 0; }

.bento-grid.grid-fail {
  display: grid; flex-grow: 1;
  grid-template-columns: 1.3fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 0.8rem; min-height: 0;
}

.bento-bottleneck { grid-column: 1; grid-row: 1 / 3; }
.bento-success-rate { grid-column: 2; grid-row: 1; }
.bento-fail-metrics { grid-column: 2; grid-row: 2; }

/* DA Standard Slide 5 for Cards */
.bento-card {
  background: var(--storm-surface-strong);
  border: 1px solid rgba(22, 34, 43, 0.06);
  border-radius: 16px; padding: 0.8rem 1rem;
  display: flex; flex-direction: column; gap: 0.4rem;
  box-shadow: 0 4px 16px rgba(18, 28, 36, 0.04);
  backdrop-filter: blur(20px); min-height: 0;
}

.bento-card-header { display: flex; align-items: center; gap: 0.4rem; }
.bento-card-header h3 { font-size: 0.9rem; font-weight: 700; color: var(--storm-ink); letter-spacing: -0.02em; margin: 0; }
.bento-icon { width: 20px; height: 20px; border-radius: 6px; display: flex; align-items: center; justify-content: center; background: rgba(22, 34, 43, 0.04); }
.cyan-glow { box-shadow: 0 0 12px rgba(13, 159, 184, 0.2); color: var(--storm-cyan); }
.orange-glow { box-shadow: 0 0 12px rgba(239, 138, 41, 0.2); color: var(--storm-orange); }
.slate-glow { box-shadow: 0 0 12px rgba(63, 85, 96, 0.2); color: var(--storm-slate); }

.bento-desc { font-size: 0.72rem; line-height: 1.25; color: var(--storm-muted); margin: 0; }

/* -- Rate Limit Visual -- */
.rl-visual {
  flex-grow: 1; display: flex; flex-direction: column; justify-content: center; gap: 1rem;
  background: rgba(22, 34, 43, 0.02); border-radius: 8px; padding: 1.2rem;
  border: 1px solid rgba(22, 34, 43, 0.05); font-family: "IBM Plex Mono", monospace;
  margin-top: 0.5rem; margin-bottom: 0.5rem;
}
.rl-label-row { display: flex; justify-content: space-between; align-items: center; padding: 0 0.5rem; }
.rl-small { font-size: 0.55rem; font-weight: 700; color: var(--storm-slate); text-transform: uppercase; }

.rl-flow { display: flex; align-items: center; justify-content: space-between; height: 80px; position: relative; }
.rl-packets { display: flex; gap: 4px; align-items: center; width: 35%; flex-wrap: wrap; padding-left: 0.5rem; }
.rl-packet { width: 18px; height: 6px; background: var(--storm-cyan); border-radius: 2px; }
.rl-packet.warning { background: var(--storm-orange); opacity: 0.8; height: 8px; }
.rl-packet.danger { background: #ef4444; opacity: 0.9; height: 10px; }

.rl-wall {
  display: flex; flex-direction: column; justify-content: center; align-items: center;
  background: var(--storm-orange); color: white; border-radius: 6px; padding: 0.4rem 0.6rem;
  height: 110%; width: 28%; font-family: var(--storm-font); text-align: center;
  box-shadow: 0 0 15px rgba(239, 138, 41, 0.3); z-index: 2; border: 1px solid #c2410c;
}
.rl-wall span { font-size: 0.6rem; font-weight: 800; text-transform: uppercase; }
.rl-limit-code { font-size: 0.5rem; background: rgba(0,0,0,0.2); padding: 0.15rem; border-radius: 3px; margin-top: 4px; font-family: "IBM Plex Mono"; }

.rl-aks { width: 30%; display: flex; justify-content: center; opacity: 0.3; }
.aks-box { border: 1px dashed var(--storm-slate); padding: 0.5rem 0.8rem; border-radius: 4px; font-size: 0.65rem; color: var(--storm-slate); font-weight: 600; }

.bento-tag-row { margin-top: 0.3rem; display: flex; align-items: center; justify-content: space-between; background: rgba(63, 85, 96, 0.08); border: 1px solid rgba(63, 85, 96, 0.2); border-radius: 6px; padding: 0.3rem 0.5rem; }
.bento-tag-label { font-size: 0.65rem; font-weight: 600; color: var(--storm-ink); }
.bento-tag-value { font-family: "IBM Plex Mono", monospace; font-size: 0.55rem; color: var(--storm-slate); font-weight: 700;}

/* -- KPI Visual -- */
.kpi-visual { margin-top: auto; display: flex; flex-direction: column; justify-content: center; align-items: center; background: #11181c; border-radius: 8px; padding: 1rem; border: 1px solid rgba(239, 138, 41, 0.3); box-shadow: inset 0 2px 8px rgba(0,0,0,0.2); }
.kpi-value { font-size: 3rem; font-weight: 800; line-height: 1; letter-spacing: -0.04em; }
.kpi-unit { font-size: 1.5rem; margin-left: 2px; }
.c-orange { color: var(--storm-orange); }
.kpi-subc { font-family: "IBM Plex Mono", monospace; font-size: 0.55rem; color: #94a3b8; margin-top: 0.3rem; text-transform: uppercase; letter-spacing: 0.05em; text-align: center; }

/* -- Metrics Visual -- */
.metrics-visual { margin-top: auto; display: flex; flex-direction: column; gap: 0.4rem; background: rgba(22, 34, 43, 0.02); border-radius: 8px; padding: 0.6rem; border: 1px solid rgba(22, 34, 43, 0.05); }
.mv-row { display: flex; justify-content: space-between; align-items: center; padding: 0.2rem 0; }
.mv-row.border-b { border-bottom: 1px dashed rgba(22, 34, 43, 0.1); padding-bottom: 0.4rem; }
.mv-label { font-size: 0.65rem; font-weight: 600; color: var(--storm-slate); text-transform: uppercase; letter-spacing: 0.02em; }
.mv-val { font-size: 1.2rem; font-weight: 800; color: var(--storm-ink); font-family: "IBM Plex Mono", monospace; }
.c-cyan { color: var(--storm-cyan); }
.code-font { color: var(--storm-cyan); font-size: 0.9rem; }
.mt-1 { margin-top: 0.2rem; }
.mv-budget { background: #11181c; padding: 0.4rem; border-radius: 6px; }
.mv-budget .mv-val { color: var(--storm-cyan); }
</style>
