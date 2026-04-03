.PHONY: up down build logs ps clean test lint tidy

COMPOSE = docker compose -f docker/docker-compose.yml

# docker compose
up:
	$(COMPOSE) up

down:
	$(COMPOSE) down

build:
	$(COMPOSE) build

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

clean:
	$(COMPOSE) down --rmi all --volumes --remove-orphans

# go
test:
	$(COMPOSE) run --rm backend go test ./...

lint:
	$(COMPOSE) run --rm backend go vet ./...

tidy:
	$(COMPOSE) run --rm backend go mod tidy

back:
	$(COMPOSE) run --rm backend sh
