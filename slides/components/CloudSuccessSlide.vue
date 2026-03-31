<script setup lang="ts">
// Cloud Step 3: Success & Final Table - Strict Bento DA (like Slide 5)
</script>

<template>
  <div class="page-shell bento-shell">
    <div class="bento-header">
      <p class="eyebrow">13 — Azure Cloud : Étape 3</p>
      <h2 class="section-title">Le Test Ultime : 10 000 VUs</h2>
    </div>

    <div class="bento-grid grid-success">
      
      <!-- KPI 1 : Succès -->
      <article class="bento-card bento-succ">
         <div class="bento-card-header">
           <div class="bento-icon green-glow"><img src="https://api.iconify.design/carbon/checkmark-filled.svg" style="width:14px;filter:grayscale(1) brightness(0);"/></div>
           <h3>Taux de Succès</h3>
         </div>
         <div class="kpi-visual border-green">
            <div class="kpi-value">
               <span class="c-green font-block">~100</span><span class="kpi-unit c-green-dim">%</span>
            </div>
            <span class="kpi-sub">FIABLE</span>
         </div>
      </article>

      <!-- KPI 2 : Latence -->
      <article class="bento-card bento-lat">
         <div class="bento-card-header">
           <div class="bento-icon cyan-glow"><img src="https://api.iconify.design/carbon/flash.svg" style="width:14px;filter:grayscale(1) brightness(0);"/></div>
           <h3>Latence Messages</h3>
         </div>
         <div class="kpi-visual border-cyan">
            <div class="kpi-value">
               <span class="c-cyan font-block">85</span><span class="kpi-unit c-cyan-dim">ms</span>
            </div>
            <span class="kpi-sub">MÉDIANE</span>
         </div>
      </article>

      <!-- Infra -->
      <article class="bento-card bento-infra">
         <div class="bento-card-header">
           <div class="bento-icon slate-glow"><img src="https://api.iconify.design/carbon/cluster.svg" style="width:14px;filter:grayscale(1) brightness(0);"/></div>
           <h3>Artillerie Déployée</h3>
         </div>
         <div class="infra-visual">
            <div class="iv-col border-r">
               <span class="iv-label">Cluster AKS</span>
               <span class="iv-spec">26 vCPUs (D8s)</span>
            </div>
            <div class="iv-col pl">
               <span class="iv-label">PostgreSQL</span>
               <span class="iv-spec">16 vCores (D16ds)</span>
            </div>
         </div>
         <div class="cookie-tag mt-auto">
            <span class="tag-label">Volume encaissé</span>
            <span class="tag-value text-slate">210M msgs WebSockets</span>
         </div>
      </article>

      <!-- Tableau Comparatif encapsulé dans une bento-card -->
      <article class="bento-card bento-table">
        <div class="bento-card-header mb-1">
          <div class="bento-icon slate-glow"><img src="https://api.iconify.design/carbon/comparison.svg" style="width:14px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Synthèse des Trajectoires Azure (Perf vs Budget)</h3>
        </div>
        
        <div class="storm-table-container">
           <table class="storm-table">
             <thead>
               <tr>
                 <th>Tier</th>
                 <th>Configuration Principale</th>
                 <th class="text-right">Coût / H</th>
                 <th class="text-right">Perf. Login</th>
               </tr>
             </thead>
             <tbody>
               <!-- Initial -->
               <tr>
                 <td><span class="bento-badge bg-dark">Initial</span></td>
                 <td class="font-medium text-muted">B2s (2 vCPU) + DB 4 Cores</td>
                 <td class="text-right font-mono">0,15 $</td>
                 <td class="text-right font-black c-orange">> 5 000 ms</td>
               </tr>
               <!-- Mid -->
               <tr>
                 <td><span class="bento-badge border-cyan c-cyan">Mid</span></td>
                 <td class="font-medium text-muted">5x B2s (10 vCPU) + DB 4 Cores</td>
                 <td class="text-right font-mono">~0,60 $</td>
                 <td class="text-right font-black text-ink">1 800 ms</td>
               </tr>
               <!-- Ultra -->
               <tr class="ultra-row">
                 <td><span class="bento-badge border-green c-green bg-green-subtle">Ultra</span></td>
                 <td class="font-bold text-ink">26 vCPUs (D-Series) + DB 16 Cores</td>
                 <td class="text-right font-mono font-black c-green">3,33 $</td>
                 <td class="text-right font-black c-green">221 ms</td>
               </tr>
               <!-- Veille -->
               <tr class="veille-row">
                 <td><span class="bento-badge box-bg text-muted">Veille</span></td>
                 <td class="font-medium text-muted italic">1x B2s + DB 2 Cores (Actuel)</td>
                 <td class="text-right font-mono text-muted">0,15 $</td>
                 <td class="text-right font-bold text-muted uppercase text-xxs">— (Éco)</td>
               </tr>
             </tbody>
           </table>
        </div>
      </article>

    </div>
  </div>
