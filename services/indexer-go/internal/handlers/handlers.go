package handlers

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	contractabi "github.com/0xDevNinja/titular/services/indexer-go/internal/abi"
)

// zeroAddr is used when constructing filterers for offline log parsing.
// The address is irrelevant because Parse methods only call UnpackLog,
// which does not require an RPC backend.
var zeroAddr = common.Address{}

// HandleAgentLaunched decodes a LaunchpadFactory AgentLaunched log and
// persists an agent record via store. It is idempotent: duplicate logs
// (same txHash + logIndex) are silently skipped.
func HandleAgentLaunched(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleAgentLaunched: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewLaunchpadFactoryFilterer(zeroAddr, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleAgentLaunched: new filterer: %w", err)
	}

	ev, err := f.ParseAgentLaunched(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleAgentLaunched: parse: %w", err)
	}

	rec := AgentRecord{
		AgentID:  ev.AgentId,
		Token:    ev.Token,
		Curve:    ev.Curve,
		Creator:  ev.Creator,
		LpLock:   ev.LpLock,
		Pair:     ev.Pair,
		Modules:  ev.Modules,
		BlockNum: log.BlockNumber,
		TxHash:   log.TxHash,
		LogIndex: log.Index,
	}

	if err := store.UpsertAgent(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleAgentLaunched: upsert: %w", err)
	}

	return nil
}

// HandleBought decodes a BondingCurve Bought log and persists a trade record.
// It is idempotent: duplicate logs are silently skipped.
func HandleBought(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleBought: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewBondingCurveFilterer(zeroAddr, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleBought: new filterer: %w", err)
	}

	ev, err := f.ParseBought(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleBought: parse: %w", err)
	}

	rec := TradeRecord{
		Direction: "buy",
		Trader:    ev.Buyer,
		Curve:     log.Address,
		QuoteIn:   ev.QuoteIn,
		AgentOut:  ev.AgentOut,
		Fee:       ev.Fee,
		BlockNum:  log.BlockNumber,
		TxHash:    log.TxHash,
		LogIndex:  log.Index,
	}

	if err := store.InsertTrade(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleBought: insert: %w", err)
	}

	return nil
}

// HandleSold decodes a BondingCurve Sold log and persists a trade record.
// It is idempotent: duplicate logs are silently skipped.
func HandleSold(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleSold: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewBondingCurveFilterer(zeroAddr, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleSold: new filterer: %w", err)
	}

	ev, err := f.ParseSold(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleSold: parse: %w", err)
	}

	rec := TradeRecord{
		Direction: "sell",
		Trader:    ev.Seller,
		Curve:     log.Address,
		AgentIn:   ev.AgentIn,
		QuoteOut:  ev.QuoteOut,
		Fee:       ev.Fee,
		BlockNum:  log.BlockNumber,
		TxHash:    log.TxHash,
		LogIndex:  log.Index,
	}

	if err := store.InsertTrade(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleSold: insert: %w", err)
	}

	return nil
}

// HandleGraduated decodes a Graduator Graduated log, marks the associated
// agent as graduated, and persists the Uniswap V2 pair address.
// The event is emitted by the Graduator contract, which includes both the
// curve address and the pair address as indexed fields.
// It is idempotent: duplicate logs are silently skipped.
func HandleGraduated(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleGraduated: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewGraduatorFilterer(zeroAddr, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleGraduated: new filterer: %w", err)
	}

	ev, err := f.ParseGraduated(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleGraduated: parse: %w", err)
	}

	rec := GraduationRecord{
		Curve:    ev.Curve,
		Pair:     ev.Pair,
		BlockNum: log.BlockNumber,
		TxHash:   log.TxHash,
		LogIndex: log.Index,
	}

	if err := store.MarkGraduated(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleGraduated: mark: %w", err)
	}

	return nil
}
