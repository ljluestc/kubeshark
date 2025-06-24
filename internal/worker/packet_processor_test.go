package worker

import (
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/kubeshark/kubeshark/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockConnectionService implements the ConnectionService interface for testing
type MockConnectionService struct {
	mock.Mock
}

func (m *MockConnectionService) AddRequest(id string, protocol string, source string, target string, timestamp time.Time, request interface{}) bool {
	args := m.Called(id, protocol, source, target, timestamp, request)
	return args.Bool(0)
}

func (m *MockConnectionService) AddResponse(id string, response interface{}) bool {
	args := m.Called(id, response)
	return args.Bool(0)
}

func (m *MockConnectionService) TrackHalfConnection(id string, connectionType string) {
	m.Called(id, connectionType)
}

func (m *MockConnectionService) AddResponseOnly(id string, protocol string, source string, target string, timestamp time.Time, response interface{}) {
	m.Called(id, protocol, source, target, timestamp, response)
}

func (m *MockConnectionService) GetConnections(includeHalf bool) ([]*models.Connection, error) {
	args := m.Called(includeHalf)
	return args.Get(0).([]*models.Connection), args.Error(1)
}

// MockLogger implements the Logger interface for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

// Mock implementations for PacketProcessor methods
func mockPacket() gopacket.Packet {
	// This is a placeholder - in a real test we'd create a properly mocked packet
	return &MockPacket{}
}

// MockPacket implements gopacket.Packet for testing
type MockPacket struct{}

// Layer implementations for MockPacket
func (m *MockPacket) Layer(layerType gopacket.LayerType) gopacket.Layer           { return nil }
func (m *MockPacket) Layers() []gopacket.Layer                                    { return nil }
func (m *MockPacket) LayerClass(layerType gopacket.LayerType) gopacket.LayerClass { return nil }
func (m *MockPacket) NetworkLayer() gopacket.NetworkLayer                         { return nil }
func (m *MockPacket) TransportLayer() gopacket.TransportLayer                     { return nil }
func (m *MockPacket) ApplicationLayer() gopacket.ApplicationLayer                 { return nil }
func (m *MockPacket) ErrorLayer() gopacket.ErrorLayer                             { return nil }
func (m *MockPacket) LinkLayer() gopacket.LinkLayer                               { return nil }
func (m *MockPacket) Data() []byte                                                { return nil }
func (m *MockPacket) Metadata() *gopacket.PacketMetadata                          { return nil }
func (m *MockPacket) String() string                                              { return "" }
func (m *MockPacket) Dump() string                                                { return "" }
func (m *MockPacket) CgroupID() uint64                                            { return 0 }
func (m *MockPacket) Direction() gopacket.CaptureInfo                             { return gopacket.CaptureInfo{} }

// TestProcessRequestOnly tests request-only half-connection handling
func TestProcessRequestOnly(t *testing.T) {
	mockCS := new(MockConnectionService)
	mockL := new(MockLogger)

	processor := &PacketProcessor{
		connectionService: mockCS,
		logger:            mockL,
	}

	// Override methods for testing
	processor.generateConnectionID = func(p gopacket.Packet) string { return "test-id" }
	processor.determineIfRequest = func(p gopacket.Packet) bool { return true }
	processor.determineProtocol = func(p gopacket.Packet) string { return "HTTP" }
	processor.getSource = func(p gopacket.Packet) string { return "client" }
	processor.getTarget = func(p gopacket.Packet) string { return "server" }
	processor.getTimestamp = func(p gopacket.Packet) time.Time { return time.Unix(12345, 0) }
	processor.extractRequestData = func(p gopacket.Packet) interface{} { return "request-data" }

	// Setup expectations
	mockCS.On("AddRequest", "test-id", "HTTP", "client", "server", time.Unix(12345, 0), "request-data").Return(false)
	mockCS.On("TrackHalfConnection", "test-id", "request").Return()
	mockL.On("Debug", "Tracked request-only half-connection", mock.Anything).Return()

	// Call the method under test
	processor.processPacket(mockPacket())

	// Verify expectations
	mockCS.AssertExpectations(t)
	mockL.AssertExpectations(t)
}

// TestProcessResponseOnly tests response-only half-connection handling
func TestProcessResponseOnly(t *testing.T) {
	mockCS := new(MockConnectionService)
	mockL := new(MockLogger)

	processor := &PacketProcessor{
		connectionService: mockCS,
		logger:            mockL,
	}

	// Override methods for testing
	processor.generateConnectionID = func(p gopacket.Packet) string { return "test-id" }
	processor.determineIfRequest = func(p gopacket.Packet) bool { return false }
	processor.determineProtocol = func(p gopacket.Packet) string { return "HTTP" }
	processor.getSource = func(p gopacket.Packet) string { return "server" }
	processor.getTarget = func(p gopacket.Packet) string { return "client" }
	processor.getTimestamp = func(p gopacket.Packet) time.Time { return time.Unix(12345, 0) }
	processor.extractResponseData = func(p gopacket.Packet) interface{} { return "response-data" }

	// Setup expectations
	mockCS.On("AddResponse", "test-id", "response-data").Return(false)
	mockCS.On("AddResponseOnly", "test-id", "HTTP", "server", "client", time.Unix(12345, 0), "response-data").Return()
	mockCS.On("TrackHalfConnection", "test-id", "response").Return()
	mockL.On("Debug", "Tracked response-only half-connection", mock.Anything).Return()

	// Call the method under test
	processor.processPacket(mockPacket())

	// Verify expectations
	mockCS.AssertExpectations(t)
	mockL.AssertExpectations(t)
}

// TestProcessCompleteConnection tests complete connection handling
func TestProcessCompleteConnection(t *testing.T) {
	mockCS := new(MockConnectionService)
	mockL := new(MockLogger)

	processor := &PacketProcessor{
		connectionService: mockCS,
		logger:            mockL,
	}

	// Test request that completes a connection
	processor.generateConnectionID = func(p gopacket.Packet) string { return "test-id" }
	processor.determineIfRequest = func(p gopacket.Packet) bool { return true }
	processor.determineProtocol = func(p gopacket.Packet) string { return "HTTP" }
	processor.getSource = func(p gopacket.Packet) string { return "client" }
	processor.getTarget = func(p gopacket.Packet) string { return "server" }
	processor.getTimestamp = func(p gopacket.Packet) time.Time { return time.Unix(12345, 0) }
	processor.extractRequestData = func(p gopacket.Packet) interface{} { return "request-data" }

	// Setup expectations for completed connection
	mockCS.On("AddRequest", "test-id", "HTTP", "client", "server", time.Unix(12345, 0), "request-data").Return(true)

	// Call the method under test
	processor.processPacket(mockPacket())

	// Verify expectations
	mockCS.AssertExpectations(t)

	// Now test response that completes a connection
	mockCS = new(MockConnectionService)
	processor.connectionService = mockCS

	processor.determineIfRequest = func(p gopacket.Packet) bool { return false }
	processor.extractResponseData = func(p gopacket.Packet) interface{} { return "response-data" }

	// Setup expectations for completed connection
	mockCS.On("AddResponse", "test-id", "response-data").Return(true)

	// Call the method under test
	processor.processPacket(mockPacket())

	// Verify expectations
	mockCS.AssertExpectations(t)
}
