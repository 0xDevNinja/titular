package handlers

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	contractabi "github.com/0xDevNinja/titular/services/indexer-go/internal/abi"
)

// All ACP handlers follow the same three-step pattern:
//  1. Idempotency guard via IsLogProcessed.
//  2. Parse the raw log using the generated abigen filterer (no RPC needed).
//  3. Persist via the Store interface.

// HandleACPAgentRegistered decodes an AgentRegistry.AgentRegistered log and
// persists an ACP agent identity record. Idempotent on (txHash, logIndex).
func HandleACPAgentRegistered(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPAgentRegistered: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewAgentRegistryFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPAgentRegistered: new filterer: %w", err)
	}

	ev, err := f.ParseAgentRegistered(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPAgentRegistered: parse: %w", err)
	}

	rec := ACPAgentRegisteredRecord{
		AgentID:      ev.AgentId,
		Controller:   ev.Controller,
		MetadataURI:  ev.MetadataURI,
		Capabilities: ev.Capabilities,
		BlockNum:     log.BlockNumber,
		TxHash:       log.TxHash,
		LogIndex:     log.Index,
	}

	if err := store.UpsertACPAgent(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPAgentRegistered: upsert: %w", err)
	}

	return nil
}

// HandleACPScorePosted decodes an AgentRegistry.ScorePosted log and persists
// the reputation delta. Idempotent on (txHash, logIndex).
func HandleACPScorePosted(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPScorePosted: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewAgentRegistryFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPScorePosted: new filterer: %w", err)
	}

	ev, err := f.ParseScorePosted(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPScorePosted: parse: %w", err)
	}

	rec := ACPScorePostedRecord{
		AgentID:  ev.AgentId,
		Scorer:   ev.Scorer,
		Delta:    ev.Delta,
		NewTotal: ev.NewTotal,
		Nonce:    ev.Nonce,
		BlockNum: log.BlockNumber,
		TxHash:   log.TxHash,
		LogIndex: log.Index,
	}

	if err := store.InsertACPScore(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPScorePosted: insert: %w", err)
	}

	return nil
}

// HandleACPControllerTransferAccepted decodes an
// AgentRegistry.ControllerTransferAccepted log and persists the ownership
// change. Idempotent on (txHash, logIndex).
func HandleACPControllerTransferAccepted(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPControllerTransferAccepted: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewAgentRegistryFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPControllerTransferAccepted: new filterer: %w", err)
	}

	ev, err := f.ParseControllerTransferAccepted(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPControllerTransferAccepted: parse: %w", err)
	}

	rec := ACPControllerTransferRecord{
		AgentID:       ev.AgentId,
		OldController: ev.OldController,
		NewController: ev.NewController,
		BlockNum:      log.BlockNumber,
		TxHash:        log.TxHash,
		LogIndex:      log.Index,
	}

	if err := store.UpdateACPController(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPControllerTransferAccepted: update: %w", err)
	}

	return nil
}

// HandleACPJobInitialised decodes a Job.JobInitialised log and persists the
// initial job record. Idempotent on (txHash, logIndex).
func HandleACPJobInitialised(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobInitialised: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewJobFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobInitialised: new filterer: %w", err)
	}

	ev, err := f.ParseJobInitialised(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobInitialised: parse: %w", err)
	}

	rec := ACPJobRecord{
		JobID:     ev.JobId,
		Principal: ev.Principal,
		AgentID:   ev.AgentId,
		JobType:   ev.JobType,
		Budget:    ev.Budget,
		Token:     ev.Token,
		Deadline:  ev.Deadline,
		BlockNum:  log.BlockNumber,
		TxHash:    log.TxHash,
		LogIndex:  log.Index,
	}

	if err := store.UpsertACPJob(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPJobInitialised: upsert: %w", err)
	}

	return nil
}

// HandleACPAgentAccepted decodes a Job.AgentAccepted log and records the agent
// assignment. Idempotent on (txHash, logIndex).
func HandleACPAgentAccepted(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPAgentAccepted: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewJobFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPAgentAccepted: new filterer: %w", err)
	}

	ev, err := f.ParseAgentAccepted(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPAgentAccepted: parse: %w", err)
	}

	rec := ACPJobAcceptedRecord{
		JobID:    ev.JobId,
		Agent:    ev.Agent,
		AgentID:  ev.AgentId,
		BlockNum: log.BlockNumber,
		TxHash:   log.TxHash,
		LogIndex: log.Index,
	}

	if err := store.UpdateACPJobAccepted(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPAgentAccepted: update: %w", err)
	}

	return nil
}

// HandleACPResultSubmitted decodes a Job.ResultSubmitted log and records the
// result URI. Idempotent on (txHash, logIndex).
func HandleACPResultSubmitted(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPResultSubmitted: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewJobFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPResultSubmitted: new filterer: %w", err)
	}

	ev, err := f.ParseResultSubmitted(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPResultSubmitted: parse: %w", err)
	}

	rec := ACPResultSubmittedRecord{
		JobID:     ev.JobId,
		Agent:     ev.Agent,
		ResultURI: ev.ResultURI,
		BlockNum:  log.BlockNumber,
		TxHash:    log.TxHash,
		LogIndex:  log.Index,
	}

	if err := store.UpdateACPResultSubmitted(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPResultSubmitted: update: %w", err)
	}

	return nil
}

