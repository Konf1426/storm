<script setup lang="ts">
// Final Conclusion & Q&A Slide - High Impact Design
</script>

<template>
  <div class="page-shell conclusion-shell">
    <div class="header-minimal">
      <p class="eyebrow">14 — Conclusion Finale</p>
      <h2 class="section-title">STORM : Un socle robuste et mesurable</h2>
    </div>

    <div class="main-statement shadow-xl">
       <div class="statement-glow"></div>
       <div class="quote-mark">“</div>
       <p class="fs-line">
         STORM n'est pas qu'une démonstration technique — c'est un système instrumenté, testé et documenté de façon honnête.
       </p>
       <div class="fs-footer">
          <div class="thank-you">Merci de votre attention.</div>
          <div class="qa-badge">Session de Questions — 2026</div>
       </div>
    </div>

    <div class="pillars-grid">
      <article class="pillar-mini-card">
        <div class="pillar-icon cyan-bg"><img src="https://api.iconify.design/carbon/iot-connect.svg?color=white" class="icon-xs"/></div>
        <div class="pillar-content">
          <h4>Architecture</h4>
          <p>Découplage Go/NATS/Redis réussi.</p>
        </div>
      </article>

      <article class="pillar-mini-card">
        <div class="pillar-icon green-bg"><img src="https://api.iconify.design/carbon/security-block.svg?color=white" class="icon-xs"/></div>
        <div class="pillar-content">
          <h4>Qualité</h4>
          <p>Défense en profondeur & CI/CD.</p>
        </div>
      </article>

      <article class="pillar-mini-card">
        <div class="pillar-icon orange-bg"><img src="https://api.iconify.design/carbon/area-custom.svg?color=white" class="icon-xs"/></div>
        <div class="pillar-content">
          <h4>Performance</h4>
          <p>Résilience prouvée sous charge.</p>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.conclusion-shell {
  display: flex; flex-direction: column; height: 100%; gap: 1.5rem; justify-content: center;
  padding: 0 2rem;
}

.header-minimal { margin-bottom: 0.5rem; text-align: center; }

/* Main Statement Area */
.main-statement {
  position: relative;
  background: #0f172a; color: white; border-radius: 24px; padding: 3.5rem 3rem;
  display: flex; flex-direction: column; gap: 2rem; border: 1px solid #1e293b;
  box-shadow: 0 25px 60px -12px rgba(15, 23, 42, 0.4);
  overflow: hidden;
  text-align: center;
}

.statement-glow {
  position: absolute; top: -50%; left: -50%; width: 200%; height: 200%;
  background: radial-gradient(circle at center, rgba(6, 182, 212, 0.08) 0%, transparent 40%);
  pointer-events: none;
}

.quote-mark {
  position: absolute; top: 1rem; left: 2rem;
  font-size: 8rem; font-family: serif; color: rgba(255,255,255,0.05);
  line-height: 1; pointer-events: none;
}

.fs-line {
  font-size: 1.6rem; font-weight: 700; margin: 0; line-height: 1.3; letter-spacing: -0.02em;
  position: relative; z-index: 1;
  background: linear-gradient(180deg, #ffffff 0%, #cbd5e1 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.fs-footer {
  display: flex; justify-content: space-between; align-items: center;
  padding-top: 1.8rem; border-top: 1px solid rgba(255,255,255,0.08);
  position: relative; z-index: 1;
}

.thank-you { font-size: 1.1rem; font-weight: 600; color: #06b6d4; }
.qa-badge {
  background: #f97316; color: white; padding: 0.5rem 1rem; border-radius: 8px;
  font-weight: 800; font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.05em;
  box-shadow: 0 4px 12px rgba(249, 115, 22, 0.3);
}

/* Bottom Pillars Grid */
.pillars-grid {
  display: grid; grid-template-columns: repeat(3, 1fr); gap: 1.2rem;
}

.pillar-mini-card {
  display: flex; align-items: center; gap: 1rem;
  background: white; border-radius: 12px; padding: 1rem;
  border: 1px solid #f1f5f9; box-shadow: 0 4px 16px rgba(0,0,0,0.02);
}

.pillar-icon {
  width: 32px; height: 32px; border-radius: 8px;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.icon-xs { width: 14px; height: 14px; }

.pillar-content h4 { font-size: 0.75rem; font-weight: 800; color: #1e293b; margin: 0; text-transform: uppercase; letter-spacing: 0.02em; }
.pillar-content p { font-size: 0.65rem; color: #64748b; margin: 0; font-weight: 500; }

.cyan-bg { background: #06b6d4; }
.green-bg { background: #10b981; }
.orange-bg { background: #f97316; }
</style>
