package scripting
package scripting
package scripting

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
)

// PcapHelper provides JavaScript bindings for PCAP management
type PcapHelper struct {
	manager        *worker.PcapManager
	retentionMap   map[string]time.Time
	retentionMutex sync.RWMutex
}

// NewPcapHelper creates a new PCAP helper instance
func NewPcapHelper(manager *worker.PcapManager) *PcapHelper {
	return &PcapHelper{
		manager:        manager,
		retentionMap:   make(map[string]time.Time),
		retentionMutex: sync.RWMutex{},
	}
}

// RetainPcap marks a PCAP file for retention for the specified duration in seconds
func (p *PcapHelper) RetainPcap(pcapName string, durationSec int) {
	p.retentionMutex.Lock()
	defer p.retentionMutex.Unlock()

	// Calculate expiration time
	expiration := time.Now().Add(time.Duration(durationSec) * time.Second)
	
	// Store in retention map
	p.retentionMap[pcapName] = expiration
	
	// Also use the manager's retention if available
	p.manager.RetainPcap(pcapName, time.Duration(durationSec)*time.Second)
}

// IsRetained checks if a PCAP file is currently being retained
func (p *PcapHelper) IsRetained(pcapName string) bool {
	p.retentionMutex.RLock()
	defer p.retentionMutex.RUnlock()
	
	expiration, exists := p.retentionMap[pcapName]
	if !exists {
		return false
	}
	
	// Check if retention has expired
	return time.Now().Before(expiration)
}

// GetPcapPath returns the full path to a PCAP file given a stream ID
func (p *PcapHelper) GetPcapPath(streamID string) string {
	// Format the PCAP filename based on stream ID
	pcapName := fmt.Sprintf("test_stream_%s.pcap", streamID)
	
	// Join with the base directory from the manager
	return filepath.Join(p.manager.GetPcapDirectory(), pcapName)
}

// CleanExpired removes expired entries from the retention map
func (p *PcapHelper) CleanExpired() {
	p.retentionMutex.Lock()
	defer p.retentionMutex.Unlock()
	
	now := time.Now()
	for pcapName, expiration := range p.retentionMap {
		if now.After(expiration) {
			delete(p.retentionMap, pcapName)
		}
	}
}
import (
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
	"github.com/rs/zerolog/log"
)

// PcapHelper provides PCAP-related helper functions for scripts
type PcapHelper struct {
	pcapManager *worker.PcapManager
}

// NewPcapHelper creates a new PCAP helper
func NewPcapHelper(pcapManager *worker.PcapManager) *PcapHelper {
	return &PcapHelper{
		pcapManager: pcapManager,
	}
}

// RetainPcap marks a PCAP file for extended retention
// This can be used by scripts to prevent important PCAPs from being deleted
func (ph *PcapHelper) RetainPcap(pcapName string, durationSeconds int) {
	duration := time.Duration(durationSeconds) * time.Second
	ph.pcapManager.RetainPcap(pcapName, duration)
	log.Info().
		Str("pcapName", pcapName).
		Int("durationSeconds", durationSeconds).
		Msg("Script requested PCAP retention")
}
package scripting
package scripting

import (
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
	"github.com/robertkrimen/otto"
)

// PcapHelper helps with PCAP operations in scripts
type PcapHelper struct {
	pcapManager *worker.PcapManager
}

// NewPcapHelper creates a new PcapHelper
func NewPcapHelper(pcapManager *worker.PcapManager) *PcapHelper {
	return &PcapHelper{
		pcapManager: pcapManager,
	}
}
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

// CleanupExpiredRetentions removes all expired PCAP retentions
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

// GetRetainedPcaps returns a list of all currently retained PCAPs
func (p *PcapHelper) GetRetainedPcaps() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	now := time.Now()
	pcaps := make([]string, 0, len(p.retentions))
	
	for pcapName, expiration := range p.retentions {
		if now.Before(expiration) {
			pcaps = append(pcaps, pcapName)
		}
	}
	
	return pcaps
}

// CleanupExpiredRetentions removes all expired PCAP retentions
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

// GetRetainedPcaps returns a list of all currently retained PCAPs
func (p *PcapHelper) GetRetainedPcaps() []string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	now := time.Now()
	pcaps := make([]string, 0, len(p.retentions))
	
	for pcapName, expiration := range p.retentions {
		if now.Before(expiration) {
			pcaps = append(pcaps, pcapName)
		}
	}
	
	return pcaps
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
// RegisterPcapBindings registers PCAP-related functions to the JS VM
func RegisterPcapBindings(vm *otto.Otto, helper *PcapHelper) {
	pcapObj, _ := vm.Object("pcap = {}")
	
	_ = pcapObj.Set("retain", func(call otto.FunctionCall) otto.Value {
		pcapName := call.Argument(0).String()
		durationSec, _ := call.Argument(1).ToInteger()
		
		helper.pcapManager.RetainPcap(pcapName, time.Duration(durationSec)*time.Second)
		
		result, _ := otto.ToValue(true)
		return result
	})
}
import (
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
	"github.com/rs/zerolog/log"
)

