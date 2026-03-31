# Rapport de Performance STORM - Azure Load Testing (19/03/2026)

Ce document résume les tests de charge intensifs effectués sur l'infrastructure Azure `rg-storm-v5` pour faire évoluer STORM de `1 000` à `10 000` utilisateurs virtuels (`VUs`).

## 1. État initial et problèmes rencontrés

- **Infrastructure initiale** : `1` nœud AKS `Standard_B2s`, `2` réplicas Gateway, PostgreSQL `4 vCores`.
- **Blocages identifiés** :
  - **Groupe de ressources managé (`MC_`) absent** : corrigé via une réconciliation `az aks update`.
  - **Authentification ACR** : erreur `401 Unauthorized` lors du pull des images, résolue via `az aks update --attach-acr`.
  - **Rate limiting** : le test `1000 VUs` échouait à `99%` car la Gateway limitait à `5 req/min/IP`. Désactivé pour la campagne via `AUTH_RATE_LIMIT_ENABLED=false`.
  - **CPU Bcrypt** : le coût par défaut de `bcrypt` saturait le CPU pendant les inscriptions massives. Réduit à `BCRYPT_COST=4` pour la campagne de benchmark.

## 2. Optimisations techniques

### Optimisation logicielle

- **Messages asynchrones** : la Gateway publie immédiatement dans `NATS` et délègue l'écriture PostgreSQL à un pool de workers en arrière-plan.
- **Résultat** : la latence d'envoi des messages passe de `~3s` à **`< 100ms`**.

### Optimisation infrastructure

- **AKS (multi-family strategy)** : montée à **`26 vCPUs`** pour contourner les quotas Azure par famille.
  - `Standard_D8s_v4` : `8 vCPUs`
  - `Standard_D8s_v3` : `8 vCPUs`
  - `Standard_F8s_v2` : `8 vCPUs`
  - `Standard_B2s` : `2 vCPUs` pour le système
- **Base de données** : passage à **PostgreSQL Flexible `16 vCores` (`Standard_D16ds_v4`)** avec `128 Go` de stockage et `max_connections=2000`.
- **Réplicas Gateway** : montée à **`30` instances** en parallèle.

## 3. Résultats des tests de charge (k6)

| Métrique | 1 000 VUs (Initial) | 5 000 VUs (Optimisé) | 10 000 VUs (Ultra) |
| :--- | :--- | :--- | :--- |
| **Taux de succès HTTP** | `0.01%` (rate limit) | `~20%` (bottleneck DB) | **`~100%`** |
| **Latence message (médiane)** | `~2.5s` | `< 100ms` | **`~85ms`** |
| **Latence login (médiane)** | `> 4s` | `~1.1s` | **`144ms`** |
| **Latence login (moyenne)** | `> 5s` | `~1.8s` | **`221ms`** |
| **Débit WebSocket** | `~20M msg` | `~40M msg` | **`~210M msg`** |

## 4. Analyse des coûts

- **Configuration "Ultra"** : `~3,33 $ / heure` soit `~2 500 $ / mois` si l'infrastructure restait active.
- **Plus gros poste** : PostgreSQL `16 vCores`, représentant `~45%` du budget.

## 5. Mise en veille après test

Pour stopper la facturation élevée juste après la campagne :

- **PostgreSQL** : downgrade à `2 vCores`
- **AKS** : suppression des pools high-performance `hp1`, `hp3`, `hp4`
- **Kubernetes** : réduction à `1` seul réplica par service
- **Coût estimé en veille** : **`~0,15 $ / heure`**

## 6. Lecture honnête à défendre à l'oral

- Le test du **19 mars 2026** valide une **montée en charge cloud crédible à `10 000 VUs`**, pas une validation finale des `100 000` connexions du sujet.
- Les réglages `AUTH_RATE_LIMIT_ENABLED=false` et `BCRYPT_COST=4` doivent être présentés comme des **réglages de benchmark**, pas comme une configuration de production.
- Le principal gain vient d'une **évolution d'architecture** : publication immédiate dans `NATS` + persistance PostgreSQL asynchrone.
- Le principal apprentissage budgétaire est que la **persistance devient le vrai poste de coût** à forte charge.
