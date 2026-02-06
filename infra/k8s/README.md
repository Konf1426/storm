# Kubernetes Manifests

This folder provides a baseline Kubernetes deployment for STORM.
It is designed for a small cluster (dev/staging) and is not hardened for production.

## Apply
```
kubectl apply -k infra/k8s
```

## Notes
- Images are expected in a registry (set in `kustomization.yaml`).
- Secrets are placeholders; replace before apply.
- Use `secret.template.yaml` and `configmap.template.yaml` for prod values.
- Ingress is provided for gateway, grafana, prometheus (host: storm.local).
- For production, use managed Postgres/Redis and proper TLS/Ingress.
