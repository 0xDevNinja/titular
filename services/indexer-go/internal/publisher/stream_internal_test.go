package publisher

import (
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

// White-box tests for the helpers EnsureStream uses to validate an
// already-existing stream. Going through the embedded nats-server is
// fine for the happy path (see publisher_test.go), but the server
// auto-fills Duplicates with its own default when the field is zero,
// so we exercise the helper directly with a synthetic [nats.StreamInfo]
// to pin the behaviour against a hand-built bad config (which is what
// we'd actually encounter on a shared cluster).

func Test_checkDeduplication_RejectsZero(t *testing.T) {
	info := &nats.StreamInfo{
		Config: nats.StreamConfig{
			Name:       StreamName,
			Subjects:   append([]string(nil), SubjectPrefixes...),
			Duplicates: 0, // <- the misconfiguration
		},
	}
	err := checkDeduplication(info)
	if err == nil {
		t.Fatal("checkDeduplication: expected error for Duplicates=0")
	}
	if !strings.Contains(err.Error(), "Duplicates") {
		t.Errorf("err = %v, want mention of Duplicates", err)
	}
}

func Test_checkDeduplication_AcceptsNonZero(t *testing.T) {
	info := &nats.StreamInfo{
		Config: nats.StreamConfig{
			Name:       StreamName,
			Duplicates: time.Minute,
		},
	}
	if err := checkDeduplication(info); err != nil {
		t.Errorf("checkDeduplication: unexpected error: %v", err)
	}
}

// subjectCovers handles the broader-match case: an existing stream with
// `agents.>` should satisfy a `agents.*` need; a single `>` covers
// everything; an unrelated prefix does not.
func Test_subjectCovers(t *testing.T) {
	cases := []struct {
		name string
		have []string
		want string
		ok   bool
	}{
		{"exact match", []string{"agents.>"}, "agents.>", true},
		{"catch-all >", []string{">"}, "agents.>", true},
		{"broader prefix covers narrower", []string{"agents.>"}, "agents.*", true},
		{"same-token concrete subject covered", []string{"agents.>"}, "agents.foo", true},
		{"unrelated prefix", []string{"agents.>"}, "jobs.>", false},
		{"empty have", []string{}, "agents.>", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := subjectCovers(c.have, c.want); got != c.ok {
				t.Errorf("subjectCovers(%v, %q) = %v, want %v", c.have, c.want, got, c.ok)
			}
		})
	}
}
