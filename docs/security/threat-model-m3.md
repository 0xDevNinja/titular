# Threat Model — M3 ACP v2 Contracts

**Scope**: Agent Connection Protocol v2 on Base Sepolia. Covers the four
core ACP contracts (`AgentRegistry`, `Escrow`, `ReputationOracle`,
`ACPRouter`) plus the supporting `Job` / `JobFactory` / `FeeSplitter` /
`BuybackBurner` / `HookRegistry` and the registered hook implementations
that compose around the job lifecycle.
**Status**: draft, pre-implementation. Living document — updated per PR
that alters entry points, invariants, or trust boundaries.
**Audience**: solidity-principal (implementer), reviewer, external
auditor, indexer / SDK consumers downstream.

---

## 1. Assets under protection

| #  | Asset | Custodian | Value-at-risk | Notes |
|----|-------|-----------|---------------|-------|
| A1 | Per-job escrow balance (USDC / TITU / agent token) | `Escrow` clone (or `Job` contract) | full job principal + tip per active job | refund on dispute / expiry; release on `submitResult` + grace |
| A2 | Agent identity record (controller, metadata URI, capabilities) | `AgentRegistry` | reputational; downstream routing decisions | controller two-step transfer; metadata hash committed on update |
| A3 | Reputation score ledger (per agent, per dimension) | `ReputationOracle` | reputational; affects routing + caller selection | append-only; signed by whitelisted scorers; nonce-protected |
| A4 | Scorer whitelist + arbiter set | `ReputationOracle`, `Escrow` | governance — controls who can score / resolve | multisig-managed; events on add/remove |
| A5 | Protocol fee accrual (per-job basis points) | `FeeSplitter` | continuous stream | pull distribution to treasury / buyback / creator |
| A6 | Buyback / burn TITU flow | `BuybackBurner` | economic — supports TITU price floor | TWAP-gated swap; slippage guard |
| A7 | Hook registry allowlist | `HookRegistry` | privilege escalation surface | multisig-curated; only listed hook addresses can attach to a job |
| A8 | Permit2 / EIP-2612 typed-data signatures | caller wallet → `Escrow` / `Job` | one-shot fund pull authorisation | nonce + deadline + domain separator enforced on-chain |
| A9 | EIP-712 score signatures from scorers | scorer EOA → `ReputationOracle` | replayability of score posts | per-(scorer, agent, jobId, nonce) replay guard |

---

## 2. Actors

| Actor | Motivation | Capability | Trust level |
|-------|------------|------------|-------------|
| Caller | request work, fund escrow, release / dispute | calls `JobFactory.createJob`, `Job.fund`, `Job.release`, `Job.dispute` | untrusted |
| Agent (controller) | accept job, deliver result, claim payout | calls `AgentRegistry.register/update`, `Job.submitResult` | semi-trusted (rep-bonded) |
| Scorer | post EIP-712 signed reputation deltas | signs scores; transactor may be anyone (meta-tx) | semi-trusted (whitelisted) |
| Arbiter | resolve disputes during dispute window | calls `Escrow.resolveDispute(jobId, payoutBps)` | trusted (multisig) |
| Hook (registered) | run side-effects on lifecycle transitions | called by `Job` at fund / submit / release / dispute | semi-trusted (allowlisted, audited) |
| ACPRouter caller (read-only) | resolve agent runtime endpoint | view-only `resolve(agentId)` | untrusted, no state mutation |
| Indexer | mirror events to GraphQL store | read-only event stream | untrusted with respect to chain state |
| Front-runner / MEV bot | extract value on fund / release / dispute order | controls block ordering on L2 | hostile |
| Malicious hook author | exfiltrate funds, brick lifecycle | code registered in HookRegistry | hostile — primary trust-boundary concern |
| Protocol multisig | manage scorer set, arbiter set, hook registry, fee config | governance | trusted (multisig-enforced) |

---

## 3. Entry points (attack surface)

