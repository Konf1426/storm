<script setup lang="ts">
// Next Steps & Vision Slide - Premium Design
</script>

<template>
  <div class="page-shell">
    <div class="header-minimal">
      <p class="eyebrow">16 — Perspectives et conclusion</p>
      <h2 class="section-title">Perspectives : franchir un nouveau cap</h2>
    </div>

    <div class="ns-bento">
      <!-- Short Term: Tech -->
      <article class="premium-card tech-card">
        <div class="card-glow"></div>
        <header class="card-head">
          <div class="icon-glow cyan-glow">
            <img src="https://api.iconify.design/carbon/direction-fork.svg?color=%2306b6d4" class="icon-md"/>
          </div>
          <div class="head-meta">
            <span class="tag-label tech-tag">ARCHITECTURE</span>
            <h3 class="card-title">Optimisation & Noyau</h3>
          </div>
        </header>
        <div class="card-body">
          <ul class="ns-list-modern">
             <li>
                <span class="dot"></span>
                <div class="li-content">
                  <strong>Validation JSON/Crypto</strong>
                  <p>Migration vers <code>json-iter</code> et clés asymétriques pour alléger la charge CPU CPU.</p>
                </div>
             </li>
             <li>
                <span class="dot"></span>
                <div class="li-content">
                  <strong>Service Discovery Cluster</strong>
                  <p>Détection automatique des nœuds via NATS JetStream pour une scalabilité horizontale facilitée.</p>
                </div>
             </li>
          </ul>
        </div>
      </article>

      <!-- Mid Term: Features -->
      <article class="premium-card biz-card">
        <div class="card-glow"></div>
        <header class="card-head">
          <div class="icon-glow orange-glow">
            <img src="https://api.iconify.design/carbon/chart-relationship.svg?color=%23f97316" class="icon-md"/>
          </div>
          <div class="head-meta">
            <span class="tag-label biz-tag">BUSINESS</span>
            <h3 class="card-title">Fonctionnalités Produit</h3>
          </div>
        </header>
        <div class="card-body">
          <ul class="ns-list-modern">
             <li>
                <span class="dot"></span>
                <div class="li-content">
                  <strong>Chiffrement E2E</strong>
                  <p>Confidentialité totale via clé partagée de bout-en-bout entre clients.</p>
                </div>
             </li>
             <li>
                <span class="dot"></span>
                <div class="li-content">
                  <strong>Webhooks & Intégrations</strong>
                  <p>Connecteurs sortants (Go channels) pour automatiser des workflows tiers.</p>
                </div>
             </li>
          </ul>
        </div>
      </article>

      <!-- Long Term: Future -->
      <article class="premium-card ai-card">
        <div class="card-glow"></div>
        <header class="card-head">
          <div class="icon-glow slate-glow">
            <img src="https://api.iconify.design/carbon/machine-learning-model.svg?color=%2364748b" class="icon-md"/>
          </div>
          <div class="head-meta">
            <span class="tag-label ai-tag">VISION</span>
            <h3 class="card-title">IA & Analyse Live</h3>
          </div>
        </header>
        <div class="card-body">
          <ul class="ns-list-modern">
             <li>
                <span class="dot"></span>
                <div class="li-content">
                  <strong>Sentiment Analysis</strong>
                  <p>Modèle IA léger branché sur le worker pool pour une aide à la modération temps réel.</p>
                </div>
             </li>
             <li>
                <span class="dot"></span>
                <div class="li-content">
                  <strong>Filtres Toxiques</strong>
                  <p>Garantie d'un environnement sain via filtrage automatique par apprentissage.</p>
                </div>
             </li>
          </ul>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.page-shell { display: flex; flex-direction: column; height: 100%; gap: 1rem; }
.header-minimal { margin-bottom: 0.5rem; }

.ns-bento {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.2rem;
  flex-grow: 1;
}

/* Premium Card Design */
.premium-card {
  position: relative;
  background: white; border-radius: 16px; padding: 1.4rem;
  display: flex; flex-direction: column; 
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04); 
  border: 1px solid rgba(0,0,0,0.03); 
  overflow: hidden;
}

.card-glow {
  position: absolute; top: 0; left: 0; right: 0; height: 3px;
  background: linear-gradient(90deg, transparent, rgba(0,0,0,0.1), transparent);
}

.tech-card .card-glow { background: linear-gradient(90deg, #06b6d4, #67e8f9); }
.biz-card .card-glow { background: linear-gradient(90deg, #f97316, #fdba74); }
.ai-card .card-glow { background: linear-gradient(90deg, #64748b, #94a3b8); }

.card-head { display: flex; align-items: flex-start; gap: 1rem; margin-bottom: 1.2rem; }

.icon-glow {
  width: 42px; height: 42px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
  background: #f8fafc; border: 1px solid #f1f5f9;
}
.cyan-glow { box-shadow: 0 4px 12px rgba(6, 182, 212, 0.15); }
.orange-glow { box-shadow: 0 4px 12px rgba(249, 115, 22, 0.15); }
.slate-glow { box-shadow: 0 4px 12px rgba(100, 116, 139, 0.15); }

.icon-md { width: 22px; height: 22px; }

.head-meta { display: flex; flex-direction: column; gap: 0.2rem; }
.tag-label {
  font-family: "IBM Plex Mono", monospace; font-size: 0.5rem; font-weight: 800;
  letter-spacing: 0.1em; padding: 0.1rem 0.4rem; border-radius: 4px; width: fit-content;
}
.tech-tag { background: rgba(6, 182, 212, 0.1); color: #0891b2; }
.biz-tag { background: rgba(249, 115, 22, 0.1); color: #ea580c; }
.ai-tag { background: rgba(100, 116, 139, 0.1); color: #475569; }

.card-title { font-size: 0.95rem; font-weight: 800; color: var(--storm-ink); margin: 0; letter-spacing: -0.01em; }

.card-body { display: flex; flex-direction: column; flex-grow: 1; }

.ns-list-modern { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 1rem; }
.ns-list-modern li { display: flex; gap: 0.8rem; align-items: flex-start; }

.dot { width: 6px; height: 6px; border-radius: 50%; background: #cbd5e1; margin-top: 0.45rem; flex-shrink: 0; }
.tech-card .dot { background: #06b6d4; }
.biz-card .dot { background: #f97316; }
.ai-card .dot { background: #64748b; }

.li-content strong { font-size: 0.75rem; color: var(--storm-ink); display: block; margin-bottom: 0.1rem; }
.li-content p { font-size: 0.7rem; color: #64748b; margin: 0; line-height: 1.35; }
.li-content code { font-size: 0.6rem; background: #f1f5f9; padding: 0.1rem 0.3rem; border-radius: 3px; color: #334155; }

/* Vision Quote Banner */
.vision-banner {
  background: white; border: 1px solid #f1f5f9; border-radius: 16px;
  padding: 1.2rem; display: flex; align-items: center; gap: 1.5rem; justify-content: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.02); margin-top: 0.5rem;
}
.quote-icon { width: 24px; opacity: 0.3; }
.banner-quote {
  font-size: 1.1rem; font-weight: 600; color: #475569; font-style: italic;
  margin: 0; text-align: center; line-height: 1.2;
}
</style>

