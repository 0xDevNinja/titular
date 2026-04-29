# gateway-go

Public API gateway service for the Titular protocol.

## Quickstart

```bash
go build ./...
go run .
```

## Build

```bash
docker build -t titular/gateway-go .
```

## Configuration

The binary is configured entirely via environment variables. Notable security-relevant
flags:

| Variable | Purpose |
|---|---|
| `GATEWAY_ADDR` | Listen address (default `:8080`). |
| `GATEWAY_SERVICE` | Service tag emitted in structured logs. |
| `GATEWAY_CORS_ORIGINS` | Comma-separated list of permitted origins. `*` enables wildcard. |
| `GATEWAY_CORS_CREDENTIALS` | `true` to advertise `Access-Control-Allow-Credentials`. |
| `GATEWAY_CORS_REFLECT_ORIGINS` | `true` to permit the wildcard+credentials combo (echoes the request `Origin`). Refused at startup otherwise. |
| `GATEWAY_RATE_LIMIT_RPS` | Per-key token-bucket refill rate; `0` disables limiting. |
| `GATEWAY_RATE_LIMIT_BURST` | Per-key burst capacity; `0` disables limiting. |
| `GATEWAY_TRUSTED_PROXIES` | Comma-separated CIDRs forwarded to `gin.SetTrustedProxies`. |

### `GATEWAY_TRUSTED_PROXIES` — required when behind an L7 proxy

By default (env unset) the gateway trusts **no proxies**: `c.ClientIP()` returns
`Request.RemoteAddr`, which is the unspoofable peer address of the connection.
This is the safe default — it prevents trivial rate-limit bypass via a forged
`X-Forwarded-For` header.

If you deploy the gateway behind a reverse proxy (NGINX, ALB, Cloudfront,
Cloudflare, etc.) you **must** set `GATEWAY_TRUSTED_PROXIES` to the list of
proxy CIDRs you trust. Only then will the gateway honour `X-Forwarded-For` /
`X-Real-IP` for client-IP attribution. Examples:

```bash
# Single ALB in 10.0.0.0/8
GATEWAY_TRUSTED_PROXIES=10.0.0.0/8

# Cloudflare ranges (trim to actuals)
GATEWAY_TRUSTED_PROXIES=173.245.48.0/20,103.21.244.0/22,...
```

### `GATEWAY_CORS_REFLECT_ORIGINS` — leave unset in production

The combination `GATEWAY_CORS_ORIGINS=*` + `GATEWAY_CORS_CREDENTIALS=true` is
refused at startup unless `GATEWAY_CORS_REFLECT_ORIGINS=true` is set. That
combination turns CORS into a credentialed reflector that browsers warn about;
production deployments should always list explicit origins instead.
