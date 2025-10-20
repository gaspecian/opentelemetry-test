# Observability Features

## ✅ Implemented Features

### 1. Distributed Tracing
- **Automatic HTTP instrumentation** via `otelmux` middleware
- **Automatic MongoDB instrumentation** via `otelmongo` monitor
- **Custom spans** for each business operation (CreateUser, GetUser, ListUsers, UpdateUser, DeleteUser)
- **Span attributes** including user IDs and operation counts
- **Trace propagation** across all components
- **Export to Jaeger** via OTLP Collector

### 2. Custom Metrics
All metrics are exported to Prometheus via OTLP Collector:

#### http_server_requests_total (Counter)
- Tracks total number of HTTP requests
- Labels: `http.method`, `http.status_code`

#### http_server_duration_milliseconds (Histogram)
- Tracks HTTP request duration in milliseconds
- Labels: `http.method`, `http.status_code`
- Useful for latency analysis and SLO monitoring

#### http_server_errors_total (Counter)
- Tracks total number of HTTP errors (status >= 400)
- Labels: `http.method`, `http.status_code`
- Useful for error rate monitoring

### 3. Structured Logging
- **Trace context injection** in all log messages
- Logs include `TraceID` and `SpanID` for correlation
- Format: `[TraceID: xxx] [SpanID: yyy] message`
- Enables log-trace correlation in observability platforms

### 4. Resource Attributes
- `service.name`: api-sample
- `service.version`: 1.0.0
- Automatically attached to all telemetry data

## Access Points

### Jaeger UI (Traces)
- URL: http://localhost:16686
- View distributed traces
- Analyze request flows
- Identify performance bottlenecks

### Prometheus (Metrics)
- URL: http://localhost:9090
- Query custom metrics
- Create dashboards
- Set up alerts

### OTLP Collector
- gRPC endpoint: localhost:4317
- HTTP endpoint: localhost:4318
- Metrics endpoint: localhost:8889/metrics

## Example Queries

### Prometheus Queries

```promql
# Request rate by status code
rate(http_server_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(http_server_duration_milliseconds_bucket[5m]))

# Error rate
rate(http_server_errors_total[5m])

# Success rate percentage
(sum(rate(http_server_requests_total{http_status_code=~"2.."}[5m])) / sum(rate(http_server_requests_total[5m]))) * 100
```

## Log Correlation Example

When you see a log like:
```
[TraceID: eea1877b39cabe1828d1864022759dab] [SpanID: d01ae60fbfa45499] Creating new user
```

You can:
1. Copy the TraceID
2. Search for it in Jaeger UI
3. See the complete distributed trace
4. Analyze all spans and timing information