// HandleACPJobCompleted decodes a Job.JobCompleted log and marks the job
// terminal-completed. Idempotent on (txHash, logIndex).
func HandleACPJobCompleted(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobCompleted: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewJobFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobCompleted: new filterer: %w", err)
	}

	ev, err := f.ParseJobCompleted(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobCompleted: parse: %w", err)
	}

	rec := ACPJobCompletedRecord{
		JobID:      ev.JobId,
		ReleasedBy: ev.ReleasedBy,
		Amount:     ev.Amount,
		BlockNum:   log.BlockNumber,
		TxHash:     log.TxHash,
		LogIndex:   log.Index,
	}

	if err := store.UpdateACPJobCompleted(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPJobCompleted: update: %w", err)
	}

	return nil
}

// HandleACPJobCancelled decodes a Job.JobCancelled log and marks the job
// terminal-cancelled. Idempotent on (txHash, logIndex).
func HandleACPJobCancelled(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobCancelled: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewJobFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobCancelled: new filterer: %w", err)
	}

	ev, err := f.ParseJobCancelled(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobCancelled: parse: %w", err)
	}

	rec := ACPJobCancelledRecord{
		JobID:       ev.JobId,
		CancelledBy: ev.CancelledBy,
		Reason:      ev.Reason,
		BlockNum:    log.BlockNumber,
		TxHash:      log.TxHash,
		LogIndex:    log.Index,
	}

	if err := store.UpdateACPJobCancelled(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPJobCancelled: update: %w", err)
	}

	return nil
}

// HandleACPJobCreated decodes a JobFactory.JobCreated log and persists the
// factory-level record (clone address, principal). Idempotent on (txHash, logIndex).
func HandleACPJobCreated(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobCreated: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewJobFactoryFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobCreated: new filterer: %w", err)
	}

	ev, err := f.ParseJobCreated(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPJobCreated: parse: %w", err)
	}

	rec := ACPJobCreatedRecord{
		JobID:     ev.JobId,
		Clone:     ev.Clone,
		Principal: ev.Principal,
		JobType:   ev.JobType,
		Token:     ev.Token,
		Budget:    ev.Budget,
		BlockNum:  log.BlockNumber,
		TxHash:    log.TxHash,
		LogIndex:  log.Index,
	}

	if err := store.InsertACPJobCreated(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPJobCreated: insert: %w", err)
	}

	return nil
}

// HandleACPEscrowFunded decodes an Escrow.Funded log and persists the deposit.
// Idempotent on (txHash, logIndex).
func HandleACPEscrowFunded(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowFunded: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewEscrowFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowFunded: new filterer: %w", err)
	}

	ev, err := f.ParseFunded(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowFunded: parse: %w", err)
	}

	rec := ACPEscrowFundedRecord{
		Depositor: ev.Depositor,
		JobID:     ev.JobId,
		Token:     ev.Token,
		Amount:    ev.Amount,
		BlockNum:  log.BlockNumber,
		TxHash:    log.TxHash,
		LogIndex:  log.Index,
	}

	if err := store.InsertACPEscrowFunded(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowFunded: insert: %w", err)
	}

	return nil
}

// HandleACPEscrowReleased decodes an Escrow.Released log and persists the
// per-recipient release entry. Idempotent on (txHash, logIndex).
func HandleACPEscrowReleased(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowReleased: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewEscrowFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowReleased: new filterer: %w", err)
	}

	ev, err := f.ParseReleased(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowReleased: parse: %w", err)
	}

	rec := ACPEscrowReleasedRecord{
		Depositor: ev.Depositor,
		JobID:     ev.JobId,
		Token:     ev.Token,
		To:        ev.To,
		Amount:    ev.Amount,
		BlockNum:  log.BlockNumber,
		TxHash:    log.TxHash,
		LogIndex:  log.Index,
	}

	if err := store.InsertACPEscrowReleased(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowReleased: insert: %w", err)
	}

	return nil
}

// HandleACPEscrowRefunded decodes an Escrow.Refunded log and persists the
// refund entry. Idempotent on (txHash, logIndex).
func HandleACPEscrowRefunded(ctx context.Context, log types.Log, store Store) error {
	already, err := store.IsLogProcessed(ctx, log.TxHash, log.Index)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowRefunded: check processed: %w", err)
	}
	if already {
		return nil
	}

	f, err := contractabi.NewEscrowFilterer(common.Address{}, nil)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowRefunded: new filterer: %w", err)
	}

	ev, err := f.ParseRefunded(log)
	if err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowRefunded: parse: %w", err)
	}

	rec := ACPEscrowRefundedRecord{
		Depositor: ev.Depositor,
		JobID:     ev.JobId,
		Token:     ev.Token,
		Amount:    ev.Amount,
		BlockNum:  log.BlockNumber,
		TxHash:    log.TxHash,
		LogIndex:  log.Index,
	}

	if err := store.InsertACPEscrowRefunded(ctx, rec); err != nil {
		return fmt.Errorf("handlers.HandleACPEscrowRefunded: insert: %w", err)
	}

	return nil
}
