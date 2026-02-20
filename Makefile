ifneq (,$(wildcard .env))
	include .env
	export
endif

DB_URL ?= $(DATABASE_URL)
RESET ?= 0
LOGS ?= 0

run:
	@pkill -f "cmd/api" >/dev/null 2>&1 || true
	@pkill -f "cmd/worker" >/dev/null 2>&1 || true
	@$(MAKE) stop >/dev/null 2>&1 || true
ifeq ($(RESET),1)
	@echo "🧨 Resetting database..."
	@psql "$(DB_URL)" -f internal/storage/scripts/schema.sql
	@psql "$(DB_URL)" -f internal/storage/scripts/seed.sql
	@echo "✅ Database reset complete."
endif
	@echo "🚀 Starting API..."
	@go run ./cmd/api > api.log 2>&1 & echo $$! > api.pid
	@echo "🚀 Starting Worker..."
	@go run ./cmd/worker > worker.log 2>&1 & echo $$! > worker.pid
	@echo "✅ Both services started."
ifeq ($(LOGS),1)
	@echo "📄 Tailing logs..."
	@tail -f api.log worker.log
endif

stop:
	@echo "🛑 Stopping services..."
	@if [ -f api.pid ]; then kill `cat api.pid` && rm api.pid; fi
	@if [ -f worker.pid ]; then kill `cat worker.pid` && rm worker.pid; fi
	@echo "✅ Services stopped."

clean:
	@rm -f api.log worker.log api.pid worker.pid