#!/bin/bash

# Ultimate fix script to address all issues in the Kubeshark project
# This script directly rewrites problematic files with corrected content

echo "====================================="
echo "= ULTIMATE KUBESHARK PROJECT FIX   ="
echo "====================================="

# Fix 1: Correct Go version in go.mod
echo "Fixing go.mod..."
sed -i 's/go 1.24.3/go 1.24/' go.mod

# Fix 2: Directly create correct config/configStructs files
echo "Fixing config/configStructs package files..."
mkdir -p config/configStructs

# Define a function to create a file with correct package
create_config_file() {
  local file="$1"
  local content="$2"
  echo "Creating $file"
  echo "package configStructs" > "$file"
  echo "" >> "$file"
  echo "$content" >> "$file"
}

# Create all needed config files with correct content
create_config_file "config/configStructs/configStruct.go" 'import "time"

// ConfigStruct is the main configuration structure
type ConfigStruct struct {
	LogLevel     string           `yaml:"logLevel" default:"info"`
	Config       ConfigConfig     `yaml:"config"`
	Tap          TapConfig        `yaml:"tap"`
	Logs         LogsConfig       `yaml:"logs"`
	Scripting    ScriptingConfig  `yaml:"scripting"`
}'

create_config_file "config/configStructs/configConfig.go" '// ConfigConfig defines configuration-related settings
type ConfigConfig struct {
	Path string `yaml:"path" default:""`
}'

create_config_file "config/configStructs/logsConfig.go" '// LogsConfig defines logging configuration
type LogsConfig struct {
	Console bool   `yaml:"console" default:"true"`
	File    string `yaml:"file" default:""`
}'

create_config_file "config/configStructs/scriptingConfig.go" '// ScriptingConfig defines settings for the scripting engine
type ScriptingConfig struct {
	Enabled       bool   `yaml:"enabled" default:"false"`
	TimeoutMs     int    `yaml:"timeoutMs" default:"5000"`
	DefaultScript string `yaml:"defaultScript" default:""`
}'

create_config_file "config/configStructs/tapConfig.go" '// TapConfig defines the network tap configuration
type TapConfig struct {
	Debug   bool              `yaml:"debug" default:"false"`
	Insecure bool             `yaml:"insecure" default:"false"`
	Misc    TapMiscConfig     `yaml:"misc"`
}

// TapMiscConfig contains miscellaneous tap settings
type TapMiscConfig struct {
	PcapTTL      string `yaml:"pcapTTL" default:"5m"`
	PcapSizeLimit int    `yaml:"pcapSizeLimit" default:"104857600"`
}'

# Fix 3: Create proper internal/scripting/pcap_helpers.go
echo "Creating internal/scripting/pcap_helpers.go..."
mkdir -p internal/scripting
cat > internal/scripting/pcap_helpers.go << 'EOF'
package scripting

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
)

// PcapHelper provides functionality for managing PCAP files
type PcapHelper struct {
	manager    *worker.PcapManager
	retentions map[string]time.Time
	mutex      sync.RWMutex
}

// NewPcapHelper creates a new PCAP helper
func NewPcapHelper(manager *worker.PcapManager) *PcapHelper {
	return &PcapHelper{
		manager:    manager,
		retentions: make(map[string]time.Time),
		mutex:      sync.RWMutex{},
	}
}

// RetainPcap marks a PCAP for retention for the specified duration in seconds
func (p *PcapHelper) RetainPcap(pcapName string, durationSec int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	// Calculate the expiration time
	expiration := time.Now().Add(time.Duration(durationSec) * time.Second)
	
	// Store in our retention map
	p.retentions[pcapName] = expiration
}

// IsRetained checks if a PCAP is currently retained
func (p *PcapHelper) IsRetained(pcapName string) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	expiration, exists := p.retentions[pcapName]
	if !exists {
		return false
	}
	
	// Check if the retention has expired
	return time.Now().Before(expiration)
}

// GetPcapPath returns the full path to a PCAP file based on stream ID
func (p *PcapHelper) GetPcapPath(streamID string) string {
	filename := streamID + ".pcap"
	return filepath.Join(p.manager.GetPcapDir(), filename)
}

// CleanupExpiredRetentions removes expired retentions
func (p *PcapHelper) CleanupExpiredRetentions() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	now := time.Now()
	for pcapName, expiration := range p.retentions {
		if now.After(expiration) {
			delete(p.retentions, pcapName)
		}
	}
}

// GetRetainedPcaps returns a list of currently retained PCAP names
func (p *PcapHelper) GetRetainedPcaps() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	pcaps := make([]string, 0, len(p.retentions))
	now := time.Now()
	
	for pcapName, expiration := range p.retentions {
		if now.Before(expiration) {
			pcaps = append(pcaps, pcapName)
		}
	}
	
	return pcaps
}
EOF

# Fix 4: Create proper internal/scripting/pcap_helpers_test.go
echo "Creating internal/scripting/pcap_helpers_test.go..."
cat > internal/scripting/pcap_helpers_test.go << 'EOF'
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

# Fix 5: Fix worker package test files
echo "Fixing internal/worker test files..."

# First check if we need to fix the files
if [ -f "internal/worker/pcap_retention_test.go" ]; then
  # Rename the conflicting test functions
  sed -i 's/TestPcapRetention/TestPcapRetentionAlternative/g' internal/worker/pcap_retention_test.go
  sed -i 's/TestPcapRetentionExpiration/TestPcapRetentionExpirationAlternative/g' internal/worker/pcap_retention_test.go
  
  # Ensure package declaration is correct
  sed -i '1s/^.*$/package worker/' internal/worker/pcap_retention_test.go
  
  echo "Fixed internal/worker/pcap_retention_test.go"
fi

# Step 6: Format all Go code
echo "Formatting all Go code..."
go fmt ./...

# Step 7: Clean Go caches
echo "Cleaning Go caches..."
go clean -cache -modcache -testcache

# Step 8: Tidy modules
echo "Running go mod tidy..."
go mod tidy

echo ""
echo "====================================="
echo "= FIX COMPLETE                     ="
echo "====================================="
echo ""
echo "Now try running 'make test' again."
echo "If issues persist, run the tests with verbose output:"
echo "go test -v ./..."
