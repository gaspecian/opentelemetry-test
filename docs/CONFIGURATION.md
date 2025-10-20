# OpenTelemetry Configuration Guide

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      Application Layer                       │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  API Service (Go)                                       │ │
│  │  - HTTP Server (Gorilla Mux + otelmux)                 │ │
│  │  - MongoDB Client (mongo-driver + otelmongo)           │ │
│  │  - Custom Spans & Metrics                              │ │
│  └────────────────┬───────────────────────────────────────┘ │
└───────────────────┼───────────────────────────────────────┘
                    │ OTLP/gRPC (port 4317)
                    ↓
┌─────────────────────────────────────────────────────────────┐
│                   Telemetry Collection Layer                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  OpenTelemetry Collector                               │ │
│  │  - OTLP Receiver (gRPC: 4317, HTTP: 4318)             │ │
│  │  - Batch Processor                                     │ │
│  │  - Jaeger Exporter (traces)                           │ │
│  │  - Prometheus Exporter (metrics)                      │ │
│  └────────────┬──────────────────┬────────────────────────┘ │
└───────────────┼──────────────────┼──────────────────────────┘
                │                  │
                ↓                  ↓
┌──────────────────────┐  ┌──────────────────────┐
│  Jaeger (Traces)     │  │  Prometheus (Metrics)│
│  Port: 16686         │  │  Port: 9090          │
└──────────────────────┘  └──────────────────────┘
```

## Component Configurations

### 1. OTLP Collector Configuration

**File:** `monitoring-setup/otel-collector-config.yaml`

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 10s
    send_batch_size: 1024

exporters:
  jaeger:
    endpoint: jaeger:14250
    tls:
      insecure: true
  
  prometheus:
    endpoint: 0.0.0.0:8889
  
  debug:
    verbosity: detailed

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [jaeger, debug]
    
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [prometheus]
```

**Key Settings:**
- **Batch processor:** Groups telemetry data for efficient export (10s timeout, 1024 batch size)
- **OTLP receivers:** Accept both gRPC and HTTP protocols
- **Exporters:** Jaeger for traces, Prometheus for metrics

### 2. Prometheus Configuration

**File:** `monitoring-setup/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'otel-collector'
    static_configs:
      - targets: ['otel-collector:8889']
```

**Key Settings:**
- **Scrape interval:** 15 seconds (balance between freshness and load)
- **Target:** OTLP Collector metrics endpoint

### 3. Application Configuration

**Environment Variables:**

| Variable | Value | Purpose |
|----------|-------|---------|
| `MONGODB_URI` | `mongodb://admin:password@mongodb:27017/apidb?authSource=admin` | Database connection |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `host.docker.internal:4317` | OTLP Collector endpoint |
| `SERVICE_NAME` | `api-sample` | Service identifier in traces |
| `SERVICE_VERSION` | `1.0.0` | Service version in traces |
| `PORT` | `8080` | API server port |

**OpenTelemetry SDK Initialization:**

```go
// Resource attributes
resource.New(ctx,
    resource.WithAttributes(
        semconv.ServiceNameKey.String(cfg.ServiceName),
        semconv.ServiceVersionKey.String(cfg.ServiceVersion),
    ),
)

// Trace provider with batch exporter
trace.NewTracerProvider(
    trace.WithBatcher(traceExporter),
    trace.WithResource(res),
)

// Metric provider with periodic reader (10s interval)
metric.NewMeterProvider(
    metric.WithReader(metric.NewPeriodicReader(metricExporter, 
        metric.WithInterval(10*time.Second))),
    metric.WithResource(res),
)
```

## Instrumentation Details

### Automatic Instrumentation

**HTTP Server (otelmux):**
```go
router.Use(otelmux.Middleware(cfg.ServiceName))
```
- Creates spans for all HTTP requests
- Captures HTTP method, status code, route
- Propagates trace context via headers

**MongoDB (otelmongo):**
```go
opts.Monitor = otelmongo.NewMonitor()
```
- Creates spans for all database operations
- Captures operation type, collection, database name
- Links to parent HTTP span

### Custom Instrumentation

**Custom Spans:**
```go
tracer := otel.Tracer("api-sample")
ctx, span := tracer.Start(ctx, "CreateUser")
defer span.End()

span.SetAttributes(attribute.String("user.id", id.Hex()))
```

**Custom Metrics:**
- `http_server_requests_total` - Counter
- `http_server_duration_milliseconds` - Histogram
- `http_server_errors_total` - Counter

**Structured Logging:**
```go
span := trace.SpanFromContext(ctx)
spanCtx := span.SpanContext()
log.Printf("[TraceID: %s] [SpanID: %s] %s", 
    spanCtx.TraceID().String(), 
    spanCtx.SpanID().String(), 
    message)
```

## Network Configuration

### Docker Networks

**Monitoring Stack:** Bridge network for Jaeger, Prometheus, OTLP Collector

**Application Stack:** Bridge network for API and MongoDB

**Cross-stack Communication:** 
- API uses `host.docker.internal:4317` to reach OTLP Collector
- Requires `extra_hosts: - "host.docker.internal:host-gateway"` in docker-compose

### Port Mappings

| Service | Internal Port | External Port | Protocol |
|---------|--------------|---------------|----------|
| API | 8080 | 8080 | HTTP |
| MongoDB | 27017 | 27017 | TCP |
| OTLP Collector (gRPC) | 4317 | 4317 | gRPC |
| OTLP Collector (HTTP) | 4318 | 4318 | HTTP |
| OTLP Collector (Metrics) | 8889 | 8889 | HTTP |
| Jaeger UI | 16686 | 16686 | HTTP |
| Jaeger Collector | 14250 | - | gRPC |
| Prometheus | 9090 | 9090 | HTTP |

## Performance Tuning

### Batch Processing
- **Timeout:** 10s - Balance between latency and efficiency
- **Batch size:** 1024 - Optimal for network efficiency

### Metric Collection
- **Export interval:** 10s - Reduces overhead while maintaining visibility
- **Prometheus scrape:** 15s - Aligned with metric export

### Resource Limits
No explicit limits set in PoC. For production:
- Set memory limits on OTLP Collector
- Configure queue sizes for backpressure handling
- Enable sampling for high-volume traces

## Security Considerations

**Current PoC Configuration (Development Only):**
- MongoDB uses basic auth (admin/password)
- OTLP uses insecure connections (no TLS)
- Jaeger exporter uses insecure connection

**Production Recommendations:**
- Enable TLS for all connections
- Use secrets management for credentials
- Implement authentication for OTLP endpoints
- Enable RBAC in Prometheus and Jaeger
- Use network policies to restrict access
