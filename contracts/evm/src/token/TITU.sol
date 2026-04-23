// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import {ERC20Votes} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Votes.sol";
import {Nonces} from "@openzeppelin/contracts/utils/Nonces.sol";

/// @title TITU
/// @notice Titular native governance token. Fixed 1B supply, no inflation.
/// @dev ERC20 + Permit (EIP-2612) + Votes (ERC-5805, timestamp clock) + Burnable.
///      Entire supply is minted to `initialMinter` at construction and no further
///      mint function is exposed. Burns are permitted via ERC20Burnable but do not
///      affect the invariant contract's supply test which only exercises transfers.
contract TITU is ERC20, ERC20Permit, ERC20Votes, ERC20Burnable {
    /// @notice Fixed initial (and maximum) supply: 1,000,000,000 TITU with 18 decimals.
    uint256 public constant INITIAL_SUPPLY = 1_000_000_000 * 1e18;

    /// @dev Thrown when the zero address is passed to the constructor.
    error ZeroAddress();

    /// @param initialMinter Recipient of the full 1B supply. Must be non-zero and is
    ///                      expected to be a Safe multisig in production.
    constructor(address initialMinter) ERC20("Titular", "TITU") ERC20Permit("Titular") {
        if (initialMinter == address(0)) revert ZeroAddress();
        _mint(initialMinter, INITIAL_SUPPLY);
    }

    // ---------------------------------------------------------------------
    // ERC-6372 clock (timestamp mode)
    // ---------------------------------------------------------------------

    /// @notice Clock used for checkpoints. Timestamp mode per ERC-6372.
    /// @return Current block timestamp cast to uint48.
    function clock() public view override returns (uint48) {
        return uint48(block.timestamp);
    }

    /// @notice Machine-readable clock mode description.
    /// @return Always `"mode=timestamp"`.
    // solhint-disable-next-line func-name-mixedcase
    function CLOCK_MODE() public pure override returns (string memory) {
        return "mode=timestamp";
    }

    // ---------------------------------------------------------------------
    // Required overrides (OZ v5 ERC20Votes + ERC20Permit composition)
    // ---------------------------------------------------------------------

    function _update(address from, address to, uint256 value) internal override(ERC20, ERC20Votes) {
        super._update(from, to, value);
    }

    function nonces(address owner) public view override(ERC20Permit, Nonces) returns (uint256) {
        return super.nonces(owner);
    }
}
