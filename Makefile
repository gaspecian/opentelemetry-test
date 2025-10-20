.PHONY: help start stop restart status logs clean test build monitoring api all

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
BLUE   := \033[0;34m
NC     := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

##@ General

help: ## Display this help message
	@echo "$(BLUE)OpenTelemetry PoC - Available Commands$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make $(YELLOW)<target>$(NC)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(BLUE)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Lifecycle

all: start ## Start all services (monitoring + API)
	@echo "$(GREEN)✓ All services started$(NC)"
	@make status

start: monitoring api ## Start monitoring stack and API
	@echo "$(GREEN)✓ System ready$(NC)"
	@echo ""
	@echo "$(BLUE)Access Points:$(NC)"
	@echo "  API:        http://localhost:8080"
	@echo "  Grafana:    http://localhost:3000 (admin/admin)"
	@echo "  Jaeger:     http://localhost:16686"
	@echo "  Prometheus: http://localhost:9090"

stop: ## Stop all services
	@echo "$(YELLOW)⏸ Stopping all services...$(NC)"
	@cd api-sample && docker compose down
	@cd monitoring-setup && docker compose down
	@echo "$(GREEN)✓ All services stopped$(NC)"

restart: stop start ## Restart all services

clean: stop ## Stop services and remove volumes
	@echo "$(YELLOW)🧹 Cleaning up volumes...$(NC)"
	@cd api-sample && docker compose down -v
	@cd monitoring-setup && docker compose down -v
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

##@ Services

monitoring: ## Start monitoring stack (Grafana, Prometheus, Jaeger, Loki)
	@echo "$(BLUE)🚀 Starting monitoring stack...$(NC)"
	@cd monitoring-setup && docker compose up -d
	@echo "$(YELLOW)⏳ Waiting for services to be ready...$(NC)"
	@sleep 5
	@echo "$(GREEN)✓ Monitoring stack ready$(NC)"

api: ## Start API and database
	@echo "$(BLUE)🚀 Starting API services...$(NC)"
	@cd api-sample && docker compose up -d
	@echo "$(YELLOW)⏳ Waiting for API to be ready...$(NC)"
	@sleep 5
	@echo "$(GREEN)✓ API services ready$(NC)"

##@ Status & Logs

status: ## Show status of all services
	@echo "$(BLUE)📊 Service Status:$(NC)"
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "NAME|api-sample|mongodb|promtail|grafana|loki|cadvisor|jaeger|prometheus|otel-collector" || echo "$(RED)No services running$(NC)"

logs: ## Show logs from all services (use SERVICE=name for specific service)
ifdef SERVICE
	@echo "$(BLUE)📋 Logs for $(SERVICE):$(NC)"
	@docker logs -f $(SERVICE)
else
	@echo "$(YELLOW)Usage: make logs SERVICE=<service-name>$(NC)"
	@echo "Available services:"
	@docker ps --format "  - {{.Names}}" | grep -E "api-sample|mongodb|promtail|grafana|loki|cadvisor|jaeger|prometheus|otel-collector"
endif

logs-api: ## Show API logs
	@echo "$(BLUE)📋 API Logs:$(NC)"
	@docker logs -f api-sample

logs-monitoring: ## Show monitoring stack logs
	@echo "$(BLUE)📋 Monitoring Logs:$(NC)"
	@docker compose -f monitoring-setup/docker-compose.yml logs -f

##@ Testing

test: ## Run load tests with Locust
	@echo "$(BLUE)🧪 Running load tests...$(NC)"
	@cd api-sample/testing && ./run_tests.sh

test-headless: ## Run headless load tests (10 users, 30s)
	@echo "$(BLUE)🧪 Running headless load tests...$(NC)"
	@cd api-sample/testing && locust -f locustfile.py --host=http://localhost:8080 --headless -u 10 -r 2 -t 30s

test-api: ## Test API endpoints
	@echo "$(BLUE)🧪 Testing API endpoints...$(NC)"
	@echo "$(YELLOW)GET /users:$(NC)"
	@curl -s http://localhost:8080/users | jq '.' || echo "$(RED)Failed$(NC)"
	@echo ""
	@echo "$(YELLOW)POST /users:$(NC)"
	@curl -s -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"Test User","email":"test@example.com"}' | jq '.' || echo "$(RED)Failed$(NC)"

##@ Build

build: ## Build API Docker image
	@echo "$(BLUE)🔨 Building API image...$(NC)"
	@cd api-sample && docker compose build
	@echo "$(GREEN)✓ Build complete$(NC)"

