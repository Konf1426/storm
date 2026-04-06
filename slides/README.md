# STORM — Slides de Soutenance

Slides de présentation pour la soutenance MT4 du projet STORM, construites avec [Slidev](https://sli.dev).

## Prérequis

- **Node.js** ≥ 18
- **npm**

## Installation

```bash
cd slides
npm install
```

## Lancer les slides en local

```bash
npm run dev
```

Les slides s'ouvrent automatiquement sur [http://localhost:3030](http://localhost:3030).

> **Mode présentateur** : ajouter `/presenter` à l'URL pour voir les speaker notes et le timer.

## Exporter

```bash
# Export PDF
npm run export:pdf

# Export PPTX
npm run export:pptx

# Export des speaker notes
npm run export:notes
```

Les fichiers exportés se trouvent dans le dossier `exports/`.
