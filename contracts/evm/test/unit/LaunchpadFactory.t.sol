// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test, Vm} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {LaunchpadFactory} from "../../src/launchpad/LaunchpadFactory.sol";
import {AgentToken} from "../../src/launchpad/AgentToken.sol";
import {BondingCurve} from "../../src/launchpad/BondingCurve.sol";
import {LPLock} from "../../src/launchpad/LPLock.sol";
import {Graduator} from "../../src/launchpad/Graduator.sol";

import {IUniswapV2Router02} from "../../src/launchpad/interfaces/IUniswapV2Router02.sol";

import {MockUniswapV2Factory} from "../mocks/MockUniswapV2Factory.sol";
import {MockUniswapV2Router} from "../mocks/MockUniswapV2Router.sol";
import {MockMintableERC20} from "../mocks/MockERC20.sol";

/// @dev Per-id mock module that records every `configure(...)` call so the
///      dispatch matrix can be asserted from tests. Each instance is bound to
///      a single canonical id so the Solidity ABI shape exactly matches the
///      typed interface the factory dispatches against. The mock asserts the
///      caller is the factory (Ownable-style) so we exercise the
///      "factory holds Ownable on every module" invariant.
contract MockConfigurableModule {
    address public immutable owner;
    bool public shouldRevert;
    bytes public lastSelector;
    uint256 public callCount;
    address public lastAgent;
    bytes public lastPayload;
    bool public reentryAttempted;
    address public reentrancyTarget;

    error NotOwner();
    error Boom();

    constructor(address owner_) {
        owner = owner_;
    }

    /// @dev Switch the next `configure*` call to revert with `Boom()`.
    function setShouldRevert(bool flag) external {
        shouldRevert = flag;
    }

    /// @dev Arm the module to attempt a reentrant `launchAgent` inside its
    ///      next configure call. Used to verify the factory's
    ///      `nonReentrant` guard.
    function armReentry(address factory) external {
        reentrancyTarget = factory;
    }

    function _onCall(string memory selector, address agent, bytes memory payload) private {
        if (msg.sender != owner) revert NotOwner();
        if (shouldRevert) revert Boom();
        if (reentrancyTarget != address(0)) {
            reentryAttempted = true;
            address t = reentrancyTarget;
            // Reset so we attempt only once even if the outer tx is retried.
            reentrancyTarget = address(0);
            // Build a minimal LaunchParams payload and call back into the
            // factory. The factory's `nonReentrant` guard MUST revert this.
            LaunchpadFactory.LaunchParams memory p = LaunchpadFactory.LaunchParams({
                name: "Reentry",
                symbol: "REE",
                imageURI: "",
                soulURI: "",
                enabledModules: new bytes32[](0),
                moduleData: new bytes[](0)
            });
            LaunchpadFactory(t).launchAgent(p);
        }
        lastSelector = bytes(selector);
        lastAgent = agent;
        lastPayload = payload;
        unchecked {
            ++callCount;
        }
    }

    // --- IAntiSniperModule ---
    function configure(address agent, uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps) external {
        _onCall("anti-sniper", agent, abi.encode(startTime, duration, startBps, endBps));
    }

    // --- ILaunchRadarModule ---
    function configure(address agent, uint64 startTime, bool whitelistOn, address agentAdmin) external {
        _onCall("launch-radar", agent, abi.encode(startTime, whitelistOn, agentAdmin));
    }

    // --- ISixtyDaysModule ---
    function configure(address agent, address escrowToken, address curve, address creator, uint64 startTime) external {
        _onCall("sixty-days", agent, abi.encode(escrowToken, curve, creator, startTime));
    }

    // --- IAirdropModule ---
    function configure(
        address agent,
        bytes32 root,
        uint64 snapshotBlock,
        uint128 allocation,
        uint64 deadline,
        address admin
    ) external {
        _onCall("airdrop", agent, abi.encode(root, snapshotBlock, allocation, deadline, admin));
    }

    // --- ICapitalFormationModule ---
    function configure(
        address agent,
        address agentToken_,
        address payoutToken,
        address pair,
        address creator,
        address agentAdmin,
        uint256[] calldata fdvThresholds,
        uint256[] calldata payoutsUsdc
    ) external {
        _onCall("capital-formation", agent, abi.encode(agentToken_, payoutToken, pair, creator, agentAdmin));
        // Silence unused params (event-encoded only as a smoke test).
        fdvThresholds;
        payoutsUsdc;
    }

    // --- IPreBuyModule ---
    function configure(
        address agent,
        address creator,
        uint256 vestAmount,
        uint64 cliffSeconds,
        uint64 durationSeconds,
        uint256 titanIn,
        address curve
    ) external returns (address clone) {
        _onCall("pre-buy", agent, abi.encode(creator, vestAmount, cliffSeconds, durationSeconds, titanIn, curve));
        return address(this); // sentinel — never read by the factory
    }

    // --- IExistingTokenModule ---
    function configure(address externalToken, address curve, uint256 supply, address admin) external {
        _onCall("existing-token", externalToken, abi.encode(curve, supply, admin));
    }
}

