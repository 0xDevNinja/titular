package graph

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql"

	"github.com/0xDevNinja/titular/services/gateway-go/internal/store"
)

// agentKindEnum mirrors store.AgentKind. The values list is closed and
// case-sensitive on the wire; clients sending an unknown value get a
// GraphQL validation error before the resolver runs.
var agentKindEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "AgentKind",
	Description: "Origin of the agent (launchpad bonding curve or ACP factory).",
	Values: graphql.EnumValueConfigMap{
		"launchpad": {Value: store.AgentKindLaunchpad},
		"acp":       {Value: store.AgentKindACP},
	},
})

// jobStatusEnum mirrors store.JobPhase. Same closed-set semantics.
var jobStatusEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "JobStatus",
	Description: "Current phase of an ACP job state machine.",
	Values: graphql.EnumValueConfigMap{
		"created":   {Value: store.JobPhaseCreated},
		"funded":    {Value: store.JobPhaseFunded},
		"active":    {Value: store.JobPhaseActive},
		"completed": {Value: store.JobPhaseCompleted},
		"cancelled": {Value: store.JobPhaseCancelled},
		"disputed":  {Value: store.JobPhaseDisputed},
		"released":  {Value: store.JobPhaseReleased},
		"resolved":  {Value: store.JobPhaseResolved},
	},
})

// agentType is the GraphQL representation of store.Agent.
//
// Field resolvers extract the corresponding struct field — we do not
// rely on graphql-go's reflection mode because it does not understand
// the store.Agent JSON tags / our scalar bindings. Explicit resolvers
// also let us format BigInt values consistently and silently elide
// optional fields when empty.
var agentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Agent",
	Fields: graphql.Fields{
		"id":           {Type: graphql.NewNonNull(graphql.ID), Resolve: agentField("id")},
		"kind":         {Type: graphql.NewNonNull(agentKindEnum), Resolve: agentField("kind")},
		"agentId":      {Type: bigIntScalar, Resolve: agentField("agentId")},
		"controller":   {Type: addressScalar, Resolve: agentField("controller")},
		"creator":      {Type: addressScalar, Resolve: agentField("creator")},
		"metadataUri":  {Type: graphql.String, Resolve: agentField("metadataUri")},
		"capabilities": {Type: bigIntScalar, Resolve: agentField("capabilities")},
		"blockNumber":  {Type: graphql.NewNonNull(graphql.Int), Resolve: agentField("blockNumber")},
		"txHash":       {Type: graphql.NewNonNull(bytes32Scalar), Resolve: agentField("txHash")},
		"logIndex":     {Type: graphql.NewNonNull(graphql.Int), Resolve: agentField("logIndex")},
		"createdAt":    {Type: graphql.NewNonNull(timeScalar), Resolve: agentField("createdAt")},
		"updatedAt":    {Type: graphql.NewNonNull(timeScalar), Resolve: agentField("updatedAt")},
	},
})

// tradeType mirrors store.Trade.
var tradeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Trade",
	Fields: graphql.Fields{
		"id":           {Type: graphql.NewNonNull(graphql.ID), Resolve: tradeField("id")},
		"agentTokenPK": {Type: graphql.Int, Resolve: tradeField("agentTokenPK")},
		"side":         {Type: graphql.NewNonNull(graphql.String), Resolve: tradeField("side")},
		"trader":       {Type: graphql.NewNonNull(addressScalar), Resolve: tradeField("trader")},
		"curve":        {Type: graphql.NewNonNull(addressScalar), Resolve: tradeField("curve")},
		"quoteIn":      {Type: bigIntScalar, Resolve: tradeField("quoteIn")},
		"agentOut":     {Type: bigIntScalar, Resolve: tradeField("agentOut")},
		"agentIn":      {Type: bigIntScalar, Resolve: tradeField("agentIn")},
		"quoteOut":     {Type: bigIntScalar, Resolve: tradeField("quoteOut")},
		"fee":          {Type: graphql.NewNonNull(bigIntScalar), Resolve: tradeField("fee")},
		"blockNumber":  {Type: graphql.NewNonNull(graphql.Int), Resolve: tradeField("blockNumber")},
		"txHash":       {Type: graphql.NewNonNull(bytes32Scalar), Resolve: tradeField("txHash")},
		"logIndex":     {Type: graphql.NewNonNull(graphql.Int), Resolve: tradeField("logIndex")},
		"createdAt":    {Type: graphql.NewNonNull(timeScalar), Resolve: tradeField("createdAt")},
	},
})

