// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {VestingVault} from "../../src/vault/VestingVault.sol";
import {VestingHandler} from "../handlers/VestingHandler.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("M", "M") {
        _mint(msg.sender, 1_000_000_000e18);
    }
}

contract VestingInvariant is StdInvariant, Test {
    VestingVault internal vault;
    MockERC20 internal token;
    VestingHandler internal handler;
    address internal admin = address(this);

    function setUp() public {
        token = new MockERC20();
        vault = new VestingVault(IERC20(address(token)), admin);
        handler = new VestingHandler(vault, IERC20(address(token)), admin);
        token.transfer(address(handler), 500_000_000e18);
        // re-authorize handler as admin on vault for it to add grants; we kept admin=this
        // and instead grant admin role on the handler so it can call admin-gated funcs.
        vault.grantRole(vault.DEFAULT_ADMIN_ROLE(), address(handler));
        targetContract(address(handler));
    }

    /// @notice Vault must always hold at least the unreleased portion of every grant.
    function invariant_vestingSumLeDeposit() public view {
        // We only test the aggregate: deposited - released == outstanding and that's <= balance.
        uint256 bal = token.balanceOf(address(vault));
        uint256 outstanding = handler.totalDeposited() - handler.totalReleased();
        assertGe(bal, outstanding);
    }
}
