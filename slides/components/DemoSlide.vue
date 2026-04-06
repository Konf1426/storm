<script setup lang="ts">
// Slide – Live Demo backdrop. Shown while the presenter demos the app live.
</script>

<template>
  <div class="demo-backdrop">
    <!-- Pulsing rings -->
    <div class="ring ring-1"></div>
    <div class="ring ring-2"></div>
    <div class="ring ring-3"></div>

    <!-- Center content -->
    <div class="demo-center">
      <div class="demo-icon-wrap">
        <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <polygon points="5 3 19 12 5 21 5 3"/>
        </svg>
      </div>

      <h1 class="demo-title">DÉMO LIVE</h1>
      <div class="demo-divider"></div>
      <p class="demo-sub">Application STORM — Messagerie temps réel</p>
    </div>

    <!-- Bottom indicator -->
    <div class="demo-indicator">
      <span class="indicator-dot"></span>
      <span class="indicator-text">Démonstration en cours</span>
    </div>
  </div>
</template>

<style scoped>
.demo-backdrop {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

/* ── Pulsing concentric rings ── */
.ring {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  border: 1px solid rgba(13, 159, 184, 0.1);
  animation: ringPulse 4s ease-in-out infinite;
}

.ring-1 {
  width: 220px;
  height: 220px;
  border-color: rgba(13, 159, 184, 0.14);
}

.ring-2 {
  width: 360px;
  height: 360px;
  animation-delay: 0.8s;
  border-color: rgba(13, 159, 184, 0.08);
}

.ring-3 {
  width: 520px;
  height: 520px;
  animation-delay: 1.6s;
  border-color: rgba(13, 159, 184, 0.04);
}

@keyframes ringPulse {
  0%, 100% { transform: translate(-50%, -50%) scale(0.96); opacity: 0.3; }
  50% { transform: translate(-50%, -50%) scale(1.04); opacity: 1; }
}

/* ── Center block ── */
.demo-center {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.8rem;
}

.demo-icon-wrap {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: rgba(13, 159, 184, 0.1);
  border: 1px solid rgba(13, 159, 184, 0.18);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--storm-cyan);
  animation: iconGlow 3s ease-in-out infinite;
}

@keyframes iconGlow {
  0%, 100% { box-shadow: 0 0 16px rgba(13, 159, 184, 0.06); }
  50% { box-shadow: 0 0 32px rgba(13, 159, 184, 0.18); }
}

.demo-title {
  font-size: 3.6rem;
  font-weight: 900;
  letter-spacing: 0.12em;
  color: var(--storm-ink);
  margin: 0;
  background: linear-gradient(135deg, #16222b 0%, #3f5560 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.demo-divider {
  width: 60px;
  height: 4px;
  background: var(--storm-orange);
  border-radius: 2px;
}

.demo-sub {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--storm-muted);
  letter-spacing: 0.04em;
  margin: 0;
}

/* ── Bottom indicator ── */
.demo-indicator {
  position: absolute;
  bottom: 1.8rem;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 0.55rem;
  border-radius: 999px;
  border: 1px solid rgba(22, 34, 43, 0.08);
  background: rgba(255, 255, 255, 0.65);
  padding: 0.35rem 0.85rem;
  box-shadow: 0 10px 30px rgba(18, 28, 36, 0.06);
}

.indicator-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--storm-cyan);
  animation: blink 2s ease-in-out infinite;
}

@keyframes blink {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 1; }
}

.indicator-text {
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.68rem;
  font-weight: 600;
  color: var(--storm-muted);
  letter-spacing: 0.08em;
  text-transform: uppercase;
}
</style>
