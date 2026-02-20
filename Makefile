run:
	@echo "🚀 Starting API..."
	@go run ./cmd/api > api.log 2>&1 & echo $$! > api.pid
	@echo "🚀 Starting Worker..."
	@go run ./cmd/worker > worker.log 2>&1 & echo $$! > worker.pid
	@echo "✅ Both services started."
	@echo "Use 'make stop' to stop them."

stop:
	@echo "🛑 Stopping services..."
	@if [ -f api.pid ]; then kill `cat api.pid` && rm api.pid; fi
	@if [ -f worker.pid ]; then kill `cat worker.pid` && rm worker.pid; fi
	@echo "✅ Services stopped."

logs:
	@echo "📄 Tailing logs (Ctrl+C to exit)..."
	@tail -f api.log worker.log

clean:
	@rm -f api.log worker.log api.pid worker.pid