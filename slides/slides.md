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
01 - 1 min 30

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
02 - 30 secondes

Parcourir rapidement les sections.
Transition :
« On commence par la cible du sujet, les chiffres à atteindre et les critères d'évaluation. »
-->

---

<ObjectivesSlide />

<!--
03 - 2 min

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
04 - 2 min 30

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
05 - 2 min

Dérouler dans l'ordre :
1. Go — concurrence native et profilage intégré.
2. WebSocket — bidirectionnel, SSE retiré.
3. NATS — pub/sub ultra-léger, zéro stockage.
4. Postgres + Redis — séparation durable / éphémère.
5. Observabilité — métriques natives dès le gateway.

Transition :
<!--
05 - 2 min

Dérouler dans l'ordre :
1. Go — concurrence native et profilage intégré.
2. WebSocket — bidirectionnel, SSE retiré.
3. NATS — pub/sub ultra-léger, zéro stockage.
4. Postgres + Redis — séparation durable / éphémère.
5. Observabilité — métriques natives dès le gateway.

Transition :
« L'architecture étant posée, on passe à la sécurité, parce qu'un temps réel sans garde-fous reste fragile. »
-->

---

<QualitySecuritySlide />

<!--
06 - 2 min

Point à faire passer :
La couverture n'est pas juste un chiffre. Elle recouvre les parcours qui intéressent le jury : auth, refresh, logout, WebSocket, health checks, persistance.

Nuance :
« 80 % de couverture ne garantit pas la performance. C'est précisément pour ça qu'on enchaîne avec l'observabilité puis les tests de charge. »
-->

---

<ObservabilitySlide />

<!--
07 - 1 min 30

Faire simple :
Quelles métriques on expose, à quoi elles servent, et quels SLOs elles permettent de vérifier.

Transition :
« Maintenant qu'on peut mesurer, on passe à la méthodologie de tests de performance. »
-->

---

<PerfMethodologySlide />

<!--
08 - 2 min

Bien cadrer la slide :
Le message n'est pas « on a atteint les 100k en local ». Le message est « on a une méthode de test reproductible, de lecture des résultats et d'optimisation ».

Transition :
« On peut maintenant regarder les chiffres concrets à commenter. »
-->

---

<PerfResultsSlide />

<!--
09 - 2 min 30

Conseil oral :
Lire les trois graphiques comme une histoire : warm-up, montée en charge, puis scénario dégradé.

Phrase utile :
« Ce qui compte ici, c'est moins la valeur brute que la façon dont le système dégrade gracieusement tout en restant observable. »

Transition :
« Maintenant on passe à la résilience — c'est-à-dire ce qui se passe quand on injecte volontairement des incidents. »
-->

---

<ChaosResilienceSlide />

<!--
10 - 2 min

Bien assumer la nuance :
Un crash test avec 0,36 % de succès HTTP pendant l'arrêt n'est pas un échec narratif. C'est la preuve qu'on mesure la panne et qu'on sait expliquer le rétablissement.

Transition :
« Pour passer de cette base locale à l'objectif du sujet, il faut parler de scalabilité et de budget. »
-->

---

<CloudAnalysisSlide />

<!--
11 - 2 min

Contexte à poser immédiatement :
« Le 19 mars 2026, on a fait un vrai test cloud Azure, parce qu'Azure nous a été imposé en fin de projet. »

Fil narratif :
Le premier run ne mesure pas STORM, il révèle des "faux plafonds" de configuration :
1. MC_ absent (Azure Resource Group)
2. ACR 401 (Problème de pull d'images)
3. Rate Limiting applicatif trop strict pour un benchmark.
4. Coût CPU Bcrypt qui sature les nœuds.

Transition :
« Une fois ces verrous levés, on a pu mesurer la vraie puissance de l'architecture. »
-->

---

<CloudResultsSlide />

<!--
12 - 2 min

Le saut de performance :
1. Découplage NATS / Postgres : Publication immédiate, écriture différée.
2. Ultra scaling : 26 vCPUs, Postgres 16 vCores, 30 instances Gateway.
3. Résultat : ~100 % de succès à 10 000 VUs, latence message divisée par 30 (< 100 ms).

Point budgétaire :
Coût massif pendant le pic (3,33 $/h) mais standby très faible (0,15 $/h). Le système est élastique.

Transition :
« Ces tests massifs nous ouvrent déjà des pistes d'amélioration concrètes pour le futur. »
-->

---

<NextStepsSlide />

<!--
13 - 1 min

Présenter chaque axe d'amélioration comme une action concrète, pas comme un aveu de faiblesse.

Transition :
« On conclut en une minute. »
-->

---

<ConclusionSlide />

<!--
14 - 1 min

Finir sobrement :
« Notre valeur n'est pas d'avoir tout industrialisé, mais d'avoir pris des décisions solides, mesuré leurs effets et documenté nos limites. »

Passage aux questions :
« On est prêts à répondre sur la sécurité, la charge, la reprise après incident ou le budget. »
-->
