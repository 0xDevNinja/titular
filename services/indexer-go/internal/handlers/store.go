// Package handlers decodes on-chain events and persists them via a Store.
//
// Each handler is idempotent: a duplicate (txHash, logIndex) pair is silently
// skipped. Block numbers are stored so callers can detect and rewind reorgs
// without assuming logs arrive in monotonic order.
package handlers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// AgentRecord holds the fields extracted from an AgentLaunched event.
type AgentRecord struct {
	AgentID  *big.Int
	Token    common.Address
	Curve    common.Address
	Creator  common.Address
	LpLock   common.Address
	Pair     common.Address
	Modules  [][32]byte
	BlockNum uint64
	TxHash   common.Hash
	LogIndex uint
}

// TradeRecord holds the fields from a Bought or Sold event.
// Direction is "buy" or "sell".
type TradeRecord struct {
	Direction string
	Trader    common.Address
	// Curve contract that emitted the event — used to look up the agent.
	Curve    common.Address
	QuoteIn  *big.Int // non-nil for buy; nil for sell
	AgentOut *big.Int // non-nil for buy; nil for sell
	AgentIn  *big.Int // non-nil for sell; nil for buy
	QuoteOut *big.Int // non-nil for sell; nil for buy
	Fee      *big.Int
	BlockNum uint64
	TxHash   common.Hash
	LogIndex uint
}

// GraduationRecord holds the fields from a Graduator.Graduated event.
type GraduationRecord struct {
	Curve    common.Address
	Pair     common.Address
	BlockNum uint64
	TxHash   common.Hash
	LogIndex uint
}

// ── ACP v2 record types ───────────────────────────────────────────────────────

// ACPAgentRegisteredRecord holds the fields from AgentRegistry.AgentRegistered.
type ACPAgentRegisteredRecord struct {
	AgentID      *big.Int
	Controller   common.Address
	MetadataURI  string
	Capabilities *big.Int
	BlockNum     uint64
	TxHash       common.Hash
	LogIndex     uint
}

// ACPScorePostedRecord holds the fields from AgentRegistry.ScorePosted.
type ACPScorePostedRecord struct {
	AgentID  *big.Int
	Scorer   common.Address
	Delta    *big.Int
	NewTotal *big.Int
	Nonce    *big.Int
	BlockNum uint64
	TxHash   common.Hash
	LogIndex uint
}

// ACPControllerTransferRecord holds the fields from
// AgentRegistry.ControllerTransferAccepted.
type ACPControllerTransferRecord struct {
	AgentID       *big.Int
	OldController common.Address
	NewController common.Address
	BlockNum      uint64
	TxHash        common.Hash
	LogIndex      uint
}

// ACPJobRecord holds the fields from Job.JobInitialised.
type ACPJobRecord struct {
	JobID     *big.Int
	Principal common.Address
	AgentID   *big.Int
	JobType   uint8
	Budget    *big.Int
	Token     common.Address
	Deadline  uint64
	BlockNum  uint64
	TxHash    common.Hash
	LogIndex  uint
}

// ACPJobAcceptedRecord holds the fields from Job.AgentAccepted.
type ACPJobAcceptedRecord struct {
	JobID    *big.Int
	Agent    common.Address
	AgentID  *big.Int
	BlockNum uint64
	TxHash   common.Hash
	LogIndex uint
}

// ACPResultSubmittedRecord holds the fields from Job.ResultSubmitted.
type ACPResultSubmittedRecord struct {
	JobID     *big.Int
	Agent     common.Address
	ResultURI string
	BlockNum  uint64
	TxHash    common.Hash
	LogIndex  uint
}

// ACPJobCompletedRecord holds the fields from Job.JobCompleted.
type ACPJobCompletedRecord struct {
	JobID      *big.Int
	ReleasedBy common.Address
	Amount     *big.Int
	BlockNum   uint64
	TxHash     common.Hash
	LogIndex   uint
}

// ACPJobCancelledRecord holds the fields from Job.JobCancelled.
type ACPJobCancelledRecord struct {
	JobID       *big.Int
	CancelledBy common.Address
	Reason      string
	BlockNum    uint64
	TxHash      common.Hash
	LogIndex    uint
}

// ACPJobCreatedRecord holds the fields from JobFactory.JobCreated.
type ACPJobCreatedRecord struct {
	JobID     *big.Int
	Clone     common.Address
	Principal common.Address
	JobType   uint8
	Token     common.Address
	Budget    *big.Int
	BlockNum  uint64
	TxHash    common.Hash
	LogIndex  uint
}

