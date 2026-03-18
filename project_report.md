# 🧑‍💻 Mes Contributions — Projet STORM

> Résumé de mes commits personnels sur le projet de plateforme de messagerie temps-réel (100k users).

---

## 📐 Architecture & Documentation
- Diagramme d'architecture complet du système avec assets visuels (Go, Docker, NATS, Grafana, k6, Vue...)
- Ajout des consignes et du cahier des charges au repo

## 🔐 Authentification & Sécurité
- **Rate limiting IP-based** sur `/auth/*` — 5 req/min, burst de 3 — via `golang.org/x/time/rate`
- Toggle `AUTH_RATE_LIMIT_ENABLED` pour désactiver pendant les tests de charge k6
- **Auth JWT complète** : access token (15min) + refresh token (24h) en cookies HttpOnly

## ♻️ Refactoring & Qualité
- Suppression du legacy SSE au profit du **WebSocket exclusif** (~80 lignes de code mort supprimées)
- Plan de migration documenté (`plans/sse-to-websockets.md`)

## 📊 Observabilité
- Export **Prometheus** depuis le Gateway Go :
  - `storm_active_websockets`, `storm_auth_rate_limited_total`, `storm_nats_publish_duration_seconds`, `storm_save_queue_length`
- **Provisionnement automatique Grafana** (datasource + dashboard "STORM Overview") via `provisioning/`

## ⚡ Worker Pool V2 (Performance)
- **Persistance asynchrone** : messages et refresh tokens écrits en arrière-plan, sans bloquer les handlers HTTP
- Config dynamique via env vars (`WORKER_POOL_SIZE`, `BCRYPT_COST`, `POSTGRES_POOL_MAX_CONNS`)
- Métriques k6 granulaires par type de requête (p50, p95)

## 🌪️ CI/CD & Chaos Engineering
- Pipeline **GitHub Actions** avec scans **Gosec** (Go) et **Trivy** (Docker)
- Fix du trigger CI sur toutes les branches
- Scripts PowerShell de chaos (`kill` / `restore` / `latency`) pour NATS, Redis, Postgres, Gateway
- 3 scénarios validés : Kill NATS ✅ · Slow DB 500ms ✅ · Gateway Crash + auto-reconnect ✅

## 🛡️ Sécurité Statique (Gosec)
- Fonction `sanitize()` pour prévenir les injections de logs (CWE-117 / G706)
- Directives `// #nosec G117` justifiées sur les champs `Secret` / `Password`

## 🧪 Tests & Couverture
- **80.8% de couverture** sur le Gateway (objectif >80% ✅)
- Correction des race conditions sur `/auth/refresh` via `testMode`
- Tests : Worker Pool, Rate Limiter, `/healthz`, `/ping-nats`, `/logout`, `sanitize()`

## 🖥️ Robustesse Frontend
- Gestion des codes HTTP critiques : `429` (rate limit) · `503` (service down) · `401` (session expirée)
- Auto-refresh des tokens toutes les 12min, reconnexion WebSocket silencieuse

---

## 📊 Chiffres clés

| Indicateur | Valeur |
|---|---|
| Commits personnels | **~18** |
| Couverture tests Gateway | **80.8%** |
| Latence WS p50 @ 1000 VUs | **32ms** |
| Latence Login p50 @ 1000 VUs | **106ms** |
| Messages traités / 60s | **3,6 millions** |
| Scans sécurité CI | Gosec ✅ Trivy ✅ |
| Scénarios chaos validés | 3/3 ✅ |
