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
| `GATEWAY_JWT_SECRET` | Base64-encoded HMAC key (>=32 bytes after decode). Setting this enables the SIWE auth endpoints. |
| `GATEWAY_JWT_ISSUER` | `iss` claim minted into and required on every JWT. Defaults to `titular-gateway`. |
| `GATEWAY_JWT_TTL` | Session lifetime as a Go duration. Defaults to `24h`. |
| `GATEWAY_REDIS_URL` | Redis connection URL (`redis://host:port/db`). Required when `GATEWAY_JWT_SECRET` is set. |
| `GATEWAY_SIWE_DOMAIN` | The value the SIWE message MUST declare in its `domain` line. Pinned server-side to prevent cross-site replay. |
| `GATEWAY_SIWE_CHAIN_ID` | The chain id the SIWE message MUST declare. Defaults to `84532` (Base Sepolia). |

### SIWE auth (`/auth/*`)

When `GATEWAY_JWT_SECRET` is set the gateway exposes:

- `POST /auth/siwe/nonce` — issues a single-use nonce, stored in Redis with a 5
  minute TTL.
- `POST /auth/siwe/verify` — accepts `{ "message": "<EIP-4361>", "signature":
  "0x..." }`, validates domain, chain id, time window and signature, atomically
  consumes the nonce, mints a JWT and seeds a Redis session.
- `POST /auth/logout` — invalidates the session keyed by the JWT's `jti`.

Protected routes can opt in to authentication via `auth.RequireAuth(handlers)`.
The middleware double-checks the JWT signature **and** the Redis session, so
deleting a session immediately revokes every token it backs.

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
