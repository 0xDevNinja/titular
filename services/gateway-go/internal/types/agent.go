package types

import "time"

// Curve holds bonding-curve reserve state for an agent.
type Curve struct {
	VirtualTITU  string `json:"virtual_titu"`
	VirtualAgent string `json:"virtual_agent"`
	RealTITU     string `json:"real_titu"`
	RealAgent    string `json:"real_agent"`
	Graduated    bool   `json:"graduated"`
}

// Agent mirrors the on-chain AgentLaunched event plus derived off-chain fields.
type Agent struct {
	ID           string    `json:"id"`
	Slug         string    `json:"slug"`
	Name         string    `json:"name"`
	Symbol       string    `json:"symbol"`
	ImageURI     string    `json:"image_uri"`
	Creator      string    `json:"creator"`       // hex address
	ModuleBitmap uint16    `json:"module_bitmap"` // bit flags for active modules
	CreatedAt    time.Time `json:"created_at"`
	Curve        Curve     `json:"curve"`
	TokenAddress string    `json:"token_address"`
	CurveAddress string    `json:"curve_address"`
}

// Trade mirrors the on-chain Bought/Sold events.
type Trade struct {
	TxHash      string    `json:"tx_hash"`
	BlockNumber uint64    `json:"block_number"`
	LogIndex    uint      `json:"log_index"`
	AgentID     string    `json:"agent_id"`
	Trader      string    `json:"trader"`       // hex address
	Side        string    `json:"side"`         // "buy" | "sell"
	TITUAmount  string    `json:"titu_amount"`  // decimal string
	AgentAmount string    `json:"agent_amount"` // decimal string
	PriceAfter  string    `json:"price_after"`  // decimal string
	Timestamp   time.Time `json:"timestamp"`
}

// ErrorCode is a machine-readable error identifier.
type ErrorCode string

const (
	ErrNotFound       ErrorCode = "not_found"
	ErrBadRequest     ErrorCode = "bad_request"
	ErrInternalServer ErrorCode = "internal_server_error"
)

// APIError is the inner error object returned in the envelope.
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// ErrorEnvelope wraps API errors in a consistent shape.
type ErrorEnvelope struct {
	Error APIError `json:"error"`
}

// ListAgentsResponse is the paginated agents response.
type ListAgentsResponse struct {
	Agents     []Agent `json:"agents"`
	NextCursor string  `json:"next_cursor,omitempty"`
}

// ListTradesResponse is the paginated trades response.
type ListTradesResponse struct {
	Trades     []Trade `json:"trades"`
	NextCursor string  `json:"next_cursor,omitempty"`
}
