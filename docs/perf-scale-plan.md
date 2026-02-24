# Plan de Tests Grande Echelle (100k connexions / 500k msg/s)

## Objectif
Valider la tenue en charge proche des cibles consignes:
- 100 000 connexions simultanees
- 500 000 messages par seconde

## Strategie generale
1) Tests progressifs (10k -> 50k -> 100k connexions).
2) Separation des charges: publication vs reception (WS/SSE).
3) Observabilite active (latence p95/p99, erreurs, throughput).
4) Itineraires de rollback (seuils d'arret si erreurs > 1%).

## Outillage
- k6 (ou k6 cloud) pour VUs distribues.
- Instances generatrices (EC2/GCE) a proximite de la region cible.
- Scripts existants: `scripts/perf-load.sh`, `scripts/k6-load.js`.

## Dimensionnement initial (exemple)
- 4 generateurs x 25k connexions = 100k.
- Publication cible: 500k msg/s (125k msg/s par generateur).
- Duree: 10-15 min par palier.

## Scenarios
1) Warm-up: 10k connexions, 1-2 min.
2) Palier 1: 50k connexions, 5 min.
3) Palier 2: 100k connexions, 10 min.
4) Chaos: latence + restart gateway pendant palier 2.

## Metriques a capturer
- HTTP/WS connect p95/p99
- Erreur HTTP/WS
- Throughput (req/s, msg/s)
- CPU/RAM gateway/messages
- Saturation broker/DB/cache

## Resultats attendus
- p95 publish < 200ms
- erreur < 1%
- WS connect p95 < 200ms

## Livrables
- Logs k6
- Screenshots dashboards
- Rapport synthese (tableau des paliers)
