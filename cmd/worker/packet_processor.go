package worker

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// ConnectionService defines the interface for managing connection data
type ConnectionService interface {
	AddRequest(connectionID, protocol, source, target string, timestamp time.Time, request interface{}) bool
	AddResponse(connectionID string, response interface{}) bool
	TrackHalfConnection(connectionID, connectionType string, timestamp time.Time, data interface{})
	AddResponseOnly(connectionID, protocol, source, target string, timestamp time.Time, response interface{})
	GetHalfConnections() []HalfConnection
}

// HalfConnection represents an incomplete transaction with either request or response missing
type HalfConnection struct {
	ConnectionID string
	Type         string // "request" or "response"
	Timestamp    time.Time
	Data         interface{}
}

// Logger defines the interface for logging
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// PacketProcessor handles packet processing and analysis
type PacketProcessor struct {
	connectionService ConnectionService
	logger            Logger

	// Exported fields for testing
	GenerateConnectionID func(packet gopacket.Packet) string
	DetermineIfRequest   func(packet gopacket.Packet) bool
	DetermineProtocol    func(packet gopacket.Packet) string
	GetSource            func(packet gopacket.Packet) string
	GetTarget            func(packet gopacket.Packet) string
	GetTimestamp         func(packet gopacket.Packet) time.Time
	ExtractRequestData   func(packet gopacket.Packet) interface{}
	ExtractResponseData  func(packet gopacket.Packet) interface{}
}

// NewPacketProcessor creates a new PacketProcessor
func NewPacketProcessor(connectionService ConnectionService, logger Logger) *PacketProcessor {
	return &PacketProcessor{
		connectionService: connectionService,
		logger:            logger,
	}
}

// ProcessPacket processes a single network packet
func (p *PacketProcessor) ProcessPacket(packet gopacket.Packet) {
	// Extract connection info
	var connectionID string
	var isRequest bool
	var protocol string
	var source string
	var target string
	var timestamp time.Time

	// Use exported fields if set, otherwise use the internal methods
	if p.GenerateConnectionID != nil {
		connectionID = p.GenerateConnectionID(packet)
	} else {
		connectionID = p.generateConnectionID(packet)
	}

	if connectionID == "" {
		return
	}

	if p.DetermineIfRequest != nil {
		isRequest = p.DetermineIfRequest(packet)
	} else {
		isRequest = p.determineIfRequest(packet)
	}

	if p.DetermineProtocol != nil {
		protocol = p.DetermineProtocol(packet)
	} else {
		protocol = p.determineProtocol(packet)
	}

	if p.GetSource != nil {
		source = p.GetSource(packet)
	} else {
		source = p.getSource(packet)
	}

	if p.GetTarget != nil {
		target = p.GetTarget(packet)
	} else {
		target = p.getTarget(packet)
	}

	if p.GetTimestamp != nil {
		timestamp = p.GetTimestamp(packet)
	} else {
		timestamp = p.getTimestamp(packet)
	}

	if isRequest {
		// Process request
		var request interface{}
		if p.ExtractRequestData != nil {
			request = p.ExtractRequestData(packet)
		} else {
			request = p.extractRequestData(packet)
		}

		added := p.connectionService.AddRequest(
			connectionID,
			protocol,
			source,
			target,
			timestamp,
			request,
		)
		if !added {
			// Request without matching response - track as half-connection
			p.connectionService.TrackHalfConnection(connectionID, "request", timestamp, request)
			p.logger.Debug("Tracked request-only half-connection", "connectionID", connectionID, "type", "request")
		}
	} else {
		// Process response
		var response interface{}
		if p.ExtractResponseData != nil {
			response = p.ExtractResponseData(packet)
		} else {
			response = p.extractResponseData(packet)
		}

		if response == nil {
			p.logger.Debug("Skipping empty response", "connectionID", connectionID)
			return
		}

		success := p.connectionService.AddResponse(connectionID, response)
		if !success {
			// Response without matching request - track as half-connection
			p.connectionService.TrackHalfConnection(connectionID, "response", timestamp, response)
			// Also create a response-only connection entry
			p.connectionService.AddResponseOnly(connectionID, protocol, source, target, timestamp, response)
			p.logger.Debug("Tracked response-only half-connection", "connectionID", connectionID, "type", "response")
		}
	}
}

// processPacket is an alias for ProcessPacket for backward compatibility
func (p *PacketProcessor) processPacket(packet gopacket.Packet) {
	p.ProcessPacket(packet)
}

