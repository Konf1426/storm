# Slides - Soutenance (Draft)

Slide 1: Contexte
- Messagerie temps reel distribuee
- Objectifs: 100k connexions, 500k msg/s, budget 700 EUR

Slide 2: Architecture
- Gateway + Messages
- NATS, Postgres, Redis
- Observabilite Prometheus/Grafana

Slide 3: Choix techniques
- Go pour performance
- NATS pour pub/sub faible latence
- JWT stateless

Slide 4: Qualite & Securite
- Tests >80%
- Scans securite

Slide 5: Performance
- k6 + pprof
- Storm Day local
- Plan grande echelle

Slide 6: Resilience
- Chaos engineering
- Incidents simules

Slide 7: Infra & DevOps
- IaC Terraform
- CI + CD manuel

Slide 8: Budget
- Estimation + scenarios
- Compromis

Slide 9: Conclusion
- Etat du projet
- Prochaines etapes
