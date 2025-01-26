# Имя сети и файл docker-compose
NETWORK=projects_network
COMPOSE_FILE=docker-compose.yml
SERVICES=postgres redis
# elasticsearch app cli

all: network up

up:
	@docker compose -f $(COMPOSE_FILE) up -d $(SERVICES)

rebuild: network
	@docker compose -f $(COMPOSE_FILE) down -v
	@docker compose -f $(COMPOSE_FILE) build --no-cache $(SERVICES)
	@docker compose -f $(COMPOSE_FILE) up -d $(SERVICES)

down:
	@docker compose -f $(COMPOSE_FILE) down -v

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
