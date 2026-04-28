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

// TestACPEventSurface guards the indexer's M3 event subscriptions: regenerated
// bindings that drop one of these names would fail here before reaching CI.
//
// Each event entry also pins its topic0 (keccak256 of the canonical signature).
// A signature change on the Solidity side flips the hash — this test catches
// silent ABI drift even when the event name is preserved. Hashes were computed
// with `cast keccak "<sig>"`; the corresponding signature is documented inline.
func TestACPEventSurface(t *testing.T) {
	type ev struct {
		name string
		// topic0 = keccak256(canonical signature). Empty string skips the
		// hash assertion (use only for events whose signature is unstable).
		topic0 string
		// sig is the canonical signature used to derive topic0; carried as a
		// comment-equivalent so future readers can reproduce the hash.
		sig string
	}
	cases := []struct {
		name   string
		meta   *bind.MetaData
		events []ev
	}{
		{"AgentRegistry", bindings.AgentRegistryMetaData, []ev{
			{"AgentRegistered", "0xc29f819ac362ff9c94de06666235808451aafd8894a2dffb86a080a965efeae3", "AgentRegistered(uint256,address,string,uint256)"},
			{"MetadataUpdated", "0x459157ba24c7ab9878b165ef465fa6ae2ab42bcd8445f576be378768b0c47309", "MetadataUpdated(uint256,string)"},
			{"CapabilitiesUpdated", "0x2852a2d71dd37e385f378b5a9e262f9d2ed0cf51b7312a5f9229cc5cd39ea4c1", "CapabilitiesUpdated(uint256,uint256)"},
			{"ActiveStatusChanged", "0x632239944907e2f939f78bb06a2d210ac62cc0fbe3bff1b6244b4883721df129", "ActiveStatusChanged(uint256,bool)"},
			{"ScorePosted", "0xa19792167f6e5259706dc8bacefe76b932fe47d8538e98761711b78e49129e3a", "ScorePosted(uint256,address,int256,int256,uint256)"},
			{"ControllerTransferProposed", "0x41b8c1345e7590a4eb37c55b8c1cef180d9d1c694be1ae96d015bfe11a278e22", "ControllerTransferProposed(uint256,address)"},
			{"ControllerTransferAccepted", "0x969d6b1394c5147c28489f50960712e680e34e1e14ea24ff72bd52f315930c3a", "ControllerTransferAccepted(uint256,address,address)"},
			{"ControllerTransferCancelled", "0x5d4293db10304d02137eb59859f20bef23395587c83ae50ada99be07f7aa9aad", "ControllerTransferCancelled(uint256)"},
		}},
		{"Job", bindings.JobMetaData, []ev{
			{"JobInitialised", "0x6396f9d24200dbcd88fb1965753f38f76c647a2a0a793a00875957937869ce48", "JobInitialised(uint256,address,uint256,uint8,uint256,address,uint64)"},
			{"AgentAccepted", "0xa08a9d61f61bc14cb9b5ef024d79b5dbbb271dfebaaa5a22ed231c25d5d4ee8e", "AgentAccepted(uint256,address,uint256)"},
			{"ResultSubmitted", "0xc06b551d984e333ac851ab20b1454a08d92740468f52ff54c0cd5270817f20a9", "ResultSubmitted(uint256,address,string)"},
			{"ResultApproved", "0xdfa2a378c6a00b39f4d2b4247ba44e90a894960a78603b2bf5af5de40dbe183d", "ResultApproved(uint256,address)"},
			{"ResultRejected", "0x8e23b50499d65aea2724488d5fbc5924cafc75e6613cea60f4ec898b84123d90", "ResultRejected(uint256,address,string)"},
			{"JobCompleted", "0x4dbba472101e9f148a3a5ecbe793f0ee16a7efe4bf3f8cbfab5330e1642ef955", "JobCompleted(uint256,address,uint256)"},
			{"JobCancelled", "0x141454c3e63a845c54bcfebc9a4f9843c448cd8c0bb7077e5f70235e4f441356", "JobCancelled(uint256,address,string)"},
			{"DisputeRaised", "0x84a477df8a28a4276ca6dee4458a06c3015f30c477d9c949ede4e13ff8a552b4", "DisputeRaised(uint256,address)"},
			{"DisputeResolved", "0x8fdd4548a8481406b6e29c0d6f25e27cd72502f79f4adf409468502e7920dabc", "DisputeResolved(uint256,address,bool)"},
			{"EvaluatorAssigned", "0x50a93d710505e6f207121334c60e2a4c6312fdbae71f879f5abee6488e20b131", "EvaluatorAssigned(uint256,address)"},
		}},
		{"JobFactory", bindings.JobFactoryMetaData, []ev{
			{"JobCreated", "0xf956496ab0e37cf09cc3df7fd2b643c4a458b72e0d06e6cb9544a68553fabdd5", "JobCreated(uint256,address,address,uint8,address,uint256)"},
			{"ImplementationUpdated", "0xaa3f731066a578e5f39b4215468d826cdd15373cbc0dfc9cb9bdc649718ef7da", "ImplementationUpdated(address,address)"},
			{"DefaultArbiterUpdated", "0x8f4ec914918528531be3ff56f55a057f73eb540484e6fe6f63c7bd02a483a957", "DefaultArbiterUpdated(address,address)"},
		}},
		{"Escrow", bindings.EscrowMetaData, []ev{
			{"Funded", "0x6c1892d9d3736322680cc71c620c3fa6ef6e4302b30512327d66987eee7a602c", "Funded(address,uint256,address,uint256)"},
			{"Released", "0xf50168be4db3b63b5055f81fb26a4ca82261896dc1881b27a68920f764e44f0e", "Released(address,uint256,address,address,uint256)"},
			{"Refunded", "0x45751971a4d14bc4cd653a41eb84cff467a04fd5311450db55257120987005f7", "Refunded(address,uint256,address,uint256)"},
		}},
		{"FeeSplitter", bindings.FeeSplitterMetaData, []ev{
			{"FeeSplit", "0x396f3c0c854fc40380f8870c3e6c468a076d3510a63ecbf8c82c38fa647c098c", "FeeSplit(address,address,uint8,uint256,uint256,uint256)"},
			{"TreasuryUpdated", "0x4ab5be82436d353e61ca18726e984e561f5c1cc7c6d38b29d2553c790434705a", "TreasuryUpdated(address,address)"},
			{"BuybackBurnerUpdated", "0x2413e4696976dabb94237cb746c767cd0118c63bef644fdd08141569bee6331c", "BuybackBurnerUpdated(address,address)"},
		}},
		{"BuybackBurner", bindings.BuybackBurnerMetaData, []ev{
			{"BuybackAndBurn", "0x0f995e1de8b2d83fb2db107b6a0cb5b3a5bb1cca0d54dbcb01dd63f19ffed3bc", "BuybackAndBurn(address,address,uint256,uint256)"},
			{"RouterUpdated", "0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684", "RouterUpdated(address,address)"},
			{"SwapPathUpdated", "0x6ac8d8fd2cc5d6c950ab25024970e2ab3d81ff6e6e5ebeb21d992ee119b8e6ce", "SwapPathUpdated(address[])"},
			{"MinTituOutUpdated", "0x53fd101c8f6ba0f0fb4f44cc48b8e00b504b4f1e56073a9a170509f802488d84", "MinTituOutUpdated(uint256,uint256)"},
		}},
		{"HookRegistry", bindings.HookRegistryMetaData, []ev{
			{"HookRegistered", "0x8096ae10f94d9c9266fe844741548a8afe2a55900eb4607c7a5c3dd5560651a9", "HookRegistered(address,string)"},
			{"HookDeregistered", "0x5f04ed00bba892fd984432cd96beabe664a2de332ff480f08ec0c135c0372d55", "HookDeregistered(address)"},
		}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			parsed, err := tc.meta.GetAbi()
			if err != nil {
				t.Fatalf("GetAbi: %v", err)
			}
			for _, e := range tc.events {
				got, ok := parsed.Events[e.name]
				if !ok {
					t.Errorf("%s: missing event %q", tc.name, e.name)
					continue
				}
				if e.topic0 == "" {
					continue
				}
				if !strings.EqualFold(got.ID.Hex(), e.topic0) {
					t.Errorf("%s.%s: topic0 drift\n  want %s (%s)\n  got  %s",
						tc.name, e.name, e.topic0, e.sig, got.ID.Hex())
				}
			}
		})
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

		// M3 — ACP v2.
		{"AgentRegistry", bindings.AgentRegistryMetaData},
		{"BuybackBurner", bindings.BuybackBurnerMetaData},
		{"Escrow", bindings.EscrowMetaData},
		{"FeeSplitter", bindings.FeeSplitterMetaData},
		{"FundTransferHook", bindings.FundTransferHookMetaData},
		{"HookRegistry", bindings.HookRegistryMetaData},
		{"Job", bindings.JobMetaData},
		{"JobFactory", bindings.JobFactoryMetaData},
		{"MilestoneHook", bindings.MilestoneHookMetaData},
		{"RoyaltyHook", bindings.RoyaltyHookMetaData},
		{"SubscriptionHook", bindings.SubscriptionHookMetaData},
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
