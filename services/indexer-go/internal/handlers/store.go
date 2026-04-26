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

// Store is the persistence interface used by all handlers.
// Implementations must be safe for concurrent use.
type Store interface {
	// IsLogProcessed returns true when a log identified by (txHash, logIndex)
	// has already been persisted. It is called before every write to enforce
	// idempotency.
	IsLogProcessed(ctx context.Context, txHash common.Hash, logIndex uint) (bool, error)

	// UpsertAgent inserts or updates an agent record. On conflict the existing
	// row is left unchanged (insert-if-not-exists semantics).
	UpsertAgent(ctx context.Context, rec AgentRecord) error

	// InsertTrade persists a single trade. The (txHash, logIndex) uniqueness
	// is enforced by the caller via IsLogProcessed before this is called.
	InsertTrade(ctx context.Context, rec TradeRecord) error

	// MarkGraduated marks an agent as graduated and stores the Uniswap V2
	// pair address. Identified by curve address.
	MarkGraduated(ctx context.Context, rec GraduationRecord) error
}
