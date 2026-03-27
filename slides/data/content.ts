export type Tone = "cyan" | "orange" | "slate"

export interface Objective {
  label: string
  value: string
  note: string
  tone?: Tone
}

export interface Decision {
  title: string
  detail: string
  why: string
}

export interface Signal {
  title: string
  detail: string
}

export interface MetricItem {
  label: string
  value: number
  display: string
  meta?: string
  tone?: Tone
}

export interface BudgetScenario {
  name: string
  estimate: string
  note: string
}

export const deckMeta = {
  title: "STORM",
  subtitle: "Backend de messagerie temps réel distribué",
  oralWindow: "25 min de présentation + 5 min de questions",
  teamLabel: "Équipe STORM — noms à compléter",
  positioning:
    "Concevoir le backend d'une plateforme de messagerie comparable à Slack ou Discord, capable d'absorber un pic de trafic viral tout en restant observable, testable et défendable à l'oral.",
}

export const objectives: Objective[] = [
  {
    label: "Connexions simultanées",
    value: "100 000",
    note: "cible de dimensionnement fixée par le sujet",
    tone: "cyan",
  },
  {
    label: "Débit cible",
    value: "500 000 msg/s",
    note: "cap de conception pour la scalabilité horizontale",
    tone: "orange",
  },
  {
    label: "Budget plafond",
    value: "<= 700 EUR",
    note: "enveloppe mensuelle AWS ou équivalent cloud",
    tone: "slate",
  },
]

export const evaluationCriteria: Signal[] = [
  { title: "Architecture & code", detail: "35 %" },
  { title: "Tests & sécurité", detail: "20 % — couverture > 80 %" },
  { title: "Storm Day & SLOs", detail: "25 %" },
  { title: "Documentation", detail: "10 %" },
  { title: "Collaboration", detail: "10 %" },
]

export const architecture = {
  edge: [
    "Frontend Vue 3 / Vite servant de démonstrateur fonctionnel",
    "Gateway Go : point d'entrée unique pour le HTTP et le WebSocket",
  ],
  core: [
    "NATS comme bus de messages temps réel inter-services",
    "Service messages en Go : consommation NATS et logique de persistance asynchrone",
  ],
  state: [
    "Postgres pour la persistance durable (utilisateurs, canaux, messages)",
    "Redis pour la présence en ligne et l'état éphémère des sessions",
    "Prometheus + Grafana pour la collecte de métriques et les dashboards",
  ],
  flow: [
    {
      title: "1. Authentification",
      detail: "Login, émission JWT, cookie HttpOnly et renouvellement transparent via refresh token",
    },
    {
      title: "2. Temps réel",
      detail: "Connexion WebSocket sur /ws, publication et diffusion des événements via NATS",
    },
    {
      title: "3. Persistance",
      detail: "Worker pool asynchrone pour l'écriture des messages et des tokens en arrière-plan",
    },
    {
      title: "4. Supervision",
      detail: "Métriques Prometheus exposées par le gateway, dashboards Grafana et endpoints de santé",
    },
  ],
}

export const technicalChoices: Decision[] = [
  {
    title: "Go",
    detail: "Empreinte mémoire réduite, goroutines natives et profilage intégré via pprof.",
    why: "Adapté au temps réel à forte concurrence et au diagnostic de performance sous charge.",
  },
  {
    title: "WebSocket",
    detail: "Canal bidirectionnel persistant exposé sur /ws, remplaçant l'ancien SSE supprimé du repo.",
    why: "Communication full-duplex indispensable pour une messagerie interactive.",
  },
  {
    title: "NATS",
    detail: "Broker pub/sub ultra-léger entre le gateway et le service messages.",
    why: "Latence sub-milliseconde, zéro dépendance de stockage, simple à opérer.",
  },
  {
    title: "Postgres + Redis",
    detail: "Postgres pour les données durables, Redis pour l'état volatile des connexions.",
    why: "Séparation nette des responsabilités : persistance fiable d'un côté, accès rapide de l'autre.",
  },
  {
    title: "Prometheus + Grafana",
    detail: "Métriques custom exposées par le gateway, scrapées et visualisées en temps réel.",
    why: "Permet de défendre les SLOs à l'oral avec des données concrètes, pas des estimations.",
  },
]

