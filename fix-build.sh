#!/bin/bash
# This script fixes the build issues in the Kubeshark project

echo "Fixing build issues in Kubeshark..."

# Step 1: Create the SecurityContext type in a separate file
cat > config/configStructs/securityContext.go << 'EOL'
package configStructs

// SecurityContext defines security settings for containers
type SecurityContext struct {
	RunAsUser  int64 `json:"runAsUser" yaml:"runAsUser" default:"0"`
	RunAsGroup int64 `json:"runAsGroup" yaml:"runAsGroup" default:"0"`
}
EOL

echo "Created SecurityContext definition file"

# Step 2: Fix the duplicate securityContext in Helm chart
sed -i '/securityContext:/{ N; N; N; N; N; N; d; }' helm-chart/templates/08-front-deployment.yaml
echo "Fixed duplicate securityContext in Helm chart"

# Step 3: Fix the redeclaration issues in tapConfig.go
# First, backup the file
cp config/configStructs/tapConfig.go config/configStructs/tapConfig.go.bak

# Fix the Hub field issue - replace "Hub" with "HubSecurity" for the duplicate field
sed -i 's/Hub                SecurityContext/HubSecurity       SecurityContext/g' config/configStructs/tapConfig.go

# Fix the SecurityContext initialization
sed -i '/sc := DefaultTapConfig.HubSecurity.(SecurityContext)/,+6d' config/configStructs/tapConfig.go
# Add the proper initialization
sed -i '/defaults.Set(DefaultTapConfig)/a \\n\t// Set default values for SecurityContext\n\tsc := SecurityContext{\n\t\tRunAsUser:  0,\n\t\tRunAsGroup: 0,\n\t}\n\tDefaultTapConfig.HubSecurity = sc\n\tDefaultTapConfig.WorkerSecurity = sc' config/configStructs/tapConfig.go

echo "Fixed Hub redeclaration and SecurityContext initialization in tapConfig.go"

# Step 4: Fix the packet_processor_test.go issues
# First, backup the file
cp cmd/worker/packet_processor_test.go cmd/worker/packet_processor_test.go.bak

# Add MockPacket implementation at the top of the file
sed -i '1s/^/package worker\n\nimport (\n\t"github.com\/google\/gopacket"\n\t"testing"\n)\n\n\/\/ MockPacket implements the gopacket.Packet interface for testing\ntype MockPacket struct {\n\tMockNetworkLayer    gopacket.NetworkLayer\n\tMockTransportLayer  gopacket.TransportLayer\n\tMockApplicationLayer gopacket.ApplicationLayer\n\tMockMetadata        *gopacket.PacketMetadata\n}\n\nfunc (m *MockPacket) String() string { return "MockPacket" }\nfunc (m *MockPacket) Dump() string { return "MockPacket Dump" }\nfunc (m *MockPacket) Layers() []gopacket.Layer { return nil }\nfunc (m *MockPacket) Layer(gopacket.LayerType) gopacket.Layer { return nil }\nfunc (m *MockPacket) LayerClass(gopacket.LayerClass) gopacket.Layer { return nil }\nfunc (m *MockPacket) NetworkLayer() gopacket.NetworkLayer { return m.MockNetworkLayer }\nfunc (m *MockPacket) TransportLayer() gopacket.TransportLayer { return m.MockTransportLayer }\nfunc (m *MockPacket) ApplicationLayer() gopacket.ApplicationLayer { return m.MockApplicationLayer }\nfunc (m *MockPacket) ErrorLayer() gopacket.ErrorLayer { return nil }\nfunc (m *MockPacket) LinkLayer() gopacket.LinkLayer { return nil }\nfunc (m *MockPacket) Data() []byte { return nil }\nfunc (m *MockPacket) Metadata() *gopacket.PacketMetadata { return m.MockMetadata }\n\n/' cmd/worker/packet_processor_test.go

# Rename the second TestGenerateConnectionID
sed -i '0,/TestGenerateConnectionID/!s/TestGenerateConnectionID/TestGenerateConnectionIDWithMock/g' cmd/worker/packet_processor_test.go

echo "Fixed duplicate test function and added MockPacket in packet_processor_test.go"

echo "All fixes applied. Now trying to build..."

# Step 5: Try to build
go mod tidy
go build ./config/configStructs/...
go build ./cmd/worker/...

echo "Build complete. Try running tests with 'make test'"

chmod +x fix-build.sh
