# STORM

Backend de messagerie temps reel distribue (projet HETIC) avec gateway Go, NATS, persistance Postgres/Redis, frontend Vue, observabilite Prometheus/Grafana et cibles de deploiement Docker/Kubernetes/Terraform.

## Objectifs du projet

- 100 000 connexions simultanees
- 500 000 messages/seconde
- Budget cible <= 700 EUR

## Stack

- Backend: Go (`services/gateway`, `services/messages`)
- Broker: NATS
- Donnees: Postgres + Redis
- Frontend: Vue 3 + Vite (`frontend`)
- Observabilite: Prometheus + Grafana
- Infra: Docker Compose, Kubernetes manifests, Terraform skeletons

## Architecture (resume)

1. Le client s'authentifie sur le gateway (`/auth/*`) via JWT en cookie HttpOnly.
2. Les messages sont publies via HTTP (`/publish`) ou WebSocket (`/ws`).
3. Le gateway publie sur NATS.
4. Le service `messages` consomme NATS.
5. Les messages de channels sont persistes en Postgres, la presence en Redis.
6. Les metriques sont exposees sur `/metrics` et visualisees dans Grafana.

## Arborescence utile

- `services/gateway`: API principale (auth, publish, ws, users, channels)
- `services/messages`: consommateur NATS + health/metrics
- `infra/docker/docker-compose.yml`: stack locale complete
- `infra/k8s`: manifests Kubernetes
- `infra/iac`: Terraform (EC2 simple + AWS prod baseline)
- `scripts`: smoke tests, coverage, scans securite, chaos, perf
- `docs`: architecture, runbook, perf, CI/CD, API OpenAPI

## Prerequis

- Docker + Docker Compose
- Go 1.24.x (pour tests locaux Go)
- Node.js 20+ et npm (frontend)
- `bash`, `curl`, `python3` (pour scripts)

## Demarrage rapide (local)

### 1) Lancer la stack backend

```bash
docker compose -f infra/docker/docker-compose.yml up -d --build
```

Services exposes en local:

- Gateway: `http://localhost:8080`
- Messages: `http://localhost:8081`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000`
- SonarQube: `http://localhost:9000`
- NATS monitoring: `http://localhost:8222`

### 2) Lancer le frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend: `http://localhost:5173`

## Verification rapide

```bash
curl -fsS http://localhost:8080/healthz
curl -fsS http://localhost:8080/ping-nats
bash scripts/smoke-test.sh
```

Option "tout-en-un" locale:

```bash
bash scripts/test-local.sh
```

## API principale (gateway)

Auth:

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `POST /auth/logout`
- `GET /auth/me`

Realtime / messaging:

- `POST /publish?subject=storm.events`
- `GET /ws?subject=storm.events`
- `GET /ws?channel_id=<id>`

Channels / users:

- `GET/POST /channels`
- `GET/POST /channels/{id}/messages`
- `GET/POST /users`
- `GET/PATCH/DELETE /users/{id}`

Ops:

- `GET /healthz`
- `GET /ping-nats`
- `GET /metrics`

Spec OpenAPI: `docs/api/openapi.yml` (note: endpoint SSE `/events` marque legacy).

## Tests, qualite, securite

```bash
bash scripts/coverage.sh
bash scripts/coverage-docker.sh
bash scripts/security-scan.sh
bash scripts/ci-local.sh
```

Objectif coverage: `MIN_COVERAGE` (defaut 80%).

## SonarQube (local)

Demarrer SonarQube:

```bash
docker compose -f infra/docker/docker-compose.yml up -d sonarqube
```

Scanner depuis Windows PowerShell (sans WSL):

```powershell
$env:SONAR_TOKEN="ton_token_sonarqube"
powershell -ExecutionPolicy Bypass -File .\scripts\sonar-scan.ps1
```

Notes:

- Config du projet Sonar: `sonar-project.properties`
- Coverage Go generee automatiquement dans `.sonar/`
- URL dashboard local: `http://localhost:9000`

## Reglages perf utiles (gateway)

- `WORKER_POOL_SIZE`: nombre de workers async DB (defaut docker local: `64`)
- `ENSURE_CACHE_TTL_SECONDS`: cache TTL des `EnsureUser/EnsureMember` (defaut `300`)
- `AUTH_RATE_LIMIT_ENABLED`: desactiver en bench local si besoin (`false`)

## Performance et chaos

```bash
bash scripts/perf-load.sh
bash scripts/chaos.sh full
```

Voir aussi:

- `docs/performance.md`
- `docs/chaos.md`
- `docs/storm-day.md`

## Deploiement

Kubernetes:

```bash
kubectl apply -k infra/k8s
```

Deploiement Kubernetes local recommande:

```bash
bash scripts/k8s-deploy.sh
bash scripts/k8s-port-forward.sh
```

Equivalent PowerShell (Windows, sans WSL):

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\k8s-deploy.ps1
powershell -ExecutionPolicy Bypass -File .\scripts\k8s-port-forward.ps1
```

Nettoyage:

```bash
bash scripts/k8s-cleanup.sh
```

PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\k8s-cleanup.ps1
```

Notes:

- Les scripts buildent `storm-gateway:latest` et `storm-messages:latest` avant apply.
- Variables utiles: `NAMESPACE`, `BUILD_IMAGES=0`, `GATEWAY_IMAGE`, `MESSAGES_IMAGE`.
- Si tu utilises un registre distant, pousse les images puis adapte `infra/k8s/kustomization.yaml`.

IaC Terraform:

- `infra/iac/terraform/aws-ec2`
- `infra/iac/terraform/aws-prod`

Details: `infra/iac/README.md` et `infra/k8s/README.md`.

## Documentation complementaire

- `docs/architecture.md`
- `docs/runbook.md`
- `docs/ci-cd-local.md`
- `docs/ci-cd-k8s.md`
- `docs/final-report.md`
