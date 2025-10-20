# Performance Testing with Locust

This directory contains Locust performance tests for the API.

## Prerequisites

```bash
pip install -r requirements.txt
```

## Running Tests

### Option 1: Using the script (with UI)
```bash
./run_tests.sh
```

Then open http://localhost:8089 in your browser and configure:
- Number of users
- Spawn rate
- Host (pre-configured to http://localhost:8080)

### Option 2: Headless mode
```bash
locust -f locustfile.py --host=http://localhost:8080 --headless -u 10 -r 2 -t 60s
```

Parameters:
- `-u 10`: 10 concurrent users
- `-r 2`: Spawn 2 users per second
- `-t 60s`: Run for 60 seconds

## Test Scenarios

The test suite includes weighted tasks simulating realistic API usage:

- **Create User** (weight: 3) - Creates new users with random data
- **List Users** (weight: 5) - Retrieves all users
- **Get User** (weight: 4) - Retrieves specific user by ID
- **Update User** (weight: 2) - Updates existing user data
- **Delete User** (weight: 1) - Deletes a user

## Observability

While tests are running, you can observe:
- **Traces** in Jaeger: http://localhost:16686
- **Metrics** in Prometheus: http://localhost:9090
- **Logs** with trace context: `docker logs api-sample`

## Example Prometheus Queries During Load Test

```promql
# Request rate
rate(http_server_requests_total[1m])

# P95 latency
histogram_quantile(0.95, rate(http_server_duration_milliseconds_bucket[1m]))

# Error rate
rate(http_server_errors_total[1m])
```
