#!/bin/bash

# Script to fix go.mod file

echo "=== Fixing go.mod file ==="

# Check if go.mod exists
if [ -f "go.mod" ]; then
  echo "Found go.mod, updating it..."
  
  # Back up the original file
  cp go.mod go.mod.bak
  
  # Fix the Go version directive
  # Go version should be 1.24, not 1.24.3
  sed -i 's/go 1.24.3/go 1.24/' go.mod
  
  echo "go.mod has been updated."
else
  echo "go.mod not found. No changes made."
fi

echo "=== go.mod fix complete ==="
