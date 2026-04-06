<script setup lang="ts">
// Slide 06 – Real product showcase with app screenshot and notification flow highlights.
</script>

<template>
  <div class="page-shell product-shell">
    <div class="product-header">
      <p class="eyebrow">06 — Produit réel &amp; notifications</p>
      <h2 class="section-title">Une vraie messagerie, pas seulement une infra</h2>
    </div>

    <div class="product-layout">
      <!-- Left: Real app screenshot -->
      <div class="screenshot-zone">
        <div class="browser-chrome">
          <div class="chrome-dots">
            <span class="dot dot-red"></span>
            <span class="dot dot-amber"></span>
            <span class="dot dot-green"></span>
          </div>
          <span class="chrome-url">localhost:5173</span>
        </div>
        <div class="screenshot-frame">
          <img src="/storm-app-screenshot.png" alt="STORM Dashboard" class="app-screenshot" />
          <!-- Callout overlays -->
          <div class="callout callout-toast">
            <div class="callout-line"></div>
            <div class="callout-badge">
              <span class="callout-num">1</span>
              Toast de lecture
            </div>
          </div>
          <div class="callout callout-feed">
            <div class="callout-line"></div>
            <div class="callout-badge">
              <span class="callout-num">2</span>
              Feed temps réel
            </div>
          </div>
          <div class="callout callout-send">
            <div class="callout-line"></div>
            <div class="callout-badge">
              <span class="callout-num">3</span>
              Statut d'envoi
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Feature highlights -->
      <div class="features-col">
        <div class="feature-pill">
          <div class="pill-head">
            <span class="pill-num">1</span>
            <strong>Notification de lecture</strong>
          </div>
          <p>Quand un utilisateur lit un message, l'auteur reçoit un toast en temps réel via WebSocket.</p>
          <div class="pill-route">
            <code>POST /channels/:id/messages/:msgId/read</code>
          </div>
        </div>

        <div class="feature-pill">
          <div class="pill-head">
            <span class="pill-num">2</span>
            <strong>Flux live bidirectionnel</strong>
          </div>
          <p>Les messages arrivent instantanément via WebSocket. Mes messages sont en vert, ceux des autres en blanc.</p>
        </div>

        <div class="feature-pill">
          <div class="pill-head">
            <span class="pill-num">3</span>
            <strong>Feedback d'envoi</strong>
          </div>
          <p>Après un clic sur Send, l'UI affiche <code>sending</code> puis <code>sent</code> — confirmation immédiate.</p>
        </div>

        <div class="flow-banner">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#0d9fb8" stroke-width="2" stroke-linecap="round"><path d="M22 2L11 13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/></svg>
          <span>Envoi → NATS → Diffusion WS → Toast de lecture</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.product-shell {
  display: flex;
  flex-direction: column;
  height: 100%;
  gap: 0.5rem;
}

.product-header {
  display: grid;
  gap: 0.1rem;
}

/* ── Main layout ── */
.product-layout {
  display: grid;
  grid-template-columns: 1.4fr 0.85fr;
  gap: 0.9rem;
  align-items: start;
  flex: 1;
  min-height: 0;
}

/* ── Screenshot zone ── */
.screenshot-zone {
  display: flex;
  flex-direction: column;
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid rgba(22, 34, 43, 0.1);
  box-shadow: 0 8px 32px rgba(18, 28, 36, 0.08);
  position: relative;
  max-height: 100%;
}

.browser-chrome {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.4rem 0.65rem;
  background: #1e293b;
}

.chrome-dots {
  display: flex;
  gap: 0.25rem;
}

.dot {
  width: 7px;
  height: 7px;
  border-radius: 999px;
}

.dot-red { background: #fb7185; }
.dot-amber { background: #fbbf24; }
.dot-green { background: #34d399; }

.chrome-url {
  margin-left: 0.4rem;
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.48rem;
  color: #94a3b8;
  letter-spacing: 0.04em;
}

.screenshot-frame {
  position: relative;
  background: #f5f0e7;
  overflow: hidden;
  max-height: 380px;
}

.app-screenshot {
  width: 100%;
  display: block;
}

/* ── Callout overlays ── */
.callout {
  position: absolute;
  display: flex;
  align-items: center;
  gap: 0.3rem;
  pointer-events: none;
}

.callout-badge {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  background: rgba(13, 159, 184, 0.92);
  color: white;
  border-radius: 999px;
  padding: 0.18rem 0.5rem 0.18rem 0.25rem;
  font-size: 0.48rem;
  font-weight: 700;
  white-space: nowrap;
  box-shadow: 0 2px 8px rgba(13, 159, 184, 0.3);
}

.callout-num {
  width: 16px;
  height: 16px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.44rem;
  font-weight: 800;
}

.callout-line {
  width: 20px;
  height: 2px;
  background: rgba(13, 159, 184, 0.6);
  border-radius: 2px;
}

.callout-toast {
  top: 4%;
  right: -2%;
  flex-direction: row-reverse;
}

.callout-feed {
  top: 48%;
  right: -2%;
  flex-direction: row-reverse;
}

.callout-send {
  bottom: 12%;
  left: 2%;
}

/* ── Features column ── */
.features-col {
  display: flex;
  flex-direction: column;
  gap: 0.55rem;
}

.feature-pill {
  background: var(--storm-surface-strong);
  border: 1px solid rgba(22, 34, 43, 0.06);
  border-radius: 16px;
  padding: 0.65rem 0.75rem;
  display: grid;
  gap: 0.25rem;
  box-shadow: 0 4px 12px rgba(18, 28, 36, 0.03);
}

.pill-head {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.pill-head strong {
  font-size: 0.7rem;
  color: var(--storm-ink);
}

.pill-num {
  width: 20px;
  height: 20px;
  border-radius: 999px;
  background: rgba(13, 159, 184, 0.12);
  color: var(--storm-cyan);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.5rem;
  font-weight: 800;
  flex-shrink: 0;
}

.feature-pill p {
  font-size: 0.56rem;
  line-height: 1.35;
  color: var(--storm-muted);
  margin: 0;
}

.feature-pill code {
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.5rem;
  background: rgba(22, 34, 43, 0.05);
  padding: 0.08rem 0.2rem;
  border-radius: 4px;
}

.pill-route {
  margin-top: 0.15rem;
}

.pill-route code {
  font-size: 0.46rem;
  background: #16222b;
  color: #d9e4ea;
  padding: 0.2rem 0.45rem;
  border-radius: 6px;
  display: inline-block;
}

/* ── Flow banner ── */
.flow-banner {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border-radius: 14px;
  background: rgba(13, 159, 184, 0.06);
  border: 1px solid rgba(13, 159, 184, 0.14);
  padding: 0.5rem 0.7rem;
}

.flow-banner span {
  font-size: 0.56rem;
  font-weight: 600;
  color: var(--storm-cyan);
  letter-spacing: -0.01em;
}
</style>
