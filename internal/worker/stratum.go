package worker

import "encoding/json"

// JSONRPCRequest represents a JSON-RPC 1.0 request to the pool.
type JSONRPCRequest struct {
	ID     int           `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

// JSONRPCResponse represents a generic JSON-RPC 1.0 response from the pool.
type JSONRPCResponse struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  interface{}     `json:"error"`
}

// JSONRPCNotification represents a notification from the pool (usually id = null).
type JSONRPCNotification struct {
	ID     *int            `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}


