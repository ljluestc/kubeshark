#!/bin/bash

# CHECK AND FIX LINE ENDINGS
# This script checks for and fixes potential line ending issues

echo "=== CHECKING LINE ENDINGS ==="

# List of critical files to check
FILES=(
  "config/configStructs/configConfig.go"
  "config/configStructs/configStruct.go"
  "config/configStructs/logsConfig.go"
  "config/configStructs/scriptingConfig.go"
  "config/configStructs/tapConfig.go"
  "internal/scripting/bindings.go"
  "internal/scripting/engine.go"
  "internal/scripting/pcap_helpers.go"
  "internal/scripting/scripting_service.go"
)

# Function to check and fix a file
check_and_fix() {
  local file=$1
  
  if [ ! -f "$file" ]; then
    echo "  $file does not exist, skipping"
    return
  fi
  
  echo "Checking $file..."
  
  # Check for weird characters or line endings
  if file "$file" | grep -q "CRLF"; then
    echo "  Found CRLF line endings in $file, converting to LF"
    tr -d '\r' < "$file" > "$file.tmp"
    mv "$file.tmp" "$file"
  fi
  
  # Check for BOM
  if hexdump -C "$file" | head -1 | grep -q "ef bb bf"; then
    echo "  Found BOM in $file, removing"
    tail -c +4 "$file" > "$file.tmp"
    mv "$file.tmp" "$file"
  fi
  
  # Check for null bytes or other non-printable characters
  if grep -q '[[:cntrl:]]' "$file"; then
    echo "  Found control characters in $file, cleaning"
    tr -cd '[:print:]\n' < "$file" > "$file.tmp"
    mv "$file.tmp" "$file"
  fi
  
  echo "  Ensuring file starts with 'package'"
  # Extract first non-empty line and check if it's a package declaration
  first_line=$(grep -v '^$' "$file" | head -1)
  if [[ ! "$first_line" =~ ^package ]]; then
    echo "  First line is not a package declaration, fixing"
    # Get package name from filename or directory
    pkg_name=$(basename $(dirname "$file"))
    # Create a new file with proper package declaration
    echo "package $pkg_name" > "$file.tmp"
    echo "" >> "$file.tmp"
    grep -v '^package' "$file" >> "$file.tmp"
    mv "$file.tmp" "$file"
  fi
  
  echo "  $file processed"
}

# Check go.mod first
echo "Fixing go.mod..."
sed -i 's/go 1.24.3/go 1.24/' go.mod

# Process each file
for file in "${FILES[@]}"; do
  check_and_fix "$file"
done

# Ensure directories exist
mkdir -p config/configStructs
mkdir -p internal/scripting

# Format all Go files
echo "Formatting Go files..."
go fmt ./...

echo "=== LINE ENDINGS CHECK COMPLETE ==="
echo "Run: chmod +x verify_line_endings.sh && ./verify_line_endings.sh"
echo "Then try: make test"
