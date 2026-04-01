<script setup lang="ts">
// Minimalist Final Conclusion & Q&A Slide - High Impact & Simple
import { deckMeta } from "../data/content"
</script>

<template>
  <div class="page-shell conclusion-minimal">
    <div class="main-content">
      <!-- Project Branding -->
      <div class="brand-section">
        <h1 class="storm-logo">STORM</h1>
      </div>

      <!-- Centered Q&A Action -->
      <div class="qa-action-container">
        <div class="qa-icon-box">
          <div class="pulse-ring ring-1"></div>
          <div class="pulse-ring ring-2"></div>
          <img src="https://api.iconify.design/carbon/chat.svg?color=%2306b6d4" class="qa-icon"/>
        </div>
        
        <div class="text-group">
          <h2 class="qa-title">Questions ?</h2>
          <p class="qa-prompt">C'est à vous.</p>
        </div>
      </div>

      <!-- Footer Branding -->
      <div class="footer-meta">
        <div class="team-pills">
          <span v-for="member in deckMeta.teamMembers" :key="member" class="team-pill">{{ member }}</span>
        </div>
        <p class="thank-you">Merci de votre attention.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.conclusion-minimal {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white; /* Clean white background consistent with IntroCover */
}

.main-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 4rem;
  width: 100%;
  max-width: 800px;
}

/* Brand Section */
.brand-section {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.storm-logo {
  font-size: 5rem;
  font-weight: 900;
  letter-spacing: -0.06em;
  line-height: 1;
  color: var(--storm-ink);
  background: linear-gradient(135deg, #16222b 0%, #3f5560 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.hero-divider {
  width: 40px;
  height: 4px;
  background: var(--storm-orange);
  margin: 1rem 0;
  border-radius: 2px;
}

.hero-subtitle {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--storm-cyan);
  text-transform: uppercase;
  letter-spacing: 0.2em;
  margin: 0;
}

/* Q&A Action */
.qa-action-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
}

.qa-icon-box {
  position: relative;
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.qa-icon {
  width: 40px;
  height: 40px;
  z-index: 2;
}

.text-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.qa-title {
  font-size: 2.5rem;
  font-weight: 900;
  color: var(--storm-ink);
  margin: 0;
  letter-spacing: -0.02em;
}

.qa-prompt {
  font-size: 1.1rem;
  font-weight: 600;
  color: #64748b;
  margin: 0;
}

/* Animation */
.pulse-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 2px solid var(--storm-cyan);
  opacity: 0;
}

.ring-1 {
  animation: pulse 3s infinite cubic-bezier(0.215, 0.61, 0.355, 1);
}

.ring-2 {
  animation: pulse 3s infinite cubic-bezier(0.215, 0.61, 0.355, 1) 1.5s;
}

@keyframes pulse {
  0% { transform: scale(0.6); opacity: 0; }
  20% { opacity: 0.4; }
  100% { transform: scale(1.8); opacity: 0; }
}

/* Footer */
.footer-meta {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.team-pills {
  display: flex;
  gap: 0.5rem;
}

.team-pill {
  font-size: 0.6rem;
  font-weight: 700;
  color: #94a3b8;
  background: #f8fafc;
  padding: 0.15rem 0.5rem;
  border-radius: 999px;
  border: 1px solid #f1f5f9;
}

.thank-you {
  font-size: 0.75rem;
  font-weight: 600;
  color: #cbd5e1;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
</style>