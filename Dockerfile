#initial build
FROM golang:1.20 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o -gcflags "all=-N -l" /entrypoint

#dev stage
FROM build-stage as dev
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
WORKDIR /app
CMD ["air", "-c", ".air.toml"]

#debug stage
FROM build-stage as debug
WORKDIR /app
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
COPY . .
COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -gcflags "all=-N -l" -o /entrypoint ./*.go
CMD ["dlv", "--listen=127.0.0.1:40000", "--headless=true", "--api-version=2", "exec", "--accept-multiclient",  "/aviation-tracker"]

# release
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /
COPY --from=build-stage /entrypoint /entrypoint
COPY --from=build-stage /app/controller/html /controller/html
COPY --from=build-stage /app/controller/static /controller/static
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/entrypoint"]
