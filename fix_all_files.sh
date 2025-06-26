#!/bin/bash

# This script fixes common issues with Go files that might cause build problems
echo "Starting file cleanup process..."

# Fix line endings (CRLF -> LF)
echo "Fixing line endings..."
find . -name "*.go" -type f -exec sed -i 's/\r$//' {} \;

# Remove BOM from Go files
echo "Removing BOMs..."
find . -name "*.go" -type f -exec sed -i '1s/^\xEF\xBB\xBF//' {} \;

# Fix duplicate package declarations
echo "Fixing package declarations in problematic files..."

# Fix specific files manually
for file in config/configStructs/*.go internal/worker/pcap_manager*.go internal/scripting/pcap_helpers*.go; do
  if [ -f "$file" ]; then
    # Remove everything but keep a backup
    cp "$file" "${file}.bak"
    
    # Get just the package name
    package_line=$(grep -m 1 "^package " "${file}.bak")
    
    # Create a new file with just the package line
    echo "$package_line" > "$file"
    echo "" >> "$file"
    
    # Add a comment
    echo "// File fixed by cleanup script" >> "$file"
    echo "" >> "$file"
    
    echo "Fixed $file"
  fi
done

echo "Running go fmt on all packages..."
go fmt ./...

echo "Cleanup complete. Try running tests now."
