# Plan de Soutenance (Outline)

## 1) Contexte et objectifs
- Problematique: messagerie temps reel distribuee
- Cibles: 100k connexions / 500k msg/s / budget 700 EUR

## 2) Architecture
- Services: gateway + messages
- Broker: NATS
- DB: Postgres
- Cache: Redis
- Observabilite: Prometheus/Grafana

## 3) Choix techniques
- Go pour performance
- NATS pour faible latence
- JWT pour auth stateless

## 4) Qualite et securite
- Tests >80%
- Scans securite

## 5) Performance
- Tests k6 + pprof
- Resultats Storm Day local
- Plan d'echelle vers 100k/500k

## 6) Resilience
- Chaos engineering
- Incidents simules + post-mortem

## 7) Infrastructure
- IaC Terraform (aws-prod)
- CI/CD (CI + CD manuel)

## 8) Budget
- Estimation et scenarios
- Compromis

## 9) Conclusion
- Etat du projet
- Prochaines etapes
