# STORM Day Results - 2026-02-24

**Date :** 24 février 2026
**Environnement :** Local Docker-Compose
**Objectif :** Validation de la résilience (Chaos Engineering) - Scénarios 2 (Slow DB) et 3 (Gateway Crash).

## 📊 Résumé des Scénarios

### Scénario 2 : Latence Base de Données (Slow DB)
- **Configuration :** 500ms de latence artificielle via `SIMULATE_DB_DELAY`.
- **Charge :** 100 VUs via k6.
- **Observations :**
  - Le Worker Pool (20 workers) a permis de maintenir le service mais a subi une forte pression.
  - La latence d'envoi de messages a augmenté proportionnellement au délai DB.
  - Le système est resté stable sans crash.

### Scénario 3 : Crash du Gateway
- **Configuration :** Arrêt brutal du conteneur `gateway` pendant un test de charge.
- **Charge :** 100 VUs via k6.
- **Observations :**
  - **Taux de succès HTTP :** 0.36% (quasiment toutes les requêtes ont échoué pendant l'arrêt).
  - **Durée de l'interruption :** Environ 60 secondes (durée du test k6 pendant laquelle le gateway était hors service).
  - **Rétablissement :** Le service a redémarré avec succès. Des erreurs de contraintes d'intégrité (Foreign Key) ont été observées dans les logs au redémarrage, suggérant des tentatives de reconnexion de clients avec des données inconsistantes ou expirées.
  - **Validation :** Un smoke test post-crash confirme que le Gateway est à nouveau opérationnel (`/healthz` et `/ping-nats` OK).

## ⚡ Métriques k6 (Scenario 3)
- **Messages échangés (Total) :** 68,855
- **Temps de connexion WS (médiane) :** 1.72 ms
- **Taux de succès global :** < 1% (attendu pour un crash test)

## ✅ Conclusions
Le système STORM démontre une capacité de rétablissement automatique après un crash du Gateway. Cependant, la gestion de la consistance de la base de données au redémarrage immédiat sous charge pourrait être affinée pour éviter les erreurs de clés étrangères observées.
