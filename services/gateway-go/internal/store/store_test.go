package store_test

import (
	"strings"
	"testing"
	"time"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
)

func TestClampLimit(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{0, store.DefaultLimit},
		{-5, store.DefaultLimit},
		{1, 1},
		{store.MaxLimit, store.MaxLimit},
		{store.MaxLimit + 1, store.MaxLimit},
		{1_000_000, store.MaxLimit},
	}
	for _, c := range cases {
		if got := store.ClampLimit(c.in); got != c.want {
			t.Errorf("ClampLimit(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestCursorRoundTrip(t *testing.T) {
	in := store.Cursor{
		ID:        42,
		CreatedAt: time.Date(2026, 4, 28, 13, 0, 0, 0, time.UTC),
	}
	encoded := store.EncodeCursor(in)
	if encoded == "" {
		t.Fatal("encoded cursor must not be empty")
	}
	out, err := store.DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.ID != in.ID || !out.CreatedAt.Equal(in.CreatedAt) {
		t.Fatalf("round-trip mismatch: got %+v want %+v", out, in)
	}
}

func TestDecodeCursor_EmptyIsZero(t *testing.T) {
	c, err := store.DecodeCursor("")
	if err != nil {
		t.Fatalf("empty cursor must decode without error: %v", err)
	}
	if !c.CreatedAt.IsZero() || c.ID != 0 {
		t.Fatalf("empty cursor should be zero, got %+v", c)
	}
}

func TestDecodeCursor_RejectsGarbage(t *testing.T) {
	cases := []string{
		"not-base64!",
		"AAAA", // valid base64 but not valid JSON
	}
	for _, s := range cases {
		s := s
		t.Run(s, func(t *testing.T) {
			_, err := store.DecodeCursor(s)
			if err == nil {
				t.Fatalf("expected error for %q", s)
			}
			if !strings.HasPrefix(err.Error(), "cursor ") {
				t.Fatalf("error must start with 'cursor ' so the handler can map to 400, got %q", err.Error())
			}
		})
	}
}

func TestAgentKindValid(t *testing.T) {
	cases := map[store.AgentKind]bool{
		store.AgentKindLaunchpad: true,
		store.AgentKindACP:       true,
		"":                       false,
		"unknown":                false,
		"ACP":                    false, // case-sensitive on purpose
	}
	for k, want := range cases {
		if got := k.Valid(); got != want {
			t.Errorf("AgentKind(%q).Valid() = %v, want %v", k, got, want)
		}
	}
}

func TestJobPhaseValid(t *testing.T) {
	good := []store.JobPhase{
		store.JobPhaseCreated, store.JobPhaseFunded, store.JobPhaseActive,
		store.JobPhaseCompleted, store.JobPhaseCancelled, store.JobPhaseDisputed,
		store.JobPhaseReleased, store.JobPhaseResolved,
	}
	for _, p := range good {
		if !p.Valid() {
			t.Errorf("phase %q should be valid", p)
		}
	}
	bad := []store.JobPhase{"", "ARCHIVED", "active "}
	for _, p := range bad {
		if p.Valid() {
			t.Errorf("phase %q should be invalid", p)
		}
	}
}
