#!/bin/bash

# Script to specifically fix config/configStructs package files
# This is a targeted fix for the persistent issues in these files

CONFIG_STRUCTS_DIR="config/configStructs"

echo "=== Fixing config/configStructs package files ==="

# Process each file in the configStructs directory
for file in "$CONFIG_STRUCTS_DIR"/*.go; do
  echo "Processing $file"
  
  # Create a temporary file with correct content
  cat > "${file}.new" << EOF
package configStructs

$(tail -n +2 "$file" | grep -v "^package ")
EOF

  # Replace the original file with the fixed version
  mv "${file}.new" "$file"
  
  echo "Fixed $file"
done

echo "=== Config structs fix complete ==="

# Make sure the permissions are correct
chmod 644 "$CONFIG_STRUCTS_DIR"/*.go

# Run go fmt on the package
go fmt "./$CONFIG_STRUCTS_DIR"

echo "Config structs files have been fixed and formatted."
