#!/bin/bash

# Master script to fix all issues in the project
# This script orchestrates all the fixes in the proper sequence

echo "====================================="
echo "= KUBESHARK PROJECT REPAIR UTILITY ="
echo "====================================="

# Step 1: Make all utility scripts executable
echo "Making utility scripts executable..."
chmod +x fix_file_encoding.sh fix_config_structs.sh fix_worker_tests.sh

# Step 2: Fix file encoding issues (general pass)
echo "Running first pass file encoding fixes..."
./fix_file_encoding.sh

# Step 3: Fix config structs specifically
echo "Fixing config/configStructs package..."
./fix_config_structs.sh

# Step 4: Fix worker tests specifically
echo "Fixing worker test conflicts..."
./fix_worker_tests.sh

# Step 5: Fix the Go module
echo "Updating go.mod to ensure correct version syntax..."
# Go 1.24.3 should be 1.24 in go.mod, the patch version is omitted
sed -i 's/go 1.24.3/go 1.24/' go.mod

# Step 6: Clean up Go caches
echo "Cleaning Go caches..."
go clean -cache -modcache -testcache

# Step 7: Tidy modules
echo "Running go mod tidy..."
go mod tidy

# Step 8: Format all Go code
echo "Running go fmt on all packages..."
go fmt ./...

echo ""
echo "====================================="
echo "= REPAIR COMPLETE                  ="
echo "====================================="
echo ""
echo "Now try running 'make test' again."
echo "If issues persist, run 'go test -v ./...' to see detailed error messages."
