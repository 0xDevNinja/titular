// Package integration holds chain-level end-to-end tests for the indexer
// pipeline (subscriber → decoder → publisher). The whole package is gated
// behind the `integration` build tag so plain `go test ./...` skips it
// entirely; CI opts in via `go test -tags integration`.
//
// These tests stand up real moving parts:
//
//   - foundry's `anvil` as the EVM execution layer (out-of-process, WS+HTTP),
//   - an embedded `nats-server` with JetStream for the publisher boundary,
//   - the subscriber and decoder packages exactly as wired in production.
//
// The point is to exercise reorg semantics end-to-end with the same code
// paths a deployed indexer follows. Mocks are reserved for the per-package
// unit tests; this package does not import any of them.
//
// Local devs need foundry installed (`anvil` on $PATH). CI uses
// `foundry-rs/foundry-toolchain@v1` pinned to `version: stable`. Tests skip
// rather than fail when anvil is missing so a developer without foundry
// can still run the rest of the suite.
package integration
