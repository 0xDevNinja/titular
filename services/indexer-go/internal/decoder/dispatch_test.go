package decoder_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	contractabi "github.com/0xDevNinja/titular/services/indexer-go/internal/abi"
	"github.com/0xDevNinja/titular/services/indexer-go/internal/decoder"
)

// Tests in this file build synthetic types.Log values from the embedded
// abigen ABI and assert that:
//
//   - decoder.Dispatch resolves topic0 -> the right contract-qualified name.
//   - the returned Payload is the abigen-generated struct type for that event.
//   - non-indexed fields round-trip (we sample one or two per case to keep the
//     suite small while still proving the decoder did real work).
//
// The fixture builders (`encodeData`, the topic helpers) are local to keep
// the decoder package self-contained; the equivalent helpers in
// internal/handlers/handlers_test.go are not exported.

// ──────────────────────────────────────────────────────────────────────────
// Fixture helpers
// ──────────────────────────────────────────────────────────────────────────

func mustAbi(t *testing.T, m *bind.MetaData) abi.ABI {
	t.Helper()
	parsed, err := m.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	return *parsed
}

// encodeData packs the non-indexed fields of an event in declaration order.
// Caller must pass the correctly-typed Go values (uint8, *big.Int, bool, etc.).
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

// addrTopic left-pads an address to 32 bytes for use as a topic.
func addrTopic(a common.Address) common.Hash { return common.BytesToHash(a.Bytes()) }

// uintTopic converts a *big.Int to a 32-byte topic (used for both uint256
// and signed integers — for negatives the caller must pre-encode two's
// complement, which big.Int's BigToHash does for non-negative values).
func uintTopic(n *big.Int) common.Hash { return common.BigToHash(n) }

// uint8Topic encodes a uint8 as a 32-byte topic (left-padded big-endian).
func uint8Topic(v uint8) common.Hash {
	return common.BigToHash(new(big.Int).SetUint64(uint64(v)))
}

// fixedTxHash returns a deterministic tx hash for a given seed.
func fixedTxHash(seed byte) common.Hash {
	var h common.Hash
	h[31] = seed
	return h
}

// topicID returns keccak256 topic0 for an event in the given ABI.
func topicID(t *testing.T, m *bind.MetaData, eventName string) common.Hash {
	t.Helper()
	parsed := mustAbi(t, m)
	ev, ok := parsed.Events[eventName]
	if !ok {
		t.Fatalf("event %q not in ABI", eventName)
	}
	return ev.ID
}

// ──────────────────────────────────────────────────────────────────────────
// Dispatch table coverage — every event pinned in abi_test.go must round-trip.
//
// Each case constructs a synthetic log, runs it through Dispatch, asserts
// the qualified name, and type-asserts the payload back to its concrete
// abigen struct. A targeted field check on the payload proves the decoder
// actually consumed the topics/data instead of returning a zero value.
// ──────────────────────────────────────────────────────────────────────────

