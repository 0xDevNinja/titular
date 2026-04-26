// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";

import {DeployPhase2} from "../../script/DeployPhase2.s.sol";

import {LaunchpadFactory} from "../../src/launchpad/LaunchpadFactory.sol";
import {Graduator} from "../../src/launchpad/Graduator.sol";
import {AgentToken} from "../../src/launchpad/AgentToken.sol";
import {FeeRouter} from "../../src/launchpad/FeeRouter.sol";

import {MockUniswapV2Factory} from "../mocks/MockUniswapV2Factory.sol";
import {MockUniswapV2Router} from "../mocks/MockUniswapV2Router.sol";

/// @title DeployPhase2 fixture-driven tests.
/// @notice Exercises the full Phase-2 deploy script against test-controlled
///         Phase-1 JSON fixtures and mocked Uniswap V2 endpoints. Pinned to
///         `test/script/__phase2_*.json` so the canonical
///         `deployments/<chain>.json` is never touched.
///
///         The script exposes three entry points:
///           - {DeployPhase2.run}           — production; uses env vars.
///           - {DeployPhase2.runWithPath}   — path override; uses env vars.
///           - {DeployPhase2.runWith}       — fully explicit; bypasses env.
///
///         Tests use {runWith} for the happy paths because Forge runs tests
///         in parallel by default and `vm.setEnv` mutates a process-shared
///         env table — racy across cases. {runWithPath} is reserved for the
///         env-revert tests where reading the env IS the system under test.
contract DeployPhase2Test is Test {
    DeployPhase2 internal deployer;

    MockUniswapV2Factory internal uniFactory;
    MockUniswapV2Router internal uniRouter;

    /// @dev Fixture paths. Each test owns its own file so parallel-run
    ///      Forge does not race on the same fixture.
    string internal constant PATH_TOPOLOGY = "test/script/__phase2_topology.json";
    string internal constant PATH_IDEMPOTENT = "test/script/__phase2_idempotent.json";
    string internal constant PATH_REVERT_PRIOR = "test/script/__phase2_revert-prior.json";
    string internal constant PATH_REVERT_UNI = "test/script/__phase2_revert-uni.json";
    string internal constant PATH_MERGE = "test/script/__phase2_merged.json";
    string internal constant PATH_LEGACY = "test/script/__phase2_legacy.json";
    string internal constant PATH_REVERT_NULL = "test/script/__phase2_revert-null.json";

    /// @dev Sentinel Phase-1 addresses written to the fixture JSONs. Only
    ///      their non-zeroness matters to the script — using deterministic
    ///      hashes keeps assertions readable in trace output.
    address internal constant FIXTURE_TITU = address(0x1111111111111111111111111111111111111111);
    address internal constant FIXTURE_TREASURY = address(0x2222222222222222222222222222222222222222);
    address internal constant FIXTURE_TREASURY_IMPL = address(0x3333333333333333333333333333333333333333);
    address internal constant FIXTURE_VESTING_VAULT = address(0x4444444444444444444444444444444444444444);
    address internal constant FIXTURE_VE_TITU = address(0x5555555555555555555555555555555555555555);
    address internal constant FIXTURE_FEE_DISTRIBUTOR = address(0x6666666666666666666666666666666666666666);

    function setUp() public {
        deployer = new DeployPhase2();

        // Stand in for the real V2 factory / router. The graduator constructor
        // calls `router.factory()` so we need a wired pair, not bare addresses.
        uniFactory = new MockUniswapV2Factory();
        uniRouter = new MockUniswapV2Router(address(uniFactory));

        // Clean any state left behind by a previous run so we never read a
        // stale fixture and silently pass.
        _wipe(PATH_TOPOLOGY);
        _wipe(PATH_IDEMPOTENT);
        _wipe(PATH_REVERT_PRIOR);
        _wipe(PATH_REVERT_UNI);
        _wipe(PATH_MERGE);
        _wipe(PATH_LEGACY);
        _wipe(PATH_REVERT_NULL);
    }

    function _wipe(string memory path) internal {
        if (vm.exists(path)) vm.removeFile(path);
    }

    // -----------------------------------------------------------------------
    // Test — full topology
    // -----------------------------------------------------------------------

    /// @dev Confirms every Phase-2 contract is deployed, wired, and surfaced
    ///      on the returned struct. Also asserts the Graduator owner is the
    ///      LaunchpadFactory and every module slot is registered.
    function test_topology_deploysAndWiresEverything() public {
        _writePhase1(PATH_TOPOLOGY);

        DeployPhase2.Phase2Deployment memory d =
            deployer.runWith(PATH_TOPOLOGY, address(uniFactory), address(uniRouter));

        // Every slot non-zero.
        assertTrue(d.agentTokenImpl != address(0), "agentTokenImpl");
        assertTrue(d.bondingCurveImpl != address(0), "bondingCurveImpl");
        assertTrue(d.feeRouter != address(0), "feeRouter");
        assertTrue(d.graduator != address(0), "graduator");
        assertTrue(d.launchpadFactory != address(0), "launchpadFactory");
        assertTrue(d.antiSniper != address(0), "antiSniper");
        assertTrue(d.launchRadar != address(0), "launchRadar");
        assertTrue(d.sixtyDays != address(0), "sixtyDays");
        assertTrue(d.airdrop != address(0), "airdrop");
        assertTrue(d.capitalFormation != address(0), "capitalFormation");
        assertTrue(d.preBuy != address(0), "preBuy");
        assertTrue(d.existingToken != address(0), "existingToken");

        // Factory carries the Phase-1 references and the Phase-2 wiring.
        LaunchpadFactory factory = LaunchpadFactory(d.launchpadFactory);
        assertEq(factory.tituToken(), FIXTURE_TITU, "factory.titu");
        assertEq(factory.treasury(), FIXTURE_TREASURY, "factory.treasury");
        assertEq(factory.agentTokenImpl(), d.agentTokenImpl, "factory.agentTokenImpl");
        assertEq(factory.bondingCurveImpl(), d.bondingCurveImpl, "factory.bondingCurveImpl");
        assertEq(factory.feeRouter(), d.feeRouter, "factory.feeRouter");
        assertEq(address(factory.graduator()), d.graduator, "factory.graduator");
        assertEq(address(factory.uniV2Factory()), address(uniFactory), "factory.uniFactory");
        assertEq(address(factory.uniV2Router()), address(uniRouter), "factory.uniRouter");

        // Module registry — every canonical id resolves to its module.
        assertEq(factory.modules(factory.MODULE_ID_ANTI_SNIPER()), d.antiSniper, "module.antiSniper");
        assertEq(factory.modules(factory.MODULE_ID_LAUNCH_RADAR()), d.launchRadar, "module.launchRadar");
        assertEq(factory.modules(factory.MODULE_ID_SIXTY_DAYS()), d.sixtyDays, "module.sixtyDays");
        assertEq(factory.modules(factory.MODULE_ID_AIRDROP()), d.airdrop, "module.airdrop");
        assertEq(factory.modules(factory.MODULE_ID_CAPITAL_FORMATION()), d.capitalFormation, "module.capitalFormation");
        assertEq(factory.modules(factory.MODULE_ID_PRE_BUY()), d.preBuy, "module.preBuy");
        assertEq(factory.modules(factory.MODULE_ID_EXISTING_TOKEN()), d.existingToken, "module.existingToken");

        // Graduator ownership transferred to the factory so launches can
        // invoke `registerLPLock`.
        Graduator graduator = Graduator(d.graduator);
        assertEq(graduator.owner(), d.launchpadFactory, "graduator.owner");

        // FeeRouter wired to the Phase-1 treasury.
        FeeRouter fr = FeeRouter(payable(d.feeRouter));
        assertEq(fr.treasury(), FIXTURE_TREASURY, "feeRouter.treasury");

        // AgentToken impl is initializer-disabled, confirming we deployed the
        // master copy and not a clone.
        AgentToken impl = AgentToken(d.agentTokenImpl);
        vm.expectRevert();
        impl.initialize("X", "X", address(0xCAFE), d.feeRouter, address(0xBABE), d.graduator);
    }

    // -----------------------------------------------------------------------
    // Test — idempotent re-run
    // -----------------------------------------------------------------------

    /// @dev A second `runWith` against the same JSON must reuse every
    ///      Phase-2 address. We assert each slot is byte-for-byte identical
    ///      across runs.
    function test_idempotent_reusesEveryAddress() public {
        _writePhase1(PATH_IDEMPOTENT);

        DeployPhase2.Phase2Deployment memory first =
            deployer.runWith(PATH_IDEMPOTENT, address(uniFactory), address(uniRouter));
        DeployPhase2.Phase2Deployment memory second =
            deployer.runWith(PATH_IDEMPOTENT, address(uniFactory), address(uniRouter));

        assertEq(first.agentTokenImpl, second.agentTokenImpl, "agentTokenImpl");
        assertEq(first.bondingCurveImpl, second.bondingCurveImpl, "bondingCurveImpl");
        assertEq(first.feeRouter, second.feeRouter, "feeRouter");
        assertEq(first.graduator, second.graduator, "graduator");
        assertEq(first.launchpadFactory, second.launchpadFactory, "launchpadFactory");
        assertEq(first.antiSniper, second.antiSniper, "antiSniper");
        assertEq(first.launchRadar, second.launchRadar, "launchRadar");
        assertEq(first.sixtyDays, second.sixtyDays, "sixtyDays");
        assertEq(first.airdrop, second.airdrop, "airdrop");
        assertEq(first.capitalFormation, second.capitalFormation, "capitalFormation");
        assertEq(first.preBuy, second.preBuy, "preBuy");
        assertEq(first.existingToken, second.existingToken, "existingToken");

        // Graduator ownership remains pinned to the factory across re-runs.
        assertEq(Graduator(second.graduator).owner(), second.launchpadFactory, "graduator.owner");
    }

    // -----------------------------------------------------------------------
    // Test — revert paths
    // -----------------------------------------------------------------------

    /// @dev Missing prior-deployment file is the FIRST precondition the
    ///      script checks. Pre-condition: pass valid Uniswap addresses so
    ///      the revert can only come from the Phase-1 check.
    function test_reverts_whenPriorDeploymentMissing() public {
        // setUp's wipe already removed the file but be explicit so the test
        // reads cleanly.
        _wipe(PATH_REVERT_PRIOR);

        vm.expectRevert(abi.encodeWithSelector(DeployPhase2.Phase1DeploymentMissing.selector, PATH_REVERT_PRIOR));
        deployer.runWith(PATH_REVERT_PRIOR, address(uniFactory), address(uniRouter));
    }

    /// @dev Phase-1 file present but with null TITU/Treasury also surfaces
    ///      `Phase1DeploymentMissing`. Catches the half-populated case where
    ///      a previous Phase-1 dry-run wrote a skeleton.
    function test_reverts_whenPhase1AddressesUnpopulated() public {
        // Skeleton file with all-null Phase-1 contracts.
        string memory skeleton = string.concat(
            "{\n",
            '  "chainId": 31337,\n',
            '  "deployedAt": null,\n',
            '  "contracts": {\n',
            '    "TITU": null,\n',
            '    "Treasury": null\n',
            "  }\n",
            "}\n"
        );
        vm.writeFile(PATH_REVERT_NULL, skeleton);

        vm.expectRevert(abi.encodeWithSelector(DeployPhase2.Phase1DeploymentMissing.selector, PATH_REVERT_NULL));
        deployer.runWith(PATH_REVERT_NULL, address(uniFactory), address(uniRouter));
    }

    /// @dev With both Uniswap addresses zero, the script reverts BEFORE any
    ///      contract is deployed. Pre-condition: valid Phase-1 file.
    function test_reverts_whenUniswapAddressesUnset() public {
        _writePhase1(PATH_REVERT_UNI);

        vm.expectRevert(DeployPhase2.UniswapV2NotConfigured.selector);
        deployer.runWith(PATH_REVERT_UNI, address(0), address(0));
    }

    /// @dev Half-set Uniswap (factory only) is the silent-footgun shape; the
    ///      script must still revert.
    function test_reverts_whenUniswapRouterMissing() public {
        _writePhase1(PATH_REVERT_UNI);

        vm.expectRevert(DeployPhase2.UniswapV2NotConfigured.selector);
        deployer.runWith(PATH_REVERT_UNI, address(uniFactory), address(0));
    }

    /// @dev Half-set Uniswap (router only) — symmetric coverage of the
    ///      factory-missing branch.
    function test_reverts_whenUniswapFactoryMissing() public {
        _writePhase1(PATH_REVERT_UNI);

        vm.expectRevert(DeployPhase2.UniswapV2NotConfigured.selector);
        deployer.runWith(PATH_REVERT_UNI, address(0), address(uniRouter));
    }

    /// @dev `runWithPath` (env-driven) reverts with `UniswapV2NotConfigured`
    ///      when neither env var is set. Mirrors the production wiring path
    ///      and keeps the env-driven entry point covered. Locks env-mutation
    ///      to this single test so it does not race other cases.
    function test_runWithPath_revertsWhenUniswapEnvUnset() public {
        _writePhase1(PATH_REVERT_UNI);

        // Force both vars to the empty-address sentinel; `vm.setEnv` cannot
        // delete a var, but the script's parser folds zero into the
        // `UniswapV2NotConfigured` revert.
        vm.setEnv("UNISWAP_V2_FACTORY", "0x0000000000000000000000000000000000000000");
        vm.setEnv("UNISWAP_V2_ROUTER", "0x0000000000000000000000000000000000000000");

        vm.expectRevert(DeployPhase2.UniswapV2NotConfigured.selector);
        deployer.runWithPath(PATH_REVERT_UNI);
    }

    /// @dev Phase-1-missing must out-rank uni-config-missing. The user's
    ///      prior WIP failed on this ordering — pinning it now so we never
    ///      regress.
    function test_revertOrdering_phase1ChecksBeforeUniswap() public {
        // No Phase-1 file AND no Uniswap addresses: the Phase-1 selector wins.
        _wipe(PATH_REVERT_PRIOR);

        vm.expectRevert(abi.encodeWithSelector(DeployPhase2.Phase1DeploymentMissing.selector, PATH_REVERT_PRIOR));
        deployer.runWith(PATH_REVERT_PRIOR, address(0), address(0));
    }

    // -----------------------------------------------------------------------
    // Test — JSON merge
    // -----------------------------------------------------------------------

    /// @dev Reading the merged file after a successful run yields BOTH the
    ///      Phase-1 keys (verbatim from the input fixture) and the freshly
    ///      written Phase-2 keys.
    function test_writesMergedJson_withPhase1AndPhase2() public {
        _writePhase1(PATH_MERGE);

        DeployPhase2.Phase2Deployment memory d = deployer.runWith(PATH_MERGE, address(uniFactory), address(uniRouter));

        string memory merged = vm.readFile(PATH_MERGE);

        // Phase-1 — every populated input key survives the round trip.
        assertEq(vm.parseJsonAddress(merged, ".phase1.contracts.TITU"), FIXTURE_TITU, "merged.titu");
        assertEq(vm.parseJsonAddress(merged, ".phase1.contracts.Treasury"), FIXTURE_TREASURY, "merged.treasury");
        assertEq(
            vm.parseJsonAddress(merged, ".phase1.contracts.TreasuryImpl"), FIXTURE_TREASURY_IMPL, "merged.treasuryImpl"
        );
        assertEq(
            vm.parseJsonAddress(merged, ".phase1.contracts.VestingVault"), FIXTURE_VESTING_VAULT, "merged.vestingVault"
        );
        assertEq(vm.parseJsonAddress(merged, ".phase1.contracts.VeTITU"), FIXTURE_VE_TITU, "merged.veTitu");
        assertEq(
            vm.parseJsonAddress(merged, ".phase1.contracts.FeeDistributor"),
            FIXTURE_FEE_DISTRIBUTOR,
            "merged.feeDistributor"
        );

        // Phase-2 — every deployed contract is surfaced under `phase2.contracts`.
        assertEq(
            vm.parseJsonAddress(merged, ".phase2.contracts.AgentTokenImpl"), d.agentTokenImpl, "merged.agentTokenImpl"
        );
        assertEq(
            vm.parseJsonAddress(merged, ".phase2.contracts.BondingCurveImpl"),
            d.bondingCurveImpl,
            "merged.bondingCurveImpl"
        );
        assertEq(vm.parseJsonAddress(merged, ".phase2.contracts.FeeRouter"), d.feeRouter, "merged.feeRouter");
        assertEq(vm.parseJsonAddress(merged, ".phase2.contracts.Graduator"), d.graduator, "merged.graduator");
        assertEq(
            vm.parseJsonAddress(merged, ".phase2.contracts.LaunchpadFactory"),
            d.launchpadFactory,
            "merged.launchpadFactory"
        );
        assertEq(vm.parseJsonAddress(merged, ".phase2.contracts.AntiSniperModule"), d.antiSniper, "merged.antiSniper");
        assertEq(
            vm.parseJsonAddress(merged, ".phase2.contracts.LaunchRadarModule"), d.launchRadar, "merged.launchRadar"
        );
        assertEq(vm.parseJsonAddress(merged, ".phase2.contracts.SixtyDaysModule"), d.sixtyDays, "merged.sixtyDays");
        assertEq(vm.parseJsonAddress(merged, ".phase2.contracts.AirdropModule"), d.airdrop, "merged.airdrop");
        assertEq(
            vm.parseJsonAddress(merged, ".phase2.contracts.CapitalFormationModule"),
            d.capitalFormation,
            "merged.capitalFormation"
        );
        assertEq(vm.parseJsonAddress(merged, ".phase2.contracts.PreBuyModule"), d.preBuy, "merged.preBuy");
        assertEq(
            vm.parseJsonAddress(merged, ".phase2.contracts.ExistingTokenModule"),
            d.existingToken,
            "merged.existingToken"
        );

        // Both top-level keys present with `deployedAt` populated.
        assertTrue(vm.keyExistsJson(merged, ".phase1.deployedAt"), "phase1.deployedAt key");
        assertTrue(vm.keyExistsJson(merged, ".phase2.deployedAt"), "phase2.deployedAt key");
        assertEq(vm.parseJsonUint(merged, ".phase2.deployedAt"), block.timestamp, "phase2.deployedAt value");
        assertEq(vm.parseJsonUint(merged, ".phase1.deployedAt"), 1_700_000_000, "phase1.deployedAt value");
    }

    /// @dev Legacy Phase-1 files (no `phase1` key, top-level `contracts`) are
    ///      auto-promoted on first Phase-2 run. Confirms we never lose data
    ///      on a deployment-file shape upgrade.
    function test_legacyPhase1Shape_isPromotedToPhase1Key() public {
        _writeLegacyPhase1(PATH_LEGACY);

        deployer.runWith(PATH_LEGACY, address(uniFactory), address(uniRouter));

        string memory merged = vm.readFile(PATH_LEGACY);
        assertEq(vm.parseJsonAddress(merged, ".phase1.contracts.TITU"), FIXTURE_TITU, "legacy.titu");
        assertEq(vm.parseJsonAddress(merged, ".phase1.contracts.Treasury"), FIXTURE_TREASURY, "legacy.treasury");
        // Legacy `deployedAt` carried forward.
        assertEq(vm.parseJsonUint(merged, ".phase1.deployedAt"), 1_700_000_000, "legacy.deployedAt");
        // `phase2` block lands alongside the promoted `phase1` block.
        assertTrue(vm.keyExistsJson(merged, ".phase2.contracts.LaunchpadFactory"), "legacy -> phase2 key");
    }

    // -----------------------------------------------------------------------
    // Helpers
    // -----------------------------------------------------------------------

    /// @dev Write a populated Phase-1 deployment fixture in the merged shape
    ///      so the script's preflight passes.
    function _writePhase1(string memory path) internal {
        string memory json = string.concat(
            "{\n",
            '  "chainId": 31337,\n',
            '  "phase1": {\n',
            '    "deployedAt": 1700000000,\n',
            '    "contracts": {\n',
            '      "TITU": "',
            vm.toString(FIXTURE_TITU),
            '",\n',
            '      "Treasury": "',
            vm.toString(FIXTURE_TREASURY),
            '",\n',
            '      "TreasuryImpl": "',
            vm.toString(FIXTURE_TREASURY_IMPL),
            '",\n',
            '      "VestingVault": "',
            vm.toString(FIXTURE_VESTING_VAULT),
            '",\n',
            '      "VeTITU": "',
            vm.toString(FIXTURE_VE_TITU),
            '",\n',
            '      "FeeDistributor": "',
            vm.toString(FIXTURE_FEE_DISTRIBUTOR),
            '"\n',
            "    }\n",
            "  }\n",
            "}\n"
        );
        vm.writeFile(path, json);
    }

    /// @dev Write a legacy Phase-1 fixture (top-level `contracts`, no `phase1`
    ///      wrapper) — the shape DeployPhase1.s.sol emits today.
    function _writeLegacyPhase1(string memory path) internal {
        string memory json = string.concat(
            "{\n",
            '  "chainId": 31337,\n',
            '  "deployedAt": 1700000000,\n',
            '  "contracts": {\n',
            '    "TITU": "',
            vm.toString(FIXTURE_TITU),
            '",\n',
            '    "Treasury": "',
            vm.toString(FIXTURE_TREASURY),
            '",\n',
            '    "TreasuryImpl": "',
            vm.toString(FIXTURE_TREASURY_IMPL),
            '",\n',
            '    "VestingVault": "',
            vm.toString(FIXTURE_VESTING_VAULT),
            '",\n',
            '    "VeTITU": "',
            vm.toString(FIXTURE_VE_TITU),
            '",\n',
            '    "FeeDistributor": "',
            vm.toString(FIXTURE_FEE_DISTRIBUTOR),
            '"\n',
            "  }\n",
            "}\n"
        );
        vm.writeFile(path, json);
    }
}
