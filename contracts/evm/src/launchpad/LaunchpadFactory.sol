// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {Clones} from "@openzeppelin/contracts/proxy/Clones.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {AgentToken} from "./AgentToken.sol";
import {BondingCurve} from "./BondingCurve.sol";
import {LPLock} from "./LPLock.sol";
import {Graduator} from "./Graduator.sol";

import {IBondingCurve} from "./interfaces/IBondingCurve.sol";
import {IUniswapV2Factory} from "./interfaces/IUniswapV2Factory.sol";
import {IUniswapV2Router02} from "./interfaces/IUniswapV2Router02.sol";

import {IAntiSniperModule} from "./interfaces/IAntiSniperModule.sol";
import {ILaunchRadarModule} from "./interfaces/ILaunchRadarModule.sol";
import {ISixtyDaysModule} from "./interfaces/ISixtyDaysModule.sol";
import {IAirdropModule} from "./interfaces/IAirdropModule.sol";
import {ICapitalFormationModule} from "./interfaces/ICapitalFormationModule.sol";
import {IPreBuyModule} from "./interfaces/IPreBuyModule.sol";
import {IExistingTokenModule} from "./interfaces/IExistingTokenModule.sol";

/// @title LaunchpadFactory
/// @notice Single audited launch entry point for the Titular launchpad fleet.
///         A `launchAgent` call atomically:
///           1. Clones an `AgentToken` via EIP-1167 minimal proxy (45-byte
///              runtime), so every per-agent ERC-20 carries identical bytecode
///              and the deploy gas cost stays bounded.
///           2. Deploys a fresh `BondingCurve` via `new BondingCurve(...)`.
///              The curve is constructor-bound (immutables for token sides,
///              reserves, threshold) and therefore not eligible for cloning
///              today; this is documented as a deferred refactor and is the
///              only `new` site on the launch path. The `bondingCurveImpl`
///              constructor argument is reserved for the future cloneable
///              variant so deploy scripts can pin the impl in advance.
///           3. Pre-creates the Uniswap V2 pair so the per-agent {LPLock} can
///              bind the LP ERC-20 immutably at construction.
///           4. Deploys a per-agent {LPLock} whose beneficiary is the creator.
///           5. Wires the curve to the shared {Graduator} via
///              {BondingCurve.setGraduator} and registers the lock via
///              {Graduator.registerLPLock} so the graduator knows where to
///              mint LP at threshold.
///           6. Dispatches `configure(...)` into every enabled module in the
///              registry, using a typed interface per known module id and an
///              explicit `if/else` block — no generic decoder.
///         Free to launch (gas-only). Per-module fees, if any, are collected
///         inside the modules themselves on their own configure paths.
/// @dev    Owner is expected to be a Safe multisig that ALSO holds Ownable
///         rights on every module the factory dispatches into AND on the
///         shared {Graduator}. Wiring those ownerships is a deploy-script
///         responsibility and is deliberately out of scope here so the
///         factory's surface stays minimal. If the factory is not the
///         Graduator's owner at {launchAgent} time, the inner
///         {Graduator.registerLPLock} reverts with `OwnableUnauthorizedAccount`
///         — fail-loud by design.
///
///         Modules are bound to ids by {setModule} so new modules can be
///         added (or hot-swapped, after audit) without redeploying. Module
///         dispatch is implemented with explicit `if/else` branches on the
///         canonical id rather than a generic `bytes` decoder, so each path
///         is statically auditable. Adding a new module means adding a new
///         `MODULE_ID_*` constant and a new branch.
///
///         CEI is preserved on {launchAgent}: validate -> deploy -> wire ->
///         persist record + emit -> dispatch into modules. The agent record
///         is stored before any module call, so an indexer following the
///         {AgentLaunched} event can rely on {getAgent} immediately. A
///         per-module revert bubbles up and rolls the whole launch back —
///         atomicity is preserved. {ReentrancyGuard} covers the pathological
///         case of a malicious module re-entering {launchAgent} from inside
///         its own `configure`.
contract LaunchpadFactory is Ownable, ReentrancyGuard {
    // ---------------------------------------------------------------------
    // Types
    // ---------------------------------------------------------------------

    /// @notice On-chain record of every agent the factory has launched.
    /// @param token      ERC-20 clone for this agent.
    /// @param curve      Bonding curve paired with `token`.
    /// @param lpLock     Per-agent LP lock holding the post-graduation LP.
    /// @param pair       Uniswap V2 pair pre-created at launch (empty until
    ///                   graduation seeds it).
    /// @param creator    Account that called {launchAgent}; recipient of the
    ///                   locked LP at unlock.
    /// @param createdAt  Unix timestamp at which the agent was launched.
    /// @param modules    Module ids the factory dispatched at launch.
    struct Agent {
        address token;
        address curve;
        address lpLock;
        address pair;
        address creator;
        uint64 createdAt;
        bytes32[] modules;
    }

    /// @notice Inbound launch parameters. `enabledModules` and `moduleData`
    ///         are parallel arrays of equal length; `moduleData[i]` is the
    ///         abi-encoded configure payload for `enabledModules[i]`.
    /// @param name            ERC-20 name; non-empty, <= {MAX_NAME_LEN}.
    /// @param symbol          ERC-20 symbol; non-empty, <= {MAX_SYMBOL_LEN}.
    /// @param imageURI        Off-chain image pointer (may be empty).
    /// @param soulURI         Off-chain persona pointer (may be empty).
    /// @param enabledModules  Canonical module ids to dispatch at launch.
    /// @param moduleData      Per-module abi-encoded configure payloads.
    struct LaunchParams {
        string name;
        string symbol;
        string imageURI;
        string soulURI;
        bytes32[] enabledModules;
        bytes[] moduleData;
    }

    // ---------------------------------------------------------------------
    // Module ids
    // ---------------------------------------------------------------------

    /// @notice Canonical id for the AntiSniper module (`keccak256("AntiSniper")`).
    bytes32 public constant MODULE_ID_ANTI_SNIPER = keccak256("AntiSniper");

    /// @notice Canonical id for the LaunchRadar module (`keccak256("LaunchRadar")`).
    bytes32 public constant MODULE_ID_LAUNCH_RADAR = keccak256("LaunchRadar");

    /// @notice Canonical id for the SixtyDays module (`keccak256("SixtyDays")`).
    bytes32 public constant MODULE_ID_SIXTY_DAYS = keccak256("SixtyDays");

    /// @notice Canonical id for the Airdrop module (`keccak256("Airdrop")`).
    bytes32 public constant MODULE_ID_AIRDROP = keccak256("Airdrop");

    /// @notice Canonical id for the CapitalFormation module
    ///         (`keccak256("CapitalFormation")`).
    bytes32 public constant MODULE_ID_CAPITAL_FORMATION = keccak256("CapitalFormation");

    /// @notice Canonical id for the PreBuy module (`keccak256("PreBuy")`).
    bytes32 public constant MODULE_ID_PRE_BUY = keccak256("PreBuy");

    /// @notice Canonical id for the ExistingToken module
    ///         (`keccak256("ExistingToken")`).
    bytes32 public constant MODULE_ID_EXISTING_TOKEN = keccak256("ExistingToken");

    // ---------------------------------------------------------------------
    // Bounds
    // ---------------------------------------------------------------------

    /// @notice Hard upper bound on the ERC-20 name length, in bytes. The OZ
    ///         ERC-20 implementation does not bound `name` itself; we cap
    ///         here so an attacker cannot pad the name with megabytes of
    ///         calldata to grief the launch path's gas cost.
    uint256 public constant MAX_NAME_LEN = 64;

    /// @notice Hard upper bound on the ERC-20 symbol length, in bytes.
    uint256 public constant MAX_SYMBOL_LEN = 16;

    /// @notice Spec-pinned graduation threshold: real TITU at which trading
    ///         freezes and graduation may fire. 42_000 TITU at 18 decimals.
    uint256 public constant GRADUATION_THRESHOLD = 42_000e18;

    /// @notice Virtual quote reserve seeded into every curve. Bounds the
    ///         initial slope so the first buy does not deliver an absurd
    ///         agent allocation.
    uint256 public constant VIRTUAL_QUOTE_RESERVE = 30_000e18;

    /// @notice Virtual agent reserve seeded into every curve. Strictly
    ///         greater than {INITIAL_AGENT_RESERVE} so the curve always has
    ///         LP seed left at graduation.
    uint256 public constant VIRTUAL_AGENT_RESERVE = 1_073_000_000e18;

    /// @notice Real agent inventory the curve is loaded with at construction.
    ///         Equals the AgentToken total supply (1B) since the entire
    ///         minted balance lands on the curve.
    uint256 public constant INITIAL_AGENT_RESERVE = 1_000_000_000e18;

    // ---------------------------------------------------------------------
    // Immutables
    // ---------------------------------------------------------------------

    /// @notice TITU quote token used by every bonding curve as the base asset.
    address public immutable tituToken;

    /// @notice Protocol treasury — surfaced for downstream consumers; not
    ///         used on the launch path itself (FeeRouter holds the route).
    address public immutable treasury;

    /// @notice EIP-1167 implementation that every per-agent ERC-20 clones from.
    address public immutable agentTokenImpl;

    /// @notice EIP-1167 implementation slot reserved for a future cloneable
    ///         BondingCurve. Stored even though we currently `new` the curve
    ///         so the constructor signature matches the brief's wiring and
    ///         the deploy script can pin a real impl when the refactor lands.
    address public immutable bondingCurveImpl;

    /// @notice Shared FeeRouter every agent token routes its 1% transfer
    ///         tax to. Bound on every curve and every clone at launch.
    address public immutable feeRouter;

    /// @notice Shared graduation engine. The factory must already be (or be
    ///         made) the Ownable owner of this contract before any
    ///         {launchAgent} call lands; otherwise the inner
    ///         {Graduator.registerLPLock} reverts.
    Graduator public immutable graduator;

    /// @notice Uniswap V2 factory used to pre-create the (agentToken, TITU)
    ///         pair so the LP token address is known at LPLock construction.
    IUniswapV2Factory public immutable uniV2Factory;

    /// @notice Uniswap V2 router; surfaced for downstream consumers / future
    ///         modules. The graduator already binds its own router immutably.
    IUniswapV2Router02 public immutable uniV2Router;

    // ---------------------------------------------------------------------
    // Storage
    // ---------------------------------------------------------------------

    /// @notice Module id -> module address. Set by owner via {setModule}.
    ///         Empty mapping on deploy: the owner registers modules as they
    ///         go live so the factory can ship before every module exists.
    mapping(bytes32 moduleId => address module) public modules;

    /// @notice Sequential agent identifier. Always > 0; 0 is reserved for
    ///         "unset" so `getAgent(0)` reverts.
    uint256 private _agentIdCounter;

    /// @notice agentId -> on-chain record.
    mapping(uint256 agentId => Agent) private _agents;

    // ---------------------------------------------------------------------
    // Events
    // ---------------------------------------------------------------------

    /// @dev Emitted on {setModule}. `prev` is `address(0)` on first bind.
    event ModuleSet(bytes32 indexed moduleId, address indexed prev, address indexed module);

    /// @dev Emitted exactly once per successful {launchAgent}. EVM caps at
    ///      three indexed args, so we index the three values external
    ///      indexers most commonly query by — `agentId`, `token`, `curve`
    ///      — and pass the rest as data. `creator` is left non-indexed
    ///      because creators are queried via the agent record, not by
    ///      iterating launch events.
    event AgentLaunched(
        uint256 indexed agentId,
        address indexed token,
        address indexed curve,
        address creator,
        address lpLock,
        address pair,
        bytes32[] modules
    );

    // ---------------------------------------------------------------------
    // Errors
    // ---------------------------------------------------------------------

    /// @dev Thrown when a constructor or setter receives `address(0)`.
    error ZeroAddress();
    /// @dev Thrown when {LaunchParams.name} is empty.
    error EmptyName();
    /// @dev Thrown when {LaunchParams.symbol} is empty.
    error EmptySymbol();
    /// @dev Thrown when {LaunchParams.name} exceeds {MAX_NAME_LEN}.
    error NameTooLong();
    /// @dev Thrown when {LaunchParams.symbol} exceeds {MAX_SYMBOL_LEN}.
    error SymbolTooLong();
    /// @dev Thrown when `enabledModules.length != moduleData.length`.
    error LengthMismatch();
    /// @dev Thrown when {launchAgent} encounters a module id with no
    ///      registered address. Carries the offending id for indexer triage.
    error UnknownModule(bytes32 moduleId);
    /// @dev Thrown when {launchAgent} encounters the same module id twice in
    ///      a single launch's `enabledModules` array. Carries the duplicate.
    error DuplicateModule(bytes32 moduleId);
    /// @dev Thrown when {setModule} is called with the zero id.
    error InvalidModuleId();
    /// @dev Thrown by {getAgent} for an id outside `[1, agentCount]`.
    error UnknownAgent(uint256 agentId);

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    /// @param tituToken_         TITU quote token. Non-zero.
    /// @param treasury_          Protocol treasury. Non-zero.
    /// @param agentTokenImpl_    EIP-1167 implementation for AgentToken.
    /// @param bondingCurveImpl_  Reserved impl slot for a future cloneable
    ///                           curve. Required non-zero so deploy scripts
    ///                           cannot silently leave it unset; not consumed
    ///                           by the launch path today.
    /// @param feeRouter_         Shared FeeRouter that receives every
    ///                           agent's 1% transfer tax. Non-zero.
    /// @param graduator_         Shared Graduator. The factory must be (or
    ///                           be made) the Ownable owner of this contract
    ///                           before {launchAgent} is called.
    /// @param uniV2Factory_      Uniswap V2 factory used to pre-create pairs.
    /// @param uniV2Router_       Uniswap V2 router; pinned for downstream use.
    /// @param owner_             Initial owner; gates {setModule}.
    constructor(
        address tituToken_,
        address treasury_,
        address agentTokenImpl_,
        address bondingCurveImpl_,
        address feeRouter_,
        address graduator_,
        address uniV2Factory_,
        address uniV2Router_,
        address owner_
    ) Ownable(_validatedOwner(owner_)) {
        if (
            tituToken_ == address(0) || treasury_ == address(0) || agentTokenImpl_ == address(0)
                || bondingCurveImpl_ == address(0) || feeRouter_ == address(0) || graduator_ == address(0)
                || uniV2Factory_ == address(0) || uniV2Router_ == address(0)
        ) revert ZeroAddress();

        tituToken = tituToken_;
        treasury = treasury_;
        agentTokenImpl = agentTokenImpl_;
        bondingCurveImpl = bondingCurveImpl_;
        feeRouter = feeRouter_;
        graduator = Graduator(graduator_);
        uniV2Factory = IUniswapV2Factory(uniV2Factory_);
        uniV2Router = IUniswapV2Router02(uniV2Router_);
    }

    /// @dev Hoist the zero-address check above the OZ Ownable constructor so
    ///      callers see our canonical {ZeroAddress} error instead of OZ's
    ///      `OwnableInvalidOwner`.
    function _validatedOwner(address owner_) private pure returns (address) {
        if (owner_ == address(0)) revert ZeroAddress();
        return owner_;
    }

    // ---------------------------------------------------------------------
    // Admin
    // ---------------------------------------------------------------------

    /// @notice Bind a module address to a canonical id. Owner-gated.
    /// @dev    Setting `module = address(0)` is rejected: clearing a binding
    ///         once a module has gone live is an audit-relevant operation
    ///         and should go through a deliberate redeploy + governance
    ///         action rather than a single-tx flip. Reverts on the empty id
    ///         so `bytes32(0)` cannot be smuggled into a launch as a
    ///         sentinel.
    /// @param  moduleId Canonical id (e.g. {MODULE_ID_ANTI_SNIPER}).
    /// @param  module   Module address to bind. Must be non-zero.
    function setModule(bytes32 moduleId, address module) external onlyOwner {
        if (moduleId == bytes32(0)) revert InvalidModuleId();
        if (module == address(0)) revert ZeroAddress();
        address prev = modules[moduleId];
        modules[moduleId] = module;
        emit ModuleSet(moduleId, prev, module);
    }

    // ---------------------------------------------------------------------
    // Launch
    // ---------------------------------------------------------------------

    /// @notice Atomically launch a new agent.
    /// @dev    Steps (CEI):
    ///           1. Validate {LaunchParams} and the enabled-module set.
    ///           2. Clone AgentToken; deploy BondingCurve, pair, LPLock.
    ///           3. Initialize the cloned AgentToken (mints supply to curve).
    ///           4. Wire graduator: register the lock + bind on the curve.
    ///           5. Persist the agent record and bump the id counter.
    ///           6. Emit `AgentLaunched`.
    ///           7. Dispatch into every enabled module via typed interface.
    ///         Module dispatch is the only block of "external interaction
    ///         that may revert per-module"; the agent record is already
    ///         persisted by then so an indexer following `AgentLaunched`
    ///         sees the agent regardless of which modules ran — and a
    ///         per-module revert bubbles up and rolls back the entire
    ///         launch (atomicity preserved). `nonReentrant` covers the
    ///         pathological case of a malicious module re-entering
    ///         {launchAgent} from inside its own `configure`.
    /// @param  params Launch parameters; see {LaunchParams}.
    /// @return agentToken    The per-agent ERC-20 clone.
    /// @return bondingCurve  The per-agent BondingCurve.
    /// @return agentId       Sequential id starting from 1.
    function launchAgent(LaunchParams calldata params)
        external
        nonReentrant
        returns (address agentToken, address bondingCurve, uint256 agentId)
    {
        // --- Checks ---
        _validateParams(params);

        address creator = msg.sender;
        address lpLock;
        address pair;

        // --- Effects (deploy + wire per-agent artifacts) ---
        // Pulled out of this function so the local-variable footprint of
        // {launchAgent} fits the EVM's 16-slot stack budget when via-ir
        // is disabled (the project default; see foundry.toml).
        (agentToken, bondingCurve, lpLock, pair) = _deployAgentArtifacts(params.name, params.symbol, creator);

        // Persist the record before any module dispatch so an indexer
        // reading {AgentLaunched} can trust {getAgent} immediately.
        unchecked {
            // Counter increments by 1 per launch; cannot overflow within
            // any realistic block budget (each launch is ~1M gas).
            agentId = ++_agentIdCounter;
        }
        _writeAgentRecord(agentId, agentToken, bondingCurve, lpLock, pair, creator, params.enabledModules);

        emit AgentLaunched(agentId, agentToken, bondingCurve, creator, lpLock, pair, params.enabledModules);

        // --- Interactions (module dispatch) ---
        // Routed through a private helper to keep this function's stack
        // under the EVM's 16-slot budget when via-ir is disabled.
        _dispatchAll(agentToken, bondingCurve, creator, params.enabledModules, params.moduleData);
    }

    /// @dev Walk the parallel `enabledModules`/`moduleData` arrays and dispatch
    ///      into each module's typed configure surface. Pulled out of
    ///      {launchAgent} so the parent stays under the EVM's stack limit.
    ///      A per-module revert bubbles up and rolls back the entire launch.
    function _dispatchAll(
        address agentToken,
        address bondingCurve,
        address creator,
        bytes32[] calldata enabledModules,
        bytes[] calldata moduleData
    ) private {
        uint256 n = enabledModules.length;
        for (uint256 i = 0; i < n;) {
            bytes32 id = enabledModules[i];
            // _validateParams already asserted every id is registered, so a
            // zero check here would be redundant. We keep the explicit
            // revert path on dispatch for the case where an unhandled id
            // sneaks through (defence-in-depth — see {_dispatchModule}).
            _dispatchModule(id, modules[id], agentToken, bondingCurve, creator, moduleData[i]);
            unchecked {
                ++i;
            }
        }
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    /// @notice Total number of agents the factory has launched.
    function agentCount() external view returns (uint256) {
        return _agentIdCounter;
    }

    /// @notice Return the on-chain record for `agentId`.
    /// @dev    Reverts on out-of-range ids (`0` or `> agentCount`) so callers
    ///         do not have to disambiguate the zero struct from a real record.
    function getAgent(uint256 agentId) external view returns (Agent memory) {
        if (agentId == 0 || agentId > _agentIdCounter) revert UnknownAgent(agentId);
        return _agents[agentId];
    }

    // ---------------------------------------------------------------------
    // Internal — deployment helpers
    // ---------------------------------------------------------------------

    /// @dev Deploy the four per-agent artifacts in dependency order:
    ///      AgentToken clone -> BondingCurve -> V2 pair -> LPLock, then wire
    ///      the curve's graduator and the graduator's LP-lock recipient.
    ///      Pulled out of {launchAgent} purely to keep the caller's stack
    ///      under the EVM's 16-slot limit; no policy lives here.
    function _deployAgentArtifacts(string calldata name_, string calldata symbol_, address creator)
        private
        returns (address agentToken, address curve, address lpLock, address pair)
    {
        // 1) Clone the AgentToken implementation. EIP-1167 minimal proxy:
        //    constant 45-byte runtime, deterministic deployment cost.
        agentToken = Clones.clone(agentTokenImpl);

        // 2) Deploy a fresh BondingCurve. NOTE: the merged BondingCurve uses
        //    constructor + immutables and is therefore not cloneable today.
        //    The factory's `bondingCurveImpl` slot is reserved for a future
        //    cloneable refactor; until then `new BondingCurve(...)` is the
        //    only `new` site on the launch path and is the documented
        //    deviation from the EIP-1167-everywhere policy.
        BondingCurve curveContract = new BondingCurve(
            address(this), // owner — factory will call setGraduator below
            agentToken,
            tituToken,
            feeRouter,
            VIRTUAL_QUOTE_RESERVE,
            VIRTUAL_AGENT_RESERVE,
            GRADUATION_THRESHOLD,
            INITIAL_AGENT_RESERVE
        );
        curve = address(curveContract);

        // 3) Initialize the AgentToken AFTER the curve exists so the freshly
        //    minted supply lands directly on the curve. AgentToken's
        //    `_disableInitializers` call in the impl constructor blocks the
        //    impl itself from being initialized; only clones see this entry.
        AgentToken(agentToken).initialize(name_, symbol_, creator, feeRouter, curve);

        // 4) Pre-create the V2 pair so the LP ERC-20 address is known at
        //    LPLock construction. Idempotent on subsequent runs because we
        //    gate on `getPair` first.
        pair = _ensurePair(agentToken);

        // 5) Per-agent LPLock binds the LP token immutably; the beneficiary
        //    is the creator (sole withdrawal recipient at unlock).
        lpLock = address(new LPLock(IERC20(pair), creator));

        // 6) Wire the graduator: bind the lock target so {Graduator.graduate}
        //    knows where to mint LP, then bind the graduator on the curve so
        //    it is authorized to drain the reserves at graduation.
        //    NOTE: both calls require the factory to be the respective
        //    contract's owner. {Graduator.registerLPLock} is `onlyOwner` and
        //    {BondingCurve.setGraduator} is `onlyOwner` (the curve's
        //    constructor wired the factory as owner above).
        graduator.registerLPLock(curve, lpLock);
        curveContract.setGraduator(address(graduator));
    }

    /// @dev Returns the V2 pair address for `(agentToken, tituToken)`,
    ///      creating it via the Uniswap V2 factory if missing.
    function _ensurePair(address agentToken_) private returns (address pair) {
        pair = uniV2Factory.getPair(agentToken_, tituToken);
        if (pair == address(0)) {
            pair = uniV2Factory.createPair(agentToken_, tituToken);
        }
    }

    /// @dev Persist the on-chain agent record. Pulled out of {launchAgent}
    ///      so the parent stays under the stack limit; copies the enabled
    ///      module ids one-by-one because Solidity does not let us assign
    ///      a calldata `bytes32[]` directly into a storage `bytes32[]`.
    function _writeAgentRecord(
        uint256 agentId,
        address agentToken,
        address curve,
        address lpLock,
        address pair,
        address creator,
        bytes32[] calldata enabledModules
    ) private {
        Agent storage rec = _agents[agentId];
        rec.token = agentToken;
        rec.curve = curve;
        rec.lpLock = lpLock;
        rec.pair = pair;
        rec.creator = creator;
        // `block.timestamp` truncated to uint64. Won't wrap before year ~2554.
        // slither-disable-next-line timestamp
        rec.createdAt = uint64(block.timestamp);
        uint256 n = enabledModules.length;
        for (uint256 i = 0; i < n;) {
            rec.modules.push(enabledModules[i]);
            unchecked {
                ++i;
            }
        }
    }

    // ---------------------------------------------------------------------
    // Internal — validation
    // ---------------------------------------------------------------------

    /// @dev Validate {LaunchParams} up-front. Splits into a private function
    ///      to keep the caller's stack small and group the revert paths.
    ///      O(n^2) duplicate scan is acceptable: launches are gas-bounded
    ///      and `n` is bounded by the registry size (<10 today, no
    ///      realistic growth path that would push it past a few dozen).
    function _validateParams(LaunchParams calldata params) private view {
        bytes calldata nameBytes = bytes(params.name);
        bytes calldata symbolBytes = bytes(params.symbol);
        if (nameBytes.length == 0) revert EmptyName();
        if (symbolBytes.length == 0) revert EmptySymbol();
        if (nameBytes.length > MAX_NAME_LEN) revert NameTooLong();
        if (symbolBytes.length > MAX_SYMBOL_LEN) revert SymbolTooLong();

        uint256 n = params.enabledModules.length;
        if (n != params.moduleData.length) revert LengthMismatch();
        for (uint256 i = 0; i < n;) {
            bytes32 id = params.enabledModules[i];
            if (modules[id] == address(0)) revert UnknownModule(id);
            for (uint256 j = i + 1; j < n;) {
                if (params.enabledModules[j] == id) revert DuplicateModule(id);
                unchecked {
                    ++j;
                }
            }
            unchecked {
                ++i;
            }
        }
    }

    // ---------------------------------------------------------------------
    // Internal — module dispatch
    // ---------------------------------------------------------------------

    /// @dev Decode `data` per `id` and forward to the module's typed
    ///      configure surface. Adding a new module means adding a new
    ///      branch here AND a `MODULE_ID_*` constant — by design, so that
    ///      every dispatch path is statically auditable.
    /// @param  id           Canonical module id.
    /// @param  module       Registered module address.
    /// @param  agentToken_  Agent address — first argument to every module's
    ///                      `configure` (except {ExistingTokenModule}).
    /// @param  curve        Bonding curve for this launch; passed to modules
    ///                      that need a per-agent curve binding.
    /// @param  creator      Agent creator; used as the agent admin for
    ///                      modules that delegate per-agent control.
    /// @param  data         abi-encoded module-specific payload.
    function _dispatchModule(
        bytes32 id,
        address module,
        address agentToken_,
        address curve,
        address creator,
        bytes calldata data
    ) private {
        if (id == MODULE_ID_ANTI_SNIPER) {
            (uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps) =
                abi.decode(data, (uint64, uint32, uint16, uint16));
            IAntiSniperModule(module).configure(agentToken_, startTime, duration, startBps, endBps);
        } else if (id == MODULE_ID_LAUNCH_RADAR) {
            // LaunchRadar takes a per-agent admin; the factory hands the
            // creator full mutation rights post-bootstrap so the protocol
            // owner does not hold per-agent keys.
            (uint64 startTime, bool whitelistOn) = abi.decode(data, (uint64, bool));
            ILaunchRadarModule(module).configure(agentToken_, startTime, whitelistOn, creator);
        } else if (id == MODULE_ID_SIXTY_DAYS) {
            // SixtyDays binds the curve as the only authorized accrual / contribution
            // source and the creator as the sole party authorized to commit/refund.
            (address escrowToken, uint64 startTime) = abi.decode(data, (address, uint64));
            ISixtyDaysModule(module).configure(agentToken_, escrowToken, curve, creator, startTime);
        } else if (id == MODULE_ID_AIRDROP) {
            // Airdrop pre-supply MUST already sit on the module before this dispatch
            // lands; the module's configure validates the balance covers `allocation`.
            (bytes32 root, uint64 snapshotBlock, uint128 allocation, uint64 deadline) =
                abi.decode(data, (bytes32, uint64, uint128, uint64));
            IAirdropModule(module).configure(agentToken_, root, snapshotBlock, allocation, deadline, creator);
        } else if (id == MODULE_ID_CAPITAL_FORMATION) {
            // CapitalFormation needs the V2 pair as its FDV oracle. The factory
            // does not know which payout token / threshold schedule the creator
            // wants, so the full payload is abi-decoded from `data`.
            (address payoutToken, address pair, uint256[] memory thresholds, uint256[] memory payouts) =
                abi.decode(data, (address, address, uint256[], uint256[]));
            ICapitalFormationModule(module)
                .configure(agentToken_, agentToken_, payoutToken, pair, creator, creator, thresholds, payouts);
        } else if (id == MODULE_ID_PRE_BUY) {
            // PreBuy requires `titanIn` of TITU to already sit on the module.
            // The returned vesting-clone address is intentionally discarded —
            // it is emitted by the module's own `Configured` event so an
            // indexer can stitch it back to this agent off-chain.
            (uint256 vestAmount, uint64 cliffSeconds, uint64 durationSeconds, uint256 titanIn) =
                abi.decode(data, (uint256, uint64, uint64, uint256));
            // slither-disable-next-line unused-return
            IPreBuyModule(module)
                .configure(
                    agentToken_, creator, vestAmount, cliffSeconds, durationSeconds, titanIn, IBondingCurve(curve)
                );
        } else if (id == MODULE_ID_EXISTING_TOKEN) {
            // ExistingToken's first argument is the external token, not the
            // agent clone. The supply must already be on the module.
            (address externalToken, address externalCurve, uint256 supply) =
                abi.decode(data, (address, address, uint256));
            IExistingTokenModule(module).configure(externalToken, externalCurve, supply, creator);
        } else {
            // Module is registered but no dispatch branch exists yet.
            // Treat as misconfiguration: revert rather than silently no-op.
            revert UnknownModule(id);
        }
    }
}
