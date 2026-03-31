---
theme: default
title: STORM - Soutenance
colorSchema: light
mdc: true
presenter: true
transition: slide-left
lineNumbers: false
drawings:
  persist: false
---

<IntroCover />

<!--
1 min 30

Ouverture :
« Bonjour, nous vous présentons STORM — une application de messagerie temps réel pour laquelle nous avons conçu un backend capable d'absorber un pic de trafic important. »

Message clé :
L'objectif n'est pas de prétendre que le système est industrialisé à 100 %, mais de montrer une architecture cohérente, instrumentée et défendable avec des preuves.

Transition :
« On va commencer par le sommaire de la présentation. »
-->

---

<SummarySlide />

<!--
30 secondes

Parcourir rapidement les sections.
Transition :
« On commence par la cible du sujet, les chiffres à atteindre et les critères d'évaluation. »
-->

---

<ObjectivesSlide />

<!--
2 min

Insister sur les 3 chiffres clés du sujet :
100 000 connexions, 500 000 msg/s, budget ≤ 700 EUR.

Dire clairement :
« Ces chiffres servent de cap de conception. Aujourd'hui, on vous montre comment on s'en approche avec une architecture, des tests et un plan de scaling. »

Transition :
« Maintenant que la cible est posée, on passe à l'architecture réelle du repo. »
-->

---

<ArchitectureFlow />

<!--
2 min 30

Message clé :
On part volontairement d'une architecture lisible. Le jury doit comprendre en 30 secondes comment un message circule du client jusqu'à la base.

Point à verbaliser :
« Le repo actuel est full WebSocket. Le SSE a été supprimé — inutile de le mentionner sauf si on nous pose la question. »

Transition :
« L'architecture étant posée, on va maintenant justifier chaque choix technique et les compromis qu'on a faits. »
-->

---

<ArchDetailsSlide />

<!--
2 min

Dérouler dans l'ordre :
1. Go — concurrence native et profilage intégré.
2. WebSocket — bidirectionnel, SSE retiré.
3. NATS — pub/sub ultra-léger, zéro stockage.
4. Postgres + Redis — séparation durable / éphémère.
5. Observabilité — métriques natives dès le gateway.

Transition :
« Maintenant qu'on a posé les choix, on passe à la sécurité, parce qu'un temps réel sans garde-fous reste fragile. »
-->

---

<QualitySecuritySlide />

<!--
2 min

Point à faire passer :
La couverture n'est pas juste un chiffre. Elle recouvre les parcours qui intéressent le jury : auth, refresh, logout, WebSocket, health checks, persistance.

Nuance :
« 80 % de couverture ne garantit pas la performance. C'est précisément pour ça qu'on enchaîne avec l'observabilité puis les tests de charge. »
-->

---

<div class="page-shell">
  <p class="eyebrow">05 — Observabilité & SLOs</p>
  <h2 class="section-title">Impossible de parler de résilience sans mesurer ce qui se passe</h2>

  <div class="content-grid-2">
    <div class="callout-panel">
      <p class="divider-title">Métriques custom exposées</p>
      <div class="pill-grid">
        <span class="pill" v-for="metric in $storm.observability.metrics" :key="metric">{{ metric }}</span>
      </div>
    </div>
    <div class="callout-panel">
      <p class="divider-title">SLOs définis pour le Storm Day</p>
      <BulletStackList :items="$storm.observability.slos" />
    </div>
  </div>

  <SourceStrip :items="$storm.sources.observability" />
</div>

<!--
1 min 30

Faire simple :
Quelles métriques on expose, à quoi elles servent, et quels SLOs elles permettent de vérifier.

Transition :
« Maintenant qu'on peut mesurer, on passe à la méthodologie de tests de performance. »
-->

---

<div class="page-shell">
  <p class="eyebrow">07 — Méthodologie performance</p>
  <h2 class="section-title">Notre approche : charger, profiler, puis analyser ce qu'on a appris</h2>

  <div class="k6-grid">
    <div class="callout-panel">
      <p class="divider-title">k6 pour la charge</p>
      <ul class="signal-list">
        <li v-for="item in $storm.perfMethodology" :key="item.title">
          <strong>{{ item.title }}</strong> - {{ item.detail }}
        </li>
      </ul>
    </div>
    <div class="callout-panel">
      <p class="divider-title">pprof pour le diagnostic</p>
      <BulletStackList :items="$storm.perfResults.hotspotSummary" />
    </div>
  </div>

  <div class="accent-banner orange">
    <span class="accent-badge">Important</span>
    <span>
      Les chiffres présentés ensuite proviennent d'un environnement Docker Compose local. Ils
      illustrent une tendance et une méthode, pas une validation finale des 100 000 connexions.
    </span>
  </div>

  <SourceStrip :items="$storm.sources.performance" />
