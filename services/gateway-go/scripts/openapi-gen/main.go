// Command openapi-gen post-processes the swag-generated Swagger 2.0
// document into the OpenAPI 3.0 form actually served at runtime.
//
// Pipeline:
//
//   1. swag init   →  docs/swagger.json + docs/swagger.yaml (Swagger 2.0).
//   2. openapi-gen →  docs/openapi.json + docs/openapi.yaml (OpenAPI 3.0).
//
// We keep both forms on disk for two reasons:
//
//   * Swagger 2.0 is what swag emits and what the drift gate diffs.
//     Re-running `swag init` should produce a byte-stable file in CI.
//
//   * OpenAPI 3.0 is what SDK generators consume in 2025 and what the
//     gateway embeds + serves at /swagger/doc.{json,yaml}. The 3.0 form
//     is generated deterministically from the 2.0 form via
//     getkin/kin-openapi/openapi2conv.
//
// Custom scalar formats (BigInt, Address, Bytes32) are patched here
// because swag's annotation surface cannot describe them. Patches live
// in applyCustomScalars and are intentionally narrow: only types we own
// are touched.
//
// Usage:
//
//	go run ./scripts/openapi-gen
//
// Inputs:  services/gateway-go/docs/swagger.json
// Outputs: services/gateway-go/docs/openapi.json
//          services/gateway-go/docs/openapi.yaml
//
// Exit codes:
//
//	0 — files written, in sync with input
//	1 — input missing, malformed, or post-processing failed
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func ctxBackground() context.Context { return context.Background() }

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "openapi-gen:", err)
		os.Exit(1)
	}
}

func run() error {
	root, err := repoRoot()
	if err != nil {
		return err
	}
	docsDir := filepath.Join(root, "docs")

	in := filepath.Join(docsDir, "swagger.json")
	raw, err := os.ReadFile(in)
	if err != nil {
		return fmt.Errorf("read %s: %w", in, err)
	}

	var v2 openapi2.T
	if err := json.Unmarshal(raw, &v2); err != nil {
		return fmt.Errorf("parse swagger 2.0: %w", err)
	}

	v3, err := openapi2conv.ToV3(&v2)
	if err != nil {
		return fmt.Errorf("convert 2.0 → 3.0: %w", err)
	}

	applyCustomScalars(v3)

	if err := v3.Validate(ctxBackground()); err != nil {
		return fmt.Errorf("converted spec failed validation: %w", err)
	}

	jsonOut, err := marshalJSONStable(v3)
	if err != nil {
		return fmt.Errorf("marshal openapi.json: %w", err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "openapi.json"), jsonOut, 0o644); err != nil {
		return fmt.Errorf("write openapi.json: %w", err)
	}

	// Round-trip JSON → YAML so the YAML form is byte-stable across runs
	// and is structurally identical to the JSON form.
	var generic any
	if err := json.Unmarshal(jsonOut, &generic); err != nil {
		return fmt.Errorf("re-decode for yaml: %w", err)
	}
	yamlOut, err := yaml.Marshal(generic)
	if err != nil {
		return fmt.Errorf("marshal openapi.yaml: %w", err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "openapi.yaml"), yamlOut, 0o644); err != nil {
		return fmt.Errorf("write openapi.yaml: %w", err)
	}

	return nil
}

// applyCustomScalars patches the converted OpenAPI 3.0 document to give
// the protocol's bespoke string scalars precise format / pattern hints.
//
// swag has no annotation for custom formats (it can only emit standard
// OpenAPI primitives + a handful of well-known formats), so we layer
// the protocol-specific shape here. The patch is identity for any
// (schema, field) pair we don't recognise — adding a new wire type
// does NOT silently inherit an unrelated format hint.
//
// Mapping:
//
//	BigInt (uint256-as-decimal-string)  → format: bigint, pattern: ^[0-9]+$
//	Address (0x-hex 20 bytes)           → format: address, pattern: ^0x[a-fA-F0-9]{40}$
//	Bytes32 (0x-hex 32 bytes)           → format: bytes32, pattern: ^0x[a-fA-F0-9]{64}$
//	Time                                → format: date-time
//
// Classification is keyed by the (schemaName, fieldName) PAIR — not by
// fieldName alone. A field-name-only patcher would mis-classify
// homonyms across schemas: e.g. store.Job.token (ERC-20 contract addr,
// genuinely 0x-hex 20 bytes) and openapi.AuthVerifyResponse.token
// (opaque JWT bearer string) share the JSON name "token" but have
// completely different wire shapes. Keying on schema name avoids the
// collision and makes every classification explicit.
func applyCustomScalars(doc *openapi3.T) {
	if doc == nil || doc.Components == nil {
		return
	}
	for schemaName, schema := range doc.Components.Schemas {
		if schema == nil || schema.Value == nil {
			continue
		}
		for fieldName, prop := range schema.Value.Properties {
			if prop == nil || prop.Value == nil {
				continue
			}
			patchFieldFormat(schemaName, fieldName, prop.Value)
		}
	}
}

// scalarKind is the bespoke format we patch onto a (schema, field) pair.
type scalarKind int

