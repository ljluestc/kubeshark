#!/bin/bash
# Script to fix securityContext.go issues in the Kubeshark project

echo "Fixing SecurityContext issues..."

# Remove the problematic file if it exists
if [ -f "config/configStructs/securityContext.go" ]; then
  rm config/configStructs/securityContext.go
  echo "Removed existing securityContext.go file"
fi

# Create a new file with the correct content
cat > config/configStructs/securityContext.go << 'EOL'
package configStructs

// SecurityContext defines security settings for containers
type SecurityContext struct {
	RunAsUser  int64 `json:"runAsUser" yaml:"runAsUser" default:"0"`
	RunAsGroup int64 `json:"runAsGroup" yaml:"runAsGroup" default:"0"`
}
EOL
echo "Created new securityContext.go file"

# Fix references in tapConfig.go
# This assumes the ResourcesConfig struct has exactly the form shown in the error message
sed -i 's/\tHub     SecurityContext            `json:"hub" yaml:"hub"`/\tHubSecurity  SecurityContext            `json:"hubSecurity" yaml:"hubSecurity"`/g' config/configStructs/tapConfig.go
sed -i 's/\tFront   SecurityContext            `json:"front" yaml:"front"`/\tFrontSecurity SecurityContext           `json:"frontSecurity" yaml:"frontSecurity"`/g' config/configStructs/tapConfig.go
echo "Fixed references in tapConfig.go"

echo "Fix complete. Try running 'go build ./config/configStructs/...' to verify."
