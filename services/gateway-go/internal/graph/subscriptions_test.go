package graph

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
	srvtest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

// runEmbeddedNATS boots an in-process nats-server with JetStream
// enabled so subscription tests exercise the same wire path the
// indexer's publisher uses. Storage is on-disk-temp so successive
// runs do not share state. Mirrors the helper in indexer-go's
// publisher_test.go but lives here because we cannot import it across
// services.
func runEmbeddedNATS(t *testing.T) *nats.Conn {
	t.Helper()
	dir := t.TempDir()

	opts := srvtest.DefaultTestOptions
	opts.Port = -1
	opts.JetStream = true
	opts.StoreDir = filepath.Join(dir, "jetstream")

	srv := srvtest.RunServer(&opts)
	t.Cleanup(func() {
		srv.Shutdown()
		srv.WaitForShutdown()
	})
	if !srv.ReadyForConnections(2 * time.Second) {
		t.Fatal("nats-server not ready")
	}
	nc, err := nats.Connect(srv.ClientURL())
	if err != nil {
		t.Fatalf("nats.Connect: %v", err)
	}
	t.Cleanup(nc.Close)
	return nc
}

// silence the imported package so future maintainers see the import
// is intentional even when the in-test helpers don't reference its
// types directly.
var _ = natsserver.RANDOM_PORT

