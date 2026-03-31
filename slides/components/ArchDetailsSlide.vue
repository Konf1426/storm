<script setup lang="ts">
// Vue component for Architecture Detailed Grid
</script>

<template>
  <div class="page-shell details-shell">
    <div class="header-minimal">
      <p class="eyebrow">03 — Détails techniques</p>
      <h2 class="section-title">Anatomie du pipeline temps réel</h2>
    </div>

    <div class="clean-grid">
      <!-- Auth -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box cyan-bg"><img src="/icons/vue.svg" class="icon-sm"/></div>
          <h3 class="card-title">Sécurité & JWT</h3>
        </header>
        <div class="card-body">
          <p class="text-sm">Token d'accès court (15m) + Refresh Token en Cookie <strong>HttpOnly</strong> pour bloquer le vol par failles XSS.</p>
          <div class="jwt-mini">
            <div class="jwt-row"><span class="badge bg-red">En-tête</span> <code class="mono-text">{ "alg": "HS256" }</code></div>
            <div class="jwt-row"><span class="badge bg-purple">Payload</span> <code class="mono-text">{ "sub": "usr_9x2", "nosnif": true, "exp": 17... }</code></div>
            <div class="jwt-row"><span class="badge bg-blue">Signature</span> <code class="mono-text">Générée par le serveur (HMAC)</code></div>
          </div>
        </div>
      </article>

      <!-- Temps Réel -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box orange-bg"><img src="/icons/nats.png" class="icon-sm"/></div>
          <h3 class="card-title">Routage (Pub/Sub)</h3>
        </header>
        <div class="card-body">
          <p class="text-sm">Le Broker <strong>NATS</strong> distribue instantanément les messages entre tous les serveurs Gateways et les clients.</p>
          <div class="flow-mini">
            <div class="flow-pill">Alice<br/><span class="sub text-muted">WebSocket</span></div>
            <div class="flow-arr text-orange">Publie ➔</div>
            <div class="flow-broker dark-bg">
              NATS <br/> <span>room.142</span>
            </div>
            <div class="flow-arr text-cyan">➔ Reçoit</div>
            <div class="flow-pill">Bob<br/><span class="sub text-muted">WebSocket</span></div>
          </div>
        </div>
      </article>

      <!-- Worker Pool -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box slate-bg"><img src="/icons/postgres.svg" class="icon-sm"/></div>
          <h3 class="card-title">Worker Pool</h3>
        </header>
        <div class="card-body">
          <p class="text-sm">L'écriture lente dans PostgreSQL est détachée des WebSockets grâce à un Channel en RAM et des Goroutines.</p>
          <div class="flow-mini">
             <div class="flow-pill border-cyan">WebSocket<br/><span class="speed">< 1ms</span></div>
             <div class="flow-arr">➔</div>
             <div class="flow-broker outline">File Go<br/><span class="sub text-slate">(RAM)</span></div>
             <div class="flow-arr">➔</div>
             <div class="flow-pill border-orange">Workers ➔ DB<br/><span class="speed">> 10ms</span></div>
          </div>
        </div>
      </article>

      <!-- Observabilité -->
      <article class="clean-card shadow-sm">
        <header class="card-head">
          <div class="icon-box cyan-bg"><img src="/icons/prometheus.svg" class="icon-sm"/></div>
          <h3 class="card-title">Observabilité</h3>
        </header>
        <div class="card-body">
          <p class="text-sm">Métriques natives exposées sur <code>/metrics</code> (Prometheus) et profilage mémoire (pprof) intégrés au binaire Go.</p>
          <div class="metrics-mini">
             <div class="spark-bars">
               <div class="bar h-40"></div><div class="bar h-60"></div><div class="bar h-30"></div><div class="bar h-80"></div><div class="bar h-100 peak"></div>
             </div>
             <code class="metric-text">storm_active_websockets</code>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
/* Base Layout */
.details-shell {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.header-minimal {
  margin-bottom: 1.2rem;
}

.clean-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 1.5rem;
  flex-grow: 1;
  min-height: 0;
}

/* Cards */
.clean-card {
  background: white;
  border-radius: 12px;
  padding: 1.4rem;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0,0,0,0.04);
}

