#!/bin/bash

# EMERGENCY FIX - Fix only the specific bindings.go file that's causing issues
echo "=== EMERGENCY SINGLE FILE FIX ==="

# Check if the file exists
if [ -f "internal/scripting/bindings.go" ]; then
  echo "Fixing internal/scripting/bindings.go..."
  
  # Fix the file - keep only one package declaration
  grep -v "^package scripting$" internal/scripting/bindings.go | sed '1s/^/package scripting\n\n/' > internal/scripting/bindings.go.new
  mv internal/scripting/bindings.go.new internal/scripting/bindings.go
  
  echo "bindings.go fixed."
else
  echo "Creating internal/scripting/bindings.go..."
  
  # Create the directory if it doesn't exist
  mkdir -p internal/scripting
  
  # Create a minimal version of the file
  cat > internal/scripting/bindings.go << 'EOF'
package scripting

import (
	"github.com/robertkrimen/otto"
)

// RegisterPcapBindings registers PCAP-related functions to the JavaScript VM
func RegisterPcapBindings(vm *otto.Otto, pcapHelper *PcapHelper) {
	// Minimal implementation
}
EOF
  
  echo "bindings.go created."
fi

echo "=== EMERGENCY FIX COMPLETE ==="
echo "Run: chmod +x single_file_fix.sh && ./single_file_fix.sh"
echo "Then try: make test"
