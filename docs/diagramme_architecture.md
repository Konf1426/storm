# System Architecture

## Overview
The system follows a microservices architecture using **Go** for backend services, **Vue 3** for the frontend, and **NATS** for event-driven communication. Infrastructure is containerized via Docker and orchestrated (e.g., K8s) for high availability.

## Components

### CI/CD (GitHub)
- **Source Control**: GitHub.
- **CI/CD**: GitHub Actions.
- **Artifacts**: Docker Images stored in GitHub Container Registry (GHCR) or Docker Hub.
- **Workflow**:
    1.  **Push**: Developer pushes code to main/staging branch.
    2.  **Test**: Actions run unit tests and linters.
    3.  **Build**: Actions build Docker images.
    4.  **Deploy**: Updates infrastructure (e.g., via SSH, GitOps, or Webhook).

### Frontend
- **Tech Stack**: Vue 3, Vite, TailwindCSS.
- **Role**: User Interface.
- **Communication**: Talks to the **Gateway** service via HTTPS (REST & WebSocket).

### Backend Services
1.  **Gateway** (`services/gateway`) -> **Autoscaling (HPA)**
    -   **Role**: API Gateway. Entry point for the frontend.
    -   **Responsibilities**: Authentication, routing, aggregation.
    -   **Dependencies**: Postgres, Redis, NATS.
    -   **Scalability**: Stateless service, scales horizontally based on CPU/Memory and active WS connections.

2.  **Messages** (`services/messages`) -> **Autoscaling (HPA)**
    -   **Role**: Event Consumer / specialized service.
    -   **Responsibilities**: Handling message-related logic, listening to `storm.events` on NATS.
    -   **Dependencies**: NATS, Postgres.
    -   **Processing**: Uses JetStream Pull Consumers for backpressure control.

3.  **Auth Service** (Proposed/External)
    -   **Role**: Identity Provider (IdP).
    -   **Tech**: Keycloak, Auth0, or custom Go service issuing JWTs (RS256).

### Infrastructure
-   **NATS JetStream (Cluster)**:
    -   **Role**: Event Bus / Message Broker.
    -   **Configuration**: 3-node cluster for HA and stream replication (R=3).
    -   **Features**: Durable streams, subject partitioning.
-   **PostgreSQL**:
    -   **Role**: Primary Relational Database.
    -   **Usage**: Persistent storage (users, channels, messages).
-   **Redis (Cluster)**:
    -   **Role**: Distributed Cache / Ratelimiter.
    -   **Usage**: Session cache, API rate limiting (Sliding Window), ephemeral state.

### Observability
-   **Prometheus**: Metrics collection (scrape intervals: 15s).
-   **Grafana**: Dashboards (SLO tracking).
-   **Targets (SLOs)**:
    -   **Availability**: 99.9%.
    -   **Latency (p99)**: < 100ms for API calls.
    -   **Throughput**: Support 500k msg/s peak.

## Scalability & Performance Strategy
1.  **Horizontal Pod Autoscaling (HPA)**: Both Gateway and Worker services scale based on load.
2.  **NATS Partitioning**: Message streams partitioned by `channel_id` to distribute load across consumers.
3.  **Batching**: Producers batch messages (up to 1ms or 100 msgs) to reduce syscall overhead.
4.  **Connection Pooling**: Aggressive pooling for Postgres and Redis connections.

## Security
-   **Authentication**: JWT Token passed in headers (Bearer). Validated at Gateway.
-   **Rate Limiting**: Redis-backed limiter per IP/User to prevent abuse/DDoS.
-   **TLS**: All external traffic encrypted (HTTPS/WSS).

## Testing Strategy
To ensure quality (Requirement: >80% coverage) and resilience (Storm Day):
1.  **Unit Logic (Go)**:
    -   Use standard `testing` package + mocks (`mockery`).
    -   Focus on business logic coverage.
2.  **Integration (Testcontainers)**:
    -   Spin up ephemeral NATS/Postgres/Redis for handler testing.
    -   Verify correct message publishing and persistence.
3.  **Load Testing (k6)**:
    -   Simulate 100k concurrent WebSocket connections.
    -   Verify latency SLO (<100ms) under load.
4.  **Chaos Engineering**:
    -   Randomly kill Gateway/Messages pods to verify recovery.
    -   Simulate network partitions in NATS cluster.

## Diagram (Mermaid)

