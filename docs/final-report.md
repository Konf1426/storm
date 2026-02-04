# Rapport Final — Projet STORM

## Résumé
Le projet STORM met en place un backend de messagerie temps réel distribué avec un gateway HTTP/WebSocket, un service messages, NATS comme broker, Postgres pour la persistance, Redis pour le cache/presence, Prometheus/Grafana pour l’observabilité, et une CI GitHub Actions. Les phases Développement, Qualité/Sécurité, Performance et Chaos Engineering ont été couvertes avec tests, scans, load tests, pprof, et scénarios d’incidents. Storm Day est planifié et documenté.

## Architecture & Stack
- Langage: Go
- Broker: NATS
- DB: Postgres
- Cache: Redis
- Temps réel: WebSocket (SSE legacy)
- Auth: JWT (access + refresh)
- Observabilité: Prometheus/Grafana
- Infra: Docker compose local

Docs:
- `docs/architecture.md`
- `docs/api/openapi.yml`
- `docs/performance.md`
- `docs/chaos.md`
- `docs/storm-day.md`
- `docs/post-mortem-template.md`

## Qualité & Sécurité
- Tests unitaires + intégration
- Coverage >= 80%:
  - gateway: 80.0%
  - messages: 81.6%
- Scan dépendances: `scripts/security-scan.sh`
- CI: `.github/workflows/ci.yml`

## Performance
Load tests k6 + pprof activés.
Résultats notables:
- Baseline p95 ~3–4ms, erreurs ~0%
- Chaos run (200 VUs): p95 ~31ms, erreurs ~3.31% (sous seuil chaos)

## Chaos Engineering
Scénarios disponibles via `scripts/chaos.sh`:
- restart gateway/messages
- pause NATS
- injection de latence
- CPU spike

Résultats documentés dans `docs/performance.md` et `docs/post-mortem-template.md`.

## Storm Day
Plan défini dans `docs/storm-day.md`, SLOs et timeline détaillés.

## Écart / Reste à faire
- Infrastructure cloud / IaC (Terraform/Pulumi/CloudFormation)
- Déploiement production (K8s/Swarm/Serverless)
- Budget cloud détaillé
- Post‑mortem réel après un Storm Day complet
- Rapport de soutenance (slides + récit)

## Prochaines étapes recommandées
1. Définir IaC + environnements (dev/staging/prod).
2. Mettre en place un pipeline de déploiement.
3. Exécuter un Storm Day complet et produire le post‑mortem final.