</div>

<!--
2 min

Bien cadrer la slide :
Le message n'est pas « on a atteint les 100k en local ». Le message est « on a une méthode de test reproductible, de lecture des résultats et d'optimisation ».

Transition :
« On peut maintenant regarder les chiffres concrets à commenter. »
-->

---

<div class="page-shell">
  <p class="eyebrow">08 — Résultats de charge</p>
  <h2 class="section-title">Résultats du Storm Day local : tenue correcte sous charge, même avec chaos</h2>

  <div class="accent-banner">
    <span class="accent-badge">Local Docker</span>
    <span>
      Pire scénario observé lors du Storm Day du 4 février : <strong>p95 = 141,5 ms</strong>,
      <strong>error rate = 0,00 %</strong>, <strong>WS connect p95 = 28,42 ms</strong>.
    </span>
  </div>

  <div class="content-grid-3">
    <div class="callout-panel">
      <p class="divider-title">Latency p95</p>
      <MetricBars :items="$storm.perfResults.latency" :max="160" />
    </div>
    <div class="callout-panel">
      <p class="divider-title">Debit HTTP</p>
      <MetricBars :items="$storm.perfResults.reqRate" :max="2600" />
    </div>
    <div class="callout-panel">
      <p class="divider-title">Debit WS</p>
      <MetricBars :items="$storm.perfResults.wsRate" :max="175000" />
    </div>
  </div>
  <div class="chip-row takeaways-row">
    <span class="chip">Latence contenue même sous charge + chaos</span>
    <span class="chip">Le run local de référence ne s'effondre pas</span>
    <span class="chip">Lecture croisée k6 + pprof + persistance</span>
  </div>
  <SourceStrip :items="$storm.sources.performance" />
</div>

<!--
2 min 30

Conseil oral :
Lire les trois graphiques comme une histoire : warm-up, montée en charge, puis scénario dégradé.

Phrase utile :
« Ce qui compte ici, c'est moins la valeur brute que la façon dont le système dégrade gracieusement tout en restant observable. »

Transition :
« Maintenant on passe à la résilience — c'est-à-dire ce qui se passe quand on injecte volontairement des incidents. »
-->

---

<div class="page-shell">
  <p class="eyebrow">09 — Chaos & résilience</p>
  <h2 class="section-title">Le chaos engineering donne une lecture honnête de la robustesse</h2>

  <div class="signal-grid">
    <article class="signal-card" v-for="scenario in $storm.chaosResults.scenarios" :key="scenario.title">
      <p class="card-title">{{ scenario.title }}</p>
      <p class="signal-detail">{{ scenario.detail }}</p>
    </article>
  </div>

  <div class="callout-panel">
    <p class="divider-title">Ce qu'on valide</p>
    <ul class="signal-list">
      <li v-click>Le service se rétablit après coupure NATS ou redémarrage du gateway</li>
      <li v-click>Le worker pool absorbe un ralentissement Postgres</li>
      <li v-click>Le projet produit des post-mortems structurés</li>
    </ul>
  </div>

  <SourceStrip :items="$storm.sources.chaos" />
</div>

<!--
2 min

Bien assumer la nuance :
Un crash test avec 0,36 % de succès HTTP pendant l'arrêt n'est pas un échec narratif. C'est la preuve qu'on mesure la panne et qu'on sait expliquer le rétablissement.

Transition :
« Pour passer de cette base locale à l'objectif du sujet, il faut parler de scalabilité et de budget. »
-->

---

