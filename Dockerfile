FROM node:latest as assets
WORKDIR /app
COPY package.json ./
COPY package-lock.json ./
COPY postcss.config.cjs ./
COPY fonts.css ./
RUN mkdir -p controller/static/css controller/static/fonts
RUN npm install --ci
RUN npm run fonts

FROM golang:latest
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
COPY --from=assets /app/controller/static/css/* ./controller/static/css/
COPY --from=assets /app/controller/static/fonts/* ./controller/static/fonts/
RUN CGO_ENABLED=0 go build -o /app/server
EXPOSE 6969
ENTRYPOINT ["/app/server"]
