//go:build integration

package integration

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// emitterRuntimeBytecode is the runtime bytecode of the test-only
// `Emitter` contract:
//
//	pragma solidity 0.8.25;
//	contract Emitter {
//	    event CapabilitiesUpdated(uint256 indexed agentId, uint256 capabilities);
//	    function emitOne(uint256 agentId, uint256 capabilities) external {
//	        emit CapabilitiesUpdated(agentId, capabilities);
//	    }
//	}
//
// Compiled with `solc --optimize --bin-runtime`. We embed the runtime
// bytes (not the deploy bytes) so the test installs the contract via
// `anvil_setCode` without paying for a deploy transaction. The bytecode
// is short enough (~200 bytes) that pinning it inline is cheaper than
// dragging a forge build into the integration suite.
//
// Re-deriving the bytecode (when bumping solc): run
//
//	solc --optimize --bin-runtime contracts/test/Emitter.sol
//
// and copy the "Binary of the runtime part" into this constant. The
// embedded topic0
//
//	0x2852a2d71dd37e385f378b5a9e262f9d2ed0cf51b7312a5f9229cc5cd39ea4c1
//
// is `keccak256("CapabilitiesUpdated(uint256,uint256)")` and matches the
// real AgentRegistry ABI in `internal/abi/acp_agent_registry.go`. Tests
// assert on this equality so a compiler/event-signature drift surfaces
// as a clean failure rather than as an unmapped-event skip.
const emitterRuntimeHex = "" +
	"6080604052348015600e575f80fd5b50600436106026575f3560e01c806378d81c5214602a575b5f80fd5b603960353660046078565b603b565b005b817f2852a2d71dd37e385f378b5a9e262f9d2ed0cf51b7312a5f9229cc5cd39ea4c182604051606c91815260200190565b60405180910390a25050565b5f80604083850312156088575f80fd5b5050803592602090910135915056fea2646970667358221220d9a9f6be3aee55cf182adea87c4b9edace7c52102299d172f287d844158a909964736f6c63430008190033"

// emitterRuntimeCode decodes [emitterRuntimeHex] once at package init so
// the per-test setUp path does not redo the work. Panics on a bad const,
// which can only happen if a developer hand-edits the hex string into
// invalid form — the test suite would fail at the first import either
// way, so a panic at init is the clearer error.
var emitterRuntimeCode = mustHex(emitterRuntimeHex)

func mustHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic("integration: bad emitterRuntimeHex: " + err.Error())
	}
	return b
}

// emitterAddress is the deterministic address used by every test in
// this package. Picked far from the dev-account range so debugging
// receipts is unambiguous (`0xE...` is "the emitter", not "account 9").
var emitterAddress = common.HexToAddress("0xE0000000000000000000000000000000000000E0")

// emitOneSelector is the 4-byte function selector for
// `emitOne(uint256,uint256)`. We compute it at init rather than pin it
// to avoid the foot-gun where a future renamed function silently
// matches the cached selector.
var emitOneSelector = func() [4]byte {
	h := crypto.Keccak256([]byte("emitOne(uint256,uint256)"))
	var out [4]byte
	copy(out[:], h[:4])
	return out
}()