const (
	scalarNone scalarKind = iota
	scalarAddress
	scalarBytes32
	scalarBigInt
	scalarDateTime
)

// scalarClassification is the explicit allowlist of (schema, field)
// pairs that get a custom format. Adding a new wire type means adding
// rows here — there is no name-only fallback. SIWE/SIWS shared response
// shapes (openapi.AuthVerifyRequest, openapi.AuthVerifyResponse) are
// intentionally absent because their string fields are heterogeneous:
//
//   - AuthVerifyResponse.token is a JWT, not an address.
//   - AuthVerifyResponse.address can be a 0x-hex EVM address (SIWE) or
//     a base58 Solana pubkey (SIWS); neither fits the "address" format
//     (which is EVM-specific, ^0x[a-fA-F0-9]{40}$). Leave it as plain
//     string so SDK consumers don't reject valid SIWS responses.
//   - AuthVerifyRequest.signature is hex (SIWE) or base58 (SIWS), same
//     polymorphism.
//
// See TestSpecLeavesAuthVerifyTokenUnformatted /
// TestSpecLeavesAuthVerifyAddressUnformatted for the regression guards.
var scalarClassification = map[string]map[string]scalarKind{
	"store.Agent": {
		"agent_id":     scalarBigInt,
		"capabilities": scalarBigInt,
		"controller":   scalarAddress,
		"creator":      scalarAddress,
		"tx_hash":      scalarBytes32,
		"created_at":   scalarDateTime,
		"updated_at":   scalarDateTime,
	},
	"store.Trade": {
		"trader":     scalarAddress,
		"quote_in":   scalarBigInt,
		"quote_out":  scalarBigInt,
		"agent_in":   scalarBigInt,
		"agent_out":  scalarBigInt,
		"fee":        scalarBigInt,
		"tx_hash":    scalarBytes32,
		"created_at": scalarDateTime,
	},
	"store.Job": {
		"agent_id":      scalarBigInt,
		"budget":        scalarBigInt,
		"clone_address": scalarAddress,
		"principal":     scalarAddress,
		"token":         scalarAddress,
		"tx_hash":       scalarBytes32,
		"created_at":    scalarDateTime,
		"updated_at":    scalarDateTime,
		"deadline":      scalarDateTime,
	},
}

// patchFieldFormat mutates s in place when (schemaName, fieldName) is
// in scalarClassification. The function is exercised by unit tests so
// the mapping table stays visible to assertions.
func patchFieldFormat(schemaName, fieldName string, s *openapi3.Schema) {
	// We only patch string fields; numeric / bool / object fields
	// inherit whatever the converter chose. An unset Type (nil
	// *Types) shows up on bare $ref children — leave those alone too.
	if s.Type == nil {
		return
	}
	if !s.Type.Is("string") {
		return
	}
	fields, ok := scalarClassification[schemaName]
	if !ok {
		return
	}
	kind, ok := fields[fieldName]
	if !ok {
		return
	}
	stringType := &openapi3.Types{"string"}
	switch kind {
	case scalarBytes32:
		s.Type = stringType
		s.Format = "bytes32"
		s.Pattern = `^0x[a-fA-F0-9]{64}$`
	case scalarAddress:
		s.Type = stringType
		s.Format = "address"
		s.Pattern = `^0x[a-fA-F0-9]{40}$`
	case scalarBigInt:
		s.Type = stringType
		s.Format = "bigint"
		s.Pattern = `^[0-9]+$`
	case scalarDateTime:
		s.Type = stringType
		s.Format = "date-time"
	}
}

// marshalJSONStable returns a deterministic, indent-2 JSON encoding of
// the spec. Determinism matters because the drift gate diffs the bytes;
// kin-openapi's default Marshal already sorts map keys.
func marshalJSONStable(doc *openapi3.T) ([]byte, error) {
	b, err := doc.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var generic any
	if err := json.Unmarshal(b, &generic); err != nil {
		return nil, err
	}
	out, err := json.MarshalIndent(generic, "", "  ")
	if err != nil {
		return nil, err
	}
	// Trailing newline matches go's gofmt convention and most editors'
	// auto-newline behaviour, keeping `git diff` quiet on round-trip.
	return append(out, '\n'), nil
}

// repoRoot resolves the gateway service root regardless of the cwd the
// script was launched from. We climb until we find the go.mod whose
// module path is the gateway, so e.g. `cd services/gateway-go &&
// go run ./scripts/openapi-gen` and `go run ./services/gateway-go/...`
// both work.
func repoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		mod := filepath.Join(dir, "go.mod")
		if data, err := os.ReadFile(mod); err == nil {
			if isGatewayModule(data) {
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find services/gateway-go from %s", wd)
		}
		dir = parent
	}
}

// isGatewayModule returns true when the given go.mod content declares
// the gateway-go module, anchored at the start of the file.
func isGatewayModule(data []byte) bool {
	prefix := []byte("module github.com/0xDevNinja/titular/services/gateway-go")
	if len(data) < len(prefix) {
		return false
	}
	return string(data[:len(prefix)]) == string(prefix)
}
