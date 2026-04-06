# Infrastructure Azure – STORM

Terraform pour provisionner l'infrastructure Azure complète du projet STORM.

## Ressources créées

| Ressource | Description |
|---|---|
| AKS Cluster | 3-10 nodes `Standard_D4s_v5`, autoscaling, OS tuning WebSocket |
| ACR | Container Registry pour les images gateway/messages |
| PostgreSQL Flexible Server | `B_Standard_B2s`, Postgres 16, subnet privé |
| Azure Cache for Redis | Standard C1 (1 Go) |
| Application Gateway v2 | Ingress AGIC, support WebSocket natif |
| Log Analytics | Monitoring Azure Monitor |
| VNet | 3 subnets (AKS, App Gateway, DB) |

## Prérequis

- Azure CLI (`az login`)
- Terraform >= 1.5
- Droits `Contributor` sur la subscription Azure

## Usage

```bash
cd infra/azure/terraform
terraform init
terraform apply -var "db_admin_password=<MOT_DE_PASSE_SECURISE>"
```

## Récupérer le kubeconfig

```bash
az aks get-credentials --resource-group rg-storm --name $(terraform output -raw aks_cluster_name)
```

## Coût estimé

| Scénario | Estimation mensuelle |
|---|---|
| Test ponctuel (quelques heures) | 5-15 € |
| Fonctionnement continu | 150-300 € |

## Nettoyage

```bash
terraform destroy -var "db_admin_password=<MOT_DE_PASSE>"
```
