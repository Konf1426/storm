# Budget Estimate (<= 700 EUR)

This is a template for a minimal cloud budget. Fill in real prices for the chosen provider.

## Assumptions
- Provider:
- Region:
- Target concurrency: 100,000 connections
- Target throughput: 500,000 messages/sec
- Autoscaling enabled

## Minimal Production Stack (example)
- 2x Gateway instances (autoscale)
- 3x NATS nodes (cluster)
- 1x Postgres (managed)
- 1x Redis (managed)
- 1x Observability stack (Prometheus/Grafana or managed)
- Load balancer + bandwidth

## Cost Breakdown (placeholder)
- Compute (gateway + messages): ___ EUR / month
- NATS cluster: ___ EUR / month
- Postgres managed: ___ EUR / month
- Redis managed: ___ EUR / month
- Monitoring: ___ EUR / month
- Load balancer + bandwidth: ___ EUR / month
- Storage + backups: ___ EUR / month

Total: ___ EUR / month

## Notes
- If budget exceeds 700 EUR, reduce managed services and use fewer nodes.
- If budget is too low for targets, document trade-offs (SLOs, availability, or capacity).
