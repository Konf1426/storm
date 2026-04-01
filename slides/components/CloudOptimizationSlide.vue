<script setup lang="ts">
// Cloud Step 2: Optimization Pivot - Strict Bento DA (like Slide 5)
</script>

<template>
  <div class="page-shell bento-shell">
    <div class="bento-header">
      <p class="eyebrow">12 — Scalabilité cloud, budget et limites</p>
      <h2 class="section-title">Étape 2 : 5 000 VUs et pivot d'architecture</h2>
    </div>

    <div class="bento-grid grid-opti">
      
      <!-- Le Pivot NATS (Carte large en haut) -->
      <article class="bento-card bento-rt-pivot">
        <div class="bento-card-header">
          <div class="bento-icon cyan-glow"><img src="https://api.iconify.design/carbon/rotate.svg" style="width:14px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Le Pivot : Flux Asynchrone</h3>
        </div>
        <div class="rt-content">
          <p class="bento-desc">
            Pour tenir la charge, on élargit l'infrastructure mais surtout on <strong>découple</strong> la base de données. La Gateway ne bloque plus, elle confie la validation à <span class="c-cyan font-bold">NATS</span>.
          </p>
          <div class="rt-visual">
            <div class="rt-user">Gateway</div>
            <div class="rt-arrow">
               <span class="rt-pub c-orange">Publie msg</span>
               <div class="rt-line"><div class="rt-dot"></div></div>
            </div>
            <div class="rt-nats">
               Broker NATS
               <code>Asynchrone</code>
            </div>
            <div class="rt-arrow">
               <span class="rt-pub c-cyan">Reçoit</span>
               <div class="rt-line dashed"></div>
            </div>
            <div class="rt-user">Worker PG</div>
          </div>
        </div>
      </article>

      <!-- Optimisation CPU (Bas Gauche) -->
      <article class="bento-card bento-cpu">
        <div class="bento-card-header">
          <div class="bento-icon slate-glow"><img src="https://api.iconify.design/carbon/microchip.svg" style="width:14px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Optimisation CPU (Bcrypt)</h3>
        </div>
        <p class="bento-desc">
          Les nœuds (B2s) saturaient lors du hashage massif des mots de passe.
        </p>
        
        <div class="cpu-visual">
          <div class="cpu-step bg-dark">
            <span class="wv-label c-orange">Cost = 10</span>
            <span class="wv-detail">CPU Surchargé</span>
          </div>
          <div class="wv-arrow">➔</div>
          <div class="cpu-step box-bg">
            <span class="wv-label c-cyan">Cost = 4</span>
            <span class="wv-detail">Charge Tolérable</span>
          </div>
        </div>
      </article>

      <!-- Résultats Latence (Bas Droite) -->
      <article class="bento-card bento-opti-metrics">
        <div class="bento-card-header">
          <div class="bento-icon cyan-glow"><img src="https://api.iconify.design/carbon/chart-line-smooth.svg" style="width:14px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Bénéfices & Nouveau Limiteur</h3>
        </div>
        
        <div class="metrics-visual">
          <div class="mv-row">
            <span class="mv-label c-cyan">Latence MSG</span>
            <span class="mv-val code-font">&lt; 100 ms</span>
          </div>
          <div class="mv-row">
            <span class="mv-label">Budget Avoisiné</span>
            <span class="mv-val">~0,60 $ / H</span>
          </div>
        </div>

        <div class="cookie-tag mt-auto">
          <span class="tag-label">Limite Actuelle</span>
          <span class="tag-value text-slate">CPU AKS (10 cores)</span>
        </div>
      </article>

    </div>
  </div>
</template>

<style scoped>
.bento-shell { display: flex; flex-direction: column; height: 100%; gap: 0.5rem; }
.bento-header { display: grid; gap: 0.1rem; margin-bottom: 0; }

.bento-grid.grid-opti {
  display: grid; flex-grow: 1;
  grid-template-columns: 1fr 1.2fr;
  grid-template-rows: auto 1fr;
  gap: 0.8rem; min-height: 0;
}

.bento-rt-pivot { grid-column: 1 / 3; grid-row: 1; }
.bento-cpu { grid-column: 1; grid-row: 2; }
.bento-opti-metrics { grid-column: 2; grid-row: 2; }

/* DA Standard Slide 5 */
.bento-card {
  background: var(--storm-surface-strong);
  border: 1px solid rgba(22, 34, 43, 0.06);
  border-radius: 16px; padding: 0.8rem 1rem;
  display: flex; flex-direction: column; gap: 0.4rem;
  box-shadow: 0 4px 16px rgba(18, 28, 36, 0.04);
  backdrop-filter: blur(20px); min-height: 0;
}
.bento-card-header { display: flex; align-items: center; gap: 0.4rem; }
.bento-card-header h3 { font-size: 0.8rem; font-weight: 700; color: var(--storm-ink); letter-spacing: -0.02em; margin: 0; }
.bento-icon { width: 18px; height: 18px; border-radius: 6px; display: flex; align-items: center; justify-content: center; background: rgba(22, 34, 43, 0.04); }
.cyan-glow { box-shadow: 0 0 10px rgba(13, 159, 184, 0.15); color: var(--storm-cyan); }
.orange-glow { box-shadow: 0 0 10px rgba(239, 138, 41, 0.15); color: var(--storm-orange); }
.slate-glow { box-shadow: 0 0 10px rgba(63, 85, 96, 0.15); color: var(--storm-slate); }
.bento-desc { font-size: 0.65rem; line-height: 1.2; color: var(--storm-muted); margin: 0; }

