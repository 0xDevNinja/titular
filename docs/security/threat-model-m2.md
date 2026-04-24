# Threat Model — M2 Launchpad

**Scope**: modular launchpad on Base Sepolia. Covers 7 core contracts + 7 optional
modules, bonding-curve trading against TITU, graduation to Uniswap V2, 10-year LP
lock, 1% transfer tax with 70/30 split.
**Status**: draft, pre-implementation. Living document — updated per PR that alters
entry points, invariants, or trust boundaries.
**Audience**: solidity-principal (implementer), reviewer, external auditor.

---

## 1. Assets under protection

| # | Asset | Custodian | Value-at-risk | Notes |
|---|-------|-----------|---------------|-------|
| A1 | `AgentToken` total supply (1e27 wei) | `BondingCurve` at mint | economic — dilutes all holders if inflatable | supply MUST be fixed at launch |
| A2 | Bonding-curve TITU real reserve | `BondingCurve` | up to 42,000 TITU per agent | drained at graduation; invariants enforce monotone growth pre-grad |
| A3 | LP tokens (Uniswap V2) | `LPLock` | creator + protocol economic value | immutable 10-year lock, non-upgradeable |
| A4 | Creator fees (70% of 1% tax) | `FeeRouter` | continuous stream | pull-based distribution |
| A5 | Treasury fees (30% of 1% tax) | `FeeRouter` → Treasury | continuous stream | same surface as A4 |
| A6 | veTITU holder snapshot (Merkle root) | `AirdropModule` | up to 5% of agent supply per agent | snapshot integrity depends on root provenance |
| A7 | 60-day escrow balance | `SixtyDaysModule` | up to creator-fee accrual window | refundable if creator defaults |
| A8 | Pre-buy vesting allocations | `PreBuyModule` → `TokenVesting` | up to 100% of agent supply | clone storage; admin-key risk |
| A9 | CapitalFormationModule USDC float | module-owned | milestone payouts | TWAP-gated |

---

## 2. Actors

| Actor | Motivation | Capability | Trust level |
|-------|------------|------------|-------------|
| Creator | launch + earn fee stream | calls `launchAgent`, configures modules, receives 70% fee, withdraws LP at T+10y | semi-trusted (can rug if modules allow) |
| Trader | price exposure | calls `buy`/`sell` on curve or pool | untrusted |
| Sniper | extract alpha from launch | same surface as trader; may use MEV infra | hostile |
| Graduator caller | triggers `Graduator.graduate` | any EOA/contract at threshold | untrusted (callable by anyone) |
| Module-gated caller | e.g. Merkle claimant, FDV payout caller | limited to module-specific entry | untrusted |
| veTITU holder | airdrop claim | holds veTITU at snapshot block | untrusted but eligibility-gated |
| MEV bot | sandwich, back-run, atomic arb | controls block ordering on L2 sequencer queue | hostile |
| Malicious module dev | privilege escalation | authored module registered in factory | hostile — primary trust-boundary concern |
| Protocol multisig | pause / upgrade `LaunchpadFactory` | governance | trusted (multisig-enforced) |

---

## 3. Entry points (attack surface)

File paths below are expected layout under `contracts/evm/src/launchpad/`.

| # | Contract.function | External? | Guarded by | Touches asset |
|---|-------------------|-----------|------------|---------------|
| E1 | `LaunchpadFactory.launchAgent` | external | param validation, module allowlist | A1, creates A2 |
| E2 | `BondingCurve.buy` | external | ReentrancyGuard, CEI, slippage | A1, A2 |
| E3 | `BondingCurve.sell` | external | ReentrancyGuard, CEI, slippage | A1, A2 |
| E4 | `FeeRouter.distribute` | external | permissionless, stateless | A4, A5 |
| E5 | `IModule.onLaunch` | external (factory-only) | factory-authenticated caller | depends on module |
| E6 | `IModule.onTrade` | external (curve-only) | curve-authenticated caller | A1, A2 indirectly |
| E7 | `Graduator.graduate` | external (curve-only) | one-shot guard, threshold check | A2 → A3 |
| E8 | `LPLock.withdraw` | external | `block.timestamp >= unlockTime`, recipient == creator | A3 |
| E9 | `AirdropModule.claim` | external | Merkle proof, nullifier bitmap | A6 |
| E10 | `SixtyDaysModule.commit / refund` | external | role (creator for commit, anyone post-T60) | A7 |
| E11 | `CapitalFormationModule.payout` | external | FDV TWAP check | A9 |
| E12 | `PreBuyModule` + `TokenVesting.release` | external | vesting schedule, beneficiary | A8 |

