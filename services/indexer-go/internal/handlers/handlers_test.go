package handlers_test

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	contractabi "github.com/0xDevNinja/titular/services/indexer-go/internal/abi"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/handlers"
)

// ----------------------------------------------------------------------------
// In-memory Store implementation used by all tests
// ----------------------------------------------------------------------------

type memStore struct {
	mu          sync.Mutex
	processed   map[string]bool // key: txHash+logIndex
	agents      []handlers.AgentRecord
	trades      []handlers.TradeRecord
	graduations []handlers.GraduationRecord
}

func newMemStore() *memStore {
	return &memStore{processed: make(map[string]bool)}
}

func logKey(txHash common.Hash, idx uint) string {
	return txHash.Hex() + ":" + string(rune(idx))
}

func (s *memStore) IsLogProcessed(_ context.Context, txHash common.Hash, logIndex uint) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.processed[logKey(txHash, logIndex)], nil
}

func (s *memStore) UpsertAgent(_ context.Context, rec handlers.AgentRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.agents = append(s.agents, rec)
	return nil
}

func (s *memStore) InsertTrade(_ context.Context, rec handlers.TradeRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.trades = append(s.trades, rec)
	return nil
}

func (s *memStore) MarkGraduated(_ context.Context, rec handlers.GraduationRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.graduations = append(s.graduations, rec)
	return nil
}

// ACP v2 Store methods — no-op stubs so memStore satisfies the full Store interface.
// ACP-specific tests use acpMemStore (defined in acp_test.go) which overrides these.

func (s *memStore) UpsertACPAgent(_ context.Context, _ handlers.ACPAgentRegisteredRecord) error {
	return nil
}

func (s *memStore) InsertACPScore(_ context.Context, _ handlers.ACPScorePostedRecord) error {
	return nil
}

func (s *memStore) UpdateACPController(_ context.Context, _ handlers.ACPControllerTransferRecord) error {
	return nil
}

func (s *memStore) UpsertACPJob(_ context.Context, _ handlers.ACPJobRecord) error {
	return nil
}

func (s *memStore) UpdateACPJobAccepted(_ context.Context, _ handlers.ACPJobAcceptedRecord) error {
	return nil
}

func (s *memStore) UpdateACPResultSubmitted(_ context.Context, _ handlers.ACPResultSubmittedRecord) error {
	return nil
}

func (s *memStore) UpdateACPJobCompleted(_ context.Context, _ handlers.ACPJobCompletedRecord) error {
	return nil
}

func (s *memStore) UpdateACPJobCancelled(_ context.Context, _ handlers.ACPJobCancelledRecord) error {
	return nil
}

func (s *memStore) InsertACPJobCreated(_ context.Context, _ handlers.ACPJobCreatedRecord) error {
	return nil
}

func (s *memStore) InsertACPEscrowFunded(_ context.Context, _ handlers.ACPEscrowFundedRecord) error {
	return nil
}

func (s *memStore) InsertACPEscrowReleased(_ context.Context, _ handlers.ACPEscrowReleasedRecord) error {
	return nil
}

func (s *memStore) InsertACPEscrowRefunded(_ context.Context, _ handlers.ACPEscrowRefundedRecord) error {
	return nil
}

// ----------------------------------------------------------------------------
// ABI helpers — parse event signatures once and reuse across tests
// ----------------------------------------------------------------------------

func mustGetABI(t *testing.T, meta *bind.MetaData) abi.ABI {
	t.Helper()
	parsed, err := meta.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	return *parsed
}

// encodeData packs non-indexed fields using the ABI encoder. fieldNames must
// match the non-indexed parameters in declaration order.
func encodeData(t *testing.T, parsed abi.ABI, eventName string, values ...interface{}) []byte {
	t.Helper()
	ev, ok := parsed.Events[eventName]
	if !ok {
		t.Fatalf("event %q not in ABI", eventName)
	}

	var args abi.Arguments
	for _, a := range ev.Inputs {
		if !a.Indexed {
			args = append(args, a)
		}
	}

	data, err := args.Pack(values...)
	if err != nil {
		t.Fatalf("pack %q: %v", eventName, err)
	}
	return data
}

// addrTopic converts an Ethereum address to a log topic (left-padded to 32 bytes).
func addrTopic(addr common.Address) common.Hash {
	return common.BytesToHash(addr.Bytes())
}

// uint256Topic converts a *big.Int to a log topic.
func uint256Topic(n *big.Int) common.Hash {
	return common.BigToHash(n)
}

// fixedTxHash returns a deterministic tx hash for a given seed.
func fixedTxHash(seed byte) common.Hash {
	var h common.Hash
	h[31] = seed
	return h
}