// Test_Dispatch_AgentRegistry covers all 8 AgentRegistry events.
func Test_Dispatch_AgentRegistry(t *testing.T) {
	parsed := mustAbi(t, contractabi.AgentRegistryMetaData)
	registry := contractabi.AgentRegistryMetaData

	t.Run("AgentRegistered", func(t *testing.T) {
		agentID := big.NewInt(7)
		controller := common.HexToAddress("0xAAAA000000000000000000000000000000000001")
		caps := big.NewInt(0x1F)
		uri := "ipfs://bafybeifoo"

		log := types.Log{
			Topics: []common.Hash{
				topicID(t, registry, "AgentRegistered"),
				uintTopic(agentID),
				addrTopic(controller),
			},
			Data:        encodeData(t, parsed, "AgentRegistered", uri, caps),
			BlockNumber: 100, TxHash: fixedTxHash(1), Index: 0,
		}

		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if ev.Name != "AgentRegistry.AgentRegistered" {
			t.Errorf("Name = %q", ev.Name)
		}
		got, ok := ev.Payload.(*contractabi.AgentRegistryAgentRegistered)
		if !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
		if got.AgentId.Cmp(agentID) != 0 {
			t.Errorf("AgentId = %s, want %s", got.AgentId, agentID)
		}
		if got.Controller != controller {
			t.Errorf("Controller = %s, want %s", got.Controller, controller)
		}
		if got.MetadataURI != uri {
			t.Errorf("MetadataURI = %q, want %q", got.MetadataURI, uri)
		}
	})

	t.Run("MetadataUpdated", func(t *testing.T) {
		agentID := big.NewInt(2)
		uri := "ipfs://meta-v2"
		log := types.Log{
			Topics: []common.Hash{topicID(t, registry, "MetadataUpdated"), uintTopic(agentID)},
			Data:   encodeData(t, parsed, "MetadataUpdated", uri),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if ev.Name != "AgentRegistry.MetadataUpdated" {
			t.Errorf("Name = %q", ev.Name)
		}
		got := ev.Payload.(*contractabi.AgentRegistryMetadataUpdated)
		if got.MetadataURI != uri {
			t.Errorf("MetadataURI = %q", got.MetadataURI)
		}
	})

	t.Run("CapabilitiesUpdated", func(t *testing.T) {
		caps := big.NewInt(0xFF)
		log := types.Log{
			Topics: []common.Hash{topicID(t, registry, "CapabilitiesUpdated"), uintTopic(big.NewInt(3))},
			Data:   encodeData(t, parsed, "CapabilitiesUpdated", caps),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.AgentRegistryCapabilitiesUpdated)
		if got.Capabilities.Cmp(caps) != 0 {
			t.Errorf("Capabilities = %s", got.Capabilities)
		}
	})

	t.Run("ActiveStatusChanged", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{topicID(t, registry, "ActiveStatusChanged"), uintTopic(big.NewInt(4))},
			Data:   encodeData(t, parsed, "ActiveStatusChanged", true),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.AgentRegistryActiveStatusChanged)
		if !got.Active {
			t.Errorf("Active = false, want true")
		}
	})

	t.Run("ScorePosted", func(t *testing.T) {
		agentID := big.NewInt(8)
		scorer := common.HexToAddress("0xBBBB000000000000000000000000000000000002")
		delta := big.NewInt(15)
		newTotal := big.NewInt(115)
		nonce := big.NewInt(1)
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, registry, "ScorePosted"),
				uintTopic(agentID),
				addrTopic(scorer),
			},
			Data: encodeData(t, parsed, "ScorePosted", delta, newTotal, nonce),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.AgentRegistryScorePosted)
		if got.Delta.Cmp(delta) != 0 || got.NewTotal.Cmp(newTotal) != 0 {
			t.Errorf("score fields mismatch: delta=%s, newTotal=%s", got.Delta, got.NewTotal)
		}
	})

	t.Run("ControllerTransferProposed", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, registry, "ControllerTransferProposed"),
				uintTopic(big.NewInt(5)),
				addrTopic(common.HexToAddress("0x1111111111111111111111111111111111111111")),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.AgentRegistryControllerTransferProposed); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})

	t.Run("ControllerTransferAccepted", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, registry, "ControllerTransferAccepted"),
				uintTopic(big.NewInt(6)),
				addrTopic(common.HexToAddress("0x2222222222222222222222222222222222222222")),
				addrTopic(common.HexToAddress("0x3333333333333333333333333333333333333333")),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.AgentRegistryControllerTransferAccepted); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})

	t.Run("ControllerTransferCancelled", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, registry, "ControllerTransferCancelled"),
				uintTopic(big.NewInt(7)),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.AgentRegistryControllerTransferCancelled); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})
}

