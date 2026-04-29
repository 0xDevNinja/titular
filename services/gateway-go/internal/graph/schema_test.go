package graph

import (
	"context"
	"testing"
)

// Schema-level tests pin the introspection result so a refactor that
// silently renames a type or field is caught at CI time. The set is
// the same one documented in schema.graphqls — the SDL is the
// human-readable contract; this test is the executable contract.

func Test_Introspection_RootTypes(t *testing.T) {
	h := newTestHandler(t, &fakeStore{})
	res := h.ExecuteForTest(context.Background(), `
{
  __schema {
    queryType { name }
    subscriptionType { name }
    types {
      name
    }
  }
}
`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("introspection errors: %v", res.Errors)
	}
	schema := res.Data.(map[string]any)["__schema"].(map[string]any)
	if schema["queryType"].(map[string]any)["name"] != "Query" {
		t.Errorf("queryType: got %v", schema["queryType"])
	}
	if schema["subscriptionType"].(map[string]any)["name"] != "Subscription" {
		t.Errorf("subscriptionType: got %v", schema["subscriptionType"])
	}

	// Asserting the entire types list would be brittle (graphql-go
	// adds many built-in __* types automatically). Instead, assert
	// the entity types are present.
	want := map[string]bool{
		"Agent": false, "Trade": false, "Job": false, "Stats": false,
		"AgentPage": false, "TradePage": false, "JobPage": false,
		"AgentKind": false, "JobStatus": false,
		"BigInt": false, "Address": false, "Bytes32": false, "Time": false,
	}
	for _, raw := range schema["types"].([]any) {
		name := raw.(map[string]any)["name"].(string)
		if _, ok := want[name]; ok {
			want[name] = true
		}
	}
	for n, ok := range want {
		if !ok {
			t.Errorf("type %q missing from schema", n)
		}
	}
}

func Test_Introspection_QueryFields(t *testing.T) {
	h := newTestHandler(t, &fakeStore{})
	res := h.ExecuteForTest(context.Background(), `
{
  __type(name: "Query") {
    fields {
      name
      args { name type { name kind ofType { name } } }
      type { name kind ofType { name kind } }
    }
  }
}
`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("introspection errors: %v", res.Errors)
	}
	fields := res.Data.(map[string]any)["__type"].(map[string]any)["fields"].([]any)
	got := make(map[string]bool)
	for _, raw := range fields {
		got[raw.(map[string]any)["name"].(string)] = true
	}
	for _, want := range []string{"agent", "agents", "trades", "jobs", "stats"} {
		if !got[want] {
			t.Errorf("Query.%s missing", want)
		}
	}
}

func Test_Introspection_SubscriptionFields(t *testing.T) {
	h := newTestHandler(t, &fakeStore{})
	res := h.ExecuteForTest(context.Background(), `
{
  __type(name: "Subscription") {
    fields { name }
  }
}
`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("introspection errors: %v", res.Errors)
	}
	fields := res.Data.(map[string]any)["__type"].(map[string]any)["fields"].([]any)
	got := make(map[string]bool)
	for _, raw := range fields {
		got[raw.(map[string]any)["name"].(string)] = true
	}
	for _, want := range []string{"newTrade", "jobUpdated", "agentRegistered"} {
		if !got[want] {
			t.Errorf("Subscription.%s missing", want)
		}
	}
}

// Test_Introspection_AgentFields pins the Agent object's field set.
// Catches the case where a refactor drops a field on store.Agent and
// the schema silently goes out of sync.
func Test_Introspection_AgentFields(t *testing.T) {
	h := newTestHandler(t, &fakeStore{})
	res := h.ExecuteForTest(context.Background(), `
{
  __type(name: "Agent") {
    fields { name }
  }
}
`, nil)
	if len(res.Errors) > 0 {
		t.Fatalf("introspection errors: %v", res.Errors)
	}
	fields := res.Data.(map[string]any)["__type"].(map[string]any)["fields"].([]any)
	got := make(map[string]bool)
	for _, raw := range fields {
		got[raw.(map[string]any)["name"].(string)] = true
	}
	for _, want := range []string{
		"id", "kind", "agentId", "controller", "creator", "metadataUri",
		"capabilities", "blockNumber", "txHash", "logIndex",
		"createdAt", "updatedAt",
	} {
		if !got[want] {
			t.Errorf("Agent.%s missing", want)
		}
	}
}

// Test_NewSchema_RejectsNilStore ensures the constructor refuses to
// build a schema without a store; this is a classic foot-gun where a
// nil-deref only surfaces at the first request.
func Test_NewSchema_RejectsNilStore(t *testing.T) {
	if _, err := NewSchema(&Resolver{Store: nil}); err == nil {
		t.Fatal("expected error for nil store")
	}
	if _, err := NewSchema(nil); err == nil {
		t.Fatal("expected error for nil resolver")
	}
}
