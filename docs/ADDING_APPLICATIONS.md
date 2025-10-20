# Adding Applications to Monitoring

This guide explains how to add new applications to the Grafana monitoring dashboard.

## Overview

The monitoring setup uses a **category-based filtering system** that allows you to control which containers appear in the Grafana dashboard dropdown.

## Quick Start

To add a container to monitoring, add this label to your `docker-compose.yml`:

```yaml
labels:
  - "monitor_category=application"
```

## Step-by-Step Guide

### 1. Add Monitoring Label

Edit your application's `docker-compose.yml` and add the label:

```yaml
services:
  my-app:
    image: my-app:latest
    container_name: my-app
    labels:
      - "monitor_category=application"  # Add this line
    # ... rest of configuration
```

### 2. Configure OpenTelemetry (Optional)

If you want application-level metrics (requests, errors, latency), add OpenTelemetry instrumentation:

**Environment Variables:**
```yaml
environment:
  OTEL_EXPORTER_OTLP_ENDPOINT: host.docker.internal:4317
  SERVICE_NAME: my-app
  SERVICE_VERSION: 1.0.0
```

**Network Configuration:**
```yaml
extra_hosts:
  - "host.docker.internal:host-gateway"
```

### 3. Configure Log Shipping (Optional)

To send logs to Loki, update Promtail configuration in `api-sample/promtail-config.yaml`:

```yaml
scrape_configs:
  - job_name: docker
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
        filters:
          - name: label
            values: 
              - "com.docker.compose.service=api"
              - "com.docker.compose.service=my-app"  # Add your service
    relabel_configs:
      - source_labels: ['__meta_docker_container_label_com_docker_compose_service']
        target_label: 'container_name'
      - source_labels: ['__meta_docker_container_log_stream']
        target_label: 'stream'
```

### 4. Restart Services

```bash
# Restart your application
cd your-app-directory
docker compose up -d

# Restart Promtail (if logs configured)
cd api-sample
docker compose restart promtail

# Restart Prometheus to pick up new labels
cd monitoring-setup
docker compose restart prometheus
```

### 5. Verify in Grafana

1. Open Grafana: http://localhost:3000
2. Go to "Application Monitoring" dashboard
3. Check the dropdown - your application should appear
4. Select it to view metrics

## Available Categories

You can organize containers into different categories:

| Category | Purpose | Example Containers |
|----------|---------|-------------------|
| `application` | Application services | APIs, microservices, web apps |
| `database` | Database systems | MongoDB, PostgreSQL, MySQL |
| `cache` | Caching layers | Redis, Memcached |
| `infrastructure` | Infrastructure tools | Message queues, proxies |

**Example:**
```yaml
labels:
  - "monitor_category=database"  # For databases
```

## Complete Example

Here's a complete example for adding a new API service:

```yaml
version: '3.8'

services:
  my-new-api:
    build: .
    container_name: my-new-api
    ports:
      - "8081:8080"
    environment:
      # Application config
      DATABASE_URL: mongodb://mongodb:27017/mydb
      
      # OpenTelemetry config
      OTEL_EXPORTER_OTLP_ENDPOINT: host.docker.internal:4317
      SERVICE_NAME: my-new-api
      SERVICE_VERSION: 1.0.0
    extra_hosts:
      - "host.docker.internal:host-gateway"
    labels:
      - "monitor_category=application"
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
```

## Metrics Available

Once configured, you'll see these metrics in Grafana:

### Container Metrics (Automatic)
- ✅ CPU Usage
- ✅ Memory Usage

### Application Metrics (Requires OpenTelemetry)
- ✅ Transaction Rate (requests/sec)
- ✅ Error Rate (errors/sec)
- ✅ Request Latency
- ✅ HTTP Status Codes

### Logs (Requires Promtail Configuration)
- ✅ Application logs
- ✅ Error logs with TraceID correlation

## Troubleshooting

### Container not appearing in dropdown

**Check label:**
```bash
docker inspect my-app | grep monitor_category
```

**Restart Prometheus:**
```bash
cd monitoring-setup
docker compose restart prometheus
```

### No application metrics

**Verify OpenTelemetry connection:**
```bash
docker logs my-app | grep -i otel
```

**Check OTLP Collector:**
```bash
docker logs otel-collector
```

### No logs appearing

**Check Promtail configuration:**
```bash
docker logs promtail | grep my-app
```

**Verify Loki connection:**
```bash
curl http://localhost:3100/ready
```

## Best Practices

1. **Use descriptive service names** - They appear in the dropdown
2. **Add structured logging** - Include TraceID in logs for correlation
3. **Set appropriate categories** - Organize containers logically
4. **Monitor resource limits** - Set memory/CPU limits in docker-compose
5. **Use semantic versioning** - Track SERVICE_VERSION for deployments

## Next Steps

- [OpenTelemetry Instrumentation Guide](OBSERVABILITY_FEATURES.md)
- [Dashboard Customization](CONFIGURATION.md)
- [Performance Tuning](RESOURCE_CONSUMPTION.md)
