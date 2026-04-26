package handlers_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	contractabi "github.com/0xDevNinja/titular/services/indexer-go/internal/abi"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/handlers"
)

// ──────────────────────────────────────────────────────────────────────────────
// In-memory ACP store (extends the base memStore via composition)
// ──────────────────────────────────────────────────────────────────────────────

type acpMemStore struct {
	*memStore
	acpAgents     []handlers.ACPAgentRegisteredRecord
	acpScores     []handlers.ACPScorePostedRecord
	acpTransfers  []handlers.ACPControllerTransferRecord
	acpJobs       []handlers.ACPJobRecord
	acpAccepted   []handlers.ACPJobAcceptedRecord
	acpResults    []handlers.ACPResultSubmittedRecord
	acpCompleted  []handlers.ACPJobCompletedRecord
	acpCancelled  []handlers.ACPJobCancelledRecord
	acpJobCreated []handlers.ACPJobCreatedRecord
	acpFunded     []handlers.ACPEscrowFundedRecord
	acpReleased   []handlers.ACPEscrowReleasedRecord
	acpRefunded   []handlers.ACPEscrowRefundedRecord
}

func newACPMemStore() *acpMemStore {
	return &acpMemStore{memStore: newMemStore()}
}

func (s *acpMemStore) UpsertACPAgent(_ context.Context, rec handlers.ACPAgentRegisteredRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpAgents = append(s.acpAgents, rec)
	return nil
}

func (s *acpMemStore) InsertACPScore(_ context.Context, rec handlers.ACPScorePostedRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpScores = append(s.acpScores, rec)
	return nil
}

func (s *acpMemStore) UpdateACPController(_ context.Context, rec handlers.ACPControllerTransferRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpTransfers = append(s.acpTransfers, rec)
	return nil
}

func (s *acpMemStore) UpsertACPJob(_ context.Context, rec handlers.ACPJobRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpJobs = append(s.acpJobs, rec)
	return nil
}

func (s *acpMemStore) UpdateACPJobAccepted(_ context.Context, rec handlers.ACPJobAcceptedRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpAccepted = append(s.acpAccepted, rec)
	return nil
}

func (s *acpMemStore) UpdateACPResultSubmitted(_ context.Context, rec handlers.ACPResultSubmittedRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpResults = append(s.acpResults, rec)
	return nil
}

func (s *acpMemStore) UpdateACPJobCompleted(_ context.Context, rec handlers.ACPJobCompletedRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpCompleted = append(s.acpCompleted, rec)
	return nil
}

func (s *acpMemStore) UpdateACPJobCancelled(_ context.Context, rec handlers.ACPJobCancelledRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpCancelled = append(s.acpCancelled, rec)
	return nil
}

func (s *acpMemStore) InsertACPJobCreated(_ context.Context, rec handlers.ACPJobCreatedRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpJobCreated = append(s.acpJobCreated, rec)
	return nil
}

func (s *acpMemStore) InsertACPEscrowFunded(_ context.Context, rec handlers.ACPEscrowFundedRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpFunded = append(s.acpFunded, rec)
	return nil
}

func (s *acpMemStore) InsertACPEscrowReleased(_ context.Context, rec handlers.ACPEscrowReleasedRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpReleased = append(s.acpReleased, rec)
	return nil
}