// Test_Dispatch_Job covers all 10 Job events.
func Test_Dispatch_Job(t *testing.T) {
	parsed := mustAbi(t, contractabi.JobMetaData)
	jm := contractabi.JobMetaData

	t.Run("JobInitialised", func(t *testing.T) {
		jobID := big.NewInt(42)
		principal := common.HexToAddress("0xCC0C000000000000000000000000000000000001")
		agentID := big.NewInt(0)
		var jobType uint8 = 1
		budget := big.NewInt(1_000_000)
		token := common.HexToAddress("0xDD0D000000000000000000000000000000000002")
		var deadline uint64 = 9999999999

		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jm, "JobInitialised"),
				uintTopic(jobID),
				addrTopic(principal),
				uintTopic(agentID),
			},
			Data: encodeData(t, parsed, "JobInitialised", jobType, budget, token, deadline),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if ev.Name != "Job.JobInitialised" {
			t.Errorf("Name = %q", ev.Name)
		}
		got := ev.Payload.(*contractabi.JobJobInitialised)
		if got.JobId.Cmp(jobID) != 0 || got.JobType != jobType || got.Deadline != deadline {
			t.Errorf("fields: jobId=%s jobType=%d deadline=%d", got.JobId, got.JobType, got.Deadline)
		}
	})

	t.Run("AgentAccepted", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jm, "AgentAccepted"),
				uintTopic(big.NewInt(11)),
				addrTopic(common.HexToAddress("0xEEEE000000000000000000000000000000000003")),
				uintTopic(big.NewInt(7)),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.JobAgentAccepted); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})

	t.Run("ResultSubmitted", func(t *testing.T) {
		uri := "ipfs://result-1"
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jm, "ResultSubmitted"),
				uintTopic(big.NewInt(12)),
				addrTopic(common.HexToAddress("0xFEED000000000000000000000000000000000004")),
			},
			Data: encodeData(t, parsed, "ResultSubmitted", uri),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.JobResultSubmitted)
		if got.ResultURI != uri {
			t.Errorf("ResultURI = %q", got.ResultURI)
		}
	})

	t.Run("ResultApproved", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jm, "ResultApproved"),
				uintTopic(big.NewInt(13)),
				addrTopic(common.HexToAddress("0xAAAA000000000000000000000000000000000005")),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.JobResultApproved); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})

	t.Run("ResultRejected", func(t *testing.T) {
		reason := "did not match spec"
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jm, "ResultRejected"),
				uintTopic(big.NewInt(14)),
				addrTopic(common.HexToAddress("0xBBBB000000000000000000000000000000000006")),
			},
			Data: encodeData(t, parsed, "ResultRejected", reason),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.JobResultRejected)
		if got.Reason != reason {
			t.Errorf("Reason = %q", got.Reason)
		}
	})

	t.Run("JobCompleted", func(t *testing.T) {
		amount := big.NewInt(900_000)
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jm, "JobCompleted"),
				uintTopic(big.NewInt(15)),
				addrTopic(common.HexToAddress("0xCCCC000000000000000000000000000000000007")),
			},
			Data: encodeData(t, parsed, "JobCompleted", amount),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.JobJobCompleted)
		if got.Amount.Cmp(amount) != 0 {
			t.Errorf("Amount = %s", got.Amount)
		}
	})

	t.Run("JobCancelled", func(t *testing.T) {
		reason := "buyer cancelled"
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jm, "JobCancelled"),
				uintTopic(big.NewInt(16)),
				addrTopic(common.HexToAddress("0xDDDD000000000000000000000000000000000008")),
			},
			Data: encodeData(t, parsed, "JobCancelled", reason),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.JobJobCancelled)
		if got.Reason != reason {
			t.Errorf("Reason = %q", got.Reason)
		}
	})

	t.Run("DisputeRaised", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jm, "DisputeRaised"),
				uintTopic(big.NewInt(17)),
				addrTopic(common.HexToAddress("0xEEEE000000000000000000000000000000000009")),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.JobDisputeRaised); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})

	t.Run("DisputeResolved", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jm, "DisputeResolved"),
				uintTopic(big.NewInt(18)),
				addrTopic(common.HexToAddress("0xFEEE00000000000000000000000000000000000A")),
			},
			Data: encodeData(t, parsed, "DisputeResolved", true),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.JobDisputeResolved)
		if !got.AgentFavoured {
			t.Errorf("AgentFavoured = false, want true")
		}
	})

	t.Run("EvaluatorAssigned", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jm, "EvaluatorAssigned"),
				uintTopic(big.NewInt(19)),
				addrTopic(common.HexToAddress("0xAAAA00000000000000000000000000000000000B")),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.JobEvaluatorAssigned); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})
}

