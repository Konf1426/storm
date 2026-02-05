# Runbook (Operations)

## Prerequis
- Docker + Compose
- Variables env (.env) ou secrets

## Demarrage local
```
docker compose -f infra/docker/docker-compose.yml up -d --build
```

## Verification
```
curl -fsS http://localhost:8080/healthz
curl -fsS http://localhost:8080/ping-nats
```

## Smoke tests
```
bash scripts/smoke-test.sh
```

## Deploiement (cloud)
1) Pousser les images en registry (ECR).
2) Appliquer l'IaC.
3) Verifier ALB /healthz.

## Rollback
1) Re-deployer l'image precedente (tag N-1).
2) Redemarrer les services.

## Incident standard
- Symptomes: erreurs 5xx, latence p95 elevee.
- Actions:
  - verifier NATS / DB / Redis
  - redemarrer gateway si besoin
  - activer logs debug temporairement

## Sauvegardes
- RDS: snapshots planifies
- Export logs avant redemarrage majeur
