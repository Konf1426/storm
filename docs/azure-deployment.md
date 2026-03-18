# Déploiement STORM sur Azure AKS

Guide complet pour déployer STORM sur Azure Kubernetes Service (AKS) et lancer des tests de charge à 100 000 connexions WebSocket simultanées.

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                    Azure Cloud                       │
│                                                      │
│  ┌──────────────┐    ┌──────────────────────────┐   │
│  │ App Gateway  │───▶│     AKS Cluster           │   │
│  │ v2 (WS+TLS) │    │  ┌─────────┐ ×3-20       │   │
│  └──────────────┘    │  │ Gateway │──────┐       │   │
│                      │  └─────────┘      │       │   │
│                      │  ┌─────────┐      │       │   │
│                      │  │  NATS   │◀─────┤       │   │
│                      │  └─────────┘      │       │   │
│                      │  ┌─────────┐      │       │   │
│                      │  │Messages │      │       │   │
│                      │  └─────────┘      │       │   │
│                      │  ┌─────────┐      │       │   │
│                      │  │k6 Runners│     │       │   │
│                      │  └─────────┘      │       │   │
│                      └───────────────────┼───────┘   │
│                                          │           │
│  ┌──────────────┐   ┌──────────────┐     │           │
│  │  PostgreSQL  │◀──┤  Azure Cache │◀────┘           │
│  │  Flexible    │   │  for Redis   │                 │
│  └──────────────┘   └──────────────┘                 │
└─────────────────────────────────────────────────────┘
```

## Prérequis

- **Azure CLI** : `az login`
- **Terraform** >= 1.5
- **Docker** : pour build les images
- **kubectl** : pour interagir avec AKS
- Droits **Contributor** sur la subscription Azure

## Déploiement rapide

### Option 1 : Script tout-en-un

```powershell
powershell -ExecutionPolicy Bypass -File scripts/azure-deploy.ps1
```

Ce script fait tout : Terraform → build Docker → push ACR → kubectl apply.

### Option 2 : Étape par étape

#### 1. Provisionner l'infrastructure

```bash
cd infra/azure/terraform
terraform init
terraform apply -var "db_admin_password=<MOT_DE_PASSE>"
```

#### 2. Configurer kubectl

```bash
az aks get-credentials --resource-group rg-storm --name $(terraform output -raw aks_cluster_name)
```

#### 3. Build et push les images

```bash
ACR=$(terraform output -raw acr_login_server)
az acr login --name ${ACR%%.*}
docker build -t $ACR/storm-gateway:latest services/gateway
docker push $ACR/storm-gateway:latest
docker build -t $ACR/storm-messages:latest services/messages
docker push $ACR/storm-messages:latest
```

#### 4. Déployer sur AKS

```bash
kubectl apply -k infra/k8s-azure
```

#### 5. Vérifier

```bash
kubectl get pods -n storm
kubectl get hpa -n storm
```

## Tests de charge

### Lancer un test à 100k connexions

```powershell
# 4 pods × 25k VUs = 100k connexions simultanées
powershell -ExecutionPolicy Bypass -File scripts/azure-load-test.ps1

# Test progressif : 10k d'abord
powershell -ExecutionPolicy Bypass -File scripts/azure-load-test.ps1 -TargetVUs 2500

# Puis 50k
powershell -ExecutionPolicy Bypass -File scripts/azure-load-test.ps1 -TargetVUs 12500
```

### Monitoring pendant le test

```bash
# HPA (autoscaling)
kubectl get hpa -n storm -w

# Pods gateway
kubectl get pods -n storm -w -l app=gateway

# Grafana
kubectl port-forward svc/grafana 3000:3000 -n storm
```

## Estimation des coûts

| Ressource | SKU | Estimation mensuelle |
|---|---|---|
| AKS (3 nodes D4s_v5) | Standard | ~100 €/mois |
| PostgreSQL Flexible | B_Standard_B2s | ~30 €/mois |
| Azure Cache Redis | Standard C1 | ~40 €/mois |
| Application Gateway v2 | Standard_v2 | ~80 €/mois |
| ACR | Basic | ~5 €/mois |
| **Total** | | **~255 €/mois** |

> Pour un test ponctuel de quelques heures : **5-15 €**.

## Nettoyage

```powershell
powershell -ExecutionPolicy Bypass -File scripts/azure-cleanup.ps1
```

## Fichiers clés

| Fichier | Description |
|---|---|
| `infra/azure/terraform/main.tf` | Infrastructure Terraform |
| `infra/k8s-azure/kustomization.yaml` | Overlay K8s pour Azure |
| `infra/k8s-azure/gateway-hpa.yaml` | Autoscaling gateway |
| `infra/k8s-azure/ingress-azure.yaml` | Ingress AGIC WebSocket |
| `scripts/k6/storm-azure.js` | Script k6 distribué |
| `infra/k8s-azure/k6-runner/k6-job.yaml` | Job k6 parallèle |
| `scripts/azure-deploy.ps1` | Déploiement tout-en-un |
| `scripts/azure-load-test.ps1` | Lancement tests de charge |
