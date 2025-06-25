package worker

import (
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/kubeshark/kubeshark/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConnectionService implements services.ConnectionService for testing
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

// MockLogger implements Logger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	// When called with variadic arguments, we need to pass them as individual arguments to Called
	if len(args) == 0 {
		m.Called(msg)
	} else if len(args) == 1 {
		m.Called(msg, args[0])
	} else {
		// We need to handle the case where args are passed as a slice
		callArgs := make([]interface{}, 0, len(args)+1)
		callArgs = append(callArgs, msg)
		for _, arg := range args {
			callArgs = append(callArgs, arg)
		}
		m.Called(callArgs...)
	}
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	// When using testify/mock with variadic arguments,
	// we need to be careful about how the arguments are passed
	if len(args) == 0 {
		m.Called(msg)
		return
	}

	// When debug is called with direct args like Debug(msg, "key", "value")
	if len(args) == 2 && args[0] == "connectionID" {
		m.Called(msg, "connectionID", args[1])
		return
	}

	// Handle the case where args are passed as an array
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	// Mock arguments need to match exactly how they're passed in code
	callArgs := make([]interface{}, 0, len(args)+1)
	callArgs = append(callArgs, msg)
	for _, arg := range args {
		callArgs = append(callArgs, arg)
	}
	m.Called(callArgs...)
}

// MockPacket is a simple implementation of gopacket.Packet for testing
type MockPacket struct{}

func (p MockPacket) String() string                                   { return "MockPacket" }
func (p MockPacket) Dump() string                                     { return "MockPacket Dump" }
func (p MockPacket) Layers() []gopacket.Layer                         { return nil }
func (p MockPacket) Layer(lt gopacket.LayerType) gopacket.Layer       { return nil }
func (p MockPacket) LayerClass(lc gopacket.LayerClass) gopacket.Layer { return nil }
func (p MockPacket) NetworkLayer() gopacket.NetworkLayer              { return nil }
func (p MockPacket) TransportLayer() gopacket.TransportLayer          { return nil }
func (p MockPacket) ApplicationLayer() gopacket.ApplicationLayer      { return nil }
func (p MockPacket) ErrorLayer() gopacket.ErrorLayer                  { return nil }
func (p MockPacket) LinkLayer() gopacket.LinkLayer                    { return nil }
func (p MockPacket) Data() []byte                                     { return nil }
func (p MockPacket) Metadata() *gopacket.PacketMetadata               { return nil }
func (p MockPacket) CgroupID() uint64                                 { return 0 }
func (p MockPacket) Direction() gopacket.CaptureInfo                  { return gopacket.CaptureInfo{} }

// TestProcessPacketRequest tests processing a request packet
func TestProcessPacketRequest(t *testing.T) {
	mockCS := new(MockConnectionService)
	mockLogger := new(MockLogger)

	// We'll use MockPacket{} directly in the call
	timestamp := time.Unix(123456789, 0)

	processor := NewPacketProcessor(mockCS, mockLogger)

	// Set expectations
	mockCS.On("AddRequest", "test-conn-id", "http", "10.0.0.1", "10.0.0.2", timestamp, "request data").Return(true)

	// Override functions with test values (using accessor methods)
	setProcessorTestFunctions(processor, "test-conn-id", true, "http", "10.0.0.1", "10.0.0.2", timestamp, "request data", nil)

	// Call the method being tested
	processor.processPacket(MockPacket{})

	// Verify expectations
	mockCS.AssertExpectations(t)
	assert.NotNil(t, processor)
}

// TestProcessPacketResponseWithoutRequest tests processing a response packet without a matching request
func TestProcessPacketResponseWithoutRequest(t *testing.T) {
	mockCS := new(MockConnectionService)
	mockLogger := new(MockLogger)

	// Set the test timestamp
	timestamp := time.Unix(123456789, 0)

	processor := NewPacketProcessor(mockCS, mockLogger)

	// Set expectations
	mockCS.On("AddResponse", "test-conn-id", "response data").Return(false)
	mockCS.On("AddResponseOnly", "test-conn-id", "http", "10.0.0.2", "10.0.0.1", timestamp, "response data").Return()
	// Match the exact way Debug is called with individual arguments
	mockLogger.On("Debug", "Tracked response-only half-connection", mock.Anything).Return()

	// Override functions with test values
	setProcessorTestFunctions(processor, "test-conn-id", false, "http", "10.0.0.2", "10.0.0.1", timestamp, nil, "response data")

	// Call the method being tested
	processor.processPacket(&MockPacket{})

	// Verify expectations
	mockCS.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestProcessPacketResponseWithRequest tests processing a response packet with a matching request
func TestProcessPacketResponseWithRequest(t *testing.T) {
	mockCS := new(MockConnectionService)
	mockLogger := new(MockLogger)

	// Set the test timestamp
	timestamp := time.Unix(123456789, 0)

	processor := NewPacketProcessor(mockCS, mockLogger)

	// Set expectations
	mockCS.On("AddResponse", "test-conn-id", "response data").Return(true)

	// Override functions with test values
	setProcessorTestFunctions(processor, "test-conn-id", false, "http", "10.0.0.2", "10.0.0.1", timestamp, nil, "response data")

	// Call the method being tested
	processor.processPacket(MockPacket{})

	// Verify expectations
	mockCS.AssertExpectations(t)
	// No Debug calls should be made since AddResponse returned true
}

// Helper function to set test functions on processor
func setProcessorTestFunctions(
	p *PacketProcessor,
	connID string,
	isRequest bool,
	protocol string,
	source string,
	target string,
	timestamp time.Time,
	requestData interface{},
	responseData interface{},
) {
	// We need to set the exported fields via the setters
	p.generateConnectionID = func(packet gopacket.Packet) string { return connID }
	p.determineIfRequest = func(packet gopacket.Packet) bool { return isRequest }
	p.determineProtocol = func(packet gopacket.Packet) string { return protocol }
	p.getSource = func(packet gopacket.Packet) string { return source }
	p.getTarget = func(packet gopacket.Packet) string { return target }
	p.getTimestamp = func(packet gopacket.Packet) time.Time { return timestamp }
	p.extractRequestData = func(packet gopacket.Packet) interface{} { return requestData }
	p.extractResponseData = func(packet gopacket.Packet) interface{} { return responseData }
}
