// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";

import {AgentRegistry} from "../src/acp/AgentRegistry.sol";
import {Job} from "../src/acp/Job.sol";
import {JobFactory} from "../src/acp/JobFactory.sol";
import {Escrow} from "../src/acp/Escrow.sol";
import {FeeSplitter} from "../src/acp/FeeSplitter.sol";
import {BuybackBurner} from "../src/acp/BuybackBurner.sol";
import {HookRegistry} from "../src/acp/HookRegistry.sol";
import {FundTransferHook} from "../src/acp/FundTransferHook.sol";
import {SubscriptionHook} from "../src/acp/SubscriptionHook.sol";
import {MilestoneHook} from "../src/acp/MilestoneHook.sol";
import {RoyaltyHook} from "../src/acp/RoyaltyHook.sol";

/// @notice Phase-3 deployment: ACP v2 contracts on Base Sepolia.
/// @dev    Deploy order:
///           1. AgentRegistry
///           2. Job implementation (EIP-1167 master copy)
///           3. JobFactory (clones Job)
///           4. Escrow (with Permit2)
///           5. FeeSplitter
///           6. BuybackBurner
///           7. HookRegistry + 4 hook implementations
///
///         Required env vars:
///           DEPLOYER_ADDRESS    — deployer EOA address (for broadcast)
///           ADMIN_ADDRESS       — admin address granted all admin roles
///           TREASURY_ADDRESS    — treasury recipient for fee splits
///           BUYBACK_BURNER_TOKEN — payment token address for BuybackBurner
///           TITU_ADDRESS        — TITU token address for BuybackBurner
///           UNI_V2_ROUTER       — Uniswap V2 router on Base Sepolia
///           PERMIT2_ADDRESS     — Permit2 canonical address
///           ARBITER_ADDRESS     — default arbiter for job disputes
///
///         All addresses are written to deployments/local.json after deploy.
///
/// @custom:example
///   forge script script/DeployPhase3.s.sol:DeployPhase3 \
///     --rpc-url $BASE_SEPOLIA_RPC --broadcast --verify \
///     --etherscan-api-key $BASESCAN_API_KEY \
///     --sender $DEPLOYER_ADDRESS
contract DeployPhase3 is Script {
    // ─────────────────────────────────────────────────────────────
    // Deployed contract references
    // ─────────────────────────────────────────────────────────────

    AgentRegistry public agentRegistry;
    Job public jobImpl;
    JobFactory public jobFactory;
    Escrow public escrow;
    FeeSplitter public feeSplitter;
    BuybackBurner public buybackBurner;
    HookRegistry public hookRegistry;
    FundTransferHook public fundTransferHook;
    SubscriptionHook public subscriptionHook;
    MilestoneHook public milestoneHook;
    RoyaltyHook public royaltyHook;

    // ─────────────────────────────────────────────────────────────
    // run()
    // ─────────────────────────────────────────────────────────────

    function run() external {
        address admin = _envAddr("ADMIN_ADDRESS");
        address treasury = _envAddr("TREASURY_ADDRESS");
        address permit2 = _envAddr("PERMIT2_ADDRESS");
        address arbiter = _envAddr("ARBITER_ADDRESS");
        address paymentToken = _envAddr("BUYBACK_BURNER_TOKEN");
        address titu = _envAddr("TITU_ADDRESS");
        address uniRouter = _envAddr("UNI_V2_ROUTER");

        vm.startBroadcast();
        _deployCore(admin, arbiter, permit2, treasury);
        _deployBuybackAndFees(admin, paymentToken, titu, uniRouter, treasury);
        _deployFeeSplitter(admin, treasury);
        _deployJobFactory(admin, arbiter);
        _deployHooks(admin);
        vm.stopBroadcast();

        _writeDeployments(admin, permit2, arbiter);
    }

    function _deployCore(address admin, address arbiter, address permit2, address treasury) internal {
        agentRegistry = new AgentRegistry(admin);
        console2.log("AgentRegistry:", address(agentRegistry));

        // Deploy Escrow and FeeSplitter first (needed by JobFactory constructor)
        escrow = new Escrow(admin, permit2);
        console2.log("Escrow:", address(escrow));

        // BuybackBurner is deployed in _deployBuybackAndFees; use address(1) as placeholder
        // and update after that step. FeeSplitter requires buybackBurner address at construction.
        // Deploy order: Escrow → BuybackBurner → FeeSplitter → HookRegistry → Job impl → JobFactory

        jobImpl = new Job();
        console2.log("Job (impl):", address(jobImpl));
    }

    function _deployFeeSplitter(address admin, address treasury) internal {
        feeSplitter = new FeeSplitter(admin, treasury, address(buybackBurner));
        console2.log("FeeSplitter:", address(feeSplitter));

        hookRegistry = new HookRegistry(admin);
        console2.log("HookRegistry (pre-hooks):", address(hookRegistry));
    }

    function _deployJobFactory(address admin, address arbiter) internal {
        jobFactory = new JobFactory(admin, address(jobImpl), agentRegistry, arbiter, escrow, feeSplitter, hookRegistry);
        console2.log("JobFactory:", address(jobFactory));

        // Grant JobFactory DEFAULT_ADMIN_ROLE on Escrow so it can grant RELEASER_ROLE to clones
        escrow.grantRole(escrow.DEFAULT_ADMIN_ROLE(), address(jobFactory));
        // Grant JobFactory DEFAULT_ADMIN_ROLE on FeeSplitter so it can grant CALLER_ROLE to clones
        feeSplitter.grantRole(feeSplitter.DEFAULT_ADMIN_ROLE(), address(jobFactory));
        console2.log("Roles wired: JobFactory -> Escrow + FeeSplitter");
    }

    function _deployBuybackAndFees(address admin, address paymentToken, address titu, address uniRouter, address)
        internal
    {
        address[] memory swapPath = new address[](2);
        swapPath[0] = paymentToken;
        swapPath[1] = titu;
        // minTituOut = 0 initially (admin should set a sensible floor after deployment)
        buybackBurner = new BuybackBurner(admin, uniRouter, paymentToken, titu, swapPath, 0, 300);
        console2.log("BuybackBurner:", address(buybackBurner));
    }

    function _deployHooks(address) internal {
        // HookRegistry was already deployed in _deployFeeSplitter
        fundTransferHook = new FundTransferHook();
        subscriptionHook = new SubscriptionHook();
        milestoneHook = new MilestoneHook();
        royaltyHook = new RoyaltyHook();

        hookRegistry.registerHook(address(fundTransferHook));
        hookRegistry.registerHook(address(subscriptionHook));
        hookRegistry.registerHook(address(milestoneHook));
        hookRegistry.registerHook(address(royaltyHook));

        console2.log("HookRegistry:", address(hookRegistry));
        console2.log("FundTransferHook:", address(fundTransferHook));
        console2.log("SubscriptionHook:", address(subscriptionHook));
        console2.log("MilestoneHook:", address(milestoneHook));
        console2.log("RoyaltyHook:", address(royaltyHook));
    }

    // ─────────────────────────────────────────────────────────────
    // Helpers
    // ─────────────────────────────────────────────────────────────

    function _envAddr(string memory key) internal view returns (address addr) {
        addr = vm.envAddress(key);
        require(addr != address(0), string.concat("DeployPhase3: missing env ", key));
    }

    function _writeDeployments(address admin, address permit2, address arbiter) internal {
        string memory obj = "phase3";
        vm.serializeAddress(obj, "agentRegistry", address(agentRegistry));
        vm.serializeAddress(obj, "jobImpl", address(jobImpl));
        vm.serializeAddress(obj, "jobFactory", address(jobFactory));
        vm.serializeAddress(obj, "escrow", address(escrow));
        vm.serializeAddress(obj, "feeSplitter", address(feeSplitter));
        vm.serializeAddress(obj, "buybackBurner", address(buybackBurner));
        vm.serializeAddress(obj, "hookRegistry", address(hookRegistry));
        vm.serializeAddress(obj, "fundTransferHook", address(fundTransferHook));
        vm.serializeAddress(obj, "subscriptionHook", address(subscriptionHook));
        vm.serializeAddress(obj, "milestoneHook", address(milestoneHook));
        vm.serializeAddress(obj, "royaltyHook", address(royaltyHook));
        vm.serializeAddress(obj, "admin", admin);
        vm.serializeAddress(obj, "permit2", permit2);
        string memory finalJson = vm.serializeAddress(obj, "arbiter", arbiter);
        vm.writeJson(finalJson, "deployments/phase3.json");
        console2.log("Deployment manifest written to deployments/phase3.json");
    }
}
