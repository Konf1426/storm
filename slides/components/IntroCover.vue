<script setup lang="ts">
import { deckMeta, objectives } from "../data/content"
</script>

<template>
  <div class="page-shell intro-shell">
    <div class="intro-main intro-main-full">
      <div class="intro-topline">
        <p class="eyebrow">MT4 — Soutenance de projet</p>
        <div class="intro-teamline">
          <span v-for="member in deckMeta.teamMembers" :key="member" class="intro-team-member">
            {{ member }}
          </span>
        </div>
      </div>

      <h1 class="hero-title">{{ deckMeta.title }}</h1>
      <p class="hero-subtitle">{{ deckMeta.subtitle }}</p>
      <p class="intro-lead">{{ deckMeta.positioning }}</p>

      <div class="hero-tags">
        <span class="hero-tag">Go</span>
        <span class="hero-tag">WebSocket</span>
        <span class="hero-tag">NATS</span>
        <span class="hero-tag">Postgres</span>
        <span class="hero-tag">Redis</span>
        <span class="hero-tag">Prometheus</span>
        <span class="hero-tag">Grafana</span>
      </div>

      <div class="intro-targets">
        <article v-for="objective in objectives" :key="objective.label" class="intro-target">
          <p class="intro-target-label">{{ objective.label }}</p>
          <p class="intro-target-value">{{ objective.value }}</p>
          <p class="intro-target-note">{{ objective.note }}</p>
        </article>
      </div>
    </div>
  </div>
</template>
