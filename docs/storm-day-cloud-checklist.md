# Storm Day - Checklist Cloud

## Avant l'execution
- [ ] Infra deployee (ALB/RDS/Redis OK)
- [ ] Images push ECR (tag correct)
- [ ] Observabilite active (dashboards ouverts)
- [ ] Generateurs distribues prets

## Pendant l'execution
- [ ] Warm-up lance
- [ ] Spike 1 lance
- [ ] Incident NATS pause
- [ ] Incident latency gateway
- [ ] Spike 2 lance

## Apres l'execution
- [ ] Logs k6 sauvegardes
- [ ] Export Prometheus / dashboards
- [ ] Post-mortem rempli
