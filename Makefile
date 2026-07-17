.PHONY: assets build dev generate infra-up infra-down test verify

assets:
	bun install --frozen-lockfile
	bun run build

generate:
	go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate

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
