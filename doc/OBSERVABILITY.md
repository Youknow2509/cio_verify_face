# Observability Stack for CIO Verify Face

This documentation explains how to set up and use Jaeger (distributed tracing), Prometheus (metrics), and Grafana (visualization) for monitoring all backend services.

## Quick Start

### Starting the Observability Stack

```bash
# Start only the observability stack
cd server
docker-compose -f docker-compose-observability.yml up -d

# Or start services with full observability integration
docker-compose -f docker-compose-with-observability.yml up -d
```

### Accessing the UIs

| Service    | URL                        | Credentials              |
|------------|----------------------------|--------------------------|
| Grafana    | http://localhost:3000      | admin / admin            |
| Prometheus | http://localhost:9091      | -                        |
| Jaeger     | http://localhost:16686     | -                        |

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Requests                          │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                       NGINX Gateway                              │
└────────────────────────────┬────────────────────────────────────┘
                             │
     ┌───────────────────────┼───────────────────────┐
     │                       │                       │
     ▼                       ▼                       ▼
┌─────────────┐      ┌─────────────┐         ┌─────────────┐
│ service-auth│      │service-ai   │         │ service-xxx │
│ (Go)        │      │ (Python)    │         │ (Go/Node)   │
│             │      │             │         │             │
│ ┌─────────┐ │      │ ┌─────────┐ │         │ ┌─────────┐ │
│ │Metrics  │ │      │ │Metrics  │ │         │ │Metrics  │ │
│ │:9090    │ │      │ │:8000    │ │         │ │:9090    │ │
│ └─────────┘ │      │ └─────────┘ │         │ └─────────┘ │
└──────┬──────┘      └──────┬──────┘         └──────┬──────┘
       │                    │                       │
       └────────────────────┼───────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │      Prometheus         │
              │      (Scraping)         │
              │       :9091             │
              └────────────┬────────────┘
                           │
                           ▼
              ┌─────────────────────────┐
              │        Grafana          │
              │     (Dashboards)        │
              │        :3000            │
              └─────────────────────────┘

       ┌────────────────────────────────────────┐
       │                                        │
       ▼                                        ▼
┌─────────────┐                          ┌─────────────┐
│ Jaeger      │ ◄────────────────────────│ All Services│
│ (Tracing)   │      OTLP Traces         │ send traces │
│ :16686      │                          │             │
└─────────────┘                          └─────────────┘
```

## Metrics Collected

### HTTP Metrics
- `http_requests_total` - Total number of HTTP requests (labels: method, path, status)
- `http_request_duration_seconds` - Request duration histogram
- `http_requests_in_flight` - Current number of requests being processed
- `http_response_size_bytes` - Response size histogram
- `http_errors_total` - Total number of HTTP errors (4xx, 5xx)

### gRPC Metrics
- `grpc_requests_total` - Total number of gRPC requests (labels: method, status)
- `grpc_request_duration_seconds` - gRPC request duration histogram
- `grpc_errors_total` - Total number of gRPC errors

## Pre-configured Grafana Dashboards

### CIO Verify Face - Services Dashboard
Includes the following panels:
- HTTP Request Rate by Service
- HTTP Request Duration (p50, p95) by Service
- HTTP Error Rate by Service (%)
- Service Health Status
- gRPC Request Rate by Service
- gRPC Request Duration by Service

## Service Configuration

### Go Services Configuration (YAML)

Add the following to your service's config.yaml:

```yaml
observability:
    enabled: true
    metrics_path: '/metrics'
    metrics_port: 9090
    tracing_enabled: true
    otlp_endpoint: 'http://jaeger:4318/v1/traces'
```

### Python Service Configuration (ENV)

```env
TRACING_ENABLED=true
OTLP_ENDPOINT=http://jaeger:4318/v1/traces
```

### Node.js Service Configuration (ENV)

```env
TRACING_ENABLED=true
JAEGER_ENDPOINT=http://jaeger:4318/v1/traces
```

## Tracing with Jaeger

### Viewing Traces
1. Open Jaeger UI at http://localhost:16686
2. Select a service from the dropdown
3. Click "Find Traces"
4. Click on a trace to see the full distributed trace

### Understanding Trace Data
- Each span represents a unit of work
- Spans include timing, status, and custom attributes
- Distributed traces show the full request flow across services

## Prometheus Queries Examples

### Request Rate
```promql
sum(rate(http_requests_total[5m])) by (service)
```

### 95th Percentile Latency
```promql
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (service, le))
```

### Error Rate
```promql
(sum(rate(http_errors_total[5m])) by (service) / sum(rate(http_requests_total[5m])) by (service)) * 100
```

### Service Availability
```promql
up
```

## Alerting

Pre-configured alerts in `config/prometheus/alerts.yml`:

| Alert Name        | Condition                              | Severity  |
|-------------------|----------------------------------------|-----------|
| HighErrorRate     | Error rate > 5% for 5 minutes          | Critical  |
| HighLatency       | P95 latency > 2s for 5 minutes         | Warning   |
| ServiceDown       | Service unavailable for 1 minute       | Critical  |
| HighRequestRate   | Request rate > 1000/s for 5 minutes    | Warning   |
| GRPCHighErrorRate | gRPC error rate > 5% for 5 minutes     | Critical  |
| GRPCHighLatency   | gRPC P95 latency > 2s for 5 minutes    | Warning   |

## Troubleshooting

### No metrics appearing in Prometheus
1. Check if the service is running: `docker-compose ps`
2. Verify the metrics endpoint is accessible: `curl http://localhost:9090/metrics`
3. Check Prometheus targets: http://localhost:9091/targets

### No traces appearing in Jaeger
1. Verify OTLP endpoint is configured correctly
2. Check service logs for tracing errors
3. Ensure Jaeger is running: `docker-compose logs jaeger`

### Grafana dashboard shows no data
1. Verify Prometheus datasource is configured
2. Check the time range in Grafana
3. Ensure services are generating traffic

## Best Practices

1. **Use meaningful span names** - Include the operation name and relevant context
2. **Add custom attributes** - Include user IDs, request IDs, and other contextual data
3. **Sample traces in production** - Use probabilistic sampling to reduce overhead
4. **Set up alerts** - Configure alerts for critical metrics
5. **Use dashboards** - Create dashboards for different audiences (ops, dev, business)

## Files Structure

```
server/
├── config/
│   ├── prometheus/
│   │   ├── prometheus.yml       # Prometheus configuration
│   │   └── alerts.yml           # Alert rules
│   └── grafana/
│       ├── provisioning/
│       │   ├── datasources/
│       │   │   └── datasources.yml
│       │   └── dashboards/
│       │       └── dashboards.yml
│       └── dashboards/
│           └── services-dashboard.json
├── pkg/
│   └── observability/           # Shared observability package for Go services
│       ├── config.go
│       ├── metrics.go
│       ├── tracing.go
│       └── grpc.go
├── docker-compose-observability.yml       # Standalone observability stack
└── docker-compose-with-observability.yml  # Full stack with observability
```
