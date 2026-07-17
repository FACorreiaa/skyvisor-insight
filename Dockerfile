# syntax=docker/dockerfile:1.7

FROM oven/bun:1.3.14-alpine AS assets
WORKDIR /src
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile
COPY app ./app
RUN bun run build

FROM golang:1.26.5-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=assets /src/app/static ./app/static
RUN go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/skyvisor ./cmd/web
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/importer ./cmd/importer

FROM alpine:3.22 AS runtime
RUN apk add --no-cache ca-certificates tzdata && addgroup -S skyvisor && adduser -S -G skyvisor skyvisor
WORKDIR /app
COPY --from=build /out/skyvisor /usr/local/bin/skyvisor
COPY --from=build /out/migrate /usr/local/bin/skyvisor-migrate
COPY --from=build /out/importer /usr/local/bin/skyvisor-importer
USER skyvisor
EXPOSE 6969
ENTRYPOINT ["/usr/local/bin/skyvisor"]