contract LaunchpadFactoryTest is Test {
    // ---------------------------------------------------------------------
    // Fixtures
    // ---------------------------------------------------------------------
    LaunchpadFactory internal factory;
    AgentToken internal agentTokenImpl;
    BondingCurve internal bondingCurveImpl; // sentinel; not consumed by launch path
    Graduator internal graduator;
    MockMintableERC20 internal titu;
    MockUniswapV2Factory internal uniFactory;
    MockUniswapV2Router internal router;

    address internal owner = address(0xA0);
    address internal treasury = address(0x7);
    address internal feeRouter = address(0xFEE);
    address internal alice = address(0xA11CE);
    address internal bob = address(0xB0B);

    // ---------------------------------------------------------------------
    // Events (mirror the contract for vm.expectEmit)
    // ---------------------------------------------------------------------
    event ModuleSet(bytes32 indexed moduleId, address indexed prev, address indexed module);
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
        agentTokenImpl = new AgentToken();
        // Deploy a sentinel bonding-curve impl. Using address(this) as owner
        // satisfies the constructor's non-zero check; the slot is reserved on
        // the factory and never invoked by launch.
        bondingCurveImpl = new BondingCurve(
            address(this),
            address(uint160(uint256(keccak256("agent")))),
            address(uint160(uint256(keccak256("titu")))),
            address(uint160(uint256(keccak256("router")))),
            30_000e18,
            1_073_000_000e18,
            42_000e18,
            1_000_000_000e18
        );

        titu = new MockMintableERC20("TITU", "TITU");
        uniFactory = new MockUniswapV2Factory();
        router = new MockUniswapV2Router(address(uniFactory));

        // Graduator owned by the factory-to-be: easiest to construct first
        // with a placeholder, then transfer ownership after the factory is
        // up. We use the more direct path: pre-compute the factory address
        // is unnecessary — the brief explicitly says "deploy script
        // transfers ownership". So `owner` deploys the graduator and then
        // hands ownership off in the test path.
        graduator = new Graduator(IUniswapV2Router02(address(router)), owner);

        factory = new LaunchpadFactory(
            address(titu),
            treasury,
            address(agentTokenImpl),
            address(bondingCurveImpl),
            feeRouter,
            address(graduator),
            address(uniFactory),
            address(router),
            owner
        );

        // Hand graduator ownership to the factory so the launch path can
        // call `registerLPLock`.
        vm.prank(owner);
        graduator.transferOwnership(address(factory));
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

    function _setMockModule(bytes32 id) internal returns (MockConfigurableModule m) {
        m = new MockConfigurableModule(address(factory));
        vm.prank(owner);
        factory.setModule(id, address(m));
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_sets_state() public view {
        assertEq(factory.tituToken(), address(titu));
        assertEq(factory.treasury(), treasury);
        assertEq(factory.agentTokenImpl(), address(agentTokenImpl));
        assertEq(factory.bondingCurveImpl(), address(bondingCurveImpl));
        assertEq(factory.feeRouter(), feeRouter);
        assertEq(address(factory.graduator()), address(graduator));
        assertEq(address(factory.uniV2Factory()), address(uniFactory));
        assertEq(address(factory.uniV2Router()), address(router));
        assertEq(factory.owner(), owner);
        assertEq(factory.agentCount(), 0);
        assertEq(factory.GRADUATION_THRESHOLD(), 42_000e18);
        assertEq(factory.MAX_NAME_LEN(), 64);
        assertEq(factory.MAX_SYMBOL_LEN(), 16);
    }

    function test_constructor_revert_zero_titu() public {
        vm.expectRevert(LaunchpadFactory.ZeroAddress.selector);
        new LaunchpadFactory(
            address(0),
            treasury,
            address(agentTokenImpl),
            address(bondingCurveImpl),
            feeRouter,
            address(graduator),
            address(uniFactory),
            address(router),
            owner
        );
    }

    function test_constructor_revert_zero_treasury() public {
        vm.expectRevert(LaunchpadFactory.ZeroAddress.selector);
        new LaunchpadFactory(
            address(titu),
            address(0),
            address(agentTokenImpl),
            address(bondingCurveImpl),
            feeRouter,
            address(graduator),
            address(uniFactory),
            address(router),
            owner
        );
    }

    function test_constructor_revert_zero_agent_token_impl() public {
        vm.expectRevert(LaunchpadFactory.ZeroAddress.selector);
        new LaunchpadFactory(
            address(titu),
            treasury,
            address(0),
            address(bondingCurveImpl),
            feeRouter,
            address(graduator),
            address(uniFactory),
            address(router),
            owner
        );
    }

    function test_constructor_revert_zero_bonding_curve_impl() public {
        vm.expectRevert(LaunchpadFactory.ZeroAddress.selector);
        new LaunchpadFactory(
            address(titu),
            treasury,
            address(agentTokenImpl),
            address(0),
            feeRouter,
            address(graduator),
            address(uniFactory),
            address(router),
            owner
        );
    }

    function test_constructor_revert_zero_fee_router() public {
        vm.expectRevert(LaunchpadFactory.ZeroAddress.selector);
        new LaunchpadFactory(
            address(titu),
            treasury,
            address(agentTokenImpl),
            address(bondingCurveImpl),
            address(0),
            address(graduator),
            address(uniFactory),
            address(router),
            owner
        );
    }

    function test_constructor_revert_zero_graduator() public {
        vm.expectRevert(LaunchpadFactory.ZeroAddress.selector);
        new LaunchpadFactory(
            address(titu),
            treasury,
            address(agentTokenImpl),
            address(bondingCurveImpl),
            feeRouter,
            address(0),
            address(uniFactory),
            address(router),
            owner
        );
    }

    function test_constructor_revert_zero_uni_factory() public {
        vm.expectRevert(LaunchpadFactory.ZeroAddress.selector);
        new LaunchpadFactory(
            address(titu),
            treasury,
            address(agentTokenImpl),
            address(bondingCurveImpl),
            feeRouter,
            address(graduator),
            address(0),
            address(router),
            owner
        );
    }

    function test_constructor_revert_zero_uni_router() public {
        vm.expectRevert(LaunchpadFactory.ZeroAddress.selector);
        new LaunchpadFactory(
            address(titu),
            treasury,
            address(agentTokenImpl),
            address(bondingCurveImpl),
            feeRouter,
            address(graduator),
            address(uniFactory),
            address(0),
            owner
        );
    }

    function test_constructor_revert_zero_owner() public {
        vm.expectRevert(LaunchpadFactory.ZeroAddress.selector);
        new LaunchpadFactory(
            address(titu),
            treasury,
            address(agentTokenImpl),
            address(bondingCurveImpl),
            feeRouter,
            address(graduator),
            address(uniFactory),
            address(router),
            address(0)
        );
    }

    // ---------------------------------------------------------------------
    // setModule
    // ---------------------------------------------------------------------

    function test_setModule_happy_path_emits_event() public {
        MockConfigurableModule m = new MockConfigurableModule(address(factory));
        bytes32 id = factory.MODULE_ID_ANTI_SNIPER();

        vm.expectEmit(true, true, true, true, address(factory));
        emit ModuleSet(id, address(0), address(m));

        vm.prank(owner);
        factory.setModule(id, address(m));

        assertEq(factory.modules(id), address(m));
    }

    function test_setModule_overwrite_emits_prev() public {
        MockConfigurableModule a = new MockConfigurableModule(address(factory));
        MockConfigurableModule b = new MockConfigurableModule(address(factory));
        bytes32 id = factory.MODULE_ID_ANTI_SNIPER();

        vm.prank(owner);
        factory.setModule(id, address(a));

        vm.expectEmit(true, true, true, true, address(factory));
        emit ModuleSet(id, address(a), address(b));

        vm.prank(owner);
        factory.setModule(id, address(b));

        assertEq(factory.modules(id), address(b));
    }

    function test_setModule_revert_when_not_owner() public {
        MockConfigurableModule m = new MockConfigurableModule(address(factory));
        bytes32 id = factory.MODULE_ID_ANTI_SNIPER();
        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        factory.setModule(id, address(m));
    }

    function test_setModule_revert_zero_address() public {
        bytes32 id = factory.MODULE_ID_ANTI_SNIPER();
        vm.prank(owner);
        vm.expectRevert(LaunchpadFactory.ZeroAddress.selector);
        factory.setModule(id, address(0));
    }

    function test_setModule_revert_zero_id() public {
        MockConfigurableModule m = new MockConfigurableModule(address(factory));
        vm.prank(owner);
        vm.expectRevert(LaunchpadFactory.InvalidModuleId.selector);
        factory.setModule(bytes32(0), address(m));
    }

    // ---------------------------------------------------------------------
    // launchAgent — happy paths
    // ---------------------------------------------------------------------

    function test_launchAgent_no_modules_wires_state() public {
        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");

        vm.prank(alice);
        (address agentToken, address curve, uint256 agentId) = factory.launchAgent(p);

        assertEq(agentId, 1);
        assertEq(factory.agentCount(), 1);

        // AgentToken state.
        AgentToken at = AgentToken(agentToken);
        assertEq(at.name(), "Agent A");
        assertEq(at.symbol(), "AGA");
        assertEq(at.creator(), alice);
        assertEq(at.feeRouter(), feeRouter);
        assertEq(at.bondingCurve(), curve);
        assertEq(at.totalSupply(), at.TOTAL_SUPPLY());
        assertEq(at.balanceOf(curve), at.TOTAL_SUPPLY());

        // BondingCurve state.
        BondingCurve bc = BondingCurve(curve);
        assertEq(bc.agentToken(), agentToken);
        assertEq(bc.quoteToken(), address(titu));
        assertEq(bc.feeRouter(), feeRouter);
        assertEq(bc.graduationThreshold(), 42_000e18);
        assertEq(bc.graduator(), address(graduator));
        assertEq(bc.realAgentReserve(), 1_000_000_000e18);

        // Per-agent record.
        LaunchpadFactory.Agent memory rec = factory.getAgent(1);
        assertEq(rec.token, agentToken);
        assertEq(rec.curve, curve);
        assertEq(rec.creator, alice);
        assertTrue(rec.lpLock != address(0));
        assertTrue(rec.pair != address(0));
        assertEq(rec.modules.length, 0);

        // LPLock binds the pair as its LP token, beneficiary = creator.
        LPLock lock = LPLock(rec.lpLock);
        assertEq(address(lock.lpToken()), rec.pair);
        assertEq(lock.beneficiary(), alice);

        // Graduator wired to the right lock.
        assertEq(graduator.lpLocks(curve), rec.lpLock);
    }

    function test_launchAgent_emits_event() public {
        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");

        // We only assert the indexed fields strictly because the address
        // values for `lpLock` / `pair` aren't known until launch.
        vm.recordLogs();
        vm.prank(alice);
        (address agentToken, address curve, uint256 agentId) = factory.launchAgent(p);
        Vm.Log[] memory logs = vm.getRecordedLogs();

        bytes32 sig = keccak256("AgentLaunched(uint256,address,address,address,address,address,bytes32[])");
        bool found;
        for (uint256 i = 0; i < logs.length; ++i) {
            if (logs[i].topics.length == 0) continue;
            if (logs[i].topics[0] != sig) continue;
            // topics: [sig, agentId, token, curve]
            assertEq(uint256(logs[i].topics[1]), agentId);
            assertEq(address(uint160(uint256(logs[i].topics[2]))), agentToken);
            assertEq(address(uint160(uint256(logs[i].topics[3]))), curve);
            // data: (creator, lpLock, pair, modules)
            (address creator, address lpLock_, address pair_, bytes32[] memory mods) =
                abi.decode(logs[i].data, (address, address, address, bytes32[]));
            assertEq(creator, alice);
            assertEq(mods.length, 0);
            assertGt(uint160(lpLock_), 0);
            assertGt(uint160(pair_), 0);
            found = true;
            break;
        }
        assertTrue(found, "AgentLaunched event missing");
    }

    function test_launchAgent_two_consecutive_increment_id() public {
        LaunchpadFactory.LaunchParams memory p1 = _bareParams("Agent A", "AGA");
        LaunchpadFactory.LaunchParams memory p2 = _bareParams("Agent B", "AGB");

        vm.prank(alice);
        (,, uint256 id1) = factory.launchAgent(p1);
        vm.prank(bob);
        (,, uint256 id2) = factory.launchAgent(p2);

        assertEq(id1, 1);
        assertEq(id2, 2);
        assertEq(factory.agentCount(), 2);

        LaunchpadFactory.Agent memory rec1 = factory.getAgent(1);
        LaunchpadFactory.Agent memory rec2 = factory.getAgent(2);
        assertEq(rec1.creator, alice);
        assertEq(rec2.creator, bob);
        assertTrue(rec1.token != rec2.token, "tokens must differ");
        assertTrue(rec1.curve != rec2.curve, "curves must differ");
    }

    function test_launchAgent_creates_eip1167_clone_with_45_byte_runtime() public {
        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");
        vm.prank(alice);
        (address agentToken,,) = factory.launchAgent(p);

        bytes memory rt = agentToken.code;
        assertEq(rt.length, 45, "EIP-1167 clone runtime must be 45 bytes");

        // EIP-1167 minimal proxy bytes:
        //   3d602d80600a3d3981f3 363d3d373d3d3d363d73 <impl 20 bytes> 5af43d82803e903d91602b57fd5bf3
        // Verify the prefix and suffix; the middle 20 bytes are the impl.
        bytes memory prefix = hex"363d3d373d3d3d363d73";
        bytes memory suffix = hex"5af43d82803e903d91602b57fd5bf3";
        // Skip the leading constructor (10 bytes) since we read runtime, not creation.
        for (uint256 i = 0; i < prefix.length; ++i) {
            assertEq(rt[i], prefix[i]);
        }
        bytes20 impl;
        assembly {
            // rt[10..30] is the impl address. rt is `bytes memory`, so its
            // payload starts at offset 0x20 from the pointer.
            impl := mload(add(add(rt, 0x20), 10))
        }
        assertEq(address(impl), address(agentTokenImpl));
        for (uint256 i = 0; i < suffix.length; ++i) {
            assertEq(rt[30 + i], suffix[i]);
        }
    }

    // ---------------------------------------------------------------------
    // launchAgent — module dispatch
    // ---------------------------------------------------------------------

    function test_launchAgent_with_one_module_dispatches() public {
        MockConfigurableModule m = _setMockModule(factory.MODULE_ID_ANTI_SNIPER());

        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");
        bytes32[] memory ids = new bytes32[](1);
        ids[0] = factory.MODULE_ID_ANTI_SNIPER();
        bytes[] memory data = new bytes[](1);
        data[0] = abi.encode(uint64(100), uint32(60), uint16(9900), uint16(100));
        p.enabledModules = ids;
        p.moduleData = data;

        vm.prank(alice);
        (address agentToken,,) = factory.launchAgent(p);

        assertEq(m.callCount(), 1);
        assertEq(m.lastAgent(), agentToken);
        assertEq(string(m.lastSelector()), "anti-sniper");
    }

    function test_launchAgent_dispatch_all_seven_modules() public {
        // Bind one mock per id and dispatch all seven in a single launch.
        MockConfigurableModule mAS = _setMockModule(factory.MODULE_ID_ANTI_SNIPER());
        MockConfigurableModule mLR = _setMockModule(factory.MODULE_ID_LAUNCH_RADAR());
        MockConfigurableModule m60 = _setMockModule(factory.MODULE_ID_SIXTY_DAYS());
        MockConfigurableModule mAD = _setMockModule(factory.MODULE_ID_AIRDROP());
        MockConfigurableModule mCF = _setMockModule(factory.MODULE_ID_CAPITAL_FORMATION());
        MockConfigurableModule mPB = _setMockModule(factory.MODULE_ID_PRE_BUY());
        MockConfigurableModule mET = _setMockModule(factory.MODULE_ID_EXISTING_TOKEN());

        bytes32[] memory ids = new bytes32[](7);
        ids[0] = factory.MODULE_ID_ANTI_SNIPER();
        ids[1] = factory.MODULE_ID_LAUNCH_RADAR();
        ids[2] = factory.MODULE_ID_SIXTY_DAYS();
        ids[3] = factory.MODULE_ID_AIRDROP();
        ids[4] = factory.MODULE_ID_CAPITAL_FORMATION();
        ids[5] = factory.MODULE_ID_PRE_BUY();
        ids[6] = factory.MODULE_ID_EXISTING_TOKEN();

        bytes[] memory data = new bytes[](7);
        data[0] = abi.encode(uint64(100), uint32(60), uint16(9900), uint16(100));
        data[1] = abi.encode(uint64(200), true);
        data[2] = abi.encode(address(titu), uint64(300));
        data[3] = abi.encode(bytes32(uint256(0xABCD)), uint64(123), uint128(5_000_000e18), uint64(456));
        uint256[] memory thresholds = new uint256[](2);
        thresholds[0] = 2_000_000e18;
        thresholds[1] = 10_000_000e18;
        uint256[] memory payouts = new uint256[](2);
        payouts[0] = 100_000e6;
        payouts[1] = 500_000e6;
        data[4] = abi.encode(address(titu), address(0xCAFE), thresholds, payouts);
        data[5] = abi.encode(uint256(50_000e18), uint64(60), uint64(86_400), uint256(100e18));
        data[6] = abi.encode(address(0xEEEE), address(0xDDDD), uint256(10_000e18));

        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");
        p.enabledModules = ids;
        p.moduleData = data;

        vm.prank(alice);
        factory.launchAgent(p);

        assertEq(mAS.callCount(), 1);
        assertEq(mLR.callCount(), 1);
        assertEq(m60.callCount(), 1);
        assertEq(mAD.callCount(), 1);
        assertEq(mCF.callCount(), 1);
        assertEq(mPB.callCount(), 1);
        assertEq(mET.callCount(), 1);

        assertEq(string(mAS.lastSelector()), "anti-sniper");
        assertEq(string(mLR.lastSelector()), "launch-radar");
        assertEq(string(m60.lastSelector()), "sixty-days");
        assertEq(string(mAD.lastSelector()), "airdrop");
        assertEq(string(mCF.lastSelector()), "capital-formation");
        assertEq(string(mPB.lastSelector()), "pre-buy");
        assertEq(string(mET.lastSelector()), "existing-token");
    }

    function test_launchAgent_module_revert_propagates() public {
        MockConfigurableModule m = _setMockModule(factory.MODULE_ID_ANTI_SNIPER());
        m.setShouldRevert(true);

        bytes32[] memory ids = new bytes32[](1);
        ids[0] = factory.MODULE_ID_ANTI_SNIPER();
        bytes[] memory data = new bytes[](1);
        data[0] = abi.encode(uint64(0), uint32(1), uint16(0), uint16(0));

        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");
        p.enabledModules = ids;
        p.moduleData = data;

        vm.prank(alice);
        vm.expectRevert(MockConfigurableModule.Boom.selector);
        factory.launchAgent(p);
        // Counter should not advance on revert.
        assertEq(factory.agentCount(), 0);
    }

    function test_launchAgent_unknown_module_revert() public {
        bytes32 weirdId = keccak256("never-registered");
        bytes32[] memory ids = new bytes32[](1);
        ids[0] = weirdId;
        bytes[] memory data = new bytes[](1);
        data[0] = bytes("");

        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");
        p.enabledModules = ids;
        p.moduleData = data;

        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSelector(LaunchpadFactory.UnknownModule.selector, weirdId));
        factory.launchAgent(p);
    }

    function test_launchAgent_duplicate_module_revert() public {
        _setMockModule(factory.MODULE_ID_ANTI_SNIPER());

        bytes32 id = factory.MODULE_ID_ANTI_SNIPER();
        bytes32[] memory ids = new bytes32[](2);
        ids[0] = id;
        ids[1] = id;
        bytes[] memory data = new bytes[](2);
        data[0] = abi.encode(uint64(0), uint32(1), uint16(0), uint16(0));
        data[1] = abi.encode(uint64(0), uint32(1), uint16(0), uint16(0));

        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");
        p.enabledModules = ids;
        p.moduleData = data;

        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSelector(LaunchpadFactory.DuplicateModule.selector, id));
        factory.launchAgent(p);
    }

    function test_launchAgent_length_mismatch_revert() public {
        _setMockModule(factory.MODULE_ID_ANTI_SNIPER());

        bytes32[] memory ids = new bytes32[](1);
        ids[0] = factory.MODULE_ID_ANTI_SNIPER();
        bytes[] memory data = new bytes[](2);
        data[0] = bytes("");
        data[1] = bytes("");

        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");
        p.enabledModules = ids;
        p.moduleData = data;

        vm.prank(alice);
        vm.expectRevert(LaunchpadFactory.LengthMismatch.selector);
        factory.launchAgent(p);
    }

    // ---------------------------------------------------------------------
    // launchAgent — input validation
    // ---------------------------------------------------------------------

    function test_launchAgent_revert_empty_name() public {
        LaunchpadFactory.LaunchParams memory p = _bareParams("", "AGA");
        vm.prank(alice);
        vm.expectRevert(LaunchpadFactory.EmptyName.selector);
        factory.launchAgent(p);
    }

    function test_launchAgent_revert_empty_symbol() public {
        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "");
        vm.prank(alice);
        vm.expectRevert(LaunchpadFactory.EmptySymbol.selector);
        factory.launchAgent(p);
    }

    function test_launchAgent_revert_name_too_long() public {
        // 65 bytes — exceeds MAX_NAME_LEN (64).
        bytes memory longName = new bytes(65);
        for (uint256 i = 0; i < 65; ++i) {
            longName[i] = "a";
        }
        LaunchpadFactory.LaunchParams memory p = _bareParams(string(longName), "AGA");
        vm.prank(alice);
        vm.expectRevert(LaunchpadFactory.NameTooLong.selector);
        factory.launchAgent(p);
    }

    function test_launchAgent_accepts_max_name_len() public {
        bytes memory n = new bytes(64);
        for (uint256 i = 0; i < 64; ++i) {
            n[i] = "a";
        }
        LaunchpadFactory.LaunchParams memory p = _bareParams(string(n), "AGA");
        vm.prank(alice);
        (,, uint256 id) = factory.launchAgent(p);
        assertEq(id, 1);
    }

    function test_launchAgent_revert_symbol_too_long() public {
        bytes memory s = new bytes(17);
        for (uint256 i = 0; i < 17; ++i) {
            s[i] = "x";
        }
        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", string(s));
        vm.prank(alice);
        vm.expectRevert(LaunchpadFactory.SymbolTooLong.selector);
        factory.launchAgent(p);
    }

    function test_launchAgent_accepts_max_symbol_len() public {
        bytes memory s = new bytes(16);
        for (uint256 i = 0; i < 16; ++i) {
            s[i] = "x";
        }
        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", string(s));
        vm.prank(alice);
        (,, uint256 id) = factory.launchAgent(p);
        assertEq(id, 1);
    }

    // ---------------------------------------------------------------------
    // launchAgent — reentrancy
    // ---------------------------------------------------------------------

    function test_launchAgent_reentrancy_blocked() public {
        // Bind the mock and arm it to call launchAgent inside its configure.
        MockConfigurableModule m = _setMockModule(factory.MODULE_ID_ANTI_SNIPER());
        m.armReentry(address(factory));

        bytes32[] memory ids = new bytes32[](1);
        ids[0] = factory.MODULE_ID_ANTI_SNIPER();
        bytes[] memory data = new bytes[](1);
        data[0] = abi.encode(uint64(0), uint32(1), uint16(0), uint16(0));

        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");
        p.enabledModules = ids;
        p.moduleData = data;

        // OZ ReentrancyGuard reverts with ReentrancyGuardReentrantCall().
        vm.prank(alice);
        vm.expectRevert(); // tolerate either OZ selector or wrapper
        factory.launchAgent(p);
    }

    // ---------------------------------------------------------------------
    // Views
    // ---------------------------------------------------------------------

    function test_getAgent_revert_zero() public {
        vm.expectRevert(abi.encodeWithSelector(LaunchpadFactory.UnknownAgent.selector, uint256(0)));
        factory.getAgent(0);
    }

    function test_getAgent_revert_out_of_range() public {
        vm.expectRevert(abi.encodeWithSelector(LaunchpadFactory.UnknownAgent.selector, uint256(99)));
        factory.getAgent(99);
    }

    function test_agentCount_starts_zero() public view {
        assertEq(factory.agentCount(), 0);
    }

    // ---------------------------------------------------------------------
    // launchAgent — propagation of inner reverts
    // ---------------------------------------------------------------------

    function test_launchAgent_revert_when_factory_not_graduator_owner() public {
        // Re-deploy a graduator NOT owned by the factory and a new factory
        // pointing at it. Launch should revert at registerLPLock.
        Graduator orphan = new Graduator(IUniswapV2Router02(address(router)), owner);
        LaunchpadFactory f2 = new LaunchpadFactory(
            address(titu),
            treasury,
            address(agentTokenImpl),
            address(bondingCurveImpl),
            feeRouter,
            address(orphan),
            address(uniFactory),
            address(router),
            owner
        );

        LaunchpadFactory.LaunchParams memory p = _bareParams("Agent A", "AGA");
        vm.prank(alice);
        // OZ Ownable reverts with OwnableUnauthorizedAccount(address).
        vm.expectRevert();
        f2.launchAgent(p);
    }
}
