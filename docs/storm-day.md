# STORM DAY Plan

## Objectives
- Simulate viral launch with rapid user growth.
- Validate resiliency under incidents.
- Track SLOs and communicate clearly during outages.

## SLO Targets (baseline)
- Availability: 99.5% during the event
- Publish latency p95: < 200ms
- WebSocket connect p95: < 200ms
- Error rate: < 1%

## Timeline (example 60-90 min)
1. Warm-up (10 min): ramp traffic to steady load.
2. Growth spike (15 min): 3-5x throughput.
3. Incident 1 (10 min): NATS pause / gateway restart.
4. Recovery (10 min): monitor stabilization.
5. Incident 2 (10 min): latency injection on gateway/messages.
6. Growth spike 2 (15 min): peak load + error budget check.
7. Cooldown + data capture (10 min).

## Traffic Simulation
- Use `scripts/perf-load.sh` with increased VUs/duration.
- Or run the orchestrated script: `bash scripts/storm-day-runner.sh`
- Example:
  ```
  ACCESS_TOKEN=... DURATION=10m PUB_VUS=200 WS_VUS=200 PUB_RATE=50 bash scripts/perf-load.sh
  ```

## Cloud Execution (when infra ready)
1) Pousser images en registry (ECR).
2) Deployer infra (Terraform aws-prod).
3) Lancer le Storm Day depuis des generateurs distribues.
4) Capturer metriques + incidents + post-mortem.

## Incident Scenarios
- `bash scripts/chaos.sh nats-pause`
- `bash scripts/chaos.sh gateway-restart`
- `bash scripts/chaos.sh gateway-latency`
- `bash scripts/chaos.sh messages-restart`

## Observability Checklist
- Prometheus up: `http://localhost:9090`
- Grafana up: `http://localhost:3000`
- Monitor:
  - HTTP errors, latency (p95)
  - WS connect success rate
  - NATS and Postgres health
  - CPU / memory

## Data to capture
- k6 results (stdout)
- pprof snapshots (CPU + heap)
- Prometheus screenshots or exported metrics
- Incident timeline notes
- Save results in `docs/storm-day-results.md`

## Roles during Storm Day
- Incident commander
- Comms lead
- Operator (infra)
- Scribe (timeline + metrics)
