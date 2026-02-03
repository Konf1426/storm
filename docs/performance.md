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
