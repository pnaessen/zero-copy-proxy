# Zero-Copy TCP Load Balancer

Small TCP load balancer written in Go.

The proxy listens on `:8123`, forwards traffic to several backends, and exposes Prometheus metrics on `:8081/metrics`.

## What it does

- Forwards TCP connections to backends with round-robin selection.
- Copies traffic in both directions with `io.Copy`
- Runs a background health check loop to keep only reachable backends.
- Exposes basic metrics:
	- `proxy_requests_total`
	- `proxy_active_servers`

## Project layout

- `main.go`: starts the TCP listener and metrics endpoint.
- `proxy/tcp.go`: handles bidirectional proxying between client and backend.
- `proxy/balancer.go`: round-robin logic and backend health checks.
- `docker-compose.yml`: local test setup (3 Nginx backends + Prometheus + Grafana).

## Run locally

1. Start the local stack:

```bash
docker compose up -d
```

2. Run the proxy:

```bash
go run main.go
```

3. Send traffic through the proxy:

```bash
curl http://localhost:8123
```

## Metrics and dashboards

- App metrics: http://localhost:8081/metrics
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (default login: admin/admin)

## Notes

- This project is Linux-oriented because the zero-copy behavior depends on kernel capabilities.
- If all backends are down, new client connections are rejected until a backend becomes healthy again.