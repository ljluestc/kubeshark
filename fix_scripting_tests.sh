#!/bin/bash

# Script to fix the scripting package test files

SCRIPTING_DIR="internal/scripting"

echo "=== Fixing scripting package test files ==="

# Create backup of pcap_helpers_test.go
if [ -f "$SCRIPTING_DIR/pcap_helpers_test.go" ]; then
  cp "$SCRIPTING_DIR/pcap_helpers_test.go" "$SCRIPTING_DIR/pcap_helpers_test.go.bak"
  
  # Create a new file with proper structure
  cat > "$SCRIPTING_DIR/pcap_helpers_test.go" << 'EOF'
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
}

func TestPcapRetentionExpiration(t *testing.T) {
	// Create a PCAP manager for testing
	manager := worker.NewPcapManager("/tmp/pcaps", 10*time.Second, 1024*1024)
	
	// Create a PCAP helper
	helper := NewPcapHelper(manager)
	
	// Add some retentions
	helper.RetainPcap("test1.pcap", 60)
	helper.RetainPcap("test2.pcap", 1)
	
	// Wait for test2 to expire
	time.Sleep(2 * time.Second)
	
	// Verify test1 is still retained, but test2 is not
	if !helper.IsRetained("test1.pcap") {
		t.Errorf("test1.pcap should still be retained")
	}
	
	if helper.IsRetained("test2.pcap") {
		t.Errorf("test2.pcap should not be retained anymore")
	}
}
EOF

  echo "Fixed pcap_helpers_test.go"
  go fmt "./$SCRIPTING_DIR/pcap_helpers_test.go"
else
  echo "pcap_helpers_test.go not found. No changes made."
fi

# Create pcap_helpers.go if it doesn't exist
if [ ! -f "$SCRIPTING_DIR/pcap_helpers.go" ] || [ ! -s "$SCRIPTING_DIR/pcap_helpers.go" ]; then
  cp "$SCRIPTING_DIR/pcap_helpers_fix.go" "$SCRIPTING_DIR/pcap_helpers.go"
  # Remove build tag from the copied file
  sed -i '/\/\/ +build ignore/d' "$SCRIPTING_DIR/pcap_helpers.go"
  echo "Created/updated pcap_helpers.go"
fi

echo "=== Scripting package test fix complete ==="
