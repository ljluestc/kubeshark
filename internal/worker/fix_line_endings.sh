#!/bin/bash

# This script fixes line endings and removes BOM from Go files
# Run from project root: bash internal/worker/fix_line_endings.sh

echo "Fixing line endings in Go files..."

# Find all Go files and convert CRLF to LF
find . -name "*.go" -type f -exec sed -i 's/\r$//' {} \;

# Remove BOM from Go files
find . -name "*.go" -type f -exec sed -i '1s/^\xEF\xBB\xBF//' {} \;

echo "Done. Now run 'go fmt ./...' to ensure proper formatting."
