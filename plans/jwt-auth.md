# Plan : Authentification JWT

## État actuel — ✅ Déjà implémenté

Le backend **implémente déjà** un système JWT complet dans `services/gateway/server.go` :

### Structures & configuration

```go
type AuthConfig struct {
    Secret        []byte        // secret access token (HS256)
    RefreshSecret []byte        // secret refresh token
    Enabled       bool          // active/inactive via env
    AccessTTL     time.Duration // durée de vie access token
    RefreshTTL    time.Duration // durée de vie refresh token
    CookieDomain  string
    CookieSecure  bool
    CorsOrigin    string
}
```

### Endpoints disponibles

| Méthode | Route | Description |
| :--- | :--- | :--- |
| `POST` | `/auth/register` | Créer un utilisateur (bcrypt password) |
| `POST` | `/auth/login` | Login → émet access + refresh token (cookies HttpOnly) |
| `POST` | `/auth/refresh` | Échange le refresh token contre un nouvel access token |
| `POST` | `/auth/logout` | Révoque le refresh token, clear les cookies |
| `GET` | `/auth/me` | Retourne l'utilisateur connecté (protégé) |

### Middleware

`authMiddleware(cfg AuthConfig)` : protège toutes les routes sous `/` en vérifiant le JWT depuis :
1. Header `Authorization: Bearer <token>`
2. Cookie `access_token`
3. Query param `?token=<token>` (utilisé pour le handshake WebSocket)

### Stockage des refresh tokens

Les refresh tokens sont stockés en base Postgres (`refresh_tokens` table) via `Store.SaveRefreshToken` / `GetRefreshToken` / `RevokeRefreshToken`. La révocation est bien implémentée.

---

## Ce qui reste à faire

### 1. Activer JWT en production (variable d'environnement)

JWT est contrôlé par `AuthConfig.Enabled`. Il faut s'assurer que `AUTH_ENABLED=true` est bien positionné dans les configs K8s / Docker Compose.

**Fichiers à vérifier :**
- `infra/k8s/` → ConfigMap ou Secret pour `JWT_SECRET`, `JWT_REFRESH_SECRET`, `AUTH_ENABLED`
- `infra/docker/` → `docker-compose.yml` pour le développement local

### 2. Gérer la rotation des secrets JWT

- Utiliser des **Kubernetes Secrets** (et non des ConfigMaps) pour `JWT_SECRET` et `JWT_REFRESH_SECRET`
- Documenter la procédure de rotation (invalidation des sessions actives)

### 3. Intégrer l'auth dans le frontend

Le frontend Vue doit :
- [ ] Appeler `POST /auth/register` et `POST /auth/login`
- [ ] Stocker le token (cookies HttpOnly gérés automatiquement par le navigateur)
- [ ] Inclure le header `Authorization: Bearer <token>` ou laisser le cookie passer automatiquement sur les requêtes
- [ ] Passer `?token=<access_token>` dans l'URL de connexion WebSocket (contrainte du protocole WS qui ne supporte pas les headers custom à l'upgrade)
- [ ] Implémenter le refresh automatique : détecter une 401 → appeler `/auth/refresh` → rejouer la requête
- [ ] Gérer le logout (appel à `/auth/logout`)

### 4. Tests à écrire / compléter

- [ ] Test unitaire `authMiddleware` : token valide, token expiré, token absent, mauvais algo
- [ ] Test intégration `/auth/login` → `/auth/refresh` → `/auth/logout`
- [ ] Test que les routes protégées retournent 401 sans token
- [ ] Viser > 80% coverage (exigence consignes)

### 5. Sécurité à vérifier

| Point | Status |
| :--- | :--- |
| Cookies `HttpOnly` | ✅ Implémenté |
| Cookies `Secure` (HTTPS seulement) | ⚠️ Configurable via `CookieSecure` — à activer en prod |
| `SameSite=Lax` | ✅ Implémenté |
| Signature HS256 | ✅ Implémenté |
| Vérification algo (`!ok` check) | ✅ Implémenté |
| Révocation refresh token | ✅ Implémenté |
| Rate limiting sur `/auth/login` | ❌ À ajouter (protection brute force) |

---

## Ordre d'exécution recommandé

1. [ ] Ajouter les Kubernetes Secrets pour `JWT_SECRET` / `JWT_REFRESH_SECRET`
2. [ ] Activer `AUTH_ENABLED=true` dans tous les environnements
3. [ ] Activer `COOKIE_SECURE=true` en prod (HTTPS)
4. [ ] Intégrer l'auth dans le frontend Vue
5. [ ] Ajouter du rate limiting sur les routes `/auth/*`
6. [ ] Compléter les tests unitaires et d'intégration

---

## Critères de done

- [ ] `AUTH_ENABLED=true` en production (K8s secrets configurés)
- [ ] Les WebSockets utilisent le JWT pour s'authentifier
- [ ] Le frontend implémente le flux login / refresh / logout
- [ ] Rate limiting sur `/auth/login` et `/auth/register`
- [ ] Tests auth avec > 80% de coverage
