# Chaos Engineering

This project includes a simple chaos script to simulate failures and latency.

## Prereqs
- Stack running via `docker compose -f infra/docker/docker-compose.yml up -d`
- Docker engine available

## Scenarios
```
bash scripts/chaos.sh gateway-restart
bash scripts/chaos.sh messages-restart
bash scripts/chaos.sh nats-pause
bash scripts/chaos.sh gateway-latency
bash scripts/chaos.sh messages-latency
bash scripts/chaos.sh gateway-cpu
bash scripts/chaos.sh full
```

## Expected observations
- Gateway restart: WS reconnects, publishing resumes after restart.
- Messages restart: consumers resume, live feed continues.
- NATS pause: temporary disconnects, recovery after unpause.
- Latency: increased response times and delayed WS events.
- CPU spike: higher latency and lower throughput.

## Notes
- The latency scenario uses `tc` inside the container (installs `iproute2` temporarily).
- Containers need `NET_ADMIN` (added in compose) for `tc` to work.
- This is a local chaos harness; for production, use a dedicated tool (e.g. Litmus, Chaos Mesh).
