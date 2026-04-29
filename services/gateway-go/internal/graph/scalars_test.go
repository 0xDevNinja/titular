package graph

import (
	"math/big"
	"testing"
	"time"

	"github.com/graphql-go/graphql/language/ast"
)

// Scalar tests pin the wire form for BigInt, Address, Bytes32 and Time.
// The wire form is part of the public contract: SDK consumers, the
// REST surface, and the NATS publisher all rely on identical encoding
// for these types, so any drift here is an interop break.

func Test_BigInt_Serialize(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want any
	}{
		{"nil", nil, nil},
		{"string", "12345", "12345"},
		{"big.Int ptr", big.NewInt(42), "42"},
		{"big.Int val", *big.NewInt(7), "7"},
		{"int", 100, "100"},
		{"int64", int64(1<<40), "1099511627776"},
		{"large > 2^53", new(big.Int).Lsh(big.NewInt(1), 200), "1606938044258990275541962092341162602522202993782792835301376"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := bigIntScalar.Serialize(c.in)
			if got != c.want {
				t.Errorf("Serialize(%v) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func Test_BigInt_ParseValue(t *testing.T) {
	cases := []struct {
		in   any
		want any
	}{
		{"42", "42"},
		{42, "42"},
		{int64(99), "99"},
		{float64(7), "7"},
		{float64(7.5), nil}, // fractional must reject
		{"not a number", nil},
		{true, nil},
	}
	for _, c := range cases {
		got := bigIntScalar.ParseValue(c.in)
		if got != c.want {
			t.Errorf("ParseValue(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func Test_BigInt_ParseLiteral(t *testing.T) {
	if got := bigIntScalar.ParseLiteral(&ast.IntValue{Value: "42"}); got != "42" {
		t.Errorf("IntValue: got %v", got)
	}
	if got := bigIntScalar.ParseLiteral(&ast.StringValue{Value: "1234567890"}); got != "1234567890" {
		t.Errorf("StringValue: got %v", got)
	}
	if got := bigIntScalar.ParseLiteral(&ast.StringValue{Value: "abc"}); got != nil {
		t.Errorf("non-numeric string: got %v, want nil", got)
	}
	if got := bigIntScalar.ParseLiteral(&ast.BooleanValue{Value: true}); got != nil {
		t.Errorf("BooleanValue: got %v, want nil", got)
	}
}

func Test_Address_Serialize_Lowercases(t *testing.T) {
	got := addressScalar.Serialize("0xABCDEF0000000000000000000000000000000000")
	if got != "0xabcdef0000000000000000000000000000000000" {
		t.Errorf("got %v", got)
	}
}

func Test_Address_ParseValue(t *testing.T) {
	cases := []struct {
		in   any
		want any
	}{
		{"0xabcdef0000000000000000000000000000000000", "0xabcdef0000000000000000000000000000000000"},
		{"0xABCDEF0000000000000000000000000000000000", "0xabcdef0000000000000000000000000000000000"},
		{"abcdef", nil}, // missing prefix
		{"0xZZZZ000000000000000000000000000000000000", nil},
		{"0xabc", nil}, // too short
		{42, nil},
	}
	for _, c := range cases {
		got := addressScalar.ParseValue(c.in)
		if got != c.want {
			t.Errorf("ParseValue(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func Test_Bytes32_Validates32Bytes(t *testing.T) {
	good := "0x" + repeatA(64)
	if got := bytes32Scalar.ParseValue(good); got != good {
		t.Errorf("32-byte ok: got %v", got)
	}
	if got := bytes32Scalar.ParseValue("0x" + repeatA(60)); got != nil {
		t.Errorf("short hash: got %v, want nil", got)
	}
}

func Test_Time_RoundTrip(t *testing.T) {
	now := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
	wire := timeScalar.Serialize(now)
	got := timeScalar.ParseValue(wire)
	gotT, ok := got.(time.Time)
	if !ok {
		t.Fatalf("ParseValue: got %T", got)
	}
	if !gotT.Equal(now) {
		t.Errorf("round-trip mismatch: got %v, want %v", gotT, now)
	}
}

func Test_Time_RejectsGarbage(t *testing.T) {
	if got := timeScalar.ParseValue("yesterday"); got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

// validHex is exercised via the scalars but pin the explicit edge
// cases here so a future refactor that inlines or rewrites it doesn't
// silently break the contract.
func Test_validHex(t *testing.T) {
	cases := []struct {
		in   string
		want bool
		size int
	}{
		{"0x" + repeatA(40), true, 20},
		{"0X" + repeatA(40), true, 20}, // uppercase prefix tolerated
		{"00" + repeatA(40), false, 20},
		{"0x" + repeatA(38), false, 20},
		{"0xZZ" + repeatA(38), false, 20},
		{"", false, 20},
	}
	for _, c := range cases {
		if got := validHex(c.in, c.size); got != c.want {
			t.Errorf("validHex(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

// repeatA returns "aaaa..." n times. Local helper so tests don't import
// strings just for this.
func repeatA(n int) string {
	out := make([]byte, n)
	for i := range out {
		out[i] = 'a'
	}
	return string(out)
}
