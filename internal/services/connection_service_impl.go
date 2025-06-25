package services

import (
	"time"

	"github.com/kubeshark/kubeshark/internal/models"
)

// ConnectionServiceImpl implements the ConnectionService interface
type ConnectionServiceImpl struct {
	store *ConnectionStore
}

// NewConnectionServiceImpl creates a new connection service implementation
func NewConnectionServiceImpl(store *ConnectionStore) *ConnectionServiceImpl {
	return &ConnectionServiceImpl{
		store: store,
	}
}

// AddRequest adds a request to the connection store
func (s *ConnectionServiceImpl) AddRequest(id string, protocol string, source string, target string, timestamp time.Time, request interface{}) bool {
	return s.store.AddRequest(id, protocol, source, target, timestamp, request)
}

// AddResponse adds a response to the connection store
func (s *ConnectionServiceImpl) AddResponse(id string, response interface{}) bool {
	return s.store.AddResponse(id, response)
}

// TrackHalfConnection tracks half-connection status
func (s *ConnectionServiceImpl) TrackHalfConnection(id string, connectionType string) {
	s.store.TrackHalfConnection(id, connectionType)
}

// AddResponseOnly creates a new response-only half-connection
func (s *ConnectionServiceImpl) AddResponseOnly(id string, protocol string, source string, target string, timestamp time.Time, response interface{}) {
	s.store.AddResponseOnly(id, protocol, source, target, timestamp, response)
}

// GetConnections returns connections based on the includeHalf parameter
func (s *ConnectionServiceImpl) GetConnections(includeHalf bool) ([]*models.Connection, error) {
	connections := s.store.GetAllConnections()

	if !includeHalf {
		// Filter out unpaired connections
		pairedConnections := make([]*models.Connection, 0)
		for _, conn := range connections {
			if conn.Request != nil && conn.Response != nil {
				conn.Status = models.ConnectionStatusComplete
				pairedConnections = append(pairedConnections, conn)
			}
		}
		return pairedConnections, nil
	}

	// Include all connections and ensure status is set correctly
	for _, conn := range connections {
		if conn.Request != nil && conn.Response != nil {
			conn.Status = models.ConnectionStatusComplete
		} else if conn.Request != nil {
			conn.Status = models.ConnectionStatusRequestOnly
		} else if conn.Response != nil {
			conn.Status = models.ConnectionStatusResponseOnly
		}
	}

	return connections, nil
}