File paths below are the expected layout under `contracts/evm/src/acp/`.

| #   | Contract.function | External? | Guarded by | Touches asset |
|-----|-------------------|-----------|------------|---------------|
| E1  | `AgentRegistry.register` | external | `msg.sender` becomes controller; one-shot per agent ID | A2 |
| E2  | `AgentRegistry.updateMetadata` | external | controller-only; emits hash | A2 |
| E3  | `AgentRegistry.transferController` (two-step) | external | controller-only initiates; recipient accepts | A2 |
| E4  | `JobFactory.createJob(params, hooks[])` | external | param validation, hook allowlist, deadline floor | A1, A7 |
| E5  | `Job.fund(amount, sigOrPermit2)` | external | nonReentrant, CEI, sig verification, deadline | A1, A8 |
| E6  | `Job.submitResult(resultHash, uri)` | external | agent-only, state == funded | A1 |
| E7  | `Job.release()` | external | grace expired OR caller explicit ack; nonReentrant | A1, A5 |
| E8  | `Job.dispute()` | external | caller-only OR agent-only (if no submit by deadline); within window | A1 |
| E9  | `Escrow.resolveDispute(jobId, payoutBps)` | external | arbiter-only; bps ≤ 10000; one-shot | A1, A5 |
| E10 | `Job.expire()` | external | permissionless after deadline; one-shot | A1 |
| E11 | `ReputationOracle.postScore(scoreSig)` | external | EIP-712 verify, scorer ∈ whitelist, nonce match | A3, A9 |
| E12 | `ReputationOracle.{addScorer,removeScorer}` | external | onlyOwner (multisig) | A4 |
| E13 | `ACPRouter.resolve(agentId)` | external view | none — read-only | none (read) |
| E14 | `IHook.onFund / onSubmit / onRelease / onDispute / onExpire` | external | `msg.sender == job` (job-authenticated caller) | depends on hook |
| E15 | `HookRegistry.{register,deregister}` | external | onlyOwner (multisig) | A7 |
| E16 | `FeeSplitter.distribute()` | external | permissionless, stateless split | A5 |
| E17 | `BuybackBurner.executeBuyback()` | external | TWAP staleness check, slippage min-out, pause flag | A6 |

---

## 4. STRIDE per entry point

S = Spoofing, T = Tampering, R = Repudiation, I = Info-disclosure, D =
DoS, E = Elevation of privilege. On-chain info-disclosure is generally
N/A because state is public; the I rows below address off-chain leaks
(metadata URIs, EIP-712 typed data) where relevant.

### E1 / E2 / E3. `AgentRegistry`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | attacker registers an agent ID claiming a known controller's identity | first-writer-wins on agent ID; controller is `msg.sender` at register; identity binding is off-chain (signed metadata) |
| T | metadata URI swapped post-register to a hostile payload | `updateMetadata` requires controller; emits `AgentMetadataUpdated(agentId, metadataHash)`; downstream consumers MUST pin `metadataHash` |
| R | controller disputes a transfer they did not initiate | two-step: `transferController(newCtl)` then `acceptControl()` from `newCtl`; both legs emit; immutable event log |
| D | spam registrations exhaust storage / inflate ID-space | per-register fee in TITU OR rate-limit via veTITU gate; ID-space cap documented |
| E | hostile controller-transfer to address(0) bricks agent | reject `newController == address(0)` and self-transfer; pre-empts accidental rug |

### E4. `JobFactory.createJob`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | impostor creates a job referencing a victim caller | caller binds `msg.sender`; job storage records `caller`, `agentId`, `params`, `createdAt` |
| T | hook list contains an unregistered or hostile address | factory MUST iterate `hooks[]` and check `HookRegistry.isRegistered(hook)`; revert on miss |
| T | unbounded `hooks[]` length griefs gas / blocks creation | hard cap `MAX_HOOKS_PER_JOB` (e.g. 8); enforced in factory |
| T | params encode invalid deadline / dispute window (zero, MAX_UINT, past) | param validation: `deadline > now + minDeadline`, `disputeWindow >= MIN_WINDOW`, `disputeWindow <= MAX_WINDOW` |
| R | caller disputes job config post-creation | `JobCreated(jobId, caller, agent, paramsHash, hooksHash)` event commits the full config |
| D | deterministic-clone collision used to grief deploy | use `CREATE2` with `keccak256(caller, agentId, salt)`; revert on collision |

