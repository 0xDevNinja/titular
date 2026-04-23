# Threat model — M1 (TITU + Treasury + VestingVault + VeTITU + FeeDistributor)

Scope: first on-chain deployment of the Titular protocol to Base Sepolia.
Contracts live in `contracts/evm/src/`.

## 1. Assets protected

| Asset | Custody | Impact if lost |
|-------|---------|----------------|
| TITU total supply (1B) | Initial recipient (Safe multisig) | Irrecoverable dilution; protocol credibility loss |
| Treasury balances (ETH + arbitrary ERC-20) | `Treasury` proxy, owner = Safe | Direct financial loss |
| Vesting allocations | `VestingVault` | Breach of contractual grants; legal exposure |
| Active ve locks (TITU principal + time) | `VeTITU` | Loss of principal + governance weight |
| Fee pools (weekly buckets) | `FeeDistributor` | Mispaid rewards, theft, DoS of claims |

## 2. Trust model

- Treasury **owner = Safe multisig** (placeholder on Sepolia; 2-of-3 on mainnet).
- VestingVault **admin = Safe multisig** (via `DEFAULT_ADMIN_ROLE`). Admin can grant + revoke.
- TITU is ownerless post-deploy. The multisig holds the initial mint but has no contract-level privilege after construction.
- VeTITU is ownerless and non-upgradeable.
- FeeDistributor is ownerless; anyone can checkpoint and anyone can trigger a claim on behalf of a user.
- **Deployer has no residual authority**: all ownership is handed to `MULTISIG` at construction/init.

## 3. Per-contract threats + mitigations

### TITU.sol

| Threat | Mitigation |
|--------|-----------|
| Post-construction mint exploit | No `mint` function; supply is fixed at constructor. `_mint` reachable only via OZ `ERC20Burnable.burn*` which subtracts, not adds. |
| Permit replay across chains | `ERC20Permit` uses domain separator with chainId, name, version. Nonces monotonic per owner. |
| Vote double-counting via fork | `clock()` returns timestamp (ERC-6372 `mode=timestamp`), not block number. Historical checkpoints anchored to timestamps. |
| Delegation griefing | Standard OZ votes; user opts in via `delegate`. No third-party force-delegate. |
| Unsafe casts | None in contract (OZ v5 checkpoint library already uses SafeCast). |

### Treasury.sol

| Threat | Mitigation |
|--------|-----------|
| Stranger drains funds | All mutating functions gated on `onlyOwner` (Safe). |
| Re-init hijack | `_disableInitializers()` in constructor; `initialize` guarded by OZ `initializer` modifier. |
| Native transfer revert in `withdraw` | `(bool,)` return check; reverts with `NativeTransferFailed`. |
| Non-standard ERC-20 (USDT / USDC legacy) | All transfers go through `SafeERC20.safeTransfer`. |
| UUPS malicious upgrade | `_authorizeUpgrade` is `onlyOwner` + zero-address check. Any upgrade requires Safe sig threshold. |
| Reentrancy via `withdraw` native call | No state follows the `.call{value}`; event emit is after, but no balance is re-read. CEI observed. |
| Pause missing | Flagged for M10 audit hardening. Treasury in v1 has **no pause**; trade-off is intentional (Safe multisig is the emergency brake). |

### VestingVault.sol

| Threat | Mitigation |
|--------|-----------|
| Re-grant to same beneficiary bypasses accounting | `GrantExists` revert on any pre-existing `total != 0`. |
| Admin steals grant by re-granting after revoke | `revoke` wipes storage; fresh grant requires fresh `safeTransferFrom`. Not a vulnerability — admin is trust-assumed. |
| Rounding siphons dust | Integer `total * elapsed / duration` floors conservatively; any dust stays as the last release is `total - g.released` cumulative. |
| Beneficiary griefs by DoS on `release(self)` | Anyone can call `release(beneficiary)` — removes grief. |
| Missing cliff check | Enforced: `cliff <= duration`, pre-cliff returns 0. |

Invariant: `sum(grants.total - grants.released) ≤ token.balanceOf(vault)` — enforced by `VestingInvariant.sol` (50 runs × 50 calls, 0 reverts).

### VeTITU.sol

| Threat | Mitigation |
|--------|-----------|
| Early withdrawal of locked principal | `withdraw` requires `block.timestamp >= lock.end`. |
| Decreasing lock length | `increaseUnlockTime` strictly later-only (`end <= old.end` reverts). No shorten path. |
| Lock-end > MAXTIME | Explicit `UnlockTimeTooLong` guard. |
| Re-entrancy via token callbacks (hostile TITU fork?) | TITU whitelist only; `nonReentrant` on all mutators. `_checkpoint` runs BEFORE `safeTransferFrom`. |
| Lock transfer evades decay | `transfer` / `transferFrom` always revert (`NoTransfersAllowed`). |
| `getPastVotes` returns future | Reverts if `timepoint >= block.timestamp` per ERC-5805. |
| `totalSupply` drift post-expiry | Documented approximation — global point decays but doesn't replay per-week slope changes for expired locks. Corrected at every mutation. **Full `slope_changes` replay deferred to M10.** |
| int128 overflow on lock amount | `_toI128` guards against upper bound. |

