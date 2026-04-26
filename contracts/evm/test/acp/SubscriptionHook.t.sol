// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {SubscriptionHook} from "../../src/acp/SubscriptionHook.sol";

contract MockToken is ERC20 {
    constructor() ERC20("MockSub", "MSUB") {}
    function mint(address to, uint256 amount) external { _mint(to, amount); }
}

contract SubscriptionHookTest is Test {
    SubscriptionHook internal hook;
    MockToken internal token;

    address internal subscriber = makeAddr("subscriber");
    address internal provider = makeAddr("provider");
    uint256 internal constant PERIOD_AMOUNT = 100e18;
    uint64 internal constant PERIOD_DURATION = 30 days;
    uint256 internal constant JOB_ID = 5;

    event SubscriptionRenewed(
        uint256 indexed jobId, address indexed subscriber, address indexed provider,
        uint256 amount, uint64 nextRenewalAt
    );
    event SubscriptionCancelled(uint256 indexed jobId);

    function _initCtx() internal view returns (bytes memory) {
        return abi.encode(SubscriptionHook.RenewContext({
            subscriber: subscriber, provider: provider, token: address(token),
            periodAmount: PERIOD_AMOUNT, periodDuration: PERIOD_DURATION
        }));
    }

    function setUp() public {
        hook = new SubscriptionHook();
        token = new MockToken();
        token.mint(subscriber, PERIOD_AMOUNT * 100);
        vm.prank(subscriber);
        token.approve(address(hook), type(uint256).max);
    }

    function test_hookName() public view {
        assertEq(hook.hookName(), "SubscriptionHook");
    }

    function test_onAccept_initialisesState() public {
        hook.onAccept(JOB_ID, _initCtx());
        (uint64 nextRenewal, bool active) = hook.subscriptions(JOB_ID);
        assertTrue(active);
        assertGt(nextRenewal, 0);
    }

    function test_onApprove_renewsSubscription() public {
        hook.onAccept(JOB_ID, _initCtx());
        uint256 before = token.balanceOf(provider);
        hook.onApprove(JOB_ID, _initCtx());
        assertEq(token.balanceOf(provider), before + PERIOD_AMOUNT);
    }

    function test_onApprove_updatesNextRenewal() public {
        hook.onAccept(JOB_ID, _initCtx());
        hook.onApprove(JOB_ID, _initCtx());
        (uint64 next,) = hook.subscriptions(JOB_ID);
        assertApproxEqAbs(next, uint64(block.timestamp) + PERIOD_DURATION, 1);
    }

    function test_onApprove_revert_notActive() public {
        // Never initialised — active=false by default
        vm.expectRevert(abi.encodeWithSelector(SubscriptionHook.SubscriptionNotActive.selector, JOB_ID));
        hook.onApprove(JOB_ID, _initCtx());
    }

    function test_onCancel_deactivatesSubscription() public {
        hook.onAccept(JOB_ID, _initCtx());
        vm.expectEmit(true, false, false, false);
        emit SubscriptionCancelled(JOB_ID);
        hook.onCancel(JOB_ID, "");
        (, bool active) = hook.subscriptions(JOB_ID);
        assertFalse(active);
    }

    function test_onCancel_noop_ifNotInitialised() public {
        hook.onCancel(99, ""); // no revert
    }

    function testFuzz_onApprove_multiPeriod(uint8 periods) public {
        periods = uint8(bound(periods, 1, 10));
        hook.onAccept(JOB_ID, _initCtx());
        uint256 beforeBalance = token.balanceOf(provider);
        for (uint8 i = 0; i < periods; ++i) {
            hook.onApprove(JOB_ID, _initCtx());
        }
        assertEq(token.balanceOf(provider), beforeBalance + PERIOD_AMOUNT * periods);
    }
}
