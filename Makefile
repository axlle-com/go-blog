NETWORK_NAME=projects_network

# Имя файлов Docker Compose
DOCKER_COMPOSE_PROJECT_BLOG=docker-compose.yml
SERVICES=postgres  redis
# elasticsearch app cli

# Цель по умолчанию
.PHONY: all
all: network up

# Цель для запуска всех проектов
.PHONY: up
up: network
	@echo "Starting Docker Compose projects..."
	@docker compose -f $(DOCKER_COMPOSE_PROJECT_BLOG) up -d $(SERVICES)

# Цель для пересборки Docker-образов и перезапуска контейнеров
.PHONY: rebuild
rebuild: network
	@echo "Stopping Docker Compose projects..."
	@docker compose -f $(DOCKER_COMPOSE_PROJECT_BLOG) down -v
	@echo "Rebuilding Docker Compose images..."
	@docker compose -f $(DOCKER_COMPOSE_PROJECT_BLOG) build --no-cache $(SERVICES)
	@echo "Restarting Docker Compose projects..."
	@docker compose -f $(DOCKER_COMPOSE_PROJECT_BLOG) up -d $(SERVICES)

# Цель для остановки и удаления контейнеров
.PHONY: down
down:
	@echo "Stopping Docker Compose projects..."
	@docker compose -f $(DOCKER_COMPOSE_PROJECT_BLOG) down -v

# Цель для создания сети, если она не существует
.PHONY: network
network:
	@if [ -z "$$(docker network ls --filter name=$(NETWORK_NAME) --format '{{ .Name }}')" ]; then \
		echo "Creating Docker network $(NETWORK_NAME)..."; \
		docker network create --driver bridge $(NETWORK_NAME); \
	else \
		echo "Docker network $(NETWORK_NAME) already exists."; \
	fi

# Цель для очистки сети
.PHONY: network-clean
network-clean:
	@if [ -n "$$(docker network ls --filter name=$(NETWORK_NAME) --format '{{ .Name }}')" ]; then \
		echo "Removing Docker network $(NETWORK_NAME)..."; \
		docker network rm $(NETWORK_NAME); \
	else \
		echo "Docker network $(NETWORK_NAME) does not exist."; \
	fi