func (s *acpMemStore) InsertACPEscrowRefunded(_ context.Context, rec handlers.ACPEscrowRefundedRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[logKey(rec.TxHash, rec.LogIndex)] = true
	s.acpRefunded = append(s.acpRefunded, rec)
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers (ABI-level)
// ──────────────────────────────────────────────────────────────────────────────

func mustGetACPABI(t *testing.T, meta *bind.MetaData) abi.ABI {
	t.Helper()
	parsed, err := meta.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	return *parsed
}

// int256Topic converts a signed *big.Int to a log topic (two's complement, 32 bytes).
func int256Topic(n *big.Int) common.Hash {
	return common.BigToHash(n)
}

// uint64Topic converts a uint64 to a log topic.
func uint64Topic(v uint64) common.Hash {
	return common.BigToHash(new(big.Int).SetUint64(v))
}

// uint8Topic converts a uint8 to a log topic.
func uint8Topic(v uint8) common.Hash {
	return common.BigToHash(new(big.Int).SetUint64(uint64(v)))
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPAgentRegistered
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPAgentRegistered(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.AgentRegistryMetaData)
	eventID := parsed.Events["AgentRegistered"].ID

	agentID := big.NewInt(7)
	controller := common.HexToAddress("0xAAAA000000000000000000000000000000000001")
	capabilities := big.NewInt(0xFF)
	const metadataURI = "ipfs://bafybeifoo"

	data := encodeData(t, parsed, "AgentRegistered", metadataURI, capabilities)

	log := types.Log{
		Address:     common.HexToAddress("0x1111111111111111111111111111111111111111"),
		BlockNumber: 1000,
		TxHash:      fixedTxHash(0x10),
		Index:       0,
		Topics: []common.Hash{
			eventID,
			uint256Topic(agentID),
			addrTopic(controller),
		},
		Data: data,
	}

	store := newACPMemStore()
	ctx := context.Background()

	// First call — should persist.
	if err := handlers.HandleACPAgentRegistered(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpAgents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(store.acpAgents))
	}
	got := store.acpAgents[0]
	if got.AgentID.Cmp(agentID) != 0 {
		t.Errorf("AgentID: got %s, want %s", got.AgentID, agentID)
	}
	if got.Controller != controller {
		t.Errorf("Controller: got %s, want %s", got.Controller, controller)
	}
	if got.MetadataURI != metadataURI {
		t.Errorf("MetadataURI: got %q, want %q", got.MetadataURI, metadataURI)
	}
	if got.Capabilities.Cmp(capabilities) != 0 {
		t.Errorf("Capabilities: got %s, want %s", got.Capabilities, capabilities)
	}
	if got.BlockNum != 1000 {
		t.Errorf("BlockNum: got %d, want 1000", got.BlockNum)
	}

	// Second call — duplicate, must be skipped.
	if err := handlers.HandleACPAgentRegistered(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpAgents) != 1 {
		t.Errorf("duplicate: expected 1 agent, got %d", len(store.acpAgents))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPScorePosted
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPScorePosted(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.AgentRegistryMetaData)
	eventID := parsed.Events["ScorePosted"].ID

	agentID := big.NewInt(3)
	scorer := common.HexToAddress("0xBBBB000000000000000000000000000000000002")
	delta := big.NewInt(-10)
	newTotal := big.NewInt(90)
	nonce := big.NewInt(1)

	data := encodeData(t, parsed, "ScorePosted", delta, newTotal, nonce)

	log := types.Log{
		Address:     common.HexToAddress("0x2222222222222222222222222222222222222222"),
		BlockNumber: 2000,
		TxHash:      fixedTxHash(0x11),
		Index:       1,
		Topics: []common.Hash{
			eventID,
			uint256Topic(agentID),
			addrTopic(scorer),
		},
		Data: data,
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPScorePosted(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpScores) != 1 {
		t.Fatalf("expected 1 score, got %d", len(store.acpScores))
	}
	got := store.acpScores[0]
	if got.AgentID.Cmp(agentID) != 0 {
		t.Errorf("AgentID: got %s, want %s", got.AgentID, agentID)
	}
	if got.Scorer != scorer {
		t.Errorf("Scorer: got %s, want %s", got.Scorer, scorer)
	}
	if got.Delta.Cmp(delta) != 0 {
		t.Errorf("Delta: got %s, want %s", got.Delta, delta)
	}
	if got.NewTotal.Cmp(newTotal) != 0 {
		t.Errorf("NewTotal: got %s, want %s", got.NewTotal, newTotal)
	}
	if got.Nonce.Cmp(nonce) != 0 {
		t.Errorf("Nonce: got %s, want %s", got.Nonce, nonce)
	}

	// Duplicate is skipped.
	if err := handlers.HandleACPScorePosted(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpScores) != 1 {
		t.Errorf("duplicate: expected 1 score, got %d", len(store.acpScores))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPControllerTransferAccepted
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPControllerTransferAccepted(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.AgentRegistryMetaData)
	eventID := parsed.Events["ControllerTransferAccepted"].ID

	agentID := big.NewInt(5)
	oldCtrl := common.HexToAddress("0xCCCC000000000000000000000000000000000003")
	newCtrl := common.HexToAddress("0xDDDD000000000000000000000000000000000004")

	// ControllerTransferAccepted has all three fields indexed — no non-indexed data.
	log := types.Log{
		Address:     common.HexToAddress("0x3333333333333333333333333333333333333333"),
		BlockNumber: 3000,
		TxHash:      fixedTxHash(0x12),
		Index:       2,
		Topics: []common.Hash{
			eventID,
			uint256Topic(agentID),
			addrTopic(oldCtrl),
			addrTopic(newCtrl),
		},
		Data: []byte{},
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPControllerTransferAccepted(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpTransfers) != 1 {
		t.Fatalf("expected 1 transfer, got %d", len(store.acpTransfers))
	}
	got := store.acpTransfers[0]
	if got.AgentID.Cmp(agentID) != 0 {
		t.Errorf("AgentID: got %s, want %s", got.AgentID, agentID)
	}
	if got.OldController != oldCtrl {
		t.Errorf("OldController: got %s, want %s", got.OldController, oldCtrl)
	}
	if got.NewController != newCtrl {
		t.Errorf("NewController: got %s, want %s", got.NewController, newCtrl)
	}

	// Duplicate skipped.
	if err := handlers.HandleACPControllerTransferAccepted(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpTransfers) != 1 {
		t.Errorf("duplicate: expected 1 transfer, got %d", len(store.acpTransfers))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPJobInitialised
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPJobInitialised(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.JobMetaData)
	eventID := parsed.Events["JobInitialised"].ID

	jobID := big.NewInt(42)
	principal := common.HexToAddress("0xEEEE000000000000000000000000000000000005")
	agentID := big.NewInt(0) // open to any
	var jobType uint8 = 0    // Direct
	budget := big.NewInt(1_000_000)
	token := common.HexToAddress("0xFFFF000000000000000000000000000000000006")
	var deadline uint64 = 9999999999

	data := encodeData(t, parsed, "JobInitialised", jobType, budget, token, deadline)

	log := types.Log{
		Address:     common.HexToAddress("0x4444444444444444444444444444444444444444"),
		BlockNumber: 4000,
		TxHash:      fixedTxHash(0x13),
		Index:       3,
		Topics: []common.Hash{
			eventID,
			uint256Topic(jobID),
			addrTopic(principal),
			uint256Topic(agentID),
		},
		Data: data,
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPJobInitialised(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpJobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(store.acpJobs))
	}
	got := store.acpJobs[0]
	if got.JobID.Cmp(jobID) != 0 {
		t.Errorf("JobID: got %s, want %s", got.JobID, jobID)
	}
	if got.Principal != principal {
		t.Errorf("Principal: got %s, want %s", got.Principal, principal)
	}
	if got.Budget.Cmp(budget) != 0 {
		t.Errorf("Budget: got %s, want %s", got.Budget, budget)
	}
	if got.Token != token {
		t.Errorf("Token: got %s, want %s", got.Token, token)
	}
	if got.JobType != jobType {
		t.Errorf("JobType: got %d, want %d", got.JobType, jobType)
	}
	if got.Deadline != deadline {
		t.Errorf("Deadline: got %d, want %d", got.Deadline, deadline)
	}

	// Duplicate skipped.
	if err := handlers.HandleACPJobInitialised(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpJobs) != 1 {
		t.Errorf("duplicate: expected 1 job, got %d", len(store.acpJobs))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPAgentAccepted
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPAgentAccepted(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.JobMetaData)
	eventID := parsed.Events["AgentAccepted"].ID

	jobID := big.NewInt(42)
	agent := common.HexToAddress("0xAAAA100000000000000000000000000000000007")
	agentID := big.NewInt(9)

	// All three fields are indexed — no non-indexed data.
	log := types.Log{
		Address:     common.HexToAddress("0x5555555555555555555555555555555555555555"),
		BlockNumber: 5000,
		TxHash:      fixedTxHash(0x14),
		Index:       4,
		Topics: []common.Hash{
			eventID,
			uint256Topic(jobID),
			addrTopic(agent),
			uint256Topic(agentID),
		},
		Data: []byte{},
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPAgentAccepted(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpAccepted) != 1 {
		t.Fatalf("expected 1 accepted record, got %d", len(store.acpAccepted))
	}
	got := store.acpAccepted[0]
	if got.JobID.Cmp(jobID) != 0 {
		t.Errorf("JobID: got %s, want %s", got.JobID, jobID)
	}
	if got.Agent != agent {
		t.Errorf("Agent: got %s, want %s", got.Agent, agent)
	}
	if got.AgentID.Cmp(agentID) != 0 {
		t.Errorf("AgentID: got %s, want %s", got.AgentID, agentID)
	}

	// Duplicate skipped.
	if err := handlers.HandleACPAgentAccepted(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpAccepted) != 1 {
		t.Errorf("duplicate: expected 1 accepted record, got %d", len(store.acpAccepted))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPResultSubmitted
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPResultSubmitted(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.JobMetaData)
	eventID := parsed.Events["ResultSubmitted"].ID

	jobID := big.NewInt(42)
	agent := common.HexToAddress("0xAAAA200000000000000000000000000000000008")
	const resultURI = "ipfs://bafybeijobresult"

	data := encodeData(t, parsed, "ResultSubmitted", resultURI)

	log := types.Log{
		Address:     common.HexToAddress("0x6666666666666666666666666666666666666666"),
		BlockNumber: 6000,
		TxHash:      fixedTxHash(0x15),
		Index:       5,
		Topics: []common.Hash{
			eventID,
			uint256Topic(jobID),
			addrTopic(agent),
		},
		Data: data,
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPResultSubmitted(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpResults) != 1 {
		t.Fatalf("expected 1 result, got %d", len(store.acpResults))
	}
	got := store.acpResults[0]
	if got.JobID.Cmp(jobID) != 0 {
		t.Errorf("JobID: got %s, want %s", got.JobID, jobID)
	}
	if got.Agent != agent {
		t.Errorf("Agent: got %s, want %s", got.Agent, agent)
	}
	if got.ResultURI != resultURI {
		t.Errorf("ResultURI: got %q, want %q", got.ResultURI, resultURI)
	}

	// Duplicate skipped.
	if err := handlers.HandleACPResultSubmitted(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpResults) != 1 {
		t.Errorf("duplicate: expected 1 result, got %d", len(store.acpResults))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPJobCompleted
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPJobCompleted(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.JobMetaData)
	eventID := parsed.Events["JobCompleted"].ID

	jobID := big.NewInt(42)
	releasedBy := common.HexToAddress("0xAAAA300000000000000000000000000000000009")
	amount := big.NewInt(1_000_000)

	data := encodeData(t, parsed, "JobCompleted", amount)

	log := types.Log{
		Address:     common.HexToAddress("0x7777777777777777777777777777777777777777"),
		BlockNumber: 7000,
		TxHash:      fixedTxHash(0x16),
		Index:       6,
		Topics: []common.Hash{
			eventID,
			uint256Topic(jobID),
			addrTopic(releasedBy),
		},
		Data: data,
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPJobCompleted(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpCompleted) != 1 {
		t.Fatalf("expected 1 completed record, got %d", len(store.acpCompleted))
	}
	got := store.acpCompleted[0]
	if got.JobID.Cmp(jobID) != 0 {
		t.Errorf("JobID: got %s, want %s", got.JobID, jobID)
	}
	if got.ReleasedBy != releasedBy {
		t.Errorf("ReleasedBy: got %s, want %s", got.ReleasedBy, releasedBy)
	}
	if got.Amount.Cmp(amount) != 0 {
		t.Errorf("Amount: got %s, want %s", got.Amount, amount)
	}

	// Duplicate skipped.
	if err := handlers.HandleACPJobCompleted(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpCompleted) != 1 {
		t.Errorf("duplicate: expected 1 completed record, got %d", len(store.acpCompleted))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPJobCancelled
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPJobCancelled(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.JobMetaData)
	eventID := parsed.Events["JobCancelled"].ID

	jobID := big.NewInt(42)
	cancelledBy := common.HexToAddress("0xAAAA400000000000000000000000000000000010")
	const reason = "principal cancelled"

	data := encodeData(t, parsed, "JobCancelled", reason)

	log := types.Log{
		Address:     common.HexToAddress("0x8888888888888888888888888888888888888888"),
		BlockNumber: 8000,
		TxHash:      fixedTxHash(0x17),
		Index:       7,
		Topics: []common.Hash{
			eventID,
			uint256Topic(jobID),
			addrTopic(cancelledBy),
		},
		Data: data,
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPJobCancelled(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpCancelled) != 1 {
		t.Fatalf("expected 1 cancelled record, got %d", len(store.acpCancelled))
	}
	got := store.acpCancelled[0]
	if got.JobID.Cmp(jobID) != 0 {
		t.Errorf("JobID: got %s, want %s", got.JobID, jobID)
	}
	if got.CancelledBy != cancelledBy {
		t.Errorf("CancelledBy: got %s, want %s", got.CancelledBy, cancelledBy)
	}
	if got.Reason != reason {
		t.Errorf("Reason: got %q, want %q", got.Reason, reason)
	}

	// Duplicate skipped.
	if err := handlers.HandleACPJobCancelled(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpCancelled) != 1 {
		t.Errorf("duplicate: expected 1 cancelled record, got %d", len(store.acpCancelled))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPJobCreated
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPJobCreated(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.JobFactoryMetaData)
	eventID := parsed.Events["JobCreated"].ID

	jobID := big.NewInt(1)
	clone := common.HexToAddress("0xCCCC100000000000000000000000000000000011")
	principal := common.HexToAddress("0xDDDD100000000000000000000000000000000012")
	var jobType uint8 = 1 // Evaluated
	token := common.HexToAddress("0xEEEE100000000000000000000000000000000013")
	budget := big.NewInt(500_000)

	data := encodeData(t, parsed, "JobCreated", jobType, token, budget)

	log := types.Log{
		Address:     common.HexToAddress("0x9999999999999999999999999999999999999999"),
		BlockNumber: 9000,
		TxHash:      fixedTxHash(0x18),
		Index:       8,
		Topics: []common.Hash{
			eventID,
			uint256Topic(jobID),
			addrTopic(clone),
			addrTopic(principal),
		},
		Data: data,
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPJobCreated(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpJobCreated) != 1 {
		t.Fatalf("expected 1 job-created record, got %d", len(store.acpJobCreated))
	}
	got := store.acpJobCreated[0]
	if got.JobID.Cmp(jobID) != 0 {
		t.Errorf("JobID: got %s, want %s", got.JobID, jobID)
	}
	if got.Clone != clone {
		t.Errorf("Clone: got %s, want %s", got.Clone, clone)
	}
	if got.Principal != principal {
		t.Errorf("Principal: got %s, want %s", got.Principal, principal)
	}
	if got.JobType != jobType {
		t.Errorf("JobType: got %d, want %d", got.JobType, jobType)
	}
	if got.Token != token {
		t.Errorf("Token: got %s, want %s", got.Token, token)
	}
	if got.Budget.Cmp(budget) != 0 {
		t.Errorf("Budget: got %s, want %s", got.Budget, budget)
	}

	// Duplicate skipped.
	if err := handlers.HandleACPJobCreated(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpJobCreated) != 1 {
		t.Errorf("duplicate: expected 1 job-created record, got %d", len(store.acpJobCreated))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPEscrowFunded
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPEscrowFunded(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.EscrowMetaData)
	eventID := parsed.Events["Funded"].ID

	depositor := common.HexToAddress("0xFFFF100000000000000000000000000000000014")
	jobID := big.NewInt(1)
	token := common.HexToAddress("0xFFFF200000000000000000000000000000000015")
	amount := big.NewInt(250_000)

	data := encodeData(t, parsed, "Funded", amount)

	log := types.Log{
		Address:     common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		BlockNumber: 10000,
		TxHash:      fixedTxHash(0x19),
		Index:       9,
		Topics: []common.Hash{
			eventID,
			addrTopic(depositor),
			uint256Topic(jobID),
			addrTopic(token),
		},
		Data: data,
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPEscrowFunded(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpFunded) != 1 {
		t.Fatalf("expected 1 funded record, got %d", len(store.acpFunded))
	}
	got := store.acpFunded[0]
	if got.Depositor != depositor {
		t.Errorf("Depositor: got %s, want %s", got.Depositor, depositor)
	}
	if got.JobID.Cmp(jobID) != 0 {
		t.Errorf("JobID: got %s, want %s", got.JobID, jobID)
	}
	if got.Token != token {
		t.Errorf("Token: got %s, want %s", got.Token, token)
	}
	if got.Amount.Cmp(amount) != 0 {
		t.Errorf("Amount: got %s, want %s", got.Amount, amount)
	}

	// Duplicate skipped.
	if err := handlers.HandleACPEscrowFunded(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpFunded) != 1 {
		t.Errorf("duplicate: expected 1 funded record, got %d", len(store.acpFunded))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPEscrowReleased
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPEscrowReleased(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.EscrowMetaData)
	eventID := parsed.Events["Released"].ID

	depositor := common.HexToAddress("0xFFFF300000000000000000000000000000000016")
	jobID := big.NewInt(1)
	token := common.HexToAddress("0xFFFF400000000000000000000000000000000017")
	to := common.HexToAddress("0xFFFF500000000000000000000000000000000018")
	amount := big.NewInt(250_000)

	data := encodeData(t, parsed, "Released", to, amount)

	log := types.Log{
		Address:     common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
		BlockNumber: 11000,
		TxHash:      fixedTxHash(0x1A),
		Index:       10,
		Topics: []common.Hash{
			eventID,
			addrTopic(depositor),
			uint256Topic(jobID),
			addrTopic(token),
		},
		Data: data,
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPEscrowReleased(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpReleased) != 1 {
		t.Fatalf("expected 1 released record, got %d", len(store.acpReleased))
	}
	got := store.acpReleased[0]
	if got.Depositor != depositor {
		t.Errorf("Depositor: got %s, want %s", got.Depositor, depositor)
	}
	if got.To != to {
		t.Errorf("To: got %s, want %s", got.To, to)
	}
	if got.Amount.Cmp(amount) != 0 {
		t.Errorf("Amount: got %s, want %s", got.Amount, amount)
	}

	// Duplicate skipped.
	if err := handlers.HandleACPEscrowReleased(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpReleased) != 1 {
		t.Errorf("duplicate: expected 1 released record, got %d", len(store.acpReleased))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestHandleACPEscrowRefunded
// ──────────────────────────────────────────────────────────────────────────────

func TestHandleACPEscrowRefunded(t *testing.T) {
	parsed := mustGetACPABI(t, contractabi.EscrowMetaData)
	eventID := parsed.Events["Refunded"].ID

	depositor := common.HexToAddress("0xFFFF600000000000000000000000000000000019")
	jobID := big.NewInt(2)
	token := common.HexToAddress("0xFFFF700000000000000000000000000000000020")
	amount := big.NewInt(100_000)

	data := encodeData(t, parsed, "Refunded", amount)

	log := types.Log{
		Address:     common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"),
		BlockNumber: 12000,
		TxHash:      fixedTxHash(0x1B),
		Index:       11,
		Topics: []common.Hash{
			eventID,
			addrTopic(depositor),
			uint256Topic(jobID),
			addrTopic(token),
		},
		Data: data,
	}

	store := newACPMemStore()
	ctx := context.Background()

	if err := handlers.HandleACPEscrowRefunded(ctx, log, store); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if len(store.acpRefunded) != 1 {
		t.Fatalf("expected 1 refunded record, got %d", len(store.acpRefunded))
	}
	got := store.acpRefunded[0]
	if got.Depositor != depositor {
		t.Errorf("Depositor: got %s, want %s", got.Depositor, depositor)
	}
	if got.JobID.Cmp(jobID) != 0 {
		t.Errorf("JobID: got %s, want %s", got.JobID, jobID)
	}
	if got.Token != token {
		t.Errorf("Token: got %s, want %s", got.Token, token)
	}
	if got.Amount.Cmp(amount) != 0 {
		t.Errorf("Amount: got %s, want %s", got.Amount, amount)
	}

	// Duplicate skipped.
	if err := handlers.HandleACPEscrowRefunded(ctx, log, store); err != nil {
		t.Fatalf("duplicate call: %v", err)
	}
	if len(store.acpRefunded) != 1 {
		t.Errorf("duplicate: expected 1 refunded record, got %d", len(store.acpRefunded))
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestACPHandlers_StoreError
// Verify that an error from IsLogProcessed is propagated correctly.
// ──────────────────────────────────────────────────────────────────────────────

func TestACPHandlers_StoreError(t *testing.T) {
	ctx := context.Background()
	es := &errStore{}
	emptyLog := types.Log{TxHash: fixedTxHash(0xEE), Index: 77}

	cases := []struct {
		name    string
		handler func(context.Context, types.Log, handlers.Store) error
	}{
		{"ACPAgentRegistered", handlers.HandleACPAgentRegistered},
		{"ACPScorePosted", handlers.HandleACPScorePosted},
		{"ACPControllerTransferAccepted", handlers.HandleACPControllerTransferAccepted},
		{"ACPJobInitialised", handlers.HandleACPJobInitialised},
		{"ACPAgentAccepted", handlers.HandleACPAgentAccepted},
		{"ACPResultSubmitted", handlers.HandleACPResultSubmitted},
		{"ACPJobCompleted", handlers.HandleACPJobCompleted},
		{"ACPJobCancelled", handlers.HandleACPJobCancelled},
		{"ACPJobCreated", handlers.HandleACPJobCreated},
		{"ACPEscrowFunded", handlers.HandleACPEscrowFunded},
		{"ACPEscrowReleased", handlers.HandleACPEscrowReleased},
		{"ACPEscrowRefunded", handlers.HandleACPEscrowRefunded},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.handler(ctx, emptyLog, es); err == nil {
				t.Errorf("%s: expected error from store, got nil", tc.name)
			}
		})
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// TestACPHandlers_ParseError
// Verify that a malformed log causes a parse error, not a panic.
// ──────────────────────────────────────────────────────────────────────────────

func TestACPHandlers_ParseError(t *testing.T) {
	ctx := context.Background()
	badLog := types.Log{
		TxHash: fixedTxHash(0xFF),
		Index:  66,
		Topics: []common.Hash{},
		Data:   []byte{0xDE, 0xAD},
	}

	cases := []struct {
		name    string
		handler func(context.Context, types.Log, handlers.Store) error
	}{
		{"ACPAgentRegistered", handlers.HandleACPAgentRegistered},
		{"ACPScorePosted", handlers.HandleACPScorePosted},
		{"ACPControllerTransferAccepted", handlers.HandleACPControllerTransferAccepted},
		{"ACPJobInitialised", handlers.HandleACPJobInitialised},
		{"ACPAgentAccepted", handlers.HandleACPAgentAccepted},
		{"ACPResultSubmitted", handlers.HandleACPResultSubmitted},
		{"ACPJobCompleted", handlers.HandleACPJobCompleted},
		{"ACPJobCancelled", handlers.HandleACPJobCancelled},
		{"ACPJobCreated", handlers.HandleACPJobCreated},
		{"ACPEscrowFunded", handlers.HandleACPEscrowFunded},
		{"ACPEscrowReleased", handlers.HandleACPEscrowReleased},
		{"ACPEscrowRefunded", handlers.HandleACPEscrowRefunded},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			store := newACPMemStore()
			err := tc.handler(ctx, badLog, store)
			if err == nil {
				t.Errorf("%s: expected parse error, got nil", tc.name)
			}
		})
	}
}
