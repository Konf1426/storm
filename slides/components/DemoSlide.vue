<script setup lang="ts">
// Demo & Fallback Slide
</script>

<template>
  <div class="page-shell">
    <div class="header-minimal">
      <p class="eyebrow">11 — Démo Live & Fallback</p>
      <h2 class="section-title">Démonstration maîtrisée et plan de secours prêt</h2>
    </div>

    <div class="demo-grid">
      <!-- Demo Script -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box cyan-bg"><img src="https://api.iconify.design/carbon/screen.svg" class="icon-sm"/></div>
          <h3 class="card-title">Script de Démonstration</h3>
        </header>
        <div class="card-body">
          <ol class="demo-list">
             <li>
                <strong>Auth & Session</strong>
                <span>Connexion sécurisée et émission des cookies HttpOnly.</span>
             </li>
             <li>
                <strong>Temps Réel Broadcast</strong>
                <span>Échange de messages instantané entre deux fenêtres via NATS.</span>
             </li>
             <li>
                <strong>Persistance & Sync</strong>
                <span>Vérification de l'insertion asynchrone en base PostgreSQL.</span>
             </li>
             <li>
                <strong>Supervision Grafana</strong>
                <span>Visualisation immédiate du pic de messages sur le dashboard.</span>
             </li>
          </ol>
        </div>
      </article>

      <!-- Fallback -->
      <article class="clean-card shadow-sm border-orange-light">
        <header class="card-head">
          <div class="icon-box orange-bg"><img src="https://api.iconify.design/carbon/life-buoy.svg" class="icon-sm"/></div>
          <h3 class="card-title">Plan de Secours (Fallback)</h3>
        </header>
        <div class="card-body">
          <p class="text-sm mb-4">En cas d'instabilité réseau ou serveur durant l'oral :</p>
          <ul class="fallback-list">
             <li v-click>
                <strong>Preuves d'Exécution</strong>
                <span>Résultats du Storm Day déjà versionnés et horodatés dans le repo.</span>
             </li>
             <li v-click>
                <strong>Scripts d'Automatisation</strong>
                <span>Présentation des scripts <code>perf-load.sh</code> et <code>storm-day-runner.sh</code>.</span>
             </li>
             <li v-click>
                <strong>Analyse Post-Mortem</strong>
                <span>Focus sur les métriques et logs capturés lors des derniers stress-tests.</span>
             </li>
          </ul>
        </div>
      </article>
    </div>

    <div class="demo-note shadow-sm">
       <img src="https://api.iconify.design/carbon/information.svg?color=%230ea5e9" style="width:14px;"/>
       <span>"On ne fait pas reposer la note sur un succès réseau, mais sur la solidité des preuves apportées."</span>
    </div>
  </div>
</template>

<style scoped>
.page-shell { display: flex; flex-direction: column; height: 100%; }
.header-minimal { margin-bottom: 1rem; }

.demo-grid {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 1.5rem;
  flex-grow: 1;
}

.clean-card {
  background: white; border-radius: 12px; padding: 1.4rem;
  display: flex; flex-direction: column; box-shadow: 0 4px 16px rgba(0, 0, 0, 0.03); border: 1px solid rgba(0,0,0,0.04);
}
.border-orange-light { border: 1px solid #fed7aa; background: #fffcf9; }

.card-head { display: flex; align-items: center; gap: 0.6rem; margin-bottom: 1.2rem; }
.card-title { font-size: 1rem; font-weight: 700; color: var(--storm-ink); margin: 0; }

.icon-box { width: 30px; height: 30px; border-radius: 8px; display: flex; align-items: center; justify-content: center; }
.icon-sm { width: 18px; filter: grayscale(1) brightness(0) invert(1); }

.cyan-bg { background: var(--storm-cyan); }
.orange-bg { background: var(--storm-orange); }

.card-body { display: flex; flex-direction: column; flex-grow: 1; }

.demo-list { padding: 0 0 0 1.2rem; margin: 0; display: flex; flex-direction: column; gap: 1rem; }
.demo-list li { font-size: 0.8rem; color: #475569; position: relative; }
.demo-list li strong { display: block; color: var(--storm-ink); font-size: 0.85rem; }
.demo-list li span { font-size: 0.75rem; color: #64748b; line-height: 1.4; }

.fallback-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 1.2rem; justify-content: center; flex-grow: 1; }
.fallback-list li { background: white; border: 1px solid #f1f5f9; padding: 0.8rem; border-radius: 8px; display: flex; flex-direction: column; gap: 0.2rem; border-left: 3px solid var(--storm-orange); }
.fallback-list strong { font-size: 0.8rem; color: var(--storm-ink); }
.fallback-list span { font-size: 0.7rem; color: #64748b; line-height: 1.3; font-style: italic; }

.text-sm { font-size: 0.75rem; color: #64748b; font-weight: 500; }
.mb-4 { margin-bottom: 1rem; }

.demo-note {
  margin-top: 1.2rem; background: #f0f9ff; border: 1px solid #e0f2fe; border-radius: 8px;
  padding: 0.8rem 1.2rem; display: flex; align-items: center; gap: 0.8rem;
  font-size: 0.75rem; color: #0369a1; font-weight: 500; text-align: center; justify-content: center;
}
</style>
