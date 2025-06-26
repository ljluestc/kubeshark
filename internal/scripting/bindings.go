package scripting
package scripting

import (
	"fmt"
	"time"

	"github.com/robertkrimen/otto"
)

// RegisterBindings registers all JavaScript bindings for the scripting engine
func RegisterBindings(vm *otto.Otto, pcapHelper *PcapHelper) error {
	// Register PCAP helper functions
	if err := registerPcapHelperBindings(vm, pcapHelper); err != nil {
		return fmt.Errorf("failed to register PCAP helper bindings: %w", err)
	}

	// Register console functions
	if err := registerConsoleBindings(vm); err != nil {
		return fmt.Errorf("failed to register console bindings: %w", err)
	}

	// Register utility functions
	if err := registerUtilBindings(vm); err != nil {
		return fmt.Errorf("failed to register utility bindings: %w", err)
	}

	return nil
}

// registerPcapHelperBindings registers PCAP helper functions to the JavaScript VM
func registerPcapHelperBindings(vm *otto.Otto, pcapHelper *PcapHelper) error {
	// Create pcap object
	pcapObj, err := vm.Object("pcap = {}")
	if err != nil {
		return err
	}

	// Register retain function
	err = pcapObj.Set("retain", func(call otto.FunctionCall) otto.Value {
		pcapName := call.Argument(0).String()
		durationSec, _ := call.Argument(1).ToInteger()
		
		pcapHelper.RetainPcap(pcapName, int(durationSec))
		
		return otto.UndefinedValue()
	})
	if err != nil {
		return err
	}

	// Register isRetained function
	err = pcapObj.Set("isRetained", func(call otto.FunctionCall) otto.Value {
		pcapName := call.Argument(0).String()
		isRetained := pcapHelper.IsRetained(pcapName)
		
		result, _ := vm.ToValue(isRetained)
		return result
	})
	if err != nil {
		return err
	}

	// Register getPcapPath function
	err = pcapObj.Set("getPcapPath", func(call otto.FunctionCall) otto.Value {
		streamID := call.Argument(0).String()
		path := pcapHelper.GetPcapPath(streamID)
		
		result, _ := vm.ToValue(path)
		return result
	})
	if err != nil {
		return err
	}

	return nil
}

// registerConsoleBindings registers console.log and similar functions
func registerConsoleBindings(vm *otto.Otto) error {
	console, err := vm.Object("console = {}")
	if err != nil {
		return err
	}

	err = console.Set("log", func(call otto.FunctionCall) otto.Value {
		for _, arg := range call.ArgumentList {
			fmt.Print(arg.String() + " ")
		}
		fmt.Println()
		return otto.UndefinedValue()
	})
	if err != nil {
		return err
	}

	return nil
}

// registerUtilBindings registers utility functions like sleep
func registerUtilBindings(vm *otto.Otto) error {
	// Register global sleep function
	err := vm.Set("sleep", func(call otto.FunctionCall) otto.Value {
		milliseconds, _ := call.Argument(0).ToInteger()
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
		return otto.UndefinedValue()
	})
	if err != nil {
		return err
	}

	return nil
}
import (
	"github.com/robertkrimen/otto"
	"github.com/rs/zerolog/log"
)

// RegisterPcapBindings registers PCAP-related functions to the JavaScript VM
func RegisterPcapBindings(vm *otto.Otto, pcapHelper *PcapHelper) {
	pcapObj, _ := vm.Object("({})")
	
	// Register the retain function
	_ = pcapObj.Set("retain", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 {
			log.Error().Msg("pcap.retain requires two arguments: streamID and duration in seconds")
			return otto.UndefinedValue()
		}
		
		streamID, err := call.Argument(0).ToString()
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse streamID argument")
			return otto.UndefinedValue()
		}
		
		duration, err := call.Argument(1).ToInteger()
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse duration argument")
			return otto.UndefinedValue()
		}
		
		pcapHelper.RetainPcap(streamID, int(duration))
		
		result, _ := otto.ToValue(true)
		return result
	})
	
	// Register the isRetained function
	_ = pcapObj.Set("isRetained", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 {
			log.Error().Msg("pcap.isRetained requires one argument: streamID")
			return otto.UndefinedValue()
		}
		
		streamID, err := call.Argument(0).ToString()
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse streamID argument")
			return otto.UndefinedValue()
		}
		
		isRetained := pcapHelper.IsRetained(streamID)
		
		result, _ := otto.ToValue(isRetained)
		return result
	})
	
	// Register the path function
	_ = pcapObj.Set("path", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 {
			log.Error().Msg("pcap.path requires one argument: streamID")
			return otto.UndefinedValue()
		}
		
		streamID, err := call.Argument(0).ToString()
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse streamID argument")
			return otto.UndefinedValue()
		}
		
		path := pcapHelper.GetPcapPath(streamID)
		
		result, _ := otto.ToValue(path)
		return result
	})
	
	// Set the pcap object as a property of the global object
	vm.Set("pcap", pcapObj)
}
import (
	"github.com/robertkrimen/otto"
	"github.com/rs/zerolog/log"
)

// RegisterPcapBindings registers PCAP-related functions to the JavaScript VM
func RegisterPcapBindings(vm *otto.Otto, pcapHelper *PcapHelper) {
	pcapObj, _ := vm.Object("({})")
	
	// Register the retain function
	_ = pcapObj.Set("retain", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 {
			log.Error().Msg("pcap.retain requires two arguments: streamID and duration in seconds")
			return otto.UndefinedValue()
		}
		
		streamID, err := call.Argument(0).ToString()
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse streamID argument")
			return otto.UndefinedValue()
		}
		
		duration, err := call.Argument(1).ToInteger()
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse duration argument")
			return otto.UndefinedValue()
		}
		
		pcapHelper.RetainPcap(streamID, int(duration))
		
		result, _ := otto.ToValue(true)
		return result
	})
	
	// Register the isRetained function
	_ = pcapObj.Set("isRetained", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 {
			log.Error().Msg("pcap.isRetained requires one argument: streamID")
			return otto.UndefinedValue()
		}
		
		streamID, err := call.Argument(0).ToString()
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse streamID argument")
			return otto.UndefinedValue()
		}
		
		isRetained := pcapHelper.IsRetained(streamID)
		
		result, _ := otto.ToValue(isRetained)
		return result
	})
	
	// Register the path function
	_ = pcapObj.Set("path", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 {
			log.Error().Msg("pcap.path requires one argument: streamID")
			return otto.UndefinedValue()
		}
		
		streamID, err := call.Argument(0).ToString()
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse streamID argument")
			return otto.UndefinedValue()
		}
		
		path := pcapHelper.GetPcapPath(streamID)
		
		result, _ := otto.ToValue(path)
		return result
	})
	
	// Set the pcap object as a property of the global object
	vm.Set("pcap", pcapObj)
}
