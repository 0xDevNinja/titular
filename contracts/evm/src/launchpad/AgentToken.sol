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
/// @dev    Tax is bypassed on:
///           - mint (`from == address(0)`) and burn (`to == address(0)`),
///           - transfers where either counterparty is in the {taxExempt} allowlist —
///             populated at {initialize} with `feeRouter`, `bondingCurve`, and `graduator`.
///         The graduator entry is critical: graduation pulls the curve's full agent reserve
///         into the graduator and immediately forwards it into Uniswap V2 `addLiquidity`. If
///         the 1% tax fired on either leg, the router's `transferFrom` would underflow the
///         graduator's balance and graduation would revert. Exemption is set ONCE at
///         {initialize} — there is no setter and no upgrade path. Curve, graduator, and
///         feeRouter are all protocol-internal contracts; their exemption does not change
///         end-user economics, since user-to-user transfers are still taxed in full.
///
///         The implementation contract disables initializers in its constructor; each
///         clone calls {initialize} exactly once.
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

    /// @notice Allowlist of addresses for which the 1% transfer tax is skipped on either
    ///         leg of a transfer. Populated only at {initialize}; no runtime setter.
    mapping(address account => bool exempt) public taxExempt;

    /// @dev Emitted once per clone, at the end of {initialize}.
    event AgentTokenInitialized(
        address indexed creator, address indexed feeRouter, address indexed bondingCurve, string name, string symbol
    );

    /// @dev Emitted once per exempt address at {initialize} so off-chain indexers can mirror
    ///      the allowlist without having to re-derive it from the launchpad wiring.
    event TaxExemptSet(address indexed addr, bool exempt);

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
    ///      Populates {taxExempt} with `feeRouter_`, `bondingCurve_`, and `graduator_` so
    ///      protocol-internal hops (curve <-> graduator <-> Uniswap router) skip the 1%
    ///      transfer tax. User-to-user transfers remain taxed.
    /// @param name_         ERC-20 name; must be non-empty.
    /// @param symbol_       ERC-20 symbol; must be non-empty.
    /// @param creator_      Agent creator; becomes the contract owner. Must be non-zero.
    /// @param feeRouter_    Shared fee router that receives the 1% transfer tax. Non-zero.
    /// @param bondingCurve_ Initial supply recipient (the paired bonding curve). Non-zero.
    /// @param graduator_    Shared graduator that drains both reserves at graduation. Non-zero.
    function initialize(
        string memory name_,
        string memory symbol_,
        address creator_,
        address feeRouter_,
        address bondingCurve_,
        address graduator_
    ) external initializer {
        if (bytes(name_).length == 0 || bytes(symbol_).length == 0) revert EmptyMetadata();
        if (
            creator_ == address(0) || feeRouter_ == address(0) || bondingCurve_ == address(0)
                || graduator_ == address(0)
        ) revert ZeroAddress();

        __ERC20_init(name_, symbol_);
        __Ownable_init(creator_);

        creator = creator_;
        feeRouter = feeRouter_;
        bondingCurve = bondingCurve_;

        // Allowlist the three protocol-internal contracts that hop the supply during
        // bonding-curve trading and graduation. End-user transfers are still taxed.
        taxExempt[feeRouter_] = true;
        taxExempt[bondingCurve_] = true;
        taxExempt[graduator_] = true;

        _mint(bondingCurve_, TOTAL_SUPPLY);

        emit TaxExemptSet(feeRouter_, true);
        emit TaxExemptSet(bondingCurve_, true);
        emit TaxExemptSet(graduator_, true);
        emit AgentTokenInitialized(creator_, feeRouter_, bondingCurve_, name_, symbol_);
    }

    /// @dev Routes 1% of every peer-to-peer transfer to {feeRouter}. Tax is skipped for:
    ///      - mint (`from == address(0)`) so {TOTAL_SUPPLY} lands intact on the bonding curve,
    ///      - burn (`to == address(0)`) to preserve supply semantics,
    ///      - transfers where either counterparty is on the {taxExempt} allowlist (set
    ///        once at {initialize} for `feeRouter`, `bondingCurve`, and `graduator`).
    ///        This keeps the protocol-internal graduation hand-off lossless: the curve ->
    ///        graduator pull and the graduator -> Uniswap router transfer both skip tax,
    ///        so the router receives the full agent leg and `addLiquidity` does not
    ///        underflow. User-to-user transfers, including any path that does not touch
    ///        an exempt address, are still taxed in full.
    ///
    ///      Math: `tax = value * 100 / 10_000` is safe under Solidity 0.8 checked arithmetic
    ///      for any `value <= TOTAL_SUPPLY`. The remainder `value - tax` can never underflow
    ///      because `tax <= value`.
    function _update(address from, address to, uint256 value) internal override {
        // Skip tax on mint, burn, and any transfer touching an exempt address.
        if (from == address(0) || to == address(0) || taxExempt[from] || taxExempt[to]) {
            super._update(from, to, value);
            return;
        }

        uint256 tax = (value * TRADE_TAX_BPS) / BPS_DENOMINATOR;
        if (tax == 0) {
            // Value too small for tax to register; forward full amount (dust-safe).
            super._update(from, to, value);
            return;
        }

        super._update(from, feeRouter, tax);
        super._update(from, to, value - tax);

        emit TaxCollected(from, to, tax);
    }
}
