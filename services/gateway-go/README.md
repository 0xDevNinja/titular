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
| `GATEWAY_SIWS_DOMAIN` | The value the SIWS (Sign-In With Solana) message MUST declare. Falls back to `GATEWAY_SIWE_DOMAIN` when unset so single-domain deployments configure once. |
| `GATEWAY_SIWS_CLUSTER` | Solana cluster the SIWS message MUST declare. One of `devnet`, `mainnet-beta`, `testnet`. Setting this enables `/auth/siws/*`. No default — operators must opt in explicitly. |
| `GATEWAY_DATABASE_URL` | Read-only Postgres DSN (`postgres://user:pw@host:port/db`). Setting this enables the M4 read-only REST API under `/api/v1`. Unset disables the API; the binary still serves the M2/M3 fixture surface. |

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

### SIWS auth (`/auth/siws/*`)

When `GATEWAY_JWT_SECRET` and `GATEWAY_SIWS_CLUSTER` are both set, the gateway
additionally exposes the Sign-In-With-Solana flow:

- `POST /auth/siws/nonce` — same shape and TTL as the SIWE nonce. The nonce
  namespace is shared between SIWE and SIWS so a single nonce is single-use
  across both paths.
- `POST /auth/siws/verify` — accepts `{ "message": "<SIWS>", "signature":
  "<base58>" }`, validates domain, cluster (`devnet`/`mainnet-beta`/`testnet`),
  time window and the ed25519 signature, atomically consumes the nonce, mints a
  JWT and seeds a Redis session. Verification uses Go's stdlib
  `crypto/ed25519`, which is constant-time per the language contract.
- Logout is shared with the SIWE path (`POST /auth/logout`); a SIWS-minted token
  logs out via the same endpoint.

The SIWS message format is a port of EIP-4361 with a `Chain ID` (or `Cluster`)
field naming a Solana network instead of an EVM chain id; the handler accepts
either field name. Phantom and Backpack clients can be wired in directly: their
`signMessage` returns base58, which is the format the verify endpoint
expects.

### Read-only REST API (`/api/v1/*`)

When `GATEWAY_DATABASE_URL` is set the gateway opens a pgx connection pool and
mounts a small read-only API for indexer-derived state:

| Endpoint | Description |
|---|---|
| `GET /api/v1/agents?cursor=&limit=&kind=` | Paginated agent list. `kind` filters to `launchpad` or `acp`. |
| `GET /api/v1/agents/:id` | Lookup by primary key OR by `kind:agent_id` slug (e.g. `acp:42`). |
| `GET /api/v1/trades?agent_token=&from=&to=&cursor=&limit=` | Paginated trade list. `agent_token` is a 0x-prefixed 20-byte hex address; `from` / `to` are RFC 3339 timestamps. |
| `GET /api/v1/jobs?status=&cursor=&limit=` | Paginated job list. `status` is one of `created`, `funded`, `active`, `completed`, `cancelled`, `disputed`, `released`, `resolved`. |
| `GET /api/v1/stats` | Aggregate counts plus `last_block_indexed`. |

Pagination is opaque-cursor based: pass `next_cursor` from the previous
response back as `cursor`. Limits default to 50 and are capped at 200; values
outside `[1, 200]` return 400. Unknown query parameters return 400 to surface
client typos. The endpoints are read-only and live outside the SIWE wall —
deployments that want them gated should wrap the group in
`auth.RequireAuth(cfg.Auth)`.

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