.card-head {
  display: flex;
  align-items: center;
  gap: 0.8rem;
  margin-bottom: 0.6rem;
}

.card-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--storm-ink);
  margin: 0;
}

.icon-box {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.icon-sm { width: 16px; filter: grayscale(1) brightness(0) invert(1); }

.cyan-bg { background: var(--storm-cyan); box-shadow: 0 4px 12px rgba(13, 159, 184, 0.2); }
.orange-bg { background: var(--storm-orange); box-shadow: 0 4px 12px rgba(239, 138, 41, 0.2); }
.slate-bg { background: var(--storm-slate); box-shadow: 0 4px 12px rgba(63, 85, 96, 0.2); }

.card-body {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  gap: 0.8rem;
}

.text-sm {
  font-size: 0.8rem;
  line-height: 1.4;
  color: #475569;
  margin: 0;
}

/* JWT Mini */
.jwt-mini {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  background: #f8fafc;
  padding: 0.8rem;
  border-radius: 8px;
  border: 1px solid #f1f5f9;
}

.jwt-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.badge {
  font-size: 0.55rem;
  font-family: var(--storm-font);
  font-weight: 700;
  text-transform: uppercase;
  padding: 0.2rem;
  border-radius: 4px;
  color: white;
  width: 65px;
  text-align: center;
}

.bg-red { background: #fb7185; }
.bg-purple { background: #c084fc; }
.bg-blue { background: #38bdf8; }

.mono-text {
  font-family: inherit;
  font-size: 0.65rem;
  color: #334155;
  background: white;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  border: 1px solid #e2e8f0;
  flex-grow: 1;
}

/* Flow Visuals (NATS & Worker Pool) */
.flow-mini {
  margin-top: auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #f8fafc;
  padding: 0.8rem 0.6rem;
  border-radius: 8px;
  border: 1px solid #f1f5f9;
}

.flow-pill {
  background: white;
  padding: 0.4rem 0.5rem;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  color: var(--storm-ink);
  border: 1px solid #e2e8f0;
  text-align: center;
  box-shadow: 0 1px 2px rgba(0,0,0,0.02);
  line-height: 1.2;
}

.flow-arr {
  font-size: 0.55rem;
  font-weight: 700;
  color: var(--storm-muted);
  text-align: center;
  text-transform: uppercase;
}
.text-orange { color: #d97706; }
.text-cyan { color: #0891b2; }

.flow-broker {
  background: white;
  padding: 0.4rem 0.8rem;
  border-radius: 6px;
  font-size: 0.7rem;
  font-weight: 700;
  text-align: center;
  line-height: 1.3;
}

.flow-broker.dark-bg {
  background: #0f172a;
  color: white;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.15);
}
.flow-broker.dark-bg span { color: var(--storm-cyan); font-size: 0.55rem; font-family: monospace; }
.flow-broker.outline { border: 1px dashed #94a3b8; color: #475569; }

.border-cyan { border-left: 3px solid var(--storm-cyan); }
.border-orange { border-left: 3px solid var(--storm-orange); }

.sub { font-size: 0.5rem; font-weight: 600; text-transform: uppercase; }
.text-muted { color: #94a3b8; }
.text-slate { color: #64748b; }
.speed { display: block; font-family: monospace; font-size: 0.6rem; color: #64748b; margin-top: 0.15rem; }

/* Metrics Mini */
.metrics-mini {
  margin-top: auto;
  background: #f8fafc;
  padding: 0.8rem;
  border-radius: 8px;
  border: 1px solid #f1f5f9;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.4rem;
}

.spark-bars {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  height: 32px;
  width: 100%;
}

.bar {
  flex: 1;
  background: var(--storm-cyan);
  border-radius: 2px 2px 0 0;
  opacity: 0.6;
}

.bar.peak { background: var(--storm-orange); opacity: 1; }
.h-40 { height: 40%; }
.h-60 { height: 60%; }
.h-30 { height: 30%; }
.h-80 { height: 80%; }
.h-100 { height: 100%; }

.metric-text {
  font-size: 0.65rem;
  color: #475569;
  font-family: monospace;
}
</style>
