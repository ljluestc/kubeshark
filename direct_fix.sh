#!/bin/bash

# Direct fix script - applies immediate fixes to specific files
# This will brute force fix the most problematic files

echo "====== STARTING DIRECT FIX ======"

# Fix 1: Correct Go version in go.mod
echo "Fixing go.mod..."
sed -i 's/go 1.24.3/go 1.24/' go.mod

# Fix 2: Replace configConfig.go with correct version
echo "Fixing config/configStructs/configConfig.go..."
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
}
EOF

# Fix 3: Replace logsConfig.go with correct version
echo "Fixing config/configStructs/logsConfig.go..."
cat > config/configStructs/logsConfig.go << 'EOF'
package configStructs

// LogsConfig defines logging configuration
type LogsConfig struct {
	Console bool   `yaml:"console" default:"true"`
	File    string `yaml:"file" default:""`
}
EOF

# Fix 4: Replace scriptingConfig.go with correct version
echo "Fixing config/configStructs/scriptingConfig.go..."
cat > config/configStructs/scriptingConfig.go << 'EOF'
package configStructs

// ScriptingConfig defines settings for the scripting engine
type ScriptingConfig struct {
	Enabled       bool   `yaml:"enabled" default:"false"`
	TimeoutMs     int    `yaml:"timeoutMs" default:"5000"`
	DefaultScript string `yaml:"defaultScript" default:""`
}
EOF

# Fix 5: Replace tapConfig.go with correct version
echo "Fixing config/configStructs/tapConfig.go..."
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

# Fix 6: Replace configStruct.go with correct version
echo "Fixing config/configStructs/configStruct.go..."
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

# Fix 7: Replace pcap_retention_test.go with renamed test functions
echo "Fixing internal/worker/pcap_retention_test.go..."
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

# Fix 8: Create a proper pcap_helpers.go file
echo "Creating internal/scripting/pcap_helpers.go..."
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

# Fix 9: Create a proper pcap_helpers_test.go file
echo "Creating internal/scripting/pcap_helpers_test.go..."
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
	
	// Test PCAP path generation
	streamID := "stream_123456"
	path := helper.GetPcapPath(streamID)
	expectedPath := "/tmp/pcaps/stream_123456.pcap"
	if path != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, path)
	}
}
EOF

# Clean Go caches
echo "Cleaning Go caches..."
go clean -cache
go clean -testcache

# Format all Go files
echo "Formatting all Go files..."
go fmt ./...

echo "====== DIRECT FIX COMPLETE ======"
echo "Now try running: make test"
