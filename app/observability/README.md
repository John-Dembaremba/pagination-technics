# Observability Stack for Pagination Application

## Overview

This observability stack monitors the **Pagination App**, its **PostgreSQL database**, and the underlying **Node/System metrics**. The solution provides real-time performance insights through Grafana dashboards powered by Prometheus metrics.

## Monitored Components

### 1. Pagination App Metrics
- **API Request Rates**: Track request volume for cursor vs offset pagination
- **Response Times**: Monitor latency by pagination type
- **Error Rates**: Failed API requests
- **Throughput**: Requests per second

### 2. PostgreSQL Metrics
- **Query Performance**:
  ```promql
  rate(pg_stat_database_tup_returned[$__rate_interval])  # Rows read
  rate(pg_stat_database_xact_commit[$__rate_interval])  # Transactions
  ```
- **Efficiency Metrics**:
  ```promql
  pg_stat_database_blks_hit / (pg_stat_database_blks_read + pg_stat_database_blks_hit)  # Cache ratio
  sum(pg_locks_count)  # Lock contention
  ```

### 3. Node/System Metrics
- **Resource Usage**:
  ```promql
  sum(irate(node_cpu_seconds_total[$__rate_interval])) / scalar(count(...))  # CPU
  node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes  # Memory
  ```
- **Disk I/O**:
  ```promql
  rate(node_disk_read_bytes_total[$__rate_interval])  # Disk reads
  ```

## Key Dashboards

1. **Pagination Performance**
   - Compare cursor vs offset pagination efficiency
   - Monitor API response times

2. **Database Health**
   - Query throughput
   - Lock contention
   - Cache hit ratios

3. **System Resources**
   - CPU/Memory/Disk usage
   - Network throughput

## Getting Started

### 1. Environment Setup

First, prepare your configuration files:

```bash
# Copy the example files
cp env-example .env
cp example_datasources.json observability/datasources/prometheus/datasources.yaml
cp example_prometheus observability/prometheus/prometheus.yml
```

### 2. Required Configuration

Edit these files with your specific values:

#### `.env` Configuration:
```bash
# Database Configuration (REQUIRED)
POSTGRES_PSW="your_secure_password_here"
POSTGRES_USER="your_db_username"
POSTGRES_DB="your_database_name"
POSTGRES_PORT="5432"
POSTGRES_HOST="pagination_app_db"

# Monitoring (REQUIRED)
GRAFANA_ADMIN_PASSWORD="your_grafana_admin_password"

# Optional (change if needed)
PROJECT_VERSION="v1"
SERVER_PORT="3025"
POSTGRES_VERSION="17"
```

#### `datasources.yaml` Configuration:
```yaml
apiVersion: 1
datasources:
  - name: prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090  # Update port if different
    isDefault: true
```

#### `prometheus.yml` Configuration:
```yaml
global:
  scrape_interval: 5s

scrape_configs:
  - job_name: "node_exporter"
    static_configs:
      - targets: ["node_exporter:9100"]

  - job_name: "cadvisor"
    static_configs:
      - targets: ["cadvisor:8080"]

  - job_name: "postgres_exporter"
    static_configs:
      - targets: ["postgres_exporter:9187"]

  - job_name: "pagination-app"
    static_configs:
      - targets: ["app:3025"]  # Matches SERVER_PORT from .env
```

## Launching the Stack

1. Start the services:
   ```bash
   docker-compose up -d
   ```

2. Verify all containers are running:
   ```bash
   docker-compose ps
   ```

## Accessing Services

- **Grafana**: http://localhost:3000
  - Username: `admin`
  - Password: Value from `GRAFANA_ADMIN_PASSWORD` in `.env`

- **Prometheus**: http://localhost:9090
- **Adminer** (DB GUI): http://localhost:8081

## Security Notes

1. These files are excluded from Git via `.gitignore`
2. For production environments:
   - Use proper secret management
   - Change all default passwords
   - Restrict network access to monitoring ports
   - Consider enabling TLS for all connections

## Troubleshooting

If services fail to start:
1. Check logs:
   ```bash
   docker-compose logs
   ```
2. Verify all required fields in `.env` are populated
3. Ensure ports aren't already in use

## Future Roadmap

- **Tracing**: Add OpenTelemetry for request tracing
- **Logging**: Integrate Loki for log correlation
- **Synthetic Monitoring**: Add API test scenarios
