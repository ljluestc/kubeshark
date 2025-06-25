package worker

import (
	"sync"
	"time"

	"github.com/kubeshark/kubeshark/internal/models"
)

// InMemoryConnectionService implements ConnectionService with in-memory storage
type InMemoryConnectionService struct {
	connections      map[string]*models.Connection
	halfConnections  map[string]string // connectionID -> type (request/response)
	mutex            sync.RWMutex
}

// NewInMemoryConnectionService creates a new InMemoryConnectionService
func NewInMemoryConnectionService() *InMemoryConnectionService {
	return &InMemoryConnectionService{
		connections:     make(map[string]*models.Connection),
		halfConnections: make(map[string]string),
	}
}

// AddRequest adds a new request to the connection service
func (s *InMemoryConnectionService) AddRequest(connectionID, protocol, source, target string, timestamp time.Time, request interface{}) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	conn, exists := s.connections[connectionID]
	if exists {
		conn.Request = request
		// If both request and response are present, mark as complete
		if conn.Response != nil {
			conn.Status = "complete"
			delete(s.halfConnections, connectionID)
			return true
		}
		return false
	}

	s.connections[connectionID] = &models.Connection{
		ID:        connectionID,
		Protocol:  protocol,
		Source:    source,
		Target:    target,
		Timestamp: timestamp,
		Request:   request,
		Status:    "request", // Mark as request-only
	}
	s.halfConnections[connectionID] = "request"
	return false
}

// AddResponse adds a response to an existing request
func (s *InMemoryConnectionService) AddResponse(connectionID string, response interface{}) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	conn, exists := s.connections[connectionID]
	if !exists {
		return false
	}

	conn.Response = response
	// If both request and response are present, mark as complete
	if conn.Request != nil {
		conn.Status = "complete"
		delete(s.halfConnections, connectionID)
		return true
	}
	return false
}

// TrackHalfConnection tracks a connection that only has a request or response
func (s *InMemoryConnectionService) TrackHalfConnection(connectionID, connectionType string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.halfConnections[connectionID] = connectionType
	
	// If the connection already exists, update its status
	if conn, exists := s.connections[connectionID]; exists {
		conn.Status = connectionType
	}
}

// AddResponseOnly adds a response without a matching request
func (s *InMemoryConnectionService) AddResponseOnly(connectionID, protocol, source, target string, timestamp time.Time, response interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.connections[connectionID] = &models.Connection{
		ID:        connectionID,
		Protocol:  protocol,
		Source:    source,
		Target:    target,
		Timestamp: timestamp,
		Response:  response,
		Status:    "response", // Mark as response-only
	}
	s.halfConnections[connectionID] = "response"
}

// GetConnections returns all tracked connections
// If includeHalf is true, half-connections are included with appropriate status
func (s *InMemoryConnectionService) GetConnections(includeHalf bool) ([]*models.Connection, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !includeHalf {
		// Only return complete connections
		completeConnections := make([]*models.Connection, 0)
		for _, conn := range s.connections {
			if conn.Status == "complete" || conn.Status == "" {
				if conn.Request != nil && conn.Response != nil {
					completeConnections = append(completeConnections, conn)
				}
			}
		}
		return completeConnections, nil
	}

	// Return all connections including half-connections
	connections := make([]*models.Connection, 0, len(s.connections))
	for _, conn := range s.connections {
		connections = append(connections, conn)
	}
	return connections, nil
}
