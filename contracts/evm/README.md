# contracts/evm

Solidity contracts for the Titular protocol. Built with Foundry.

## Quickstart

```bash
# Install dependencies (git submodules)
git submodule update --init --recursive

# Build
forge build

# Test
forge test -vvv

# Format
forge fmt
```

## Layout

```
src/        Solidity source files
test/       Foundry tests (.t.sol)
script/     Deployment scripts
lib/        Dependencies (forge-std, openzeppelin-contracts)
```

## Dependencies

- OpenZeppelin Contracts v5 (`lib/openzeppelin-contracts`)
- forge-std (`lib/forge-std`)

Remappings defined in `remappings.txt`.
