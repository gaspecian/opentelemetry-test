#!/bin/bash

echo "Starting Locust performance tests..."
echo "API endpoint: http://localhost:8080"
echo ""
echo "Access Locust UI at: http://localhost:8089"
echo ""

locust -f locustfile.py --host=http://localhost:8080