// jobType mirrors store.Job.
var jobType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Job",
	Fields: graphql.Fields{
		"id":           {Type: graphql.NewNonNull(graphql.ID), Resolve: jobField("id")},
		"onChainId":    {Type: graphql.NewNonNull(bigIntScalar), Resolve: jobField("onChainId")},
		"cloneAddress": {Type: addressScalar, Resolve: jobField("cloneAddress")},
		"agentPK":      {Type: graphql.Int, Resolve: jobField("agentPK")},
		"agentId":      {Type: bigIntScalar, Resolve: jobField("agentId")},
		"principal":    {Type: graphql.NewNonNull(addressScalar), Resolve: jobField("principal")},
		"jobType":      {Type: graphql.NewNonNull(graphql.Int), Resolve: jobField("jobType")},
		"token":        {Type: addressScalar, Resolve: jobField("token")},
		"budget":       {Type: graphql.NewNonNull(bigIntScalar), Resolve: jobField("budget")},
		"deadline":     {Type: timeScalar, Resolve: jobField("deadline")},
		"phase":        {Type: graphql.NewNonNull(jobStatusEnum), Resolve: jobField("phase")},
		"resultUri":    {Type: graphql.String, Resolve: jobField("resultUri")},
		"blockNumber":  {Type: graphql.NewNonNull(graphql.Int), Resolve: jobField("blockNumber")},
		"txHash":       {Type: graphql.NewNonNull(bytes32Scalar), Resolve: jobField("txHash")},
		"logIndex":     {Type: graphql.NewNonNull(graphql.Int), Resolve: jobField("logIndex")},
		"createdAt":    {Type: graphql.NewNonNull(timeScalar), Resolve: jobField("createdAt")},
		"updatedAt":    {Type: graphql.NewNonNull(timeScalar), Resolve: jobField("updatedAt")},
	},
})

var statsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Stats",
	Fields: graphql.Fields{
		"totalAgents": {Type: graphql.NewNonNull(graphql.Int), Resolve: func(p graphql.ResolveParams) (any, error) {
			return p.Source.(store.Stats).TotalAgents, nil
		}},
		"totalJobs": {Type: graphql.NewNonNull(graphql.Int), Resolve: func(p graphql.ResolveParams) (any, error) {
			return p.Source.(store.Stats).TotalJobs, nil
		}},
		"totalTrades": {Type: graphql.NewNonNull(graphql.Int), Resolve: func(p graphql.ResolveParams) (any, error) {
			return p.Source.(store.Stats).TotalTrades, nil
		}},
		"lastBlockIndexed": {Type: graphql.NewNonNull(graphql.Int), Resolve: func(p graphql.ResolveParams) (any, error) {
			return p.Source.(store.Stats).LastBlockIndexed, nil
		}},
	},
})

// pageType returns a closed Page<T> object around the given inner type.
// Each call constructs a distinct GraphQL type because the runtime
// rejects two types sharing a name; the inner type's name is used as a
// suffix to keep them stable across schema reflection.
func pageType(name string, inner *graphql.Object, dataResolver, cursorResolver graphql.FieldResolveFn) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"data":       {Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(inner))), Resolve: dataResolver},
			"nextCursor": {Type: graphql.String, Resolve: cursorResolver},
		},
	})
}

