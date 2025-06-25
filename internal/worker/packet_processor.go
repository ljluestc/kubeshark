package worker

import (
	"time"

	"github.com/google/gopacket"
	"github.com/kubeshark/kubeshark/internal/services"
)

// Using the Logger interface defined in logger.go

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
			p.connectionService.TrackHalfConnection(connectionID, "response")
			p.connectionService.AddResponseOnly(connectionID, protocol, source, target, timestamp, response)
			p.logger.Debug("Tracked response-only half-connection", "connectionID", connectionID)
		}
	}
}

// processPacket is the unexported version for backward compatibility
func (p *PacketProcessor) processPacket(packet gopacket.Packet) {
	p.ProcessPacket(packet)
}
