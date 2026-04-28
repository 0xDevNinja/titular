# M3 ACP v2 — Slither Audit Summary

**Date**: 2026-04-26  
**Tool**: Slither + crytic-compile (forge build-info backend)  
**Scope**: `contracts/evm/src/acp/`  

## Results per contract

| Contract | Detectors Run | Findings | Suppressions |
|---|---|---|---|
| AgentRegistry.sol | 80 | 0 | 3 `incorrect-equality,timestamp` (registeredAt==0 sentinel) |
| Job.sol | 80 | 0 | 1 `timestamp` (deadline validation) |
| JobFactory.sol | 80 | 0 | 0 (reentrancy-events resolved by emit-before-call reorder) |
| Escrow.sol | 80 | 0 | 0 |
| FeeSplitter.sol | 80 | 0 | 0 |
| BuybackBurner.sol | 80 | 0 | 0 (made paymentToken/titu/swapDeadlineBuffer immutable) |
| HookRegistry.sol | 80 | 0 | 0 |
| IHook.sol | 80 | 0 | n/a (interface only) |
| FundTransferHook.sol | 80 | 0 | 0 |
| SubscriptionHook.sol | 80 | 0 | 1 `arbitrary-send-erc20` (subscriber is a trusted param from Job contract) |
| MilestoneHook.sol | 80 | 0 | 1 `divide-before-multiply` (intentional integer remainder using %) |
| RoyaltyHook.sol | 80 | 0 | 0 |
| IERC8183.sol | 80 | 0 | n/a (interface only) |
| IPermit2.sol | n/a | n/a | n/a (interface only, not directly analyzed) |

**Total high/medium unresolved: 0**

## Suppressions rationale

### `incorrect-equality,timestamp` on `registeredAt == 0`
Using `AgentInfo.registeredAt == 0` as a sentinel for "never registered" is idiomatic and safe. `block.timestamp` is never 0 post-genesis. The comparison is not a business-logic timestamp comparison.

### `timestamp` on `p.deadline <= block.timestamp`
Intentional deadline validation in `Job.initialize`. Miner timestamp manipulation (±15 sec) is acceptable risk for a deadline check.

### `arbitrary-send-erc20` in `SubscriptionHook.onApprove`
`ctx.subscriber` is supplied by the Job contract, which is a trusted caller registered via HookRegistry. The subscriber must have pre-approved SubscriptionHook. This is not an arbitrary transfer — it is an authorized pull pattern gated by the hook+registry access control.

### `divide-before-multiply` in `MilestoneHook.onAccept`
`stageAmt = totalBudget / stages` is intentional integer truncation for per-stage amount computation. Remainder is correctly computed as `totalBudget % stages` and paid on the last stage. This is the standard approach for dust-safe budget splitting.

## Commands used
```bash
cd contracts/evm
PATH="$HOME/.foundry/bin:$PATH" slither src/acp/<Contract>.sol --filter-paths 'lib|out|cache'
```