---

## 4. STRIDE per entry point

S = Spoofing, T = Tampering, R = Repudiation, I = Info-disclosure, D = DoS, E =
Elevation of privilege. Each cell is blank when not applicable at the relevant
layer; on-chain info-disclosure is generally N/A because state is public.

### E1. `LaunchpadFactory.launchAgent`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | attacker impersonates creator to launch in their name | launch binds `msg.sender` as creator; off-chain metadata signed if linked to identity |
| T | user-supplied module flags enable hostile combinations | server-side (on-chain) allowlist of registered modules; reject unknown addresses |
| T | name/symbol/URI fields exceed length, break indexer | length caps (64 B name, 16 B symbol, 256 B URI) |
| R | creator disputes module config post-launch | `AgentLaunched` event emits module-flag bitmap + config hash; immutable |
| D | spam launches exhaust factory storage | per-launch fee OR rate-limit via veTITU gate (module-level) |
| E | hostile module registered → steals supply at `onLaunch` | only multisig can register modules; module code audited before allowlist |

### E2 / E3. `BondingCurve.buy` / `sell`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | token callback (ERC-777 / hook) re-enters | TITU + AgentToken are plain ERC-20 by construction; ReentrancyGuard nonetheless |
| T | integer rounding drains reserve (off-by-one loop) | OZ `Math.mulDiv`; invariant `x * y >= k_initial`; Echidna fuzz |
| T | slippage ignored, user gets 0 out | `minOut` parameters required and checked post-swap |
| R | user claims they were front-run; no evidence | events `Bought(...)` / `Sold(...)` with trader, amounts, reserves |
| D | curve griefed via tiny dust swaps | minimum swap size (e.g. 1e12 wei); not-at-block-zero gate from `AntiSniperModule` |
| D | `onTrade` hook in malicious module reverts → trading frozen | allowlist modules; wrap `onTrade` in try/catch with gas limit OR fail-closed by policy (see §5) |
| E | module returns `taxBpsOverride` > 10000 → overflow or full drain | cap override at 9900 bps (99%) in curve; revert on overflow |

### E4. `FeeRouter.distribute`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| T | tax-loop: FeeRouter transfer re-enters AgentToken._update → infinite recursion on tax | `FeeRouter` MUST be excluded from transfer tax (address allowlist in `AgentToken._update`) |
| D | `distribute` gas-grief via malicious recipient contract | push-to-EOA for creator; Treasury is a known contract; fall back to pull if push fails |
| D | dust accumulation prevents split precision | caller pays gas; no minimum-balance gate; rounding dust stays in router |

### E5 / E6. `IModule.onLaunch` / `onTrade`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | arbitrary caller invokes hook pretending to be factory/curve | hook reverts unless `msg.sender == factory` (for `onLaunch`) or `msg.sender == curve` (for `onTrade`) |
| T | module composition: two modules override `onTrade` tax — non-associative result | factory enforces at most one tax-override module per agent; other modules may observe but not override |
| D | long module list exceeds block gas | hard cap N modules per launch (e.g. 6); per-module gas budget on `onTrade` |
| E | malicious module reads `tx.origin` and acts on behalf of trader elsewhere | policy: auditor rejects modules that use `tx.origin`; Slither detector enabled |

### E7. `Graduator.graduate`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | anyone calls before threshold | require real TITU reserve >= 42_000e18; state flag `graduated = true` one-shot |
| T | reentrancy via Uniswap V2 `addLiquidity` callback chain | CEI: set `graduated = true` before external calls; ReentrancyGuard |
| T | attacker pre-creates the UniV2 pair with hostile reserves | call `factory.getPair`; if exists and not empty, revert or use `safeAddLiquidity` that tolerates skewed reserves by quoting both sides |
| D | `createPair` fails for non-standard token (AgentToken has tax) | exclude pair + router from tax allowlist in `AgentToken._update`; covered in unit test |
| D | graduation stuck because `addLiquidity` reverts on dust | pre-compute exact amounts; no slippage on this leg because curve is sole source |
| E | LP minted to attacker instead of `LPLock` | hard-code `LPLock` address at curve init; no mutator |

