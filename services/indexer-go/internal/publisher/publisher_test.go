package publisher_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	natsserver "github.com/nats-io/nats-server/v2/server"
	srvtest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"

	contractabi "github.com/0xDevNinja/titular/services/indexer-go/internal/abi"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/decoder"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/publisher"
)

// Tests in this file run against an embedded nats-server with JetStream
// enabled, so they exercise the same code paths as a production deploy
// without depending on docker. The server is started once per test and
// torn down via t.Cleanup; each test gets its own port via
// natsserver.RANDOM_PORT.

// runServer boots an embedded nats-server with JetStream backed by a
// temp dir. The temp dir is removed on cleanup so successive runs do
// not reuse stream state. Callers receive a connected *nats.Conn and a
// nats.JetStreamContext ready to use.
func runServer(t *testing.T) (*nats.Conn, nats.JetStreamContext) {
	t.Helper()
	dir := t.TempDir()

	opts := srvtest.DefaultTestOptions
	opts.Port = -1 // pick a random free port
	opts.JetStream = true
	opts.StoreDir = filepath.Join(dir, "jetstream")

	srv := srvtest.RunServer(&opts)
	t.Cleanup(func() {
		srv.Shutdown()
		srv.WaitForShutdown()
		_ = os.RemoveAll(opts.StoreDir)
	})

	if !srv.ReadyForConnections(2 * time.Second) {
		t.Fatal("nats-server not ready")
	}

	nc, err := nats.Connect(srv.ClientURL())
	if err != nil {
		t.Fatalf("nats.Connect: %v", err)
	}
	t.Cleanup(nc.Close)

	js, err := nc.JetStream()
	if err != nil {
		t.Fatalf("nc.JetStream: %v", err)
	}

	return nc, js
}

// silenceServerLogs is a no-op compile reference to the natsserver
// package — kept so the import stays scoped if a later change wants to
// pass a *natsserver.Logger to silence stdout in CI.
var _ = natsserver.RANDOM_PORT

// ---------------------------------------------------------------------------
// EnsureStream
// ---------------------------------------------------------------------------

// Test_EnsureStream_CreatesWithExpectedConfig asserts the stream is
// created with the four documented prefixes and that a second call is a
// no-op rather than an error. This catches the foot-gun where AddStream
// returns ErrStreamNameAlreadyInUse on second boot — EnsureStream should
// hide that.
func Test_EnsureStream_CreatesWithExpectedConfig(t *testing.T) {
	_, js := runServer(t)

	if err := publisher.EnsureStream(js, publisher.StreamOptions{
		Storage: nats.MemoryStorage, // avoid touching disk in unit tests
	}); err != nil {
		t.Fatalf("EnsureStream first call: %v", err)
	}

	info, err := js.StreamInfo(publisher.StreamName)
	if err != nil {
		t.Fatalf("StreamInfo: %v", err)
	}

	got := append([]string(nil), info.Config.Subjects...)
	want := publisher.SubjectPrefixes
	if !sameSet(got, want) {
		t.Errorf("subjects = %v, want %v", got, want)
	}

	// Second call must not error: idempotent EnsureStream is a hard
	// requirement for the indexer's restart loop.
	if err := publisher.EnsureStream(js, publisher.StreamOptions{
		Storage: nats.MemoryStorage,
	}); err != nil {
		t.Errorf("EnsureStream second call: %v", err)
	}
}

// Test_EnsureStream_DetectsMissingPrefix asserts that an existing stream
// missing one of our prefixes is rejected, rather than silently
// dropping events. Mirrors the operational failure mode where a shared
// NATS cluster has a TITULAR_EVENTS stream from an older deploy.
func Test_EnsureStream_DetectsMissingPrefix(t *testing.T) {
	_, js := runServer(t)

	// Pre-create a stream with only `agents.>`.
	if _, err := js.AddStream(&nats.StreamConfig{
		Name:     publisher.StreamName,
		Subjects: []string{"agents.>"},
		Storage:  nats.MemoryStorage,
	}); err != nil {
		t.Fatalf("AddStream: %v", err)
	}

	err := publisher.EnsureStream(js, publisher.StreamOptions{Storage: nats.MemoryStorage})
	if err == nil {
		t.Fatal("EnsureStream: expected error for partial stream")
	}
	if !strings.Contains(err.Error(), "missing subject") {
		t.Errorf("err = %v, want missing subject", err)
	}
}