// errStore is a Store that always returns an error from IsLogProcessed.
type errStore struct{}

func (e *errStore) IsLogProcessed(_ context.Context, _ common.Hash, _ uint) (bool, error) {
	return false, fmt.Errorf("store unavailable")
}

func (e *errStore) UpsertAgent(_ context.Context, _ handlers.AgentRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) InsertTrade(_ context.Context, _ handlers.TradeRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) MarkGraduated(_ context.Context, _ handlers.GraduationRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) UpsertACPAgent(_ context.Context, _ handlers.ACPAgentRegisteredRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) InsertACPScore(_ context.Context, _ handlers.ACPScorePostedRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) UpdateACPController(_ context.Context, _ handlers.ACPControllerTransferRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) UpsertACPJob(_ context.Context, _ handlers.ACPJobRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) UpdateACPJobAccepted(_ context.Context, _ handlers.ACPJobAcceptedRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) UpdateACPResultSubmitted(_ context.Context, _ handlers.ACPResultSubmittedRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) UpdateACPJobCompleted(_ context.Context, _ handlers.ACPJobCompletedRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) UpdateACPJobCancelled(_ context.Context, _ handlers.ACPJobCancelledRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) InsertACPJobCreated(_ context.Context, _ handlers.ACPJobCreatedRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) InsertACPEscrowFunded(_ context.Context, _ handlers.ACPEscrowFundedRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) InsertACPEscrowReleased(_ context.Context, _ handlers.ACPEscrowReleasedRecord) error {
	return fmt.Errorf("store unavailable")
}

func (e *errStore) InsertACPEscrowRefunded(_ context.Context, _ handlers.ACPEscrowRefundedRecord) error {
	return fmt.Errorf("store unavailable")
}

// ----------------------------------------------------------------------------
// Error-path tests: store errors and malformed log data
// ----------------------------------------------------------------------------

// TestHandlers_StoreError verifies that a store error is propagated correctly
// from each handler (covering the IsLogProcessed error branch).
func TestHandlers_StoreError(t *testing.T) {
	ctx := context.Background()
	es := &errStore{}

	// Any non-empty log — the store will fail before parse is reached.
	emptyLog := types.Log{TxHash: fixedTxHash(0xAA), Index: 99}

	if err := handlers.HandleAgentLaunched(ctx, emptyLog, es); err == nil {
		t.Error("HandleAgentLaunched: expected error from store, got nil")
	}
	if err := handlers.HandleBought(ctx, emptyLog, es); err == nil {
		t.Error("HandleBought: expected error from store, got nil")
	}
	if err := handlers.HandleSold(ctx, emptyLog, es); err == nil {
		t.Error("HandleSold: expected error from store, got nil")
	}
	if err := handlers.HandleGraduated(ctx, emptyLog, es); err == nil {
		t.Error("HandleGraduated: expected error from store, got nil")
	}
}

