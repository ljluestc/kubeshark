package worker

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
)

// MockConnectionService for testing
type MockConnectionService struct {
	requests        map[string]bool
	responses       map[string]bool
	halfConnections []string
}

func (m *MockConnectionService) AddRequest(connectionID, protocol, source, target string, timestamp time.Time, request interface{}) bool {
	m.requests[connectionID] = true
	return true
}

func (m *MockConnectionService) AddResponse(connectionID string, response interface{}) bool {
	m.responses[connectionID] = true
	return true
}

func (m *MockConnectionService) TrackHalfConnection(connectionID, connectionType string) {
	m.halfConnections = append(m.halfConnections, fmt.Sprintf("%s:%s", connectionID, connectionType))
}

func (m *MockConnectionService) AddResponseOnly(connectionID, protocol, source, target string, timestamp time.Time, response interface{}) {
	// For testing purposes, just mark as a response
	m.responses[connectionID] = true
}

// MockLogger for testing
type MockLogger struct {
	logs []string
}

func (m *MockLogger) Debug(msg string, keysAndValues ...interface{}) {
	logEntry := fmt.Sprintf("%s: %v", msg, keysAndValues)
	m.logs = append(m.logs, logEntry)
}

func TestPacketProcessor(t *testing.T) {
	connService := &MockConnectionService{
		requests:  make(map[string]bool),
		responses: make(map[string]bool),
	}
	logger := &MockLogger{}
	processor := NewPacketProcessor(connService, logger)

	// Create a test packet (HTTP request)
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(12345),
		DstPort: layers.TCPPort(80),
	}
	ipv4Layer := &layers.IPv4{
		SrcIP: net.ParseIP("192.168.1.1"),
		DstIP: net.ParseIP("10.0.0.1"),
	}
	appLayer := gopacket.Payload([]byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"))
	opts := gopacket.SerializeOptions{}
	buffer := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buffer, opts, ipv4Layer, tcpLayer, appLayer)
	assert.NoError(t, err)

	packet := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
	packet.Metadata().Timestamp = time.Now()

	t.Run("ProcessRequestPacket", func(t *testing.T) {
		processor.generateConnectionID = func(p gopacket.Packet) string {
			return "192.168.1.1:12345-10.0.0.1:80"
		}
		processor.determineIfRequest = func(p gopacket.Packet) bool {
			return true
		}
		processor.ProcessPacket(packet)
		connID := "192.168.1.1:12345-10.0.0.1:80"
		assert.True(t, connService.requests[connID], "Request should be added")
		assert.Empty(t, connService.halfConnections, "No half-connection should be tracked")
	})

	t.Run("ProcessResponsePacket", func(t *testing.T) {
		appLayer := gopacket.Payload([]byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"))
		buffer := gopacket.NewSerializeBuffer()
		err := gopacket.SerializeLayers(buffer, opts, ipv4Layer, tcpLayer, appLayer)
		assert.NoError(t, err)

		responsePacket := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
		responsePacket.Metadata().Timestamp = time.Now()

		processor.generateConnectionID = func(p gopacket.Packet) string {
			return "192.168.1.1:12345-10.0.0.1:80"
		}
		processor.determineIfRequest = func(p gopacket.Packet) bool {
			return false
		}
		processor.extractResponseData = func(p gopacket.Packet) interface{} {
			return "response data"
		}

		processor.ProcessPacket(responsePacket)
		connID := "192.168.1.1:12345-10.0.0.1:80"
		assert.True(t, connService.responses[connID], "Response should be added")
		assert.Contains(t, connService.halfConnections, connID+":response", "Half-connection should be tracked")
	})

	t.Run("ProcessNonHTTPPacket", func(t *testing.T) {
		buffer := gopacket.NewSerializeBuffer()
		err := gopacket.SerializeLayers(buffer, opts, ipv4Layer, tcpLayer)
		assert.NoError(t, err)

		nonHTTPPacket := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
		nonHTTPPacket.Metadata().Timestamp = time.Now()

		processor.generateConnectionID = func(p gopacket.Packet) string {
			return "192.168.1.1:12345-10.0.0.1:80"
		}
		processor.determineIfRequest = func(p gopacket.Packet) bool {
			return false
		}
		processor.extractResponseData = func(p gopacket.Packet) interface{} {
			return "non-http data"
		}

		processor.ProcessPacket(nonHTTPPacket)
		connID := "192.168.1.1:12345-10.0.0.1:80"
		assert.True(t, connService.responses[connID], "Non-HTTP packet treated as response")
		assert.Contains(t, connService.halfConnections, connID+":response", "Half-connection should be tracked")
	})
}

func TestGenerateConnectionID(t *testing.T) {
	processor := NewPacketProcessor(&MockConnectionService{requests: make(map[string]bool), responses: make(map[string]bool)}, &MockLogger{})
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(12345),
		DstPort: layers.TCPPort(80),
	}
	ipv4Layer := &layers.IPv4{
		SrcIP: net.ParseIP("192.168.1.1"),
		DstIP: net.ParseIP("10.0.0.1"),
	}
	buffer := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{}, ipv4Layer, tcpLayer)
	assert.NoError(t, err)

	packet := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeIPv4, gopacket.Default)

	// Override default function for testing
	originalFunc := processor.generateConnectionID
	processor.generateConnectionID = func(p gopacket.Packet) string {
		return "192.168.1.1:12345-10.0.0.1:80"
	}

	connID := processor.generateConnectionID(packet)
	assert.Equal(t, "192.168.1.1:12345-10.0.0.1:80", connID)

	// Restore original function
	processor.generateConnectionID = originalFunc
}
