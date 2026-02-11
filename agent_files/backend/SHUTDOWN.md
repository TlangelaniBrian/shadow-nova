# Graceful Shutdown Documentation

## Overview

The Shadow Nova backend implements a graceful shutdown mechanism that ensures all running tasks complete properly and resources are cleaned up before the application exits.

## Components

### 1. Server Context Management
The `Server` struct maintains a context specifically for the background collector goroutine:
- `collectorCtx`: Context for the collector goroutine lifecycle
- `collectorCancel`: Cancel function to signal shutdown to the collector

### 2. Shutdown Sequence

When a shutdown signal is received (SIGINT or SIGTERM), the application performs the following steps in order:

1. **HTTP Server Shutdown** (30-second timeout)
   - Stops accepting new connections
   - Waits for active requests to complete
   - Times out after 30 seconds if requests don't complete

2. **Collector Goroutine Shutdown**
   - Signals the collector to stop via context cancellation
   - Collector checks context in its main loop and exits cleanly

3. **Database Connection Cleanup**
   - Closes the PostgreSQL connection pool
   - Releases all database resources

4. **Flags Service Cleanup**
   - Closes the Unleash feature flags client connection

## How It Works

### Signal Handling
The main function sets up signal handlers for graceful shutdown:
```go
shutdown := make(chan os.Signal, 1)
signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
```

### Collector Goroutine Context Awareness
The collector goroutine checks for context cancellation in two places:

1. **Initial startup wait**: Uses `select` to respond immediately if shutdown happens during startup
2. **Sleep intervals**: Uses `select` with `time.After` to respond to shutdown while waiting

This ensures the collector can exit quickly without waiting for the full interval.

### Timeout Behavior
The shutdown has a 30-second timeout:
- If graceful shutdown completes within 30 seconds, the server exits cleanly
- If timeout is exceeded, the HTTP server is force-closed
- Application shutdown always completes regardless of timeout

## Testing

### Manual Testing

1. **Start the server**:
   ```bash
   cd backend
   go run main.go
   ```

2. **Send shutdown signal**:
   Press `Ctrl+C` or send SIGTERM:
   ```bash
   kill -TERM <pid>
   ```

3. **Verify expected log output**:
   ```
   Received signal: interrupt. Starting graceful shutdown...
   Initiating graceful shutdown...
   Stopping collector goroutine...
   Collector goroutine stopped gracefully
   Closing database connections...
   Closing flags service...
   Graceful shutdown complete
   Server stopped
   ```

### Automated Testing

You can test with a script:
```bash
# Start server in background
go run main.go &
SERVER_PID=$!

# Wait for server to start
sleep 2

# Send shutdown signal
kill -TERM $SERVER_PID

# Wait for graceful shutdown
wait $SERVER_PID
```

## Expected Log Messages

### Normal Operation
- `Server running on :3000`
- `Running initial content collection...`
- `Next collection in 24h0m0s (Runs per day: 1)`

### During Shutdown
- `Received signal: interrupt. Starting graceful shutdown...`
- `Initiating graceful shutdown...`
- `Stopping collector goroutine...`
- `Collector goroutine stopped gracefully`
- `Closing database connections...`
- `Closing flags service...`
- `Graceful shutdown complete`
- `Server stopped`

### Error Conditions
- `HTTP server shutdown error: <error>` - If HTTP server fails to shut down gracefully
- `HTTP server force close error: <error>` - If force close also fails
- `Application shutdown error: <error>` - If background task cleanup fails

## Timeout Behavior

### 30-Second Timeout
If the HTTP server has long-running requests that don't complete within 30 seconds:
1. The shutdown context will time out
2. The server will be force-closed
3. Background tasks will still complete their shutdown
4. You'll see: `HTTP server shutdown error: context deadline exceeded`

### Collector Shutdown
The collector goroutine responds immediately to shutdown signals:
- If collecting content when shutdown occurs, it will complete the current operation
- The context is checked between operations, not mid-operation
- Maximum delay is the time to complete one collection cycle

## Best Practices

1. **Development**: Always use `Ctrl+C` to stop the server for clean shutdown
2. **Production**: Use process managers (systemd, supervisor) that send SIGTERM
3. **Docker**: Ensure proper signal forwarding with `STOPSIGNAL SIGTERM`
4. **Kubernetes**: Configure appropriate `terminationGracePeriodSeconds` (default 30s is good)

## Deployment Configuration

### Docker
```dockerfile
STOPSIGNAL SIGTERM
```

### Kubernetes
```yaml
spec:
  terminationGracePeriodSeconds: 30
```

### Systemd
```ini
[Service]
KillSignal=SIGTERM
TimeoutStopSec=30
```

## Troubleshooting

### Server Doesn't Stop
- Check if goroutines are blocked on I/O operations
- Verify context is being passed to all long-running operations
- Look for goroutines not checking context cancellation

### Database Errors on Shutdown
- Check if database operations are context-aware
- Verify connection pool is properly closed
- Look for leaked connections

### Collector Doesn't Stop
- Verify context is passed to CollectAll and ProcessUnprocessedItems
- Check if those functions respect context cancellation
- Ensure select statements include context.Done() case
