# STORM Architecture

## Overview
STORM is a distributed real-time messaging backend. The current MVP is built around a simple HTTP gateway that publishes events to NATS, a consumer service that subscribes to NATS, and a minimal observability stack (Prometheus/Grafana). A Vue demo frontend consumes events over SSE from the gateway to demonstrate real-time behavior.

## Goals (from consignes)
- High concurrency: 100,000 simultaneous connections
- High throughput: 500,000 messages/second
- Budget: <= 700 EUR (AWS or equivalent)

## Current Components
- Gateway service (Go): HTTP ingress, publish to NATS, SSE stream for clients, metrics
- Messages service (Go): NATS subscriber, logs messages, metrics
- NATS: message broker
- Prometheus + Grafana: metrics and dashboards
- Postgres + Redis: used for users/channels/messages and presence
- Frontend (Vue): realtime dashboard using SSE
- Kubernetes manifests: baseline deployment in `infra/k8s`

## System Flow (Current)
1) Client sends HTTP POST to `/publish?subject=storm.events`.
2) Gateway publishes payload to NATS.
3) Messages service subscribes to `storm.events` and logs payload.
4) Gateway SSE endpoint `/events?subject=storm.events` streams incoming events to clients.

## Target Architecture (Planned)
- Ingress: HTTP API + WebSocket gateway
- Broker: NATS
- Auth: JWT (stateless, signed)
- User state: Redis (presence, fan-out state)
- Persistence: Postgres (users, channels, messages history)
- Observability: metrics + logs + traces
- CI/CD + IaC for deployment

## Key Design Decisions
- Broker: NATS for low-latency pub/sub and horizontal scaling
- Realtime: WebSocket for bidirectional messaging and high fan-out control
- Auth: JWT for stateless auth and easy integration
- DB: Postgres for durable data and relational modeling
- Cache: Redis for presence and fan-out state
- Prometheus as primary metrics system

## Scaling Strategy (Planned)
- Stateless gateway scaled horizontally behind a load balancer
- NATS clustering for throughput and resilience
- Fan-out handled by gateways; users routed to a gateway instance
- Redis for presence and per-user routing state
- WebSocket connection draining and backpressure on slow clients

## Performance & SLO (Draft)
- Publish p95 latency: < 50ms (in-cluster)
- WebSocket fan-out latency p95: < 100ms
- Availability: 99.9% for gateway and broker

## Risks / Constraints
- High fan-out requires careful backpressure handling
- WebSocket scaling requires connection management and resource tuning
- Cost target requires lean infra (autoscaling + minimal managed services)

## Open Questions
- Data retention and storage requirements
