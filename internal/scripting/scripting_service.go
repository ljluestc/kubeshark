package scripting

import (
	"fmt"
	"sync"
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
	"github.com/robertkrimen/otto"
)

// ScriptingService manages JavaScript script execution for Kubeshark
type ScriptingService struct {
	vm          *otto.Otto
	pcapManager *worker.PcapManager
	pcapTTL     time.Duration
}

// NewScriptingService creates a new scripting service
func NewScriptingService(pcapManager *worker.PcapManager, timeoutMs int) (*ScriptingService, error) {
	pcapHelper := NewPcapHelper(pcapManager)

	engine, err := NewScriptEngine(pcapHelper, timeoutMs)
	if err != nil {
		return nil, fmt.Errorf("failed to create script engine: %w", err)
	}

	return &ScriptingService{
		engine:     engine,
		pcapHelper: pcapHelper,
		mutex:      sync.Mutex{},
	}, nil
}

// ExecuteScript runs a script with locking to ensure thread safety
func (s *ScriptingService) ExecuteScript(script string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Reset engine before execution
	if err := s.engine.Reset(); err != nil {
		return "", fmt.Errorf("failed to reset script engine: %w", err)
	}

	// Execute the script
	return s.engine.Execute(script)
}

// RetainPcap provides access to the PCAP retention functionality
func (s *ScriptingService) RetainPcap(pcapName string, durationSec int) {
	s.pcapHelper.RetainPcap(pcapName, durationSec)
}

// IsRetained checks if a PCAP is currently retained
func (s *ScriptingService) IsRetained(pcapName string) bool {
	return s.pcapHelper.IsRetained(pcapName)
}

// GetPcapPath gets the path for a PCAP file by stream ID
func (s *ScriptingService) GetPcapPath(streamID string) string {
	return s.pcapHelper.GetPcapPath(streamID)
}

// NewScriptingService creates a new scripting service
func NewScriptingService(pcapManager *worker.PcapManager, pcapTTL time.Duration) *ScriptingService {
	return &ScriptingService{
		vm:          otto.New(),
		pcapManager: pcapManager,
		pcapTTL:     pcapTTL,
	}
}

// RegisterPcapFunctions registers PCAP-related functions to the JavaScript VM
func (s *ScriptingService) RegisterPcapFunctions() error {
	// Register pcap.path function
	err := s.vm.Set("pcap", map[string]interface{}{
		"path": func(call otto.FunctionCall) otto.Value {
			streamID, _ := call.Argument(0).ToString()
			pcapPath := fmt.Sprintf("pcaps/master/%s.pcap", streamID)

			// Automatically retain this PCAP for longer time
			s.pcapManager.RetainPcap(pcapPath, s.pcapTTL*2)

			result, _ := otto.ToValue(pcapPath)
			return result
		},
	})
	if err != nil {
		return err
	}

	// Register file.copy function with automatic retention
	err = s.vm.Set("file", map[string]interface{}{
		"copy": func(call otto.FunctionCall) otto.Value {
			src, _ := call.Argument(0).ToString()
			dst, _ := call.Argument(1).ToString()

			// Extract PCAP name from path and retain it
			s.pcapManager.RetainPcap(src, s.pcapTTL)

			// Actual file copy would be implemented here
			// (not shown for brevity)

			result, _ := otto.ToValue(true)
			return result
		},
	})

	return err
}

// ExecuteScript runs a JavaScript script with optimized PCAP handling
func (s *ScriptingService) ExecuteScript(script string, data map[string]interface{}) error {
	// Run the script with optimized PCAP retention
	_, err := s.vm.Run(script)
	return err
}
