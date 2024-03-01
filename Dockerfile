# syntax=docker/dockerfile:1

# Build.
FROM golang:1.22 AS build-stage

LABEL Author="FC a11199"


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /entrypoint ./*.go

# dev stage
FROM build-stage AS app-dev
#we test
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64
RUN chmod +x tailwindcss-linux-arm64
RUN mv tailwindcss-linux-arm64 tailwindcss

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
RUN GOOS=linux go install github.com/a-h/templ/cmd/templ@latest
WORKDIR /app
CMD ["air"]

#FROM build-stage AS app-debug
#RUN CGO_ENABLED=0 go get -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv
#WORKDIR /app
# Debug Stage
FROM build-stage AS app-debug
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /app
#COPY --from=build-stage /entrypoint /entrypoint

# Deploy.
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /
COPY --from=build-stage /entrypoint /entrypoint
COPY --from=build-stage /app/controller/static /app/static
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/entrypoint"]


#FROM golang:alpine as build
#ARG project
#ARG version
#RUN apk add git
#WORKDIR /build
#COPY . .
#RUN go build -ldflags="-X 'bt.BuildHash=$version'" -o server bt/cmd/$project
#
#FROM alpine:latest
#RUN apk --no-cache add tzdata ca-certificates
#COPY --from=build /build/server /server
#
#ENTRYPOINT ["/server"]