// Test_NATSBus_TradesRoundTrip publishes a `trades.bought` message on
// the bus and asserts that a subscriber receives a decoded store.Trade.
// This is the smoke test required by #89 — the subscription path can
// be assumed working in the rest of the suite.
func Test_NATSBus_TradesRoundTrip(t *testing.T) {
	nc := runEmbeddedNATS(t)
	bus, err := NewNATSBus(nc)
	if err != nil {
		t.Fatalf("NewNATSBus: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch, err := bus.SubscribeTrades(ctx, nil)
	if err != nil {
		t.Fatalf("SubscribeTrades: %v", err)
	}

	// Publish a trade envelope shaped like the indexer's publisher.
	env := NATSEnvelope{
		Name:        "BondingCurve.Bought",
		BlockNumber: 100,
		TxHash:      "0xabcd",
		LogIndex:    3,
		Payload:     json.RawMessage(`{"Trader":"0x1111111111111111111111111111111111111111","Curve":"0x2222222222222222222222222222222222222222","QuoteIn":"1000","AgentOut":"999","Fee":"10"}`),
	}
	body, _ := json.Marshal(env)
	if err := nc.Publish("trades.bought", body); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	select {
	case <-ctx.Done():
		t.Fatal("timed out waiting for trade")
	case got, ok := <-ch:
		if !ok {
			t.Fatal("channel closed without a value")
		}
		if got.Trader != "0x1111111111111111111111111111111111111111" {
			t.Errorf("trader: got %q", got.Trader)
		}
		if got.Side != "buy" {
			t.Errorf("side: got %q, want buy", got.Side)
		}
		if got.Fee != "10" {
			t.Errorf("fee: got %q", got.Fee)
		}
		if got.BlockNumber != 100 {
			t.Errorf("blockNumber: got %d", got.BlockNumber)
		}
	}
}

// Test_NATSBus_FilterDropsOtherToken pins the server-side filter
// behaviour: a subscriber asking for a specific bonding curve must
// not see trades on a different curve. Catches the regression where
// the SubjectFilter closure is dropped when wiring new resolvers.
func Test_NATSBus_FilterDropsOtherToken(t *testing.T) {
	nc := runEmbeddedNATS(t)
	bus, err := NewNATSBus(nc)
	if err != nil {
		t.Fatalf("NewNATSBus: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	want := "0x2222222222222222222222222222222222222222"
	filter := func(_ string, env *NATSEnvelope) bool {
		return strings.Contains(string(env.Payload), `"Curve":"`+want+`"`)
	}
	ch, err := bus.SubscribeTrades(ctx, filter)
	if err != nil {
		t.Fatalf("SubscribeTrades: %v", err)
	}

	// Trade for the wrong curve — should be filtered out.
	wrong := NATSEnvelope{
		Name:    "BondingCurve.Bought",
		Payload: json.RawMessage(`{"Trader":"0x9","Curve":"0x9999999999999999999999999999999999999999","Fee":"1"}`),
	}
	wb, _ := json.Marshal(wrong)
	_ = nc.Publish("trades.bought", wb)

	// Trade for the right curve — must arrive.
	right := NATSEnvelope{
		Name:    "BondingCurve.Bought",
		Payload: json.RawMessage(`{"Trader":"0x9","Curve":"` + want + `","Fee":"1"}`),
	}
	rb, _ := json.Marshal(right)
	_ = nc.Publish("trades.bought", rb)

	select {
	case <-ctx.Done():
		t.Fatal("timed out — filter may have dropped both")
	case got, ok := <-ch:
		if !ok {
			t.Fatal("channel closed")
		}
		if got.Curve != want {
			t.Errorf("curve: got %q, want %q", got.Curve, want)
		}
	}
}

// Test_NATSBus_JobsDecodesSubject asserts the job event decoder honours
// the publisher's name -> phase mapping. Funded events come with a
// budget; Released events do not.
func Test_NATSBus_JobsDecodesSubject(t *testing.T) {
	nc := runEmbeddedNATS(t)
	bus, err := NewNATSBus(nc)
	if err != nil {
		t.Fatalf("NewNATSBus: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch, err := bus.SubscribeJobs(ctx, nil)
	if err != nil {
		t.Fatalf("SubscribeJobs: %v", err)
	}

	env := NATSEnvelope{
		Name:        "Escrow.Funded",
		BlockNumber: 500,
		TxHash:      "0xfeed",
		LogIndex:    4,
		Payload:     json.RawMessage(`{"Depositor":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","JobId":"42","Token":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","Amount":"1000000"}`),
	}
	body, _ := json.Marshal(env)
	if err := nc.Publish("jobs.funded", body); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	select {
	case <-ctx.Done():
		t.Fatal("timed out")
	case got := <-ch:
		if got.OnChainID != "42" {
			t.Errorf("onChainId: got %q", got.OnChainID)
		}
		if got.Budget != "1000000" {
			t.Errorf("budget: got %q", got.Budget)
		}
	}
}

func Test_NilBus_IsClosedImmediately(t *testing.T) {
	bus := NilBus()
	ch, err := bus.SubscribeTrades(context.Background(), nil)
	if err != nil {
		t.Fatalf("SubscribeTrades: %v", err)
	}
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("nilBus emitted a value")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("nilBus channel was not closed")
	}
}

// Test_NATSBus_RejectsNilConn pins the constructor's nil guard.
func Test_NATSBus_RejectsNilConn(t *testing.T) {
	if _, err := NewNATSBus(nil); err == nil {
		t.Fatal("expected error")
	}
}

// Test_NewSchema_DefaultsToNilBus confirms a Resolver constructed
// without a Bus is fine — the schema constructor injects NilBus()
// rather than crashing later.
func Test_NewSchema_DefaultsToNilBus(t *testing.T) {
	r := &Resolver{Store: &fakeStore{}}
	if _, err := NewSchema(r); err != nil {
		t.Fatalf("NewSchema: %v", err)
	}
	if r.Bus == nil {
		t.Error("Bus must default to NilBus, not nil")
	}
}

// Test_RunNewTrade_NilBusEndsImmediately exercises the resolver path:
// a subscription resolver running against the no-op bus must yield
// no values and close cleanly. This is the GraphQL-side mirror of
// Test_NilBus_IsClosedImmediately.
func Test_RunNewTrade_NilBusEndsImmediately(t *testing.T) {
	r := &Resolver{Store: &fakeStore{}, Bus: NilBus()}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	ch, err := r.RunNewTrade(ctx, "")
	if err != nil {
		t.Fatalf("RunNewTrade: %v", err)
	}
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("got a value from nilBus")
		}
	case <-ctx.Done():
		t.Fatal("channel did not close")
	}
}
