<script setup lang="ts">
// Animated API & Message Flow slide
</script>

<template>
  <div class="page-shell api-flow-shell">
    <div class="header-minimal">
      <p class="eyebrow">06 — Flux API & Temps réel</p>
      <h2 class="section-title">Visualisation du pipeline : de l'Auth au WebSocket</h2>
    </div>

    <div class="api-grid">
      <!-- Section Auth -->
      <section class="api-section">
        <h3 class="section-subtitle"><span class="icon">🔐</span> Cycle d'Authentification</h3>
        <div class="flow-container">
          <div class="actor">Client</div>
          <div class="lanes">
            <!-- Login -->
            <div class="lane">
              <div class="method post">POST</div>
              <div class="path">/auth/login</div>
              <div class="connection">
                <div class="dot bg-cyan"></div>
              </div>
              <div class="result">Cookie <span class="tag">HttpOnly</span></div>
            </div>
            <!-- Refresh -->
            <div class="lane">
              <div class="method post">POST</div>
              <div class="path">/auth/refresh</div>
              <div class="connection">
                <div class="dot bg-cyan"></div>
              </div>
              <div class="result">Session prolongée</div>
            </div>
            <!-- Me -->
            <div class="lane">
              <div class="method get">GET</div>
              <div class="path">/auth/me</div>
              <div class="connection">
                <div class="dot bg-cyan"></div>
              </div>
              <div class="result">Profil utilisateur</div>
            </div>
          </div>
          <div class="actor">Gateway</div>
        </div>
      </section>

      <!-- Section Real-time -->
      <section class="api-section">
        <h3 class="section-subtitle"><span class="icon">⚡</span> Flux Temps Réel</h3>
        <div class="flow-container">
          <div class="actor">Client</div>
          <div class="lanes">
            <!-- Websocket -->
            <div class="lane">
              <div class="method get">GET</div>
              <div class="path">/ws</div>
              <div class="connection persistent">
                <div class="dot bg-orange pulse"></div>
              </div>
              <div class="result">Canal Full-Duplex</div>
            </div>
            <!-- Publish -->
            <div class="lane">
              <div class="method post">POST</div>
              <div class="path">/publish</div>
              <div class="connection">
                <div class="dot bg-orange"></div>
              </div>
              <div class="result">Injection NATS</div>
            </div>
            <!-- Metrics -->
            <div class="lane">
              <div class="method get">GET</div>
              <div class="path">/healthz</div>
              <div class="connection dashed">
                <div class="dot bg-slate"></div>
              </div>
              <div class="result">Status OK</div>
            </div>
          </div>
          <div class="actor">Gateway</div>
        </div>
      </section>
    </div>

    <div class="api-footer">
      <p class="footer-note">
        <strong>Stateless :</strong> Le Gateway ne stocke pas de session en RAM. Tout repose sur la validité du <strong>JWT</strong> et du <strong>Refresh Token</strong>.
      </p>
    </div>
  </div>
</template>

<style scoped>
.api-flow-shell {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.header-minimal {
  margin-bottom: 1.5rem;
}

.section-subtitle {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--storm-slate);
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.section-subtitle .icon {
  font-size: 1.1rem;
}

.api-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  flex-grow: 1;
  min-height: 0;
}

.api-section {
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.4);
  border-radius: 16px;
  padding: 1.2rem;
  border: 1px solid rgba(0, 0, 0, 0.03);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02);
}

.flow-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-grow: 1;
  padding: 1rem 0;
}

.actor {
  background: var(--storm-ink);
  color: white;
  padding: 0.5rem 0.8rem;
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 700;
  writing-mode: vertical-lr;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  transform: rotate(180deg);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.lanes {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
  flex-grow: 1;
  padding: 0 1rem;
}

.lane {
  display: grid;
  grid-template-columns: 45px 80px 1fr 100px;
  align-items: center;
  gap: 0.6rem;
}

.method {
  font-family: monospace;
  font-size: 0.65rem;
  font-weight: 800;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  text-align: center;
}

.method.post { background: #fee2e2; color: #ef4444; }
.method.get { background: #dcfce7; color: #10b981; }

.path {
  font-family: monospace;
  font-size: 0.65rem;
  color: var(--storm-slate);
  font-weight: 600;
}

.connection {
  position: relative;
  height: 2px;
  background: rgba(0, 0, 0, 0.08);
  margin: 0 0.5rem;
}

.connection.persistent {
  background: var(--storm-orange);
  opacity: 0.4;
  height: 3px;
}

.connection.dashed {
  background: transparent;
  border-bottom: 1px dashed rgba(0, 0, 0, 0.2);
  height: 0;
}

.dot {
  position: absolute;
  top: 50%;
  left: 0;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  transform: translateY(-50%);
  animation: moveDot 2s infinite cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes moveDot {
  0% { left: 0%; opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { left: 100%; opacity: 0; }
}

.dot.pulse {
  animation: moveDot 2s infinite cubic-bezier(0.4, 0, 0.2, 1), breathe 2s infinite ease-in-out;
}

@keyframes breathe {
  0%, 100% { transform: translateY(-50%) scale(1); box-shadow: 0 0 0 rgba(239, 138, 41, 0); }
  50% { transform: translateY(-50%) scale(1.5); box-shadow: 0 0 8px rgba(239, 138, 41, 0.6); }
}

.bg-cyan { background: var(--storm-cyan); box-shadow: 0 0 4px var(--storm-cyan); }
.bg-orange { background: var(--storm-orange); box-shadow: 0 0 4px var(--storm-orange); }
.bg-slate { background: var(--storm-slate); }

.result {
  font-size: 0.65rem;
  color: var(--storm-muted);
  font-weight: 500;
  white-space: nowrap;
}

.tag {
  background: rgba(13, 159, 184, 0.1);
  color: var(--storm-cyan);
  padding: 0 0.2rem;
  border-radius: 3px;
  font-weight: 700;
}

.api-footer {
  margin-top: 1.5rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.footer-note {
  font-size: 0.8rem;
  color: var(--storm-slate);
  margin: 0;
  text-align: center;
}

.footer-note strong { color: var(--storm-ink); }
</style>
