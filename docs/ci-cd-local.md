# CI/CD Local (sans cloud)

## Objectif
Reproduire localement la CI/CD de base:
- tests + coverage
- scans securite (si outils installes)
- smoke tests

## Commandes
```
bash scripts/ci-local.sh
docker compose -f infra/docker/docker-compose.yml up -d --build
bash scripts/smoke-test.sh
```

## Notes
- Les scans securite utilisent govulncheck + gosec.
- Les smoke tests supposent gateway sur :8080.
