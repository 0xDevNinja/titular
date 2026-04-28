# indexer-go / subscriber

EVM block-head subscriber for the Titular indexer.

## What it does

Subscribes to `newHeads` over a WebSocket RPC connection (typically
Alchemy / Base Sepolia), waits for blocks to age past a configurable
confirmation depth (default 12), detects reorgs at the confirmed tip, and
emits canonical headers to a downstream `Sink` in monotonic order.

## API

```go
client, err := ethclient.DialContext(ctx, wsURL)
// ...
sub, err := subscriber.New(client, mySink, subscriber.Config{
    Confirmations:  12,
    ReorgWindow:    64,
    InitialBackoff: time.Second,
    MaxBackoff:     30 * time.Second,
    Logger:         slog.Default(),
})
// Run blocks until ctx is cancelled or a fatal error occurs.
err = sub.Run(ctx)
```

`Sink` is a single-method interface:

```go
type Sink interface {
    Emit(ctx context.Context, b ConfirmedBlock) error
}
```

`ConfirmedBlock.Reorged == true` signals that the same `Number` has been
emitted before with a different `Hash`; the sink must roll back any state
derived from the previous hash.

## Behaviour

| concern             | how                                                  |
|---------------------|------------------------------------------------------|
| connection drops    | exponential backoff (`InitialBackoff` → `MaxBackoff`) and resubscribe |
| skipped heads       | always backfilled via `HeaderByNumber` so no height is ever silently dropped |
| reorg detection     | walks `[tip-ReorgWindow .. tip]` after every head; first matching hash terminates |
| reorg deeper than window | fatal — `Run` returns; operator must drain and resync |
| sink error          | fatal — `Run` returns                                |
| context cancel      | clean shutdown — `Run` returns `ctx.Err()`           |

## Health gauges

Three atomic counters are exposed for liveness probes / Prometheus:

- `LastSeenBlock()` — most recent head observed on the subscription.
- `LastConfirmedBlock()` — most recent block successfully emitted.
- `LastReorgAt()` — unix timestamp of the most recent reorg (0 if never).

## Testing

The `EthClient` interface is intentionally minimal so tests can drive the
loop with a mock — see `subscriber_test.go`. There is no integration test
in this package; chain-level integration belongs in a follow-up phase that
wires the subscriber into `cmd/indexer/main.go`.
