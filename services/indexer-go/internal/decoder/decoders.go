package decoder

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	contractabi "github.com/0xDevNinja/titular/services/indexer-go/internal/abi"
)

// This file installs one [Decoder] per Solidity event the indexer subscribes
// to. The list is the same surface pinned in `internal/abi/abi_test.go`'s
// TestACPEventSurface — every entry there must have a decoder here, and a
// decoder here without a corresponding pin in abi_test would be considered
// drift (not an error per se, but worth flagging in review).
//
// Implementation notes:
//
//   - All filterers are constructed once at init time, bound to the zero
//     address with a nil backend. abigen's Parse* methods only call
//     UnpackLog under the hood, which uses the embedded ABI and never
//     touches the backend, so this is safe and avoids per-log allocations.
//   - We deliberately surface filterer-construction errors as panics at
//     init: the only way they fail is a malformed embedded ABI, which would
//     have already blown up TestAllBindingsCompile. Panicking at process
//     start is preferable to silently emitting decode errors per log later.

// mustFiltererErr panics if a filterer constructor returns an error. This
// drains the noise from the var block below and gives a single, scoped
// panic point if abigen output ever ships with a bad ABI.
func mustFiltererErr(name string, err error) {
	if err != nil {
		panic(fmt.Sprintf("decoder.init: filterer constructor for %s failed: %v", name, err))
	}
}

// zeroAddr is the placeholder address passed to filterer constructors. It
// has no effect because Parse* operations only consume the log's own
// fields; the constructor stores the address purely for FilterLogs queries
// we never make.
var zeroAddr = common.Address{}

// Cached filterers — each constructed once and reused across every Dispatch
// call. Sharing is safe: the abigen filterer types are read-only after
// construction (they hold an ABI parse tree and a contract handle).
var (
	agentRegistryFilterer *contractabi.AgentRegistryFilterer
	jobFilterer           *contractabi.JobFilterer
	jobFactoryFilterer    *contractabi.JobFactoryFilterer
	escrowFilterer        *contractabi.EscrowFilterer
	feeSplitterFilterer   *contractabi.FeeSplitterFilterer
	buybackBurnerFilterer *contractabi.BuybackBurnerFilterer
	hookRegistryFilterer  *contractabi.HookRegistryFilterer
)

