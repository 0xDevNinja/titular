//go:build integration

package integration

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// anvilDevKey is anvil's account 0 private key. Stable across versions and
// safe to embed: it is the public dev key shipped with foundry, never used
// outside ephemeral test instances. Address: 0xf39Fd6e51aad88F6F4ce6aB8827279cfFFb92266.
const anvilDevKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" // gitleaks:allow

// anvilChainID is the default chain id anvil exposes on `--chain-id` left
// at its zero value. Pinned here so tx signing matches.
const anvilChainID = 31337

// anvil is a process-level handle to a foundry anvil instance with both
// HTTP-RPC and WebSocket-RPC listeners on the same port. It owns the
// child process lifetime via t.Cleanup so a test failure cannot leak
// daemons across subsequent runs.
type anvil struct {
	t *testing.T

	cmd  *exec.Cmd
	port int

	httpURL string
	wsURL   string

	// rpcHTTP is the HTTP JSON-RPC client used for admin (`anvil_*`,
	// `evm_*`) calls and for synchronous tx submission. The subscriber
	// receives a separate ethclient bound to wsURL — keeping them
	// separate makes it obvious which connection drives which side of
	// the pipeline.
	rpcHTTP *rpc.Client
}

// startAnvil boots anvil on a random local port and waits for it to
// accept eth_blockNumber. On any error the process is killed before
// returning so a partial start does not leak.
//
// The caller MUST register t.Cleanup callbacks elsewhere if they install
// additional state on the chain — startAnvil itself only owns the daemon
// lifetime. Block timestamps progress automatically only when blocks are
// mined; we set --block-time=0 (default) and drive blocks explicitly with
// `anvil_mine` so the test is deterministic.
func startAnvil(t *testing.T) *anvil {
	t.Helper()

	if _, err := exec.LookPath("anvil"); err != nil {
		t.Skipf("anvil not on PATH (install foundry to run this test): %v", err)
		return nil
	}

	port, err := freeTCPPort()
	if err != nil {
		t.Fatalf("freeTCPPort: %v", err)
	}

	// --silent suppresses the per-block log spam that anvil emits to
	// stderr by default. We still redirect both streams to discard so
	// the test output stays clean even on noisy versions of anvil that
	// ignore --silent for some events.
	cmd := exec.Command(
		"anvil",
		"--port", strconv.Itoa(port),
		"--host", "127.0.0.1",
		"--chain-id", strconv.Itoa(anvilChainID),
		"--accounts", "5",
		// Disable auto-mining: every block must be produced by an
		// explicit `anvil_mine` from the test. With auto-mining left
		// on, anvil seals a block as soon as `eth_sendRawTransaction`
		// arrives AND the test's own `anvil_mine` calls produce more
		// blocks, leaving the mempool/chain ordering ambiguous in a
		// way that defeats the receipt → re-read invariants this test
		// depends on.
		"--no-mining",
		"--silent",
	)
	// Anvil's stdout/stderr can be useful when debugging a CI failure
	// but we keep them off by default to keep the test output clean.
	// Set ANVIL_TEST_LOGS=1 in the environment to surface them.
	if os.Getenv("ANVIL_TEST_LOGS") != "" {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("anvil start: %v", err)
	}

	a := &anvil{
		t:       t,
		cmd:     cmd,
		port:    port,
		httpURL: fmt.Sprintf("http://127.0.0.1:%d", port),
		wsURL:   fmt.Sprintf("ws://127.0.0.1:%d", port),
	}

	t.Cleanup(func() {
		// Best-effort kill: SIGTERM, then SIGKILL after a grace period.
		// We never block the test exit on anvil shutdown — if it hangs
		// the OS will reap the process when the test binary exits.
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_ = cmd.Wait()
	})

	if err := a.waitReady(t, 10*time.Second); err != nil {
		t.Fatalf("anvil never became ready: %v", err)
	}

	c, err := rpc.DialContext(context.Background(), a.httpURL)
	if err != nil {
		t.Fatalf("rpc.Dial(http): %v", err)
	}
	a.rpcHTTP = c
	t.Cleanup(c.Close)

	return a
}

