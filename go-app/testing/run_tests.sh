#!/bin/bash

echo "Installing Locust..."
pip install -r requirements.txt

echo "Starting performance test..."
echo "Make sure your Go API is running on http://localhost:5032"
echo ""
echo "Running Locust with:"
echo "- 10 users"
echo "- 2 users spawned per second"
echo "- 60 seconds duration"
echo ""

locust -f locustfile.py --host=http://localhost:5032 --users 10 --spawn-rate 2 --run-time 60s --headless
