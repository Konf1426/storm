# CI/CD Kubernetes

## Objectif
Automatiser:
- build images
- push registry
- deploy sur cluster Kubernetes

## Pre-requis
- Cluster Kubernetes accessible
- KUBECONFIG exporte en secret GitHub `KUBECONFIG_B64`
  - `base64 -w0 ~/.kube/config`
- Registry GHCR (GitHub Container Registry)

## Workflow
- `.github/workflows/cd-k8s.yml`
- Build + push vers `ghcr.io/<owner>/storm-gateway` et `ghcr.io/<owner>/storm-messages`
- Update `kustomization.yaml` (tag + images)
- Apply `kubectl apply -k infra/k8s`

## Secrets attendus
- `KUBECONFIG_B64`: kubeconfig base64 du cluster

## Commande manuelle (local)
```
kubectl apply -k infra/k8s
```
