# Sovana Web App

Sovana is a small yet production-ready Go web application that demonstrates a clean separation between an HTTP API and a lightweight frontend. It implements a thread-safe in-memory store with JSON endpoints and a dynamic HTML/CSS/JS interface.

## Getting started

Requirements:
- Go 1.21+

Steps:
```bash
cd app
# Run the server
go run ./cmd/server
```
The server listens on `:8080` by default. Override with `ADDR=0.0.0.0:9000` as needed.

Open [http://localhost:8080](http://localhost:8080) in your browser to use the app.

## API

Base path: `/api`

- `GET /api/items` — list all items.
- `POST /api/items` — create an item. Body example:
  ```json
  {"title": "My item", "description": "Optional text"}
  ```
- `DELETE /api/items/{id}` — delete an item by ID.

Errors are returned as plain text with appropriate HTTP status codes.

## Project layout

```
app/
├── cmd/server         # Entrypoint wiring HTTP server and middleware
├── internal
│   ├── api            # HTTP handlers
│   ├── middleware     # Shared HTTP middleware (logging)
│   ├── model          # Domain models
│   └── storage        # In-memory store implementation
├── web                # Frontend assets (HTML/CSS/JS)
└── go.mod
```

This layout keeps internal packages private while making it easy to extend the project (database-backed storage, authentication, etc.).

## Notes

- All data lives in memory; restarting the server clears stored items.
- Logging middleware prints method, path, status, and latency for every request.