</template>

<style scoped>
.bento-shell { display: flex; flex-direction: column; height: 100%; gap: 0.5rem; }
.bento-header { display: grid; gap: 0.1rem; margin-bottom: 0; }

.bento-grid.grid-success {
  display: grid; flex-grow: 1;
  grid-template-columns: 1fr 1fr 1.6fr;
  grid-template-rows: auto 1fr;
  gap: 0.8rem; min-height: 0;
}

.bento-succ { grid-column: 1; grid-row: 1; }
.bento-lat { grid-column: 2; grid-row: 1; }
.bento-infra { grid-column: 3; grid-row: 1; }
.bento-table { grid-column: 1 / 4; grid-row: 2; padding: 0.8rem 1.2rem; }

/* DA Standard Slide 5 for Cards */
.bento-card {
  background: var(--storm-surface-strong);
  border: 1px solid rgba(22, 34, 43, 0.06);
  border-radius: 16px; padding: 0.6rem 0.8rem;
  display: flex; flex-direction: column; gap: 0.4rem;
  box-shadow: 0 4px 16px rgba(18, 28, 36, 0.04);
  backdrop-filter: blur(20px); min-height: 0;
}
.bento-card-header { display: flex; align-items: center; gap: 0.4rem; }
.bento-card-header h3 { font-size: 0.8rem; font-weight: 700; color: var(--storm-ink); letter-spacing: -0.02em; margin: 0; }
.mb-1 { margin-bottom: 0.3rem; }
.bento-icon { width: 18px; height: 18px; border-radius: 6px; display: flex; align-items: center; justify-content: center; background: rgba(22, 34, 43, 0.04); }
.cyan-glow { box-shadow: 0 0 10px rgba(13, 159, 184, 0.15); color: var(--storm-cyan); }
.orange-glow { box-shadow: 0 0 10px rgba(239, 138, 41, 0.15); color: var(--storm-orange); }
.slate-glow { box-shadow: 0 0 10px rgba(63, 85, 96, 0.15); color: var(--storm-slate); }
.green-glow { box-shadow: 0 0 10px rgba(16, 185, 129, 0.15); color: #10b981; }

.c-cyan { color: var(--storm-cyan); }
.c-cyan-dim { color: rgba(13, 159, 184, 0.6); }
.c-orange { color: var(--storm-orange); }
.c-green { color: #10b981; }
.c-green-dim { color: rgba(16, 185, 129, 0.6); }
.border-cyan { border: 1px solid rgba(13, 159, 184, 0.3); }
.border-green { border: 1px solid rgba(16, 185, 129, 0.3); }
/* -- KPI Visuals (Light Tinted Style) -- */
.kpi-visual { 
  flex-grow: 1; display: flex; flex-direction: column; justify-content: center; align-items: center; 
  background: rgba(22, 34, 43, 0.02); border-radius: 8px; padding: 0.5rem; 
  border: 1px solid rgba(22, 34, 43, 0.05); margin-top: auto; 
}
.kpi-visual.border-green { background: rgba(16, 185, 129, 0.05); border: 1px solid rgba(16, 185, 129, 0.2); }
.kpi-visual.border-cyan { background: rgba(13, 159, 184, 0.05); border: 1px solid rgba(13, 159, 184, 0.2); }

.kpi-value { display: flex; align-items: baseline; justify-content: center; line-height: 1; }
.font-block { font-family: "IBM Plex Mono", monospace; font-size: 1.8rem; font-weight: 800; letter-spacing: -0.05em; }
.kpi-unit { font-size: 1rem; font-weight: 800; margin-left: 2px; }
.text-slate { color: var(--storm-slate); }
.kpi-sub { font-family: "IBM Plex Mono", monospace; font-size: 0.45rem; font-weight: 700; color: var(--storm-slate); text-transform: uppercase; letter-spacing: 0.1em; margin-top: 0.1rem; opacity: 0.8; }

/* Infra Visual */
.infra-visual { display: flex; align-items: center; background: rgba(22, 34, 43, 0.02); border-radius: 8px; padding: 0.6rem; border: 1px solid rgba(22, 34, 43, 0.05); margin-top: 0.2rem; }
.iv-col { flex: 1; display: flex; flex-direction: column; gap: 0.2rem; }
.border-r { border-right: 1px dashed rgba(22, 34, 43, 0.1); }
.pl { padding-left: 0.8rem; }
.iv-label { font-size: 0.55rem; font-weight: 800; color: var(--storm-slate); text-transform: uppercase; letter-spacing: 0.02em; }
.iv-spec { font-family: "IBM Plex Mono", monospace; font-size: 0.7rem; font-weight: 700; color: var(--storm-ink); }

.cookie-tag { display: flex; align-items: center; justify-content: space-between; background: rgba(63, 85, 96, 0.08); border: 1px solid rgba(63, 85, 96, 0.2); border-radius: 6px; padding: 0.3rem 0.6rem; }
.tag-label { font-size: 0.6rem; font-weight: 600; color: var(--storm-ink); }
.tag-value { font-family: "IBM Plex Mono", monospace; font-size: 0.55rem; font-weight: 700;}
.mt-auto { margin-top: auto; }

/* Custom Storm Table inside Bento Card */
.storm-table-container { background: #ffffff; border-radius: 10px; border: 1px solid rgba(22, 34, 43, 0.08); overflow: hidden; height: 100%; display: flex; flex-direction: column; }
.bg-dark { background: rgba(22, 34, 43, 0.05); color: var(--storm-slate) !important; border: 1px solid rgba(22, 34, 43, 0.1); }
.box-bg { background: white; border: 1px solid rgba(22, 34, 43, 0.1); color: var(--storm-slate); }
.storm-table { width: 100%; border-collapse: collapse; text-align: left; }
.storm-table th { padding: 0.4rem 0.6rem; font-size: 0.5rem; font-weight: 700; color: var(--storm-slate); text-transform: uppercase; letter-spacing: 0.05em; background: rgba(22, 34, 43, 0.03); border-bottom: 1px solid rgba(22, 34, 43, 0.08); }
.storm-table td { padding: 0.35rem 0.6rem; border-bottom: 1px solid rgba(22, 34, 43, 0.04); vertical-align: middle; font-size: 0.65rem; }
.text-right { text-align: right; }
.font-medium { font-weight: 500; }
.font-bold { font-weight: 700; }
.font-black { font-weight: 800; font-size: 0.75rem; }
.font-mono { font-family: "IBM Plex Mono", monospace; font-size: 0.6rem; font-weight: 600; color: var(--storm-slate); }
.text-muted { color: var(--storm-muted); }
.text-ink { color: var(--storm-ink); }
.opacity-80 { opacity: 0.8; }
.italic { font-style: italic; }
.text-xxs { font-size: 0.5rem; letter-spacing: 0.05em; }

.ultra-row { background: rgba(16, 185, 129, 0.04); border-left: 3px solid #10b981; }
.veille-row { background: transparent; }

.bento-badge { font-family: var(--storm-font); font-size: 0.5rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.05em; padding: 0.15rem 0.3rem; border-radius: 4px; box-shadow: 0 1px 2px rgba(0,0,0,0.02); }
</style>
