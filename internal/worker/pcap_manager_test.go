package worker

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPcapRetention(t *testing.T) {
	// Create a temporary directory for test PCAPs
	tempDir, err := os.MkdirTemp("", "pcap-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test PCAP files
	pcap1 := filepath.Join(tempDir, "test1.pcap")
	pcap2 := filepath.Join(tempDir, "test2.pcap")

	if err := os.WriteFile(pcap1, []byte("test data 1"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	if err := os.WriteFile(pcap2, []byte("test data 2"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Set file modification times to the past
	pastTime := time.Now().Add(-30 * time.Second)
	if err := os.Chtimes(pcap1, pastTime, pastTime); err != nil {
		t.Fatalf("Failed to set file time: %v", err)
	}
	if err := os.Chtimes(pcap2, pastTime, pastTime); err != nil {
		t.Fatalf("Failed to set file time: %v", err)
	}

	// Create manager with 20s TTL
	manager := NewPcapManager(tempDir, 20*time.Second, 1024*1024)

	// Mark pcap1 for retention
	manager.RetainPcap("test1.pcap", 60*time.Second)

	// Cleanup should remove pcap2 but keep pcap1
	if err := manager.CleanupExpiredPcaps(); err != nil {
		t.Fatalf("CleanupExpiredPcaps failed: %v", err)
	}

	// Check if files exist as expected
	if _, err := os.Stat(pcap1); os.IsNotExist(err) {
		t.Fatalf("Retained PCAP was incorrectly removed")
	}
	if _, err := os.Stat(pcap2); !os.IsNotExist(err) {
		t.Fatalf("Expired PCAP was not removed")
	}
}

func TestPcapRetentionExpiration(t *testing.T) {
	// Create a temporary directory for test PCAPs
	tempDir, err := os.MkdirTemp("", "pcap-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test PCAP file
	pcap1 := filepath.Join(tempDir, "test1.pcap")

	if err := os.WriteFile(pcap1, []byte("test data 1"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Set file modification time to the past
	pastTime := time.Now().Add(-30 * time.Second)
	if err := os.Chtimes(pcap1, pastTime, pastTime); err != nil {
		t.Fatalf("Failed to set file time: %v", err)
	}

	// Create manager with 20s TTL
	manager := NewPcapManager(tempDir, 20*time.Second, 1024*1024)

	// Mark pcap1 for short retention
	manager.RetainPcap("test1.pcap", 1*time.Millisecond)

	// Wait for retention to expire
	time.Sleep(10 * time.Millisecond)

	// Cleanup should now remove the file since retention expired
	if err := manager.CleanupExpiredPcaps(); err != nil {
		t.Fatalf("CleanupExpiredPcaps failed: %v", err)
	}

	// Check if file was removed
	if _, err := os.Stat(pcap1); !os.IsNotExist(err) {
		t.Fatalf("PCAP with expired retention was not removed")
	}
}