### E5. `Job.fund`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | non-caller funds a job to skew accounting | caller-only OR open funding with `msg.sender` recorded as funder; refunds always go back to recorded funder |
| T | Permit2 / EIP-2612 signature replay | enforce `PermitWitnessTransferFrom` with per-job witness OR check Permit2 nonce after pull; deadline ≤ block.timestamp |
| T | tax-on-transfer or rebasing token drains escrow over time | reject non-allowlisted tokens at factory; for allowlisted set, measure delta post-transfer |
| T | reentrancy via ERC-777 / hook callback during fund | nonReentrant; CEI ordering: state set → external pull → hook call (or hook call last) |
| R | caller claims double-fund | `JobFunded(jobId, funder, token, amount)` event fires only after balance confirmed |
| D | griefer funds with dust to flip state if guard is `>0` | minimum-fund threshold per (token, params); reject otherwise |
| E | malicious token's `transferFrom` returns false silently | use `SafeERC20.safeTransferFrom`; revert on false |

### E6. `Job.submitResult`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | non-agent submits a result | enforce `msg.sender == AgentRegistry.controllerOf(agentId)` at call time |
| T | agent submits placeholder hash to start grace clock and abandon | grace clock starts at `submitResult`; caller has full window to dispute; eventual `release` is permissionless after grace |
| T | hook's `onSubmit` reverts and freezes the job | hooks are bounded (try / catch with gas budget), or fail-closed by policy with multisig pause path |
| R | duplicate submit | `state == funded` precondition; transition to `submitted` is one-shot |
| D | submit after deadline still credited | reject if `block.timestamp > deadline` |

### E7. `Job.release`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | unauthorized caller forces premature release | `release` is permissionless ONLY after `submittedAt + graceWindow` OR caller-explicit ack; otherwise reverts |
| T | reentrancy via agent payout transferring to malicious receiver | nonReentrant; CEI: state → effects (fee split, balance zero) → external transfers; receivers known to be EOAs or audited contracts |
| T | rounding loss on fee split sends 0 to one side | use OZ `Math.mulDiv`; assert sum == amount post-split |
| R | dispute-after-release path | `release` flips state to `released`; `dispute` requires `state ∈ {funded, submitted}` |
| D | hook `onRelease` reverts | bounded gas budget; pause path; documented fail-closed |

### E8. `Job.dispute`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | third party disputes someone else's job | restrict to `msg.sender == caller` OR `msg.sender == agentController` with documented preconditions |
| T | dispute opened repeatedly to grief lifecycle | one-shot transition to `disputed`; subsequent dispute reverts |
| T | dispute opened after window closes | enforce `block.timestamp <= submittedAt + disputeWindow` |
| D | griefing dispute right before grace expiry | dispute window MUST be a strict subset of grace window; arbiter SLA documented |

### E9. `Escrow.resolveDispute`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | non-arbiter resolves | `onlyArbiter` modifier; arbiter set is multisig-managed |
| T | arbiter sends > 100% to one side via bps overflow | enforce `payoutBps + refundBps == 10000`; revert on mismatch |
| T | arbiter resolves twice with different bps | one-shot — `state == disputed` precondition; transition to `resolved` is monotonic |
| R | controversial resolution | `DisputeResolved(jobId, arbiter, payoutBps, reasonHash)` event; reasonHash points off-chain to signed memo |
| E | arbiter quorum compromised | arbiter is a multisig contract; per-resolution requires multisig threshold |

