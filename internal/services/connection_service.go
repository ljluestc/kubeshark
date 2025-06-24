package services

import (
	"time"

	"github.com/kubeshark/kubeshark/internal/models"
)

// ConnectionService defines the interface for connection services
type ConnectionService interface {
	// AddRequest adds a request to the connection store and returns true if it completed a connection
	AddRequest(id string, protocol string, source string, target string, timestamp time.Time, request interface{}) bool

	// AddResponse adds a response to an existing request and returns true if successful
	AddResponse(id string, response interface{}) bool

	// TrackHalfConnection explicitly marks a connection as a half-connection
	TrackHalfConnection(id string, connectionType string)

	// AddResponseOnly creates a new response-only half-connection
	AddResponseOnly(id string, protocol string, source string, target string, timestamp time.Time, response interface{})

	// GetConnections returns connections based on includeHalf parameter
	GetConnections(includeHalf bool) ([]*models.Connection, error)
}