```mermaid
graph TD
    %% Styles
    classDef ci fill:#D6EAF8,stroke:#2E86C1,stroke-width:2px,color:black;
    classDef testing fill:#FCF3CF,stroke:#F1C40F,stroke-width:2px,color:black;
    classDef client fill:#E5E8E8,stroke:#7F8C8D,stroke-width:2px,color:black;
    classDef service fill:#D5F5E3,stroke:#27AE60,stroke-width:2px,color:black;
    classDef db fill:#FAD7A0,stroke:#E67E22,stroke-width:2px,color:black;
    classDef monitor fill:#E8DAEF,stroke:#8E44AD,stroke-width:2px,color:black;

    %% CI/CD Pipeline
    subgraph "CI/CD Pipeline"
        direction LR
        GitHub[("<img src='assets/icons/github.png' height='45' style='width:auto; object-fit:contain;' /> GitHub<br/>(Source Control)")]:::ci
        Actions[("<img src='assets/icons/actions.svg' height='45' style='width:auto; object-fit:contain;' /> Actions<br/>(Build & Test)")]:::ci
        Registry[("<img src='assets/icons/docker.png' height='45' style='width:auto; object-fit:contain;' /> Registry<br/>(Container Store)")]:::ci
        
        GitHub --> Actions --> Registry
    end
    
    %% Quality Assurance
    subgraph "Quality Assurance (Testing)"
        direction LR
        UnitTests["<img src='assets/icons/go.png' height='45' style='width:auto; object-fit:contain;' /> Unit Tests<br/>(Go/Mockery)"]:::testing
        Integration["<img src='assets/icons/docker.png' height='45' style='width:auto; object-fit:contain;' /> Integration<br/>(Testcontainers)"]:::testing
        LoadTesting["<img src='assets/icons/k6.png' height='45' style='width:auto; object-fit:contain;' /> Load Testing<br/>(100k conns)"]:::testing
        Chaos["<img src='assets/icons/chaos-mesh.png' height='45' style='width:auto; object-fit:contain;' /> Chaos Mesh<br/>(Resilience)"]:::testing
    end

    %% External Access
    subgraph "External Access"
        direction LR
        Client["<img src='assets/icons/vue.svg' height='45' style='width:auto; object-fit:contain;' /> Client App<br/>(Vue 3 / Vite)"]:::client
        LB["<img src='assets/icons/nginx.svg' height='45' style='width:auto; object-fit:contain;' /> Ingress LB<br/>(Nginx / Traefik)"]:::client
        
        Client --> LB
    end

    %% Backend Services
    subgraph "Backend Services"
        direction LR
        Auth["<img src='assets/icons/keycloak.png' height='45' style='width:auto; object-fit:contain;' /> Auth Service<br/>(Keycloak/OIDC)"]:::service
        Gateway["<img src='assets/icons/go.png' height='45' style='width:auto; object-fit:contain;' /> API Gateway<br/>(Go - HPA)"]:::service
        Messages["<img src='assets/icons/go.png' height='45' style='width:auto; object-fit:contain;' /> Messages Service<br/>(Go - HPA)"]:::service
    end

    %% Infrastructure
    subgraph "Infrastructure"
        direction LR
        NATS[("<img src='assets/icons/nats.png' height='45' style='width:auto; object-fit:contain;' /> NATS Cluster<br/>(JetStream - 3 Nodes)")]:::db
        Postgres[("<img src='assets/icons/postgres.svg' height='45' style='width:auto; object-fit:contain;' /> PostgreSQL<br/>(Primary + Replica)")]:::db
        Redis[("<img src='assets/icons/redis.svg' height='45' style='width:auto; object-fit:contain;' /> Redis Cluster<br/>(Cache/RateLimit)")]:::db
    end

    %% Observability
    subgraph "Observability"
        direction LR
        Prometheus["<img src='assets/icons/prometheus.svg' height='45' style='width:auto; object-fit:contain;' /> Prometheus<br/>(Metrics Scrape)"]:::monitor
        Grafana["<img src='assets/icons/grafana.svg' height='45' style='width:auto; object-fit:contain;' /> Grafana<br/>(Dashboards / SLOs)"]:::monitor
        
        Prometheus --> Grafana
    end

    %% Cross-Graph Connections
    GitHub -- "Trigger" --> Actions
    Actions -- "Run" --> UnitTests
    Actions -- "Run" --> Integration
    Actions -- "Push Images" --> Registry
    
    Registry -.-> Gateway
    Registry -.-> Messages
    
    LoadTesting -- "Stress Test" --> LB
    Chaos -.-> Gateway
    Chaos -.-> Messages
    
    Client -- "HTTPS / WSS" --> LB
    LB --> Gateway
    Gateway -- "Auth Check" --> Auth
    Gateway -- "Read/Write" --> Postgres
    Gateway -- "Rate Limit" --> Redis
    Gateway -- "Pub Events (Batch)" --> NATS
    NATS -- "Sub (Pull/Queue)" --> Messages
    Messages -- "Persist" --> Postgres
    
    Prometheus -.-> Gateway
    Prometheus -.-> Messages
    
    %% Layout Helpers
    Registry ~~~ Client
    LB ~~~ Auth
    Messages ~~~ NATS
    Redis ~~~ Prometheus
```

## Draw.io Import
You can copy the Mermaid code above and paste it into Draw.io (Arrange > Insert > Advanced > Mermaid) to get a starting diagram.