### E10. `Job.expire`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | griefer expires a job mid-submission | enforce `block.timestamp > deadline AND state == funded`; not callable once `submitted` |
| T | refund sent to wrong recipient | refund destination is the recorded funder (immutable post-fund) |
| T | reentrancy via refund hook | nonReentrant; CEI |
| D | spam expire calls | one-shot transition; subsequent reverts |

### E11. `ReputationOracle.postScore`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | impostor signs as scorer | EIP-712 typed-data domain pins `chainId`, `verifyingContract`, `name`, `version`; signer recovered and checked against whitelist |
| T | replay of a valid score signature | typed struct includes `(scorer, agentId, jobId, nonce, deadline)`; oracle stores `consumedNonces[scorer]`; reject if seen |
| T | scorer signs scores for a different chain (cross-chain replay) | domain `chainId` enforced; reject foreign-chain sigs |
| T | reputation grinding via Sybil agents | scoring is whitelisted; agent ID is one-time per controller; per-job-id binding caps double-counting; off-chain Sybil resistance is out-of-scope (see §9) |
| R | scorer disputes a score they did not sign | `ScorePosted(agentId, scorer, jobId, scoreVector, nonce)` event commits the full payload |
| D | scorer floods small scores to bloat per-agent score history | per-scorer rate limit (one score per `(scorer, jobId)`); pruning policy documented |
| E | new scorer added without governance | `addScorer` is `onlyOwner` (multisig); emits `ScorerAdded` |

### E14. `IHook.{onFund, onSubmit, onRelease, onDispute, onExpire}`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | arbitrary caller invokes hook pretending to be Job | hook reverts unless `msg.sender == job`, and `Job.isCanonical(msg.sender)` (factory-clone bytecode hash check) |
| T | hook composition order matters: two hooks both mutate fee split | factory enforces hook ordering policy; only the LAST hook's bps override is honoured, or fee-mutating hooks are mutex-tagged in registry |
| T | hook invokes `Job.release` synchronously inside `onSubmit` | nonReentrant on Job covers this; documented as out-of-policy |
| D | long hook list exceeds block gas | `MAX_HOOKS_PER_JOB`; per-hook gas budget (`gasleft()`-bounded sub-call) |
| E | hook reads `tx.origin` and impersonates caller elsewhere | policy: HookRegistry rejects hooks that reference `tx.origin`; Slither detector enabled |
| E | hook re-enters Job to flip state mid-lifecycle | nonReentrant on every Job mutator |

### E15. `HookRegistry.{register,deregister}`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | impostor registers a hook | `onlyOwner` (multisig) |
| T | deregister of an in-flight hook leaves jobs orphaned | hooks evaluated at job-create time and snapshotted into the job; deregister affects only NEW jobs |
| R | which hooks were attached to a job | `JobCreated` event emits `hooksHash`; hook list is recoverable from `JobFactory.hooksOf(jobId)` |
| E | hostile hook code listed | external audit + version-pin of hook bytecode hash; `HookRegistered(hook, codeHash)` event commits hash |

### E16. `FeeSplitter.distribute`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| T | `FeeSplitter` becomes a tax-loop endpoint if a hook routes back | FeeSplitter targets are multisig-curated; circular routing fails the same allowlist check used by `AgentToken._update` for tax-exempt set |
| D | distribute griefed by a recipient that reverts | push-to-EOA where possible; pull pattern with per-recipient claim for contracts |

### E17. `BuybackBurner.executeBuyback`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| T | TWAP manipulated via flash-loan pump | TWAP window ≥ 30 min; reject if pool liquidity < threshold; sanity bound vs spot |
| T | sandwich on the buyback swap | `minOut` derived from TWAP - slippage budget; revert otherwise |
| D | buyback called repeatedly to drain treasury via slippage | per-epoch cap on burn budget; multisig pause |
| E | swap router pinned to a hostile address | router immutable, set at deploy |

