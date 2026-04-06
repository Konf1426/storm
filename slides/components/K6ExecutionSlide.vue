<script setup lang="ts">
// Slide 10 – Why k6 over Azure Load Testing, and how we ran it.
</script>

<template>
  <div class="page-shell k6-shell">
    <div class="k6-header">
      <p class="eyebrow">10 — Outil de charge</p>
      <h2 class="section-title">Pourquoi k6 et pas Azure Load Testing ?</h2>
    </div>

    <!-- Comparison: Azure Load Testing vs k6 -->
    <div class="compare-row">
      <article class="compare-card rejected">
        <div class="compare-head">
          <div class="icon-box red-glow">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#dc2626" stroke-width="2.5" stroke-linecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </div>
          <h3>Azure Load Testing</h3>
        </div>
        <ul class="compare-list">
          <li><span class="compare-tag cost">~1 $/1k VUs</span> Facturation à l'usage</li>
          <li><span class="compare-tag lock">Fermé</span> Scénarios limités</li>
          <li><span class="compare-tag lock">Opaque</span> Pas de contrôle fin</li>
        </ul>
      </article>

      <div class="compare-vs">
        <span>VS</span>
      </div>

      <article class="compare-card chosen">
        <div class="compare-head">
          <div class="icon-box cyan-glow">
            <img src="https://api.iconify.design/simple-icons/k6.svg?color=%230d9fb8" class="icon-sm color-keep" />
          </div>
          <h3>k6 Open Source</h3>
        </div>
        <ul class="compare-list">
          <li><span class="compare-tag free">Gratuit</span> Coût = infra seulement</li>
          <li><span class="compare-tag open">JS</span> Scénarios sur mesure</li>
          <li><span class="compare-tag open">Full</span> Contrôle total</li>
        </ul>
      </article>
    </div>

    <!-- Execution pipeline -->
    <div class="pipeline-row">
      <div class="pipe-step">
        <div class="pipe-num">1</div>
        <div class="pipe-body">
          <strong>Local</strong>
          <span>Docker Compose<br/>Validation rapide</span>
        </div>
      </div>
      <div class="pipe-arrow">
        <svg width="28" height="12" viewBox="0 0 28 12"><path d="M0 6h24M20 1l5 5-5 5" fill="none" stroke="#94a3b8" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </div>
      <div class="pipe-step">
        <div class="pipe-num">2</div>
        <div class="pipe-body">
          <strong>Azure AKS</strong>
          <span>Job k6 distribué<br/>Scale des VUs</span>
        </div>
      </div>
      <div class="pipe-arrow">
        <svg width="28" height="12" viewBox="0 0 28 12"><path d="M0 6h24M20 1l5 5-5 5" fill="none" stroke="#94a3b8" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </div>
      <div class="pipe-step">
        <div class="pipe-num">3</div>
        <div class="pipe-body">
          <strong>Arrêt</strong>
          <span>Infra coupée<br/>0 coût résiduel</span>
        </div>
      </div>
    </div>

    <!-- Bottom banner -->
    <div class="k6-banner">
      <div class="banner-stats">
        <div class="banner-stat">
          <span class="stat-label">Coût bench</span>
          <strong>0,15 $/h</strong>
          <span class="stat-sub">en veille</span>
        </div>
        <div class="banner-divider"></div>
        <div class="banner-stat">
          <span class="stat-label">Coût pic</span>
          <strong>3,33 $/h</strong>
          <span class="stat-sub">10k VUs</span>
        </div>
        <div class="banner-divider"></div>
        <div class="banner-stat">
          <span class="stat-label">Économie</span>
          <strong>~90 %</strong>
          <span class="stat-sub">vs managed</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.k6-shell {
  display: flex;
  flex-direction: column;
  height: 100%;
  gap: 0.75rem;
  justify-content: center;
}

.k6-header {
  display: grid;
  gap: 0.1rem;
}

/* ── Comparison row ── */
.compare-row {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 0.6rem;
  align-items: center;
}

