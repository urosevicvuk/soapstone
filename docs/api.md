# API

Messages and errors are JSON. An error looks like `{"error": "…"}`.

A message:

```json
{"id": 1, "author": "Solaire", "text": "Praise the Sun", "createdAt": "2026-10-12T10:00:00Z"}
```

`id` starts at 1 and only grows, and `createdAt` is in UTC. Leading and trailing spaces are trimmed before validation: `author` has 1–40 characters, `text` 1–280.

| Method and path | Response |
| --- | --- |
| `GET /` | HTML page with messages and a form |
| `GET /api/messages?limit=50` | 200, list from newest; `limit` from 1 to 200, otherwise 400 |
| `POST /api/messages` | body `{"author": "…", "text": "…"}`; 201 with the message, or 400 |
| `GET /api/messages/{id}` | 200 or 404 |
| `DELETE /api/messages/{id}` | 204 with no body, or 404 |
| `GET /healthz` | 200, `ok` |

Examples:

```bash
curl -s localhost:8080/api/messages | jq
curl -s -X POST localhost:8080/api/messages \
  -H 'Content-Type: application/json' \
  -d '{"author":"Ana","text":"Try jumping"}' | jq
curl -s -X DELETE localhost:8080/api/messages/1
```

## Configuration

All settings are env variables. `make run` also loads `.env`, if it exists.

| Variable | Default | Meaning |
| --- | --- | --- |
| `PORT` | `8080` | server port |
| `APP_COLOR` | `steelblue` | page header color |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn` or `error` |
| `LOG_FORMAT` | `json` | `json` or `text`; logs go to stdout, one line per request |

Messages live in process memory and are gone when the app stops.
