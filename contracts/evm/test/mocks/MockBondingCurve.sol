// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {IBondingCurve} from "../../src/launchpad/interfaces/IBondingCurve.sol";

/// @title MockBondingCurve
/// @notice Test double implementing the full `IBondingCurve` surface used by
///         module unit tests: the trade entry point ({buy}) consumed by the
///         pre-buy module, and the graduation accessors + reserve drain
///         consumed by the graduator.
/// @dev    Two independent fixture surfaces share one mock to avoid drift:
///
///         - Pre-buy path: construct with `(agent_, quote_)`, pre-fund the
///           mock with the agent token, and let modules call {buy}. The
///           `payoutOverride` and `revertOnBuy` knobs let a single test
///           exercise the under-fund / over-deliver / revert paths without
///           swapping the mock implementation.
///
///         - Graduator path: construct with `(address(0), address(0))`, then
///           use {set} to rebind tokens and seed the recorded reserves. The
///           graduator calls {pullForGraduation} to drain into itself.
///
///         Never enforces CEI or reentrancy — those behaviours are verified
///         by the real `BondingCurve` tests. Here we only need a deterministic
///         source of tokens for module-level assertions.
contract MockBondingCurve is IBondingCurve {
    // ---------------------------------------------------------------------
    // IBondingCurve state
    // ---------------------------------------------------------------------

    address public override agentToken;
    address public override quoteToken;
    bool public override graduated;
    address public override graduator;
    uint256 public override realQuoteReserve;
    uint256 public override realAgentReserve;

    // ---------------------------------------------------------------------
    // Pre-buy buy() shim
    // ---------------------------------------------------------------------

    /// @notice If non-zero, used as the literal agent-out delivered to the
    ///         buyer instead of `minAgentOut`. Lets a single test exercise
    ///         the under-fund or over-deliver paths.
    uint256 public payoutOverride;

    /// @notice If true, {buy} reverts with a string the test recognises.
    bool public revertOnBuy;

    // ---------------------------------------------------------------------
    // Graduator instrumentation
    // ---------------------------------------------------------------------

    /// @notice Number of times {pullForGraduation} has been invoked. Used by
    ///         the graduator tests to assert one-shot semantics.
    uint256 public pullCalls;

    // ---------------------------------------------------------------------
    // Construction + configuration
    // ---------------------------------------------------------------------

    /// @notice Bind the trade tokens at deploy time. Pass `(address(0),
    ///         address(0))` and call {set} when the test fixture (e.g. the
    ///         graduator) needs to seed reserves separately.
    constructor(address agent_, address quote_) {
        agentToken = agent_;
        quoteToken = quote_;
    }

    /// @notice Rebind tokens and seed reserves in one call. Used by the
    ///         graduator fixture; pre-buy tests do not need this entry point.
    function set(
        address agentToken_,
        address quoteToken_,
        uint256 quoteReserve_,
        uint256 agentReserve_,
        bool graduated_
    ) external {
        agentToken = agentToken_;
        quoteToken = quoteToken_;
        realQuoteReserve = quoteReserve_;
        realAgentReserve = agentReserve_;
        graduated = graduated_;
    }

    function setGraduator(address g) external {
        graduator = g;
    }

    function setGraduated(bool g) external {
        graduated = g;
    }

    function setPayoutOverride(uint256 v) external {
        payoutOverride = v;
    }

    function setRevertOnBuy(bool v) external {
        revertOnBuy = v;
    }

    // ---------------------------------------------------------------------
    // IBondingCurve trade surface
    // ---------------------------------------------------------------------

    /// @inheritdoc IBondingCurve
    function buy(uint256 minAgentOut, uint256 quoteIn) external override {
        require(!revertOnBuy, "curve: forced revert");
        IERC20(quoteToken).transferFrom(msg.sender, address(this), quoteIn);
        uint256 out = payoutOverride == 0 ? minAgentOut : payoutOverride;
        IERC20(agentToken).transfer(msg.sender, out);
    }

    // ---------------------------------------------------------------------
    // IBondingCurve graduation surface
    // ---------------------------------------------------------------------

    /// @inheritdoc IBondingCurve
    function pullForGraduation() external override returns (uint256 quoteAmount, uint256 agentAmount) {
        pullCalls += 1;
        quoteAmount = realQuoteReserve;
        agentAmount = realAgentReserve;
        realQuoteReserve = 0;
        realAgentReserve = 0;
        if (quoteAmount != 0) IERC20(quoteToken).transfer(msg.sender, quoteAmount);
        if (agentAmount != 0) IERC20(agentToken).transfer(msg.sender, agentAmount);
    }

    /// @inheritdoc IBondingCurve
    function requestGraduation() external override {
        graduated = true;
    }
}
