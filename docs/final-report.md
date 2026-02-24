# Rapport Final - Projet STORM

## Resume
Le projet STORM met en place un backend de messagerie temps reel distribue avec un gateway HTTP/WebSocket, un service messages, NATS comme broker, Postgres pour la persistance, Redis pour la presence, Prometheus/Grafana pour l'observabilite, et une CI GitHub Actions. Les phases Developpement, Qualite/Securite, Performance et Chaos Engineering sont couvertes par du code, des tests, des scripts et des docs. Storm Day est planifie et documente.

## Architecture & Stack
- Langage: Go
- Broker: NATS
- DB: Postgres
- Cache: Redis
- Temps reel: WebSocket (SSE legacy)
- Auth: JWT (access + refresh)
- Observabilite: Prometheus/Grafana
- Infra locale: Docker Compose

Docs principales:
- `docs/architecture.md`
- `docs/api/openapi.yml`
- `docs/performance.md`
- `docs/chaos.md`
- `docs/storm-day.md`
- `docs/post-mortem-template.md`

## Qualite & Securite
- Tests unitaires + integration (NATS embarque, mocks Postgres/Redis).
- Coverage >= 80%:
  - gateway: 80.4% (Feb 4, 2026)
  - messages: 83.6% (Feb 4, 2026)
- Scans dependances: `scripts/security-scan.sh` (requires govulncheck + gosec).
- CI: `.github/workflows/ci.yml` (coverage + scans).

## Performance
- Load tests k6 + pprof actives.
- Resultats references dans `docs/performance.md` (runs Feb 3, 2026).

## Chaos Engineering
Scenarios disponibles via `scripts/chaos.sh`:
- restart gateway/messages
- pause NATS
- injection de latence
- CPU spike

## Storm Day
Plan defini dans `docs/storm-day.md` (SLOs + timeline + roles). Execution reelle effectuee le Feb 4, 2026. Resultats dans `docs/storm-day-results.md` et post-mortem dans `docs/post-mortem-20260204.md`.

## Ecart / Reste a faire
- IaC cloud: stack AWS completee (VPC, ALB, ASG, RDS, Redis, ECR). A durcir pour prod (TLS, securite, monitoring managé).
- Budget cloud detaille: chiffrage indicatif dans `docs/budget.md`, a valider via AWS Pricing Calculator.
- Soutenance (slides + recit final).
