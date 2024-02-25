##initial build
#FROM golang:latest AS build-stage
#WORKDIR /app
#COPY go.mod go.sum ./
#RUN go mod download
#COPY . /app
#RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint
#
##dev stage
#FROM build-stage as dev
#RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
#WORKDIR /app
#CMD ["air"]
##CMD ["air", "-c", ".air.toml"]
#
##debug stage
#FROM build-stage as debug
#WORKDIR /app
#RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
#COPY . .
#COPY go.mod go.sum ./
#RUN go mod download
#COPY --from=build-stage /entrypoint /entrypoint-debug
#COPY --from=build-stage /app/controller/html /controller/html
#COPY --from=build-stage /app/controller/static /controller/static
#EXPOSE 40000
#CMD ["dlv", "--listen=127.0.0.1:40000", "--headless=true", "--api-version=2", "exec", "--accept-multiclient",  "/entrypoint-debug"]
#
## release
#FROM gcr.io/distroless/static-debian11 AS release-stage
#WORKDIR /
#COPY --from=build-stage /entrypoint /entrypoint
#COPY --from=build-stage /app/controller/html /controller/html
#COPY --from=build-stage /app/controller/static /controller/static
#EXPOSE 8080
#USER nonroot:nonroot
#ENTRYPOINT ["/entrypoint"]

# Build.
FROM golang:1.20 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint

# Deploy.
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /
COPY --from=build-stage /entrypoint /entrypoint
COPY --from=build-stage /app/controllers/static /controllers/static
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
