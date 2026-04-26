// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IBondingCurve} from "./IBondingCurve.sol";

/// @title  IPreBuyModule
/// @notice Minimal binding the launchpad factory uses to bootstrap a creator
///         pre-buy + vesting clone via the shared PreBuy module.
/// @dev    Mirrors {PreBuyModule.configure}. Keep in lockstep with
///         `src/launchpad/modules/PreBuyModule.sol`. `titanIn` of the curve's
///         quote token MUST be pre-funded onto the module before this call.
interface IPreBuyModule {
    /// @notice Configure the pre-buy + vest for `agent`.
    /// @param  agent            Agent ERC-20 address.
    /// @param  creator          Vest beneficiary.
    /// @param  vestAmount       Agent tokens locked into the vesting clone.
    /// @param  cliffSeconds     Cliff before any vesting unlocks.
    /// @param  durationSeconds  Total vest duration (must be > 0, >= cliff).
    /// @param  titanIn          Quote (TITU) the curve will pull.
    /// @param  curve            Bonding curve for `agent`.
    /// @return clone            Newly-deployed vesting clone.
    function configure(
        address agent,
        address creator,
        uint256 vestAmount,
        uint64 cliffSeconds,
        uint64 durationSeconds,
        uint256 titanIn,
        IBondingCurve curve
    ) external returns (address clone);
}
