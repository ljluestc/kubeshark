#!/bin/bash

# This script completely cleans all Go packages and rebuilds them

echo "=== Starting comprehensive cleanup ==="

# Step 1: Fix file encoding and remove duplicate package declarations
echo "Fixing file encoding issues..."
./fix_file_encoding.sh

# Step 2: Fix Go formatting across all packages
echo "Running go fmt on all packages..."
go fmt ./...

# Step 3: Clean Go cache to ensure fresh build
echo "Cleaning Go cache..."
go clean -cache
go clean -modcache

# Step 4: Tidy Go modules
echo "Tidying Go modules..."
go mod tidy

# Step 5: Install the module
echo "Installing the module..."
go install ./...

# Step 6: Remove all existing test caches
echo "Cleaning test cache..."
go clean -testcache

echo "=== Cleanup complete ==="
echo "Now try running 'make test' again"
