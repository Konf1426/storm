# Azure Load Testing

Ce dossier permet de lancer un test de charge Azure pour STORM via `Azure CLI` avec `Locust`.

## Pourquoi Locust

Le projet contient deja des scripts `k6`, mais Azure Load Testing supporte officiellement les scripts `Locust` et `JMeter` pour les tests scripts manages.

Pour un test WebSocket de `100000` connexions simultanees, l'approche retenue ici est:

- deployer STORM sur Azure
- exposer le `gateway` derriere un endpoint public compatible WebSocket
- utiliser Azure Load Testing avec `Locust`

## Fichiers fournis

- `scripts/azure/locustfile.py`: script Locust qui ouvre et maintient des connexions WebSocket
- `scripts/azure/loadtest-100k.yaml`: configuration Azure Load Testing
- `scripts/azure/requirements.txt`: dependance `websocket-client`
- `scripts/azure/create-loadtest.ps1`: script PowerShell pour creer et lancer le test

## Prerequis

- Azure CLI installee
- extension `load` Azure CLI
- une subscription Azure active
- un endpoint public pour le gateway, par exemple `https://storm.example.com`
- un token JWT valide si l'auth STORM est active

## Parametres a renseigner

Dans `scripts/azure/loadtest-100k.yaml`, remplace:

- `TARGET_HOST` par l'URL publique du gateway
- `ACCESS_TOKEN` par un JWT valide

Exemple:

```yaml
  - name: TARGET_HOST
    value: https://storm.example.com
  - name: ACCESS_TOKEN
    value: eyJhbGciOi...
```

Le script ouvre les connexions sur:

- `GET /ws?subject=storm.events&token=...`

Tu peux aussi tester un channel en remplacant `WS_SUBJECT` par `WS_CHANNEL_ID` dans le YAML et le script.

## Lancement

Connexion Azure:

```powershell
az login
az account set --subscription "<SUBSCRIPTION_ID>"
```

Creation de la ressource et du test:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\azure\create-loadtest.ps1 `
  -SubscriptionId "<SUBSCRIPTION_ID>" `
  -CreateResource
```

Relancer uniquement le test:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\azure\create-loadtest.ps1 `
  -SubscriptionId "<SUBSCRIPTION_ID>"
```

## Dimensionnement

La documentation Azure recommande jusqu'a `500` utilisateurs Locust par engine.

Pour `100000` utilisateurs:

- `100000 / 500 = 200` engines

Le YAML fourni positionne donc `engineInstances: 200`.

## Strategie conseillee

Ne pars pas directement a `100000`.

Fais plutot:

1. `10000`
2. `25000`
3. `50000`
4. `100000`

Ca te permet de verifier:

- taux d'erreur
- latence p95
- CPU/memoire du gateway
- nombre de pods gateway necessaires
- comportement de l'ingress / Application Gateway

## Notes STORM

Le gateway supporte deja les WebSockets sur `GET /ws` et accepte le token via query string `token`, header `Authorization`, ou cookie `access_token`.

Pour un test de connexions massives, il est plus pertinent de:

- reutiliser un token deja valide
- ouvrir et maintenir les connexions
- ne pas faire un `register/login` par utilisateur

Sinon tu testes surtout la couche auth et la persistence, pas la tenue des `100000` sockets.