// generateConnectionID creates a unique ID for the connection
func (p *PacketProcessor) generateConnectionID(packet gopacket.Packet) string {
	networkLayer := packet.NetworkLayer()
	transportLayer := packet.TransportLayer()
	if networkLayer == nil || transportLayer == nil {
		return ""
	}

	var srcIP, dstIP string
	switch nl := networkLayer.(type) {
	case *layers.IPv4:
		srcIP = nl.SrcIP.String()
		dstIP = nl.DstIP.String()
	case *layers.IPv6:
		srcIP = nl.SrcIP.String()
		dstIP = nl.DstIP.String()
	}

	var srcPort, dstPort string
	switch tl := transportLayer.(type) {
	case *layers.TCP:
		srcPort = fmt.Sprintf("%d", tl.SrcPort)
		dstPort = fmt.Sprintf("%d", tl.DstPort)
	case *layers.UDP:
		srcPort = fmt.Sprintf("%d", tl.SrcPort)
		dstPort = fmt.Sprintf("%d", tl.DstPort)
	}

	return fmt.Sprintf("%s:%s-%s:%s", srcIP, srcPort, dstIP, dstPort)
}

// determineIfRequest determines if the packet is a request
func (p *PacketProcessor) determineIfRequest(packet gopacket.Packet) bool {
	appLayer := packet.ApplicationLayer()
	if appLayer == nil {
		return false
	}

	payload := string(appLayer.Payload())
	return containsHTTPMethod(payload)
}

// extractRequestData extracts request data from the packet
func (p *PacketProcessor) extractRequestData(packet gopacket.Packet) interface{} {
	appLayer := packet.ApplicationLayer()
	if appLayer == nil {
		return nil
	}
	return string(appLayer.Payload())
}

// extractResponseData extracts response data from the packet
func (p *PacketProcessor) extractResponseData(packet gopacket.Packet) interface{} {
	appLayer := packet.ApplicationLayer()
	if appLayer == nil {
		return nil
	}
	return string(appLayer.Payload())
}

// determineProtocol determines the protocol of the packet
func (p *PacketProcessor) determineProtocol(packet gopacket.Packet) string {
	transportLayer := packet.TransportLayer()
	if transportLayer == nil {
		return "unknown"
	}
	switch transportLayer.(type) {
	case *layers.TCP:
		if appLayer := packet.ApplicationLayer(); appLayer != nil {
			if containsHTTPMethod(string(appLayer.Payload())) || isHTTPResponse(string(appLayer.Payload())) {
				return "http"
			}
		}
		return "tcp"
	case *layers.UDP:
		return "udp"
	default:
		return "unknown"
	}
}

// getSource gets the source address of the packet
func (p *PacketProcessor) getSource(packet gopacket.Packet) string {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return ""
	}
	switch nl := networkLayer.(type) {
	case *layers.IPv4:
		return nl.SrcIP.String()
	case *layers.IPv6:
		return nl.SrcIP.String()
	}
	return ""
}

// getTarget gets the target address of the packet
func (p *PacketProcessor) getTarget(packet gopacket.Packet) string {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return ""
	}
	switch nl := networkLayer.(type) {
	case *layers.IPv4:
		return nl.DstIP.String()
	case *layers.IPv6:
		return nl.DstIP.String()
	}
	return ""
}

// getTimestamp gets the timestamp of the packet
func (p *PacketProcessor) getTimestamp(packet gopacket.Packet) time.Time {
	if metadata := packet.Metadata(); metadata != nil {
		return metadata.Timestamp
	}
	return time.Now()
}

// Helper function to extract source and target endpoints from a packet
func extractEndpoints(packet gopacket.Packet) (string, string) {
	ipLayer := packet.NetworkLayer()
	tcpLayer := packet.TransportLayer()

	if ipLayer == nil || tcpLayer == nil {
		return "", ""
	}

	ip, ok := ipLayer.(*layers.IPv4)
	if !ok {
		return "", ""
	}

	tcp, ok := tcpLayer.(*layers.TCP)
	if !ok {
		return "", ""
	}

	source := fmt.Sprintf("%s:%d", ip.SrcIP, tcp.SrcPort)
	target := fmt.Sprintf("%s:%d", ip.DstIP, tcp.DstPort)

	return source, target
}

// containsHTTPMethod checks if the payload contains an HTTP method
func containsHTTPMethod(payload string) bool {
	methods := []string{"GET ", "POST ", "PUT ", "DELETE ", "HEAD ", "OPTIONS ", "PATCH "}
	for _, method := range methods {
		if len(payload) >= len(method) && payload[:len(method)] == method {
			return true
		}
	}
	return false
}

// isHTTPResponse checks if the payload is an HTTP response
func isHTTPResponse(payload string) bool {
	return len(payload) >= 8 && payload[:8] == "HTTP/1."
}
