# Shadow Nova - Observability with PLG Stack

This guide explains how to use the **PLG Stack** (Prometheus, Loki, Grafana) for monitoring and observability.

## Stack Overview

- **Prometheus**: Collects and stores metrics (HTTP requests, response times, CPU, memory)
- **Loki**: Aggregates logs from all services
- **Promtail**: Ships logs to Loki
- **Grafana**: Visualizes metrics and logs in unified dashboards

## Running the Stack

### Local Development (Docker Compose)

Start all services including the PLG stack:

```bash
docker-compose up -d
```

This will start:

- **Frontend**: http://localhost:8080
- **Backend**: http://localhost:3000
- **Grafana**: http://localhost:3001 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Loki**: http://localhost:3100

### Rancher on Windows

The `docker-compose.yml` is fully compatible with Rancher Desktop on Windows. Simply:

1. Open Rancher Desktop
2. Navigate to your project directory
3. Run: `docker-compose up -d`

All volumes use named volumes (not bind mounts) for Windows compatibility.

## Accessing Dashboards

### Grafana (Primary Dashboard)

1. Open http://localhost:3001
2. Login: `admin` / `admin`
3. Navigate to **Explore** to query:
   - **Prometheus**: Select "Prometheus" datasource for metrics
   - **Loki**: Select "Loki" datasource for logs

### Prometheus (Metrics)

Direct access: http://localhost:9090

Example queries:

- `http_requests_total` - Total HTTP requests
- `rate(http_requests_total[5m])` - Request rate over 5 minutes
- `http_request_duration_seconds` - Request latency

### Backend Metrics Endpoint

The Go backend exposes metrics at: http://localhost:3000/metrics

## Custom Metrics

The backend automatically tracks:

- `http_requests_total{method, path, status}` - Request counter
- `http_request_duration_seconds{method, path}` - Request duration histogram

To add custom metrics, use the Prometheus client in your Go code:

```go
import "github.com/prometheus/client_golang/prometheus/promauto"

var myCounter = promauto.NewCounter(prometheus.CounterOpts{
    Name: "my_custom_metric",
    Help: "Description of my metric",
})

myCounter.Inc()
```

## Logs

Logs are automatically collected by Promtail and sent to Loki. View them in Grafana:

1. Go to **Explore**
2. Select **Loki** datasource
3. Query: `{job="varlogs"}`

## Production Deployment

For production:

1. Set strong Grafana admin password via `GF_SECURITY_ADMIN_PASSWORD`
2. Configure persistent storage for Prometheus/Loki data
3. Set up retention policies
4. Enable authentication for Prometheus/Loki endpoints
5. Use a reverse proxy (Nginx/Traefik) for HTTPS

## Troubleshooting

### Prometheus not scraping backend

- Check backend is running: `curl http://localhost:3000/health`
- Check Prometheus targets: http://localhost:9090/targets
- Verify network connectivity: `docker network inspect shadow-nova-network`

### Grafana datasources not working

- Ensure Prometheus/Loki containers are running: `docker-compose ps`
- Check datasource URLs in Grafana settings
- Restart Grafana: `docker-compose restart grafana`