### E8. `LPLock.withdraw`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | non-creator withdraws | stored `recipient` immutable at deposit; compare `msg.sender == recipient` |
| T | `unlockTime` manipulated via proxy upgrade | `LPLock` non-upgradeable, no owner, no admin |
| T | timestamp manipulation | `block.timestamp` has ≤ 15s drift; 10-year horizon makes this negligible; documented as accepted risk |
| D | LP tokens locked forever if recipient address becomes unreachable | recipient-of-record pattern: allow creator to assign new recipient before unlock via signed message (open question — see §9) |

### E9. `AirdropModule.claim`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| S | claimant submits proof for another address | proof commits to `(claimant, amount)`; enforced in `keccak256` leaf |
| T | replay of same proof twice | nullifier bit per leaf index |
| T | malicious Merkle root (creator allocates 100% to themselves) | root committed pre-launch and emitted in `AgentLaunched`; offchain snapshot script verified + published; audit-time check |
| D | griefing: many zero-value proofs bloat bitmap | set minimum claim amount; amounts are part of leaf so zero-value is impossible if snapshot script rejects them |
| E | root settable by attacker | setter protected by factory-authenticated call; emitted; one-shot |

### E10. `SixtyDaysModule.commit / refund`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| T | creator calls `commit` then disappears | acceptable — this is the success path |
| T | creator calls `refund` after trading to dump their fees on holders | `refund` path drains LP pre-graduation; no refund path post-graduation — document state machine |
| D | after T+60, module griefed via dust refund calls | single one-shot `commit` or `refund` state transition |

### E11. `CapitalFormationModule.payout`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| T | TWAP manipulated by flash-loan pump of UniV2 pair | TWAP window >= 30 min; reject if pool liquidity < threshold |
| T | module holds USDC — upgradable? | non-upgradeable; explicit payout schedule committed at launch |
| D | payout caller grief-reverts transfer to a DoS-prone milestone | pull-to-creator pattern; milestone stays claimable |

### E12. `PreBuyModule` + `TokenVesting`

| STRIDE | Threat | Mitigation |
|--------|--------|------------|
| T | creator changes vesting schedule post-launch | clone storage is initialized once; fields marked immutable or gated by `initializer` modifier |
| T | vesting beneficiary swapped to attacker | beneficiary immutable post-init |

---

## 5. Economic attacks

| # | Attack | Vector | Mitigation |
|---|--------|--------|------------|
| EA1 | Bonding-curve sandwich | MEV bot buys before large victim buy, sells after | slippage bounds enforced; inherent MEV risk documented; L2 sequencer partially mitigates; open risk for public mempool (see §9) |
| EA2 | Graduation front-run | sniper observes threshold-1 TITU reserve and inserts buy of exactly 1 TITU to trigger graduation, then snipes UniV2 pair first block | V2 pair bootstraps liquidity from curve at same exchange rate as last curve trade → minimal arb; still an open risk, mitigate with private mempool submission of graduation tx |
| EA3 | Tax-loop DoS on FeeRouter | `FeeRouter` recipient is AgentToken itself → re-enters `_update` infinitely | `FeeRouter` MUST be in tax-exclusion allowlist; invariant test |
| EA4 | Two modules both override `onTrade` tax | later module's return value overwrites earlier; economics non-deterministic from user viewpoint | factory enforces single tax-override role per agent; bitmap field in `LaunchParams` validated |
| EA5 | UniV2 LP snipe at threshold | first swap on fresh pair buys bottom price | see EA2; additionally, LPLock mints LP to lock immediately, so only swap-side exposure; AntiSniperModule 99→1% decay applies post-grad? NO — post-grad is V2 pair, tax still applies via AgentToken transfer tax, but decay module is curve-scoped. Documented gap. |
| EA6 | Merkle-claim replay | resubmit valid proof | nullifier per leaf index |
| EA7 | Time manipulation on 10y lock | miner/sequencer skews `block.timestamp` | ≤ 15s drift vs 10y horizon = 5e-8 relative; accepted risk |
| EA8 | Virtual-reserve exploit | launch with zero real reserve; first buy at extreme ratio | virtual reserve constants are immutable in curve init; first trade price is deterministic |
| EA9 | Malicious `ExistingTokenModule` wraps tax-on-transfer token | curve math breaks because tokens-in != tokens-received | module uses whitelist of known-good ERC-20 behaviors OR measures actual delta post-transfer |
| EA10 | Creator rug via `AirdropModule` root | creator sets root allocating 100% to self | root commitment published + snapshot script open-source; client-side verification in FE |

