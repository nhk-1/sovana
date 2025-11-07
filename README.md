# Sovana – ESN SPA

Ce dépôt contient la landing page single-page de l’ESN Sovana, construite avec React 18, Vite 5, Tailwind CSS 3 et Framer Motion. L’interface respecte la charte graphique demandée et met en avant les services, expertises et le formulaire de contact de l’entreprise.

## Démarrage

```bash
npm install
npm run dev
```

## Production

```bash
npm run build
npm run preview
```

## Structure

- `src/components`: Composants fonctionnels React organisés par section et bibliothèque UI.
- `src/assets`: Ressources SVG pour le logo et les expertises.
- `tailwind.config.js`: Configuration Tailwind (palette, typographie, breakpoints).
- `index.html`: Entrée Vite avec métadonnées SEO et chargement de la police Inter.

## Accessibilité & UX

- Navigation clavier complète, focus visible et contrastes AA.
- Défilement doux via `react-scroll` et animations subtiles avec Framer Motion.
- Formulaire de contact mock avec toast d’accusé de réception accessible.

## Licence

Projet de démonstration – droits réservés à Sovana.
