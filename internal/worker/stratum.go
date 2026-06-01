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

// SubscribeResult is the result of mining.subscribe
type SubscribeResult []interface{}

// NotifyParams represents the params array of a mining.notify request.
// ["job_id", "prevhash", "coinb1", "coinb2", [merkle_branches], "version", "nbits", "ntime", clean_jobs]
type NotifyParams []interface{}
