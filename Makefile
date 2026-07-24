.PHONY: assets build dev generate api-client migrate infra-up infra-down test verify

OPENAPI_SPEC ?= ../skyvisor-api/api/openapi.yaml

assets:
	bun install --frozen-lockfile
	bun run build

generate:
	go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate

api-client:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.0 -config oapi-codegen.yaml "$(OPENAPI_SPEC)"

# Optional one-shot; the web process also migrates on start when AUTO_MIGRATE is on.
migrate:
	go run ./cmd/migrate

build: assets generate
	go build ./cmd/...

test:
	go test -race ./...

verify: assets generate test
	go vet ./...
	go build ./cmd/...

dev:
	air -c .air.toml

infra-up:
	docker compose up -d postgres redis

infra-down:
	docker compose down
