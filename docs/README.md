# OpenTelemetry vs ELK - Proof of Concept

## Overview

This project is a Proof of Concept (PoC) to compare OpenTelemetry implementation against the traditional ELK stack (Elasticsearch, Logstash, and Kibana) for observability and monitoring.

## Project Structure

```
opentelemetry/
├── api-sample/           # Sample application
│   ├── docker-compose.yml
│   └── [Go CRUD API code]
└── monitoring-setup/     # OpenTelemetry monitoring stack
    └── docker-compose.yml
```

## Components

### api-sample
- **Technology**: GoLang
- **Type**: CRUD Web API
- **Database**: MongoDB
- **Deployment**: Docker Compose with MongoDB and the application

### monitoring-setup
- **OTLP Collector**: OpenTelemetry Protocol collector for receiving telemetry data
- **Jaeger**: Distributed tracing backend and UI
- **Prometheus**: Metrics collection and storage
- **Deployment**: Docker Compose with all monitoring components

## Objective

Compare the observability capabilities, ease of implementation, and performance characteristics between:
- **OpenTelemetry**: Modern, vendor-neutral observability framework
- **ELK Stack**: Traditional logging and monitoring solution

## Getting Started

1. Start the monitoring stack:
   ```bash
   cd monitoring-setup
   docker-compose up -d
   ```

2. Start the sample API:
   ```bash
   cd api-sample
   docker-compose up -d
   ```

## Evaluation Criteria

- Implementation complexity
- Performance overhead
- Feature completeness (traces, metrics, logs)
- Visualization capabilities
- Resource consumption
