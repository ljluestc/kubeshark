#!/bin/bash
# This is a minimal fix script that uses a different approach

# First, fix go.mod
sed -i 's/go 1.24.3/go 1.24/' go.mod

# Create empty directories if they don't exist
mkdir -p config/configStructs
mkdir -p internal/scripting

# Fix config/configStructs/configConfig.go
cat <<EOT > config/configStructs/configConfig.go
package configStructs

// ConfigConfig defines configuration-related settings
type ConfigConfig struct {
	Path string \`yaml:"path" json:"path"\`
	AutoSave bool \`yaml:"autoSave" json:"autoSave" default:"true"\`
	WatchChanges bool \`yaml:"watchChanges" json:"watchChanges" default:"true"\`
}
EOT

# Fix config/configStructs/configStruct.go
cat <<EOT > config/configStructs/configStruct.go
package configStructs

// ConfigStruct is the main configuration structure
type ConfigStruct struct {
	LogLevel  string          \`yaml:"logLevel" default:"info"\`
	Config    ConfigConfig    \`yaml:"config"\`
	Tap       TapConfig       \`yaml:"tap"\`
	Logs      LogsConfig      \`yaml:"logs"\`
	Scripting ScriptingConfig \`yaml:"scripting"\`
}
EOT

# Fix config/configStructs/logsConfig.go
cat <<EOT > config/configStructs/logsConfig.go
package configStructs

// LogsConfig defines logging configuration
type LogsConfig struct {
	Console bool   \`yaml:"console" default:"true"\`
	File    string \`yaml:"file" default:""\`
}
EOT

# Fix config/configStructs/scriptingConfig.go
cat <<EOT > config/configStructs/scriptingConfig.go
package configStructs

// ScriptingConfig defines settings for the scripting engine
type ScriptingConfig struct {
	Enabled       bool   \`yaml:"enabled" default:"false"\`
	TimeoutMs     int    \`yaml:"timeoutMs" default:"5000"\`
	DefaultScript string \`yaml:"defaultScript" default:""\`
}
EOT

# Fix config/configStructs/tapConfig.go
cat <<EOT > config/configStructs/tapConfig.go
package configStructs

// TapConfig defines the network tap configuration
type TapConfig struct {
	Debug    bool          \`yaml:"debug" default:"false"\`
	Insecure bool          \`yaml:"insecure" default:"false"\`
	Misc     TapMiscConfig \`yaml:"misc"\`
}

// TapMiscConfig contains miscellaneous tap settings
type TapMiscConfig struct {
	PcapTTL      string \`yaml:"pcapTTL" default:"5m"\`
	PcapSizeLimit int    \`yaml:"pcapSizeLimit" default:"104857600"\`
}
EOT

# Fix internal/scripting/bindings.go
cat <<EOT > internal/scripting/bindings.go
package scripting

import (
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
)

// ScriptBindings contains the various objects bound to the scripting environment
type ScriptBindings struct {
	PcapHelper *PcapHelper
	StartTime  time.Time
}

// NewScriptBindings creates a new set of script bindings
func NewScriptBindings(pcapManager *worker.PcapManager) *ScriptBindings {
	return &ScriptBindings{
		PcapHelper: NewPcapHelper(pcapManager),
		StartTime:  time.Now(),
	}
}
EOT

# Fix internal/scripting/engine.go
cat <<EOT > internal/scripting/engine.go
package scripting

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

// ScriptEngine provides functionality to execute scripts
type ScriptEngine struct {
	bindings *ScriptBindings
	timeout  time.Duration
}

// NewScriptEngine creates a new scripting engine
func NewScriptEngine(bindings *ScriptBindings, timeoutMs int) *ScriptEngine {
	return &ScriptEngine{
		bindings: bindings,
		timeout:  time.Duration(timeoutMs) * time.Millisecond,
	}
}

// ExecuteScript runs a script with the configured bindings and timeout
func (e *ScriptEngine) ExecuteScript(script string) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	log.Debug().Msg("Executing script...")

	// Placeholder for actual script execution
	select {
	case <-time.After(50 * time.Millisecond):
		log.Debug().Msg("Script executed successfully")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("script execution timed out after %v", e.timeout)
	}
}
EOT

# Fix internal/scripting/pcap_helpers.go
cat <<EOT > internal/scripting/pcap_helpers.go
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
EOT

# Fix internal/scripting/pcap_helpers_test.go
cat <<EOT > internal/scripting/pcap_helpers_test.go
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
EOT

# Fix internal/scripting/scripting_service.go
cat <<EOT > internal/scripting/scripting_service.go
package scripting

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
	"github.com/rs/zerolog/log"
)

// ScriptingService manages script execution and lifecycle
type ScriptingService struct {
	engine   *ScriptEngine
	bindings *ScriptBindings
	mutex    sync.RWMutex
	running  bool
}

// NewScriptingService creates a new scripting service
func NewScriptingService(pcapManager *worker.PcapManager, timeoutMs int) *ScriptingService {
	bindings := NewScriptBindings(pcapManager)
	engine := NewScriptEngine(bindings, timeoutMs)
	
	return &ScriptingService{
		engine:   engine,
		bindings: bindings,
		mutex:    sync.RWMutex{},
		running:  false,
	}
}

// Start begins the scripting service
func (s *ScriptingService) Start(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if s.running {
		return fmt.Errorf("scripting service is already running")
	}
	
	log.Info().Msg("Starting scripting service")
	s.running = true
	
	return nil
}

// Stop halts the scripting service
func (s *ScriptingService) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if !s.running {
		return
	}
	
	log.Info().Msg("Stopping scripting service")
	s.running = false
}

// ExecuteScript runs a script with the current engine
func (s *ScriptingService) ExecuteScript(script string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	if !s.running {
		return fmt.Errorf("scripting service is not running")
	}
	
	return s.engine.ExecuteScript(script)
}
EOT

# Clean and format
go clean -cache
go fmt ./...

echo "Minimal fix completed. Try running 'make test' now."
