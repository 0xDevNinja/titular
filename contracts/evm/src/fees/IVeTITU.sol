// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @notice Minimal interface FeeDistributor expects from a ve contract.
interface IVeTITU {
    function getPastVotes(address user, uint256 timepoint) external view returns (uint256);
    function getPastTotalSupply(uint256 timepoint) external view returns (uint256);
}