<div class="page-shell">
  <p class="eyebrow">10 — Scalabilité & cloud</p>
  <h2 class="section-title">Budget cible AWS et preuve de faisabilité Azure / AKS</h2>
  <p class="section-copy">
    Le sujet mentionne AWS ou équivalent. On distingue donc deux axes :
    un budget cible pour raisonner sur le plafond de 700 EUR, et une architecture Azure
    documentée dans le repo pour démontrer la faisabilité d'un test distribué à plus grande échelle.
  </p>

  <div class="truth-grid">
    <div class="truth-card">
      <h3>Budget cible</h3>
      <p class="muted">
        En lecture AWS, un setup minimal reste sous les 100 USD, un setup intermédiaire
        se situe autour de 250-300 USD, et la très haute charge dépasse le plafond de 700 EUR.
      </p>
    </div>
    <div class="truth-card">
      <h3>Preuve Azure</h3>
      <p class="muted">
        Le repo documente un déploiement AKS avec App Gateway v2, autoscaling, Redis, Postgres et
        k6 distribué pour cibler les 100 000 connexions WebSocket.
      </p>
    </div>
  </div>

  <div class="budget-grid">
    <article class="budget-card" v-for="scenario in $storm.budgetScenarios" :key="scenario.name">
      <h3>{{ scenario.name }}</h3>
      <p class="stat-value">{{ scenario.estimate }}</p>
      <p class="signal-detail">{{ scenario.note }}</p>
    </article>
  </div>

  <SourceStrip :items="$storm.sources.budget" />
</div>

<!--
2 min

Formule utile :
« Le budget n'est pas un chiffre absolu — c'est une zone de compromis entre débit, haute disponibilité et niveau de service managé. »

Bien distinguer :
AWS = référence budgétaire du sujet.
Azure = preuve d'implémentation cloud équivalente dans le repo.

Transition :
« Avant de conclure, on va prévoir une slide dédiée à la démo. »
-->

---

<div class="page-shell">
  <p class="eyebrow">11 — Démo live / fallback</p>
  <h2 class="section-title">Démo courte et maîtrisée, avec un plan B prêt si le live ne fonctionne pas</h2>

  <div class="truth-grid">
    <div class="demo-card">
      <h3>Script de démo</h3>
      <ol class="signal-list">
        <li v-for="step in $storm.liveDemo" :key="step">{{ step }}</li>
      </ol>
    </div>
    <div class="demo-card">
      <h3>Fallback préparé</h3>
      <ul class="truth-list">
        <li v-click>Afficher les résultats Storm Day déjà versionnés dans le repo</li>
        <li v-click>Montrer les scripts <span class="mono">perf-load.sh</span> et <span class="mono">storm-day-runner.sh</span></li>
        <li v-click>Basculer sur les métriques et les logs plutôt que de forcer une démo instable</li>
      </ul>
    </div>
  </div>


  <SourceStrip :items="$storm.sources.performance" />
</div>

<!--
1 min

Dire clairement au jury :
« On a un plan de démo, mais on ne fait pas reposer notre note sur un clic ou une connexion réseau. »

Transition :
« On termine par les limites honnêtes et ce qu'il reste à faire pour monter d'un cran. »
-->

---

<div class="page-shell">
  <p class="eyebrow">12 — Prochaines étapes</p>
  <h2 class="section-title">Ce qu'il reste à construire pour aller plus loin</h2>

  <div class="signal-grid">
    <SignalCardsStack :items="$storm.nextSteps" />
  </div>

  <SourceStrip :items="$storm.sources.budget" />
</div>

<!--
1 min

Présenter chaque axe d'amélioration comme une action concrète, pas comme un aveu de faiblesse.

Transition :
« On conclut en une minute. »
-->

---

<div class="page-shell">
  <p class="eyebrow">13 — Conclusion</p>
  <h2 class="section-title">STORM est déjà défendable comme système, pas seulement comme maquette</h2>

  <div class="signal-grid">
    <article class="signal-card">
      <p class="card-title">Architecture claire</p>
      <p class="signal-detail">
        Gateway Go, WebSocket, NATS, persistance et observabilité : chaque couche a un rôle net et démontrable.
      </p>
    </article>
    <article class="signal-card">
      <p class="card-title">Qualité mesurable</p>
      <p class="signal-detail">
        Couverture > 80 %, tests ciblés, scans CI, endpoints de santé et métriques Prometheus.
      </p>
    </article>
    <article class="signal-card">
      <p class="card-title">Performance lisible</p>
      <p class="signal-detail">
        Charge locale, profilage, chaos engineering et limites documentées pour guider la suite.
      </p>
    </article>
  </div>

  <div class="callout-panel">
    <p class="finale-line">
      Notre message final : STORM n'est pas « fini », mais il est déjà instrumenté, testé, sécurisé
      et expliqué de façon honnête.
    </p>
  </div>

  <SourceStrip :items="$storm.sources.architecture" />
