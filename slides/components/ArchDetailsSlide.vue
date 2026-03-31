<script setup lang="ts">
// Vue component for Architecture Detailed Bento Grid
</script>

<template>
  <div class="page-shell bento-shell">
    <div class="bento-header">
      <p class="eyebrow">05 — Détails techniques</p>
      <h2 class="section-title">Anatomie du pipeline temps réel</h2>
    </div>

    <div class="bento-grid">
      <!-- Auth (Grosse carte à gauche) -->
      <article class="bento-card bento-auth">
        <div class="bento-card-header">
          <div class="bento-icon cyan-glow"><img src="/icons/vue.svg" style="width:16px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Sécurité & JWT</h3>
        </div>
        <p class="bento-desc">
          Double token pour contrer les XSS tout en gardant une API stateless.
        </p>
        
        <div class="jwt-visual-grid">
          <!-- Header -->
          <div class="jwt-row c-red-zone">
            <div class="jwt-raw-part">eyJhbGci...</div>
            <div class="jwt-json-part">
              <div class="panel-head">En-tête</div>
              <div class="panel-body">{ "alg": "HS256" }</div>
            </div>
          </div>
          
          <!-- Payload -->
          <div class="jwt-row c-purple-zone">
            <div class="jwt-raw-part">eyJzdWIi...</div>
            <div class="jwt-json-part">
              <div class="panel-head">Payload</div>
              <pre class="panel-code">
{
  "sub": "usr_9x2b",
  "nosnif": true,
  "urldomaine": "http://localhost:5173",
  "iat": 172910283,
  "exp": 172911183
}</pre>
            </div>
          </div>
          
          <!-- Signature -->
          <div class="jwt-row c-blue-zone">
            <div class="jwt-raw-part">SflKxwR...</div>
            <div class="jwt-json-part">
              <div class="panel-head">Signature</div>
              <div class="panel-body">HMAC(Secret)</div>
            </div>
          </div>
        </div>
        <div class="cookie-tag">
          <span class="tag-label">Refresh Token</span>
          <span class="tag-value">Cookie HttpOnly (24h)</span>
        </div>
      </article>

      <!-- Temps Réel (Carte large en haut à droite) -->
      <article class="bento-card bento-rt">
        <div class="bento-card-header">
          <div class="bento-icon orange-glow"><img src="/icons/nats.png" style="width:16px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Routage Temps Réel (Pub/Sub)</h3>
        </div>
        <div class="rt-content">
          <p class="bento-desc">
            Si <strong>Alice</strong> envoie un message, le Gateway le confie au broker <strong>NATS</strong>, qui le redistribue <br/> instantanément à tous les utilisateurs (ex: <strong>Bob</strong>) connectés au même salon :
          </p>
          <div class="rt-visual">
            <div class="rt-user">Alice <span class="ws-tag">WS</span></div>
            <div class="rt-arrow">
               <span class="rt-pub c-orange">Publie</span>
               <div class="rt-line"><div class="rt-dot"></div></div>
            </div>
            <div class="rt-nats">
              Broker NATS
              <code>T: room.142.events</code>
            </div>
            <div class="rt-arrow">
               <span class="rt-pub c-cyan">Reçoit</span>
               <div class="rt-line dashed"></div>
            </div>
            <div class="rt-user">Bob <span class="ws-tag">WS</span></div>
          </div>
        </div>
      </article>

      <!-- Persistance (En bas au centre) -->
      <article class="bento-card bento-db">
        <div class="bento-card-header">
          <div class="bento-icon slate-glow"><img src="/icons/postgres.svg" style="width:14px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Sauvegarde Asynchrone</h3>
        </div>
        <p class="bento-desc">
          On détache l'insertion en base du flux réseau WebSocket :
        </p>
        
        <div class="worker-visual">
          <div class="wv-step">
            <span class="wv-label">WebSocket</span>
            <span class="wv-speed c-cyan">< 1ms</span>
          </div>
          <div class="wv-arrow">➔</div>
          <div class="wv-step box-bg">
            <span class="wv-detail">Channel Go</span>
          </div>
          <div class="wv-arrow">➔</div>
          <div class="wv-step">
            <span class="wv-label">Postgres (Worker)</span>
            <span class="wv-speed c-orange">> 10ms</span>
          </div>
        </div>
      </article>

      <!-- Supervision (En bas à droite) -->
      <article class="bento-card bento-obs">
        <div class="bento-card-header">
          <div class="bento-icon cyan-glow"><img src="/icons/prometheus.svg" style="width:16px;filter:grayscale(1) brightness(0);"/></div>
          <h3>Observabilité</h3>
        </div>
        <p class="bento-desc">
          Métriques Go exposées sur `/metrics`.
        </p>
        <div class="obs-visual">
          <span class="metric-name">storm_active_websockets</span>
          <div class="sparkline">
            <div class="spark-bar" style="height: 30%"></div>
            <div class="spark-bar" style="height: 45%"></div>
            <div class="spark-bar" style="height: 60%"></div>
            <div class="spark-bar" style="height: 50%"></div>
            <div class="spark-bar" style="height: 80%"></div>
            <div class="spark-bar" style="height: 100%; background: var(--storm-orange)"></div>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.bento-shell {
  display: flex;
  flex-direction: column;
  height: 100%;
  gap: 0.5rem;
}

