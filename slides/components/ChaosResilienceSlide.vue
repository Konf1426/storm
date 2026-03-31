<script setup lang="ts">
// Chaos & Resilience Slide
</script>

<template>
  <div class="page-shell">
    <div class="header-minimal">
      <p class="eyebrow">10 — Chaos & Résilience</p>
      <h2 class="section-title">Mesurer la robustesse par l'injection de fautes</h2>
    </div>

    <div class="chaos-grid">
      <!-- NATS Pause -->
      <article class="clean-card shadow-sm experiment-card">
        <div class="card-status status-success">Test Passé</div>
        <header class="card-head">
          <div class="icon-box orange-bg"><img src="https://api.iconify.design/carbon/flash-off.svg" class="icon-sm"/></div>
          <h3 class="card-title">Injection : Coupure NATS</h3>
        </header>
        <div class="card-body">
          <div class="experiment-flow">
            <div class="exp-step">
              <span class="step-lbl">Action</span>
              <p>Arrêt brutal du broker de messages (Docker Stop).</p>
            </div>
            <div class="exp-step">
              <span class="step-lbl">Observation</span>
              <p>Le Gateway Go bufferise les messages en RAM. Les clients WS sont notifiés.</p>
            </div>
            <div class="exp-step">
              <span class="step-lbl">Résultat</span>
              <div class="res-box">
                 <img src="https://api.iconify.design/carbon/checkmark-filled.svg?color=%2310b981" class="res-icon"/>
                 <span><strong>Rétablissement immédiat</strong> dès le redémarrage. Zéro message perdu.</span>
              </div>
            </div>
          </div>
        </div>
      </article>

      <!-- DB Restart -->
      <article class="clean-card shadow-sm experiment-card">
        <div class="card-status status-success">Test Passé</div>
        <header class="card-head">
          <div class="icon-box slate-bg"><img src="https://api.iconify.design/carbon/restart.svg" class="icon-sm"/></div>
          <h3 class="card-title">Injection : Redémarrage DB</h3>
        </header>
        <div class="card-body">
          <div class="experiment-flow">
            <div class="exp-step">
              <span class="step-lbl">Action</span>
              <p>Restart de l'instance PostgreSQL pendant un pic de charge.</p>
            </div>
            <div class="exp-step">
              <span class="step-lbl">Observation</span>
              <p>Le Worker Pool passe en mode <i>retry</i>. Les WebSockets restent ouverts.</p>
            </div>
            <div class="exp-step">
              <span class="step-lbl">Résultat</span>
              <div class="res-box">
                 <img src="https://api.iconify.design/carbon/checkmark-filled.svg?color=%2310b981" class="res-icon"/>
                 <span><strong>Découplage réussi</strong>. Le chat continue, la persistance rattrape son retard.</span>
              </div>
            </div>
          </div>
        </div>
      </article>
    </div>

    <div class="lesson-bar shadow-sm">
       <div class="lesson-icon"><img src="https://api.iconify.design/carbon/idea.svg?color=%23f59e0b" style="width:16px;"/></div>
       <div class="lesson-text">
          <strong>Leçon apprise :</strong> "Un crash-test réussi n'est pas un test sans panne, c'est un test dont on gère le rétablissement automatique."
       </div>
    </div>
  </div>
</template>

<style scoped>
.page-shell { display: flex; flex-direction: column; height: 100%; }
.header-minimal { margin-bottom: 1.2rem; }

.chaos-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  flex-grow: 1;
}

.experiment-card {
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.card-status {
  position: absolute; top: 0.8rem; right: 1rem;
  font-size: 0.55rem; font-weight: 800; text-transform: uppercase;
  padding: 0.2rem 0.5rem; border-radius: 4px;
}
.status-success { background: #dcfce7; color: #166534; border: 1px solid #bbf7d0; }

.clean-card {
  background: white; border-radius: 12px; padding: 1.4rem;
  display: flex; flex-direction: column; box-shadow: 0 4px 16px rgba(0, 0, 0, 0.03); border: 1px solid rgba(0,0,0,0.04);
}

.card-head { display: flex; align-items: center; gap: 0.6rem; margin-bottom: 1rem; }
.card-title { font-size: 1rem; font-weight: 700; color: var(--storm-ink); margin: 0; }

.icon-box { width: 30px; height: 30px; border-radius: 8px; display: flex; align-items: center; justify-content: center; }
.icon-sm { width: 16px; filter: grayscale(1) brightness(0) invert(1); }

.orange-bg { background: #f97316; }
.slate-bg { background: #475569; }

.card-body { display: flex; flex-direction: column; flex-grow: 1; }

.experiment-flow { font-size: 0.75rem; color: #4b5563; line-height: 1.4; display: flex; flex-direction: column; gap: 0.8rem; }
.step-lbl { display: block; font-weight: 700; font-size: 0.55rem; text-transform: uppercase; color: #9ca3af; margin-bottom: 0.1rem; }
.experiment-flow p { margin: 0; }

.res-box {
  margin-top: 0.4rem; background: #f8fafc; padding: 0.6rem; border-radius: 6px; border: 1px solid #f1f5f9;
  display: flex; align-items: center; gap: 0.6rem; font-size: 0.7rem;
}
.res-icon { width: 14px; flex-shrink: 0; }

.lesson-bar {
  margin-top: 1.2rem; background: #fffbe3; border: 1px solid #fef3c7; border-radius: 10px;
  display: flex; align-items: center; gap: 0.8rem; padding: 1rem 1.4rem;
}
.lesson-text { font-size: 0.75rem; color: #92400e; font-style: italic; }
.lesson-text strong { font-style: normal; color: #78350f; font-weight: 700; }
</style>
