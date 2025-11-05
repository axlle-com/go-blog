# Имя сети и файл docker-compose
NETWORK?=projects_network
#traefik
COMPOSE_FILE?=docker-compose.yml
SERVICES?=postgres redis
# elasticsearch app cli

all: env network up

up: env
	@docker compose -f $(COMPOSE_FILE) up -d $(SERVICES)

up-logging:
	@docker compose -f $(COMPOSE_FILE) --profile logging up -d $(SERVICES)

rebuild: env network
	@docker compose -f $(COMPOSE_FILE) down -v
	@docker compose -f $(COMPOSE_FILE) build --no-cache $(SERVICES)
	@docker compose -f $(COMPOSE_FILE) up -d $(SERVICES)

re-up: env network
	@docker compose -f $(COMPOSE_FILE) down -v
	@docker compose -f $(COMPOSE_FILE) up -d $(SERVICES)

down:
	@docker compose -f $(COMPOSE_FILE) down -v

up-dev: COMPOSE_FILE = docker-compose.dev.yml
up-dev: SERVICES = app postgres redis
up-dev: up

rebuild-dev: COMPOSE_FILE = docker-compose.dev.yml
rebuild-dev: SERVICES = app postgres redis
rebuild-dev: rebuild

re-up-dev: COMPOSE_FILE = docker-compose.dev.yml
re-up-dev: SERVICES = app postgres redis
re-up-dev: re-up

down-dev: COMPOSE_FILE = docker-compose.dev.yml
down-dev: SERVICES = app postgres redis
down-dev: down

network:
	@docker network inspect $(NETWORK) >/dev/null 2>&1 || { \
		echo "Создаю сеть $(NETWORK)..."; \
		docker network create --driver bridge $(NETWORK); \
	}

clean-network:
	@docker network inspect $(NETWORK) >/dev/null 2>&1 && { \
		echo "Удаляю сеть $(NETWORK)..."; \
		docker network rm $(NETWORK); \
	} || echo "Сеть $(NETWORK) не существует."

.PHONY: env
env:
	@if [ ! -f .env ]; then \
		if [ -f .env.example ]; then \
			echo "Создаю .env из .env.example..."; \
			cp .env.example .env; \
		else \
			echo "Файл .env отсутствует и .env.example не найден — создайте его вручную." >&2; \
			exit 1; \
		fi; \
	else \
		echo ".env уже существует"; \
	fi
