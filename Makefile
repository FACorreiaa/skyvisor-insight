.PHONY: build clean
project_name = aviation-client
image_name = aviation-client

compose-up:
	make delete-container-if-exist
	docker-compose up -d

compose-down:
	@docker compose down \
  @docker volume rm postgres_db \
  @docker compose up -d \
  @rm-rf .data

stop:
	docker stop $(project_name)

start:
	docker start $(project_name)

swag-init:
	swag init --parseDependency

go-test:
	go test -v

go-bench:
	go test -bench .

run-app:
	docker compose run --rm app air init

run-tidy:
	docker compose run --rm app go mod tidy

watch-tcss:
	./tailwindcss -i controller/static/css/main.css -o controller/static/css/output.css --watch

build-tcss:
	./tailwindcss -i controller/static/css/main.css -o controller/static/css/output.css --minify