// ---------------------------------------------------------------------------
// Publisher.OnEvent — round-trip per top-level prefix
// ---------------------------------------------------------------------------

// Test_OnEvent_RoundTripPerPrefix dispatches one synthetic event per
// top-level prefix (`agents.*`, `jobs.*`, `tokens.*`) through the
// publisher and verifies that a JetStream consumer receives the
// envelope, on the right subject, with the right payload. The `trades.*`
// prefix has no decoded events yet (M2 BondingCurve is not in the
// dispatch table) so we only assert the stream is provisioned for it
// in the EnsureStream test above.
func Test_OnEvent_RoundTripPerPrefix(t *testing.T) {
	_, js := runServer(t)
	if err := publisher.EnsureStream(js, publisher.StreamOptions{
		Storage: nats.MemoryStorage,
	}); err != nil {
		t.Fatalf("EnsureStream: %v", err)
	}

	pub, err := publisher.New(js, publisher.Options{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Subscribe synchronously so each case can grab the expected
	// message after publishing, without juggling channels.
	subAgents, err := js.PullSubscribe("agents.>", "agents-test", nats.AckExplicit())
	if err != nil {
		t.Fatalf("PullSubscribe agents: %v", err)
	}
	subJobs, err := js.PullSubscribe("jobs.>", "jobs-test", nats.AckExplicit())
	if err != nil {
		t.Fatalf("PullSubscribe jobs: %v", err)
	}
	subTokens, err := js.PullSubscribe("tokens.>", "tokens-test", nats.AckExplicit())
	if err != nil {
		t.Fatalf("PullSubscribe tokens: %v", err)
	}

	cases := []struct {
		name        string
		event       decoder.DecodedEvent
		wantSubject string
		sub         *nats.Subscription
		check       func(t *testing.T, payload json.RawMessage)
	}{
		{
			name: "agents.registered",
			event: decoder.DecodedEvent{
				Name: "AgentRegistry.AgentRegistered",
				Payload: &contractabi.AgentRegistryAgentRegistered{
					AgentId:     big.NewInt(7),
					Controller:  common.HexToAddress("0xAAAA000000000000000000000000000000000001"),
					MetadataURI: "ipfs://bafy",
				},
				Raw: types.Log{
					BlockNumber: 100,
					TxHash:      common.HexToHash("0x01"),
					Index:       3,
				},
			},
			wantSubject: "agents.registered",
			sub:         subAgents,
			check: func(t *testing.T, raw json.RawMessage) {
				// big.Int is wire-encoded as a quoted decimal string —
				// see [bigIntStr] in publisher.go for the rationale
				// (JS consumers lose precision past 2^53).
				var got struct {
					AgentId     string `json:"AgentId"`
					Controller  string `json:"Controller"`
					MetadataURI string `json:"MetadataURI"`
				}
				if err := json.Unmarshal(raw, &got); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				if got.AgentId != "7" {
					t.Errorf("AgentId = %q, want %q", got.AgentId, "7")
				}
				if got.MetadataURI != "ipfs://bafy" {
					t.Errorf("MetadataURI = %q", got.MetadataURI)
				}
			},
		},
		{
			name: "jobs.funded",
			event: decoder.DecodedEvent{
				Name: "Escrow.Funded",
				Payload: &contractabi.EscrowFunded{
					Depositor: common.HexToAddress("0xBBBB000000000000000000000000000000000002"),
					JobId:     big.NewInt(42),
					Token:     common.HexToAddress("0xCCCC000000000000000000000000000000000003"),
					Amount:    big.NewInt(1_000_000),
				},
				Raw: types.Log{
					BlockNumber: 200,
					TxHash:      common.HexToHash("0x02"),
					Index:       1,
				},
			},
			wantSubject: "jobs.funded",
			sub:         subJobs,
			check: func(t *testing.T, raw json.RawMessage) {
				var got struct {
					JobId  string `json:"JobId"`
					Amount string `json:"Amount"`
				}
				if err := json.Unmarshal(raw, &got); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				if got.JobId != "42" {
					t.Errorf("JobId = %q, want %q", got.JobId, "42")
				}
				if got.Amount != "1000000" {
					t.Errorf("Amount = %q, want %q", got.Amount, "1000000")
				}
			},
		},
		{
			name: "tokens.fee_split",
			event: decoder.DecodedEvent{
				Name: "FeeSplitter.FeeSplit",
				Payload: &contractabi.FeeSplitterFeeSplit{
					Token:         common.HexToAddress("0xDDDD000000000000000000000000000000000004"),
					Primary:       common.HexToAddress("0xEEEE000000000000000000000000000000000005"),
					Schedule:      2,
					PrimaryAmount: big.NewInt(900),
					TreasuryAmount: big.NewInt(50),
					BuybackAmount:  big.NewInt(50),
				},
				Raw: types.Log{
					BlockNumber: 300,
					TxHash:      common.HexToHash("0x03"),
					Index:       7,
				},
			},
			wantSubject: "tokens.fee_split",
			sub:         subTokens,
			check: func(t *testing.T, raw json.RawMessage) {
				var got struct {
					Schedule      uint8  `json:"Schedule"`
					PrimaryAmount string `json:"PrimaryAmount"`
				}
				if err := json.Unmarshal(raw, &got); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				if got.Schedule != 2 {
					t.Errorf("Schedule = %d, want 2", got.Schedule)
				}
				if got.PrimaryAmount != "900" {
					t.Errorf("PrimaryAmount = %q, want %q", got.PrimaryAmount, "900")
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := pub.OnEvent(context.Background(), c.event); err != nil {
				t.Fatalf("OnEvent: %v", err)
			}

			msgs, err := c.sub.Fetch(1, nats.MaxWait(2*time.Second))
			if err != nil {
				t.Fatalf("Fetch: %v", err)
			}
			if len(msgs) != 1 {
				t.Fatalf("Fetch returned %d messages", len(msgs))
			}
			msg := msgs[0]
			if msg.Subject != c.wantSubject {
				t.Errorf("Subject = %q, want %q", msg.Subject, c.wantSubject)
			}

			// Envelope shape sanity check.
			var env struct {
				Name        string          `json:"name"`
				BlockNumber uint64          `json:"block_number"`
				TxHash      string          `json:"tx_hash"`
				LogIndex    uint            `json:"log_index"`
				Payload     json.RawMessage `json:"payload"`
			}
			if err := json.Unmarshal(msg.Data, &env); err != nil {
				t.Fatalf("unmarshal envelope: %v", err)
			}
			if env.Name != c.event.Name {
				t.Errorf("envelope name = %q, want %q", env.Name, c.event.Name)
			}
			if env.BlockNumber != c.event.Raw.BlockNumber {
				t.Errorf("block = %d, want %d", env.BlockNumber, c.event.Raw.BlockNumber)
			}
			if env.LogIndex != c.event.Raw.Index {
				t.Errorf("logIndex = %d, want %d", env.LogIndex, c.event.Raw.Index)
			}

			// Per-event payload assertion.
			c.check(t, env.Payload)

			// Idempotency header round-trips on the message.
			wantID := publisher.DedupID(c.event.Raw.TxHash.Hex(), c.event.Raw.Index)
			if got := msg.Header.Get(nats.MsgIdHdr); got != wantID {
				t.Errorf("Nats-Msg-Id = %q, want %q", got, wantID)
			}

			if err := msg.Ack(); err != nil {
				t.Errorf("Ack: %v", err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Idempotency: re-publish within the dedup window must not produce a
// second message in the stream.
// ---------------------------------------------------------------------------

func Test_OnEvent_DedupsByTxAndLogIndex(t *testing.T) {
	_, js := runServer(t)
	if err := publisher.EnsureStream(js, publisher.StreamOptions{
		Storage:         nats.MemoryStorage,
		DuplicateWindow: time.Minute,
	}); err != nil {
		t.Fatalf("EnsureStream: %v", err)
	}

	pub, err := publisher.New(js, publisher.Options{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	event := decoder.DecodedEvent{
		Name: "HookRegistry.HookRegistered",
		Payload: &contractabi.HookRegistryHookRegistered{
			Hook: common.HexToAddress("0x9999000000000000000000000000000000000099"),
			Name: "fund-transfer",
		},
		Raw: types.Log{
			BlockNumber: 555,
			TxHash:      common.HexToHash("0xdedupe"),
			Index:       0,
		},
	}

	// Publish twice with the same dedup ID.
	for i := 0; i < 2; i++ {
		if err := pub.OnEvent(context.Background(), event); err != nil {
			t.Fatalf("OnEvent #%d: %v", i, err)
		}
	}

	// Stream should hold exactly one message.
	info, err := js.StreamInfo(publisher.StreamName)
	if err != nil {
		t.Fatalf("StreamInfo: %v", err)
	}
	if info.State.Msgs != 1 {
		t.Errorf("stream msgs = %d, want 1 (dedup window)", info.State.Msgs)
	}
}

// ---------------------------------------------------------------------------
// Unmapped events: default = error, opt-in skip = silent drop.
// ---------------------------------------------------------------------------

func Test_OnEvent_UnmappedEvent(t *testing.T) {
	_, js := runServer(t)
	if err := publisher.EnsureStream(js, publisher.StreamOptions{
		Storage: nats.MemoryStorage,
	}); err != nil {
		t.Fatalf("EnsureStream: %v", err)
	}

	t.Run("default returns ErrUnmappedEvent", func(t *testing.T) {
		pub, err := publisher.New(js, publisher.Options{})
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		err = pub.OnEvent(context.Background(), decoder.DecodedEvent{Name: "Mystery.Event", Raw: types.Log{}})
		if !errors.Is(err, publisher.ErrUnmappedEvent) {
			t.Fatalf("err = %v, want ErrUnmappedEvent", err)
		}
	})

	t.Run("SkipUnmapped silently drops", func(t *testing.T) {
		pub, err := publisher.New(js, publisher.Options{SkipUnmapped: true})
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		if err := pub.OnEvent(context.Background(), decoder.DecodedEvent{Name: "Mystery.Event", Raw: types.Log{}}); err != nil {
			t.Errorf("err = %v, want nil with SkipUnmapped", err)
		}
	})
}

// ---------------------------------------------------------------------------
// Subject table integrity
// ---------------------------------------------------------------------------

// Test_AllSubjects_PrefixesAndUniqueness asserts every mapped subject
// belongs to one of the registered prefixes and that no two events
// share a subject. Reverse drift (e.g. typo `agent.registered` for
// `agents.registered`) would silently route messages outside the stream
// — better caught here.
func Test_AllSubjects_PrefixesAndUniqueness(t *testing.T) {
	rev := make(map[string]string)
	for evName, subj := range publisher.AllSubjects() {
		var matched bool
		for _, p := range publisher.SubjectPrefixes {
			// p is `agents.>` etc. Strip the wildcard for the prefix
			// check; concrete subjects must start with the dot-form.
			dot := p[:len(p)-1]
			if strings.HasPrefix(subj, dot) {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("event %q -> subject %q has no registered prefix", evName, subj)
		}
		if other, dup := rev[subj]; dup {
			t.Errorf("subject %q used by both %q and %q", subj, other, evName)
		}
		rev[subj] = evName
	}
}

// Test_AllSubjects_CoversDispatchTable asserts every event the decoder
// knows about has a subject mapping. Driven off [decoder.KnownEvents]
// (the live dispatch table) rather than a duplicated literal list, so a
// regen that adds a binding + decoder but forgets the subject mapping
// fails CI here. The reverse — a stale subject for an event the decoder
// no longer emits — is also surfaced by comparing set sizes.
func Test_AllSubjects_CoversDispatchTable(t *testing.T) {
	known := decoder.KnownEvents()
	subjects := publisher.AllSubjects()

	if got, want := len(subjects), len(known); got != want {
		t.Errorf("subject table size = %d, want %d (decoder.KnownEvents)", got, want)
	}

	for _, name := range known {
		if _, ok := publisher.SubjectFor(name); !ok {
			t.Errorf("missing subject mapping for decoder event %q", name)
		}
	}

	// Reverse direction: every entry in the subject table must
	// correspond to a decoder. A leftover from a renamed event would
	// hit ErrUnmappedEvent at runtime — surface it here.
	knownSet := make(map[string]struct{}, len(known))
	for _, n := range known {
		knownSet[n] = struct{}{}
	}
	for name := range subjects {
		if _, ok := knownSet[name]; !ok {
			t.Errorf("subject mapping %q has no decoder", name)
		}
	}
}

// Test_DedupID_Format pins the wire format. Any change here is a
// breaking deploy and must be intentional. The format is lowercase so
// callers that capitalise hex differently still hash to the same key.
func Test_DedupID_Format(t *testing.T) {
	got := publisher.DedupID("0xDeadBeef", 12)
	want := "0xdeadbeef:12"
	if got != want {
		t.Fatalf("DedupID = %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// bigint precision: payloads with values > 2^53 must round-trip without loss
// ---------------------------------------------------------------------------

// Test_OnEvent_BigIntStringEncoding asserts that large *big.Int payload
// fields (here 2^200) survive the JetStream round-trip as decimal
// strings and reconstruct exactly via big.Int.SetString. This catches
// the regression we'd hit if a future refactor reverted to the default
// big.Int.MarshalJSON, which emits a JSON number and silently truncates
// past 2^53 in JS consumers (gateway, SDK).
func Test_OnEvent_BigIntStringEncoding(t *testing.T) {
	_, js := runServer(t)
	if err := publisher.EnsureStream(js, publisher.StreamOptions{
		Storage: nats.MemoryStorage,
	}); err != nil {
		t.Fatalf("EnsureStream: %v", err)
	}

	pub, err := publisher.New(js, publisher.Options{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	sub, err := js.PullSubscribe("jobs.funded", "jobs-bigint", nats.AckExplicit())
	if err != nil {
		t.Fatalf("PullSubscribe: %v", err)
	}

	// 2^200 is well above the 2^53 boundary where JSON numbers lose
	// precision in IEEE-754 doubles. We reuse Escrow.Funded because
	// it has both an indexed *big.Int (JobId) and a non-indexed one
	// (Amount), so the test covers both reflection paths.
	huge := new(big.Int).Lsh(big.NewInt(1), 200)
	want := huge.Text(10)

	event := decoder.DecodedEvent{
		Name: "Escrow.Funded",
		Payload: &contractabi.EscrowFunded{
			Depositor: common.HexToAddress("0xBBBB000000000000000000000000000000000002"),
			JobId:     huge,
			Token:     common.HexToAddress("0xCCCC000000000000000000000000000000000003"),
			Amount:    huge,
		},
		Raw: types.Log{
			BlockNumber: 999,
			TxHash:      common.HexToHash("0xbig"),
			Index:       0,
		},
	}

	if err := pub.OnEvent(context.Background(), event); err != nil {
		t.Fatalf("OnEvent: %v", err)
	}

	msgs, err := sub.Fetch(1, nats.MaxWait(2*time.Second))
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("Fetch returned %d messages", len(msgs))
	}
	msg := msgs[0]
	if err := msg.Ack(); err != nil {
		t.Fatalf("Ack: %v", err)
	}

	// The raw JSON for the payload must contain the value as a quoted
	// string, never as a bare number. We check for `"<dec>"` to fail
	// loudly if the encoder regresses to a numeric form.
	rawJSON := string(msg.Data)
	wantQuoted := `"` + want + `"`
	if !strings.Contains(rawJSON, wantQuoted) {
		t.Errorf("payload JSON does not contain quoted decimal:\n  want substring: %s\n  got: %s",
			wantQuoted, rawJSON)
	}
	// And it must NOT contain the value as a bare number — i.e. the
	// digits followed by a comma or `}` instead of a closing quote.
	if strings.Contains(rawJSON, want+",") || strings.Contains(rawJSON, want+"}") {
		t.Errorf("payload JSON contains bare numeric for big.Int: %s", rawJSON)
	}

	// Round-trip: pull the strings out and rebuild the big.Int.
	var env struct {
		Payload struct {
			JobId  string `json:"JobId"`
			Amount string `json:"Amount"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(msg.Data, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	for label, got := range map[string]string{"JobId": env.Payload.JobId, "Amount": env.Payload.Amount} {
		if got != want {
			t.Errorf("%s = %q, want %q", label, got, want)
		}
		round, ok := new(big.Int).SetString(got, 10)
		if !ok {
			t.Errorf("%s: SetString(%q) failed", label, got)
			continue
		}
		if round.Cmp(huge) != 0 {
			t.Errorf("%s: round-trip mismatch: got %s, want %s", label, round.Text(10), want)
		}
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// sameSet returns true if two slices contain the same elements
// (order-insensitive). The test asserts subject configuration equality;
// the server may return subjects in a different order than we sent.
func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	seen := make(map[string]int, len(a))
	for _, x := range a {
		seen[x]++
	}
	for _, x := range b {
		seen[x]--
	}
	for _, v := range seen {
		if v != 0 {
			return false
		}
	}
	return true
}

// _ = fmt is here to keep the import alive in case future helpers use it.
var _ = fmt.Sprintf
