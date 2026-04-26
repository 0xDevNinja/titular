// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

import {ExistingTokenModule} from "../../src/launchpad/modules/ExistingTokenModule.sol";

/// @dev Standard mintable ERC-20 with name/symbol/decimals.
contract ExistingTokenERC20 is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev ERC-20 whose `decimals()` reverts. Triggers {InvalidExternalToken}.
contract ExistingTokenNoDecimals is ERC20 {
    constructor() ERC20("X", "X") {}

    function decimals() public pure override returns (uint8) {
        revert("nope");
    }

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev ERC-20 whose `name()` reverts.
contract ExistingTokenNoName is ERC20 {
    constructor() ERC20("X", "X") {}

    function name() public pure override returns (string memory) {
        revert("nope");
    }

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev ERC-20 whose `symbol()` reverts.
contract ExistingTokenNoSymbol is ERC20 {
    constructor() ERC20("X", "X") {}

    function symbol() public pure override returns (string memory) {
        revert("nope");
    }

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev Hand-rolled ERC-20 that re-enters {ExistingTokenModule.wrap} on
///      `transferFrom`. Used to verify the module's reentrancy guard.
contract ExistingTokenReentrantToken {
    string public constant name = "Reentrant";
    string public constant symbol = "RNT";
    uint8 public constant decimals = 18;

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;
    uint256 public totalSupply;

    ExistingTokenModule public target;
    bool public attackPrimed;

    function mint(address to, uint256 amt) external {
        balanceOf[to] += amt;
        totalSupply += amt;
    }

    function approve(address spender, uint256 amt) external returns (bool) {
        allowance[msg.sender][spender] = amt;
        return true;
    }

    function setTarget(ExistingTokenModule m) external {
        target = m;
    }

    function primeAttack() external {
        attackPrimed = true;
    }

    function transfer(address to, uint256 amt) external returns (bool) {
        balanceOf[msg.sender] -= amt;
        balanceOf[to] += amt;
        return true;
    }

    function transferFrom(address from, address to, uint256 amt) external returns (bool) {
        if (attackPrimed) {
            attackPrimed = false;
            target.wrap(address(this), 1);
        }
        if (allowance[from][msg.sender] != type(uint256).max) {
            allowance[from][msg.sender] -= amt;
        }
        balanceOf[from] -= amt;
        balanceOf[to] += amt;
        return true;
    }
}

/// @title  ExistingTokenModuleSuiteTest
/// @notice Per-contract unit suite for {ExistingTokenModule}: owner-gated
///         one-shot configure path, metadata introspection edges, and the
///         permissionless `wrap` top-up with reentrancy guard.
contract ExistingTokenModuleSuiteTest is Test {
    ExistingTokenModule internal module;
    ExistingTokenERC20 internal token;
    ExistingTokenERC20 internal token2;

    address internal owner = address(0xA0);
    address internal attacker = address(0xBAD);
    address internal admin = address(0xAD);
    address internal curve = address(0xC0DE);
    address internal curve2 = address(0xC0DF);
    address internal depositor = address(0xDE);

    uint256 internal constant SUPPLY = 1_000_000_000e18;
    uint256 internal constant TOPUP = 5000e18;

    function setUp() public {
        module = new ExistingTokenModule(owner);
        token = new ExistingTokenERC20("Existing", "EXT");
        token2 = new ExistingTokenERC20("Other", "OTH");
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_zero_owner_reverts() public {
        vm.expectRevert(ExistingTokenModule.ZeroAddress.selector);
        new ExistingTokenModule(address(0));
    }

    function test_constructor_owner_set() public view {
        assertEq(module.owner(), owner);
    }

    // ---------------------------------------------------------------------
    // configure
    // ---------------------------------------------------------------------

    function test_configure_only_owner() public {
        token.mint(address(module), SUPPLY);
        vm.prank(attacker);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, attacker));
        module.configure(address(token), curve, SUPPLY, admin);
    }

    function test_configure_seeds_curve() public {
        token.mint(address(module), SUPPLY);

        vm.expectEmit(true, true, true, true, address(module));
        emit ExistingTokenModule.Configured(address(token), curve, SUPPLY, admin);

        vm.prank(owner);
        module.configure(address(token), curve, SUPPLY, admin);

        (address gotCurve, address gotAdmin, uint256 gotSupply, bool wrapped) = module.configs(address(token));
        assertEq(gotCurve, curve);
        assertEq(gotAdmin, admin);
        assertEq(gotSupply, SUPPLY);
        assertTrue(wrapped);
        assertEq(token.balanceOf(curve), SUPPLY);
        assertEq(token.balanceOf(address(module)), 0);
    }

    function test_configure_excess_balance_remains() public {
        uint256 excess = 1234e18;
        token.mint(address(module), SUPPLY + excess);
        vm.prank(owner);
        module.configure(address(token), curve, SUPPLY, admin);
        assertEq(token.balanceOf(curve), SUPPLY);
        assertEq(token.balanceOf(address(module)), excess);
    }

    function test_configure_zero_external_token_reverts() public {
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.ZeroAddress.selector);
        module.configure(address(0), curve, SUPPLY, admin);
    }

    function test_configure_zero_curve_reverts() public {
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.ZeroAddress.selector);
        module.configure(address(token), address(0), SUPPLY, admin);
    }

    function test_configure_zero_admin_reverts() public {
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.ZeroAddress.selector);
        module.configure(address(token), curve, SUPPLY, address(0));
    }

    function test_configure_zero_supply_reverts() public {
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.ZeroSupply.selector);
        module.configure(address(token), curve, 0, admin);
    }

    function test_configure_one_shot() public {
        token.mint(address(module), SUPPLY);
        vm.prank(owner);
        module.configure(address(token), curve, SUPPLY, admin);

        token.mint(address(module), SUPPLY);
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.AlreadyConfigured.selector);
        module.configure(address(token), curve2, SUPPLY, admin);
    }

    function test_configure_underfunded_reverts() public {
        uint256 held = SUPPLY / 2;
        token.mint(address(module), held);
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(ExistingTokenModule.InsufficientPreDeposit.selector, held, SUPPLY));
        module.configure(address(token), curve, SUPPLY, admin);
    }

    function test_configure_no_predeposit_reverts() public {
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(ExistingTokenModule.InsufficientPreDeposit.selector, 0, SUPPLY));
        module.configure(address(token), curve, SUPPLY, admin);
    }

    function test_configure_reverting_decimals_rejected() public {
        ExistingTokenNoDecimals bad = new ExistingTokenNoDecimals();
        bad.mint(address(module), SUPPLY);
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.InvalidExternalToken.selector);
        module.configure(address(bad), curve, SUPPLY, admin);
    }

    function test_configure_reverting_name_rejected() public {
        ExistingTokenNoName bad = new ExistingTokenNoName();
        bad.mint(address(module), SUPPLY);
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.InvalidExternalToken.selector);
        module.configure(address(bad), curve, SUPPLY, admin);
    }

    function test_configure_reverting_symbol_rejected() public {
        ExistingTokenNoSymbol bad = new ExistingTokenNoSymbol();
        bad.mint(address(module), SUPPLY);
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.InvalidExternalToken.selector);
        module.configure(address(bad), curve, SUPPLY, admin);
    }

    function test_configure_eoa_token_rejected() public {
        // No code at the address — staticcall returns (true, "") which fails
        // the data-length guard.
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.InvalidExternalToken.selector);
        module.configure(address(0xDEAD), curve, SUPPLY, admin);
    }

    function test_configure_independent_per_token() public {
        token.mint(address(module), SUPPLY);
        token2.mint(address(module), SUPPLY);
        vm.startPrank(owner);
        module.configure(address(token), curve, SUPPLY, admin);
        module.configure(address(token2), curve2, SUPPLY, admin);
        vm.stopPrank();

        (address c1,,, bool w1) = module.configs(address(token));
        (address c2,,, bool w2) = module.configs(address(token2));
        assertEq(c1, curve);
        assertEq(c2, curve2);
        assertTrue(w1);
        assertTrue(w2);
    }

    // ---------------------------------------------------------------------
    // wrap
    // ---------------------------------------------------------------------

    function test_wrap_forwards_to_curve() public {
        _configureDefault();
        token.mint(depositor, TOPUP);
        vm.prank(depositor);
        token.approve(address(module), TOPUP);

        vm.expectEmit(true, true, false, true, address(module));
        emit ExistingTokenModule.Wrapped(address(token), TOPUP, depositor);
        vm.prank(depositor);
        module.wrap(address(token), TOPUP);

        assertEq(token.balanceOf(curve), SUPPLY + TOPUP);
        assertEq(token.balanceOf(address(module)), 0);
    }

    function test_wrap_callable_by_anyone() public {
        _configureDefault();
        address d2 = address(0xDE2);
        token.mint(depositor, TOPUP);
        token.mint(d2, TOPUP);
        vm.prank(depositor);
        token.approve(address(module), TOPUP);
        vm.prank(d2);
        token.approve(address(module), TOPUP);

        vm.prank(depositor);
        module.wrap(address(token), TOPUP);
        vm.prank(d2);
        module.wrap(address(token), TOPUP);

        assertEq(token.balanceOf(curve), SUPPLY + 2 * TOPUP);
    }

    function test_wrap_zero_amount_reverts() public {
        _configureDefault();
        vm.prank(depositor);
        vm.expectRevert(ExistingTokenModule.ZeroSupply.selector);
        module.wrap(address(token), 0);
    }

    function test_wrap_unconfigured_reverts() public {
        token.mint(depositor, TOPUP);
        vm.prank(depositor);
        token.approve(address(module), TOPUP);
        vm.prank(depositor);
        vm.expectRevert(ExistingTokenModule.NotConfigured.selector);
        module.wrap(address(token), TOPUP);
    }

    function test_wrap_reentrancy_guarded() public {
        ExistingTokenReentrantToken rt = new ExistingTokenReentrantToken();
        rt.setTarget(module);
        rt.mint(address(module), SUPPLY);
        vm.prank(owner);
        module.configure(address(rt), curve, SUPPLY, admin);

        rt.mint(depositor, TOPUP);
        vm.prank(depositor);
        rt.approve(address(module), TOPUP);
        rt.primeAttack();

        vm.prank(depositor);
        vm.expectRevert(ReentrancyGuard.ReentrancyGuardReentrantCall.selector);
        module.wrap(address(rt), TOPUP);
    }

    // ---------------------------------------------------------------------
    // getConfig
    // ---------------------------------------------------------------------

    function test_getConfig_zero_struct_pre_configure() public view {
        ExistingTokenModule.Config memory cfg = module.getConfig(address(token));
        assertEq(cfg.curve, address(0));
        assertEq(cfg.agentAdmin, address(0));
        assertEq(cfg.supply, 0);
        assertFalse(cfg.wrapped);
    }

    function test_getConfig_post_configure() public {
        _configureDefault();
        ExistingTokenModule.Config memory cfg = module.getConfig(address(token));
        assertEq(cfg.curve, curve);
        assertEq(cfg.agentAdmin, admin);
        assertEq(cfg.supply, SUPPLY);
        assertTrue(cfg.wrapped);
    }

    function _configureDefault() internal {
        token.mint(address(module), SUPPLY);
        vm.prank(owner);
        module.configure(address(token), curve, SUPPLY, admin);
    }
}
