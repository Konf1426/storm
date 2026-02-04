# Performance Notes

## Load test (k6)

Prerequisites:
- Stack running (gateway on `:8080`)
- k6 installed
- Access token (use `bash scripts/gen-jwt.sh user-1`)

Example:
```
ACCESS_TOKEN=... bash scripts/perf-load.sh
```

Optional env vars:
- `GATEWAY_URL` (default `http://localhost:8080`)
- `SUBJECT` (default `storm.events`)
- `CHANNEL_ID` (use channel WS stream)
- `DURATION` (default `30s`)
- `PUB_VUS` / `WS_VUS` / `PUB_RATE`
- `MODE` (`baseline` or `chaos`)
- `THRESH_FAIL` / `THRESH_P95`

### Example results (Feb 3, 2026)
- Duration: 30s
- VUs: 10 (5 HTTP + 5 WS)
- HTTP p95: ~3.53ms
- HTTP error rate: 0%
- WS connect p95: ~11.54ms
- HTTP req/s: ~48.3/s
- WS msgs received: ~132/s

### Latest run (Feb 3, 2026)
- Duration: 15s
- VUs: 10 (5 HTTP + 5 WS)
- HTTP p95: ~3.52ms
- HTTP error rate: 0%
- WS connect p95: ~1.84ms
- HTTP req/s: ~45.6/s
- WS msgs received: ~134.5/s

### Chaos mode example
```
ACCESS_TOKEN=... MODE=chaos DURATION=5m PUB_VUS=100 WS_VUS=100 PUB_RATE=50 bash scripts/perf-load.sh
```

### Chaos run results (Feb 3, 2026)
- Duration: 5m
- VUs: 200 (100 HTTP + 100 WS)
- HTTP p95: ~31.28ms
- HTTP error rate: ~3.31% (below 10% threshold)
- HTTP req/s: ~2744.7/s
- WS connect p95: ~20.84ms
- WS msgs received: ~115,667/s

## Profiling (pprof)

Enable:
```
PPROF_ADDR=:6060   # gateway
PPROF_ADDR=:6061   # messages
```

Use:
```
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
go tool pprof http://localhost:6060/debug/pprof/heap
```

Checklist:
- Capture CPU profile (30s)
- Capture heap profile
- Note top hot paths and allocations

### Gateway pprof snapshot (Feb 3, 2026)
CPU top (10s):
- `internal/runtime/syscall.Syscall6` ~34.6%
- `pgx` Exec/pgconn read path ~26.9% cum
- `chi` router/middleware ~53.9% cum
- `jwt` parse ~11.5% cum
- `websocket` write ~11.5% cum

Heap top (inuse):
- `compress/flate.NewWriter` ~19.5%
- `bufio.NewReaderSize` ~14.1%
- `runtime.allocm` ~22.2%
- `redis` client connect path ~14.1%

### Messages pprof snapshot (Feb 3, 2026)
CPU top (10s):
- `net/http.(*chunkWriter).writeHeader` ~50%
- `runtime.write1` ~50%

Heap top (inuse):
- `compress/flate.NewWriter` ~49.7%
- `compress/flate.(*compressor).initDeflate` ~15.0%
- `runtime.allocm` ~14.1%
