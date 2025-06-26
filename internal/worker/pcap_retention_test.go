package worker

import (
	"testing"
	"time"
)

// Renamed test functions to avoid conflict with pcap_manager_test.go
func TestPcapRetentionObject(t *testing.T) {
	// Create a new PCAP retention manager with 20s default TTL
	retention := NewPcapRetention(20 * time.Second)

	// Test retaining a PCAP
	pcapPath := "pcaps/master/000000000123_udp.pcap"
	retention.RetainPcap(pcapPath, 60*time.Second)

	// Check if the PCAP should be retained
	if !retention.ShouldRetain(pcapPath) {
		t.Errorf("PCAP should be retained but wasn't")
	}

	// Test non-retained PCAP
	nonRetainedPath := "pcaps/master/000000000456_tcp.pcap"
	if retention.ShouldRetain(nonRetainedPath) {
		t.Errorf("PCAP should not be retained but was")
	}
}

func TestPcapRetentionObjectExpiration(t *testing.T) {
	// Create a new PCAP retention manager with 20s default TTL
	retention := NewPcapRetention(20 * time.Second)

	// Test retaining a PCAP with a very short TTL
	pcapPath := "pcaps/master/000000000123_udp.pcap"
	retention.RetainPcap(pcapPath, 1*time.Millisecond)

	// Wait for the retention to expire
	time.Sleep(10 * time.Millisecond)

	// PCAP should no longer be retained
	if retention.ShouldRetain(pcapPath) {
		t.Errorf("PCAP should not be retained after expiration but was")
	}
}

func TestPcapRetentionCleanup(t *testing.T) {
	// Create a new PCAP retention manager with 20s default TTL
	retention := NewPcapRetention(20 * time.Second)

	// Add several PCAPs with different expiration times
	retention.RetainPcap("pcaps/master/000000000123_udp.pcap", 1*time.Millisecond)
	retention.RetainPcap("pcaps/master/000000000456_tcp.pcap", 1*time.Hour)

	// Wait for the first one to expire
	time.Sleep(10 * time.Millisecond)

	// Run cleanup
	retention.CleanupExpired()

	// Check retention status
	if retention.ShouldRetain("pcaps/master/000000000123_udp.pcap") {
		t.Errorf("Expired PCAP should be removed after cleanup")
	}
	if !retention.ShouldRetain("pcaps/master/000000000456_tcp.pcap") {
		t.Errorf("Non-expired PCAP should be retained after cleanup")
	}
}
