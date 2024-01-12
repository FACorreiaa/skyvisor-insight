.PHONY: build clean
project_name = aviation-client
image_name = aviation-client

templ:
	@templ generate

build:
	npm run fonts
	templ generate
	go build -o server-exe

clean:
	find controller/static/fonts -type f ! -name "ionicons*" -delete
	rm controller/static/css/fonts.css
	rm server-exe

run-local:
	@templ generate
	@go run main.go

requirements:
	make clean-packages
	go mod tidy

clean-packages:
	go clean -modcache

up:
	make up-silent
	make shell

build:
	docker build -t $(image_name) .

push:
	docker build -t $(image_name) .

build-no-cache:
	docker build --no-cache -t $(image_name) .

#up-silent:
#	make delete-container-if-exist
#	docker run -d -p 3000:3000 --name $(project_name) $(image_name) ./main
#
#up-silent-prefork:
#	make delete-container-if-exist
#	docker run -d -p 3000:3000 --name $(project_name) $(image_name) ./app -prod
#
delete-container-if-exist:
	docker stop $(project_name) || true && docker rm $(project_name) || true

#shell:
#	docker exec -it $(project_name) /bin/sh

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