// TestHandlers_ParseError verifies that malformed log data (empty topics + data)
// causes a parse error that is wrapped and returned.
func TestHandlers_ParseError(t *testing.T) {
	ctx := context.Background()

	// Log with wrong topic count / no data — will fail ABI unpack.
	badLog := types.Log{
		TxHash: fixedTxHash(0xBB),
		Index:  88,
		Topics: []common.Hash{}, // no topics → parse will fail
		Data:   []byte{0xBA, 0xAD},
	}

	table := []struct {
		name    string
		handler func(context.Context, types.Log, handlers.Store) error
	}{
		{"AgentLaunched", handlers.HandleAgentLaunched},
		{"Bought", handlers.HandleBought},
		{"Sold", handlers.HandleSold},
		{"Graduated", handlers.HandleGraduated},
	}

	for _, tc := range table {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			store := newMemStore()
			err := tc.handler(ctx, badLog, store)
			if err == nil {
				t.Errorf("%s: expected parse error, got nil", tc.name)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// TestHandleAgentLaunched
// ----------------------------------------------------------------------------

func TestHandleAgentLaunched(t *testing.T) {
	factoryABI := mustGetABI(t, contractabi.LaunchpadFactoryMetaData)
	eventID := factoryABI.Events["AgentLaunched"].ID

	agentID := big.NewInt(1)
	token := common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	curve := common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	creator := common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")
	lpLock := common.HexToAddress("0xDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	pair := common.HexToAddress("0xEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
	modules := [][32]byte{}

	data := encodeData(t, factoryABI, "AgentLaunched", creator, lpLock, pair, modules)

	log := types.Log{
		Address:     common.HexToAddress("0x1111111111111111111111111111111111111111"),
		BlockNumber: 42,
		TxHash:      fixedTxHash(0x01),
		Index:       0,
		Topics: []common.Hash{
			eventID,
			uint256Topic(agentID),
			addrTopic(token),
			addrTopic(curve),
		},
		Data: data,
	}

	table := []struct {
		name      string
		duplicate bool
	}{
		{name: "first call persists record", duplicate: false},
		{name: "duplicate call is skipped", duplicate: true},
	}

	store := newMemStore()
	ctx := context.Background()

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			before := len(store.agents)
			if err := handlers.HandleAgentLaunched(ctx, log, store); err != nil {
				t.Fatalf("HandleAgentLaunched: %v", err)
			}
			after := len(store.agents)

			if tc.duplicate {
				if after != before {
					t.Errorf("duplicate: expected no new agent, got %d new", after-before)
				}
				return
			}

			if after != before+1 {
				t.Fatalf("expected 1 new agent, got %d", after-before)
			}

			got := store.agents[before]
			if got.AgentID.Cmp(agentID) != 0 {
				t.Errorf("AgentID: got %s, want %s", got.AgentID, agentID)
			}
			if got.Token != token {
				t.Errorf("Token: got %s, want %s", got.Token, token)
			}
			if got.Curve != curve {
				t.Errorf("Curve: got %s, want %s", got.Curve, curve)
			}
			if got.Creator != creator {
				t.Errorf("Creator: got %s, want %s", got.Creator, creator)
			}
			if got.LpLock != lpLock {
				t.Errorf("LpLock: got %s, want %s", got.LpLock, lpLock)
			}
			if got.Pair != pair {
				t.Errorf("Pair: got %s, want %s", got.Pair, pair)
			}
			if got.BlockNum != 42 {
				t.Errorf("BlockNum: got %d, want 42", got.BlockNum)
			}
			if got.TxHash != log.TxHash {
				t.Errorf("TxHash mismatch")
			}
			if got.LogIndex != log.Index {
				t.Errorf("LogIndex: got %d, want %d", got.LogIndex, log.Index)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// TestHandleBought
// ----------------------------------------------------------------------------

func TestHandleBought(t *testing.T) {
	curveABI := mustGetABI(t, contractabi.BondingCurveMetaData)
	eventID := curveABI.Events["Bought"].ID

	buyer := common.HexToAddress("0x1100000000000000000000000000000000000001")
	quoteIn := big.NewInt(1_000_000)
	agentOut := big.NewInt(500_000)
	fee := big.NewInt(3_000)
	curveAddr := common.HexToAddress("0x2200000000000000000000000000000000000002")

	data := encodeData(t, curveABI, "Bought", quoteIn, agentOut, fee)

	log := types.Log{
		Address:     curveAddr,
		BlockNumber: 100,
		TxHash:      fixedTxHash(0x02),
		Index:       1,
		Topics: []common.Hash{
			eventID,
			addrTopic(buyer),
		},
		Data: data,
	}

	table := []struct {
		name      string
		duplicate bool
	}{
		{name: "first call persists trade", duplicate: false},
		{name: "duplicate call is skipped", duplicate: true},
	}

	store := newMemStore()
	ctx := context.Background()

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			before := len(store.trades)
			if err := handlers.HandleBought(ctx, log, store); err != nil {
				t.Fatalf("HandleBought: %v", err)
			}
			after := len(store.trades)

			if tc.duplicate {
				if after != before {
					t.Errorf("duplicate: expected no new trade, got %d new", after-before)
				}
				return
			}

			if after != before+1 {
				t.Fatalf("expected 1 new trade, got %d", after-before)
			}

			got := store.trades[before]
			if got.Direction != "buy" {
				t.Errorf("Direction: got %q, want %q", got.Direction, "buy")
			}
			if got.Trader != buyer {
				t.Errorf("Trader: got %s, want %s", got.Trader, buyer)
			}
			if got.Curve != curveAddr {
				t.Errorf("Curve: got %s, want %s", got.Curve, curveAddr)
			}
			if got.QuoteIn.Cmp(quoteIn) != 0 {
				t.Errorf("QuoteIn: got %s, want %s", got.QuoteIn, quoteIn)
			}
			if got.AgentOut.Cmp(agentOut) != 0 {
				t.Errorf("AgentOut: got %s, want %s", got.AgentOut, agentOut)
			}
			if got.Fee.Cmp(fee) != 0 {
				t.Errorf("Fee: got %s, want %s", got.Fee, fee)
			}
			if got.BlockNum != 100 {
				t.Errorf("BlockNum: got %d, want 100", got.BlockNum)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// TestHandleSold
// ----------------------------------------------------------------------------

func TestHandleSold(t *testing.T) {
	curveABI := mustGetABI(t, contractabi.BondingCurveMetaData)
	eventID := curveABI.Events["Sold"].ID

	seller := common.HexToAddress("0x3300000000000000000000000000000000000003")
	agentIn := big.NewInt(250_000)
	quoteOut := big.NewInt(490_000)
	fee := big.NewInt(1_500)
	curveAddr := common.HexToAddress("0x4400000000000000000000000000000000000004")

	data := encodeData(t, curveABI, "Sold", agentIn, quoteOut, fee)

	log := types.Log{
		Address:     curveAddr,
		BlockNumber: 200,
		TxHash:      fixedTxHash(0x03),
		Index:       2,
		Topics: []common.Hash{
			eventID,
			addrTopic(seller),
		},
		Data: data,
	}

	table := []struct {
		name      string
		duplicate bool
	}{
		{name: "first call persists trade", duplicate: false},
		{name: "duplicate call is skipped", duplicate: true},
	}

	store := newMemStore()
	ctx := context.Background()

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			before := len(store.trades)
			if err := handlers.HandleSold(ctx, log, store); err != nil {
				t.Fatalf("HandleSold: %v", err)
			}
			after := len(store.trades)

			if tc.duplicate {
				if after != before {
					t.Errorf("duplicate: expected no new trade, got %d new", after-before)
				}
				return
			}

			if after != before+1 {
				t.Fatalf("expected 1 new trade, got %d", after-before)
			}

			got := store.trades[before]
			if got.Direction != "sell" {
				t.Errorf("Direction: got %q, want %q", got.Direction, "sell")
			}
			if got.Trader != seller {
				t.Errorf("Trader: got %s, want %s", got.Trader, seller)
			}
			if got.Curve != curveAddr {
				t.Errorf("Curve: got %s, want %s", got.Curve, curveAddr)
			}
			if got.AgentIn.Cmp(agentIn) != 0 {
				t.Errorf("AgentIn: got %s, want %s", got.AgentIn, agentIn)
			}
			if got.QuoteOut.Cmp(quoteOut) != 0 {
				t.Errorf("QuoteOut: got %s, want %s", got.QuoteOut, quoteOut)
			}
			if got.Fee.Cmp(fee) != 0 {
				t.Errorf("Fee: got %s, want %s", got.Fee, fee)
			}
			if got.BlockNum != 200 {
				t.Errorf("BlockNum: got %d, want 200", got.BlockNum)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// TestHandleGraduated
// ----------------------------------------------------------------------------

func TestHandleGraduated(t *testing.T) {
	graduatorABI := mustGetABI(t, contractabi.GraduatorMetaData)
	eventID := graduatorABI.Events["Graduated"].ID

	curve := common.HexToAddress("0x5500000000000000000000000000000000000005")
	pair := common.HexToAddress("0x6600000000000000000000000000000000000006")
	liquidity := big.NewInt(1_000_000_000)
	quoteAmount := big.NewInt(500_000_000)
	agentAmount := big.NewInt(200_000_000)

	data := encodeData(t, graduatorABI, "Graduated", liquidity, quoteAmount, agentAmount)

	log := types.Log{
		Address:     common.HexToAddress("0x7777777777777777777777777777777777777777"),
		BlockNumber: 300,
		TxHash:      fixedTxHash(0x04),
		Index:       3,
		Topics: []common.Hash{
			eventID,
			addrTopic(curve),
			addrTopic(pair),
		},
		Data: data,
	}

	table := []struct {
		name      string
		duplicate bool
	}{
		{name: "first call marks graduated", duplicate: false},
		{name: "duplicate call is skipped", duplicate: true},
	}

	store := newMemStore()
	ctx := context.Background()

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			before := len(store.graduations)
			if err := handlers.HandleGraduated(ctx, log, store); err != nil {
				t.Fatalf("HandleGraduated: %v", err)
			}
			after := len(store.graduations)

			if tc.duplicate {
				if after != before {
					t.Errorf("duplicate: expected no new graduation, got %d new", after-before)
				}
				return
			}

			if after != before+1 {
				t.Fatalf("expected 1 new graduation, got %d", after-before)
			}

			got := store.graduations[before]
			if got.Curve != curve {
				t.Errorf("Curve: got %s, want %s", got.Curve, curve)
			}
			if got.Pair != pair {
				t.Errorf("Pair: got %s, want %s", got.Pair, pair)
			}
			if got.BlockNum != 300 {
				t.Errorf("BlockNum: got %d, want 300", got.BlockNum)
			}
			if got.TxHash != log.TxHash {
				t.Errorf("TxHash mismatch")
			}
			if got.LogIndex != log.Index {
				t.Errorf("LogIndex: got %d, want %d", got.LogIndex, log.Index)
			}
		})
	}
}
