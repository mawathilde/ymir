ENV_FILE_DEV=.env.dev
ENV_FILE_PROD=.env.prod
COMPOSE_BASE=docker-compose.yml

up-dev:
	docker compose --env-file $(ENV_FILE_DEV) -f $(COMPOSE_BASE) -f docker-compose.dev.yml up --build

down-dev:
	docker compose --env-file $(ENV_FILE_DEV) -f $(COMPOSE_BASE) -f docker-compose.dev.yml down

up-prod:
	docker compose --env-file $(ENV_FILE_PROD) -f $(COMPOSE_BASE) -f docker-compose.prod.yml up --build -d

down-prod:
	docker compose --env-file $(ENV_FILE_PROD) -f $(COMPOSE_BASE) -f docker-compose.prod.yml down

logs:
	docker compose logs  -f

ps:
	docker compose ps