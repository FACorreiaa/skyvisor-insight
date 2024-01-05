FROM node:latest as assets
WORKDIR /app
COPY package.json ./
COPY package-lock.json ./
COPY postcss.config.cjs ./
COPY fonts.css ./
RUN mkdir -p controller/static/css controller/static/fonts
RUN npm install --ci
RUN npm run fonts

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
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
CMD ["air"]

# Define the final stage
FROM base as final
COPY --from=assets /app/controller/static/css/* ./controller/static/css/
COPY --from=assets /app/controller/static/fonts/* ./controller/static/fonts/
RUN CGO_ENABLED=0 go build -o /app/server
EXPOSE 6969
ENTRYPOINT ["/app/server"]
