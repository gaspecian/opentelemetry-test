# Performance Testing with Locust

This folder contains Locust performance tests for the TODO API.

## Setup

1. Install Python dependencies:
```bash
pip install -r requirements.txt
```

## Running Tests

### Quick Test (Automated)
```bash
./run_tests.sh
```

### Manual Test with Web UI
```bash
locust -f locustfile.py --host=http://localhost:8080
```
Then open http://localhost:8089 in your browser.

### Headless Test (Custom Parameters)
```bash
locust -f locustfile.py --host=http://localhost:8080 --users 50 --spawn-rate 5 --run-time 120s --headless
```

## Test Scenarios

The test simulates realistic user behavior:
- **Get all todos** (30% of requests) - Most common operation
- **Create todo** (20% of requests) - Regular creation
- **Get specific todo** (20% of requests) - Individual lookups
- **Update todo** (10% of requests) - Occasional updates
- **Delete todo** (10% of requests) - Cleanup operations

Each user waits 1-3 seconds between requests to simulate real usage patterns.
