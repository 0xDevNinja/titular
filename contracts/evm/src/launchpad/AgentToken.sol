// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {ERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/// @title AgentToken
/// @notice Per-agent ERC-20 deployed as an EIP-1167 minimal proxy clone by the launchpad factory.
///         Fixed 1B supply minted to the bonding curve on initialization. Every peer-to-peer
///         transfer is taxed 1% (100 bps) to a shared `feeRouter` that splits between the
///         agent creator and the protocol treasury.
/// @dev    Tax is bypassed on mint, burn, and on any transfer that touches `feeRouter` itself
///         so the router can sweep collected fees without paying tax on itself (which would
///         recurse and round down to zero, wasting gas). The implementation contract disables
///         initializers in its constructor; each clone calls {initialize} exactly once.
contract AgentToken is Initializable, ERC20Upgradeable, OwnableUpgradeable {
    /// @notice Fixed per-agent supply: 1,000,000,000 tokens with 18 decimals.
    uint256 public constant TOTAL_SUPPLY = 1_000_000_000e18;

    /// @notice Transfer tax, expressed in basis points (1% = 100 bps).
    uint16 public constant TRADE_TAX_BPS = 100;

    /// @notice Denominator for basis-point math.
    uint16 public constant BPS_DENOMINATOR = 10_000;

    /// @notice Agent creator address (for off-chain attribution; not used in transfer logic).
    address public creator;

    /// @notice Shared fee router that receives the 1% transfer tax.
    address public feeRouter;

    /// @notice Address holding the full initial supply (typically the paired bonding curve).
    address public bondingCurve;

    /// @dev Emitted once per clone, at the end of {initialize}.
    event AgentTokenInitialized(
        address indexed creator, address indexed feeRouter, address indexed bondingCurve, string name, string symbol
    );

    /// @dev Emitted on every taxed transfer. Complements the standard ERC20 `Transfer` events.
    event TaxCollected(address indexed from, address indexed to, uint256 taxAmount);

    /// @dev Thrown when a required address argument is the zero address.
    error ZeroAddress();
    /// @dev Thrown when the token name or symbol is empty.
    error EmptyMetadata();

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initialize a freshly cloned `AgentToken`.
    /// @dev Mints the full {TOTAL_SUPPLY} to `bondingCurve_`. Can only be invoked once per
    ///      proxy. `creator_` is stored for attribution and set as the Ownable owner so that
    ///      downstream module calls (e.g. metadata updates, if added later) can be gated.
    /// @param name_         ERC-20 name; must be non-empty.
    /// @param symbol_       ERC-20 symbol; must be non-empty.
    /// @param creator_      Agent creator; becomes the contract owner. Must be non-zero.
    /// @param feeRouter_    Shared fee router that receives the 1% transfer tax. Non-zero.
    /// @param bondingCurve_ Initial supply recipient (the paired bonding curve). Non-zero.
    function initialize(
        string memory name_,
        string memory symbol_,
        address creator_,
        address feeRouter_,
        address bondingCurve_
    ) external initializer {
        if (bytes(name_).length == 0 || bytes(symbol_).length == 0) revert EmptyMetadata();
        if (creator_ == address(0) || feeRouter_ == address(0) || bondingCurve_ == address(0)) revert ZeroAddress();

        __ERC20_init(name_, symbol_);
        __Ownable_init(creator_);

        creator = creator_;
        feeRouter = feeRouter_;
        bondingCurve = bondingCurve_;

        _mint(bondingCurve_, TOTAL_SUPPLY);

        emit AgentTokenInitialized(creator_, feeRouter_, bondingCurve_, name_, symbol_);
    }

    /// @dev Routes 1% of every peer-to-peer transfer to {feeRouter}. Tax is skipped for:
    ///      - mint (`from == address(0)`) so {TOTAL_SUPPLY} lands intact on the bonding curve,
    ///      - burn (`to == address(0)`) to preserve supply semantics,
    ///      - transfers where either counterparty is `feeRouter`, so the router can sweep
    ///        its own balance without paying tax on tax (which would round to zero and waste
    ///        gas on every sweep).
    ///
    ///      Math: `tax = value * 100 / 10_000` is safe under Solidity 0.8 checked arithmetic
    ///      for any `value <= TOTAL_SUPPLY`. The remainder `value - tax` can never underflow
    ///      because `tax <= value`.
    function _update(address from, address to, uint256 value) internal override {
        address router = feeRouter;

        // Skip tax on mint, burn, and router-involved transfers.
        if (from == address(0) || to == address(0) || from == router || to == router) {
            super._update(from, to, value);
            return;
        }

        uint256 tax = (value * TRADE_TAX_BPS) / BPS_DENOMINATOR;
        if (tax == 0) {
            // Value too small for tax to register; forward full amount (dust-safe).
            super._update(from, to, value);
            return;
        }

        super._update(from, router, tax);
        super._update(from, to, value - tax);

        emit TaxCollected(from, to, tax);
    }
}
