// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test, Vm} from "forge-std/Test.sol";

import {LaunchpadFactory} from "../../src/launchpad/LaunchpadFactory.sol";
import {AgentToken} from "../../src/launchpad/AgentToken.sol";
import {BondingCurve} from "../../src/launchpad/BondingCurve.sol";
import {LPLock} from "../../src/launchpad/LPLock.sol";
import {Graduator} from "../../src/launchpad/Graduator.sol";
import {FeeRouter} from "../../src/launchpad/FeeRouter.sol";

import {AntiSniperModule} from "../../src/launchpad/modules/AntiSniperModule.sol";

import {IUniswapV2Router02} from "../../src/launchpad/interfaces/IUniswapV2Router02.sol";

import {MockUniswapV2Factory} from "../mocks/MockUniswapV2Factory.sol";
import {MockUniswapV2Router} from "../mocks/MockUniswapV2Router.sol";
import {MockMintableERC20} from "../mocks/MockERC20.sol";

/// @title Launch.t.sol
/// @notice End-to-end integration test for {LaunchpadFactory.launchAgent}.
///         Wires every shared contract (factory, graduator, fee router,
///         AntiSniper module) the way the production deploy script does,
///         then drives a single launch and asserts:
///           - the AgentToken is an EIP-1167 clone of the impl;
///           - the BondingCurve is a fresh deploy with the right immutables
///             and the full 1B agent supply on hand;
///           - the LPLock binds the V2 pair as its LP token and the creator
///             as beneficiary;
///           - the Graduator has the lock registered for the curve;
///           - the AntiSniperModule received its `configure` call;
///           - {AgentLaunched} fired with the expected indexed + data fields.
/// @dev    Uses the in-tree Mock V2 factory + router so no live RPC is
///         required. The mock router records add-liquidity calls; this file
///         does not graduate, only launches.
contract LaunchIntegrationTest is Test {
    // ---------------------------------------------------------------------
    // Shared fleet contracts
    // ---------------------------------------------------------------------

    LaunchpadFactory internal factory;
    AgentToken internal agentTokenImpl;
    BondingCurve internal bondingCurveImpl; // sentinel; reserved-impl slot only
    Graduator internal graduator;
    FeeRouter internal feeRouter;
    AntiSniperModule internal antiSniper;

    MockMintableERC20 internal titu;
    MockUniswapV2Factory internal uniFactory;
    MockUniswapV2Router internal router;

    // ---------------------------------------------------------------------
    // Actors
    // ---------------------------------------------------------------------

    address internal owner = address(0xA0);
    address internal treasury = address(0x7);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    // ---------------------------------------------------------------------
    // Spec constants (mirror LaunchpadFactory immutables for clarity)
    // ---------------------------------------------------------------------

    uint256 internal constant VIRTUAL_QUOTE = 30_000e18;
    uint256 internal constant VIRTUAL_AGENT = 1_073_000_000e18;
    uint256 internal constant THRESHOLD = 42_000e18;
    uint256 internal constant INITIAL_AGENT = 1_000_000_000e18;

    // ---------------------------------------------------------------------
    // Local mirrors of the contract events for vm.expectEmit
    // ---------------------------------------------------------------------

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
    // Setup
    // ---------------------------------------------------------------------

    function setUp() public {
        // Shared TITU + Uni V2 mocks.
        titu = new MockMintableERC20("TITU", "TITU");
        uniFactory = new MockUniswapV2Factory();
        router = new MockUniswapV2Router(address(uniFactory));

        // AgentToken impl: cloned per launch, never initialized directly.
        agentTokenImpl = new AgentToken();

        // BondingCurve sentinel impl. The factory's `bondingCurveImpl` slot
        // is reserved for a future cloneable variant; today the launch path
        // still does `new BondingCurve(...)` per agent. We deploy a sentinel
        // here purely to satisfy the factory's non-zero check.
        bondingCurveImpl = new BondingCurve(
            address(this),
            address(uint160(uint256(keccak256("agent")))),
            address(uint160(uint256(keccak256("titu")))),
            address(uint160(uint256(keccak256("router")))),
            VIRTUAL_QUOTE,
            VIRTUAL_AGENT,
            THRESHOLD,
            INITIAL_AGENT
        );

        // FeeRouter: shared protocol-side fee sink. Owned by `owner` (Safe).
        feeRouter = new FeeRouter(treasury, owner);

        // Graduator: shared graduation sink. Constructed with `owner` so we
        // can transfer ownership to the factory after the factory is up.
        graduator = new Graduator(IUniswapV2Router02(address(router)), owner);

        // Factory wiring matches the production deploy script.
        factory = new LaunchpadFactory(
            address(titu),
            treasury,
            address(agentTokenImpl),
            address(bondingCurveImpl),
            address(feeRouter),
            address(graduator),
            address(uniFactory),
            address(router),
            owner
        );

        // Hand graduator ownership to the factory so the launch path can
        // call `registerLPLock`.
        vm.prank(owner);
        graduator.transferOwnership(address(factory));

        // AntiSniper module: owned by the factory so the factory's launch
        // path can call its `configure` exactly once per agent.
        antiSniper = new AntiSniperModule(address(factory));

        // Bind the AntiSniper module on the factory.
        bytes32 antiSniperId = factory.MODULE_ID_ANTI_SNIPER();
        vm.prank(owner);
        factory.setModule(antiSniperId, address(antiSniper));
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    function _bareParams(string memory name, string memory symbol)
        internal
        pure
        returns (LaunchpadFactory.LaunchParams memory p)
    {
        p.name = name;
        p.symbol = symbol;
        p.imageURI = "ipfs://image";
        p.soulURI = "ipfs://soul";
        p.enabledModules = new bytes32[](0);
        p.moduleData = new bytes[](0);
    }

    function _antiSniperParams(string memory name, string memory symbol, uint64 startTime)
        internal
        view
        returns (LaunchpadFactory.LaunchParams memory p)
    {
        p = _bareParams(name, symbol);
        bytes32[] memory ids = new bytes32[](1);
        ids[0] = factory.MODULE_ID_ANTI_SNIPER();
        bytes[] memory data = new bytes[](1);
        // 99% -> 1% over 60 seconds — matches the spec'd snipe window.
        data[0] = abi.encode(startTime, uint32(60), uint16(9900), uint16(100));
        p.enabledModules = ids;
        p.moduleData = data;
    }

    // ---------------------------------------------------------------------
    // launchAgent — happy-path full wiring
    // ---------------------------------------------------------------------

    function test_launch_clones_agent_token_via_eip1167() public {
        LaunchpadFactory.LaunchParams memory p = _bareParams("Alpha", "AGA");

        vm.prank(alice);
        (address agentToken, address curve, uint256 agentId) = factory.launchAgent(p);

        // EIP-1167 minimal proxy: 45-byte runtime, embedded impl address.
        bytes memory rt = agentToken.code;
        assertEq(rt.length, 45, "EIP-1167 clone runtime is 45 bytes");

        bytes20 impl;
        assembly {
            // First 10 bytes are the prefix; impl address sits at offset 10.
            impl := mload(add(add(rt, 0x20), 10))
        }
        assertEq(address(impl), address(agentTokenImpl), "clone impl pointer mismatch");

        // BondingCurve is a fresh deploy (not a clone) — its runtime is far
        // larger than 45 bytes.
        assertGt(curve.code.length, 45, "BondingCurve must not be a clone");

        // First-ever agent.
        assertEq(agentId, 1);
        assertEq(factory.agentCount(), 1);
    }

    function test_launch_wires_agent_token_state() public {
        vm.prank(alice);
        (address agentToken, address curve,) = factory.launchAgent(_bareParams("Alpha", "AGA"));

        AgentToken at = AgentToken(agentToken);

        // ERC-20 metadata + immutable wiring.
        assertEq(at.name(), "Alpha");
        assertEq(at.symbol(), "AGA");
        assertEq(at.creator(), alice);
        assertEq(at.feeRouter(), address(feeRouter));
        assertEq(at.bondingCurve(), curve);

        // Full 1B supply minted to the curve at initialization.
        assertEq(at.totalSupply(), at.TOTAL_SUPPLY());
        assertEq(at.balanceOf(curve), at.TOTAL_SUPPLY());
        assertEq(at.balanceOf(alice), 0);
    }

    function test_launch_wires_bonding_curve_state() public {
        vm.prank(alice);
        (address agentToken, address curve,) = factory.launchAgent(_bareParams("Alpha", "AGA"));

        BondingCurve bc = BondingCurve(curve);

        assertEq(bc.agentToken(), agentToken);
        assertEq(bc.quoteToken(), address(titu));
        assertEq(bc.feeRouter(), address(feeRouter));
        assertEq(bc.virtualQuoteReserve(), VIRTUAL_QUOTE);
        assertEq(bc.virtualAgentReserve(), VIRTUAL_AGENT);
        assertEq(bc.graduationThreshold(), THRESHOLD);
        assertEq(bc.initialAgentReserve(), INITIAL_AGENT);
        assertEq(bc.realAgentReserve(), INITIAL_AGENT);
        assertEq(bc.realQuoteReserve(), 0);
        assertFalse(bc.graduated());

        // Curve owner is the factory (one-shot setGraduator already fired).
        assertEq(bc.owner(), address(factory));
        assertEq(bc.graduator(), address(graduator));
    }

    function test_launch_wires_lp_lock() public {
        vm.prank(alice);
        (,, uint256 agentId) = factory.launchAgent(_bareParams("Alpha", "AGA"));

        LaunchpadFactory.Agent memory rec = factory.getAgent(agentId);

        assertTrue(rec.lpLock != address(0), "lpLock must be deployed");
        assertTrue(rec.pair != address(0), "pair must be pre-created");

        LPLock lock = LPLock(rec.lpLock);
        assertEq(address(lock.lpToken()), rec.pair, "LP token bound to pair");
        assertEq(lock.beneficiary(), alice, "beneficiary is the creator");
        assertFalse(lock.withdrawn());
        // 10-year lock anchored to launch block timestamp.
        assertEq(lock.unlockTime(), block.timestamp + 365 days * 10);
    }

    function test_launch_registers_lp_lock_on_graduator() public {
        vm.prank(alice);
        (, address curve,) = factory.launchAgent(_bareParams("Alpha", "AGA"));

        LaunchpadFactory.Agent memory rec = factory.getAgent(1);
        assertEq(graduator.lpLocks(curve), rec.lpLock, "graduator must know the lock");
    }

    function test_launch_pre_creates_uniswap_pair() public {
        // Before launch, the factory has not asked for any pair.
        assertEq(uniFactory.createPairCalls(), 0);

        vm.prank(alice);
        factory.launchAgent(_bareParams("Alpha", "AGA"));

        // After launch, exactly one pair was created.
        assertEq(uniFactory.createPairCalls(), 1);
        assertTrue(uniFactory.lastPair() != address(0));
    }

    // ---------------------------------------------------------------------
    // launchAgent — module dispatch (AntiSniper)
    // ---------------------------------------------------------------------

    function test_launch_with_anti_sniper_module_configures_decay() public {
        uint64 start = uint64(block.timestamp);
        LaunchpadFactory.LaunchParams memory p = _antiSniperParams("Alpha", "AGA", start);

        vm.prank(alice);
        (address agentToken,,) = factory.launchAgent(p);

        // The AntiSniper config is keyed on the agent token address. Verify
        // every field landed exactly as the launch payload specified.
        (uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps, bool configured) =
            antiSniper.configs(agentToken);
        assertTrue(configured, "AntiSniper config not written");
        assertEq(startTime, start);
        assertEq(uint256(duration), 60);
        assertEq(uint256(startBps), 9900);
        assertEq(uint256(endBps), 100);
    }

    function test_launch_with_anti_sniper_module_records_id_in_record() public {
        uint64 start = uint64(block.timestamp);
        LaunchpadFactory.LaunchParams memory p = _antiSniperParams("Alpha", "AGA", start);

        vm.prank(alice);
        factory.launchAgent(p);

        LaunchpadFactory.Agent memory rec = factory.getAgent(1);
        assertEq(rec.modules.length, 1, "exactly one module recorded");
        assertEq(rec.modules[0], factory.MODULE_ID_ANTI_SNIPER());
    }

    // ---------------------------------------------------------------------
    // launchAgent — events
    // ---------------------------------------------------------------------

    function test_launch_emits_AgentLaunched_with_full_payload() public {
        LaunchpadFactory.LaunchParams memory p = _bareParams("Alpha", "AGA");

        // We capture logs and decode by signature so we can assert the data
        // fields (lpLock / pair) without pre-computing addresses.
        vm.recordLogs();
        vm.prank(alice);
        (address agentToken, address curve, uint256 agentId) = factory.launchAgent(p);
        Vm.Log[] memory logs = vm.getRecordedLogs();

        bytes32 sig = keccak256("AgentLaunched(uint256,address,address,address,address,address,bytes32[])");
        bool found;
        for (uint256 i = 0; i < logs.length; ++i) {
            if (logs[i].emitter != address(factory)) continue;
            if (logs[i].topics.length == 0 || logs[i].topics[0] != sig) continue;

            // Indexed: agentId, token, curve.
            assertEq(uint256(logs[i].topics[1]), agentId);
            assertEq(address(uint160(uint256(logs[i].topics[2]))), agentToken);
            assertEq(address(uint160(uint256(logs[i].topics[3]))), curve);

            (address creator, address lpLock_, address pair_, bytes32[] memory mods) =
                abi.decode(logs[i].data, (address, address, address, bytes32[]));
            assertEq(creator, alice);
            assertEq(mods.length, 0, "no modules in bare launch");
            // Cross-check pair against the per-agent record.
            LaunchpadFactory.Agent memory rec = factory.getAgent(agentId);
            assertEq(lpLock_, rec.lpLock);
            assertEq(pair_, rec.pair);
            found = true;
            break;
        }
        assertTrue(found, "AgentLaunched event missing");
    }

    function test_launch_emits_AgentLaunched_with_module_id() public {
        uint64 start = uint64(block.timestamp);
        LaunchpadFactory.LaunchParams memory p = _antiSniperParams("Alpha", "AGA", start);

        vm.recordLogs();
        vm.prank(alice);
        factory.launchAgent(p);
        Vm.Log[] memory logs = vm.getRecordedLogs();

        bytes32 sig = keccak256("AgentLaunched(uint256,address,address,address,address,address,bytes32[])");
        for (uint256 i = 0; i < logs.length; ++i) {
            if (logs[i].emitter != address(factory)) continue;
            if (logs[i].topics.length == 0 || logs[i].topics[0] != sig) continue;
            (,,, bytes32[] memory mods) = abi.decode(logs[i].data, (address, address, address, bytes32[]));
            assertEq(mods.length, 1);
            assertEq(mods[0], factory.MODULE_ID_ANTI_SNIPER());
            return;
        }
        revert("AgentLaunched event missing");
    }

    // ---------------------------------------------------------------------
    // launchAgent — multiple launches share the fleet
    // ---------------------------------------------------------------------

    function test_two_launches_share_graduator_and_fee_router() public {
        vm.prank(alice);
        (address tokenA, address curveA,) = factory.launchAgent(_bareParams("Alpha", "AGA"));
        vm.prank(bob);
        (address tokenB, address curveB,) = factory.launchAgent(_bareParams("Beta", "AGB"));

        // Different per-agent artifacts.
        assertTrue(tokenA != tokenB, "token clones must be distinct");
        assertTrue(curveA != curveB, "curve deploys must be distinct");

        // Same shared fleet contracts.
        assertEq(BondingCurve(curveA).feeRouter(), BondingCurve(curveB).feeRouter());
        assertEq(BondingCurve(curveA).graduator(), BondingCurve(curveB).graduator());
        assertEq(AgentToken(tokenA).feeRouter(), AgentToken(tokenB).feeRouter());

        // Each agent has its own lock + pair.
        LaunchpadFactory.Agent memory recA = factory.getAgent(1);
        LaunchpadFactory.Agent memory recB = factory.getAgent(2);
        assertTrue(recA.lpLock != recB.lpLock);
        assertTrue(recA.pair != recB.pair);
    }
}
