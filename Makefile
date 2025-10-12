# ========== CONFIG ==========
PROJECT_NAME = open-illustrations
DOCKER_COMPOSE = docker compose

# ========== COMMANDS ==========

# Build dan jalankan semua container
up:
	@echo "🚀 Starting $(PROJECT_NAME) containers..."
	$(DOCKER_COMPOSE) up -d --build
	@echo "✅ All services are up and running!"

# Stop semua container
down:
	@echo "🛑 Stopping $(PROJECT_NAME) containers..."
	$(DOCKER_COMPOSE) down
	@echo "✅ All containers stopped."

# Lihat log aplikasi (realtime)
logs:
	@echo "📜 Showing logs from app container..."
	$(DOCKER_COMPOSE) logs -f app

# Rebuild image (tanpa cache)
rebuild:
	@echo "🔄 Rebuilding Docker images..."
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) up -d
	@echo "✅ Rebuild complete!"

# Hapus semua container + volume (reset total)
clean:
	@echo "🧹 Removing all containers and volumes..."
	$(DOCKER_COMPOSE) down -v --remove-orphans
	@echo "✅ Clean complete!"

# Jalankan command di dalam container app
shell:
	@echo "🐚 Opening shell in app container..."
	$(DOCKER_COMPOSE) exec app sh

# Jalankan Go test (kalau kamu punya unit test)
test:
	@echo "🧪 Running Go tests..."
	$(DOCKER_COMPOSE) exec app go test ./...

# Jalankan migration SQL (kalau kamu ada init.sql)
migrate:
	@echo "📦 Running database migration..."
	$(DOCKER_COMPOSE) exec mysql bash -c "mysql -u$$DB_USER -p$$DB_PASS $$DB_NAME < /docker-entrypoint-initdb.d/init.sql"
	@echo "✅ Migration complete!"

# Full reset (hapus semua container + build ulang)
reset:
	@echo "♻️  Resetting everything (containers, images, volumes)..."
	$(DOCKER_COMPOSE) down -v --rmi all --remove-orphans
	$(DOCKER_COMPOSE) up -d --build
	@echo "✅ Full reset complete!"

# Check container status
status:
	@echo "📦 Container status:"
	$(DOCKER_COMPOSE) ps

