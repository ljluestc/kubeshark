#!/bin/bash

# Script to fix conflicting test functions in the worker package

WORKER_DIR="internal/worker"

echo "=== Fixing worker package test files ==="

# Check if the files exist
if [ -f "$WORKER_DIR/pcap_manager_test.go" ] && [ -f "$WORKER_DIR/pcap_retention_test.go" ]; then
  echo "Both test files exist, handling the conflict..."
  
  # Create a backup
  cp "$WORKER_DIR/pcap_retention_test.go" "$WORKER_DIR/pcap_retention_test.go.bak"
  
  # Replace the conflicting function names in pcap_retention_test.go
  sed -i 's/TestPcapRetention/TestPcapRetentionAlternative/g' "$WORKER_DIR/pcap_retention_test.go"
  sed -i 's/TestPcapRetentionExpiration/TestPcapRetentionExpirationAlternative/g' "$WORKER_DIR/pcap_retention_test.go"
  
  # Ensure package declaration is correct
  # This creates a temporary file with a correct package declaration at the top
  echo "package worker" > "$WORKER_DIR/pcap_retention_test.go.new"
  grep -v "^package " "$WORKER_DIR/pcap_retention_test.go" >> "$WORKER_DIR/pcap_retention_test.go.new"
  mv "$WORKER_DIR/pcap_retention_test.go.new" "$WORKER_DIR/pcap_retention_test.go"
  
  # Format the file
  go fmt "./$WORKER_DIR/pcap_retention_test.go"
  
  echo "Worker test files fixed."
else
  echo "One or both test files are missing. No changes made."
fi

echo "=== Worker package test fix complete ==="