// Test_Dispatch_JobFactory covers all 3 JobFactory events.
func Test_Dispatch_JobFactory(t *testing.T) {
	parsed := mustAbi(t, contractabi.JobFactoryMetaData)
	jfm := contractabi.JobFactoryMetaData

	t.Run("JobCreated", func(t *testing.T) {
		jobID := big.NewInt(99)
		clone := common.HexToAddress("0x0001000000000000000000000000000000000001")
		principal := common.HexToAddress("0x0002000000000000000000000000000000000002")
		var jobType uint8 = 0
		token := common.HexToAddress("0x0003000000000000000000000000000000000003")
		budget := big.NewInt(2_500_000)
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jfm, "JobCreated"),
				uintTopic(jobID),
				addrTopic(clone),
				addrTopic(principal),
			},
			Data: encodeData(t, parsed, "JobCreated", jobType, token, budget),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if ev.Name != "JobFactory.JobCreated" {
			t.Errorf("Name = %q", ev.Name)
		}
		got := ev.Payload.(*contractabi.JobFactoryJobCreated)
		if got.Clone != clone || got.Budget.Cmp(budget) != 0 {
			t.Errorf("fields: clone=%s budget=%s", got.Clone, got.Budget)
		}
	})

	t.Run("ImplementationUpdated", func(t *testing.T) {
		oldImpl := common.HexToAddress("0x0004000000000000000000000000000000000004")
		newImpl := common.HexToAddress("0x0005000000000000000000000000000000000005")
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jfm, "ImplementationUpdated"),
				addrTopic(oldImpl),
				addrTopic(newImpl),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.JobFactoryImplementationUpdated); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})

	t.Run("DefaultArbiterUpdated", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, jfm, "DefaultArbiterUpdated"),
				addrTopic(common.HexToAddress("0x0006000000000000000000000000000000000006")),
				addrTopic(common.HexToAddress("0x0007000000000000000000000000000000000007")),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.JobFactoryDefaultArbiterUpdated); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})
}

// Test_Dispatch_Escrow covers all 3 Escrow events.
func Test_Dispatch_Escrow(t *testing.T) {
	parsed := mustAbi(t, contractabi.EscrowMetaData)
	em := contractabi.EscrowMetaData

	t.Run("Funded", func(t *testing.T) {
		depositor := common.HexToAddress("0x1100000000000000000000000000000000000001")
		jobID := big.NewInt(20)
		token := common.HexToAddress("0x2200000000000000000000000000000000000002")
		amount := big.NewInt(500_000)
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, em, "Funded"),
				addrTopic(depositor),
				uintTopic(jobID),
				addrTopic(token),
			},
			Data: encodeData(t, parsed, "Funded", amount),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if ev.Name != "Escrow.Funded" {
			t.Errorf("Name = %q", ev.Name)
		}
		got := ev.Payload.(*contractabi.EscrowFunded)
		if got.Amount.Cmp(amount) != 0 || got.Depositor != depositor {
			t.Errorf("fields: depositor=%s amount=%s", got.Depositor, got.Amount)
		}
	})

	t.Run("Released", func(t *testing.T) {
		depositor := common.HexToAddress("0x1100000000000000000000000000000000000010")
		jobID := big.NewInt(21)
		token := common.HexToAddress("0x2200000000000000000000000000000000000020")
		to := common.HexToAddress("0x3300000000000000000000000000000000000030")
		amount := big.NewInt(123)
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, em, "Released"),
				addrTopic(depositor),
				uintTopic(jobID),
				addrTopic(token),
			},
			Data: encodeData(t, parsed, "Released", to, amount),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.EscrowReleased)
		if got.To != to || got.Amount.Cmp(amount) != 0 {
			t.Errorf("fields: to=%s amount=%s", got.To, got.Amount)
		}
	})

	t.Run("Refunded", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, em, "Refunded"),
				addrTopic(common.HexToAddress("0x4400000000000000000000000000000000000040")),
				uintTopic(big.NewInt(22)),
				addrTopic(common.HexToAddress("0x5500000000000000000000000000000000000050")),
			},
			Data: encodeData(t, parsed, "Refunded", big.NewInt(77)),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.EscrowRefunded)
		if got.Amount.Cmp(big.NewInt(77)) != 0 {
			t.Errorf("Amount = %s", got.Amount)
		}
	})
}