// freeTCPPort asks the kernel for an unused port via a short-lived
// listener. We close immediately and pass the port to anvil; there is a
// small race window in which another process could grab the port, but
// for a developer-machine test the convenience outweighs the risk and
// `cmd.Start` would fail loudly on collision.
func freeTCPPort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// waitReady polls the HTTP RPC endpoint until eth_blockNumber returns a
// non-error response. Anvil typically opens its listener before it has
// finalised the genesis block, so a naive Dial-and-go can race the first
// request.
func (a *anvil) waitReady(t *testing.T, timeout time.Duration) error {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		c, err := rpc.DialContext(context.Background(), a.httpURL)
		if err == nil {
			var hex string
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			err = c.CallContext(ctx, &hex, "eth_blockNumber")
			cancel()
			c.Close()
			if err == nil {
				return nil
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	return errors.New("anvil did not respond to eth_blockNumber within deadline")
}

// dialEth returns a fresh ethclient bound to the anvil WebSocket. Each
// caller gets its own connection so a subscription drop cannot affect
// the admin RPC path.
func (a *anvil) dialEth(t *testing.T) *ethclient.Client {
	t.Helper()
	c, err := ethclient.DialContext(context.Background(), a.wsURL)
	if err != nil {
		t.Fatalf("ethclient.Dial(ws): %v", err)
	}
	t.Cleanup(c.Close)
	return c
}

// dialEthHTTP returns a fresh ethclient bound to the anvil HTTP RPC
// endpoint. Used for non-subscription calls (eth_getLogs, eth_call,
// eth_sendRawTransaction) so the WebSocket connection that drives the
// subscriber stays free of head-of-line blocking from request bursts.
func (a *anvil) dialEthHTTP(t *testing.T) *ethclient.Client {
	t.Helper()
	c, err := ethclient.DialContext(context.Background(), a.httpURL)
	if err != nil {
		t.Fatalf("ethclient.Dial(http): %v", err)
	}
	t.Cleanup(c.Close)
	return c
}

// mine advances anvil by `n` blocks via `anvil_mine`. n=0 is a no-op
// (anvil treats it as mine-one in some versions; we explicitly forbid
// the ambiguity).
func (a *anvil) mine(ctx context.Context, n int) error {
	if n <= 0 {
		return fmt.Errorf("anvil.mine: n must be > 0, got %d", n)
	}
	hex := "0x" + strconv.FormatInt(int64(n), 16)
	return a.rpcHTTP.CallContext(ctx, nil, "anvil_mine", hex)
}

// snapshot captures full chain state and returns the snapshot id.
// `evm_revert` consumes the id; subsequent reverts to the same id fail.
func (a *anvil) snapshot(ctx context.Context) (string, error) {
	var id string
	if err := a.rpcHTTP.CallContext(ctx, &id, "evm_snapshot"); err != nil {
		return "", err
	}
	return id, nil
}

// revert rolls the chain back to the given snapshot id. Returns the
// boolean ack from anvil (false means the id was already consumed or
// invalid). After a successful revert, any open WebSocket subscription
// continues to work; anvil keeps the connection but the chain head is
// the snapshot's recorded head.
func (a *anvil) revert(ctx context.Context, id string) (bool, error) {
	var ok bool
	if err := a.rpcHTTP.CallContext(ctx, &ok, "evm_revert", id); err != nil {
		return false, err
	}
	return ok, nil
}

// setCode installs raw EVM bytecode at addr via `anvil_setCode`. Useful
// for deploying a test contract without paying for a full deploy tx —
// it skips the constructor, so the caller must pass runtime bytecode
// (the part after the constructor returns).
func (a *anvil) setCode(ctx context.Context, addr common.Address, runtimeCode []byte) error {
	hexCode := "0x" + hex.EncodeToString(runtimeCode)
	return a.rpcHTTP.CallContext(ctx, nil, "anvil_setCode", addr.Hex(), hexCode)
}

// blockNumber returns the current head height via eth_blockNumber.
func (a *anvil) blockNumber(ctx context.Context) (uint64, error) {
	var raw string
	if err := a.rpcHTTP.CallContext(ctx, &raw, "eth_blockNumber"); err != nil {
		return 0, err
	}
	raw = strings.TrimPrefix(raw, "0x")
	n, err := strconv.ParseUint(raw, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("parse block number %q: %w", raw, err)
	}
	return n, nil
}

// devKey returns the parsed dev private key + its address. Tests that
// sign their own transactions reuse this so we do not allocate a fresh
// keypair per call.
func devKey(t *testing.T) (*ecdsa.PrivateKey, common.Address) {
	t.Helper()
	pk, err := crypto.HexToECDSA(anvilDevKey)
	if err != nil {
		t.Fatalf("dev key parse: %v", err)
	}
	pub, ok := pk.Public().(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("dev key public part is not ECDSA")
	}
	return pk, crypto.PubkeyToAddress(*pub)
}

// signedCall builds and signs a legacy transaction calling `to` with
// `data`. Legacy (rather than EIP-1559) keeps the test independent of
// anvil's basefee schedule — anvil accepts both, but legacy txs do not
// carry a `chainId == 0` ambiguity in older versions. Nonce is fetched
// from the HTTP RPC; gas is set to a comfortable 200k for a single
// LOG2 plus calldata, well below anvil's 30M block gas limit.
func signedCall(
	ctx context.Context,
	a *anvil,
	pk *ecdsa.PrivateKey,
	from, to common.Address,
	data []byte,
) (*types.Transaction, error) {
	cli := ethclient.NewClient(a.rpcHTTP)
	nonce, err := cli.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, fmt.Errorf("pending nonce: %w", err)
	}
	tx := types.NewTransaction(
		nonce,
		to,
		big.NewInt(0),  // value
		200_000,        // gas limit
		big.NewInt(1e9), // gas price (1 gwei is fine on anvil)
		data,
	)
	signer := types.NewEIP155Signer(big.NewInt(anvilChainID))
	signed, err := types.SignTx(tx, signer, pk)
	if err != nil {
		return nil, fmt.Errorf("sign tx: %w", err)
	}
	return signed, nil
}

// sendAndMine submits tx via eth_sendRawTransaction, mines exactly one
// block to include it, and returns the receipt. Returning the receipt
// rather than the hash spares callers a second RPC hop and lets them
// assert on (BlockNumber, BlockHash, Status) directly.
func sendAndMine(
	ctx context.Context,
	a *anvil,
	tx *types.Transaction,
) (*types.Receipt, error) {
	cli := ethclient.NewClient(a.rpcHTTP)
	if err := cli.SendTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("send raw: %w", err)
	}
	if err := a.mine(ctx, 1); err != nil {
		return nil, fmt.Errorf("mine after send: %w", err)
	}
	// Anvil includes the tx in the just-mined block; receipts are
	// available immediately. We retry briefly to absorb the race
	// between send and mine when the local CPU is busy.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		r, err := cli.TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			return r, nil
		}
		time.Sleep(20 * time.Millisecond)
	}
	return nil, fmt.Errorf("receipt for %s not available within deadline", tx.Hash().Hex())
}
