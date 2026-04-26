// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {ExistingTokenModule} from "../../src/launchpad/modules/ExistingTokenModule.sol";

/// @dev Minimal mintable mock ERC-20 implementing the standard metadata
///      extension (name / symbol / decimals). Stands in for any pre-existing
///      project token that wants to wrap onto the launchpad.
contract MockERC20 is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev Token whose `decimals()` reverts. Used to assert {configure} rejects
///      tokens that don't honour the standard ERC-20 metadata interface.
contract MockNoDecimals is ERC20 {
    constructor() ERC20("NoDec", "ND") {}

    function decimals() public pure override returns (uint8) {
        revert("nope");
    }

    function mint(address to, uint256 amt) external {
        _mint(to, amt);
    }
}

/// @dev Reentrancy attacker: a fake "ERC-20" whose `transferFrom` re-enters
///      {ExistingTokenModule.wrap}. Verifies the {ReentrancyGuard} on the
///      top-up path holds.
contract ReentrantToken {
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
        // First call from configure(): allow it through cleanly so the module
        // is set up. Subsequent calls from the test's wrap() probe trigger the
        // re-entry attempt.
        if (attackPrimed) {
            attackPrimed = false;
            // Re-enter wrap(); the guard should catch this.
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

/// @dev Unit coverage for {ExistingTokenModule}. Exercises the owner-gated
///      one-shot {configure} path, the open {wrap} top-up path with reentrancy
///      probe, and the {getConfig} view. Custom error selectors only.
contract ExistingTokenModuleTest is Test {
    ExistingTokenModule internal module;

    address internal owner = address(0xA0);
    address internal attacker = address(0xBAD);
    address internal admin = address(0xAD);
    address internal curve = address(0xC0DE);
    address internal curve2 = address(0xC0DF);
    address internal depositor = address(0xDE);

    MockERC20 internal token;
    MockERC20 internal token2;

    uint256 internal constant SUPPLY = 1_000_000_000e18;
    uint256 internal constant TOPUP = 5000e18;

    // Re-declared events so {vm.expectEmit} can match by signature.
    event Configured(address indexed externalToken, address indexed curve, uint256 supply, address indexed agentAdmin);
    event Wrapped(address indexed externalToken, uint256 amount, address indexed depositor);

    function setUp() public {
        module = new ExistingTokenModule(owner);
        token = new MockERC20("Existing", "EXT");
        token2 = new MockERC20("Other", "OTH");
    }

    // ---------------------------------------------------------------------
    // Constructor
    // ---------------------------------------------------------------------

    function test_constructor_sets_owner() public view {
        assertEq(module.owner(), owner);
    }

    function test_constructor_reverts_on_zero_owner() public {
        vm.expectRevert(ExistingTokenModule.ZeroAddress.selector);
        new ExistingTokenModule(address(0));
    }

    // ---------------------------------------------------------------------
    // configure — access control
    // ---------------------------------------------------------------------

    function test_configure_only_owner() public {
        token.mint(address(module), SUPPLY);
        vm.prank(attacker);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, attacker));
        module.configure(address(token), curve, SUPPLY, admin);
    }

    // ---------------------------------------------------------------------
    // configure — happy path
    // ---------------------------------------------------------------------

    function test_configure_happy_writes_storage_and_seeds_curve() public {
        token.mint(address(module), SUPPLY);

        vm.expectEmit(true, true, true, true, address(module));
        emit Configured(address(token), curve, SUPPLY, admin);

        vm.prank(owner);
        module.configure(address(token), curve, SUPPLY, admin);

        // Storage matches.
        (address gotCurve, address gotAdmin, uint256 gotSupply, bool gotWrapped) = module.configs(address(token));
        assertEq(gotCurve, curve);
        assertEq(gotAdmin, admin);
        assertEq(gotSupply, SUPPLY);
        assertTrue(gotWrapped);

        // Tokens flushed to the curve.
        assertEq(token.balanceOf(curve), SUPPLY);
        assertEq(token.balanceOf(address(module)), 0);
    }

    function test_configure_excess_balance_stays_on_module() public {
        // Pre-deposit MORE than `supply`. Module forwards exactly `supply`; the
        // surplus sits idle (and is available for a future wrap top-up).
        uint256 excess = 1234e18;
        token.mint(address(module), SUPPLY + excess);

        vm.prank(owner);
        module.configure(address(token), curve, SUPPLY, admin);

        assertEq(token.balanceOf(curve), SUPPLY);
        assertEq(token.balanceOf(address(module)), excess);
    }

    // ---------------------------------------------------------------------
    // configure — validation reverts
    // ---------------------------------------------------------------------

    function test_configure_reverts_zero_external_token() public {
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.ZeroAddress.selector);
        module.configure(address(0), curve, SUPPLY, admin);
    }

    function test_configure_reverts_zero_curve() public {
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.ZeroAddress.selector);
        module.configure(address(token), address(0), SUPPLY, admin);
    }

    function test_configure_reverts_zero_admin() public {
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.ZeroAddress.selector);
        module.configure(address(token), curve, SUPPLY, address(0));
    }

    function test_configure_reverts_zero_supply() public {
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.ZeroSupply.selector);
        module.configure(address(token), curve, 0, admin);
    }

    function test_configure_reverts_double_configure() public {
        token.mint(address(module), SUPPLY);
        vm.prank(owner);
        module.configure(address(token), curve, SUPPLY, admin);

        // Even with another fresh pre-deposit, second configure on the same
        // token reverts.
        token.mint(address(module), SUPPLY);
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.AlreadyConfigured.selector);
        module.configure(address(token), curve2, SUPPLY, admin);
    }

    function test_configure_reverts_when_module_underfunded() public {
        // Pre-deposit less than `supply` — must revert with the held vs required
        // figures so indexers can debug.
        uint256 held = SUPPLY / 2;
        token.mint(address(module), held);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(ExistingTokenModule.InsufficientPreDeposit.selector, held, SUPPLY));
        module.configure(address(token), curve, SUPPLY, admin);
    }

    function test_configure_reverts_when_no_predeposit() public {
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(ExistingTokenModule.InsufficientPreDeposit.selector, 0, SUPPLY));
        module.configure(address(token), curve, SUPPLY, admin);
    }

    function test_configure_reverts_on_non_metadata_token() public {
        MockNoDecimals bad = new MockNoDecimals();
        bad.mint(address(module), SUPPLY);

        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.InvalidExternalToken.selector);
        module.configure(address(bad), curve, SUPPLY, admin);
    }

    function test_configure_reverts_on_eoa_external_token() public {
        // Address with no code at all — staticcall to decimals() returns
        // (true, "") which fails our `data.length < 32` guard.
        address noCode = address(0xDEAD);
        vm.prank(owner);
        vm.expectRevert(ExistingTokenModule.InvalidExternalToken.selector);
        module.configure(noCode, curve, SUPPLY, admin);
    }

    function test_configure_independent_for_distinct_tokens() public {
        token.mint(address(module), SUPPLY);
        token2.mint(address(module), SUPPLY);

        vm.startPrank(owner);
        module.configure(address(token), curve, SUPPLY, admin);
        module.configure(address(token2), curve2, SUPPLY, admin);
        vm.stopPrank();

        (address gotCurve1,,, bool wrapped1) = module.configs(address(token));
        (address gotCurve2,,, bool wrapped2) = module.configs(address(token2));
        assertEq(gotCurve1, curve);
        assertEq(gotCurve2, curve2);
        assertTrue(wrapped1);
        assertTrue(wrapped2);
        assertEq(token.balanceOf(curve), SUPPLY);
        assertEq(token2.balanceOf(curve2), SUPPLY);
    }

    // ---------------------------------------------------------------------
    // wrap — happy path
    // ---------------------------------------------------------------------

    function test_wrap_happy_forwards_to_curve_and_emits() public {
        _configureDefault();

        // Fund the depositor and approve the module to pull.
        token.mint(depositor, TOPUP);
        vm.prank(depositor);
        token.approve(address(module), TOPUP);

        uint256 curveBefore = token.balanceOf(curve);

        vm.expectEmit(true, true, false, true, address(module));
        emit Wrapped(address(token), TOPUP, depositor);

        vm.prank(depositor);
        module.wrap(address(token), TOPUP);

        // Tokens passed straight through; module retains nothing.
        assertEq(token.balanceOf(curve), curveBefore + TOPUP);
        assertEq(token.balanceOf(address(module)), 0);
        assertEq(token.balanceOf(depositor), 0);
    }

    function test_wrap_callable_by_anyone() public {
        _configureDefault();

        // Two distinct depositors top the same curve up — each succeeds.
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

    // ---------------------------------------------------------------------
    // wrap — validation reverts
    // ---------------------------------------------------------------------

    function test_wrap_reverts_zero_amount() public {
        _configureDefault();
        vm.prank(depositor);
        vm.expectRevert(ExistingTokenModule.ZeroSupply.selector);
        module.wrap(address(token), 0);
    }

    function test_wrap_reverts_when_not_configured() public {
        // No configure for `token`. Even with full balance / allowance, wrap reverts.
        token.mint(depositor, TOPUP);
        vm.prank(depositor);
        token.approve(address(module), TOPUP);

        vm.prank(depositor);
        vm.expectRevert(ExistingTokenModule.NotConfigured.selector);
        module.wrap(address(token), TOPUP);
    }

    // ---------------------------------------------------------------------
    // wrap — reentrancy probe
    // ---------------------------------------------------------------------

    function test_wrap_reverts_on_reentrancy() public {
        // Build a fake ERC-20 whose transferFrom re-enters wrap(). Configure it
        // through the module first (with attack disarmed), then arm the attack
        // and attempt a wrap from the depositor.
        ReentrantToken rt = new ReentrantToken();
        rt.setTarget(module);
        rt.mint(address(module), SUPPLY);

        vm.prank(owner);
        module.configure(address(rt), curve, SUPPLY, admin);

        // Fund + approve depositor.
        rt.mint(depositor, TOPUP);
        vm.prank(depositor);
        rt.approve(address(module), TOPUP);

        rt.primeAttack();

        // Outer wrap call reverts because the re-entrant inner wrap trips the
        // guard. forge-std bubbles the inner ReentrancyGuardReentrantCall up
        // through the SafeERC20 call frame.
        vm.prank(depositor);
        vm.expectRevert();
        module.wrap(address(rt), TOPUP);
    }

    // ---------------------------------------------------------------------
    // getConfig
    // ---------------------------------------------------------------------

    function test_getConfig_returns_zero_struct_pre_configure() public view {
        ExistingTokenModule.Config memory cfg = module.getConfig(address(token));
        assertEq(cfg.curve, address(0));
        assertEq(cfg.agentAdmin, address(0));
        assertEq(cfg.supply, 0);
        assertFalse(cfg.wrapped);
    }

    function test_getConfig_returns_full_struct_post_configure() public {
        _configureDefault();
        ExistingTokenModule.Config memory cfg = module.getConfig(address(token));
        assertEq(cfg.curve, curve);
        assertEq(cfg.agentAdmin, admin);
        assertEq(cfg.supply, SUPPLY);
        assertTrue(cfg.wrapped);
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    function _configureDefault() internal {
        token.mint(address(module), SUPPLY);
        vm.prank(owner);
        module.configure(address(token), curve, SUPPLY, admin);
    }
}
