# Database Connection Pooling Configuration

## Overview

Shadow Nova uses PostgreSQL connection pooling via `pgxpool` to efficiently manage database connections in production. Connection pooling reduces the overhead of creating new connections and ensures optimal database resource utilization.

## Configuration

### Environment Variables

Configure the connection pool using these environment variables in your `.env` file:

```bash
# Maximum concurrent connections
DB_MAX_CONNS=25

# Minimum idle connections in the pool
DB_MIN_CONNS=5

# Connection timeout in seconds
DB_CONNECT_TIMEOUT=5
```

### Default Values

If not specified, the following defaults are used:
- `DB_MAX_CONNS`: 25 connections
- `DB_MIN_CONNS`: 5 connections
- `DB_CONNECT_TIMEOUT`: 5 seconds
- `MaxConnLifetime`: 1 hour
- `MaxConnIdleTime`: 30 minutes
- `HealthCheckPeriod`: 1 minute

## Connection Pool Parameters Explained

### MaxConns (Maximum Connections)

The maximum number of concurrent connections the pool can maintain.

**Guidelines:**
- **Small apps (< 1000 concurrent users)**: 10-25 connections
- **Medium apps (1000-10000 users)**: 25-50 connections
- **Large apps (> 10000 users)**: 50-100 connections

**Formula:**
```
MaxConns = ((core_count * 2) + effective_spindle_count)
```

For most PostgreSQL servers: `MaxConns = (CPU cores * 2) + 4`

**Example:**
- 4 CPU cores: MaxConns = 12
- 8 CPU cores: MaxConns = 20
- 16 CPU cores: MaxConns = 36

### MinConns (Minimum Idle Connections)

The minimum number of idle connections to maintain in the pool.

**Guidelines:**
- Set to 20-25% of MaxConns
- Too low: Connection creation overhead during traffic spikes
- Too high: Wasted database resources

**Example:**
- MaxConns=25: MinConns=5-7
- MaxConns=50: MinConns=10-15

### MaxConnLifetime

Maximum duration a connection can be reused before it's closed and recreated.

**Default:** 1 hour

**Purpose:**
- Prevents stale connections
- Ensures connections are periodically refreshed
- Helps with database server restarts

### MaxConnIdleTime

Maximum time a connection can remain idle before being closed.

**Default:** 30 minutes

**Purpose:**
- Frees up unused connections
- Prevents holding resources during low-traffic periods
- Balances between connection reuse and resource conservation

### HealthCheckPeriod

Frequency of background health checks on idle connections.

**Default:** 1 minute

**Purpose:**
- Detects broken connections early
- Removes dead connections from the pool
- Maintains pool health

### ConnectTimeout

Maximum time to wait when establishing a new connection.

**Default:** 5 seconds

**Purpose:**
- Prevents indefinite hanging on database connection
- Fails fast when database is unreachable
- Allows application to handle connection failures gracefully

## Workload-Specific Configurations

### High-Traffic Web Application

```bash
DB_MAX_CONNS=50
DB_MIN_CONNS=15
DB_CONNECT_TIMEOUT=3
```

**Characteristics:**
- Many concurrent requests
- Short-lived transactions
- Predictable traffic patterns

### Background Job Processor

```bash
DB_MAX_CONNS=20
DB_MIN_CONNS=3
DB_CONNECT_TIMEOUT=10
```

**Characteristics:**
- Long-running transactions
- Lower concurrency
- CPU-intensive operations

### API Gateway / Microservice

```bash
DB_MAX_CONNS=30
DB_MIN_CONNS=5
DB_CONNECT_TIMEOUT=5
```

**Characteristics:**
- Moderate concurrency
- Mix of read and write operations
- Burst traffic patterns

### Development Environment

```bash
DB_MAX_CONNS=10
DB_MIN_CONNS=2
DB_CONNECT_TIMEOUT=5
```

**Characteristics:**
- Single developer
- Low traffic
- Resource conservation

## Monitoring

### Prometheus Metrics

The following Prometheus metrics are automatically collected every 15 seconds:

| Metric | Type | Description |
|--------|------|-------------|
| `db_connections_active` | Gauge | Number of active connections in use |
| `db_connections_idle` | Gauge | Number of idle connections in the pool |
| `db_connections_max` | Gauge | Maximum allowed connections |
| `db_connections_total` | Gauge | Total connections in the pool |
| `db_connections_constructing` | Gauge | Connections currently being created |
| `db_connections_acquire_total` | Counter | Total successful connection acquisitions |
| `db_connections_acquire_duration_seconds` | Histogram | Connection acquisition latency |
| `db_connections_canceled_acquire_total` | Counter | Canceled connection acquisitions |
| `db_connections_new_total` | Counter | Total new connections opened |

### Health Check Endpoint

The `/health` endpoint includes connection pool statistics:

```json
{
  "status": "up",
  "database": "connected",
  "active_conns": "5",
  "idle_conns": "10",
  "max_conns": "25"
}
```

### Grafana Dashboard Queries

**Active Connections Over Time:**
```promql
db_connections_active
```

**Connection Pool Utilization:**
```promql
(db_connections_active / db_connections_max) * 100
```

**Connection Acquisition Rate:**
```promql
rate(db_connections_acquire_total[5m])
```

**Average Acquisition Duration:**
```promql
rate(db_connections_acquire_duration_seconds_sum[5m]) / rate(db_connections_acquire_duration_seconds_count[5m])
```