export const securityControls: Signal[] = [
  {
    title: "JWT stateless",
    detail: "Access token de 15 min et refresh token de 24 h, signés côté gateway avec rotation automatique.",
  },
  {
    title: "Cookies HttpOnly",
    detail: "Le refresh token n'est jamais exposé au JavaScript client — protégé en cookie HttpOnly sécurisé.",
  },
  {
    title: "Rate limiting",
    detail: "5 requêtes/min par IP avec burst de 3 sur /auth/* — freine efficacement le brute-force.",
  },
  {
    title: "Headers de sécurité",
    detail: "CSP, HSTS, X-Frame-Options, no-referrer, COOP et cache-control configurés au niveau du gateway.",
  },
  {
    title: "Scans CI automatisés",
    detail: "Gosec analyse le code Go et Trivy scanne les images Docker à chaque push dans le pipeline.",
  },
]

export const testEvidence = {
  coverage: [
    {
      label: "Gateway",
      value: "80.4 %",
      note: "handlers HTTP, auth, WebSocket et worker pool",
      tone: "cyan" as Tone,
    },
    {
      label: "Messages",
      value: "83.6 %",
      note: "consommation NATS et persistance asynchrone",
      tone: "orange" as Tone,
    },
  ],
  suites: [
    "Tests unitaires Go couvrant les deux services (gateway + messages)",
    "Parcours complets HTTP : login, refresh, logout, healthz, ping-nats",
    "Tests WebSocket : connexion, émission, réception et déconnexion propre",
    "Tests du worker pool et de la persistance asynchrone en base",
    "Pipeline GitHub Actions : couverture + scans sécurité à chaque commit",
  ],
  caveats: [
    "Les scans Gosec/Trivy n'étaient pas activés sur l'environnement local au moment du rapport du 4 février — corrigé depuis dans le pipeline CI.",
    "La couverture valide le comportement fonctionnel ; elle ne constitue pas à elle seule une preuve de tenue à 100 000 connexions.",
  ],
}

export const observability = {
  metrics: [
    "storm_active_websockets",
    "storm_auth_rate_limited_total",
    "storm_nats_publish_duration_seconds",
    "storm_save_queue_length",
  ],
  slos: [
    "Availability ≥ 99,5 % pendant le Storm Day",
    "Latence de publication p95 < 200 ms",
    "Connexion WebSocket p95 < 200 ms",
    "Taux d'erreur HTTP < 1 %",
  ],
  tools: [
    "Prometheus pour le scrape des métriques exposées",
    "Grafana pour la visualisation et le suivi des SLOs",
    "/healthz et /ping-nats pour les vérifications de santé rapides",
  ],
}

export const perfMethodology: Signal[] = [
  {
    title: "k6",
    detail: "Génération de charge combinée HTTP + WebSocket avec scénarios progressifs : warm-up, spikes et injection de chaos.",
  },
  {
    title: "pprof",
    detail: "Captures CPU et heap en conditions de charge pour identifier les hotspots : syscall, routeur, parsing JWT, écriture WebSocket.",
  },
  {
    title: "Environnement local",
    detail: "Tous les chiffres présentés proviennent de runs Docker Compose et sont étiquetés comme tels — pas de prétention à une validation cloud finale.",
  },
]

export const perfResults = {
  latency: <MetricItem[]>[
    {
      label: "Warm-up",
      value: 6.3,
      display: "6,3 ms",
      meta: "50 VUs HTTP + 50 VUs WS pendant 2 min",
      tone: "cyan",
    },
    {
      label: "Spike 1",
      value: 81.31,
      display: "81,3 ms",
      meta: "150 VUs HTTP + 150 VUs WS pendant 3 min",
      tone: "slate",
    },
    {
      label: "Spike 2 + chaos",
      value: 141.5,
      display: "141,5 ms",
      meta: "200 VUs HTTP + 200 VUs WS pendant 3 min avec injection de pannes",
      tone: "orange",
    },
  ],
  reqRate: <MetricItem[]>[
    {
      label: "Warm-up",
      value: 923,
      display: "923 req/s",
      tone: "cyan",
    },
    {
      label: "Spike 1",
      value: 2501.5,
      display: "2 502 req/s",
      tone: "slate",
    },
    {
      label: "Spike 2 + chaos",
      value: 2007.7,
      display: "2 008 req/s",
      tone: "orange",
    },
  ],
  wsRate: <MetricItem[]>[
    {
      label: "Warm-up",
      value: 23408,
      display: "23 408 msg/s",
      tone: "cyan",
    },
    {
      label: "Spike 1",
      value: 144984.9,
      display: "144 985 msg/s",
      tone: "slate",
    },
    {
      label: "Spike 2 + chaos",
      value: 171438.9,
      display: "171 439 msg/s",
      tone: "orange",
    },
  ],
  hotspotSummary: [
    "Gateway CPU : syscall, driver Postgres (pgx), routeur Chi, parsing JWT, écriture WebSocket",
    "Gateway heap : compress/flate, bufio, allocations goroutines, chemins Redis",
    "Messages CPU/heap : écritures HTTP et compression dominent le profil",
  ],
}

