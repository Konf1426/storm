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
« L'architecture étant posée, on montre maintenant le produit réel et les feedbacks visibles pour l'utilisateur. »
-->

---

<ProjectFeatureSlide />

<!--
06 - 1 min 30

Message clé :
STORM n'est pas seulement une infra benchmarkée. C'est aussi une vraie messagerie avec channels, envoi, historique et feedback utilisateur.

Points à verbaliser :
- Feedback immédiat après envoi.
- Notification de lecture visible en temps réel.
- Cette couche produit rend la démo plus crédible face au jury.

Transition :
« Maintenant qu'on a vu le produit réel, on revient aux garde-fous qualité et sécurité. »
-->

---

<QualitySecuritySlide />

<!--
07 - 2 min

Point à faire passer :
La couverture n'est pas juste un chiffre. Elle recouvre les parcours qui intéressent le jury : auth, refresh, logout, WebSocket, health checks, persistance.

Nuance :
« 80 % de couverture ne garantit pas la performance. C'est précisément pour ça qu'on enchaîne avec l'observabilité puis les tests de charge. »
-->

---

<ObservabilitySlide />

<!--
08 - 1 min 30

Faire simple :
Quelles métriques on expose, à quoi elles servent, et quels SLOs elles permettent de vérifier.

Transition :
« Maintenant qu'on peut mesurer, on passe à la méthodologie de tests de performance. »
-->

---

<PerfMethodologySlide />

<!--
09 - 2 min

Bien cadrer la slide :
Le message n'est pas « on a atteint les 100k en local ». Le message est « on a une méthode de test reproductible, de lecture des résultats et d'optimisation ».

Transition :
« Avant les chiffres, on montre rapidement l'outil réel qui a servi à générer la charge. »
-->

---

<K6ExecutionSlide />

<!--
10 - 1 min 30

Message clé :
On a vraiment utilisé k6 pour générer la charge, en local puis sur Azure, afin de garder la main sur le scénario et limiter les coûts.

Points à verbaliser :
- `scripts/perf-load.sh` pour les runs ciblés
- `scripts/storm-day-runner.sh` pour les campagnes orchestrées
- job Kubernetes k6 pour distribuer les VUs sur Azure

Transition :
« Une fois la méthode et l'outil posés, on peut lire les résultats obtenus sous charge. »
-->

---

<PerfResultsSlide />

<!--
11 - 2 min 30

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
12 - 2 min

Bien assumer la nuance :
Un crash test avec 0,36 % de succès HTTP pendant l'arrêt n'est pas un échec narratif. C'est la preuve qu'on mesure la panne et qu'on sait expliquer le rétablissement.

Transition :
« Pour passer de cette base locale à l'objectif du sujet, il faut parler de scalabilité et de budget. »
-->

---

<CloudInitialFailSlide />

<!--
13 - 2 min

Le test de départ (1 000 VUs) :
Expliquer l'échec critique initial. Ce n'était pas un problème de code STORM, mais une sous-configuration Azure et un Rate Limiter trop agressif.
- 99 % d'échecs.
- Latence Login > 4s.
Budget minimal : 0,15 $/h.

Transition :
« On a donc dû pivoter, à la fois sur le code et sur l'infrastructure. »
-->

---

<CloudOptimizationSlide />

<!--
14 - 2 min

Le cap des 5 000 VUs (optimisation) :
Détailler le pivot technique :
- Messages asynchrones (NATS libère la gateway).
- Bcrypt Cost=4 (libère le CPU).
Résultat : la latence message chute sous les 100ms. On commence à voir le potentiel.

Transition :
« Pour valider la trajectoire finale du sujet, on a sorti l'artillerie lourde. »
-->

---

<CloudSuccessSlide />

<!--
15 - 2 min

Le test ultime (10 000 VUs) :
Commenter la configuration "Ultra" : 26 vCPUs multi-famille, DB 16 Cores.
Résultat : 100 % de succès, 144ms de login médian. 210M de messages WebSocket.

Tableau de synthèse :
Montrer qu'on a testé plusieurs trajectoires (B2s vs D-Series) et qu'on sait combien coûte la performance (3,33 $/h).

Transition :
« Cette expérience cloud nous donne une vision claire pour la suite du projet. »
-->

---

<NextStepsSlide />

<!--
16 - 1 min

Présenter les axes d'amélioration comme des ouvertures stratégiques.

Transition :
« On termine par la conclusion finale. »
-->

---

<ConclusionSlide />

<!--
17 - 1 min

Synthèse finale :
STORM n'est pas qu'un repo de code, c'est une preuve de concept mesurée et capable de scaler si on y met le prix.

Passage aux questions :
« Nous sommes prêts pour vos questions. »
-->