### E13. `ACPRouter.resolve` (read-only)

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| I | resolver returns a stale or hostile endpoint | endpoint is the metadata URI committed by controller; consumers MUST treat as untrusted and validate at higher tier |
| D | view function gas-griefs an integrator | view is O(1) read; bounded |

---

## 5. Economic and protocol attacks

| #   | Attack | Vector | Mitigation |
|-----|--------|--------|------------|
| EA1 | Reputation grinding | caller + agent collude across N jobs to inflate score | scorer whitelist + per-(scorer, jobId) nonce; min-fee-per-job raises grinding cost; off-chain Sybil scoring deferred (see §9) |
| EA2 | Escrow re-entry on payout | malicious agent payout token re-enters Job | token allowlist (no ERC-777); nonReentrant; CEI |
| EA3 | EIP-712 score replay across chains | scorer signs once, tx posted on Sepolia + mainnet | domain separator pins `chainId`; per-scorer nonce |
| EA4 | Permit2 signature replay across jobs | victim signs `permit(funder→escrow, amount, deadline)`; attacker reuses for a second job | use Permit2's `PermitWitnessTransferFrom` with `witness = keccak256(jobId)`; OR consume Permit2 nonce explicitly |
| EA5 | Hook composition order reorder | two hooks both override fee bps; reorder changes payout | factory commits `hooksHash` at create; reorder is impossible post-create; mutex tag in registry on fee-override hook category |
| EA6 | Escrow drain via callback | hook's `onRelease` calls `Job.dispute` to roll state back after payout | Job's `release` flips state to `released` BEFORE calling external; subsequent `dispute` reverts on state precondition; nonReentrant covers re-entry too |
| EA7 | Front-run release between submit + grace | observer sees `submitResult` and races to call `release` before caller can `dispute` | dispute window MUST be a strict subset of grace window; release is gated on `submittedAt + graceWindow <= now` |
| EA8 | Front-run dispute on dispute-window edge | caller submits dispute at last block; observer races a release tx at block boundary | both branches gated on monotone time; `release` requires state ∈ {funded post-submit-and-grace}; once `state == disputed`, release reverts |
| EA9 | Arbiter griefing — never resolves | arbiter multisig idle; funds locked indefinitely | per-job arbiter timeout — fallback split (e.g. 50/50) callable by anyone after `disputedAt + maxArbitrationWindow` |
| EA10 | TWAP pump on `BuybackBurner` thin pool | flash-pump TITU price, trigger buyback, dump | TWAP window ≥ 30 min; liquidity floor; epoch-budget cap |
| EA11 | Fee bps inflation by hostile hook | hook sets `feeBps = 10000` on payout | cap `feeBps <= MAX_FEE_BPS` (e.g. 1500) in Job regardless of hook return |
| EA12 | Hook bytecode swap via SELFDESTRUCT then re-create at same address | post-Cancun, SELFDESTRUCT no longer redeploys, but `CREATE2` collisions remain | HookRegistry pins `extcodehash(hook)` at register; `isRegistered(hook)` re-checks codeHash on attach |
| EA13 | Caller refunds via expire while result is in-flight | network latency lets caller `expire` between submit and on-chain confirmation | `expire` requires `state == funded` AND `now > deadline`; once `submitted`, expire reverts |
| EA14 | Agent metadata-URI poisoning | hostile URI serves malware to indexer / SDK | indexer / SDK MUST treat URI content as untrusted; on-chain commit is `metadataHash` only |
| EA15 | Time skew on dispute window | sequencer drift extends or shortens window | drift ≤ 15s vs window ≥ 1 day = ≤ 2e-4 relative; accepted risk |
| EA16 | Malicious arbiter steals via 100/0 split | compromised arbiter signs full payout to themselves | arbiter is multisig — single key compromise insufficient; escrow MUST send to caller / agent only, never to arbiter |
| EA17 | DOS via dust dispute fee | griefer sends 1-wei `dispute()` calls every block | dispute is one-shot per job; subsequent reverts |

