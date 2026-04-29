package main

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

// TestPatchFieldFormat is the unit-level coverage for the protocol
// scalar mapping. The map of names → formats is the only place a wire
// type can drift away from its documented form, so the test names
// every field explicitly.
func TestPatchFieldFormat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		fieldName   string
		inputType   string
		wantFormat  string
		wantPattern string
	}{
		// 32-byte digests
		{"tx_hash → bytes32", "tx_hash", "string", "bytes32", `^0x[a-fA-F0-9]{64}$`},

		// 20-byte addresses
		{"trader → address", "trader", "string", "address", `^0x[a-fA-F0-9]{40}$`},
		{"creator → address", "creator", "string", "address", `^0x[a-fA-F0-9]{40}$`},
		{"controller → address", "controller", "string", "address", `^0x[a-fA-F0-9]{40}$`},
		{"principal → address", "principal", "string", "address", `^0x[a-fA-F0-9]{40}$`},
		{"token → address", "token", "string", "address", `^0x[a-fA-F0-9]{40}$`},
		{"clone_address → address", "clone_address", "string", "address", `^0x[a-fA-F0-9]{40}$`},

		// uint256 decimals
		{"agent_id → bigint", "agent_id", "string", "bigint", `^[0-9]+$`},
		{"budget → bigint", "budget", "string", "bigint", `^[0-9]+$`},
		{"fee → bigint", "fee", "string", "bigint", `^[0-9]+$`},
		{"quote_in → bigint", "quote_in", "string", "bigint", `^[0-9]+$`},
		{"quote_out → bigint", "quote_out", "string", "bigint", `^[0-9]+$`},
		{"agent_in → bigint", "agent_in", "string", "bigint", `^[0-9]+$`},
		{"agent_out → bigint", "agent_out", "string", "bigint", `^[0-9]+$`},
		{"capabilities → bigint", "capabilities", "string", "bigint", `^[0-9]+$`},

		// timestamps
		{"created_at → date-time", "created_at", "string", "date-time", ""},
		{"updated_at → date-time", "updated_at", "string", "date-time", ""},
		{"deadline → date-time", "deadline", "string", "date-time", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := &openapi3.Schema{Type: &openapi3.Types{tc.inputType}}
			patchFieldFormat(tc.fieldName, s)
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
// hypothetical field with the wrong type or an unrecognised name MUST
// pass through untouched so the patcher can never silently mis-classify
// a future addition.
func TestPatchFieldFormat_LeavesUnknownAlone(t *testing.T) {
	t.Parallel()

	t.Run("non-string field", func(t *testing.T) {
		// agent_id is in the bigint case, but the patcher is supposed
		// to bail on non-string types since bigint encoding is by
		// definition decimal-string.
		s := &openapi3.Schema{Type: &openapi3.Types{"integer"}}
		patchFieldFormat("agent_id", s)
		if s.Format != "" {
			t.Errorf("format = %q; want untouched", s.Format)
		}
	})

	t.Run("unrecognised field", func(t *testing.T) {
		s := &openapi3.Schema{Type: &openapi3.Types{"string"}}
		patchFieldFormat("some_brand_new_field", s)
		if s.Format != "" {
			t.Errorf("format = %q; want untouched on unknown field", s.Format)
		}
	})

	t.Run("nil type", func(t *testing.T) {
		s := &openapi3.Schema{Type: nil}
		patchFieldFormat("tx_hash", s)
		if s.Format != "" {
			t.Errorf("format = %q; want untouched on nil-type schema", s.Format)
		}
	})
}
