// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import {TITU} from "../src/token/TITU.sol";
import {Treasury} from "../src/treasury/Treasury.sol";
import {VestingVault} from "../src/vault/VestingVault.sol";
import {VeTITU} from "../src/governance/VeTITU.sol";
import {FeeDistributor} from "../src/fees/FeeDistributor.sol";
import {IVeTITU} from "../src/fees/IVeTITU.sol";

/// @notice Phase-1 deployment script.
/// @dev Deployment order: TITU -> VestingVault -> VeTITU -> Treasury(proxy) -> FeeDistributor
///      -> Treasury.setFeeDistributor. Writes `deployments/base-sepolia.json`.
///      Reads `MULTISIG` and `DEPLOYER_PRIVATE_KEY` from env. Missing PK runs in
///      dry-run mode (no broadcast) so the script remains safe to invoke from CI.
contract DeployPhase1 is Script {
    struct Deployment {
        address titu;
        address treasury;
        address treasuryImpl;
        address vestingVault;
        address ve;
        address feeDistributor;
    }

    function run() external returns (Deployment memory d) {
        uint256 pk = vm.envOr("DEPLOYER_PRIVATE_KEY", uint256(0));
        address multisig = vm.envOr("MULTISIG", address(0));
        if (multisig == address(0)) {
            multisig = pk == 0 ? address(0xDEADBEEF) : vm.addr(pk);
        }

        if (pk != 0) vm.startBroadcast(pk);

        TITU titu = new TITU(multisig);
        VestingVault vv = new VestingVault(IERC20(address(titu)), multisig);
        VeTITU ve = new VeTITU(IERC20(address(titu)));

        Treasury impl = new Treasury();
        bytes memory initData = abi.encodeCall(Treasury.initialize, (multisig, address(0)));
        ERC1967Proxy proxy = new ERC1967Proxy(address(impl), initData);
        Treasury treasury = Treasury(payable(address(proxy)));

        FeeDistributor fd = new FeeDistributor(IVeTITU(address(ve)), address(treasury));

        // Wire FeeDistributor into Treasury. Caller must be `multisig` when broadcasting,
        // otherwise we skip this step — the prod ceremony will set it via Safe later.
        if (pk != 0 && vm.addr(pk) == multisig) {
            treasury.setFeeDistributor(address(fd));
        }

        if (pk != 0) vm.stopBroadcast();

        d = Deployment({
            titu: address(titu),
            treasury: address(treasury),
            treasuryImpl: address(impl),
            vestingVault: address(vv),
            ve: address(ve),
            feeDistributor: address(fd)
        });

        _writeDeployment(d);
        console2.log("TITU", d.titu);
        console2.log("Treasury (proxy)", d.treasury);
        console2.log("Treasury (impl)", d.treasuryImpl);
        console2.log("VestingVault", d.vestingVault);
        console2.log("VeTITU", d.ve);
        console2.log("FeeDistributor", d.feeDistributor);
    }

    function _writeDeployment(Deployment memory d) internal {
        string memory path = "deployments/base-sepolia.json";
        string memory json = string.concat(
            "{\n",
            '  "chainId": 84532,\n',
            '  "deployedAt": ',
            vm.toString(block.timestamp),
            ",\n",
            '  "contracts": {\n',
            '    "TITU": "',
            vm.toString(d.titu),
            '",\n',
            '    "Treasury": "',
            vm.toString(d.treasury),
            '",\n',
            '    "TreasuryImpl": "',
            vm.toString(d.treasuryImpl),
            '",\n',
            '    "VestingVault": "',
            vm.toString(d.vestingVault),
            '",\n',
            '    "VeTITU": "',
            vm.toString(d.ve),
            '",\n',
            '    "FeeDistributor": "',
            vm.toString(d.feeDistributor),
            '"\n',
            "  }\n",
            "}\n"
        );
        vm.writeFile(path, json);
    }
}