// Test_Dispatch_FeeSplitter covers all 3 FeeSplitter events.
func Test_Dispatch_FeeSplitter(t *testing.T) {
	parsed := mustAbi(t, contractabi.FeeSplitterMetaData)
	fm := contractabi.FeeSplitterMetaData

	t.Run("FeeSplit", func(t *testing.T) {
		token := common.HexToAddress("0x6600000000000000000000000000000000000060")
		primary := common.HexToAddress("0x7700000000000000000000000000000000000070")
		var schedule uint8 = 2
		primaryAmt := big.NewInt(900)
		treasuryAmt := big.NewInt(50)
		buybackAmt := big.NewInt(50)
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, fm, "FeeSplit"),
				addrTopic(token),
				addrTopic(primary),
				uint8Topic(schedule),
			},
			Data: encodeData(t, parsed, "FeeSplit", primaryAmt, treasuryAmt, buybackAmt),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.FeeSplitterFeeSplit)
		if got.Schedule != schedule || got.PrimaryAmount.Cmp(primaryAmt) != 0 {
			t.Errorf("fields: schedule=%d primaryAmount=%s", got.Schedule, got.PrimaryAmount)
		}
	})

	t.Run("TreasuryUpdated", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, fm, "TreasuryUpdated"),
				addrTopic(common.HexToAddress("0x8800000000000000000000000000000000000080")),
				addrTopic(common.HexToAddress("0x9900000000000000000000000000000000000090")),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.FeeSplitterTreasuryUpdated); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})

	t.Run("BuybackBurnerUpdated", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, fm, "BuybackBurnerUpdated"),
				addrTopic(common.HexToAddress("0xAA000000000000000000000000000000000000A0")),
				addrTopic(common.HexToAddress("0xBB000000000000000000000000000000000000B0")),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.FeeSplitterBuybackBurnerUpdated); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})
}