---

## 6. Trust boundaries

```
   [Untrusted inputs: trader calldata, URIs, module configs]
             |
             v
  +----------+----------+
  |  LaunchpadFactory   |  <-- MULTISIG owns module allowlist (trust boundary)
  +----------+----------+
             |
   clones    |      registers
             v
  +----------+----------+       +---------------------+
  |    AgentToken       |<----->|     FeeRouter       | <-- allowlist in _update
  +---------------------+       +---------------------+
             |
  +---------------------+
  |    BondingCurve     |  <-- ReentrancyGuard; threshold check triggers Graduator
  +----------+----------+
             |  onTrade (boundary: curve -> modules)
             v
  +---------------------+      +---------------------+
  |  Module set (N)     |      |     Graduator       |
  |  (factory-gated)    |      |  (one-shot)         |
  +---------------------+      +----------+----------+
                                          |
                                          v
                              +-----------+---------+
                              |        LPLock       |
                              |  (non-upgradeable)  |
                              +---------------------+
```

**Critical boundaries**:

1. **Factory ↔ Modules**: factory delegates business logic to modules. A module
   registered in the allowlist can veto or tax trades, hold escrow, and receive
   agent-token supply (pre-buy). Any bug here is privilege escalation.
   Mitigation: allowlist only, multisig-controlled, module code externally
   audited before listing, version-pinned.

2. **Curve ↔ Graduator**: threshold-triggered external call. Classic reentrancy
   surface. Mitigation: state flag set before call; `nonReentrant`; Halmos
   symbolic proof on graduation atomicity.

3. **AirdropModule root provenance**: off-chain snapshot → on-chain root. Bridge
   between trusted data and untrusted callers. Mitigation: open-source snapshot
   tool, public verification, root hash pinned in `AgentLaunched` event.

4. **`AgentToken._update` tax-exclusion allowlist**: must include FeeRouter,
   BondingCurve, Graduator, UniV2 pair, UniV2 router. Any miss causes either
   double-tax or infinite-loop reentry. Covered by invariant tests.

---

## 7. Invariants the implementation MUST preserve

| ID | Invariant | Contract | Test type |
|----|-----------|----------|-----------|
| I1 | `curveK`: `realTITU * realAgent >= k_initial` pre-graduation; buy and sell both preserve | `BondingCurve` | Echidna + Foundry invariant |
| I2 | `lpLocked10y`: after `Graduator.graduate`, `LPLock.unlockTime - depositBlockTime == 10 years` | `LPLock` | unit + invariant |
| I3 | `feeSplit70_30`: every `distribute` call sends exactly 70/30 ± rounding dust to creator/treasury | `FeeRouter` | unit + fuzz |
| I4 | `totalAgentSupply`: `totalSupply() == TOTAL_SUPPLY` forever after init | `AgentToken` | invariant |
| I5 | `graduation is one-shot`: `Graduator.graduate` cannot be called twice for same agent | `Graduator` | unit + Halmos |
| I6 | `taxBpsOverride <= 9900` | `BondingCurve` | unit |
| I7 | `moduleCount <= MAX_MODULES` per agent | `LaunchpadFactory` | unit |
| I8 | Merkle claim nullifier monotonic — once set, never cleared | `AirdropModule` | invariant |

---

## 8. Mitigations — cross-cutting

- **ReentrancyGuard** on every state-mutating external entry: `buy`, `sell`,
  `graduate`, `distribute` (defence in depth even with pull-payment), `claim`.
- **Checks-effects-interactions** strictly on curve and graduator.
- **Pull-over-push** where recipient is external contract (Treasury, creator
  SCA); push to EOA is acceptable if recipient code-size is 0.
