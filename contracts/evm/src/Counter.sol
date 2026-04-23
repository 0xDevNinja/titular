// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title Counter
/// @notice Minimal counter contract used for workspace verification.
contract Counter {
    uint256 public number;

    function setNumber(uint256 newNumber) public {
        number = newNumber;
    }

    function increment() public {
        number++;
    }
}
