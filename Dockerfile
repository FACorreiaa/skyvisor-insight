# syntax=docker/dockerfile:1

# Build.
FROM golang:1.22 AS build-stage

LABEL Author="FC a11199"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /entrypoint ./*.go

FROM build-stage AS app-dev
#we test
#RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64
#RUN chmod +x tailwindcss-linux-arm64
#RUN mv tailwindcss-linux-arm64 tailwindcss

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
RUN GOOS=linux go install github.com/a-h/templ/cmd/templ@latest
WORKDIR /app
CMD ["air"]

FROM build-stage AS app-debug
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /app

# Deploy.
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /
COPY --from=build-stage /entrypoint /entrypoint
COPY --from=build-stage /app/app/static /app/static

COPY ./config /app/config

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/entrypoint"]