.bento-header {
  display: grid;
  gap: 0.1rem;
  margin-bottom: 0;
}

.bento-grid {
  display: grid;
  flex-grow: 1;
  grid-template-columns: 1.25fr 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 0.8rem;
  min-height: 0;
}

.bento-auth { grid-column: 1; grid-row: 1 / 3; }
.bento-rt { grid-column: 2 / 4; grid-row: 1; }
.bento-db { grid-column: 2; grid-row: 2; }
.bento-obs { grid-column: 3; grid-row: 2; }

.bento-card {
  background: var(--storm-surface-strong);
  border: 1px solid rgba(22, 34, 43, 0.06);
  border-radius: 16px;
  padding: 0.8rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  box-shadow: 0 4px 16px rgba(18, 28, 36, 0.04);
  backdrop-filter: blur(20px);
  min-height: 0;
}

.bento-card-header {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.bento-card-header h3 {
  font-size: 0.9rem;
  font-weight: 700;
  color: var(--storm-ink);
  letter-spacing: -0.02em;
  margin: 0;
}

.bento-icon {
  width: 20px;
  height: 20px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(22, 34, 43, 0.04);
}

.cyan-glow { box-shadow: 0 0 12px rgba(13, 159, 184, 0.2); color: var(--storm-cyan); }
.orange-glow { box-shadow: 0 0 12px rgba(239, 138, 41, 0.2); color: var(--storm-orange); }
.slate-glow { box-shadow: 0 0 12px rgba(63, 85, 96, 0.2); color: var(--storm-slate); }

.bento-desc {
  font-size: 0.72rem;
  line-height: 1.25;
  color: var(--storm-muted);
  margin: 0;
}

/* -- Auth Visual Grid -- */
.jwt-visual-grid {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.55rem;
}

.jwt-row {
  display: grid;
  grid-template-columns: 60px 1fr;
  gap: 0.4rem;
  background: #11181c;
  border-radius: 6px;
  padding: 0.4rem;
  box-shadow: inset 0 2px 8px rgba(0,0,0,0.2);
  align-items: center;
}

.jwt-raw-part {
  word-break: break-all;
  font-weight: 700;
  letter-spacing: 0.02em;
  opacity: 0.9;
  border-right: 1px dashed rgba(255,255,255,0.15);
  padding-right: 0.3rem;
  font-size: 0.5rem;
}

.c-red-zone { color: #fb7185; }
.c-purple-zone { color: #c084fc; }
.c-blue-zone { color: #38bdf8; }

.jwt-json-part {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  color: #e2e8f0;
}

.c-red-zone .jwt-json-part { color: #fb7185; }
.c-purple-zone .jwt-json-part { color: #c084fc; }
.c-blue-zone .jwt-json-part { color: #38bdf8; }

.panel-head {
  font-family: var(--storm-font);
  font-size: 0.6rem;
  font-weight: 700;
  text-transform: uppercase;
  opacity: 0.8;
}

.panel-body {
  font-weight: 600;
}

.panel-code {
  margin: 0;
  background: rgba(0,0,0,0.2);
  padding: 0.2rem 0.3rem;
  border-radius: 4px;
  border: 1px solid rgba(255,255,255,0.05);
  color: #e2e8f0;
}

.cookie-tag {
  margin-top: 0.1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(13, 159, 184, 0.08);
  border: 1px solid rgba(13, 159, 184, 0.2);
  border-radius: 6px;
  padding: 0.3rem 0.5rem;
}

.tag-label { font-size: 0.65rem; font-weight: 600; color: var(--storm-ink); }
.tag-value { font-family: "IBM Plex Mono", monospace; font-size: 0.6rem; color: var(--storm-cyan); font-weight: 700;}

/* -- Real Time Visual -- */
.rt-content {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
}

.rt-visual {
  display: flex;
  align-items: center;
  background: rgba(22, 34, 43, 0.03);
  border-radius: 8px;
  padding: 0.5rem;
  margin-top: 0.5rem;
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.65rem;
}

.rt-user {
  background: white;
  border: 1px solid rgba(22, 34, 43, 0.1);
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-weight: 600;
  box-shadow: 0 2px 4px rgba(0,0,0,0.02);
  display: flex;
  align-items: center;
  gap: 0.2rem;
}

.ws-tag { color: var(--storm-muted); font-size: 0.45rem; background: rgba(0,0,0,0.05); padding: 0.1rem 0.2rem; border-radius: 3px; }

.rt-nats {
  background: #11181c;
  color: #e2e8f0;
  padding: 0.3rem 0.5rem;
  border-radius: 6px;
  text-align: center;
  border: 1px solid var(--storm-cyan);
  box-shadow: 0 0 8px rgba(13, 159, 184, 0.2);
  font-size: 0.55rem;
  font-weight: 600;
}

.rt-nats code { color: var(--storm-cyan); display: block; margin-top: 0.1rem; font-size: 0.5rem; }

.rt-arrow {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-grow: 1;
  padding: 0 0.2rem;
}

.rt-pub {
  font-size: 0.45rem;
  font-weight: 700;
  text-transform: uppercase;
  margin-bottom: 3px;
}

.rt-line {
  flex-grow: 1;
  width: 100%;
  height: 2px;
  background: rgba(22, 34, 43, 0.1);
  position: relative;
  margin: 0;
}

.rt-line.dashed {
  background: transparent;
  border-top: 2px dashed rgba(22, 34, 43, 0.2);
}

.rt-dot {
  position: absolute;
  top: -3px; left: 0;
  width: 8px; height: 8px;
  background: var(--storm-orange);
  border-radius: 50%;
  box-shadow: 0 0 8px var(--storm-orange);
  animation: pingPong 2s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

@keyframes pingPong {
  0% { transform: translateX(0); }
  50% { transform: translateX(100%); left: auto; right: 0; }
  100% { transform: translateX(0); }
}

/* -- Persistence / Worker Visual -- */
.worker-visual {
  margin-top: auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(22, 34, 43, 0.02);
  border-radius: 8px;
  padding: 0.4rem;
  border: 1px solid rgba(22, 34, 43, 0.05);
}

.wv-step {
  display: flex;
  flex-direction: column;
  gap: 0.1rem;
  font-family: var(--storm-font);
  align-items: center;
}

.wv-label {
  font-size: 0.6rem;
  font-weight: 700;
  color: var(--storm-ink);
  text-align: center;
}

.wv-speed {
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.5rem;
  font-weight: 600;
}

.c-cyan { color: var(--storm-cyan); }
.c-orange { color: var(--storm-orange); }

.box-bg {
  background: white;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  border: 1px solid var(--storm-slate);
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.wv-detail {
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.55rem;
  color: var(--storm-slate);
  font-weight: 600;
  text-align: center;
}

.wv-arrow {
  color: var(--storm-muted);
  font-size: 0.6rem;
  opacity: 0.5;
}

/* -- Observability Visual -- */
.obs-visual {
  margin-top: auto;
  background: #11181c;
  border-radius: 8px;
  padding: 0.5rem;
}

.metric-name {
  color: #c9d1d9;
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.55rem;
  margin-bottom: 0.3rem;
  display: block;
}

.sparkline {
  display: flex;
  align-items: flex-end;
  gap: 2px;
  height: 24px;
}

.spark-bar {
  flex-grow: 1;
  background: var(--storm-cyan);
  border-radius: 2px 2px 0 0;
  opacity: 0.8;
  transition: height 0.3s ease;
}

.sparkline:hover .spark-bar {
  opacity: 1;
}

@media (max-width: 840px) {
  .bento-grid {
    grid-template-columns: 1fr;
    grid-template-rows: auto;
  }
  .bento-auth, .bento-rt, .bento-db, .bento-obs {
    grid-column: 1; grid-row: auto;
  }
}
</style>
