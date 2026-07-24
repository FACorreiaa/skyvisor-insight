# SkyVisor Insight Web

SkyVisor Insight is a server-driven flight and travel companion. The web application combines fast, progressively enhanced pages with flight discovery, aviation reference data, interactive maps, secure accounts, and a foundation for trips, alerts, AI assistance, and native mobile clients.

This repository contains the web experience. It renders accessible HTML on the server and adds focused interactivity without a client-side application framework.

## What is included

- Public flight lookup with HTMX fragments and graceful non-JavaScript behavior
- Flight, airline, airport, city, and country exploration
- Interactive MapLibre and OpenLayers maps with locally bundled browser assets
- Cookie sessions backed by Redis, CSRF protection, and authenticated account pages
- Responsive light and dark themes using current TemplUI components
- Explicit migration and Aviationstack import commands
- Production-oriented HTTP timeouts, graceful shutdown, structured logging, and a non-root container

Trips, flight watches, soft Pro paywall, and AI itinerary import are integrated against the companion API. The web app exposes `/trips`, `/watches`, track-page Watch CTAs, and Stripe/dev upgrade.

## Stack

- Go 1.26.5, Chi, pgx, PostgreSQL, and Redis
- Templ and TemplUI for type-safe server-rendered components
- HTMX 2, Alpine.js 3, Tailwind CSS 4, MapLibre GL, and OpenLayers
- Bun for reproducible frontend asset builds
- Docker for production packaging

Browser dependencies and fonts are compiled into `app/static`; production pages do not depend on a JavaScript CDN.

## Local development

Requirements:

- Go 1.26.5 or newer
- Bun 1.3 or newer
- Docker with Compose, or locally available PostgreSQL and Redis instances
- An Aviationstack key only when running the importer or live flight lookups

Create the local configuration:

```sh
cp .env.example .env
```

Set `SESSION_KEY` to at least 32 random characters and provide `DB_PASS`. The defaults expect PostgreSQL on `127.0.0.1:5435`, Redis on `127.0.0.1:6381`, and the web server on `127.0.0.1:6969`.

Start local dependencies and prepare the application:

```sh
docker compose up -d postgres redis
bun install --frozen-lockfile
bun run build
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate
go run ./cmd/web
# or: make dev
```

Open `http://127.0.0.1:6969`.

**Schema migrations** under `db/migrations/*.sql` run automatically on web startup (`AUTO_MIGRATE` defaults to on). You can still run them alone with `make migrate` / `go run ./cmd/migrate`. Set `AUTO_MIGRATE=false` to leave schema to the Helm Job only.

Bulk **data** import is still separate (not schema):

```sh
go run ./cmd/importer
```

## Verification

```sh
bun run check
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate
go test -race ./...
go vet ./...
go build ./cmd/...
```

The same checks and a container build run in GitHub Actions.

## Container

```sh
docker build -t skyvisor-insight-web .
docker run --rm -p 6969:6969 --env-file .env skyvisor-insight-web
```

Use `ADDR=0.0.0.0:6969` and `COOKIE_SECURE=true` behind the production ingress.

## Platform repositories

The modernization is split into independently deployable, repo-ready foundations beside this directory:

| Repository | Responsibility |
| --- | --- |
| `skyvisor-api` | Go API, trips, SSE events, Aviationstack adapter, and optional Google GenAI assistant |
| `skyvisor-mcp` | Go MCP server exposing flight, trip, and assistant tools over stdio |
| `skyvisor-infra` | Hetzner Cloud Terraform, k3s, Helm, Argo CD, and GitOps configuration |
| `skyvisor-ios` | Native SwiftUI feature package and typed API client |
| `skyvisor-android` | Native Kotlin and Jetpack Compose application foundation |

The web application remains a deliberate server-driven frontend. Shared product data should move through `skyvisor-api`; native applications and the MCP server should not connect directly to the web database.

## Security notes

- Never commit `.env`, API keys, database credentials, session keys, or bearer tokens.
- Use HTTPS for Aviationstack and all non-loopback API traffic.
- Store production secrets with SOPS or an external secret manager and inject them at deployment time.
- The companion API validates OIDC JWTs and scopes durable trips by the verified `sub` claim. The web application's current cookie login must complete an Authorization Code with PKCE integration before it calls those authenticated endpoints on a user's behalf.

## License

Proprietary. Copyright (c) 2026 Fernando Correia. All rights reserved.
See `LICENSE` for the commercial licensing terms.
