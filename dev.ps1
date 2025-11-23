# Shadow Nova - Development Script
# Runs both frontend and backend concurrently

Write-Host "Starting Shadow Nova Development Environment..." -ForegroundColor Green

# Start Infrastructure (Database, Unleash)
Write-Host "Starting Docker containers (DB, Unleash)..." -ForegroundColor Cyan
docker-compose up -d db unleash unleash-db


# Start backend in background
$backend = Start-Job -ScriptBlock {
    Set-Location "c:\Project.v2\view nova\shadow\shadow-nova\backend"
    go run main.go
} -Name "ShadowNova-Backend"

# Start frontend in background
$frontend = Start-Job -ScriptBlock {
    Set-Location "c:\Project.v2\view nova\shadow\shadow-nova\frontend"
    pnpm dev
} -Name "ShadowNova-Frontend"

Write-Host "`nBackend started (Job ID: $($backend.Id))" -ForegroundColor Cyan
Write-Host "Frontend started (Job ID: $($frontend.Id))" -ForegroundColor Cyan
Write-Host "`nPress Ctrl+C to stop all services`n" -ForegroundColor Yellow

# Display logs from both jobs
try {
    while ($true) {
        Receive-Job -Job $backend,$frontend
        Start-Sleep -Milliseconds 100
    }
}
finally {
    Write-Host "`nStopping services..." -ForegroundColor Yellow
    Stop-Job -Job $backend,$frontend
    Remove-Job -Job $backend,$frontend
    Write-Host "All services stopped." -ForegroundColor Green
}