---

## 6. Trust boundaries

```
   [Untrusted: caller, agent, arbiter quorum, hooks code, scorer EOAs]
                      |
                      v
        +-------------+-------------+
        |       JobFactory          |  <-- MULTISIG owns hook registry
        +-------------+-------------+
                      |
            CREATE2   |   snapshots hooks[] + paramsHash
                      v
        +-------------+-------------+        +-----------------------+
        |          Job              |<------>|     Escrow            |
        |  (one per (caller,agent)) |        |  (per-job balance)    |
        +-------------+-------------+        +-----------------------+
                      |
            calls     |  onFund / onSubmit / onRelease / onDispute / onExpire
                      v
        +-------------+-------------+        +-----------------------+
        |     IHook (N, bounded)    |        |  ReputationOracle     |
        |  (registry-allowlisted)   |        |  (EIP-712 + nonces)   |
        +---------------------------+        +-----------+-----------+
                                                         |
                                              +----------+----------+
                                              |   AgentRegistry      |
                                              |   (controller addr)  |
                                              +----------------------+
        +---------------------------+        +-----------------------+
        |       FeeSplitter         |        |     BuybackBurner     |
        |  (multisig recipient set) |        |  (TWAP, slippage,     |
        +---------------------------+        |   pause, epoch cap)   |
                                              +-----------------------+
```

**Critical boundaries**:

1. **Job ↔ Hook**: Job delegates side-effects to allowlisted hooks. A
   hook can read job state, mutate fee splits, and fire follow-on
   transfers. Any bug here is privilege escalation.
   Mitigation: `HookRegistry` is multisig-only, codeHash-pinned at
   registration, and re-validated on attach. Hook entrypoints accept
   only canonical Job clones. NonReentrant on every Job mutator.

2. **Caller ↔ Permit2 / EIP-2612 fund pull**: signature is the
   caller's only authorisation for fund movement. Mitigation:
   per-job witness OR explicit Permit2 nonce consumption; deadline
   bound to job createdAt window.

3. **Scorer ↔ ReputationOracle**: signed messages constitute the only
   trust input for reputation. Mitigation: EIP-712 domain pins
   `chainId` and `verifyingContract`; per-scorer nonce; whitelist
   managed by multisig.

4. **Arbiter ↔ Escrow**: arbiter has unilateral payout authority on
   disputed jobs. Mitigation: arbiter set is multisig; per-job arbiter
   timeout fallback; arbiter can only direct funds to caller / agent,
   not to arbiter or third parties.

5. **AgentRegistry controller-binding**: downstream contracts derive
   "is this agent" by reading `controllerOf(agentId)`. Mitigation:
   two-step transfer; events on every change; non-zero recipient
   required.

---

## 7. Invariants the implementation MUST preserve

| ID | Invariant | Contract | Test type |
|----|-----------|----------|-----------|
| I1 | `funded → exactly one of {released, refunded-via-resolve, refunded-via-expire}` | `Job` / `Escrow` | Halmos symbolic — issue #70 |
| I2 | `state` is monotonic: `created → funded → submitted → (released ∣ disputed → resolved ∣ expired)` | `Job` | Echidna + Foundry invariant |
| I3 | `escrow.balanceOf(jobId) == sum(funded) - sum(released) - sum(refunded)` | `Escrow` | Echidna |
| I4 | `controllerOf(agentId)` is monotonically owned: only current controller can transfer; recipient must accept | `AgentRegistry` | unit + invariant |
| I5 | `consumedNonces[scorer][nonce]` set ⇒ never cleared (no replay) | `ReputationOracle` | invariant |
| I6 | `feeBps + payoutBps == 10000` on every release | `Job` / `FeeSplitter` | unit + fuzz |
| I7 | `hookCount(jobId) <= MAX_HOOKS_PER_JOB` | `JobFactory` | unit |
| I8 | every hook attached to a job is `HookRegistry.isRegistered(hook) == true` AND its `extcodehash` matches the registered hash at attach time | `JobFactory`, `HookRegistry` | unit + invariant |
| I9 | `arbiter` of `Escrow` has `payoutBps + refundBps == 10000` and recipients are subset of `{caller, agent}` | `Escrow` | unit + Halmos |
| I10 | `Job.release` reverts when `block.timestamp < submittedAt + graceWindow` AND no caller ack | `Job` | unit + Halmos |
| I11 | `Job.expire` reverts when `state != funded` OR `block.timestamp <= deadline` | `Job` | unit + Halmos |

