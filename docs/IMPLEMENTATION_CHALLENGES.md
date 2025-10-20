# Implementation Challenges & Solutions

## Challenges Encountered

### 1. Docker Network Communication

**Challenge:**  
API container couldn't reach OTLP Collector using `host.docker.internal` on Linux.

**Error:**
```
ping: bad address 'host.docker.internal'
```

**Root Cause:**  
`host.docker.internal` is a Docker Desktop feature (macOS/Windows) and doesn't work natively on Linux.

**Solution:**  
Added `extra_hosts` configuration in docker-compose.yml:
```yaml
extra_hosts:
  - "host.docker.internal:host-gateway"
```

**Lesson Learned:**  
Always consider platform differences when using Docker networking features. Alternative solutions include:
- Using Docker network bridges
- Using host network mode
- Using container names for service discovery

---

### 2. OpenTelemetry Metric API Changes

**Challenge:**  
Metric recording failed with type mismatch errors.

**Error:**
```
cannot use attrs (variable of type []attribute.KeyValue) as []metric.AddOption
```

**Root Cause:**  
OpenTelemetry Go SDK changed the metric API. Attributes must be passed as `AttributeSet` with `WithAttributeSet` option.

**Solution:**
```go
// Before (incorrect)
h.metrics.RequestCounter.Add(ctx, 1, attrs...)

// After (correct)
attrs := attribute.NewSet(
    attribute.String("http.method", method),
    attribute.Int("http.status_code", statusCode),
)
h.metrics.RequestCounter.Add(ctx, 1, metric.WithAttributeSet(attrs))
```

**Lesson Learned:**  
OpenTelemetry APIs are still evolving. Always check the latest documentation and examples for the specific SDK version.

---

### 3. OTLP Endpoint Configuration

**Challenge:**  
Initial configuration used HTTP URL format for gRPC endpoint.

**Error:**
```
invalid target address http://host.docker.internal:4317, error info: address http://host.docker.internal:4317:443: too many colons in address
```

**Root Cause:**  
gRPC client expects host:port format, not HTTP URL format.

**Solution:**
```yaml
# Incorrect
OTEL_EXPORTER_OTLP_ENDPOINT: http://host.docker.internal:4317

# Correct
OTEL_EXPORTER_OTLP_ENDPOINT: host.docker.internal:4317
```

**Lesson Learned:**  
Different protocols have different endpoint format requirements. gRPC uses host:port, HTTP uses full URLs.

---

### 4. Trace Context Propagation

**Challenge:**  
Ensuring trace context flows through all layers (HTTP → Business Logic → Database).

**Solution:**  
- Used context.Context throughout the application
- Passed context from HTTP handlers to business logic
- MongoDB instrumentation automatically extracts context
- Custom spans use the same context

**Code Pattern:**
```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // Get context from HTTP request
    
    ctx, span := tracer.Start(ctx, "CreateUser")  // Create child span
    defer span.End()
    
    // Pass context to database operations
    result, err := h.db.DB.Collection("users").InsertOne(ctx, user)
}
```

**Lesson Learned:**  
Context propagation is critical for distributed tracing. Always pass context through function calls.

---

### 5. Metric Export Timing

**Challenge:**  
Metrics didn't appear immediately in Prometheus after requests.

**Root Cause:**  
Multiple timing factors:
- Metric export interval: 10 seconds
- Prometheus scrape interval: 15 seconds
- Potential delay: up to 25 seconds

**Solution:**  
This is expected behavior. For testing:
- Wait 30 seconds after generating traffic
- Use shorter intervals in development (trade-off: more overhead)

**Lesson Learned:**  
Understand the telemetry pipeline timing. Metrics are not real-time but near-real-time.

---

## Best Practices Discovered

### 1. Structured Logging with Trace Context

**Implementation:**
```go
func LogWithTrace(ctx context.Context, message string) {
    span := trace.SpanFromContext(ctx)
    spanCtx := span.SpanContext()
    
    if spanCtx.IsValid() {
        log.Printf("[TraceID: %s] [SpanID: %s] %s", 
            spanCtx.TraceID().String(), 
            spanCtx.SpanID().String(), 
            message)
    }
}
```

**Benefits:**
- Easy correlation between logs and traces
- Copy TraceID from logs to search in Jaeger
- Unified observability experience

---

### 2. Metric Naming Conventions

**Followed OpenTelemetry Semantic Conventions:**
- `http.server.requests` - Counter for requests
- `http.server.duration` - Histogram for latency
- `http.server.errors` - Counter for errors

**Benefits:**
- Consistent with industry standards
- Compatible with pre-built dashboards
- Easy to understand and query

---

### 3. Span Attributes

**Added meaningful attributes:**
```go
span.SetAttributes(
    attribute.String("user.id", id.Hex()),
    attribute.Int("user.count", len(users)),
)
```

**Benefits:**
- Rich context in traces
- Better filtering in Jaeger
- Easier debugging

---

### 4. Error Handling in Instrumentation

**Pattern:**
```go
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
    return err
}
```

**Benefits:**
- Errors visible in traces
- Automatic error tracking
- Better observability of failures

---

## Recommendations for Future Implementations

### 1. Start Simple
- Begin with automatic instrumentation
- Add custom spans only where needed
- Avoid over-instrumenting initially

### 2. Use Context Everywhere
- Make context.Context the first parameter
- Never break the context chain
- Use context for cancellation and timeouts

### 3. Test Observability Early
- Verify traces appear in Jaeger
- Check metrics in Prometheus
- Validate log correlation
- Don't wait until the end

### 4. Consider Production Early
- Plan for sampling strategies
- Think about data retention
- Consider security (TLS, auth)
- Plan for high availability

### 5. Document as You Go
- Note configuration decisions
- Document custom metrics
- Explain span hierarchies
- Keep troubleshooting notes

---

## Time Investment

**Total Implementation Time:** ~4 hours

**Breakdown:**
- Infrastructure setup: 30 minutes
- API development: 45 minutes
- OpenTelemetry instrumentation: 60 minutes
- Testing and validation: 45 minutes
- Documentation: 60 minutes

**Note:** Time includes learning curve and troubleshooting. Subsequent implementations would be faster.

---

## Complexity Assessment

**Overall Complexity:** Medium

**Easy Parts:**
- Automatic instrumentation (HTTP, MongoDB)
- Basic metric collection
- Docker deployment

**Moderate Parts:**
- Custom span creation
- Metric API usage
- Network configuration

**Challenging Parts:**
- Understanding trace propagation
- Debugging missing telemetry
- Optimizing for production

**Verdict:** OpenTelemetry has a learning curve but provides powerful capabilities once understood.