rebuild: clean build start ## Clean, rebuild, and start everything

##@ Monitoring

dashboard: ## Open Grafana dashboard in browser
	@echo "$(BLUE)🌐 Opening Grafana...$(NC)"
	@xdg-open http://localhost:3000 2>/dev/null || open http://localhost:3000 2>/dev/null || echo "$(YELLOW)Open http://localhost:3000 manually$(NC)"

jaeger: ## Open Jaeger UI in browser
	@echo "$(BLUE)🌐 Opening Jaeger...$(NC)"
	@xdg-open http://localhost:16686 2>/dev/null || open http://localhost:16686 2>/dev/null || echo "$(YELLOW)Open http://localhost:16686 manually$(NC)"

prometheus: ## Open Prometheus UI in browser
	@echo "$(BLUE)🌐 Opening Prometheus...$(NC)"
	@xdg-open http://localhost:9090 2>/dev/null || open http://localhost:9090 2>/dev/null || echo "$(YELLOW)Open http://localhost:9090 manually$(NC)"

##@ Utilities

health: ## Check health of all services
	@echo "$(BLUE)🏥 Health Check:$(NC)"
	@echo -n "API:        " && curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/users && echo " $(GREEN)✓$(NC)" || echo " $(RED)✗$(NC)"
	@echo -n "Grafana:    " && curl -s -o /dev/null -w "%{http_code}" http://localhost:3000 && echo " $(GREEN)✓$(NC)" || echo " $(RED)✗$(NC)"
	@echo -n "Prometheus: " && curl -s -o /dev/null -w "%{http_code}" http://localhost:9090 && echo " $(GREEN)✓$(NC)" || echo " $(RED)✗$(NC)"
	@echo -n "Jaeger:     " && curl -s -o /dev/null -w "%{http_code}" http://localhost:16686 && echo " $(GREEN)✓$(NC)" || echo " $(RED)✗$(NC)"
	@echo -n "Loki:       " && curl -s -o /dev/null -w "%{http_code}" http://localhost:3100/ready && echo " $(GREEN)✓$(NC)" || echo " $(RED)✗$(NC)"

metrics: ## Show current metrics from Prometheus
	@echo "$(BLUE)📊 Current Metrics:$(NC)"
	@echo "$(YELLOW)Request Rate:$(NC)"
	@curl -s "http://localhost:9090/api/v1/query?query=rate(http_server_requests_total[1m])" | jq -r '.data.result[] | "  \(.metric.http_method) \(.metric.http_status_code): \(.value[1])"' 2>/dev/null || echo "  No data"
	@echo "$(YELLOW)Error Rate:$(NC)"
	@curl -s "http://localhost:9090/api/v1/query?query=rate(http_server_errors_total[1m])" | jq -r '.data.result[] | "  Status \(.metric.http_status_code): \(.value[1])"' 2>/dev/null || echo "  No errors"

traffic: ## Generate test traffic (50 requests)
	@echo "$(BLUE)🚦 Generating traffic...$(NC)"
	@for i in $$(seq 1 50); do \
		curl -s http://localhost:8080/users > /dev/null && echo -n "$(GREEN).$(NC)"; \
		sleep 0.1; \
	done
	@echo ""
	@echo "$(GREEN)✓ Generated 50 requests$(NC)"

##@ Documentation

docs: ## Open documentation
	@echo "$(BLUE)📚 Available Documentation:$(NC)"
	@echo "  - README.md"
	@echo "  - docs/ADDING_APPLICATIONS.md"
	@echo "  - docs/IMPLEMENTATION_PLAN.md"
	@echo "  - docs/CONFIGURATION.md"
	@echo "  - docs/USAGE_EXAMPLES.md"

info: ## Show project information
	@echo "$(BLUE)OpenTelemetry PoC$(NC)"
	@echo ""
	@echo "$(YELLOW)Components:$(NC)"
	@echo "  - Go API with MongoDB"
	@echo "  - OpenTelemetry Collector"
	@echo "  - Jaeger (Tracing)"
	@echo "  - Prometheus (Metrics)"
	@echo "  - Grafana (Dashboards)"
	@echo "  - Loki (Logs)"
	@echo ""
	@echo "$(YELLOW)Ports:$(NC)"
	@echo "  - 8080:  API"
	@echo "  - 3000:  Grafana"
	@echo "  - 9090:  Prometheus"
	@echo "  - 16686: Jaeger"
	@echo "  - 3100:  Loki"
	@echo "  - 4317:  OTLP gRPC"
	@echo "  - 4318:  OTLP HTTP"
