#!/bin/bash

# Script to fix file encoding issues
# Run this script to clean up BOM (Byte Order Mark) and ensure proper line endings in Go files

echo "=== Starting file encoding cleanup ==="

# Directories to process
DIRS_TO_PROCESS=(
  "internal/worker"
  "internal/scripting"
  "config/configStructs"
)

# Loop through directories
for dir in "${DIRS_TO_PROCESS[@]}"; do
  echo "Processing files in $dir..."
  
  # Find all Go files in the directory
  find "$dir" -name "*.go" -type f | while read -r file; do
    echo "  Checking $file"
    
    # Create a backup
    cp "$file" "$file.bak"
    
    # Remove BOM if present
    # The BOM for UTF-8 is the byte sequence EF BB BF
    if [ "$(hexdump -n 3 -e '3/1 "%02X"' "$file")" = "EFBBBF" ]; then
      echo "    Removing BOM from $file"
      tail -c +4 "$file" > "$file.nobom"
      mv "$file.nobom" "$file"
    fi
    
    # Convert CRLF to LF (Windows to Unix line endings)
    sed -i 's/\r$//' "$file"
    
    # Remove duplicate package declarations
    awk '
      BEGIN { found_pkg = 0; printed = 0; }
      /^package / {
        if (found_pkg == 0) {
          found_pkg = 1;
          printed = 1;
          print;
        } else {
          next;
        }
      }
      {
        if (printed == 0) {
          print;
        } else {
          printed = 0;
        }
      }
    ' "$file" > "$file.tmp"
    mv "$file.tmp" "$file"
  done
done

echo "=== Cleanup complete ==="
echo "Now run 'go fmt ./...' to ensure proper Go formatting"
