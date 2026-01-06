# SaaS Détection d’Abonnements (MVP)

Prototype Go qui analyse un relevé bancaire PDF pour détecter des abonnements mensuels récurrents. Frontend minimal (HTML/CSS/JS) et API REST stateless avec stockage en mémoire.

## Démarrage

```bash
# lancer l'API + frontend
cd /workspace/sovana
GO111MODULE=on go run ./cmd/server
# serveur sur http://localhost:8080
```

Options :
- `-addr` : adresse d’écoute (défaut `:8080`)
- `-static` : répertoire des assets (défaut `web`)

## API

- `GET /api/health` : santé
- `POST /api/upload` : upload PDF (champ `file` en multipart/form-data), déclenche parsing + détection. Réponse `{subscriptions: [...], count: N}`.
- `GET /api/subscriptions` : retourne la dernière détection en mémoire.

## Format PDF attendu

Le parser extrait le texte brut du PDF et tente ensuite de reconstruire des transactions quel que soit l’ordre ou la présence des colonnes :

- Détection automatique du séparateur `;` ou `,` (tabs convertis en séparateur) mais fonctionne aussi sur des lignes libres sans tableau.
- Colonnes date/libellé/montant/débit/crédit détectées par les en-têtes **ou** par inspection des valeurs (position variable acceptée).
- Dates supportées : `YYYY-MM-DD`, `DD/MM/YYYY`, `DD-MM-YYYY` et reconnues même au milieu d’une phrase.
- Montants avec `,` ou `.` (débit négatif, colonnes séparées ou montants isolés dans une phrase).
- Libellés ignorés automatiquement : `Ajout de fonds`, `Solde`, `Virement interne`.
- Les lignes invalides sont ignorées (et journalisées en debug) sans bloquer l’import.

Exemples de lignes acceptées :

```
date;libelle;debit;credit
2024-01-02;NETFLIX;13,49;
Prime Cashback;15,00;01/02/2024;
02-02-2024 Spotify 9.99 paiement
Netflix -13,49 prélèvement 01/01/2024
```

## Logique de détection

1. Normalisation du libellé (minuscules, suppression de mots courants carte/prlv/visa, nettoyage des caractères non alphanumériques).
2. Regroupement par marchand normalisé.
3. Vérification d’une périodicité mensuelle (intervalle médian entre 28 et 31 jours).
4. Vérification de montants similaires (écart ≤ max(1€, 10%)).
5. Calculs : montant moyen mensuel (arrondi 2 décimales) et coût annuel (×12).
6. Ajout d’un lien de résiliation connu si disponible (Netflix, Spotify, Apple, Amazon Prime, Adobe, Canal+).

## Frontend

- Upload PDF + bouton “Analyser”
- Affichage dynamique des abonnements détectés (montant mensuel, premières/dernières occurrences, coût annuel)
- Bouton “Comment résilier” qui ouvre la page officielle si connue

## Limitations V1

- Pas d’authentification
- Stockage en mémoire uniquement
- Périodicité mensuelle seulement
- Aucun appel bancaire ou résiliation automatique