// ACPEscrowFundedRecord holds the fields from Escrow.Funded.
type ACPEscrowFundedRecord struct {
	Depositor common.Address
	JobID     *big.Int
	Token     common.Address
	Amount    *big.Int
	BlockNum  uint64
	TxHash    common.Hash
	LogIndex  uint
}

// ACPEscrowReleasedRecord holds the fields from Escrow.Released.
type ACPEscrowReleasedRecord struct {
	Depositor common.Address
	JobID     *big.Int
	Token     common.Address
	To        common.Address
	Amount    *big.Int
	BlockNum  uint64
	TxHash    common.Hash
	LogIndex  uint
}

// ACPEscrowRefundedRecord holds the fields from Escrow.Refunded.
type ACPEscrowRefundedRecord struct {
	Depositor common.Address
	JobID     *big.Int
	Token     common.Address
	Amount    *big.Int
	BlockNum  uint64
	TxHash    common.Hash
	LogIndex  uint
}

// Store is the persistence interface used by all handlers.
// Implementations must be safe for concurrent use.
type Store interface {
	// IsLogProcessed returns true when a log identified by (txHash, logIndex)
	// has already been persisted. It is called before every write to enforce
	// idempotency.
	IsLogProcessed(ctx context.Context, txHash common.Hash, logIndex uint) (bool, error)

	// MarkLogProcessed records that a log identified by (txHash, logIndex) has
	// been fully persisted. Handlers call this immediately after a successful
	// write so that retries skip the log without re-applying delta-style
	// mutations (e.g. InsertACPScore score deltas).
	MarkLogProcessed(ctx context.Context, txHash common.Hash, logIndex uint) error

	// UpsertAgent inserts or updates an agent record. On conflict the existing
	// row is left unchanged (insert-if-not-exists semantics).
	UpsertAgent(ctx context.Context, rec AgentRecord) error

	// InsertTrade persists a single trade. The (txHash, logIndex) uniqueness
	// is enforced by the caller via IsLogProcessed before this is called.
	InsertTrade(ctx context.Context, rec TradeRecord) error

	// MarkGraduated marks an agent as graduated and stores the Uniswap V2
	// pair address. Identified by curve address.
	MarkGraduated(ctx context.Context, rec GraduationRecord) error

	// ── ACP v2 ───────────────────────────────────────────────────────────────

	// UpsertACPAgent persists an agent identity from AgentRegistry.AgentRegistered.
	UpsertACPAgent(ctx context.Context, rec ACPAgentRegisteredRecord) error

	// InsertACPScore persists a reputation delta from AgentRegistry.ScorePosted.
	InsertACPScore(ctx context.Context, rec ACPScorePostedRecord) error

	// UpdateACPController persists a controller change from
	// AgentRegistry.ControllerTransferAccepted.
	UpdateACPController(ctx context.Context, rec ACPControllerTransferRecord) error

	// UpsertACPJob persists a job record from Job.JobInitialised.
	UpsertACPJob(ctx context.Context, rec ACPJobRecord) error

	// UpdateACPJobAccepted records agent acceptance from Job.AgentAccepted.
	UpdateACPJobAccepted(ctx context.Context, rec ACPJobAcceptedRecord) error

	// UpdateACPResultSubmitted records result submission from Job.ResultSubmitted.
	UpdateACPResultSubmitted(ctx context.Context, rec ACPResultSubmittedRecord) error

	// UpdateACPJobCompleted marks a job completed from Job.JobCompleted.
	UpdateACPJobCompleted(ctx context.Context, rec ACPJobCompletedRecord) error

	// UpdateACPJobCancelled marks a job cancelled from Job.JobCancelled.
	UpdateACPJobCancelled(ctx context.Context, rec ACPJobCancelledRecord) error

	// InsertACPJobCreated persists the factory record from JobFactory.JobCreated.
	InsertACPJobCreated(ctx context.Context, rec ACPJobCreatedRecord) error

	// InsertACPEscrowFunded persists an escrow deposit from Escrow.Funded.
	InsertACPEscrowFunded(ctx context.Context, rec ACPEscrowFundedRecord) error

	// InsertACPEscrowReleased persists a release entry from Escrow.Released.
	InsertACPEscrowReleased(ctx context.Context, rec ACPEscrowReleasedRecord) error

	// InsertACPEscrowRefunded persists a refund entry from Escrow.Refunded.
	InsertACPEscrowRefunded(ctx context.Context, rec ACPEscrowRefundedRecord) error
}