// Test_Dispatch_BuybackBurner covers all 4 BuybackBurner events.
func Test_Dispatch_BuybackBurner(t *testing.T) {
	parsed := mustAbi(t, contractabi.BuybackBurnerMetaData)
	bm := contractabi.BuybackBurnerMetaData

	t.Run("BuybackAndBurn", func(t *testing.T) {
		executor := common.HexToAddress("0xCC000000000000000000000000000000000000C0")
		paymentToken := common.HexToAddress("0xDD000000000000000000000000000000000000D0")
		amountIn := big.NewInt(10_000)
		tituBurned := big.NewInt(9_000)
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, bm, "BuybackAndBurn"),
				addrTopic(executor),
				addrTopic(paymentToken),
			},
			Data: encodeData(t, parsed, "BuybackAndBurn", amountIn, tituBurned),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.BuybackBurnerBuybackAndBurn)
		if got.AmountIn.Cmp(amountIn) != 0 || got.TituBurned.Cmp(tituBurned) != 0 {
			t.Errorf("fields: amountIn=%s tituBurned=%s", got.AmountIn, got.TituBurned)
		}
	})

	t.Run("RouterUpdated", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, bm, "RouterUpdated"),
				addrTopic(common.HexToAddress("0xEE000000000000000000000000000000000000E0")),
				addrTopic(common.HexToAddress("0xFF000000000000000000000000000000000000F0")),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.BuybackBurnerRouterUpdated); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})

	t.Run("SwapPathUpdated", func(t *testing.T) {
		path := []common.Address{
			common.HexToAddress("0x1100000000000000000000000000000000000011"),
			common.HexToAddress("0x2200000000000000000000000000000000000022"),
		}
		log := types.Log{
			Topics: []common.Hash{topicID(t, bm, "SwapPathUpdated")},
			Data:   encodeData(t, parsed, "SwapPathUpdated", path),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.BuybackBurnerSwapPathUpdated)
		if len(got.Path) != 2 || got.Path[0] != path[0] {
			t.Errorf("Path = %v", got.Path)
		}
	})

	t.Run("MinTituOutUpdated", func(t *testing.T) {
		oldMin := big.NewInt(100)
		newMin := big.NewInt(200)
		log := types.Log{
			Topics: []common.Hash{topicID(t, bm, "MinTituOutUpdated")},
			Data:   encodeData(t, parsed, "MinTituOutUpdated", oldMin, newMin),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.BuybackBurnerMinTituOutUpdated)
		if got.NewMin.Cmp(newMin) != 0 {
			t.Errorf("NewMin = %s", got.NewMin)
		}
	})
}

