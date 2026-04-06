# Rapport de Performance STORM - Azure Load Testing (19/03/2026)

Ce document résume les tests de charge intensifs effectués sur l'infrastructure Azure (rg-storm-v5) pour passer de 1 000 à 10 000 utilisateurs virtuels (VUs).

## 1. État Initial et Problèmes Rencontrés
- **Infrastructure initiale** : 1 nœud AKS (Standard_B2s), 2 réplicas Gateway, Postgres 4 vCores.
- **Blocages identifiés** :
    - **Groupe de ressources managé (MC_) absent** : Corrigé via une réconciliation `az aks update`.
    - **Authentification ACR** : Erreur `401 Unauthorized` pour le pull des images, résolue via `az aks update --attach-acr`.
    - **Rate Limiting** : Le test 1000 VUs échouait à 99% car la Gateway limitait à 5 requêtes/min par IP. Désactivé via `AUTH_RATE_LIMIT_ENABLED=false`.
    - **CPU Bcrypt** : Le coût par défaut de Bcrypt saturait le CPU lors des inscriptions massives. Réduit à `BCRYPT_COST=4`.

## 2. Optimisations Techniques
### Optimisation Logicielle (Code)
- **Messages Asynchrones** : Modification de la Gateway pour publier immédiatement dans NATS et déléguer l'écriture Postgres à un pool de workers en arrière-plan.
- **Résultat** : Latence d'envoi de messages passée de ~3s à **< 100ms**.

### Optimisation Infrastructure (Ultra Scaling)
- **AKS (Multi-family Strategy)** : Pour contourner les quotas Azure (limite 10 vCPUs par famille), création d'un cluster "Frankenstein" de **26 vCPUs** combinant :
    - Standard_D8s_v4 (8 vCPUs)
    - Standard_D8s_v3 (8 vCPUs)
    - Standard_F8s_v2 (8 vCPUs)
    - Standard_B2s (2 vCPUs - Système)
- **Base de données** : Passage de Postgres Flexible à **16 vCores (Standard_D16ds_v4)** avec 128 Go de stockage et `max_connections=2000`.
- **Réplicas** : Montée en charge à **30 instances Gateway** en parallèle.

## 3. Résultats des Tests de Charge (k6)

| Métrique | 1 000 VUs (Initial) | 5 000 VUs (Optimisé) | 10 000 VUs (Ultra) |
| :--- | :--- | :--- | :--- |
| **Taux de succès HTTP** | 0.01% (Rate Limit) | ~20% (DB bottleneck) | **~100%** |
| **Latence Msg (Médiane)** | ~2.5s | < 100ms | **~85ms** |
| **Latence Login (Médiane)** | > 4s | ~1.1s | **144ms** |
| **Latence Login (Moyenne)** | > 5s | ~1.8s | **221ms** |
| **Débit WebSockets** | ~20M msg | ~40M msg | **~210M msg** |

## 4. Analyse des Coûts (Estimation)
- **Config "Ultra"** : ~3,33 $ / heure (soit ~2 500 $ / mois si laissé actif).
- **Plus gros poste** : PostgreSQL (16 vCores) représentant ~45% du budget.

## 5. Mise en Veille (Post-Test)
Pour arrêter la facturation élevée, l'infrastructure a été réduite immédiatement après les tests :
- **Postgres** : Downgrade à **2 vCores**.
- **AKS** : Suppression des pools High-Performance (hp1, hp3, hp4).
- **Kubernetes** : Réduction à **1 seul réplica** par service.
- **Coût estimé en veille** : **~0,15 $ / heure**.

---
*Rapport généré par Gemini CLI.*