// init populates the dispatch table. Topic0 lookups happen via
// MetaData.GetAbi at startup (one-time cost); after that the table is
// read-only. Any topic0 collision triggers a panic in [Register], which is
// the desired behaviour: a duplicate hash means the abigen output disagrees
// with itself and the indexer cannot start until that is fixed.
func init() {
	var err error

	agentRegistryFilterer, err = contractabi.NewAgentRegistryFilterer(zeroAddr, nil)
	mustFiltererErr("AgentRegistry", err)

	jobFilterer, err = contractabi.NewJobFilterer(zeroAddr, nil)
	mustFiltererErr("Job", err)

	jobFactoryFilterer, err = contractabi.NewJobFactoryFilterer(zeroAddr, nil)
	mustFiltererErr("JobFactory", err)

	escrowFilterer, err = contractabi.NewEscrowFilterer(zeroAddr, nil)
	mustFiltererErr("Escrow", err)

	feeSplitterFilterer, err = contractabi.NewFeeSplitterFilterer(zeroAddr, nil)
	mustFiltererErr("FeeSplitter", err)

	buybackBurnerFilterer, err = contractabi.NewBuybackBurnerFilterer(zeroAddr, nil)
	mustFiltererErr("BuybackBurner", err)

	hookRegistryFilterer, err = contractabi.NewHookRegistryFilterer(zeroAddr, nil)
	mustFiltererErr("HookRegistry", err)

	// AgentRegistry — agent identity & reputation lifecycle.
	registerFromMeta(contractabi.AgentRegistryMetaData, "AgentRegistry", []eventReg{
		{"AgentRegistered", decodeAgentRegistered},
		{"MetadataUpdated", decodeMetadataUpdated},
		{"CapabilitiesUpdated", decodeCapabilitiesUpdated},
		{"ActiveStatusChanged", decodeActiveStatusChanged},
		{"ScorePosted", decodeScorePosted},
		{"ControllerTransferProposed", decodeControllerTransferProposed},
		{"ControllerTransferAccepted", decodeControllerTransferAccepted},
		{"ControllerTransferCancelled", decodeControllerTransferCancelled},
	})

	// Job — per-job state machine (one Job clone per JobCreated).
	registerFromMeta(contractabi.JobMetaData, "Job", []eventReg{
		{"JobInitialised", decodeJobInitialised},
		{"AgentAccepted", decodeAgentAccepted},
		{"ResultSubmitted", decodeResultSubmitted},
		{"ResultApproved", decodeResultApproved},
		{"ResultRejected", decodeResultRejected},
		{"JobCompleted", decodeJobCompleted},
		{"JobCancelled", decodeJobCancelled},
		{"DisputeRaised", decodeDisputeRaised},
		{"DisputeResolved", decodeDisputeResolved},
		{"EvaluatorAssigned", decodeEvaluatorAssigned},
	})

	// JobFactory — clone deployments and admin updates.
	registerFromMeta(contractabi.JobFactoryMetaData, "JobFactory", []eventReg{
		{"JobCreated", decodeJobCreated},
		{"ImplementationUpdated", decodeImplementationUpdated},
		{"DefaultArbiterUpdated", decodeDefaultArbiterUpdated},
	})

	// Escrow — funds custody & release/refund flow.
	registerFromMeta(contractabi.EscrowMetaData, "Escrow", []eventReg{
		{"Funded", decodeFunded},
		{"Released", decodeReleased},
		{"Refunded", decodeRefunded},
	})

	// FeeSplitter — protocol/treasury/buyback split distribution.
	registerFromMeta(contractabi.FeeSplitterMetaData, "FeeSplitter", []eventReg{
		{"FeeSplit", decodeFeeSplit},
		{"TreasuryUpdated", decodeTreasuryUpdated},
		{"BuybackBurnerUpdated", decodeBuybackBurnerUpdated},
	})

	// BuybackBurner — TITU buyback execution and config.
	registerFromMeta(contractabi.BuybackBurnerMetaData, "BuybackBurner", []eventReg{
		{"BuybackAndBurn", decodeBuybackAndBurn},
		{"RouterUpdated", decodeRouterUpdated},
		{"SwapPathUpdated", decodeSwapPathUpdated},
		{"MinTituOutUpdated", decodeMinTituOutUpdated},
	})

	// HookRegistry — approved-hook governance.
	registerFromMeta(contractabi.HookRegistryMetaData, "HookRegistry", []eventReg{
		{"HookRegistered", decodeHookRegistered},
		{"HookDeregistered", decodeHookDeregistered},
	})
}

// eventReg pairs a Solidity event name with its in-package decoder. Used
// only inside init for tabular registration.
type eventReg struct {
	name    string
	decoder Decoder
}

// registerFromMeta resolves topic0 for each event from the embedded ABI and
// installs it in the dispatch table under the contract-qualified name
// "<contract>.<event>". Panics if the ABI is malformed or the event is
// missing — both indicate a broken bindings regeneration that must be
// fixed before the indexer can boot.
func registerFromMeta(meta *bind.MetaData, contractName string, events []eventReg) {
	parsed, err := meta.GetAbi()
	if err != nil {
		panic(fmt.Sprintf("decoder.init: GetAbi(%s): %v", contractName, err))
	}
	for _, e := range events {
		ev, ok := parsed.Events[e.name]
		if !ok {
			panic(fmt.Sprintf(
				"decoder.init: event %q missing from %s ABI — bindings drift?",
				e.name, contractName,
			))
		}
		Register(ev.ID, contractName+"."+e.name, e.decoder)
	}
}

// ──────────────────────────────────────────────────────────────────────────
// Decoders — one per tracked event.
//
// Each decoder is a thin wrapper over the corresponding abigen Parse* call.
// They return the raw event struct as `any` so [Dispatch] can hand it back
// in [DecodedEvent.Payload]; downstream handlers type-assert to the
// concrete type they need.
// ──────────────────────────────────────────────────────────────────────────

// AgentRegistry decoders.

func decodeAgentRegistered(log types.Log) (any, error) {
	return agentRegistryFilterer.ParseAgentRegistered(log)
}

func decodeMetadataUpdated(log types.Log) (any, error) {
	return agentRegistryFilterer.ParseMetadataUpdated(log)
}

