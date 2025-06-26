#!/bin/bash

# FINAL FIX SCRIPT - Direct replacement of problematic files with correct versions
# Run this script to fix all issues at once

echo "========================================================"
echo "  KUBESHARK PROJECT FINAL FIX SCRIPT                    "
echo "========================================================"

# Fix 1: Update go.mod
echo "Fixing go.mod..."
sed -i 's/go 1.24.3/go 1.24/' go.mod

# Fix 2: Replace all config/configStructs files
echo "Replacing config/configStructs files..."

# Create configConfig.go
cat > config/configStructs/configConfig.go << 'EOF'
package configStructs

// ConfigConfig defines configuration-related settings
type ConfigConfig struct {
	// Path to the config file
	Path string `yaml:"path" json:"path"`
	
	// Whether to save changes automatically
	AutoSave bool `yaml:"autoSave" json:"autoSave" default:"true"`
	
	// Whether to watch for file changes
	WatchChanges bool `yaml:"watchChanges" json:"watchChanges" default:"true"`
	
	// Whether to regenerate the config
	Regenerate bool `yaml:"regenerate,omitempty" json:"regenerate,omitempty" default:"false" readonly:""`
}

const (
	RegenerateConfigName = "regenerate"
)
EOF

# Create configStruct.go
cat > config/configStructs/configStruct.go << 'EOF'
package configStructs

// ConfigStruct is the main configuration structure
type ConfigStruct struct {
	LogLevel  string          `yaml:"logLevel" default:"info"`
	Config    ConfigConfig    `yaml:"config"`
	Tap       TapConfig       `yaml:"tap"`
	Logs      LogsConfig      `yaml:"logs"`
	Scripting ScriptingConfig `yaml:"scripting"`
}
EOF

# Create logsConfig.go
cat > config/configStructs/logsConfig.go << 'EOF'
package configStructs

// LogsConfig defines logging configuration
type LogsConfig struct {
	Console bool   `yaml:"console" default:"true"`
	File    string `yaml:"file" default:""`
}
EOF

# Create scriptingConfig.go
cat > config/configStructs/scriptingConfig.go << 'EOF'
package configStructs

// ScriptingConfig defines settings for the scripting engine
type ScriptingConfig struct {
	Enabled       bool   `yaml:"enabled" default:"false"`
	TimeoutMs     int    `yaml:"timeoutMs" default:"5000"`
	DefaultScript string `yaml:"defaultScript" default:""`
}
EOF

# Create tapConfig.go
cat > config/configStructs/tapConfig.go << 'EOF'
package configStructs

// TapConfig defines the network tap configuration
type TapConfig struct {
	Debug    bool          `yaml:"debug" default:"false"`
	Insecure bool          `yaml:"insecure" default:"false"`
	Misc     TapMiscConfig `yaml:"misc"`
}

// TapMiscConfig contains miscellaneous tap settings
type TapMiscConfig struct {
	PcapTTL      string `yaml:"pcapTTL" default:"5m"`
	PcapSizeLimit int    `yaml:"pcapSizeLimit" default:"104857600"`
}
EOF

# Fix 3: Replace worker test file to avoid conflict
echo "Replacing internal/worker/pcap_retention_test.go..."
cat > internal/worker/pcap_retention_test.go << 'EOF'
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
EOF

# Fix 4: Create/replace scripting package files
echo "Creating/replacing scripting package files..."

# pcap_helpers.go
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
EOF

# pcap_helpers_test.go
cat > internal/scripting/pcap_helpers_test.go << 'EOF'
package scripting

import (
	"testing"
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
)

func TestPcapHelperRetention(t *testing.T) {
	// Create a PCAP manager for testing
	manager := worker.NewPcapManager("/tmp/pcaps", 10*time.Second, 1024*1024)

	// Create a PCAP helper
	helper := NewPcapHelper(manager)

	// Test PCAP retention
	testPcapName := "test_stream_123456.pcap"
	helper.RetainPcap(testPcapName, 60) // Retain for 60 seconds
	
	// Verify retention
	if !helper.IsRetained(testPcapName) {
		t.Errorf("PCAP should be retained")
	}
}
EOF

# Now clean and format the code
echo "Cleaning Go caches..."
go clean -cache -modcache -testcache

echo "Formatting Go code..."
go fmt ./...

echo "Running go mod tidy..."
go mod tidy

echo "========================================================"
echo "  FIX COMPLETE - Try 'make test' now                    "
echo "========================================================"
