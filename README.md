# SaaS Détection d’Abonnements (MVP)

Prototype Go qui analyse un relevé bancaire CSV pour détecter des abonnements mensuels récurrents. Frontend minimal (HTML/CSS/JS) et API REST stateless avec stockage en mémoire.

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
- `POST /api/upload` : upload CSV (champ `file` en multipart/form-data), déclenche parsing + détection. Réponse `{subscriptions: [...], count: N}`.
- `GET /api/subscriptions` : retourne la dernière détection en mémoire.

## Format CSV attendu

Colonnes : `date, libellé, montant`. En-tête optionnelle.

- Date : `YYYY-MM-DD`, `DD/MM/YYYY`, `MM/DD/YYYY` ou `YYYY/MM/DD`
- Montant : séparateur `.` ou `,`, signe `-` pour un débit.
- Exemple :

```csv
2024-01-02,NETFLIX,-13.49
2024-02-02,NETFLIX,-13.49
2024-03-02,NETFLIX,-13.49
2024-01-10,SPOTIFY,-9.99
2024-02-10,SPOTIFY,-9.99
2024-03-10,SPOTIFY,-9.99
```

## Logique de détection

1. Normalisation du libellé (minuscules, suppression de mots courants carte/prlv/visa, nettoyage des caractères non alphanumériques).
2. Regroupement par marchand normalisé.
3. Vérification d’une périodicité mensuelle (intervalle médian entre 28 et 31 jours).
4. Vérification de montants similaires (écart ≤ max(1€, 10%)).
5. Calculs : montant moyen mensuel (arrondi 2 décimales) et coût annuel (×12).
6. Ajout d’un lien de résiliation connu si disponible (Netflix, Spotify, Apple, Amazon Prime, Adobe, Canal+).

## Frontend

- Upload CSV + bouton “Analyser”
- Affichage dynamique des abonnements détectés (montant mensuel, premières/dernières occurrences, coût annuel)
- Bouton “Comment résilier” qui ouvre la page officielle si connue

## Limitations V1

- Pas d’authentification
- Stockage en mémoire uniquement
- Périodicité mensuelle seulement
- Aucun appel bancaire ou résiliation automatique
