# Storm Day Results

Date: Feb 4, 2026
Run ID: 20260204-163421
Environment: local docker-compose

## Summary
- Overall outcome (pass/fail): pass (latency + error thresholds met)
- Key incidents: NATS pause (10s), gateway latency injection (100ms +/- 20ms for 10s)
- SLOs met: latency and error rate met; availability not computed

## SLOs
- Availability: not computed (local run)
- Publish latency p95: 141.5ms (chaos spike)
- WebSocket connect p95: 28.42ms (chaos spike)
- Error rate: 0.00%

## Load Test Results
- Warm-up (50 HTTP + 50 WS, 2m): p95 6.3ms, error 0.00%, req/s 923.0, ws/s 23408.0
- Spike 1 (150 HTTP + 150 WS, 3m): p95 81.31ms, error 0.00%, req/s 2501.5, ws/s 144984.9
- Spike 2 (200 HTTP + 200 WS, 3m, chaos): p95 141.5ms, error 0.00%, req/s 2007.7, ws/s 171438.9

## Incidents
1) Incident 1 (NATS pause)
   - Start: during Storm Day phase 3
   - Duration: 10s
   - Impact: no k6 errors observed in summary
2) Incident 2 (Gateway latency)
   - Start: during Storm Day phase 5
   - Duration: 10s
   - Impact: higher latency, no k6 errors observed in summary

## Metrics/Artifacts
- k6 logs: `out/storm-day/20260204-163421/warmup.log`, `out/storm-day/20260204-163421/spike1.log`, `out/storm-day/20260204-163421/spike2.log`
- pprof snapshots: `out/storm-day/20260204-163421/gateway-cpu.pprof`, `out/storm-day/20260204-163421/messages-cpu.pprof`
- Prometheus screenshots/exports: not captured

## Post-Mortem
Reference: `docs/post-mortem-template.md`
