# Test Results (Local)

Date: Feb 4, 2026

## Unit tests
Executed using Docker because Go was not available on the host.

Gateway:
```
docker run --rm -v "<repo>/services/gateway:/app" -w /app golang:1.24-alpine sh -lc "/usr/local/go/bin/go test ./..."
```
Result: OK

Messages:
```
docker run --rm -v "<repo>/services/messages:/app" -w /app golang:1.24-alpine sh -lc "/usr/local/go/bin/go test ./..."
```
Result: OK

## Coverage
Gateway:
```
total: (statements) 80.4%
```

Messages:
```
total: (statements) 83.6%
```

## Notes
- Security scans not executed in this run (govulncheck + gosec not installed on host).
