<script setup lang="ts">
import SourceStrip from "./SourceStrip.vue"
import { azureLoadTest, sources } from "../data/content"

const stage1000 = azureLoadTest.stages[0]
const stage5000 = azureLoadTest.stages[1]
const stage10000 = azureLoadTest.stages[2]
</script>

<template>
  <div class="page-shell">
    <span class="backup-tag">Backup</span>
    <p class="eyebrow">Azure Load Test - Detail 1k -> 10k VUs</p>
    <h2 class="section-title">Lecture detaillee de la campagne du 19 mars 2026</h2>

    <div class="table-card">
      <div class="results-table">
        <div class="th metric-col">Metrique</div>
        <div class="th">1 000 VUs</div>
        <div class="th">5 000 VUs</div>
        <div class="th">10 000 VUs</div>

        <div class="td metric-col">Succes HTTP</div>
        <div class="td">{{ stage1000.httpSuccess }}</div>
        <div class="td">{{ stage5000.httpSuccess }}</div>
        <div class="td strong">{{ stage10000.httpSuccess }}</div>

        <div class="td metric-col">Latence msg mediane</div>
        <div class="td">{{ stage1000.messageMedian }}</div>
        <div class="td">{{ stage5000.messageMedian }}</div>
        <div class="td strong">{{ stage10000.messageMedian }}</div>

        <div class="td metric-col">Login median</div>
        <div class="td">{{ stage1000.loginMedian }}</div>
        <div class="td">{{ stage5000.loginMedian }}</div>
        <div class="td strong">{{ stage10000.loginMedian }}</div>

        <div class="td metric-col">Login moyen</div>
        <div class="td">{{ stage1000.loginAverage }}</div>
        <div class="td">{{ stage5000.loginAverage }}</div>
        <div class="td strong">{{ stage10000.loginAverage }}</div>

        <div class="td metric-col">Debit WebSocket</div>
        <div class="td">{{ stage1000.websocketThroughput }}</div>
        <div class="td">{{ stage5000.websocketThroughput }}</div>
        <div class="td strong">{{ stage10000.websocketThroughput }}</div>

        <div class="td metric-col">Blocage dominant</div>
        <div class="td">{{ stage1000.dominantBlocker }}</div>
        <div class="td">{{ stage5000.dominantBlocker }}</div>
        <div class="td strong">{{ stage10000.dominantBlocker }}</div>
      </div>
    </div>

    <div class="backup-grid">
      <article class="clean-card">
        <p class="card-label">Bench settings only</p>
        <ul class="note-list">
          <li v-for="item in azureLoadTest.caveats" :key="item">{{ item }}</li>
        </ul>
      </article>

      <article class="clean-card accent-card">
        <p class="card-label">Lecture cout</p>
        <div class="cost-stack">
          <div class="cost-row">
            <span>Campagne ultra</span>
            <strong>{{ azureLoadTest.costs.stressTestBurst }}</strong>
          </div>
          <div class="cost-row">
            <span>Mise en veille</span>
            <strong>{{ azureLoadTest.costs.standby }}</strong>
          </div>
          <div class="cost-row">
            <span>Part de PostgreSQL</span>
            <strong>{{ azureLoadTest.costs.postgresShare }}</strong>
          </div>
        </div>
      </article>
    </div>

    <SourceStrip :items="sources.azureLoadTest" />
  </div>
</template>

<style scoped>
.page-shell {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.section-title {
  font-size: 2.12rem;
  line-height: 0.98;
}

.table-card,
.clean-card {
  background: white;
  border-radius: 14px;
  border: 1px solid rgba(0, 0, 0, 0.04);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.03);
}

.table-card {
  margin-top: 0.18rem;
  overflow: hidden;
}

.results-table {
  display: grid;
  grid-template-columns: 1.18fr 0.94fr 0.94fr 0.94fr;
}

.th,
.td {
  padding: 0.3rem 0.42rem;
  border-bottom: 1px solid #eef2f7;
  border-right: 1px solid #eef2f7;
  font-size: 0.58rem;
  line-height: 1.18;
}

.results-table > :nth-child(4n) {
  border-right: 0;
}

.th {
  background: #f8fafc;
  color: #475569;
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.58rem;
  font-weight: 800;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.metric-col {
  color: var(--storm-ink);
  font-weight: 700;
}

.td {
  color: #64748b;
}

.td.strong {
  color: var(--storm-ink);
  font-weight: 700;
}

.backup-grid {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 0.78rem;
  margin-top: 0.58rem;
}

.clean-card {
  padding: 0.62rem 0.72rem;
}

.accent-card {
  background: linear-gradient(180deg, rgba(239, 138, 41, 0.08), rgba(255, 255, 255, 0.96));
}

.card-label {
  margin: 0 0 0.34rem;
  color: #64748b;
  font-family: "IBM Plex Mono", monospace;
  font-size: 0.6rem;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.note-list {
  margin: 0;
  padding-left: 0.9rem;
  color: #475569;
  font-size: 0.58rem;
  line-height: 1.22;
}

.note-list li + li {
  margin-top: 0.22rem;
}

.cost-stack {
  display: flex;
  flex-direction: column;
  gap: 0.22rem;
}

.cost-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 0.8rem;
  border-bottom: 1px dashed #e2e8f0;
  padding-bottom: 0.18rem;
  color: #475569;
  font-size: 0.6rem;
}

.cost-row:last-child {
  border-bottom: 0;
  padding-bottom: 0;
}

.cost-row strong {
  color: var(--storm-ink);
  font-size: 0.68rem;
}

@media (max-width: 720px) {
  .results-table,
  .backup-grid {
    grid-template-columns: 1fr;
  }

  .results-table {
    display: block;
  }
}
</style>
