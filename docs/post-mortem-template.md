# Post-Mortem Template

## Summary
- What happened:
- Impact:
- Duration:
- Severity:

---

# Post-Mortem (Sample) - Gateway Restart During Storm Day

## Summary
- What happened: Gateway container was restarted during a high-load run (chaos scenario).
- Impact: Short burst of HTTP failures; WS reconnects dropped briefly.
- Duration: ~10s (restart window).
- Severity: Low (SLOs still met for chaos thresholds).

## Timeline
- T0: Load test running (100 HTTP + 100 WS VUs).
- T+X: Gateway restart triggered.
- T+X+10s: Gateway healthy; errors drop back to baseline.

## Root Cause
- Planned chaos action (intentional restart).

## Detection & Response
- Detected via k6 error rate spike and connection refused/EOF.
- Recovery confirmed via `/healthz` and error rate stabilization.

## Metrics
- HTTP error rate: ~3.31%
- HTTP p95: ~31.28ms
- WS connect p95: ~20.84ms
- Throughput: ~2744.7 req/s

## Corrective Actions
- None (expected chaos).
- Keep pre-check for gateway availability before load start.

## Lessons Learned
- Chaos events should be paired with relaxed thresholds.
- Capture recovery time explicitly to compare runs.

## Timeline
- T0: Incident start
- T+X: Detection
- T+Y: Mitigation actions
- T+Z: Recovery confirmed

## Root Cause
- Technical cause:
- Contributing factors:

## Detection & Response
- How we detected it:
- What worked:
- What did not:

## Metrics
- Error rate:
- Latency p95/p99:
- Throughput:

## Corrective Actions
- Short-term fixes:
- Long-term fixes:
- Owners + due dates:

## Lessons Learned
- Process:
- Architecture:
- Team:
