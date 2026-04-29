package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

// TestPatchFieldFormat is the unit-level coverage for the protocol
// scalar mapping. The map of (schema, field) → format is the only place
// a wire type can drift away from its documented form, so the test
// names every pair explicitly.
func TestPatchFieldFormat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		schemaName  string
		fieldName   string
		inputType   string
		wantFormat  string
		wantPattern string
	}{
		// 32-byte digests
		{"store.Agent.tx_hash → bytes32", "store.Agent", "tx_hash", "string", "bytes32", `^0x[a-fA-F0-9]{64}$`},
		{"store.Trade.tx_hash → bytes32", "store.Trade", "tx_hash", "string", "bytes32", `^0x[a-fA-F0-9]{64}$`},
		{"store.Job.tx_hash → bytes32", "store.Job", "tx_hash", "string", "bytes32", `^0x[a-fA-F0-9]{64}$`},

		// 20-byte addresses
		{"store.Trade.trader → address", "store.Trade", "trader", "string", "address", `^0x[a-fA-F0-9]{40}$`},
		{"store.Agent.creator → address", "store.Agent", "creator", "string", "address", `^0x[a-fA-F0-9]{40}$`},
		{"store.Agent.controller → address", "store.Agent", "controller", "string", "address", `^0x[a-fA-F0-9]{40}$`},
		{"store.Job.principal → address", "store.Job", "principal", "string", "address", `^0x[a-fA-F0-9]{40}$`},
		{"store.Job.token → address", "store.Job", "token", "string", "address", `^0x[a-fA-F0-9]{40}$`},
		{"store.Job.clone_address → address", "store.Job", "clone_address", "string", "address", `^0x[a-fA-F0-9]{40}$`},

		// uint256 decimals
		{"store.Agent.agent_id → bigint", "store.Agent", "agent_id", "string", "bigint", `^[0-9]+$`},
		{"store.Job.agent_id → bigint", "store.Job", "agent_id", "string", "bigint", `^[0-9]+$`},
		{"store.Job.budget → bigint", "store.Job", "budget", "string", "bigint", `^[0-9]+$`},
		{"store.Trade.fee → bigint", "store.Trade", "fee", "string", "bigint", `^[0-9]+$`},
		{"store.Trade.quote_in → bigint", "store.Trade", "quote_in", "string", "bigint", `^[0-9]+$`},
		{"store.Trade.quote_out → bigint", "store.Trade", "quote_out", "string", "bigint", `^[0-9]+$`},
		{"store.Trade.agent_in → bigint", "store.Trade", "agent_in", "string", "bigint", `^[0-9]+$`},
		{"store.Trade.agent_out → bigint", "store.Trade", "agent_out", "string", "bigint", `^[0-9]+$`},
		{"store.Agent.capabilities → bigint", "store.Agent", "capabilities", "string", "bigint", `^[0-9]+$`},

		// timestamps
		{"store.Agent.created_at → date-time", "store.Agent", "created_at", "string", "date-time", ""},
		{"store.Trade.created_at → date-time", "store.Trade", "created_at", "string", "date-time", ""},
		{"store.Job.created_at → date-time", "store.Job", "created_at", "string", "date-time", ""},
		{"store.Agent.updated_at → date-time", "store.Agent", "updated_at", "string", "date-time", ""},
		{"store.Job.updated_at → date-time", "store.Job", "updated_at", "string", "date-time", ""},
		{"store.Job.deadline → date-time", "store.Job", "deadline", "string", "date-time", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := &openapi3.Schema{Type: &openapi3.Types{tc.inputType}}
			patchFieldFormat(tc.schemaName, tc.fieldName, s)
			if s.Format != tc.wantFormat {
				t.Errorf("format = %q; want %q", s.Format, tc.wantFormat)
			}
			if s.Pattern != tc.wantPattern {
				t.Errorf("pattern = %q; want %q", s.Pattern, tc.wantPattern)
			}
			if !s.Type.Is("string") {
				t.Errorf("type lost during patch: %v", s.Type.Slice())
			}
		})
	}
}

