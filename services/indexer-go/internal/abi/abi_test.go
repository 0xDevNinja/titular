package abi_test

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	bindings "github.com/0xDevNinja/titular/services/indexer-go/internal/abi"
)

// TestLaunchpadFactoryABIParse verifies that the embedded ABI JSON is valid
// and that the expected event signatures are present.
func TestLaunchpadFactoryABIParse(t *testing.T) {
	parsed, err := bindings.LaunchpadFactoryMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}

	// Event signatures we rely on in the indexer.
	wantEvents := []string{
		"AgentLaunched",
		"ModuleSet",
		"OwnershipTransferred",
	}
	for _, name := range wantEvents {
		if _, ok := parsed.Events[name]; !ok {
			t.Errorf("missing event %q in LaunchpadFactory ABI", name)
		}
	}
}

// TestLaunchpadFactoryEventTopics verifies keccak256 topic IDs are non-zero
// (confirming the ABI was generated from real contract source, not an empty stub).
func TestLaunchpadFactoryEventTopics(t *testing.T) {
	parsed, err := bindings.LaunchpadFactoryMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}

	for _, name := range []string{"AgentLaunched", "ModuleSet"} {
		ev, ok := parsed.Events[name]
		if !ok {
			t.Fatalf("event %q not found", name)
		}
		id := ev.ID
		if id == (common.Hash{}) {
			t.Errorf("event %q has zero topic ID", name)
		}
		t.Logf("%s topic: %s", name, id.Hex())
	}
}

// TestNewLaunchpadFactoryNilBackend verifies the binding constructor accepts a
// nil backend without panicking (useful for offline ABI decoding paths).
func TestNewLaunchpadFactoryNilBackend(t *testing.T) {
	addr := common.HexToAddress("0x1234000000000000000000000000000000000001")
	_, err := bindings.NewLaunchpadFactory(addr, nil)
	// abigen's NewXxx returns an error only when the ABI fails to parse;
	// a nil backend is accepted and errors only on actual RPC calls.
	if err != nil {
		t.Fatalf("NewLaunchpadFactory: %v", err)
	}
}

// TestBondingCurveABIParse checks the bonding curve ABI compiles cleanly.
func TestBondingCurveABIParse(t *testing.T) {
	parsed, err := bindings.BondingCurveMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}

	wantEvents := []string{"Bought", "Sold", "Graduated"}
	for _, name := range wantEvents {
		if _, ok := parsed.Events[name]; !ok {
			t.Errorf("missing event %q in BondingCurve ABI", name)
		}
	}
}

// TestAllBindingsCompile instantiates every MetaData and parses its ABI to
// catch any future regeneration that produces invalid JSON.
func TestAllBindingsCompile(t *testing.T) {
	table := []struct {
		name string
		meta *bind.MetaData
	}{
		{"AgentToken", bindings.AgentTokenMetaData},
		{"AirdropModule", bindings.AirdropModuleMetaData},
		{"AntiSniperModule", bindings.AntiSniperModuleMetaData},
		{"BondingCurve", bindings.BondingCurveMetaData},
		{"CapitalFormationModule", bindings.CapitalFormationModuleMetaData},
		{"ExistingTokenModule", bindings.ExistingTokenModuleMetaData},
		{"FeeDistributor", bindings.FeeDistributorMetaData},
		{"FeeRouter", bindings.FeeRouterMetaData},
		{"Graduator", bindings.GraduatorMetaData},
		{"LaunchRadarModule", bindings.LaunchRadarModuleMetaData},
		{"LaunchpadFactory", bindings.LaunchpadFactoryMetaData},
		{"LPLock", bindings.LPLockMetaData},
		{"PreBuyModule", bindings.PreBuyModuleMetaData},
		{"SixtyDaysModule", bindings.SixtyDaysModuleMetaData},
		{"TITU", bindings.TITUMetaData},
		{"Treasury", bindings.TreasuryMetaData},
		{"VeTITU", bindings.VeTITUMetaData},
		{"VestingVault", bindings.VestingVaultMetaData},
	}

	for _, tc := range table {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parsed, err := tc.meta.GetAbi()
			if err != nil {
				t.Fatalf("GetAbi: %v", err)
			}
			if strings.TrimSpace(tc.meta.ABI) == "" {
				t.Fatal("ABI string is empty")
			}
			// Sanity: at least one method or event exists.
			if len(parsed.Methods)+len(parsed.Events) == 0 {
				t.Fatal("ABI has no methods or events")
			}
		})
	}
}
