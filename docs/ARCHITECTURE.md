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
              │ Jaeger  │ │Prometheus│ │  Logs  │
              │  (UI)   │ │   (UI)   │ │        │
              └─────────┘ └──────────┘ └────────┘
```

## Components

### Application Layer

**Go CRUD API**
- RESTful API with Create, Read, Update, Delete operations
- Instrumented with OpenTelemetry SDK
- Exports telemetry data via OTLP protocol
- Connects to MongoDB for data persistence

**MongoDB**
- NoSQL database for application data
- Runs in Docker container alongside the API

### Monitoring Layer

**OTLP Collector**
- Receives telemetry data from the API (traces, metrics, logs)
- Processes and routes data to appropriate backends
- Acts as a centralized collection point

**Jaeger**
- Distributed tracing backend
- Stores and visualizes trace data
- Provides UI for trace analysis and service dependency mapping

**Prometheus**
- Time-series metrics database
- Scrapes metrics from OTLP Collector
- Provides UI for metrics visualization and querying

## Data Flow

1. **Application Execution**: API receives HTTP requests and processes CRUD operations
2. **Telemetry Generation**: OpenTelemetry SDK automatically captures:
   - Traces (request spans, database operations)
   - Metrics (request count, latency, error rates)
   - Logs (application events)
3. **Data Export**: Telemetry data sent to OTLP Collector via gRPC/HTTP
4. **Data Processing**: OTLP Collector receives, processes, and routes data
5. **Storage & Visualization**:
   - Traces → Jaeger
   - Metrics → Prometheus
   - Logs → OTLP Collector (can be extended to other backends)

## Network Communication

- **API ↔ MongoDB**: Internal Docker network
- **API → OTLP Collector**: OTLP protocol (gRPC on port 4317 or HTTP on port 4318)
- **OTLP Collector → Jaeger**: Jaeger protocol
- **OTLP Collector → Prometheus**: Prometheus remote write or scraping
- **User → UIs**: HTTP (Jaeger UI, Prometheus UI)

## Instrumentation Points

### Automatic Instrumentation
- HTTP server requests/responses
- Database queries (MongoDB operations)
- Outbound HTTP calls

### Custom Instrumentation
- Business logic spans
- Custom metrics (counters, gauges, histograms)
- Structured logging with trace context

## Deployment

Both components run in isolated Docker Compose environments:
- **api-sample**: Application and database
- **monitoring-setup**: Observability stack

This separation allows independent scaling and management of application and monitoring infrastructure.
