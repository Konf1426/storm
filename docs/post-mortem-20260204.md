# Post-Mortem - Storm Day Run (Feb 4, 2026)

## Summary
- What happened: Planned chaos actions during Storm Day (NATS pause, gateway latency injection).
- Impact: Temporary latency increase; no errors observed in k6 summaries.
- Duration: ~10s per incident.
- Severity: Low (SLOs met for latency and error rate).

## Timeline
- T0: Warm-up load started (50 HTTP + 50 WS, 2m).
- T+2m: Spike 1 started (150 HTTP + 150 WS, 3m).
- T+5m: Incident 1 (NATS pause, 10s).
- T+5m30s: Recovery window.
- T+6m: Incident 2 (gateway latency, 10s).
- T+6m30s: Spike 2 started (200 HTTP + 200 WS, 3m, chaos).

## Root Cause
- Planned chaos actions (intentional).

## Detection & Response
- Detected via k6 metrics and phase logs.
- No manual mitigation required.

## Metrics
- Warm-up: p95 6.3ms, error 0.00%, req/s 923.0
- Spike 1: p95 81.31ms, error 0.00%, req/s 2501.5
- Spike 2 (chaos): p95 141.5ms, error 0.00%, req/s 2007.7
- WS connect p95 (chaos): 28.42ms

## Corrective Actions
- Capture Prometheus exports next run.
- Add availability computation from HTTP + WS metrics.

## Artifacts
- `out/storm-day/20260204-163421/warmup.log`
- `out/storm-day/20260204-163421/spike1.log`
- `out/storm-day/20260204-163421/spike2.log`
- `out/storm-day/20260204-163421/gateway-cpu.pprof`
- `out/storm-day/20260204-163421/messages-cpu.pprof`