func decodeCapabilitiesUpdated(log types.Log) (any, error) {
	return agentRegistryFilterer.ParseCapabilitiesUpdated(log)
}

func decodeActiveStatusChanged(log types.Log) (any, error) {
	return agentRegistryFilterer.ParseActiveStatusChanged(log)
}

func decodeScorePosted(log types.Log) (any, error) {
	return agentRegistryFilterer.ParseScorePosted(log)
}

func decodeControllerTransferProposed(log types.Log) (any, error) {
	return agentRegistryFilterer.ParseControllerTransferProposed(log)
}

func decodeControllerTransferAccepted(log types.Log) (any, error) {
	return agentRegistryFilterer.ParseControllerTransferAccepted(log)
}

func decodeControllerTransferCancelled(log types.Log) (any, error) {
	return agentRegistryFilterer.ParseControllerTransferCancelled(log)
}

// Job decoders.

func decodeJobInitialised(log types.Log) (any, error) {
	return jobFilterer.ParseJobInitialised(log)
}

func decodeAgentAccepted(log types.Log) (any, error) {
	return jobFilterer.ParseAgentAccepted(log)
}

func decodeResultSubmitted(log types.Log) (any, error) {
	return jobFilterer.ParseResultSubmitted(log)
}

func decodeResultApproved(log types.Log) (any, error) {
	return jobFilterer.ParseResultApproved(log)
}

func decodeResultRejected(log types.Log) (any, error) {
	return jobFilterer.ParseResultRejected(log)
}

func decodeJobCompleted(log types.Log) (any, error) {
	return jobFilterer.ParseJobCompleted(log)
}

func decodeJobCancelled(log types.Log) (any, error) {
	return jobFilterer.ParseJobCancelled(log)
}

func decodeDisputeRaised(log types.Log) (any, error) {
	return jobFilterer.ParseDisputeRaised(log)
}

func decodeDisputeResolved(log types.Log) (any, error) {
	return jobFilterer.ParseDisputeResolved(log)
}

func decodeEvaluatorAssigned(log types.Log) (any, error) {
	return jobFilterer.ParseEvaluatorAssigned(log)
}

// JobFactory decoders.

func decodeJobCreated(log types.Log) (any, error) {
	return jobFactoryFilterer.ParseJobCreated(log)
}

func decodeImplementationUpdated(log types.Log) (any, error) {
	return jobFactoryFilterer.ParseImplementationUpdated(log)
}

func decodeDefaultArbiterUpdated(log types.Log) (any, error) {
	return jobFactoryFilterer.ParseDefaultArbiterUpdated(log)
}

// Escrow decoders.

func decodeFunded(log types.Log) (any, error) {
	return escrowFilterer.ParseFunded(log)
}

func decodeReleased(log types.Log) (any, error) {
	return escrowFilterer.ParseReleased(log)
}

func decodeRefunded(log types.Log) (any, error) {
	return escrowFilterer.ParseRefunded(log)
}

// FeeSplitter decoders.

func decodeFeeSplit(log types.Log) (any, error) {
	return feeSplitterFilterer.ParseFeeSplit(log)
}

func decodeTreasuryUpdated(log types.Log) (any, error) {
	return feeSplitterFilterer.ParseTreasuryUpdated(log)
}

func decodeBuybackBurnerUpdated(log types.Log) (any, error) {
	return feeSplitterFilterer.ParseBuybackBurnerUpdated(log)
}

// BuybackBurner decoders.

func decodeBuybackAndBurn(log types.Log) (any, error) {
	return buybackBurnerFilterer.ParseBuybackAndBurn(log)
}

func decodeRouterUpdated(log types.Log) (any, error) {
	return buybackBurnerFilterer.ParseRouterUpdated(log)
}

func decodeSwapPathUpdated(log types.Log) (any, error) {
	return buybackBurnerFilterer.ParseSwapPathUpdated(log)
}

func decodeMinTituOutUpdated(log types.Log) (any, error) {
	return buybackBurnerFilterer.ParseMinTituOutUpdated(log)
}

// HookRegistry decoders.

func decodeHookRegistered(log types.Log) (any, error) {
	return hookRegistryFilterer.ParseHookRegistered(log)
}

func decodeHookDeregistered(log types.Log) (any, error) {
	return hookRegistryFilterer.ParseHookDeregistered(log)
}
