.PHONY: help up build py bot logs stop clean restart down cleanup-cache cleanup-volumes

# デフォルトターゲット
help:
	@echo "Available commands:"
	@echo "  make up      - Build and start all containers"
	@echo "  make build   - Build all Docker images"
	@echo "  make py      - Execute separate-words main.py"
	@echo "  make bot     - Start data-collect-bot with Dockerfile.bot"
	@echo "  make logs    - Show all container logs"
	@echo "  make stop    - Stop all containers"
	@echo "  make down    - Stop and remove all containers"
	@echo "  make clean   - Stop, remove containers and images"
	@echo "  make restart - Restart all containers"
	@echo "  make cleanup-cache - Clean up Docker build cache and unused resources"
	@echo "  make cleanup-volumes - Clean up Docker volumes (CAUTION: removes data)"

# 全コンテナをビルドして起動
up:
	docker-compose up --build -d

# 全イメージをビルド
build:
	docker-compose build

# separate-wordsのmain.pyを実行
py:
	@echo "Starting containers if not running..."
	@docker-compose up -d || (echo "Removing conflicting containers..." && docker-compose down && docker-compose up -d)
	@echo "Executing main.py..."
	docker-compose exec separate-words python app/main.py

# data-collect-botをDockerfile.botで起動
bot:
	@echo "Starting data-collect-bot with Dockerfile.bot..."
	@docker stop data-collect-bot 2>/dev/null || true
	@docker rm data-collect-bot 2>/dev/null || true
	docker build -f data-collect-bot/Dockerfile.bot -t data-collect-bot:latest data-collect-bot/
	docker run -d --name data-collect-bot --env-file data-collect-bot/.env \
		-e DB_HOST=data-collect-postgres \
		-e DB_PORT=5432 \
		-e DB_USER=admin \
		-e DB_PASSWORD=password123 \
		-e DB_NAME=datacollect \
		--network miteruyo-bot_miteruyo-network data-collect-bot:latest

database:
	docker exec -it data-collect-postgres psql -U admin -d datacollect
# 全コンテナのログを表示
logs:
	docker-compose logs

# 全コンテナを停止
stop:
	docker-compose stop

# 全コンテナを停止・削除
down:
	@docker stop data-collect-bot 2>/dev/null || true
	@docker rm data-collect-bot 2>/dev/null || true
	docker-compose down -v

# コンテナとイメージを削除
clean:
	docker-compose down
	docker rmi miteruyo-bot-data-collect-bot miteruyo-bot-separate-words 2>/dev/null || true

# 全コンテナを再起動
restart: down up

# Dockerキャッシュと不要なリソースを削除（安全）
cleanup-cache:
	@echo "Cleaning up Docker build cache and unused resources..."
	docker system prune -f
	docker image prune -f
	docker builder prune -f
	@echo "Docker cache cleanup completed!"

# Dockerボリュームを削除（注意：データが削除される）
cleanup-volumes:
	@echo "WARNING: This will remove all unused Docker volumes and their data!"
	@echo "Press Ctrl+C to cancel, or wait 5 seconds to continue..."
	@sleep 5
	docker volume prune -f
	@echo "Docker volumes cleanup completed!"