// Test_Dispatch_HookRegistry covers both HookRegistry events.
func Test_Dispatch_HookRegistry(t *testing.T) {
	parsed := mustAbi(t, contractabi.HookRegistryMetaData)
	hm := contractabi.HookRegistryMetaData

	t.Run("HookRegistered", func(t *testing.T) {
		hook := common.HexToAddress("0x3300000000000000000000000000000000000033")
		name := "fund-transfer-v1"
		log := types.Log{
			Topics: []common.Hash{topicID(t, hm, "HookRegistered"), addrTopic(hook)},
			Data:   encodeData(t, parsed, "HookRegistered", name),
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		got := ev.Payload.(*contractabi.HookRegistryHookRegistered)
		if got.Name != name || got.Hook != hook {
			t.Errorf("fields: hook=%s name=%q", got.Hook, got.Name)
		}
	})

	t.Run("HookDeregistered", func(t *testing.T) {
		log := types.Log{
			Topics: []common.Hash{
				topicID(t, hm, "HookDeregistered"),
				addrTopic(common.HexToAddress("0x4400000000000000000000000000000000000044")),
			},
			Data: []byte{},
		}
		ev, err := decoder.Dispatch(log)
		if err != nil {
			t.Fatalf("Dispatch: %v", err)
		}
		if _, ok := ev.Payload.(*contractabi.HookRegistryHookDeregistered); !ok {
			t.Fatalf("payload type = %T", ev.Payload)
		}
	})
}

// ──────────────────────────────────────────────────────────────────────────
// Error paths
// ──────────────────────────────────────────────────────────────────────────

// Test_Dispatch_NoTopics asserts that an empty Topics slice yields ErrNoTopics
// rather than panicking. Caller (the per-block pipeline) treats this as a
// dropped malformed log.
func Test_Dispatch_NoTopics(t *testing.T) {
	_, err := decoder.Dispatch(types.Log{Topics: nil})
	if !errors.Is(err, decoder.ErrNoTopics) {
		t.Fatalf("err = %v, want ErrNoTopics", err)
	}
}

// Test_Dispatch_UnknownTopic asserts that a topic0 we never registered is
// surfaced as ErrUnknownEvent. Production callers log-and-skip; the error
// is wrapped with the topic0 hex for debuggability.
func Test_Dispatch_UnknownTopic(t *testing.T) {
	bogus := common.HexToHash("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	_, err := decoder.Dispatch(types.Log{Topics: []common.Hash{bogus}})
	if !errors.Is(err, decoder.ErrUnknownEvent) {
		t.Fatalf("err = %v, want ErrUnknownEvent", err)
	}
}

// Test_Dispatch_MalformedData asserts that a known topic0 with garbage data
// surfaces a wrapped decode error rather than masking it. We pick AgentRegistered
// (which has non-indexed string + uint256 data) and supply data that is too
// short to unpack.
func Test_Dispatch_MalformedData(t *testing.T) {
	tID := topicID(t, contractabi.AgentRegistryMetaData, "AgentRegistered")
	log := types.Log{
		Topics: []common.Hash{
			tID,
			uintTopic(big.NewInt(1)),
			addrTopic(common.Address{}),
		},
		Data: []byte{0x01, 0x02}, // far too short for (string, uint256)
	}
	_, err := decoder.Dispatch(log)
	if err == nil {
		t.Fatal("Dispatch: expected error on malformed data")
	}
	// Should not be ErrUnknownEvent — the topic0 is known.
	if errors.Is(err, decoder.ErrUnknownEvent) {
		t.Fatalf("err classified as ErrUnknownEvent: %v", err)
	}
}

// ──────────────────────────────────────────────────────────────────────────
// Handler integration
// ──────────────────────────────────────────────────────────────────────────

// Test_DispatchAndHandle_Success verifies the convenience entry point routes
// a known log to the handler and reports (true, nil).
func Test_DispatchAndHandle_Success(t *testing.T) {
	parsed := mustAbi(t, contractabi.HookRegistryMetaData)
	tID := topicID(t, contractabi.HookRegistryMetaData, "HookRegistered")
	log := types.Log{
		Topics: []common.Hash{tID, addrTopic(common.HexToAddress("0x9999"))},
		Data:   encodeData(t, parsed, "HookRegistered", "test-hook"),
	}

	var seen *decoder.DecodedEvent
	h := decoder.HandlerFunc(func(_ context.Context, e decoder.DecodedEvent) error {
		v := e
		seen = &v
		return nil
	})
	handled, err := decoder.DispatchAndHandle(context.Background(), log, h)
	if err != nil {
		t.Fatalf("DispatchAndHandle: %v", err)
	}
	if !handled {
		t.Fatal("handled = false, want true")
	}
	if seen == nil || seen.Name != "HookRegistry.HookRegistered" {
		t.Fatalf("handler did not see expected event: %#v", seen)
	}
}

// Test_DispatchAndHandle_UnknownIsSkip verifies that an unknown topic0
// returns (false, nil) rather than propagating an error. This is the
// log-and-skip behaviour the per-block pipeline relies on.
func Test_DispatchAndHandle_UnknownIsSkip(t *testing.T) {
	bogus := common.HexToHash("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	called := false
	h := decoder.HandlerFunc(func(context.Context, decoder.DecodedEvent) error {
		called = true
		return nil
	})
	handled, err := decoder.DispatchAndHandle(context.Background(), types.Log{Topics: []common.Hash{bogus}}, h)
	if err != nil {
		t.Fatalf("DispatchAndHandle: %v", err)
	}
	if handled {
		t.Error("handled = true, want false for unknown topic")
	}
	if called {
		t.Error("handler called for unknown topic")
	}
}

// Test_DispatchAndHandle_HandlerError verifies that handler errors are
// wrapped with the qualified event name so operator dashboards can attribute
// failures to the right contract.
func Test_DispatchAndHandle_HandlerError(t *testing.T) {
	parsed := mustAbi(t, contractabi.HookRegistryMetaData)
	tID := topicID(t, contractabi.HookRegistryMetaData, "HookDeregistered")
	log := types.Log{
		Topics: []common.Hash{tID, addrTopic(common.HexToAddress("0xABCD"))},
		Data:   encodeData(t, parsed, "HookDeregistered"),
	}
	target := errors.New("downstream broken")
	h := decoder.HandlerFunc(func(context.Context, decoder.DecodedEvent) error { return target })

	handled, err := decoder.DispatchAndHandle(context.Background(), log, h)
	if handled {
		t.Error("handled = true, want false on handler error")
	}
	if !errors.Is(err, target) {
		t.Fatalf("err = %v, want wrap of %v", err, target)
	}
}

// Test_DispatchAndHandle_NilHandler asserts that a nil handler is a programmer
// error and surfaces explicitly rather than nil-deref'ing.
func Test_DispatchAndHandle_NilHandler(t *testing.T) {
	_, err := decoder.DispatchAndHandle(context.Background(), types.Log{}, nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

// ──────────────────────────────────────────────────────────────────────────
// Coverage gate: the dispatch table must contain a decoder for every event
// pinned in abi_test.go's TestACPEventSurface. This is a structural test —
// it does not exercise decoding, only registration, so a future regen that
// drops a binding cannot silently shrink the surface.
// ──────────────────────────────────────────────────────────────────────────

func Test_DispatchTable_CoversPinnedSurface(t *testing.T) {
	type ev struct {
		meta *bind.MetaData
		name string
	}
	// Mirror of abi_test.go::TestACPEventSurface — keep in sync.
	cases := []ev{
		{contractabi.AgentRegistryMetaData, "AgentRegistered"},
		{contractabi.AgentRegistryMetaData, "MetadataUpdated"},
		{contractabi.AgentRegistryMetaData, "CapabilitiesUpdated"},
		{contractabi.AgentRegistryMetaData, "ActiveStatusChanged"},
		{contractabi.AgentRegistryMetaData, "ScorePosted"},
		{contractabi.AgentRegistryMetaData, "ControllerTransferProposed"},
		{contractabi.AgentRegistryMetaData, "ControllerTransferAccepted"},
		{contractabi.AgentRegistryMetaData, "ControllerTransferCancelled"},
		{contractabi.JobMetaData, "JobInitialised"},
		{contractabi.JobMetaData, "AgentAccepted"},
		{contractabi.JobMetaData, "ResultSubmitted"},
		{contractabi.JobMetaData, "ResultApproved"},
		{contractabi.JobMetaData, "ResultRejected"},
		{contractabi.JobMetaData, "JobCompleted"},
		{contractabi.JobMetaData, "JobCancelled"},
		{contractabi.JobMetaData, "DisputeRaised"},
		{contractabi.JobMetaData, "DisputeResolved"},
		{contractabi.JobMetaData, "EvaluatorAssigned"},
		{contractabi.JobFactoryMetaData, "JobCreated"},
		{contractabi.JobFactoryMetaData, "ImplementationUpdated"},
		{contractabi.JobFactoryMetaData, "DefaultArbiterUpdated"},
		{contractabi.EscrowMetaData, "Funded"},
		{contractabi.EscrowMetaData, "Released"},
		{contractabi.EscrowMetaData, "Refunded"},
		{contractabi.FeeSplitterMetaData, "FeeSplit"},
		{contractabi.FeeSplitterMetaData, "TreasuryUpdated"},
		{contractabi.FeeSplitterMetaData, "BuybackBurnerUpdated"},
		{contractabi.BuybackBurnerMetaData, "BuybackAndBurn"},
		{contractabi.BuybackBurnerMetaData, "RouterUpdated"},
		{contractabi.BuybackBurnerMetaData, "SwapPathUpdated"},
		{contractabi.BuybackBurnerMetaData, "MinTituOutUpdated"},
		{contractabi.HookRegistryMetaData, "HookRegistered"},
		{contractabi.HookRegistryMetaData, "HookDeregistered"},
	}

	known := make(map[common.Hash]bool, len(decoder.KnownTopics()))
	for _, h := range decoder.KnownTopics() {
		known[h] = true
	}

	for _, c := range cases {
		id := topicID(t, c.meta, c.name)
		if !known[id] {
			t.Errorf("dispatch table missing decoder for event %q (topic0 %s)", c.name, id.Hex())
		}
	}
}
