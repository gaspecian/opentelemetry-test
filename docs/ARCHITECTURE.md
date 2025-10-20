# Architecture Documentation

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Client Requests                          │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
                  ┌────────────────┐
                  │   Go CRUD API  │
                  │   (api-sample) │
                  └────────┬───────┘
                           │
                ┌──────────┼──────────┐
                │          │          │
                ▼          ▼          ▼
         ┌──────────┐  ┌──────┐  ┌────────┐
         │ MongoDB  │  │ OTLP │  │  OTLP  │
         │          │  │Traces│  │Metrics │
         └──────────┘  └───┬──┘  └───┬────┘
                           │         │
                           ▼         ▼
                    ┌──────────────────────┐
                    │  OTLP Collector      │
                    │  (monitoring-setup)  │
                    └──────────┬───────────┘
                               │
                    ┌──────────┼──────────┐
                    │          │          │
                    ▼          ▼          ▼
              ┌─────────┐ ┌──────────┐ ┌────────┐
              │ Jaeger  │ │Prometheus│ │  Loki  │
              │  (UI)   │ │   (UI)   │ │ (Logs) │
              └─────────┘ └─────┬────┘ └───▲────┘
                                │          │
                                ▼          │
                          ┌─────────┐  ┌──┴─────┐
                          │ Grafana │  │Promtail│
                          │   (UI)  │  │        │
                          └─────────┘  └────────┘
                                ▲
                                │
                          ┌─────┴─────┐
                          │ cAdvisor  │
                          │(Container │
                          │ Metrics)  │
                          └───────────┘
```

## Components

### Application Layer

**Go CRUD API**
- RESTful API with Create, Read, Update, Delete operations
- Instrumented with OpenTelemetry SDK
- Exports telemetry data via OTLP protocol
- Connects to MongoDB for data persistence
- Labeled with `monitor_category=application` for filtering

**MongoDB**
- NoSQL database for application data
- Runs in Docker container alongside the API

### Monitoring Layer

**OTLP Collector**
- Receives telemetry data from the API (traces, metrics, logs)
- Processes and routes data to appropriate backends
- Acts as a centralized collection point
- Ports: 4317 (gRPC), 4318 (HTTP), 8889 (metrics)

**Jaeger**
- Distributed tracing backend
- Stores and visualizes trace data
- Provides UI for trace analysis and service dependency mapping
- Port: 16686

**Prometheus**
- Time-series metrics database
- Scrapes metrics from OTLP Collector and cAdvisor
- Provides UI for metrics visualization and querying
- Port: 9090

**Grafana**
- Unified visualization platform
- Pre-configured dashboards for application monitoring
- Integrates with Prometheus (metrics) and Loki (logs)
- Category-based container filtering
- Port: 3000

**Loki**
- Log aggregation system
- Stores and indexes application logs
- Integrated with Grafana for log visualization
- Port: 3100

**Promtail**
- Log shipping agent
- Collects logs from Docker containers
- Ships logs to Loki
- Filters by `monitor_category` label

**cAdvisor**
- Container metrics collector
- Provides CPU, memory, network, and disk metrics
- Scraped by Prometheus
- Port: 8081

## Data Flow

### Traces
1. API generates spans via OpenTelemetry SDK
2. Spans sent to OTLP Collector (port 4317)
3. OTLP Collector forwards to Jaeger
4. Traces viewable in Jaeger UI and Grafana

### Metrics
1. API generates custom metrics (requests, errors, latency)
2. Metrics sent to OTLP Collector (port 4317)
3. OTLP Collector exposes metrics endpoint (port 8889)
4. Prometheus scrapes OTLP Collector and cAdvisor
5. Metrics viewable in Prometheus UI and Grafana dashboards

### Logs
1. API writes structured logs with TraceID/SpanID
2. Promtail collects logs from Docker containers
3. Promtail ships logs to Loki (port 3100)
4. Logs viewable in Grafana with TraceID correlation

### Container Metrics
1. cAdvisor collects container resource metrics
2. Prometheus scrapes cAdvisor (port 8081)
3. Metrics filtered by `monitor_category` label
4. Displayed in Grafana dashboards

## Network Communication

- **API ↔ MongoDB**: Internal Docker network (port 27017)
- **API → OTLP Collector**: OTLP gRPC (port 4317)
- **OTLP Collector → Jaeger**: Jaeger protocol
- **OTLP Collector → Prometheus**: Metrics scraping (port 8889)
- **Promtail → Loki**: HTTP push (port 3100)
- **Prometheus → cAdvisor**: Metrics scraping (port 8081)
- **Grafana → Prometheus**: PromQL queries (port 9090)
- **Grafana → Loki**: LogQL queries (port 3100)
- **User → UIs**: HTTP (Jaeger 16686, Prometheus 9090, Grafana 3000)

## Instrumentation Points

### Automatic Instrumentation
- HTTP server requests/responses (via otelmux)
- Database queries (via otelmongo)
- MongoDB operations
- Container resource usage (via cAdvisor)

### Custom Instrumentation
- Business logic spans (CreateUser, GetUser, etc.)
- Custom metrics (request counters, histograms, error counters)
- Structured logging with trace context

## Category-Based Filtering

Containers are labeled with `monitor_category` for organized monitoring:

- **application**: API services and microservices
- **database**: MongoDB and other databases
- **cache**: Redis, Memcached (future)
- **infrastructure**: Message queues, proxies (future)

Grafana dashboards use this label to filter and display only relevant containers.
- Custom metrics (counters, gauges, histograms)
- Structured logging with trace context

## Deployment

Both components run in isolated Docker Compose environments:
- **api-sample**: Application and database
- **monitoring-setup**: Observability stack

This separation allows independent scaling and management of application and monitoring infrastructure.
