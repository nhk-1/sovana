# Sovana Monitoring Platform

Sovana fournit un serveur de monitoring minimaliste et un agent système écrits en Go. Le serveur expose une API REST pour recevoir des métriques, conserve un historique en mémoire thread-safe et diffuse les dernières valeurs à tous les clients WebSocket connectés. L’agent collecte périodiquement l’utilisation CPU, la mémoire et l’uptime en utilisant `/proc` puis envoie les mesures via HTTP avec un worker pool.

## Démarrage rapide

### Lancer le serveur
```bash
go run ./cmd/server
```
Le serveur écoute sur `:8080` et expose :
- `POST /metrics` pour ingérer des métriques JSON
- `GET /ws` pour recevoir les métriques en temps réel via WebSocket
- `GET /health` pour les vérifications de liveness

### Démarrer l’agent
```bash
go run ./cmd/agent -server http://localhost:8080/metrics -interval 5s -workers 4
```
L’agent échantillonne les métriques toutes les 5 secondes par défaut et pousse les données en parallèle avec un pool de workers HTTP.

## Modèle de données
```json
{
  "hostname": "example-host",
  "cpu_usage": 42.5,
  "memory_used": 1048576,
  "memory_total": 2097152,
  "uptime_seconds": 3600.5,
  "timestamp": "2024-01-02T15:04:05Z"
}
```

## Notes de concurrence
- Collecte : une goroutine dédiée lit `/proc` toutes les `interval` secondes et pousse les métriques sur un channel.
- Agent : un pool de workers consomme ce channel pour envoyer les données HTTP sans bloquer la collecte.
- Serveur : le Hub WebSocket utilise plusieurs goroutines/channels pour gérer l’inscription des clients et le broadcast des métriques.
