package models

import "time"

// ConnectionStatus represents the status of a connection
type ConnectionStatus string

const (
	// ConnectionStatusComplete represents a complete connection with both request and response
	ConnectionStatusComplete ConnectionStatus = "complete"
	// ConnectionStatusRequestOnly represents a connection with only a request
	ConnectionStatusRequestOnly ConnectionStatus = "request_only"
	// ConnectionStatusResponseOnly represents a connection with only a response
	ConnectionStatusResponseOnly ConnectionStatus = "response_only"
)

type Connection struct {
	ID        string           `json:"id"`
	Protocol  string           `json:"protocol"`
	Source    string           `json:"source"`
	Target    string           `json:"target"`
	Timestamp int64            `json:"timestamp"`
	Request   interface{}      `json:"request"`
	Response  interface{}      `json:"response"`
	Status    ConnectionStatus `json:"status,omitempty"`
}

// HalfConnection represents an incomplete transaction
type HalfConnection struct {
	ConnectionID string      `json:"connectionId"`
	Type         string      `json:"type"` // "request" or "response"
	Timestamp    time.Time   `json:"timestamp"`
	Data         interface{} `json:"data"`
}
