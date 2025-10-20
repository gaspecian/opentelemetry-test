# Usage Examples

## Quick Start

### 1. Start the System

```bash
# Start monitoring stack
cd monitoring-setup
docker compose up -d

# Start application
cd ../api-sample
docker compose up -d
```

### 2. Verify Services

```bash
# Check all services are running
docker ps

# Test API
curl http://localhost:8080/users

# Access UIs
# Jaeger: http://localhost:16686
# Prometheus: http://localhost:9090
```

## API Operations

### Create User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'

# Response:
# {"id":"68f653215b8b777e8a051636","name":"John Doe","email":"john@example.com"}
```

### Get User
```bash
curl http://localhost:8080/users/68f653215b8b777e8a051636

# Response:
# {"id":"68f653215b8b777e8a051636","name":"John Doe","email":"john@example.com"}
```

### List Users
```bash
curl http://localhost:8080/users

# Response:
# [{"id":"...","name":"...","email":"..."},...]
```

### Update User
```bash
curl -X PUT http://localhost:8080/users/68f653215b8b777e8a051636 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/users/68f653215b8b777e8a051636
```

## Observability Workflows

### Workflow 1: Trace a Specific Request

1. **Make a request and note the response time:**
```bash
time curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com"}'
```

2. **Check application logs for TraceID:**
```bash
docker logs api-sample | grep TraceID | tail -1
# Output: [TraceID: abc123...] [SpanID: def456...] Creating new user
```

3. **Search in Jaeger UI:**
- Open http://localhost:16686
- Select service: `api-sample`
- Paste TraceID in search box
- View complete trace with all spans

### Workflow 2: Monitor Request Rate

1. **Query current request rate in Prometheus:**
```bash
curl -s "http://localhost:9090/api/v1/query?query=rate(http_server_requests_total[1m])"
```

2. **Or use Prometheus UI:**
- Open http://localhost:9090
- Query: `rate(http_server_requests_total[1m])`
- Click "Execute"
- Switch to "Graph" tab

### Workflow 3: Analyze Latency

**P50 Latency:**
```promql
histogram_quantile(0.50, rate(http_server_duration_milliseconds_bucket[5m]))
```

**P95 Latency:**
```promql
histogram_quantile(0.95, rate(http_server_duration_milliseconds_bucket[5m]))
```

**P99 Latency:**
```promql
histogram_quantile(0.99, rate(http_server_duration_milliseconds_bucket[5m]))
```

### Workflow 4: Monitor Error Rate

**Total errors:**
```promql
sum(rate(http_server_errors_total[5m]))
```

**Error rate by status code:**
```promql
sum by (http_status_code) (rate(http_server_errors_total[5m]))
```

**Success rate percentage:**
```promql
(sum(rate(http_server_requests_total{http_status_code=~"2.."}[5m])) / 
 sum(rate(http_server_requests_total[5m]))) * 100
```

## Performance Testing

### Run Load Test with UI

```bash
cd api-sample/testing
./run_tests.sh
```

Then open http://localhost:8089 and configure:
- Number of users: 50
- Spawn rate: 5
- Host: http://localhost:8080 (pre-filled)

### Run Headless Load Test

**Light load (10 users, 30 seconds):**
```bash
cd api-sample/testing
locust -f locustfile.py --host=http://localhost:8080 \
  --headless -u 10 -r 2 -t 30s
```

**Medium load (50 users, 2 minutes):**
```bash
locust -f locustfile.py --host=http://localhost:8080 \
  --headless -u 50 -r 5 -t 2m
```

**Heavy load (100 users, 5 minutes):**
```bash
locust -f locustfile.py --host=http://localhost:8080 \
  --headless -u 100 -r 10 -t 5m
```

### Observe During Load Test

**Watch metrics in real-time:**
```bash
watch -n 1 'curl -s "http://localhost:9090/api/v1/query?query=rate(http_server_requests_total[1m])" | jq'
```

**Watch traces being created:**
- Open Jaeger UI: http://localhost:16686
- Set lookback to "Last 5 minutes"
- Click "Find Traces"
- Refresh periodically

**Watch application logs:**
```bash
docker logs -f api-sample
```

## Troubleshooting

### No traces appearing in Jaeger

1. **Check OTLP Collector logs:**
```bash
docker logs otel-collector | grep -i trace
```

2. **Verify API can reach collector:**
```bash
docker exec api-sample ping -c 2 host.docker.internal
```

3. **Check collector is receiving data:**
```bash
curl -s http://localhost:8889/metrics | grep otelcol_receiver_accepted_spans
```

### No metrics in Prometheus

1. **Check Prometheus targets:**
```bash
curl -s http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | {job, health, lastError}'
```

2. **Check collector metrics endpoint:**
```bash
curl -s http://localhost:8889/metrics | grep http_server_requests
```

### API not responding

1. **Check API logs:**
```bash
docker logs api-sample
```

2. **Check MongoDB connection:**
```bash
docker exec api-sample nc -zv mongodb 27017
```

3. **Restart services:**
```bash
cd api-sample
docker compose restart
```

## Cleanup

### Stop all services
```bash
cd monitoring-setup && docker compose down
cd ../api-sample && docker compose down
```

### Remove all data
```bash
cd monitoring-setup && docker compose down -v
cd ../api-sample && docker compose down -v
```

### Remove images
```bash
docker rmi api-sample-api
docker rmi $(docker images -f "dangling=true" -q)
```
