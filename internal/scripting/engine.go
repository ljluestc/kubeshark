package scripting

import (
	"fmt"
	"time"

	"github.com/robertkrimen/otto"
)

// ScriptEngine handles JavaScript execution in a VM
type ScriptEngine struct {
	vm         *otto.Otto
	pcapHelper *PcapHelper
	timeout    time.Duration
}

// NewScriptEngine creates a new script engine with the given timeout
func NewScriptEngine(pcapHelper *PcapHelper, timeoutMs int) (*ScriptEngine, error) {
	vm := otto.New()
	
	engine := &ScriptEngine{
		vm:         vm,
		pcapHelper: pcapHelper,
		timeout:    time.Duration(timeoutMs) * time.Millisecond,
	}
	
	// Register bindings
	if err := RegisterBindings(vm, pcapHelper); err != nil {
		return nil, fmt.Errorf("failed to register bindings: %w", err)
	}
	
	return engine, nil
}

// Execute runs a JavaScript script with timeout protection
func (e *ScriptEngine) Execute(script string) (string, error) {
	// Set up timeout protection
	done := make(chan struct{})
	var result otto.Value
	var err error

	go func() {
		result, err = e.vm.Run(script)
		close(done)
	}()

	select {
	case <-done:
		if err != nil {
			return "", err
		}
		if result.IsUndefined() || result.IsNull() {
			return "", nil
		}
		return result.ToString()
	case <-time.After(e.timeout):
		e.vm.Interrupt <- func() {
			panic("Script execution timed out")
		}
		return "", fmt.Errorf("script execution timed out after %v", e.timeout)
	}
}

// Reset clears the VM state
func (e *ScriptEngine) Reset() error {
	// Create new VM
	e.vm = otto.New()
	
	// Re-register bindings
	if err := RegisterBindings(e.vm, e.pcapHelper); err != nil {
		return fmt.Errorf("failed to re-register bindings: %w", err)
	}
	
	return nil
}
import (
	"fmt"
	"time"

	"github.com/kubeshark/kubeshark/internal/worker"
	"github.com/robertkrimen/otto"
	"github.com/rs/zerolog/log"
)

// ScriptingEngine manages the JavaScript engine and script execution
type ScriptingEngine struct {
	vm          *otto.Otto
	pcapHelper  *PcapHelper
	timeout     time.Duration
	isInitialized bool
}

// NewScriptingEngine creates a new scripting engine
func NewScriptingEngine(pcapManager *worker.PcapManager, timeoutMs int) *ScriptingEngine {
	pcapHelper := NewPcapHelper(pcapManager)
	
	return &ScriptingEngine{
		pcapHelper: pcapHelper,
		timeout:    time.Duration(timeoutMs) * time.Millisecond,
	}
}

// Initialize sets up the JavaScript VM with necessary bindings
func (se *ScriptingEngine) Initialize() error {
	if se.isInitialized {
		return nil
	}
	
	// Create a new JavaScript VM
	se.vm = otto.New()
	
	// Register PCAP-related functions
	RegisterPcapBindings(se.vm, se.pcapHelper)
	
	// Register other utility functions
	err := se.registerUtilityFunctions()
	if err != nil {
		return fmt.Errorf("failed to register utility functions: %w", err)
	}
	
	se.isInitialized = true
	return nil
}

// ExecuteScript runs a JavaScript script with timeout protection
func (se *ScriptingEngine) ExecuteScript(script string) (otto.Value, error) {
	if !se.isInitialized {
		if err := se.Initialize(); err != nil {
			return otto.UndefinedValue(), err
		}
	}
	
	// Create a channel to signal script completion
	done := make(chan struct{})
	
	// Default return value and error
	var result otto.Value
	var err error
	
	// Execute the script in a goroutine
	go func() {
		result, err = se.vm.Run(script)
		close(done)
	}()
	
	// Wait for script completion or timeout
	select {
	case <-done:
		return result, err
	case <-time.After(se.timeout):
		// Interrupt the VM to stop execution
		se.vm.Interrupt <- func() {
			panic("Script execution timed out")
		}
		
		// Wait for the script to finish
		<-done
		
		// Return a timeout error
		return otto.UndefinedValue(), fmt.Errorf("script execution timed out after %v", se.timeout)
	}
}

// registerUtilityFunctions adds utility functions to the JavaScript VM
func (se *ScriptingEngine) registerUtilityFunctions() error {
	// Add a sleep function
	err := se.vm.Set("sleep", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 {
			log.Error().Msg("sleep requires one argument: milliseconds")
			return otto.UndefinedValue()
		}
		
		ms, err := call.Argument(0).ToInteger()
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse milliseconds argument")
			return otto.UndefinedValue()
		}
		
		time.Sleep(time.Duration(ms) * time.Millisecond)
		
		return otto.UndefinedValue()
	})
	
	if err != nil {
		return err
	}
	
	// Add a log function
	err = se.vm.Set("log", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 {
			return otto.UndefinedValue()
		}
		
		message, err := call.Argument(0).ToString()
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert log message to string")
			return otto.UndefinedValue()
		}
		
		log.Info().Str("script", "user").Msg(message)
		return otto.UndefinedValue()
	})
	
	return err
}