Invariant: per-user balance ≤ locked principal — enforced by `VeInvariant.sol`.

### FeeDistributor.sol

| Threat | Mitigation |
|--------|-----------|
| Claim of in-progress week (front-run deposits) | `_computeClaim` loops `cursor < currentWeek` — never pays out the current bucket. |
| Reentrancy via malicious reward token | `nonReentrant` on `claim`; `lastTokenBalance` decremented before `safeTransfer` (CEI). |
| Direct donation griefing accounting | `_checkpointToken` uses delta from `lastTokenBalance`; a donation simply increases next-week bucket. No overflow in realistic scenarios. |
| Unbounded claim loop | `claim` loops over user-supplied `tokens[]` (bounded) and weeks since last claim (bounded by elapsed time). An inactive user first claim after long period can be gas-heavy — mitigated by user-driven claim frequency. Not an attacker-controlled unbounded loop. |
| Non-standard ERC-20 return | `SafeERC20` everywhere. |

## 4. Static analysis baseline

- Slither `0.11.5` with `slither.config.json`:
  - **0 high, 0 medium, 0 informational** (informational excluded by config).
  - Remaining Low findings are all acceptable: `block.timestamp` comparisons (intentional; rewards cadence + lock maturity), `calls-loop` on `FeeDistributor.claim` (user-bounded tokens array), one missing-zero-check flagged on `VestingVault.release(address)` (the caller never writes and is allowed to be any address by design).
- Echidna: **NOT installed in local dev env**; wired in CI only (`security.yml`). Local alternative: invariant runs via `forge test --match-path "test/invariant/**"` (50×50).
- Halmos: not yet wired; deferred to post-M2 when FeeDistributor math is exposed to adversarial inputs.

## 5. Deferred to later milestones

| Item | Reason | Target |
|------|--------|--------|
| Etherscan verification of all 5 contracts | No `BASESCAN_API_KEY` provisioned in current env; no funded deployer key | M1 post-merge task |
| Echidna baseline run | Not installed in local env; needs CI image | M2 (CI pipeline) |
| Full `slope_changes` replay in VeTITU `totalSupply()` | Complexity trade-off vs current single-point approximation | M10 audit-hardening |
| Treasury `Pausable` | v1 relies on Safe multisig as incident brake | M10 audit-hardening |
| Two-step ownership transfer | OZ `Ownable` single-step; multisig ceremony substitutes | M10 |

## 6. Invariants enforced

| Invariant | Anchor | Run budget (default profile) |
|-----------|--------|------------------------------|
| `TITU.totalSupply() == 1e27` under transfers | `test/invariant/TituInvariant.sol` | 50 runs × 50 calls |
| Vault balance ≥ outstanding grants | `test/invariant/VestingInvariant.sol` | 50 runs × 50 calls |
| `ve.balanceOf(user) ≤ lockedAmount` | `test/invariant/VeInvariant.sol` | 50 runs × 50 calls |

CI profile `[profile.ci]` bumps fuzz runs to 1024 and invariants to 100×100.

## 7. Deployment hardening

- Deployer key (`DEPLOYER_PRIVATE_KEY`) only used for contract creation; never for privileged calls afterwards.
- `MULTISIG` passed as constructor/init arg; ownership is set at construction. No post-deploy transfer.
- `DeployPhase1.s.sol` is idempotent: rerun produces new addresses and rewrites `deployments/base-sepolia.json`. Safe-signed ceremonies are required to accept the new addresses into registries.
- Sizes (reported by `forge build --sizes`, runtime bytecode):
  - TITU ≈ 11.4 KB
  - Treasury (impl) ≈ 5.5 KB
  - VestingVault ≈ 5.5 KB
  - VeTITU ≈ 8.0 KB
  - FeeDistributor ≈ 4.3 KB

All well under the 24 KB contract-size limit.

## 8. Incident response (v1)

- No on-chain pause in TITU, VestingVault, VeTITU, FeeDistributor. Owner multisig has no pause lever on Treasury either.
- Response path: Safe rotates asset custody via `Treasury.withdraw` to a cold-wallet address.
- Post-incident: emergency patch is a UUPS upgrade on Treasury; VeTITU/VestingVault/FeeDistributor need fork-and-migrate.
- **Action required before mainnet**: add `Pausable` to Treasury and FeeDistributor gated on multisig; codify handoff runbook in `docs/runbooks/incident.md` (M10).