export const chaosResults = {
  scenarios: [
    {
      title: "Pause NATS (10 s)",
      detail: "Coupure volontaire du broker pendant 10 secondes sous charge — le système bufferise puis reprend sans erreur visible côté k6.",
    },
    {
      title: "Injection de latence",
      detail: "100 ms ± 20 ms ajoutés artificiellement — la latence p95 monte mais le service ne rompt pas.",
    },
    {
      title: "Base de données lente",
      detail: "500 ms de délai artificiel sur Postgres — le worker pool absorbe la pression et le système reste stable.",
    },
    {
      title: "Crash du gateway",
      detail: "Arrêt brutal du process pendant 60 s — le conteneur redémarre automatiquement, mais le taux de succès HTTP chute à 0,36 % pendant l'interruption.",
    },
  ],
  learnings: [
    "La résilience est fonctionnelle, mais une reprise immédiate sous forte charge peut faire remonter des erreurs de cohérence qu'il faut surveiller.",
    "Le calcul formel de l'availability et la capture systématique des preuves Prometheus restent des axes d'amélioration identifiés.",
  ],
}

export const budgetScenarios: BudgetScenario[] = [
  {
    name: "AWS minimal",
    estimate: "~86–90 USD/mois",
    note: "1 instance EC2 + ALB + RDS micro + ElastiCache micro — hors transfert de données.",
  },
  {
    name: "AWS intermédiaire",
    estimate: "~250–300 USD/mois",
    note: "2-3 instances, base et cache dimensionnés pour un trafic modéré avec marge de scaling.",
  },
  {
    name: "AWS haute charge",
    estimate: "> 700 USD/mois",
    note: "Au-delà de 5 instances avec base et cache dimensionnés pour le pic — le plafond du sujet est dépassé.",
  },
  {
    name: "Preuve Azure AKS",
    estimate: "~255 EUR/mois",
    note: "Architecture de test AKS documentée dans le repo — un test ponctuel de quelques heures coûte entre 5 et 15 EUR.",
  },
]

export const nextSteps: Signal[] = [
  {
    title: "Formaliser les mesures",
    detail: "Calculer l'availability de manière automatique, exporter les dashboards Grafana et versionner les captures.",
  },
  {
    title: "Durcir l'infrastructure cloud",
    detail: "TLS de bout en bout, monitoring managé, politiques réseau et validation distribuée du seuil 100k.",
  },
  {
    title: "Fiabiliser la reprise",
    detail: "Améliorer le comportement après crash sous charge et garantir la cohérence de la persistance asynchrone.",
  },
]

export const liveDemo = [
  "Login sur le frontend, ouverture d'un canal et établissement de la connexion WebSocket.",
  "Envoi d'un message qui traverse le pipeline complet : gateway → NATS → service messages → persistance Postgres.",
  "Si le live n'est pas disponible : présentation des résultats Storm Day versionnés et de la méthodologie k6.",
]

export const sources = {
  goals: ["docs/consignes.md", "docs/project_storm_requirements.md"],
  architecture: ["services/gateway/server.go", "services/messages/server.go", "frontend/src/App.vue"],
  security: ["services/gateway/server.go", ".github/workflows/ci.yml", "project_report.md"],
  tests: ["docs/test-results.md", "services/gateway/server_test.go", "services/gateway/server_additional_test.go"],
  observability: ["services/gateway/server.go", "infra/docker/grafana/dashboards/storm-overview.json"],
  performance: ["docs/performance.md", "docs/storm-day-results.md"],
  chaos: ["docs/chaos.md", "docs/storm-day-results-20260224.md", "docs/post-mortem-20260204.md"],
  budget: ["docs/budget.md", "docs/azure-deployment.md"],
  api: ["docs/api/openapi.yml", "services/gateway/server.go"],
}
