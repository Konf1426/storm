<script setup lang="ts">
// Cloud Analysis & Blockers Slide - Premium Design
import { azureLoadTest } from "../data/content"
</script>

<template>
  <div class="page-shell">
    <div class="header-minimal">
      <p class="eyebrow">11 — Analyse Cloud : Étape 1</p>
      <h2 class="section-title">Audit Intermédiaire : Identifier les "Faux Plafonds"</h2>
    </div>

    <div class="content-split">
      <div class="analysis-diag">
        <div class="bento-diag">
          <!-- Blocker: MC_ -->
          <article class="diag-card orange-glow slide-in delay-1">
            <header class="card-head">
              <div class="icon-indicator"><img src="https://api.iconify.design/carbon/warning-alt-filled.svg?color=%23f97316" class="icon-sm"/></div>
              <h3 class="card-title">Groupe MC_</h3>
            </header>
            <p class="diag-text">Resource group AKS non réconcilié — corrigé via CLI <code>az update</code>.</p>
          </article>

          <!-- Blocker: ACR 401 -->
          <article class="diag-card orange-glow slide-in delay-2">
            <header class="card-head">
              <div class="icon-indicator"><img src="https://api.iconify.design/carbon/locked.svg?color=%23f97316" class="icon-sm"/></div>
              <h3 class="card-title">Pull ACR (401)</h3>
            </header>
            <p class="diag-text">Permission d'accès aux images Docker manquante sur AKS.</p>
          </article>

          <!-- Blocker: Rate Limit -->
          <article class="diag-card slate-glow slide-in delay-3">
            <header class="card-head">
              <div class="icon-indicator"><img src="https://api.iconify.design/carbon/direction-sharp-turn.svg?color=%2364748b" class="icon-sm"/></div>
              <h3 class="card-title">Rate Limiting</h3>
            </header>
            <p class="diag-text">Le garde-fou <code>5 req/min/IP</code> bloquait 99 % des tests de charge initiaux.</p>
          </article>

          <!-- Blocker: Crypto -->
          <article class="diag-card slate-glow slide-in delay-4">
            <header class="card-head">
              <div class="icon-indicator"><img src="https://api.iconify.design/carbon/chip.svg?color=%2364748b" class="icon-sm"/></div>
              <h3 class="card-title">Coût Bcrypt</h3>
            </header>
            <p class="diag-text">Saturations CPU massives lors des phases d'inscriptions de masse.</p>
          </article>
        </div>
      </div>

      <div class="meta-section">
        <div class="infra-panel shadow-sm">
          <h4 class="panel-title">Infrastructure Initiale</h4>
          <ul class="infra-list">
             <li v-for="item in azureLoadTest.initialInfra" :key="item">{{ item }}</li>
          </ul>
        </div>

        <div class="quote-box">
          <p class="quote-text">
            "Le premier run ne mesure pas <strong>STORM</strong>, il révèle uniquement les limites de configuration de l'environnement Azure."
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-shell { display: flex; flex-direction: column; height: 100%; gap: 1rem; }
.header-minimal { margin-bottom: 0.5rem; }

.content-split {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 1.5rem;
  flex-grow: 1;
}

/* Bento Grid Analysis */
.bento-diag {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  height: 100%;
}

.diag-card {
  background: white; border-radius: 16px; padding: 1rem;
  display: flex; flex-direction: column; gap: 0.6rem;
  border: 1px solid rgba(0,0,0,0.03);
  box-shadow: 0 4px 16px rgba(0,0,0,0.02);
  position: relative;
  transition: transform 0.3s ease;
}

.diag-card:hover { transform: translateY(-3px); }

.orange-glow { border-top: 3px solid #f97316; box-shadow: 0 4px 20px rgba(249, 115, 22, 0.05); }
.slate-glow { border-top: 3px solid #64748b; box-shadow: 0 4px 20px rgba(100, 116, 139, 0.05); }

.card-head { display: flex; align-items: center; gap: 0.6rem; }
.icon-indicator { width: 30px; height: 30px; border-radius: 8px; background: #f8fafc; display: flex; align-items: center; justify-content: center; }
.card-title { font-size: 0.85rem; font-weight: 800; color: var(--storm-ink); margin: 0; }
.diag-text { font-size: 0.7rem; color: #64748b; line-height: 1.4; margin: 0; }
.diag-text code { font-size: 0.65rem; background: #f1f5f9; padding: 0.1rem 0.2rem; border-radius: 3px; }

/* Infra panel */
.infra-panel { background: #f8fafc; border: 1px solid #f1f5f9; border-radius: 16px; padding: 1.2rem; display: flex; flex-direction: column; gap: 0.8rem; }
.panel-title { font-size: 0.75rem; font-weight: 800; color: #1e293b; text-transform: uppercase; letter-spacing: 0.04em; margin: 0; }
.infra-list { padding: 0; margin: 0; list-style: none; display: flex; flex-direction: column; gap: 0.4rem; }
.infra-list li { font-size: 0.75rem; color: #475569; position: relative; padding-left: 1rem; }
.infra-list li::before { content: "➔"; position: absolute; left: 0; color: var(--storm-cyan); }

.quote-box { margin-top: auto; padding-top: 1rem; border-top: 1px solid #f1f5f9; }
.quote-text { font-size: 0.95rem; font-weight: 600; color: #475569; font-style: italic; line-height: 1.4; opacity: 0.85; text-align: center; }

/* Animations */
.slide-in { opacity: 0; transform: translateY(15px); animation: slide 0.5s ease forwards; }
@keyframes slide { to { opacity: 1; transform: translateY(0); } }
.delay-1 { animation-delay: 0.1s; }
.delay-2 { animation-delay: 0.2s; }
.delay-3 { animation-delay: 0.3s; }
.delay-4 { animation-delay: 0.4s; }
</style>