---

## 8. Mitigations — cross-cutting

- **ReentrancyGuard** on every state-mutating Job entry: `fund`,
  `submitResult`, `release`, `dispute`, `expire`, and on
  `Escrow.resolveDispute`. Defence in depth even with strict CEI.
- **Checks-effects-interactions** strictly on Job, Escrow, and
  ReputationOracle.
- **Pull-over-push** for fee distribution where the recipient is an
  external contract; push to EOAs only when codeSize == 0.
- **Token allowlist**: escrow tokens are restricted to the
  multisig-curated set. ERC-777 and rebasing tokens explicitly
  rejected. Tax-on-transfer tokens require post-transfer delta
  measurement.
- **Permit2 / EIP-2612 binding**: per-job witness or explicit nonce
  consumption to defeat cross-job replay. Domain separator includes
  `chainId` and `verifyingContract`.
- **EIP-712 score signatures**: `(scorer, agentId, jobId, nonce,
  deadline)` typed struct; domain pins `chainId`. Per-scorer
  `consumedNonces` mapping. `deadline` enforced.
- **Hook registry**: only multisig-signed `register(hook)` adds an
  allowlisted hook address. Registration commits `extcodehash(hook)`;
  attach-time re-validates. `HookRegistered(hook, codeHash)` and
  `HookDeregistered(hook)` events.
- **Hook ordering snapshot**: `JobCreated` event commits
  `hooksHash = keccak256(abi.encodePacked(hooks))`. Re-ordering or
  swapping hooks post-create is impossible.
- **Hook gas budget**: per-hook sub-call uses `gasleft() / N`-bounded
  sub-call to prevent a single hook from monopolising block gas.
  Reverts in hooks fail-closed by policy with multisig pause path.
- **Custom errors** instead of revert strings; descriptive but no
  internal-state leakage.
- **Events on every state change**: `JobCreated`, `JobFunded`,
  `ResultSubmitted`, `JobReleased`, `DisputeOpened`,
  `DisputeResolved`, `JobExpired`, `ScorePosted`, `ScorerAdded`,
  `ScorerRemoved`, `HookRegistered`, `HookDeregistered`,
  `AgentRegistered`, `AgentMetadataUpdated`, `ControllerTransferStarted`,
  `ControllerTransferAccepted`, `BuybackExecuted`, `FeesDistributed`.
- **Two-step controller / arbiter transfers**: prevents accidental
  rug to address(0) or typo.
- **Pausable**: `JobFactory` is `Pausable` gated by multisig — pauses
  NEW job creation. Existing jobs continue (caller / agent funds must
  remain spendable). `BuybackBurner` is independently pausable.
- **Upgrade surface**: factory is NOT upgradeable in M3; Job clones
  are minimal proxies pointing at immutable implementations. Hooks
  themselves are not upgradeable; replacing a hook requires new
  registration with a fresh address.
- **Arbiter timeout fallback**: jobs in `disputed` state past
  `maxArbitrationWindow` resolve to a documented default (e.g. 50/50
  refund) callable by anyone — defeats arbiter griefing.
- **Prompt-injection surfaces**: agent `metadataURI`,
  `capabilities`, and job `params` bytes may be consumed by language
  models in the agent runtime / discovery tier. Policy: cap lengths,
  sanitise control characters, wrap in `<<UNTRUSTED>>...<</UNTRUSTED>>`
  sentinels at the consumer tier. On-chain commits remain hashes.

