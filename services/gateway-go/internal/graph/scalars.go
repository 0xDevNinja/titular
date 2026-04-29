// Package graph hosts the GraphQL surface for the gateway: the schema, the
// resolver functions, and the NATS-backed subscription stream.
//
// We use github.com/graphql-go/graphql as the runtime because it is
// dependency-light, type-safe at runtime, and avoids a code-generation
// step that would otherwise need to be wired into CI. The shape of the
// schema is documented separately in schema.graphqls (SDL form) so SDK
// consumers and tooling that want introspection-by-file have a stable
// reference point. The two MUST stay in sync: scalars_test.go and
// schema_test.go pin the field set against the SDL via introspection.
//
// # Why not gqlgen
//
// The phase brief originally suggested gqlgen. We picked graphql-go
// instead because:
//
//   - gqlgen requires a code-gen step that produces ~30k lines of
//     vendored Go in the repo, which is a maintenance tax for an alpha
//     surface that is still settling.
//   - Our schema is small (3 entity types, 5 queries, 3 subscriptions)
//     and resolves directly off store.Store + a NATS bus, neither of
//     which benefits from gqlgen's directive/middleware machinery.
//   - graphql-go gives us run-time schema introspection out of the
//     box, which is exactly what the SDK alpha consumes; gqlgen's
//     compile-time guarantees are wasted here because the resolvers
//     are thin.
//
// If a future phase needs typed resolvers / federation / dataloaders,
// we re-evaluate; the schema.graphqls file is the migration anchor.
package graph

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// bigIntScalar is the wire form for uint256-sized integers. They serialise
// as quoted decimal strings to survive JS clients (which lose precision
// past 2^53). The same encoding is used by the NATS publisher for
// chain-decoded events, so a value flowing through subscriptions matches
// the value flowing through the REST API byte-for-byte.
//
// Parsing accepts both string and integer literals so a query can say
// `agentId: 7` or `agentId: "7"` interchangeably; we always emit strings.
var bigIntScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "BigInt",
	Description: "Arbitrary-precision integer encoded as a decimal string.",
	Serialize: func(value any) any {
		switch v := value.(type) {
		case nil:
			return nil
		case string:
			// Already in canonical form; trust the producer (the store).
			return v
		case *big.Int:
			if v == nil {
				return nil
			}
			return v.Text(10)
		case big.Int:
			return v.Text(10)
		case int:
			return fmt.Sprintf("%d", v)
		case int64:
			return fmt.Sprintf("%d", v)
		case uint64:
			return fmt.Sprintf("%d", v)
		}
		return nil
	},
	ParseValue: func(value any) any {
		switch v := value.(type) {
		case string:
			if _, ok := new(big.Int).SetString(v, 10); !ok {
				return nil
			}
			return v
		case int:
			return fmt.Sprintf("%d", v)
		case int64:
			return fmt.Sprintf("%d", v)
		case float64:
			// JSON numerics arrive as float64; reject anything with a
			// fractional component to avoid silent truncation.
			if v != float64(int64(v)) {
				return nil
			}
			return fmt.Sprintf("%d", int64(v))
		}
		return nil
	},
	ParseLiteral: func(valueAST ast.Value) any {
		switch v := valueAST.(type) {
		case *ast.StringValue:
			if _, ok := new(big.Int).SetString(v.Value, 10); !ok {
				return nil
			}
			return v.Value
		case *ast.IntValue:
			return v.Value
		}
		return nil
	},
})

// addressScalar pins the wire form of a 20-byte EVM address to lowercase
// 0x-prefixed hex. We accept upper- and mixed-case input (EIP-55 is
// allowed but not required) and always emit lowercase, matching the
// store layer.
var addressScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Address",
	Description: "Lowercase 0x-prefixed 20-byte EVM address.",
	Serialize: func(value any) any {
		s, ok := value.(string)
		if !ok || s == "" {
			return nil
		}
		return strings.ToLower(s)
	},
	ParseValue: func(value any) any {
		s, ok := value.(string)
		if !ok {
			return nil
		}
		if !validHex(s, 20) {
			return nil
		}
		return strings.ToLower(s)
	},
	ParseLiteral: func(valueAST ast.Value) any {
		v, ok := valueAST.(*ast.StringValue)
		if !ok {
			return nil
		}
		if !validHex(v.Value, 20) {
			return nil
		}
		return strings.ToLower(v.Value)
	},
})

// bytes32Scalar is the wire form for fixed-32-byte payloads (transaction
// hashes, content hashes). Mirrors addressScalar's case-folding behaviour.
var bytes32Scalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Bytes32",
	Description: "Lowercase 0x-prefixed 32-byte hex string.",
	Serialize: func(value any) any {
		s, ok := value.(string)
		if !ok || s == "" {
			return nil
		}
		return strings.ToLower(s)
	},
	ParseValue: func(value any) any {
		s, ok := value.(string)
		if !ok {
			return nil
		}
		if !validHex(s, 32) {
			return nil
		}
		return strings.ToLower(s)
	},
	ParseLiteral: func(valueAST ast.Value) any {
		v, ok := valueAST.(*ast.StringValue)
		if !ok {
			return nil
		}
		if !validHex(v.Value, 32) {
			return nil
		}
		return strings.ToLower(v.Value)
	},
})

// timeScalar serialises time.Time values as RFC 3339 with UTC offset.
// Inbound values may be RFC 3339 strings; we reject anything else so a
// caller sending a Unix timestamp gets a clear validation error rather
// than a silently-wrong filter.
var timeScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Time",
	Description: "RFC 3339 UTC timestamp.",
	Serialize: func(value any) any {
		switch v := value.(type) {
		case time.Time:
			if v.IsZero() {
				return nil
			}
			return v.UTC().Format(time.RFC3339Nano)
		case *time.Time:
			if v == nil || v.IsZero() {
				return nil
			}
			return v.UTC().Format(time.RFC3339Nano)
		case string:
			return v
		}
		return nil
	},
	ParseValue: func(value any) any {
		s, ok := value.(string)
		if !ok {
			return nil
		}
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			t, err = time.Parse(time.RFC3339, s)
			if err != nil {
				return nil
			}
		}
		return t
	},
	ParseLiteral: func(valueAST ast.Value) any {
		v, ok := valueAST.(*ast.StringValue)
		if !ok {
			return nil
		}
		t, err := time.Parse(time.RFC3339Nano, v.Value)
		if err != nil {
			t, err = time.Parse(time.RFC3339, v.Value)
			if err != nil {
				return nil
			}
		}
		return t
	},
})

// validHex reports whether s is a 0x-prefixed lowercase-or-mixed-case
// hex string of exactly want bytes (so 2*want hex characters after the
// prefix). We do the check in Go rather than via a regexp because the
// scalar parser runs on every query — no per-call regexp compile.
func validHex(s string, want int) bool {
	if len(s) != 2+2*want {
		return false
	}
	if s[0] != '0' || (s[1] != 'x' && s[1] != 'X') {
		return false
	}
	for i := 2; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
		case c >= 'a' && c <= 'f':
		case c >= 'A' && c <= 'F':
		default:
			return false
		}
	}
	return true
}

// errInvalidScalar is the user-facing message emitted when a scalar fails
// validation. We keep it generic to avoid leaking which validator branch
// rejected the value; the GraphQL response carries the field path so
// clients still know where to look.
var errInvalidScalar = errors.New("invalid scalar value")
