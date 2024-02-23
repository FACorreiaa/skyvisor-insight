#.PHONY: build clean
#project_name = skyvisor-container
#image_name = skyvisor-container
#
#compose-up:
#	make delete-container-if-exist
#	docker-compose up -d
#
#compose-down:
#	@docker compose down \
#  @docker volume rm postgres_db \
#  @docker compose up -d \
#  @rm-rf .data
#
#stop:
#	docker stop $(project_name)
#
#start:
#	docker start $(project_name)
#
#swag-init:
#	swag init --parseDependency
#
#go-test:
#	go test -v
#
#go-bench:
#	go test -bench .
#
#run-app:
#	docker compose run --rm app air init
#
#run-tidy:
#	docker compose run --rm app go mod tidy
#
#watch-tcss:
#	./tailwindcss -i controller/static/css/main.css -o controller/static/css/output.css --watch
#
#build-tcss:
#	./tailwindcss -i controller/static/css/main.css -o controller/static/css/output.css --minify
ifeq ("$(wildcard .env)","")
    $(shell cp env.sample .env)
	$(shell echo "DB_NAME=$($1)" > .env)
endif

include .env

$(eval export $(grep -v '^#' .env | xargs -0))

GO_MODULE := github.com/FACorreiaa/Aviation-tracker
VERSION  ?= $(shell git describe --tags --abbrev=0)
LDFLAGS   := -X "$(GO_MODULE)/config.Version=$(VERSION)"
#DB_DSN    := "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

tools: $(MIGRATE) $(AIR) $(MOCKERY) $(GOLANGCI) $(CHGLOG)
tools:
	@echo "Required tools are installed"

setup-local: tools
	@docker-compose up -d
	@sleep 5
	@./pgdev init

run:
	@air -c .air.toml --build.cmd "go build -ldflags \"$(LDFLAGS)\" -o ./tmp/main ."

stop:
	@docker compose down

lint:
	@golangci-lint run -c .golangci.yml

list-deps:
	go list -u -m all

upgrade-deps:
	go get -u ./...


