# GoRelay

A reverse proxy written in Go with load balancing, health checks, rate limiting, and a live terminal dashboard.

## Features
- Round robin load balancing across multiple backends
- Active health checks — automatically removes dead backends
- Rate limiting per client IP (100 requests per 10 second window)
- Live terminal dashboard showing backend status, request counts, latency
- Config file driven — no recompilation needed to change backends
- Add backends at runtime from the dashboard

## How it works
Client sends request → GoRelay picks a healthy backend using round robin → forwards the request → copies response back to client. If a backend is down, it gets skipped. If a client sends too many requests, they get a 429.

## Getting Started
```bash
git clone https://github.com/Anshtyagi1729/GoRelay
cd GoRelay
go run .
```

## Configuration
Edit `config.json`:
```json
{
  "port": ":3000",
  "backends": ["localhost:8081", "localhost:8082"],
  "rate_limit": 100,
  "window_seconds": 10,
  "timeout_seconds": 5,
  "health_check_interval": 10
}
```

## Dashboard Controls
- `A` — add a backend
- `Q` — quit

## Built with
- Go stdlib — net/http, sync, atomic
- [bubbletea](https://github.com/charmbracelet/bubbletea) — terminal UI
- [bubbles](https://github.com/charmbracelet/bubbles) — text input

## What's next
- Response caching
- Weighted round robin
- HTTPS support
- Metrics export
