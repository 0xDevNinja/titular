// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {IBondingCurve} from "../../src/launchpad/interfaces/IBondingCurve.sol";

/// @title MockBondingCurve
/// @notice Test double implementing the `IBondingCurve` surface that the
///         graduator binds against. Holds both tokens as balances, tracks a
///         mutable `graduated` flag, and exposes `pullForGraduation` that
///         transfers the configured reserves to `msg.sender` (the graduator).
///         Never enforces CEI or reentrancy — those behaviours are verified
///         by the real `BondingCurve` tests. Here we only need a deterministic
///         source of tokens for add-liquidity assertions.
contract MockBondingCurve is IBondingCurve {
    address public override agentToken;
    address public override quoteToken;
    bool public override graduated;
    address public override graduator;
    uint256 public override realQuoteReserve;
    uint256 public override realAgentReserve;

    uint256 public pullCalls;

    /// @notice Configure the curve fixture. Tokens are transferred into the
    ///         mock up-front; `pullForGraduation` forwards exactly the
    ///         recorded reserves to the graduator.
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
