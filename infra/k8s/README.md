# Kubernetes Manifests

This folder provides a baseline Kubernetes deployment for STORM.
It is designed for a small cluster (dev/staging) and is not hardened for production.

## Apply
```
kubectl apply -k infra/k8s
```

## Quickstart (local cluster)
```
bash scripts/k8s-deploy.sh
bash scripts/k8s-port-forward.sh
```

Cleanup:
```
bash scripts/k8s-cleanup.sh
```

## Notes
- Images are expected in a registry (set in `kustomization.yaml`).
- Secrets are placeholders; replace before apply.
- Use `secret.template.yaml` and `configmap.template.yaml` for prod values.
- Ingress is provided for gateway, grafana, prometheus (host: storm.local).
- TLS secret expected: `storm-tls` (create with real cert).
- For production, use managed Postgres/Redis and proper TLS/Ingress.