/* -- Real Time Visual -- */
.rt-content { display: flex; flex-direction: column; justify-content: center; flex-grow: 1; gap: 1rem; padding: 0.5rem 0;}
.rt-visual { display: flex; align-items: center; background: rgba(22, 34, 43, 0.03); border-radius: 8px; padding: 0.8rem; font-family: "IBM Plex Mono", monospace; font-size: 0.65rem; border: 1px solid rgba(22, 34, 43, 0.05); }
.rt-user { background: white; border: 1px solid rgba(22, 34, 43, 0.15); padding: 0.4rem 0.6rem; border-radius: 4px; font-weight: 700; box-shadow: 0 2px 4px rgba(0,0,0,0.03); display: flex; align-items: center; gap: 0.2rem; color: var(--storm-ink); }
.rt-nats { background: #11181c; color: white !important; padding: 0.5rem 0.8rem; border-radius: 6px; text-align: center; border: 2px solid var(--storm-cyan); box-shadow: 0 0 12px rgba(13, 159, 184, 0.25); font-size: 0.6rem; font-weight: 700; }
.rt-nats code { color: var(--storm-cyan) !important; display: block; margin-top: 0.1rem; font-size: 0.5rem; }
.rt-arrow { display: flex; flex-direction: column; align-items: center; flex-grow: 1; padding: 0 0.5rem; }
.rt-pub { font-size: 0.5rem; font-weight: 800; text-transform: uppercase; margin-bottom: 4px; }
.rt-line { flex-grow: 1; width: 100%; height: 2px; background: rgba(22, 34, 43, 0.15); position: relative; margin: 0; }
.rt-line.dashed { background: transparent; border-top: 2px dashed rgba(22, 34, 43, 0.25); }
.rt-dot { position: absolute; top: -3px; left: 0; width: 8px; height: 8px; background: var(--storm-orange); border-radius: 50%; box-shadow: 0 0 10px var(--storm-orange); animation: pingPong 1.5s cubic-bezier(0.4, 0, 0.2, 1) infinite; }
@keyframes pingPong { 0% { transform: translateX(0); } 100% { transform: translateX(100%); } }
.c-cyan { color: var(--storm-cyan); font-weight: 700; }
.c-orange { color: var(--storm-orange); font-weight: 700; }

/* -- CPU Visual -- */
.cpu-visual { margin-top: auto; display: flex; align-items: center; justify-content: space-between; background: rgba(22, 34, 43, 0.02); border-radius: 8px; padding: 0.8rem; border: 1px solid rgba(22, 34, 43, 0.08); flex-grow: 1; }
.cpu-step { display: flex; flex-direction: column; gap: 0.3rem; align-items: center; width: 42%; }
.bg-dark { background: #11181c; padding: 0.5rem 0; border-radius: 6px; box-shadow: inset 0 2px 10px rgba(0,0,0,0.3); border: 1px solid rgba(255,255,255,0.05); color: white !important; }
.box-bg { background: white; padding: 0.5rem 0; border-radius: 6px; border: 2px solid var(--storm-cyan); box-shadow: 0 2px 6px rgba(13, 159, 184, 0.15); }
.wv-label { font-size: 0.7rem; font-weight: 800; text-align: center; }
.wv-detail { font-family: "IBM Plex Mono", monospace; font-size: 0.55rem; color: var(--storm-slate); font-weight: 700; text-align: center; }
.wv-arrow { color: var(--storm-slate); font-size: 1rem; font-weight: 800; opacity: 0.6; }

/* -- Metrics Visual & Tag -- */
.metrics-visual { margin-top: auto; display: flex; flex-direction: column; gap: 0.8rem; background: rgba(22, 34, 43, 0.02); border-radius: 8px; padding: 1rem; border: 1px solid rgba(22, 34, 43, 0.08); flex-grow: 1; justify-content: center; }
.mv-row { display: flex; justify-content: space-between; align-items: center; padding: 0.2rem 0; }
.mv-label { font-size: 0.7rem; font-weight: 700; color: var(--storm-slate); text-transform: uppercase; letter-spacing: 0.03em; }
.mv-val { font-size: 1.2rem; font-weight: 800; color: var(--storm-ink); font-family: "IBM Plex Mono", monospace; }
.code-font { color: var(--storm-cyan); font-size: 1.6rem; letter-spacing: -0.05em; }

.cookie-tag { display: flex; align-items: center; justify-content: space-between; background: rgba(13, 159, 184, 0.08); border: 1px solid rgba(13, 159, 184, 0.25); border-radius: 6px; padding: 0.5rem 0.8rem; }
.tag-label { font-size: 0.65rem; font-weight: 600; color: var(--storm-ink); }
.tag-value { font-family: "IBM Plex Mono", monospace; font-size: 0.55rem; color: var(--storm-cyan); font-weight: 700;}
.text-slate { color: var(--storm-slate) !important; font-family: var(--storm-font); }
.mt-auto { margin-top: auto; }
</style>
