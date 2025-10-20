# OpenTelemetry Proof of Concept

A comprehensive implementation of OpenTelemetry for distributed tracing, metrics, and logging in a Go microservice with MongoDB.

## 🎯 Project Goals

- Evaluate OpenTelemetry as an observability solution
- Implement distributed tracing across HTTP and database layers
- Collect custom metrics for monitoring
- Enable log-trace correlation
- Compare with ELK Stack (future phase)

## 🏗️ Architecture

```
┌─────────────┐
│   API App   │ (Go + MongoDB)
│  Port 8080  │
└──────┬──────┘
       │ OTLP/gRPC
       ↓
┌─────────────────┐
│ OTLP Collector  │
│  Ports 4317/18  │
└────┬────────┬───┘
     │        │
     ↓        ↓
┌─────────┐ ┌────────────┐
│ Jaeger  │ │ Prometheus │
│  16686  │ │    9090    │
└─────────┘ └────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.24+ (for local development)
- Python 3.x (for load testing)

### Start the System

```bash
# 1. Start monitoring stack
cd monitoring-setup
docker compose up -d

# 2. Start application
cd ../api-sample
docker compose up -d

# 3. Verify services
curl http://localhost:8080/users
```

### Access UIs
- **Jaeger (Traces):** http://localhost:16686
- **Prometheus (Metrics):** http://localhost:9090
- **API:** http://localhost:8080

## 📊 Features

### ✅ Distributed Tracing
- Automatic HTTP request tracing
- Automatic MongoDB operation tracing
- Custom business logic spans
- Full trace propagation
- Span attributes for rich context

### ✅ Custom Metrics
- `http_server_requests_total` - Request counter
- `http_server_duration_milliseconds` - Latency histogram
- `http_server_errors_total` - Error counter

### ✅ Structured Logging
- Trace context in all logs
- TraceID and SpanID for correlation
- Easy log-to-trace navigation

### ✅ Performance Testing
- Locust-based load testing
- Realistic workload simulation
- UI and headless modes

## 📖 Documentation

| Document | Description |
|----------|-------------|
| [Implementation Plan](docs/IMPLEMENTATION_PLAN.md) | Step-by-step implementation guide |
| [Configuration](docs/CONFIGURATION.md) | Detailed configuration reference |
| [Usage Examples](docs/USAGE_EXAMPLES.md) | API usage and observability workflows |
| [Observability Features](docs/OBSERVABILITY_FEATURES.md) | Complete feature overview |
| [E2E Verification](docs/E2E_VERIFICATION.md) | Integration test results |
| [Resource Consumption](docs/RESOURCE_CONSUMPTION.md) | Performance metrics |
| [Implementation Challenges](docs/IMPLEMENTATION_CHALLENGES.md) | Lessons learned |
| [Comparison Criteria](docs/COMPARISON_CRITERIA.md) | ELK evaluation framework |

## 🧪 Testing

### Run Load Tests

```bash
cd api-sample/testing

# Install dependencies
pip install -r requirements.txt

# Run with UI
./run_tests.sh

# Run headless
locust -f locustfile.py --host=http://localhost:8080 \
  --headless -u 10 -r 2 -t 30s
```

### API Examples

```bash
# Create user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'

# List users
curl http://localhost:8080/users

# Get user
curl http://localhost:8080/users/{id}

# Update user
curl -X PUT http://localhost:8080/users/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane","email":"jane@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/users/{id}
```

## 📈 Observability Workflows

### Trace a Request

1. Make a request and check logs:
```bash
docker logs api-sample | grep TraceID | tail -1
```

2. Copy the TraceID and search in Jaeger UI

### Monitor Metrics

```bash
# Request rate
curl -s "http://localhost:9090/api/v1/query?query=rate(http_server_requests_total[1m])"

# P95 latency
curl -s "http://localhost:9090/api/v1/query?query=histogram_quantile(0.95,rate(http_server_duration_milliseconds_bucket[5m]))"

# Error rate
curl -s "http://localhost:9090/api/v1/query?query=rate(http_server_errors_total[1m])"
```

## 📦 Project Structure

```
.
├── api-sample/              # Go API application
│   ├── config/              # Configuration and OTel setup
│   ├── database/            # MongoDB connection
│   ├── handlers/            # HTTP handlers
│   ├── models/              # Data models
│   ├── testing/             # Locust load tests
│   ├── Dockerfile
│   └── docker-compose.yml
├── monitoring-setup/        # Observability stack
│   ├── otel-collector-config.yaml
│   ├── prometheus.yml
│   └── docker-compose.yml
└── docs/                    # Documentation
```

## 🔧 Technology Stack

**Application:**
- Go 1.24
- Gorilla Mux (HTTP router)
- MongoDB 7
- OpenTelemetry Go SDK

**Observability:**
- OpenTelemetry Collector
- Jaeger (distributed tracing)
- Prometheus (metrics)

**Testing:**
- Locust (load testing)

## 📊 Resource Usage

| Component | Memory | CPU (idle) |
|-----------|--------|------------|
| API | 8 MiB | 0.04% |
| MongoDB | 80 MiB | 0.19% |
| OTLP Collector | 33 MiB | 0.03% |
| Jaeger | 20 MiB | 0.01% |
| Prometheus | 23 MiB | 0.00% |
| **Total** | **164 MiB** | **0.27%** |

## 🎓 Key Learnings

1. **Context Propagation:** Critical for distributed tracing
2. **Automatic Instrumentation:** Reduces implementation effort
3. **Metric Timing:** Near-real-time (10-25s delay)
4. **Docker Networking:** Platform-specific considerations
5. **OpenTelemetry APIs:** Still evolving, check latest docs

## 🔜 Next Steps

- [ ] Implement ELK Stack equivalent
- [ ] Compare implementations
- [ ] Add Grafana dashboards
- [ ] Implement sampling strategies
- [ ] Add alerting rules
- [ ] Production hardening (TLS, auth)

## 🤝 Contributing

This is a PoC project for evaluation purposes. Feedback and suggestions are welcome!

## 📝 License

MIT License - See LICENSE file for details

## 🙏 Acknowledgments

- OpenTelemetry Community
- CNCF Projects (Jaeger, Prometheus)
- Go Community

---

**Status:** ✅ Phase 4 Complete - Ready for ELK Comparison

**Last Updated:** 2025-10-20
