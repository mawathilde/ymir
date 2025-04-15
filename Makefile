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

# TESTS

.PHONY: test test-db-up test-db-wait test-db-down

TEST_CONTAINER=ymir-test-db
TEST_DB_PORT=5433
TEST_DB_USER=test
TEST_DB_PASS=test
TEST_DB_NAME=testdb
TEST_DB_URL=postgresql://test:test@localhost:5433/testdb

## Lance les tests unitaires avec une base temporaire Docker
test:
	@echo "🐘 Starting PostgreSQL test container..."
	@docker run --rm -d \
		--name $(TEST_CONTAINER) \
		-e POSTGRES_USER=$(TEST_DB_USER) \
		-e POSTGRES_PASSWORD=$(TEST_DB_PASS) \
		-e POSTGRES_DB=$(TEST_DB_NAME) \
		-p $(TEST_DB_PORT):5432 \
		postgres:12-alpine > /dev/null

	@echo "⏳ Waiting for PostgreSQL to be ready..."
	@until docker exec $(TEST_CONTAINER) pg_isready -U $(TEST_DB_USER) -d $(TEST_DB_NAME) > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "✅ PostgreSQL is ready!"

	@echo "🚀 Running tests in ./api..."
	@trap 'echo \"🧹 Cleaning up...\"; docker stop $(TEST_CONTAINER) > /dev/null' EXIT; \
	cd api && DB=$(TEST_DB_URL) go run gotest.tools/gotestsum@latest --format short-verbose -- ./test
