<script setup lang="ts">
// Quality & Security Slide combining redundant components
</script>

<template>
  <div class="page-shell">
    <div class="header-minimal">
      <p class="eyebrow">04 — Qualité & Sécurité</p>
      <h2 class="section-title">Une défense en profondeur garantie par l'intégration continue</h2>
    </div>

    <div class="qs-grid">
      <!-- Sécurité Périmétrique -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box red-bg"><img src="https://api.iconify.design/carbon/security.svg" class="icon-sm"/></div>
          <h3 class="card-title">Défense Périmétrique</h3>
        </header>
        <div class="card-body">
          <ul class="feature-list">
            <li>
              <strong>Rate Limiting strict</strong>
              <span>5 requêtes/min (burst=3) sur <code>/auth/*</code> pour bloquer le brute-force.</span>
            </li>
            <li>
              <strong>Security Headers</strong>
              <span>CSP stricte, HSTS, X-Frame-Options configurés nativement sur le Gateway Go.</span>
            </li>
            <li>
              <strong>CORS restrictif</strong>
              <span>Rejet par défaut de toute origine non-déclarée en environnement de production.</span>
            </li>
          </ul>
        </div>
      </article>

      <!-- Tests & Couverture -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box green-bg"><img src="https://api.iconify.design/carbon/task-approved.svg" class="icon-sm"/></div>
          <h3 class="card-title">Couverture de Tests (> 80%)</h3>
        </header>
        <div class="card-body">
          <div class="coverage-stats">
            <div class="stat-item">
              <span class="stat-val text-cyan">80.4%</span>
              <span class="stat-label">Gateway Go</span>
            </div>
            <div class="stat-item">
              <span class="stat-val text-orange">83.6%</span>
              <span class="stat-label">Service Messages</span>
            </div>
          </div>
          <p class="text-sm mt-3">
            Couverture centrée sur le comportement métier réel : handlers HTTP, cycle de vie WebSocket, pool de workers et consommation NATS.
          </p>
        </div>
      </article>

      <!-- Pipeline CI/CD -->
      <article class="clean-card shadow-sm full-span">
        <header class="card-head">
          <div class="icon-box slate-bg"><img src="https://api.iconify.design/logos/github-actions.svg" class="icon-sm color-keep"/></div>
          <h3 class="card-title">Pipeline d'Intégration Continue (GitHub Actions)</h3>
        </header>
        <div class="card-body">
          <div class="pipeline-flow">
            <div class="pl-step">Push Code</div>
            <div class="pl-arr">➔</div>
            <div class="pl-step border-cyan">Lint & Build</div>
            <div class="pl-arr">➔</div>
            <div class="pl-step border-green">Unit Tests (Go)</div>
            <div class="pl-arr">➔</div>
            <div class="pl-step border-red">Gosec (SAST)</div>
            <div class="pl-arr">➔</div>
            <div class="pl-step border-purple">Trivy (Docker)</div>
          </div>
          <p class="text-sm text-center mt-3">
            Chaque commit est automatiquement vérifié. Le pipeline valide la qualité du code (tests unitaire) et scanne les dépendances pour éviter l'introduction de vulnérabilités connues (CVE).
          </p>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
/* Same clean aesthetic as ArchDetailsSlide */
.page-shell { display: flex; flex-direction: column; height: 100%; }
.header-minimal { margin-bottom: 0.8rem; }

.qs-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: auto auto;
  gap: 0.8rem;
  flex-grow: 1;
}

.full-span { grid-column: 1 / -1; }

.clean-card {
  background: white;
  border-radius: 12px;
  padding: 0.8rem 1.2rem;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0,0,0,0.04);
}

.card-head { display: flex; align-items: center; gap: 0.6rem; margin-bottom: 0.4rem; }
.card-title { font-size: 0.9rem; font-weight: 700; color: var(--storm-ink); margin: 0; }
.icon-box { width: 24px; height: 24px; border-radius: 6px; display: flex; align-items: center; justify-content: center; }
.icon-sm { width: 14px; filter: grayscale(1) brightness(0) invert(1); }
.icon-sm.color-keep { filter: none; }

.red-bg { background: #ef4444; box-shadow: 0 4px 12px rgba(239, 68, 68, 0.2); }
.green-bg { background: #10b981; box-shadow: 0 4px 12px rgba(16, 185, 129, 0.2); }
.slate-bg { background: #f8fafc; border: 1px solid #e2e8f0; box-shadow: 0 2px 8px rgba(0,0,0,0.02); }

.card-body { display: flex; flex-direction: column; flex-grow: 1; }

.text-sm { font-size: 0.72rem; line-height: 1.35; color: #475569; margin: 0; }
.text-center { text-align: center; }
.mt-2 { margin-top: 0.4rem; }
.mt-3 { margin-top: 0.6rem; }
.text-cyan { color: var(--storm-cyan); }
.text-orange { color: var(--storm-orange); }

.feature-list {
  list-style: none; padding: 0; margin: 0;
  display: flex; flex-direction: column; gap: 0.4rem;
}
.feature-list li {
  display: flex; flex-direction: column; gap: 0.15rem;
  padding-bottom: 0.3rem; border-bottom: 1px solid #f1f5f9;
}
.feature-list li:last-child { border-bottom: none; padding-bottom: 0; }
.feature-list strong { font-size: 0.72rem; color: var(--storm-ink); }
.feature-list span { font-size: 0.65rem; color: #64748b; line-height: 1.25; }
.feature-list code { font-size: 0.65rem; background: #f8fafc; padding: 0.1rem 0.3rem; border-radius: 4px; color: var(--storm-cyan); border: 1px solid #e2e8f0; }

.coverage-stats {
  display: flex; gap: 1.5rem; justify-content: center;
  background: #f8fafc; padding: 0.8rem; border-radius: 8px; border: 1px solid #f1f5f9;
}
.stat-item {
  display: flex; flex-direction: column; gap: 0.1rem; align-items: center;
}
.stat-val { font-size: 1.4rem; font-weight: 800; font-family: monospace; letter-spacing: -0.02em; line-height: 1; }
.stat-label { font-size: 0.6rem; font-weight: 700; color: #64748b; text-transform: uppercase; }

.pipeline-flow {
  display: flex; align-items: center; justify-content: space-between;
  background: #f8fafc; padding: 0.6rem 0.8rem; border-radius: 8px; border: 1px solid #f1f5f9;
  margin-top: auto;
}
.pl-step {
  background: white; padding: 0.3rem 0.5rem; border-radius: 6px; font-size: 0.65rem; font-weight: 600; color: var(--storm-ink);
  border: 1px solid #e2e8f0; text-align: center; box-shadow: 0 1px 2px rgba(0,0,0,0.02);
}
.pl-arr { font-size: 0.7rem; font-weight: 700; color: #cbd5e1; }

.border-cyan { border-left: 3px solid var(--storm-cyan); }
.border-green { border-left: 3px solid #10b981; }
.border-red { border-left: 3px solid #ef4444; }
.border-purple { border-left: 3px solid #8b5cf6; }
</style>