// TestPatchFieldFormat_LeavesUnknownAlone is the negative case. A
// hypothetical field with the wrong type, an unrecognised name, or a
// recognised name in a schema that does NOT list it MUST pass through
// untouched so the patcher can never silently mis-classify a future
// addition.
func TestPatchFieldFormat_LeavesUnknownAlone(t *testing.T) {
	t.Parallel()

	t.Run("non-string field", func(t *testing.T) {
		// agent_id is in the bigint case, but the patcher is supposed
		// to bail on non-string types since bigint encoding is by
		// definition decimal-string.
		s := &openapi3.Schema{Type: &openapi3.Types{"integer"}}
		patchFieldFormat("store.Agent", "agent_id", s)
		if s.Format != "" {
			t.Errorf("format = %q; want untouched", s.Format)
		}
	})

	t.Run("unrecognised field", func(t *testing.T) {
		s := &openapi3.Schema{Type: &openapi3.Types{"string"}}
		patchFieldFormat("store.Agent", "some_brand_new_field", s)
		if s.Format != "" {
			t.Errorf("format = %q; want untouched on unknown field", s.Format)
		}
	})

	t.Run("unrecognised schema", func(t *testing.T) {
		// "tx_hash" is a known scalar in store.* schemas, but if it
		// shows up on a schema we don't classify, bail. This is the
		// guard against name-only collisions across schemas.
		s := &openapi3.Schema{Type: &openapi3.Types{"string"}}
		patchFieldFormat("openapi.SomeNewType", "tx_hash", s)
		if s.Format != "" {
			t.Errorf("format = %q; want untouched on unknown schema", s.Format)
		}
	})

	t.Run("nil type", func(t *testing.T) {
		s := &openapi3.Schema{Type: nil}
		patchFieldFormat("store.Agent", "tx_hash", s)
		if s.Format != "" {
			t.Errorf("format = %q; want untouched on nil-type schema", s.Format)
		}
	})

	t.Run("AuthVerifyResponse not classified", func(t *testing.T) {
		// openapi.AuthVerifyResponse.token is a JWT, not an EVM
		// address. Even though "token" is the address-shaped field
		// name on store.Job, the patcher MUST NOT touch it here.
		s := &openapi3.Schema{Type: &openapi3.Types{"string"}}
		patchFieldFormat("openapi.AuthVerifyResponse", "token", s)
		if s.Format != "" || s.Pattern != "" {
			t.Errorf("AuthVerifyResponse.token leaked format=%q pattern=%q; want plain string",
				s.Format, s.Pattern)
		}
	})
}

// TestSpecLeavesAuthVerifyTokenUnformatted is the spec-level regression
// guard for blocker 1 of PR #304: AuthVerifyResponse.token is a JWT
// bearer string and MUST NOT inherit the EVM-address format hint that
// the legacy field-name-only patcher applied via the "token" key shared
// with store.Job.token (ERC-20 contract address).
func TestSpecLeavesAuthVerifyTokenUnformatted(t *testing.T) {
	t.Parallel()

	prop := readSpecProperty(t, "openapi.AuthVerifyResponse", "token")
	if prop["format"] != nil {
		t.Errorf("AuthVerifyResponse.token has format=%v; want unset (it is a JWT, not an address)",
			prop["format"])
	}
	if prop["pattern"] != nil {
		t.Errorf("AuthVerifyResponse.token has pattern=%v; want unset",
			prop["pattern"])
	}
	if prop["type"] != "string" {
		t.Errorf("AuthVerifyResponse.token has type=%v; want \"string\"", prop["type"])
	}
}

// TestSpecLeavesAuthVerifyAddressUnformatted is the spec-level
// regression guard for blocker 2 of PR #304: AuthVerifyResponse.address
// is heterogeneous (0x-hex EVM address for SIWE; base58 Solana pubkey
// for SIWS) and the EVM-only "address" format/pattern would reject
// otherwise-valid SIWS responses. Leave it as plain string.
func TestSpecLeavesAuthVerifyAddressUnformatted(t *testing.T) {
	t.Parallel()

	prop := readSpecProperty(t, "openapi.AuthVerifyResponse", "address")
	if prop["format"] != nil {
		t.Errorf("AuthVerifyResponse.address has format=%v; want unset (SIWE/SIWS polymorphic)",
			prop["format"])
	}
	if prop["pattern"] != nil {
		t.Errorf("AuthVerifyResponse.address has pattern=%v; want unset",
			prop["pattern"])
	}
	if prop["type"] != "string" {
		t.Errorf("AuthVerifyResponse.address has type=%v; want \"string\"", prop["type"])
	}
}

// TestSpecJobTokenStillFormattedAddress is the positive regression
// guard: store.Job.token IS an ERC-20 contract address and MUST keep
// the address format/pattern after the (schema, field) scoping change.
func TestSpecJobTokenStillFormattedAddress(t *testing.T) {
	t.Parallel()

	prop := readSpecProperty(t, "store.Job", "token")
	if got, want := prop["format"], "address"; got != want {
		t.Errorf("Job.token format = %v; want %q", got, want)
	}
	if got, want := prop["pattern"], `^0x[a-fA-F0-9]{40}$`; got != want {
		t.Errorf("Job.token pattern = %v; want %q", got, want)
	}
}

// readSpecProperty loads the generated docs/openapi.json and returns
// the property map under components.schemas[schemaName].properties[fieldName].
// Failing the test (rather than a separate skip) when the spec is
// missing is intentional: openapi.sh is a required pre-test step.
func readSpecProperty(t *testing.T, schemaName, fieldName string) map[string]any {
	t.Helper()
	root, err := repoRoot()
	if err != nil {
		t.Fatalf("repoRoot: %v", err)
	}
	specPath := filepath.Join(root, "docs", "openapi.json")
	raw, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("read %s: %v", specPath, err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parse %s: %v", specPath, err)
	}
	components, _ := doc["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)
	schema, ok := schemas[schemaName].(map[string]any)
	if !ok {
		t.Fatalf("schema %q missing from spec", schemaName)
	}
	props, _ := schema["properties"].(map[string]any)
	prop, ok := props[fieldName].(map[string]any)
	if !ok {
		t.Fatalf("property %s.%s missing from spec", schemaName, fieldName)
	}
	return prop
}
