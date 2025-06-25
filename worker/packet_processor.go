package worker

import (
	"time"

	"github.com/google/gopacket"
	"github.com/kubeshark/kubeshark/internal/services"
)

// Logger interface for dependency injection and testing
type Logger interface {
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Info(msg string, args ...interface{})
}

// Function types for easier testing
type (
	GenerateConnectionIDFunc func(packet gopacket.Packet) string
	DetermineIfRequestFunc   func(packet gopacket.Packet) bool
	DetermineProtocolFunc    func(packet gopacket.Packet) string
	GetSourceFunc            func(packet gopacket.Packet) string
	GetTargetFunc            func(packet gopacket.Packet) string
	GetTimestampFunc         func(packet gopacket.Packet) time.Time
	ExtractRequestDataFunc   func(packet gopacket.Packet) interface{}
	ExtractResponseDataFunc  func(packet gopacket.Packet) interface{}
)

// PacketProcessor handles packet processing and analysis
type PacketProcessor struct {
	connectionService services.ConnectionService
	logger            Logger

	// Functions that can be overridden for testing
	generateConnectionID GenerateConnectionIDFunc
	determineIfRequest   DetermineIfRequestFunc
	determineProtocol    DetermineProtocolFunc
	getSource            GetSourceFunc
	getTarget            GetTargetFunc
	getTimestamp         GetTimestampFunc
	extractRequestData   ExtractRequestDataFunc
	extractResponseData  ExtractResponseDataFunc
}

// NewPacketProcessor creates a new packet processor
func NewPacketProcessor(cs services.ConnectionService, logger Logger) *PacketProcessor {
	return &PacketProcessor{
		connectionService:    cs,
		logger:               logger,
		generateConnectionID: generateConnectionID,
		determineIfRequest:   determineIfRequest,
		determineProtocol:    determineProtocol,
		getSource:            getSource,
		getTarget:            getTarget,
		getTimestamp:         getTimestamp,
		extractRequestData:   extractRequestData,
		extractResponseData:  extractResponseData,
	}
}

// Default implementations - in a real codebase these would be actual implementations
func generateConnectionID(packet gopacket.Packet) string     { return "" }
func determineIfRequest(packet gopacket.Packet) bool         { return false }
func determineProtocol(packet gopacket.Packet) string        { return "" }
func getSource(packet gopacket.Packet) string                { return "" }
func getTarget(packet gopacket.Packet) string                { return "" }
func getTimestamp(packet gopacket.Packet) time.Time          { return time.Now() }
func extractRequestData(packet gopacket.Packet) interface{}  { return nil }
func extractResponseData(packet gopacket.Packet) interface{} { return nil }

// ProcessPacket processes a network packet
func (p *PacketProcessor) ProcessPacket(packet gopacket.Packet) {
	// Extract connection info
	connectionID := p.generateConnectionID(packet)
	isRequest := p.determineIfRequest(packet)
	protocol := p.determineProtocol(packet)
	source := p.getSource(packet)
	target := p.getTarget(packet)
	timestamp := p.getTimestamp(packet)

	if isRequest {
		// Process request
		request := p.extractRequestData(packet)
		success := p.connectionService.AddRequest(connectionID, protocol, source, target, timestamp, request)
		if !success {
			// Request without matching response - mark as half-connection
			p.connectionService.TrackHalfConnection(connectionID, "request")
			p.logger.Debug("Tracked request-only half-connection", "connectionID", connectionID)
		}
	} else {
		// Process response
		response := p.extractResponseData(packet)
		success := p.connectionService.AddResponse(connectionID, response)
		if !success {
			// Response without matching request - create half-connection
			p.connectionService.AddResponseOnly(connectionID, protocol, source, target, timestamp, response)
			// Use map for key-value pairs to make it clearer
			p.logger.Debug("Tracked response-only half-connection", map[string]interface{}{"connectionID": connectionID})
		}
	}
}

// processPacket is the unexported version for backward compatibility
func (p *PacketProcessor) processPacket(packet gopacket.Packet) {
	p.ProcessPacket(packet)
}
