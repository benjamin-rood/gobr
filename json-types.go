package gobr

import "encoding/json"

// InMsg – typestring wrapper for generic *received* msg
type InMsg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// OutMsg - typestring wrapper for generic *exported* msg
type OutMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
