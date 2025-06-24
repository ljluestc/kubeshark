package services

import (
	"sync"
	"time"

	"github.com/kubeshark/kubeshark/internal/models"
)

type ConnectionStore struct {
	connections map[string]*models.Connection
	mu          sync.RWMutex
}

func NewConnectionStore() *ConnectionStore {
	return &ConnectionStore{
		connections: make(map[string]*models.Connection),
	}
}

func (s *ConnectionStore) AddRequest(id string, protocol string, source string, target string, timestamp time.Time, request interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if conn, exists := s.connections[id]; exists {
		conn.Request = request
		if conn.Response != nil {
			conn.Status = models.ConnectionStatusComplete
			return true
		} else {
			conn.Status = models.ConnectionStatusRequestOnly
			return false
		}
	} else {
		s.connections[id] = &models.Connection{
			ID:        id,
			Protocol:  protocol,
			Source:    source,
			Target:    target,
			Timestamp: timestamp.Unix(),
			Request:   request,
			Status:    models.ConnectionStatusRequestOnly,
		}
		return false
	}
}

func (s *ConnectionStore) AddResponse(id string, response interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if conn, exists := s.connections[id]; exists {
		conn.Response = response
		conn.Status = models.ConnectionStatusComplete
		return true
	}

	return false
}

// TrackHalfConnection tracks a half-connection by type
func (s *ConnectionStore) TrackHalfConnection(id string, connectionType string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if conn, exists := s.connections[id]; exists {
		// Update status based on connection type
		if connectionType == "request" {
			conn.Status = models.ConnectionStatusRequestOnly
		} else if connectionType == "response" {
			conn.Status = models.ConnectionStatusResponseOnly
		}
	}
}

// AddResponseOnly creates a response-only half connection
func (s *ConnectionStore) AddResponseOnly(id string, protocol string, source string, target string, timestamp time.Time, response interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.connections[id] = &models.Connection{
		ID:        id,
		Protocol:  protocol,
		Source:    source,
		Target:    target,
		Response:  response,
		Timestamp: timestamp.Unix(),
		Status:    models.ConnectionStatusResponseOnly,
	}
}

func (s *ConnectionStore) GetAllConnections() []*models.Connection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	connections := make([]*models.Connection, 0, len(s.connections))
	for _, conn := range s.connections {
		connections = append(connections, conn)
	}
	return connections
}