</div>

<!--
1 min

Finir sobrement :
« Notre valeur n'est pas d'avoir tout industrialisé, mais d'avoir pris des décisions solides, mesuré leurs effets et documenté nos limites. »

Passage aux questions :
« On est prêts à répondre sur la sécurité, la charge, la reprise après incident ou le budget. »
-->

---

<div class="page-shell">
  <span class="backup-tag">Backup</span>
  <p class="eyebrow">14 — Storm Day timeline</p>
  <h2 class="section-title">Timeline de référence pour répondre aux questions du jury</h2>

  <div class="content-grid-2">
    <div class="callout-panel">
      <p class="divider-title">Séquence type du Storm Day</p>
      <ul class="signal-list">
        <li>Warm-up 10 min</li>
        <li>Growth spike 15 min</li>
        <li>Incident 1: pause NATS ou restart gateway</li>
        <li>Recovery 10 min</li>
        <li>Incident 2: injection de latence</li>
        <li>Growth spike 2 puis cooldown</li>
      </ul>
    </div>
    <div class="callout-panel">
      <p class="divider-title">Ce qu'on capture</p>
      <ul class="signal-list">
        <li>Logs k6</li>
        <li>Profils pprof CPU et heap</li>
        <li>Etat des SLOs</li>
        <li>Chronologie des incidents et post-mortem</li>
      </ul>
    </div>
  </div>

  <SourceStrip :items="$storm.sources.performance" />
</div>

<!--
Backup uniquement.

Utiliser cette slide si le jury demande :
« Comment s'organise exactement votre Storm Day ? »
-->

---

<div class="page-shell">
  <span class="backup-tag">Backup</span>
  <p class="eyebrow">15 — Budget détaillé</p>
  <h2 class="section-title">Lecture détaillée des scénarios de coût</h2>

  <div class="budget-grid">
    <article class="budget-card" v-for="scenario in $storm.budgetScenarios" :key="scenario.name">
      <h3>{{ scenario.name }}</h3>
      <p class="stat-value">{{ scenario.estimate }}</p>
      <p class="signal-detail">{{ scenario.note }}</p>
    </article>
  </div>

  <div class="callout-panel">
    <p class="divider-title">Point de vigilance</p>
    <p class="meta-note">
      Ces chiffres restent des estimations. Le plafond est très sensible au nombre d'instances, au
      volume de transfert et au niveau de service managé retenu.
    </p>
  </div>

  <SourceStrip :items="$storm.sources.budget" />
</div>

<!--
Backup uniquement.

Utiliser cette slide si le jury veut approfondir la viabilité économique ou le choix AWS / Azure.
-->

---

<div class="page-shell">
  <span class="backup-tag">Backup</span>
  <p class="eyebrow">16 — Auth & API flow</p>
  <h2 class="section-title">Parcours d'authentification et flux de message complet</h2>

  <div class="content-grid-2">
    <div class="callout-panel">
      <p class="divider-title">Auth</p>
      <ul class="signal-list">
        <li><span class="mono">POST /auth/login</span> émet les cookies de session</li>
        <li><span class="mono">POST /auth/refresh</span> renouvelle l'access token</li>
        <li><span class="mono">POST /auth/logout</span> invalide la session</li>
        <li><span class="mono">GET /auth/me</span> retourne l'utilisateur courant</li>
      </ul>
    </div>
    <div class="callout-panel">
      <p class="divider-title">Temps reel</p>
      <ul class="signal-list">
        <li><span class="mono">GET /ws</span> ouvre la connexion WebSocket protégée par JWT</li>
        <li><span class="mono">POST /publish</span> pousse un événement vers NATS</li>
        <li>Le gateway enregistre et rediffuse, le service messages consomme et persiste</li>
        <li>Les endpoints <span class="mono">/healthz</span> et <span class="mono">/ping-nats</span> servent au diagnostic rapide</li>
      </ul>
    </div>
  </div>

  <SourceStrip :items="$storm.sources.api" />
</div>

<!--
Backup uniquement.

Utile si la discussion porte sur le détail des endpoints ou le choix d'authentification JWT.
-->