## Troubleshooting

### Problem: "Connection Pool Exhausted"

**Symptoms:**
- Requests timeout waiting for connections
- `db_connections_active` == `db_connections_max`

**Solutions:**
1. Increase `DB_MAX_CONNS`
2. Optimize slow queries
3. Check for connection leaks
4. Review transaction timeout settings

**Diagnostic Queries:**
```promql
# Check if pool is saturated
db_connections_active == db_connections_max

# Check acquisition wait time
rate(db_connections_acquire_duration_seconds_sum[5m])
```

### Problem: "Too Many Connections" on PostgreSQL

**Symptoms:**
- PostgreSQL error: "FATAL: too many connections"
- Database rejects new connections

**Solutions:**
1. Reduce `DB_MAX_CONNS` across all app instances
2. Increase `max_connections` in PostgreSQL config
3. Use connection pooler like PgBouncer

**Calculation:**
```
PostgreSQL max_connections >= Sum of all app instances' DB_MAX_CONNS + 10 (superuser reserve)
```

**Example:**
- 3 app instances with DB_MAX_CONNS=25 each
- PostgreSQL `max_connections` should be at least: (3 × 25) + 10 = 85

### Problem: High Connection Creation Rate

**Symptoms:**
- Frequent spikes in `db_connections_new_total`
- Poor performance during traffic bursts

**Solutions:**
1. Increase `DB_MIN_CONNS` to maintain more idle connections
2. Increase `MaxConnIdleTime` to keep connections longer
3. Check if connections are being closed prematurely

**Diagnostic Queries:**
```promql
# Connection creation rate
rate(db_connections_new_total[5m])

# Should be low and stable in production
```

### Problem: Idle Connections Consuming Resources

**Symptoms:**
- High `db_connections_idle` during low traffic
- Database resources being wasted

**Solutions:**
1. Decrease `DB_MIN_CONNS`
2. Decrease `MaxConnIdleTime`
3. Review traffic patterns

### Problem: Connection Timeout Errors

**Symptoms:**
- "context deadline exceeded" errors
- Unable to acquire connections

**Solutions:**
1. Increase `DB_CONNECT_TIMEOUT`
2. Check database network latency
3. Verify database is healthy and accepting connections
4. Check firewall rules

## Best Practices

### 1. Start Conservative, Scale Up

Begin with lower values and increase based on actual metrics:
```bash
DB_MAX_CONNS=15
DB_MIN_CONNS=3
```

Monitor and adjust based on:
- Connection pool saturation
- Request latency
- Database CPU usage

### 2. Match Database Capacity

Ensure your database server can handle all app instances:
```
Total MaxConns across all instances < PostgreSQL max_connections
```

### 3. Use Connection Pooling Layers

For high-scale deployments, consider:
- **PgBouncer**: Transaction-level pooling
- **pgcat**: Modern pooler with routing capabilities
- **Amazon RDS Proxy**: Managed solution for AWS

### 4. Monitor Continuously

Set up alerts for:
- Connection pool > 80% utilized
- High connection acquisition latency (> 100ms)
- Frequent connection errors

### 5. Test Under Load

Use load testing to validate configuration:
```bash
# Example using vegeta
echo "GET http://localhost:3000/api/health" | vegeta attack -duration=60s -rate=100 | vegeta report
```

### 6. Implement Graceful Degradation

Handle pool exhaustion gracefully:
- Return 503 Service Unavailable
- Implement request queuing
- Use circuit breakers

### 7. Review Regularly

Connection pool needs change over time:
- Application growth
- Feature additions
- Traffic pattern changes

Schedule quarterly reviews of your configuration.

## Security Considerations

### Connection String Security

Never hardcode connection strings:
```go
// BAD
databaseUrl := "postgres://user:password@localhost:5432/db"

// GOOD
databaseUrl := os.Getenv("DATABASE_URL")
```

### SSL/TLS Configuration

Always use SSL in production:
```bash
DATABASE_URL=postgres://user:password@host:5432/db?sslmode=require
```

SSL modes:
- `disable`: No SSL (dev only)
- `require`: SSL required but no verification
- `verify-ca`: Verify server certificate
- `verify-full`: Verify server certificate and hostname

### Credential Rotation

When rotating database credentials:
1. Update environment variable
2. Gracefully restart application
3. Existing connections will be replaced within `MaxConnLifetime`

## Performance Tuning

### Query Performance First

Before tuning the connection pool:
1. Optimize slow queries
2. Add appropriate indexes
3. Review N+1 query patterns
4. Use query caching where appropriate

### Transaction Timeout

Prevent long-running transactions:
```bash
DB_TX_TIMEOUT=30s
```

### Database Server Configuration

Tune PostgreSQL for your workload:
```sql
-- Recommended settings for connection pooling
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 16MB
maintenance_work_mem = 128MB
```

## Further Reading

- [pgxpool Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool)
- [PostgreSQL Connection Management](https://www.postgresql.org/docs/current/runtime-config-connection.html)
- [Grafana Dashboard Setup](https://grafana.com/docs/grafana/latest/dashboards/)
- [Prometheus Query Examples](https://prometheus.io/docs/prometheus/latest/querying/examples/)

## Support

For issues or questions about connection pooling:
1. Check application logs for connection errors
2. Review Prometheus metrics in Grafana
3. Consult this documentation
4. Open a GitHub issue with metrics and logs