// agentField returns a resolver for the named field on an Agent source.
// We hand-pick fields rather than relying on reflection so the wire
// names use camelCase (matching the schema) regardless of the Go field
// names. Numeric primary keys are emitted as decimal strings so they
// compose with the cursor format.
func agentField(name string) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (any, error) {
		a, ok := p.Source.(store.Agent)
		if !ok {
			return nil, fmt.Errorf("agentField: source is %T, not store.Agent", p.Source)
		}
		switch name {
		case "id":
			return strconv.FormatInt(a.ID, 10), nil
		case "kind":
			return a.Kind, nil
		case "agentId":
			return optString(a.AgentID), nil
		case "controller":
			return optString(a.Controller), nil
		case "creator":
			return optString(a.Creator), nil
		case "metadataUri":
			return optString(a.MetadataURI), nil
		case "capabilities":
			return optString(a.Capabilities), nil
		case "blockNumber":
			return int(a.BlockNumber), nil
		case "txHash":
			return strings.ToLower(a.TxHash), nil
		case "logIndex":
			return int(a.LogIndex), nil
		case "createdAt":
			return a.CreatedAt, nil
		case "updatedAt":
			return a.UpdatedAt, nil
		}
		return nil, fmt.Errorf("agentField: unknown field %q", name)
	}
}

func tradeField(name string) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (any, error) {
		t, ok := p.Source.(store.Trade)
		if !ok {
			return nil, fmt.Errorf("tradeField: source is %T, not store.Trade", p.Source)
		}
		switch name {
		case "id":
			return strconv.FormatInt(t.ID, 10), nil
		case "agentTokenPK":
			if t.AgentTokenPK == nil {
				return nil, nil
			}
			return int(*t.AgentTokenPK), nil
		case "side":
			return t.Side, nil
		case "trader":
			return strings.ToLower(t.Trader), nil
		case "curve":
			return strings.ToLower(t.Curve), nil
		case "quoteIn":
			return optString(t.QuoteIn), nil
		case "agentOut":
			return optString(t.AgentOut), nil
		case "agentIn":
			return optString(t.AgentIn), nil
		case "quoteOut":
			return optString(t.QuoteOut), nil
		case "fee":
			return t.Fee, nil
		case "blockNumber":
			return int(t.BlockNumber), nil
		case "txHash":
			return strings.ToLower(t.TxHash), nil
		case "logIndex":
			return int(t.LogIndex), nil
		case "createdAt":
			return t.CreatedAt, nil
		}
		return nil, fmt.Errorf("tradeField: unknown field %q", name)
	}
}

func jobField(name string) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (any, error) {
		j, ok := p.Source.(store.Job)
		if !ok {
			return nil, fmt.Errorf("jobField: source is %T, not store.Job", p.Source)
		}
		switch name {
		case "id":
			return strconv.FormatInt(j.ID, 10), nil
		case "onChainId":
			return j.OnChainID, nil
		case "cloneAddress":
			return optString(j.CloneAddress), nil
		case "agentPK":
			if j.AgentPK == nil {
				return nil, nil
			}
			return int(*j.AgentPK), nil
		case "agentId":
			return optString(j.AgentID), nil
		case "principal":
			return strings.ToLower(j.Principal), nil
		case "jobType":
			return int(j.JobType), nil
		case "token":
			return optString(j.Token), nil
		case "budget":
			return j.Budget, nil
		case "deadline":
			if j.Deadline == nil {
				return nil, nil
			}
			return *j.Deadline, nil
		case "phase":
			return j.Phase, nil
		case "resultUri":
			return optString(j.ResultURI), nil
		case "blockNumber":
			return int(j.BlockNumber), nil
		case "txHash":
			return strings.ToLower(j.TxHash), nil
		case "logIndex":
			return int(j.LogIndex), nil
		case "createdAt":
			return j.CreatedAt, nil
		case "updatedAt":
			return j.UpdatedAt, nil
		}
		return nil, fmt.Errorf("jobField: unknown field %q", name)
	}
}

// optString returns nil for the empty string so the GraphQL serialiser
// emits `null` rather than `""`. Avoids the foot-gun where a client
// has to differentiate "field absent" from "field set to empty string"
// in optional address / hash fields.
func optString(s string) any {
	if s == "" {
		return nil
	}
	return s
}