.compare-card {
  border-radius: 18px;
  padding: 0.8rem 0.9rem;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  border: 1px solid rgba(22, 34, 43, 0.06);
  box-shadow: 0 4px 16px rgba(18, 28, 36, 0.04);
}

.compare-card.rejected {
  background: rgba(255, 255, 255, 0.72);
  opacity: 0.72;
}

.compare-card.chosen {
  background: var(--storm-surface-strong);
  border-color: rgba(13, 159, 184, 0.22);
  box-shadow: 0 4px 20px rgba(13, 159, 184, 0.10);
}

.compare-head {
  display: flex;
  align-items: center;
  gap: 0.55rem;
}

.compare-head h3 {
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--storm-ink);
  margin: 0;
}

.compare-vs {
  display: flex;
  align-items: center;
  justify-content: center;
}

.compare-vs span {
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.62rem;
  font-weight: 800;
  letter-spacing: 0.1em;
  color: #94a3b8;
  background: rgba(22, 34, 43, 0.04);
  border-radius: 999px;
  padding: 0.3rem 0.55rem;
}

.compare-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 0.45rem;
}

.compare-list li {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.62rem;
  color: var(--storm-muted);
  line-height: 1.3;
}

.compare-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  padding: 0.12rem 0.38rem;
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.52rem;
  font-weight: 700;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  flex-shrink: 0;
  min-width: 48px;
  text-align: center;
}

.compare-tag.cost {
  background: rgba(220, 38, 38, 0.10);
  color: #dc2626;
}

.compare-tag.lock {
  background: rgba(220, 38, 38, 0.06);
  color: #b91c1c;
}

.compare-tag.free {
  background: rgba(16, 185, 129, 0.12);
  color: #047857;
}

.compare-tag.open {
  background: rgba(13, 159, 184, 0.12);
  color: #0d9fb8;
}

.icon-box {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(22, 34, 43, 0.04);
  flex-shrink: 0;
}

.icon-sm {
  width: 16px;
  filter: grayscale(1) brightness(0);
}

.color-keep { filter: none; }
.cyan-glow { box-shadow: 0 0 12px rgba(13, 159, 184, 0.18); }
.red-glow { box-shadow: 0 0 12px rgba(220, 38, 38, 0.14); }

/* ── Pipeline row ── */
.pipeline-row {
  display: flex;
  align-items: stretch;
  gap: 0.4rem;
  justify-content: center;
}

.pipe-step {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--storm-surface-strong);
  border: 1px solid rgba(22, 34, 43, 0.06);
  border-radius: 14px;
  padding: 0.55rem 0.75rem;
  min-width: 120px;
}

.pipe-num {
  width: 24px;
  height: 24px;
  border-radius: 999px;
  background: rgba(239, 138, 41, 0.12);
  color: var(--storm-orange);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.56rem;
  font-weight: 800;
  flex-shrink: 0;
}

.pipe-body {
  display: grid;
  gap: 0.1rem;
}

.pipe-body strong {
  font-size: 0.66rem;
  color: var(--storm-ink);
}

.pipe-body span {
  font-size: 0.52rem;
  color: var(--storm-muted);
  line-height: 1.28;
}

.pipe-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* ── Bottom banner ── */
.k6-banner {
  border-radius: 16px;
  background: rgba(16, 185, 129, 0.06);
  border: 1px solid rgba(16, 185, 129, 0.16);
  padding: 0.6rem 1rem;
}

.banner-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1.2rem;
}

.banner-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.05rem;
}

.banner-stat .stat-label {
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.48rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #047857;
}

.banner-stat strong {
  font-size: 0.88rem;
  font-weight: 800;
  color: #065f46;
  letter-spacing: -0.03em;
}

.banner-stat .stat-sub {
  font-size: 0.46rem;
  color: #6b7280;
}

.banner-divider {
  width: 1px;
  height: 28px;
  background: rgba(16, 185, 129, 0.22);
}
</style>
