// This file is hand-written and lives alongside the swag-generated
// docs.go. It exposes the OpenAPI 3.0 form of the gateway specification
// (produced by `go run ./scripts/openapi-gen`) as embedded byte slices
// the runtime serves at /swagger/doc.{json,yaml}.
//
// Why two artefacts on disk:
//
//   - docs/swagger.{json,yaml} — Swagger 2.0; what swag emits, what the
//     drift gate diffs.
//   - docs/openapi.{json,yaml} — OpenAPI 3.0; converted from Swagger 2.0
//     by scripts/openapi-gen, served at runtime, consumed by SDK
//     generators (#92 onwards).
//
// Both forms MUST be present at build time; CI re-runs `swag init` AND
// `go run ./scripts/openapi-gen` and fails on `git diff --exit-code`.

package docs

import _ "embed"

// SpecJSON is the OpenAPI 3.0 form of the gateway API document.
//
//go:embed openapi.json
var SpecJSON []byte

// SpecYAML is the OpenAPI 3.0 form of the gateway API document, in YAML.
// Identical surface to SpecJSON; the YAML is the human-readable form
// for code review.
//
//go:embed openapi.yaml
var SpecYAML []byte