- **Slippage params** mandatory (`minOut`) on all curve and V2 swaps.
- **Module allowlist**: only multisig-signed `registerModule` call can add a new
  module address to `LaunchpadFactory.moduleRegistry`. Emission of
  `ModuleRegistered` event is auditable.
- **`taxBpsOverride` cap ≤ 9900 bps** enforced in `BondingCurve` regardless of
  module return value.
- **Merkle root commitment pre-launch**: root set atomically inside
  `LaunchpadFactory.launchAgent` via `AirdropModule.onLaunch`; emitted in
  `AgentLaunched`; never mutated post-launch.
- **Pausable**: `LaunchpadFactory` has `Pausable` gated by multisig for new
  launches. Existing agents continue trading (curve is not pausable to avoid
  soft-rug). Graduator is pausable to halt L1-incident response.
- **Upgrade surface**: factory is NOT upgradeable in M2; clones are minimal
  proxies pointing at immutable implementations. Any future UUPS adoption
  requires re-review and ADR.
- **Prompt-injection surfaces**: agent `name`, `symbol`, `soul URI`, and module
  `data` bytes may later be consumed by natural-language-model systems
  (agent runtime, discovery). Policy: cap lengths, sanitize control characters,
  wrap any model prompt that injects these fields in `<<UNTRUSTED>>...<</UNTRUSTED>>`
  sentinels at the consumer tier (not on-chain).

---

## 9. Open risks (explicitly deferred)

| ID | Risk | Why deferred | Mitigation plan |
|----|------|--------------|-----------------|
| O1 | MEV sandwich / graduation front-run on public L2 | Base sequencer is a single private-ish lane; acceptable for testnet | integrate MEV-protected RPC or Flashbots equivalent on Base pre-mainnet |
| O2 | Creator key loss → LP unreachable at T+10y | product/UX question, not contract-level | add `LPLock.setRecipient` gated by 2-of-3 between creator-pair + DAO; ADR required |
| O3 | UniV2 pair front-creation grief | factory.getPair check covers detection; handling of pre-existing pair needs design | design doc in follow-up milestone audit-prep |
| O4 | Modules fail-open vs fail-closed on revert in `onTrade` | today: fail-closed (trade reverts) | if a module bug freezes trading, multisig pause path exists; longer-term: per-module circuit breaker |
| O5 | `CapitalFormationModule` TWAP fragility on thin pools | TWAP window + liquidity floor planned | stress test pre-mainnet with low-liquidity forks |
| O6 | `ExistingTokenModule` wrapping of tax-on-transfer or rebasing tokens | explicit whitelist planned | maintained allowlist, rejected-by-default |
| O7 | Post-graduation AntiSniper inactivity | module is curve-scoped; post-V2 behaviour by design | documented; alternative: attach a tax-decay hook to AgentToken for N blocks post-grad |
| O8 | `PreBuyModule` 100% supply vest → creator rug | product policy: cap pre-buy at (e.g.) 20% supply + min vest duration | enforce in module config validation; ADR required |

---

## 10. Review gates

- Slither clean (`contracts/evm/slither.config.json`) — no HIGH or MEDIUM findings
  unacknowledged.
- Echidna invariants green on I1, I2, I3, I4 at `--test-limit 500000`.
- Halmos symbolic proof on `Graduator.graduate` one-shot + LP recipient.
- External audit before mainnet (blocking).
- This document re-reviewed on any PR that:
  - adds a new module,
  - changes `AgentToken._update` tax logic or allowlist,
  - changes `BondingCurve` math or thresholds,
  - touches `Graduator` or `LPLock`.

## 11. Expected file layout (for cross-reference)

```
contracts/evm/src/launchpad/
  AgentToken.sol
  BondingCurve.sol
  LaunchpadFactory.sol
  Graduator.sol
  LPLock.sol
  FeeRouter.sol
  modules/
    IModule.sol
    AntiSniperModule.sol
    SixtyDaysModule.sol
    CapitalFormationModule.sol
    AirdropModule.sol
    LaunchRadarModule.sol
    PreBuyModule.sol
    ExistingTokenModule.sol
contracts/evm/test/
  invariant/LaunchpadInvariant.sol
  integration/{Launch,Trade,Graduate,Modules}.t.sol
```

---

*End of M2 threat model. Changes to this file require security-lead review.*