---

## 9. Open risks (explicitly deferred)

| ID | Risk | Why deferred | Mitigation plan |
|----|------|--------------|-----------------|
| O1 | Sybil-resistance on agent registration | identity binding is off-chain; on-chain enforcement requires KYC-style attestation | integrate veTITU stake-to-register OR World ID-style attestation in M5+ |
| O2 | MEV on `release` / `dispute` ordering at window edges | Base sequencer is private-ish; acceptable for testnet | MEV-protected RPC / Flashbots-equivalent bundle pre-mainnet |
| O3 | Arbiter compromise with full multisig | single-point governance risk | per-job arbiter timeout + fallback split; longer-term: opt-in arbiter set per job class |
| O4 | Hook fail-open vs fail-closed policy | today: fail-closed (lifecycle reverts) | per-hook circuit breaker + multisig per-hook deregister |
| O5 | Reputation oracle scorer collusion | whitelist is small, multisig-managed | post-mainnet: stake-bonded scorers with slashing; tracked in M11+ |
| O6 | Cross-chain agent identity (same `agentId` across L1/L2/Solana) | M9 milestone bridges identity; out-of-scope here | M9 ADR will define canonical-source rule |
| O7 | TWAP fragility on thin pools for `BuybackBurner` | TWAP window + liquidity floor planned | stress test pre-mainnet; pause path |
| O8 | Permit2 vs EIP-2612 token-coverage gap | not all allowlisted tokens implement Permit2 | dual-path fund: Permit2 for compatible tokens, allowance + `safeTransferFrom` fallback |
| O9 | Job state machine extension (sub-tasks, partial-release) | single-shot job in M3; partial-release deferred | M5+ adds `Job.partialRelease(bps)` with re-anchored grace |
| O10 | Indexer divergence — events vs state | indexer is a separate service (M4) | reconciliation job + `eth_getLogs` re-scan path |

---

## 10. Review gates

- Slither clean on `contracts/evm/src/acp/` — no HIGH or MEDIUM
  findings unacknowledged. (Issue #71, M3_SECURITY_REPORT.md.)
- Echidna invariants green on I2, I3, I5, I6 — config under
  `contracts/evm/test/echidna/acp.yaml`. (Issue #71.)
- Halmos symbolic proof on I1 (`funded → exactly one of {completed,
  rejected, expired}`) and I9, I10, I11. (Issue #70,
  `test/halmos/JobLifecycle.halmos.t.sol`.)
- External audit before mainnet (blocking, M10).
- This document re-reviewed on any PR that:
  - adds or alters a hook category in `HookRegistry`,
  - changes `Job` state machine or grace / dispute window semantics,
  - changes `ReputationOracle` signing struct or domain,
  - changes `Escrow.resolveDispute` arbiter semantics,
  - changes the allowlisted token / scorer set governance path.

---

## 11. Expected file layout (for cross-reference)

```
contracts/evm/src/acp/
  AgentRegistry.sol
  ACPRouter.sol
  Job.sol
  JobFactory.sol
  Escrow.sol
  ReputationOracle.sol
  FeeSplitter.sol
  BuybackBurner.sol
  hooks/
    IHook.sol
    HookRegistry.sol
    <category>Hook.sol  (multiple, registered individually)
contracts/evm/test/
  acp/                     # Foundry unit + integration (solidity-principal)
  halmos/                  # symbolic proofs (security-lead, M3)
    JobLifecycle.halmos.t.sol
  echidna/                 # fuzz configs only
    acp.yaml
contracts/evm/audit/
  M3_SECURITY_REPORT.md    # Slither + Echidna + Halmos pass (security-lead)
docs/security/
  threat-model-m3.md       # this document
```

---

*End of M3 threat model. Changes to this file require security-lead
review.*
