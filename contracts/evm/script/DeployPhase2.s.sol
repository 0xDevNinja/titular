// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";

import {AgentToken} from "../src/launchpad/AgentToken.sol";
import {BondingCurve} from "../src/launchpad/BondingCurve.sol";
import {LaunchpadFactory} from "../src/launchpad/LaunchpadFactory.sol";
import {Graduator} from "../src/launchpad/Graduator.sol";
import {FeeRouter} from "../src/launchpad/FeeRouter.sol";

import {AntiSniperModule} from "../src/launchpad/modules/AntiSniperModule.sol";
import {LaunchRadarModule} from "../src/launchpad/modules/LaunchRadarModule.sol";
import {SixtyDaysModule} from "../src/launchpad/modules/SixtyDaysModule.sol";
import {AirdropModule} from "../src/launchpad/modules/AirdropModule.sol";
import {CapitalFormationModule} from "../src/launchpad/modules/CapitalFormationModule.sol";
import {PreBuyModule} from "../src/launchpad/modules/PreBuyModule.sol";
import {ExistingTokenModule} from "../src/launchpad/modules/ExistingTokenModule.sol";

import {IUniswapV2Router02} from "../src/launchpad/interfaces/IUniswapV2Router02.sol";

/// @notice Phase-2 deployment script.
/// @dev    Deployment topology:
///           1. AgentToken implementation (EIP-1167 master copy).
///           2. BondingCurve sentinel implementation (constructor-bound; the
///              factory's `bondingCurveImpl` slot is reserved for a future
///              cloneable refactor and only requires a non-zero address today).
///           3. FeeRouter (treasury-funded creator/protocol split).
///           4. Graduator (Uniswap V2 graduation sink, owned by deployer
///              temporarily; ownership is transferred to the factory at the
///              end of the script so {LaunchpadFactory.launchAgent} can call
///              {Graduator.registerLPLock}).
///           5. Seven launchpad modules (AntiSniper, LaunchRadar, SixtyDays,
///              Airdrop, CapitalFormation, PreBuy, ExistingToken).
///           6. LaunchpadFactory wired to TITU + Treasury (Phase-1), the
///              AgentToken impl, the sentinel curve impl, FeeRouter,
///              Graduator, and the externally-supplied Uniswap V2
///              factory + router.
///           7. Wire-up: factory.setModule(...) for each module, then
///              graduator.transferOwnership(factory).
///
///         Idempotency:
///           - Re-runs read the existing `deployments/<chain>.json`. Any
///             Phase-2 contract whose address is already populated and has
///             code at that address is reused; only missing entries get
///             redeployed. The merged JSON written at the end always carries
///             both the preserved Phase-1 keys and the (possibly partially
///             reused) Phase-2 keys.
///
///         Preconditions (checked in order):
///           1. `deployments/<chain>.json` must exist and contain a
///              populated Phase-1 contract set (`TITU` + `Treasury`
///              non-zero). Reverts with {Phase1DeploymentMissing} otherwise.
///           2. `UNISWAP_V2_FACTORY` and `UNISWAP_V2_ROUTER` env vars must be
///              set. Reverts with {UniswapV2NotConfigured} otherwise.
///
///         Output JSON shape (merged):
///           ```
///           {
///             "chainId": <id>,
///             "phase1": { "deployedAt": <ts>, "contracts": { ... } },
///             "phase2": { "deployedAt": <ts>, "contracts": { ... } }
///           }
///           ```
///         Legacy Phase-1-only files (top-level `contracts`) are auto-migrated
///         into the `phase1` key on first Phase-2 run.
contract DeployPhase2 is Script {
    /// @notice Phase-2 deployed contract set. Mirrors the JSON keys.
    struct Phase2Deployment {
        address agentTokenImpl;
        address bondingCurveImpl;
        address feeRouter;
        address graduator;
        address launchpadFactory;
        address antiSniper;
        address launchRadar;
        address sixtyDays;
        address airdrop;
        address capitalFormation;
        address preBuy;
        address existingToken;
    }

    /// @notice Phase-1 contract set, surfaced for the JSON merge.
    struct Phase1Snapshot {
        address titu;
        address treasury;
        address treasuryImpl;
        address vestingVault;
        address veTitu;
        address feeDistributor;
        uint256 deployedAt;
    }

    /// @notice Spec-pinned graduation reserves; mirrored from the unit-test
    ///         helpers so the sentinel BondingCurve constructor passes.
    uint256 internal constant VIRTUAL_QUOTE_RESERVE = 30_000e18;
    uint256 internal constant VIRTUAL_AGENT_RESERVE = 1_073_000_000e18;
    uint256 internal constant GRADUATION_THRESHOLD = 42_000e18;
    uint256 internal constant INITIAL_AGENT_RESERVE = 1_000_000_000e18;

    /// @dev Thrown when the prior Phase-1 deployment file does not exist or
    ///      has not been populated with TITU + Treasury addresses.
    error Phase1DeploymentMissing(string path);
    /// @dev Thrown when `UNISWAP_V2_FACTORY` or `UNISWAP_V2_ROUTER` env vars
    ///      are missing. Phase 2 cannot be wired without both.
    error UniswapV2NotConfigured();

    /// @notice Standard chain-id -> filename mapping. Falls back to
    ///         `<chainId>.json` for chains we have not pinned a name for.
    function _deploymentPath(uint256 chainId) internal pure returns (string memory) {
        if (chainId == 84_532) return "deployments/base-sepolia.json";
        if (chainId == 8453) return "deployments/base.json";
        if (chainId == 31_337) return "deployments/local.json";
        return string.concat("deployments/", vm.toString(chainId), ".json");
    }

    /// @notice Entry point. Reads the canonical deployment path for the
    ///         active chain and pulls the Uniswap V2 addresses from
    ///         `UNISWAP_V2_FACTORY` / `UNISWAP_V2_ROUTER` env vars.
    function run() external returns (Phase2Deployment memory d) {
        (address uniFactory, address uniRouter) = _readUniswapEnv();
        return runWith(_deploymentPath(block.chainid), uniFactory, uniRouter);
    }

    /// @notice Path-explicit entry point that still pulls Uniswap addresses
    ///         from env. Useful for forked dry-runs against a non-default
    ///         deployments file.
    function runWithPath(string memory path) public returns (Phase2Deployment memory d) {
        (address uniFactory, address uniRouter) = _readUniswapEnv();
        return runWith(path, uniFactory, uniRouter);
    }

    /// @notice Fully-explicit entry point. Lets the test harness inject the
    ///         Uniswap addresses directly so we never touch the process env
    ///         (which races across parallel-run test cases).
    /// @dev    `uniFactory_` and `uniRouter_` of `address(0)` triggers
    ///         {UniswapV2NotConfigured} — same surface as the env path.
    /// @param  path        Deployment JSON to read + write.
    /// @param  uniFactory_ Uniswap V2 factory; non-zero.
    /// @param  uniRouter_  Uniswap V2 router; non-zero.
    function runWith(string memory path, address uniFactory_, address uniRouter_)
        public
        returns (Phase2Deployment memory d)
    {
        // 1. Validate Phase-1 prerequisite BEFORE Uniswap. A missing
        //    Phase-1 file is a misordered ceremony, not a misconfiguration —
        //    surface the clearer error first.
        (string memory existingJson, Phase1Snapshot memory phase1) = _requirePhase1(path);

        // 2. Validate Uniswap V2 wiring. Required for both the factory and
        //    the graduator; checked together so a partial config is not
        //    silently accepted.
        if (uniFactory_ == address(0) || uniRouter_ == address(0)) revert UniswapV2NotConfigured();
        address uniFactory = uniFactory_;
        address uniRouter = uniRouter_;

        uint256 pk = vm.envOr("DEPLOYER_PRIVATE_KEY", uint256(0));
        // When running with --unlocked / no private key, use msg.sender (the
        // default sender). address(this) is banned in script contracts by
        // foundry >=1.5.1 because script contracts are ephemeral.
        address deployer = pk == 0 ? msg.sender : vm.addr(pk);

        // 3. Recover any Phase-2 entries that already exist on the JSON so a
        //    re-run does not redeploy what is already live.
        Phase2Deployment memory existing = _readExistingPhase2(existingJson);

        if (pk != 0) vm.startBroadcast(pk);

        // --- Logic contracts -------------------------------------------------
        d.agentTokenImpl = _hasCode(existing.agentTokenImpl) ? existing.agentTokenImpl : address(new AgentToken());

        // Sentinel curve impl: constructor-bound in the merged BondingCurve.
        // We pass canonical reserves and a sentinel agent/quote so the
        // constructor passes. The slot is reserved on the factory and never
        // invoked on the launch path; this lets us pin the impl ahead of the
        // future cloneable refactor without leaving a zero address.
        if (_hasCode(existing.bondingCurveImpl)) {
            d.bondingCurveImpl = existing.bondingCurveImpl;
        } else {
            d.bondingCurveImpl = address(
                new BondingCurve(
                    deployer,
                    address(uint160(uint256(keccak256("phase2.curve.sentinel.agent")))),
                    phase1.titu,
                    address(uint160(uint256(keccak256("phase2.curve.sentinel.feeRouter")))),
                    VIRTUAL_QUOTE_RESERVE,
                    VIRTUAL_AGENT_RESERVE,
                    GRADUATION_THRESHOLD,
                    INITIAL_AGENT_RESERVE
                )
            );
        }

        // FeeRouter pulls real ERC-20 routes; treasury is the Phase-1 address.
        d.feeRouter =
            _hasCode(existing.feeRouter) ? existing.feeRouter : address(new FeeRouter(phase1.treasury, deployer));

        // Graduator constructed with the deployer as owner. We hand ownership
        // to the factory after the factory deploys so the factory can call
        // {Graduator.registerLPLock} on every launch.
        d.graduator = _hasCode(existing.graduator)
            ? existing.graduator
            : address(new Graduator(IUniswapV2Router02(uniRouter), deployer));

        // --- Modules ---------------------------------------------------------
        d.antiSniper = _hasCode(existing.antiSniper) ? existing.antiSniper : address(new AntiSniperModule(deployer));
        d.launchRadar = _hasCode(existing.launchRadar) ? existing.launchRadar : address(new LaunchRadarModule(deployer));
        d.sixtyDays = _hasCode(existing.sixtyDays) ? existing.sixtyDays : address(new SixtyDaysModule(deployer));
        d.airdrop = _hasCode(existing.airdrop) ? existing.airdrop : address(new AirdropModule(deployer));
        d.capitalFormation = _hasCode(existing.capitalFormation)
            ? existing.capitalFormation
            : address(new CapitalFormationModule(deployer));
        d.preBuy = _hasCode(existing.preBuy) ? existing.preBuy : address(new PreBuyModule(deployer));
        d.existingToken =
            _hasCode(existing.existingToken) ? existing.existingToken : address(new ExistingTokenModule(deployer));

        // --- Factory ---------------------------------------------------------
        if (_hasCode(existing.launchpadFactory)) {
            d.launchpadFactory = existing.launchpadFactory;
        } else {
            d.launchpadFactory = address(
                new LaunchpadFactory(
                    phase1.titu,
                    phase1.treasury,
                    d.agentTokenImpl,
                    d.bondingCurveImpl,
                    d.feeRouter,
                    d.graduator,
                    uniFactory,
                    uniRouter,
                    deployer
                )
            );
        }

        // --- Wiring ----------------------------------------------------------
        // Module registrations are idempotent — `setModule` overwrites with
        // the same address as a no-op (just re-emits `ModuleSet`).
        LaunchpadFactory factory = LaunchpadFactory(d.launchpadFactory);
        if (factory.owner() == deployer) {
            _wireModules(factory, d);
        }

        // Hand graduator ownership to the factory so launch-time
        // `registerLPLock` calls succeed. Skip if already transferred.
        Graduator graduator = Graduator(d.graduator);
        if (graduator.owner() == deployer) {
            graduator.transferOwnership(d.launchpadFactory);
        }

        if (pk != 0) vm.stopBroadcast();

        // --- Persist ---------------------------------------------------------
        _writeMergedDeployment(path, phase1, d);
        _logDeployment(d);
    }

    // -----------------------------------------------------------------------
    // Internal — preflight
    // -----------------------------------------------------------------------

    /// @dev Read the existing deployment file and assert Phase-1 contracts are
    ///      populated. Returns the raw JSON (for downstream merge) and a
    ///      typed Phase-1 snapshot. Reverts with {Phase1DeploymentMissing}
    ///      if the file is missing/empty or TITU/Treasury are unset.
    function _requirePhase1(string memory path) internal view returns (string memory json, Phase1Snapshot memory snap) {
        if (!vm.exists(path)) revert Phase1DeploymentMissing(path);
        json = vm.readFile(path);
        if (bytes(json).length == 0) revert Phase1DeploymentMissing(path);

        snap.titu = _readPhase1Address(json, "TITU");
        if (snap.titu == address(0)) revert Phase1DeploymentMissing(path);

        snap.treasury = _readPhase1Address(json, "Treasury");
        if (snap.treasury == address(0)) revert Phase1DeploymentMissing(path);

        // Optional — null on legacy fixtures, populated on real Phase-1 runs.
        snap.treasuryImpl = _readPhase1Address(json, "TreasuryImpl");
        snap.vestingVault = _readPhase1Address(json, "VestingVault");
        snap.veTitu = _readPhase1Address(json, "VeTITU");
        snap.feeDistributor = _readPhase1Address(json, "FeeDistributor");
        snap.deployedAt = _readPhase1DeployedAt(json);
    }

    /// @dev Pull the Uniswap V2 endpoints from process env. Returns
    ///      `(0, 0)` when either is unset or unparseable; the caller decides
    ///      whether the missing config is fatal (it is on the production
    ///      path; tests use {runWith} and skip env entirely).
    function _readUniswapEnv() internal view returns (address factory_, address router_) {
        factory_ = _envAddress("UNISWAP_V2_FACTORY");
        router_ = _envAddress("UNISWAP_V2_ROUTER");
    }

    /// @dev Best-effort env address read. Returns `address(0)` when the var
    ///      is unset or unparseable. {vm.envOr} would normally suffice, but
    ///      on recent Foundry nightlies it has been observed to read stale
    ///      values across tests in the same suite — going through
    ///      `envExists` + `envString` keeps the lookup fully scoped to the
    ///      current call.
    function _envAddress(string memory name) internal view returns (address) {
        if (!vm.envExists(name)) return address(0);
        try this.envAddressExternal(name) returns (address a) {
            return a;
        } catch {
            return address(0);
        }
    }

    /// @dev External wrapper used by {_envAddress}'s `try/catch`. Foundry's
    ///      `envAddress` reverts when the value cannot be parsed; we want a
    ///      typed `address(0)` instead.
    function envAddressExternal(string memory name) external view returns (address) {
        return vm.envAddress(name);
    }

    /// @dev Try the merged-shape key (`phase1.contracts.<name>`) first, then
    ///      the legacy shape (`contracts.<name>`). Either yields a populated
    ///      address. Returns `address(0)` when the key is absent OR present
    ///      with a JSON `null` value (skeleton fixtures from a Phase-1 dry
    ///      run land in this state).
    function _readPhase1Address(string memory json, string memory name) internal view returns (address) {
        string memory mergedKey = string.concat(".phase1.contracts.", name);
        address a = _safeParseAddress(json, mergedKey);
        if (a != address(0)) return a;

        string memory legacyKey = string.concat(".contracts.", name);
        return _safeParseAddress(json, legacyKey);
    }

    /// @dev Pull `deployedAt` from either shape. Returns 0 (rendered as JSON
    ///      `null` on output) when missing or null in the source JSON.
    function _readPhase1DeployedAt(string memory json) internal view returns (uint256) {
        uint256 v = _safeParseUint(json, ".phase1.deployedAt");
        if (v != 0) return v;
        return _safeParseUint(json, ".deployedAt");
    }

    /// @dev Best-effort address read. `keyExistsJson` returns true for keys
    ///      bound to JSON `null`, which then crashes `parseJsonAddress`. We
    ///      route through an external try/catch so the failure surfaces as
    ///      `address(0)` and the caller can decide whether the field is
    ///      required.
    function _safeParseAddress(string memory json, string memory key) internal view returns (address) {
        if (!vm.keyExistsJson(json, key)) return address(0);
        try this.parseAddressExternal(json, key) returns (address a) {
            return a;
        } catch {
            return address(0);
        }
    }

    function _safeParseUint(string memory json, string memory key) internal view returns (uint256) {
        if (!vm.keyExistsJson(json, key)) return 0;
        try this.parseUintExternal(json, key) returns (uint256 v) {
            return v;
        } catch {
            return 0;
        }
    }

    /// @dev External wrappers used by {_safeParseAddress}/{_safeParseUint}'s
    ///      `try` blocks. Marked external so the EVM dispatches them as
    ///      separate frames and a parse-time revert can be caught.
    function parseAddressExternal(string memory json, string memory key) external pure returns (address) {
        return vm.parseJsonAddress(json, key);
    }

    function parseUintExternal(string memory json, string memory key) external pure returns (uint256) {
        return vm.parseJsonUint(json, key);
    }

    /// @dev Codehash != 0 for any deployed contract; address(0) hashes to 0
    ///      here, so `_hasCode` cleanly says "is this a deployed contract
    ///      right now" without false positives.
    function _hasCode(address a) internal view returns (bool) {
        if (a == address(0)) return false;
        return a.code.length > 0;
    }

    // -----------------------------------------------------------------------
    // Internal — read existing Phase-2
    // -----------------------------------------------------------------------

    /// @dev Pull every Phase-2 address out of the existing JSON; missing keys
    ///      surface as `address(0)`. Caller compares against on-chain code to
    ///      decide reuse vs redeploy.
    function _readExistingPhase2(string memory json) internal view returns (Phase2Deployment memory d) {
        d.agentTokenImpl = _readPhase2Address(json, "AgentTokenImpl");
        d.bondingCurveImpl = _readPhase2Address(json, "BondingCurveImpl");
        d.feeRouter = _readPhase2Address(json, "FeeRouter");
        d.graduator = _readPhase2Address(json, "Graduator");
        d.launchpadFactory = _readPhase2Address(json, "LaunchpadFactory");
        d.antiSniper = _readPhase2Address(json, "AntiSniperModule");
        d.launchRadar = _readPhase2Address(json, "LaunchRadarModule");
        d.sixtyDays = _readPhase2Address(json, "SixtyDaysModule");
        d.airdrop = _readPhase2Address(json, "AirdropModule");
        d.capitalFormation = _readPhase2Address(json, "CapitalFormationModule");
        d.preBuy = _readPhase2Address(json, "PreBuyModule");
        d.existingToken = _readPhase2Address(json, "ExistingTokenModule");
    }

    function _readPhase2Address(string memory json, string memory name) internal view returns (address) {
        string memory key = string.concat(".phase2.contracts.", name);
        return _safeParseAddress(json, key);
    }

    // -----------------------------------------------------------------------
    // Internal — wiring
    // -----------------------------------------------------------------------

    /// @dev Register every module on the factory. Idempotent: re-binding the
    ///      same address is a no-op semantically and just emits a fresh
    ///      `ModuleSet` event.
    function _wireModules(LaunchpadFactory factory, Phase2Deployment memory d) internal {
        factory.setModule(factory.MODULE_ID_ANTI_SNIPER(), d.antiSniper);
        factory.setModule(factory.MODULE_ID_LAUNCH_RADAR(), d.launchRadar);
        factory.setModule(factory.MODULE_ID_SIXTY_DAYS(), d.sixtyDays);
        factory.setModule(factory.MODULE_ID_AIRDROP(), d.airdrop);
        factory.setModule(factory.MODULE_ID_CAPITAL_FORMATION(), d.capitalFormation);
        factory.setModule(factory.MODULE_ID_PRE_BUY(), d.preBuy);
        factory.setModule(factory.MODULE_ID_EXISTING_TOKEN(), d.existingToken);
    }

    // -----------------------------------------------------------------------
    // Internal — JSON merge / write
    // -----------------------------------------------------------------------

    /// @dev Build the merged JSON manually. Using `stdJson.serialize*` would
    ///      bleed state across runs (the cheatcode caches keys per object id),
    ///      which makes ordering and idempotency tricky to reason about. The
    ///      string concat is verbose but auditable.
    function _writeMergedDeployment(string memory path, Phase1Snapshot memory phase1, Phase2Deployment memory d)
        internal
    {
        string memory phase1Block = _phase1Block(phase1);
        string memory phase2Block = _phase2Block(d);

        string memory merged = string.concat(
            "{\n",
            '  "chainId": ',
            vm.toString(block.chainid),
            ",\n",
            '  "phase1": ',
            phase1Block,
            ",\n",
            '  "phase2": ',
            phase2Block,
            "\n}\n"
        );

        vm.writeFile(path, merged);
    }

    /// @dev Render the Phase-1 block from a typed snapshot. Optional Phase-1
    ///      fields render as JSON `null` when not present in the source file.
    function _phase1Block(Phase1Snapshot memory s) internal pure returns (string memory) {
        return string.concat(
            "{\n",
            '    "deployedAt": ',
            _uintLit(s.deployedAt),
            ",\n",
            '    "contracts": {\n',
            '      "TITU": ',
            _addrLit(s.titu),
            ",\n",
            '      "Treasury": ',
            _addrLit(s.treasury),
            ",\n",
            '      "TreasuryImpl": ',
            _addrLit(s.treasuryImpl),
            ",\n",
            '      "VestingVault": ',
            _addrLit(s.vestingVault),
            ",\n",
            '      "VeTITU": ',
            _addrLit(s.veTitu),
            ",\n",
            '      "FeeDistributor": ',
            _addrLit(s.feeDistributor),
            "\n    }\n  }"
        );
    }

    /// @dev Build the Phase-2 block. Every key is always present — Phase-2
    ///      is the source of truth for itself and only this script writes it.
    ///      Split into halves to keep the EVM stack under 16 slots without
    ///      `via_ir`.
    function _phase2Block(Phase2Deployment memory d) internal view returns (string memory) {
        return string.concat(
            "{\n",
            '    "deployedAt": ',
            vm.toString(block.timestamp),
            ",\n",
            '    "contracts": {\n',
            _phase2ContractsHalfA(d),
            _phase2ContractsHalfB(d),
            "\n    }\n  }"
        );
    }

    function _phase2ContractsHalfA(Phase2Deployment memory d) internal pure returns (string memory) {
        return string.concat(
            '      "AgentTokenImpl": ',
            _addrLitRequired(d.agentTokenImpl),
            ",\n",
            '      "BondingCurveImpl": ',
            _addrLitRequired(d.bondingCurveImpl),
            ",\n",
            '      "FeeRouter": ',
            _addrLitRequired(d.feeRouter),
            ",\n",
            '      "Graduator": ',
            _addrLitRequired(d.graduator),
            ",\n",
            '      "LaunchpadFactory": ',
            _addrLitRequired(d.launchpadFactory),
            ",\n",
            '      "AntiSniperModule": ',
            _addrLitRequired(d.antiSniper)
        );
    }

    function _phase2ContractsHalfB(Phase2Deployment memory d) internal pure returns (string memory) {
        return string.concat(
            ",\n",
            '      "LaunchRadarModule": ',
            _addrLitRequired(d.launchRadar),
            ",\n",
            '      "SixtyDaysModule": ',
            _addrLitRequired(d.sixtyDays),
            ",\n",
            '      "AirdropModule": ',
            _addrLitRequired(d.airdrop),
            ",\n",
            '      "CapitalFormationModule": ',
            _addrLitRequired(d.capitalFormation),
            ",\n",
            '      "PreBuyModule": ',
            _addrLitRequired(d.preBuy),
            ",\n",
            '      "ExistingTokenModule": ',
            _addrLitRequired(d.existingToken)
        );
    }

    /// @dev Quoted hex literal. Phase-2 contracts are always set by the time
    ///      we render, so a zero address indicates a logic bug; render a
    ///      visible literal instead of crashing the JSON.
    function _addrLitRequired(address a) internal pure returns (string memory) {
        return string.concat('"', vm.toString(a), '"');
    }

    function _addrLit(address a) internal pure returns (string memory) {
        if (a == address(0)) return "null";
        return string.concat('"', vm.toString(a), '"');
    }

    function _uintLit(uint256 v) internal pure returns (string memory) {
        if (v == 0) return "null";
        return vm.toString(v);
    }

    function _logDeployment(Phase2Deployment memory d) internal pure {
        console2.log("AgentTokenImpl    ", d.agentTokenImpl);
        console2.log("BondingCurveImpl  ", d.bondingCurveImpl);
        console2.log("FeeRouter         ", d.feeRouter);
        console2.log("Graduator         ", d.graduator);
        console2.log("LaunchpadFactory  ", d.launchpadFactory);
        console2.log("AntiSniper        ", d.antiSniper);
        console2.log("LaunchRadar       ", d.launchRadar);
        console2.log("SixtyDays         ", d.sixtyDays);
        console2.log("Airdrop           ", d.airdrop);
        console2.log("CapitalFormation  ", d.capitalFormation);
        console2.log("PreBuy            ", d.preBuy);
        console2.log("ExistingToken     ", d.existingToken);
    }
}
