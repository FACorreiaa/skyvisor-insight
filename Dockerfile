# Define the "base" stage
FROM golang:latest as base
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .

# Define the "dev" stage
FROM base as dev
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
WORKDIR /app
COPY /.air.toml ./
RUN go install github.com/a-h/templ/cmd/templ@latest
CMD ["air", "-c", ".air.toml"]

# Define the final stage
FROM base as final
RUN CGO_ENABLED=0 go build -o /app/server
EXPOSE 6969
ENTRYPOINT ["/app/server"]
