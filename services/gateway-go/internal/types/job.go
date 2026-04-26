package types

import "time"

// JobPhase represents the lifecycle phase of an ACP job.
// Phases follow a strict forward progression: Created → Funded → Active →
// Completed | Disputed → Released | Resolved.
type JobPhase string

const (
	JobPhaseCreated   JobPhase = "created"
	JobPhaseFunded    JobPhase = "funded"
	JobPhaseActive    JobPhase = "active"
	JobPhaseCompleted JobPhase = "completed"
	JobPhaseDisputed  JobPhase = "disputed"
	JobPhaseReleased  JobPhase = "released"
	JobPhaseResolved  JobPhase = "resolved"
)

// Job mirrors the on-chain ACP Job state as tracked by the indexer.
type Job struct {
	// ID is the gateway-assigned stable identifier (e.g. "job-001").
	ID string `json:"id"`

	// OnChainID is the uint256 job id emitted by JobFactory.JobCreated.
	OnChainID string `json:"on_chain_id"`

	// AgentID is the registered agent this job is assigned to.
	AgentID string `json:"agent_id"`

	// Requester is the hex address of the wallet that created the job.
	Requester string `json:"requester"`

	// Phase is the current lifecycle phase.
	Phase JobPhase `json:"phase"`

	// MetadataURI is the IPFS / HTTPS URI pointing to the job payload.
	MetadataURI string `json:"metadata_uri"`

	// EscrowAddress is the deployed escrow contract for this job (may be empty
	// before Funded phase).
	EscrowAddress string `json:"escrow_address,omitempty"`

	// Amount is the escrowed token amount as a decimal string (wei).
	Amount string `json:"amount,omitempty"`

	// Token is the ERC-20 address used for payment (address(0) == native ETH).
	Token string `json:"token,omitempty"`

	// Deadline is the on-chain submission deadline.
	Deadline *time.Time `json:"deadline,omitempty"`

	// ResultURI is set once the agent submits a result.
	ResultURI string `json:"result_uri,omitempty"`

	// CreatedAt is when the JobCreated event was indexed.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the job record was last mutated.
	UpdatedAt time.Time `json:"updated_at"`
}

// PrepareJobRequest is the body for POST /v1/jobs/prepare.
// It returns unsigned transaction calldata for the client to sign and broadcast.
type PrepareJobRequest struct {
	// AgentID is the registered agent to assign the job to.
	AgentID string `json:"agent_id"`

	// MetadataURI points to the off-chain job payload.
	MetadataURI string `json:"metadata_uri"`

	// Amount is the escrow amount in wei as a decimal string.
	Amount string `json:"amount"`

	// Token is the ERC-20 contract address for payment; empty string means native ETH.
	Token string `json:"token,omitempty"`

	// DeadlineSeconds is the job deadline expressed as a Unix timestamp.
	DeadlineSeconds int64 `json:"deadline_seconds"`
}

// PrepareJobResponse is the response for POST /v1/jobs/prepare.
type PrepareJobResponse struct {
	// To is the contract address to send the transaction to.
	To string `json:"to"`

	// Calldata is the ABI-encoded function call (hex, 0x-prefixed).
	Calldata string `json:"calldata"`

	// Value is the ETH value to attach (in wei, as decimal string). "0" if token payment.
	Value string `json:"value"`
}

// ListJobsResponse is the paginated jobs response.
type ListJobsResponse struct {
	Jobs       []Job  `json:"jobs"`
	NextCursor string `json:"next_cursor,omitempty"`
}