// PcapHelper provides PCAP-related helper functions for scripts
type PcapHelper struct {
	pcapManager *worker.PcapManager
}

// NewPcapHelper creates a new PCAP helper
func NewPcapHelper(pcapManager *worker.PcapManager) *PcapHelper {
	return &PcapHelper{
		pcapManager: pcapManager,
	}
}

// RetainPcap marks a PCAP file for extended retention
// This can be used by scripts to prevent important PCAPs from being deleted
func (ph *PcapHelper) RetainPcap(pcapName string, durationSeconds int) {
	duration := time.Duration(durationSeconds) * time.Second
	ph.pcapManager.RetainPcap(pcapName, duration)
	log.Info().
		Str("pcapName", pcapName).
		Int("durationSeconds", durationSeconds).
		Msg("Script requested PCAP retention")
}
package scripting

import (
	"path/filepath"
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
)

// PcapHelper provides functions for working with PCAP files in scripts
type PcapHelper struct {
	pcapManager *worker.PcapManager
}

// NewPcapHelper creates a new PcapHelper instance
func NewPcapHelper(pcapManager *worker.PcapManager) *PcapHelper {
	return &PcapHelper{
		pcapManager: pcapManager,
	}
}

// RetainPcap marks a PCAP file for retention for the specified duration in seconds
func (h *PcapHelper) RetainPcap(pcapName string, durationSec int) {
	h.pcapManager.RetainPcap(pcapName, time.Duration(durationSec)*time.Second)
}

// IsRetained checks if a PCAP file is currently marked for retention
func (h *PcapHelper) IsRetained(pcapName string) bool {
	return h.pcapManager.IsRetained(pcapName)
}

// GetPcapPath generates a PCAP path for the given stream ID
func (h *PcapHelper) GetPcapPath(streamID string) string {
	return filepath.Join(h.pcapManager.GetPcapDir(), streamID+".pcap")
}
// GetPcapPath returns the path to a PCAP file
func (ph *PcapHelper) GetPcapPath(streamID string) string {
	// Implementation depends on how PcapManager stores files
	// This is a placeholder for the actual implementation
	return streamID + ".pcap"
}

// IsRetained checks if a PCAP file is marked for extended retention
func (ph *PcapHelper) IsRetained(pcapName string) bool {
	return ph.pcapManager.IsRetained(pcapName)
}
// GetPcapPath returns the path to a PCAP file
func (ph *PcapHelper) GetPcapPath(streamID string) string {
	// Implementation depends on how PcapManager stores files
	// This is a placeholder for the actual implementation
	return streamID + ".pcap"
}

// IsRetained checks if a PCAP file is marked for extended retention
func (ph *PcapHelper) IsRetained(pcapName string) bool {
	return ph.pcapManager.IsRetained(pcapName)
}
import (
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
	"github.com/rs/zerolog/log"
)

// PcapHelper provides PCAP-related helper functions for scripts
type PcapHelper struct {
	pcapManager *worker.PcapManager
}

// NewPcapHelper creates a new PCAP helper
func NewPcapHelper(pcapManager *worker.PcapManager) *PcapHelper {
	return &PcapHelper{
		pcapManager: pcapManager,
	}
}

// RetainPcap marks a PCAP file for extended retention
// This can be used by scripts to prevent important PCAPs from being deleted
func (ph *PcapHelper) RetainPcap(pcapName string, durationSeconds int) {
	duration := time.Duration(durationSeconds) * time.Second
	ph.pcapManager.RetainPcap(pcapName, duration)
	log.Info().
		Str("pcapName", pcapName).
		Int("durationSeconds", durationSeconds).
		Msg("Script requested PCAP retention")
}

// GetPcapPath returns the path to a PCAP file
func (ph *PcapHelper) GetPcapPath(streamID string) string {
	// Implementation depends on how PcapManager stores files
	// This is a placeholder for the actual implementation
	return streamID + ".pcap"
}

// IsRetained checks if a PCAP file is marked for extended retention
func (ph *PcapHelper) IsRetained(pcapName string) bool {
	return ph.pcapManager.IsRetained(pcapName)
}
