package scripting

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
)

func TestPcapHelperRetention(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "pcap-helper-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a PCAP manager for testing
	manager := worker.NewPcapManager(tempDir, 10*time.Second, 1024*1024)

	// Create a PCAP helper
	helper := NewPcapHelper(manager)

	// Test PCAP retention
	testPcapName := "test_stream_123456.pcap"
	helper.RetainPcap(testPcapName, 60) // Retain for 60 seconds

	// Verify the PCAP is retained
	if !helper.IsRetained(testPcapName) {
		t.Errorf("PCAP should be retained")
	}

	// Create a PCAP with short retention
	shortPcapName := "short_retention.pcap"
	helper.RetainPcap(shortPcapName, 1) // Retain for 1 second

	// Wait for retention to expire
	time.Sleep(2 * time.Second)

	// Verify short retention PCAP is no longer retained
	if helper.IsRetained(shortPcapName) {
		t.Errorf("Short retention PCAP should not be retained anymore")
	}
}

func TestGetPcapPath(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "pcap-path-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a PCAP manager for testing
	manager := worker.NewPcapManager(tempDir, 10*time.Second, 1024*1024)

	// Create a PCAP helper
	helper := NewPcapHelper(manager)

	// Test getting PCAP path
	streamID := "test_stream_789012"
	expectedPath := filepath.Join(tempDir, streamID+".pcap")

	actualPath := helper.GetPcapPath(streamID)

	if actualPath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, actualPath)
	}
}

func TestCleanupExpiredRetentions(t *testing.T) {
	// Create a PCAP manager for testing
	manager := worker.NewPcapManager("/tmp/pcaps", 10*time.Second, 1024*1024)

	// Create a PCAP helper
	helper := NewPcapHelper(manager)

	// Add some retentions
	helper.RetainPcap("test1.pcap", 60)
	helper.RetainPcap("test2.pcap", 1)

	// Wait for test2 to expire
	time.Sleep(2 * time.Second)

	// Clean up
	helper.CleanupExpiredRetentions()

	// Verify test1 is still retained, but test2 is not
	if !helper.IsRetained("test1.pcap") {
		t.Errorf("test1.pcap should still be retained")
	}

	if helper.IsRetained("test2.pcap") {
		t.Errorf("test2.pcap should not be retained anymore")
	}
}

func TestGetPcapPath(t *testing.T) {
	// Create a PCAP manager for testing
	manager := worker.NewPcapManager("/tmp/pcaps", 10*time.Second, 1024*1024)

	// Create a PCAP helper
	helper := NewPcapHelper(manager)

	// Test getting PCAP path
	streamID := "test_stream_789012"
	expectedPath := "/tmp/pcaps/test_stream_789012.pcap"

	actualPath := helper.GetPcapPath(streamID)

	if actualPath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, actualPath)
	}
}

func TestCleanupExpiredRetentions(t *testing.T) {
	// Create a PCAP manager for testing
	manager := worker.NewPcapManager("/tmp/pcaps", 10*time.Second, 1024*1024)

	// Create a PCAP helper
	helper := NewPcapHelper(manager)

	// Add some retentions
	helper.RetainPcap("test1.pcap", 60)
	helper.RetainPcap("test2.pcap", 1)

	// Wait for test2 to expire
	time.Sleep(2 * time.Second)

	// Clean up
	helper.CleanupExpiredRetentions()

	// Verify test1 is still retained, but test2 is not
	if !helper.IsRetained("test1.pcap") {
		t.Errorf("test1.pcap should still be retained")
	}

	if helper.IsRetained("test2.pcap") {
		t.Errorf("test2.pcap should not be retained anymore")
	}
}

func TestGetRetainedPcaps(t *testing.T) {
	// Create a PCAP manager for testing
	manager := worker.NewPcapManager("/tmp/pcaps", 10*time.Second, 1024*1024)

	// Create a PCAP helper
	helper := NewPcapHelper(manager)

	// Initially should be empty
	initialPcaps := helper.GetRetainedPcaps()
	if len(initialPcaps) > 0 {
		t.Errorf("Initially there should be no retained PCAPs")
	}

	// Retain some PCAPs
	helper.RetainPcap("test1.pcap", 60)
	helper.RetainPcap("test2.pcap", 60)

	// Get retained PCAPs
	retainedPcaps := helper.GetRetainedPcaps()

	// Should have 2 retained PCAPs
	if len(retainedPcaps) != 2 {
		t.Errorf("Expected 2 retained PCAPs, got %d", len(retainedPcaps))
	}
}

func TestGetRetainedPcaps(t *testing.T) {
	// Create a PCAP manager for testing
	manager := worker.NewPcapManager("/tmp/pcaps", 10*time.Second, 1024*1024)

	// Create a PCAP helper
	helper := NewPcapHelper(manager)

	// Initially should be empty
	initialPcaps := helper.GetRetainedPcaps()
	if len(initialPcaps) > 0 {
		t.Errorf("Initially there should be no retained PCAPs")
	}

	// Retain some PCAPs
	helper.RetainPcap("test1.pcap", 60)
	helper.RetainPcap("test2.pcap", 60)

	// Get retained PCAPs
	retainedPcaps := helper.GetRetainedPcaps()

	// Should have 2 retained PCAPs
	if len(retainedPcaps) != 2 {
		t.Errorf("Expected 2 retained PCAPs, got %d", len(retainedPcaps))
	}
}

// Ensure all logic is inside a function, for example:
func TestSomeFeature(t *testing.T) {
	// ...existing test code...

	// Move any stray 'if' statements here
	// if ... { ... }
}
