# Plan : Migration SSE → WebSockets

## État actuel

Le gateway expose **deux** endpoints de streaming en parallèle :

| Endpoint | Protocole | Fichier |
| :--- | :--- | :--- |
| `GET /events` | SSE (Server-Sent Events) | `services/gateway/server.go` – `streamSSE()` |
| `GET /ws` | WebSocket (gorilla/websocket) | `services/gateway/server.go` – `wsHandler()` |

La librairie `gorilla/websocket` est déjà importée et le handler `/ws` est **fonctionnel** :
- Upgrade HTTP → WebSocket ✅
- Souscription NATS + fan-out vers le client ✅
- Réception de messages depuis le client + publish NATS ✅
- Ping/Pong avec keepalive 30s/90s ✅
- Présence Redis (incr/decr) ✅
- Persistance des messages en base ✅

---

## Ce qui reste à faire

### 1. Supprimer le endpoint SSE `/events`

- Retirer la route `pr.Get("/events", ...)` dans `NewRouter()`
- Supprimer la fonction `streamSSE()` (et ses helpers)
- Retirer les imports devenus orphelins si nécessaire

### 2. Migrer le frontend Vue (SSE → WS)

Le frontend actuel (`frontend/src/`) utilise l'API `EventSource` pour consommer le stream SSE.

**Actions :**
- Remplacer les appels `new EventSource(url)` par `new WebSocket(url)` (ou utiliser la classe native `WebSocket`)
- Adapter la gestion des événements :

```diff
- const es = new EventSource('/events?subject=storm.events')
- es.onmessage = (e) => handleMessage(e.data)
+ const ws = new WebSocket('ws://localhost:8080/ws?channel_id=1')
+ ws.onmessage = (e) => handleMessage(e.data)
+ ws.onopen = () => console.log('connected')
+ ws.onclose = () => console.log('disconnected')
```

- Gérer la reconnexion automatique (exponential backoff)
- Passer le JWT dans le query param `?token=<access_token>` pour l'auth sur le handshake WS (déjà supporté par `tokenFromRequest()`)

### 3. Mettre à jour les tests

- Supprimer les tests couvrant `streamSSE` / `/events`
- S'assurer que les tests `wsHandler` couvrent les cas : connexion, envoi, réception, déconnexion, ping/pong
- Viser > 80% de coverage (exigence consignes)

### 4. Vérifier la configuration Kubernetes / Docker

- S'assurer que les services K8s exposent bien le bon port et que le protocol WebSocket est supporté par l'ingress (annotations `nginx.ingress.kubernetes.io/proxy-read-timeout`, `Upgrade`, `Connection`)

---

## Protocole de message (WS)

Format JSON attendu dans les deux sens :

```json
{
  "type": "message",
  "channel_id": 1,
  "payload": "Hello world",
  "user_id": "alice"
}
```

> À définir et à documenter dans `docs/api.md` (à créer).

---

## Ordre d'exécution recommandé

1. [ ] Migrer le frontend (ne casse pas le backend)
2. [ ] Supprimer `streamSSE` + la route `/events` côté gateway
3. [ ] Nettoyer les tests du gateway
4. [ ] Vérifier l'ingress K8s pour le support WebSocket
5. [ ] Test de charge avec k6 ou locust (target : 100K connexions WS simultanées)

---

## Critères de done

- [ ] L'endpoint `/events` (SSE) n'existe plus
- [ ] Le frontend se connecte et reçoit des messages via WS
- [ ] Les tests passent avec > 80% de coverage
- [ ] Un test de charge valide la tenue à 100K connexions simultanées
