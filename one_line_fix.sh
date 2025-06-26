#!/bin/bash
# One-line-at-a-time fix script to avoid escaping issues

# Fix go.mod
sed -i 's/go 1.24.3/go 1.24/' go.mod

# Create dirs
mkdir -p config/configStructs
mkdir -p internal/scripting

# Fix configConfig.go
echo "package configStructs" > config/configStructs/configConfig.go
echo "" >> config/configStructs/configConfig.go
echo "// ConfigConfig defines configuration-related settings" >> config/configStructs/configConfig.go
echo "type ConfigConfig struct {" >> config/configStructs/configConfig.go
echo "	Path string \`yaml:\"path\" json:\"path\"\`" >> config/configStructs/configConfig.go
echo "	AutoSave bool \`yaml:\"autoSave\" json:\"autoSave\" default:\"true\"\`" >> config/configStructs/configConfig.go
echo "	WatchChanges bool \`yaml:\"watchChanges\" json:\"watchChanges\" default:\"true\"\`" >> config/configStructs/configConfig.go
echo "}" >> config/configStructs/configConfig.go

# Fix configStruct.go
echo "package configStructs" > config/configStructs/configStruct.go
echo "" >> config/configStructs/configStruct.go
echo "// ConfigStruct is the main configuration structure" >> config/configStructs/configStruct.go
echo "type ConfigStruct struct {" >> config/configStructs/configStruct.go
echo "	LogLevel  string          \`yaml:\"logLevel\" default:\"info\"\`" >> config/configStructs/configStruct.go
echo "	Config    ConfigConfig    \`yaml:\"config\"\`" >> config/configStructs/configStruct.go
echo "	Tap       TapConfig       \`yaml:\"tap\"\`" >> config/configStructs/configStruct.go
echo "	Logs      LogsConfig      \`yaml:\"logs\"\`" >> config/configStructs/configStruct.go
echo "	Scripting ScriptingConfig \`yaml:\"scripting\"\`" >> config/configStructs/configStruct.go
echo "}" >> config/configStructs/configStruct.go

# Fix logsConfig.go
echo "package configStructs" > config/configStructs/logsConfig.go
echo "" >> config/configStructs/logsConfig.go
echo "// LogsConfig defines logging configuration" >> config/configStructs/logsConfig.go
echo "type LogsConfig struct {" >> config/configStructs/logsConfig.go
echo "	Console bool   \`yaml:\"console\" default:\"true\"\`" >> config/configStructs/logsConfig.go
echo "	File    string \`yaml:\"file\" default:\"\"\`" >> config/configStructs/logsConfig.go
echo "}" >> config/configStructs/logsConfig.go

# Fix scriptingConfig.go
echo "package configStructs" > config/configStructs/scriptingConfig.go
echo "" >> config/configStructs/scriptingConfig.go
echo "// ScriptingConfig defines settings for the scripting engine" >> config/configStructs/scriptingConfig.go
echo "type ScriptingConfig struct {" >> config/configStructs/scriptingConfig.go
echo "	Enabled       bool   \`yaml:\"enabled\" default:\"false\"\`" >> config/configStructs/scriptingConfig.go
echo "	TimeoutMs     int    \`yaml:\"timeoutMs\" default:\"5000\"\`" >> config/configStructs/scriptingConfig.go
echo "	DefaultScript string \`yaml:\"defaultScript\" default:\"\"\`" >> config/configStructs/scriptingConfig.go
echo "}" >> config/configStructs/scriptingConfig.go

# Fix tapConfig.go
echo "package configStructs" > config/configStructs/tapConfig.go
echo "" >> config/configStructs/tapConfig.go
echo "// TapConfig defines the network tap configuration" >> config/configStructs/tapConfig.go
echo "type TapConfig struct {" >> config/configStructs/tapConfig.go
echo "	Debug    bool          \`yaml:\"debug\" default:\"false\"\`" >> config/configStructs/tapConfig.go
echo "	Insecure bool          \`yaml:\"insecure\" default:\"false\"\`" >> config/configStructs/tapConfig.go
echo "	Misc     TapMiscConfig \`yaml:\"misc\"\`" >> config/configStructs/tapConfig.go
echo "}" >> config/configStructs/tapConfig.go
echo "" >> config/configStructs/tapConfig.go
echo "// TapMiscConfig contains miscellaneous tap settings" >> config/configStructs/tapConfig.go
echo "type TapMiscConfig struct {" >> config/configStructs/tapConfig.go
echo "	PcapTTL      string \`yaml:\"pcapTTL\" default:\"5m\"\`" >> config/configStructs/tapConfig.go
echo "	PcapSizeLimit int    \`yaml:\"pcapSizeLimit\" default:\"104857600\"\`" >> config/configStructs/tapConfig.go
echo "}" >> config/configStructs/tapConfig.go

# Fix bindings.go
echo "package scripting" > internal/scripting/bindings.go
echo "" >> internal/scripting/bindings.go
echo "import (" >> internal/scripting/bindings.go
echo "	\"time\"" >> internal/scripting/bindings.go
echo "" >> internal/scripting/bindings.go
echo "	\"github.com/kubeshark/kubeshark/internal/worker\"" >> internal/scripting/bindings.go
echo ")" >> internal/scripting/bindings.go
echo "" >> internal/scripting/bindings.go
echo "// ScriptBindings contains the various objects bound to the scripting environment" >> internal/scripting/bindings.go
echo "type ScriptBindings struct {" >> internal/scripting/bindings.go
echo "	PcapHelper *PcapHelper" >> internal/scripting/bindings.go
echo "	StartTime  time.Time" >> internal/scripting/bindings.go
echo "}" >> internal/scripting/bindings.go
echo "" >> internal/scripting/bindings.go
echo "// NewScriptBindings creates a new set of script bindings" >> internal/scripting/bindings.go
echo "func NewScriptBindings(pcapManager *worker.PcapManager) *ScriptBindings {" >> internal/scripting/bindings.go
echo "	return &ScriptBindings{" >> internal/scripting/bindings.go
echo "		PcapHelper: NewPcapHelper(pcapManager)," >> internal/scripting/bindings.go
echo "		StartTime:  time.Now()," >> internal/scripting/bindings.go
echo "	}" >> internal/scripting/bindings.go
echo "}" >> internal/scripting/bindings.go

# Fix scripting_service.go
echo "package scripting" > internal/scripting/scripting_service.go
echo "" >> internal/scripting/scripting_service.go
echo "import (" >> internal/scripting/scripting_service.go
echo "	\"context\"" >> internal/scripting/scripting_service.go
echo "	\"fmt\"" >> internal/scripting/scripting_service.go
echo "	\"sync\"" >> internal/scripting/scripting_service.go
echo "	\"time\"" >> internal/scripting/scripting_service.go
echo "" >> internal/scripting/scripting_service.go
echo "	\"github.com/kubeshark/kubeshark/internal/worker\"" >> internal/scripting/scripting_service.go
echo "	\"github.com/rs/zerolog/log\"" >> internal/scripting/scripting_service.go
echo ")" >> internal/scripting/scripting_service.go
echo "" >> internal/scripting/scripting_service.go
echo "// ScriptingService manages script execution and lifecycle" >> internal/scripting/scripting_service.go
echo "type ScriptingService struct {" >> internal/scripting/scripting_service.go
echo "	engine   *ScriptEngine" >> internal/scripting/scripting_service.go
echo "	bindings *ScriptBindings" >> internal/scripting/scripting_service.go
echo "	mutex    sync.RWMutex" >> internal/scripting/scripting_service.go
echo "	running  bool" >> internal/scripting/scripting_service.go
echo "}" >> internal/scripting/scripting_service.go
echo "" >> internal/scripting/scripting_service.go
echo "// NewScriptingService creates a new scripting service" >> internal/scripting/scripting_service.go
echo "func NewScriptingService(pcapManager *worker.PcapManager, timeoutMs int) *ScriptingService {" >> internal/scripting/scripting_service.go
echo "	bindings := NewScriptBindings(pcapManager)" >> internal/scripting/scripting_service.go
echo "	engine := NewScriptEngine(bindings, timeoutMs)" >> internal/scripting/scripting_service.go
echo "	" >> internal/scripting/scripting_service.go
echo "	return &ScriptingService{" >> internal/scripting/scripting_service.go
echo "		engine:   engine," >> internal/scripting/scripting_service.go
echo "		bindings: bindings," >> internal/scripting/scripting_service.go
echo "		mutex:    sync.RWMutex{}," >> internal/scripting/scripting_service.go
echo "		running:  false," >> internal/scripting/scripting_service.go
echo "	}" >> internal/scripting/scripting_service.go
echo "}" >> internal/scripting/scripting_service.go

# Create basic engine.go
echo "package scripting" > internal/scripting/engine.go
echo "" >> internal/scripting/engine.go
echo "import (" >> internal/scripting/engine.go
echo "	\"context\"" >> internal/scripting/engine.go
echo "	\"fmt\"" >> internal/scripting/engine.go
echo "	\"time\"" >> internal/scripting/engine.go
echo "" >> internal/scripting/engine.go
echo "	\"github.com/rs/zerolog/log\"" >> internal/scripting/engine.go
echo ")" >> internal/scripting/engine.go
echo "" >> internal/scripting/engine.go
echo "// ScriptEngine provides functionality to execute scripts" >> internal/scripting/engine.go
echo "type ScriptEngine struct {" >> internal/scripting/engine.go
echo "	bindings *ScriptBindings" >> internal/scripting/engine.go
echo "	timeout  time.Duration" >> internal/scripting/engine.go
echo "}" >> internal/scripting/engine.go
echo "" >> internal/scripting/engine.go
echo "// NewScriptEngine creates a new scripting engine" >> internal/scripting/engine.go
echo "func NewScriptEngine(bindings *ScriptBindings, timeoutMs int) *ScriptEngine {" >> internal/scripting/engine.go
echo "	return &ScriptEngine{" >> internal/scripting/engine.go
echo "		bindings: bindings," >> internal/scripting/engine.go
echo "		timeout:  time.Duration(timeoutMs) * time.Millisecond," >> internal/scripting/engine.go
echo "	}" >> internal/scripting/engine.go
echo "}" >> internal/scripting/engine.go
echo "" >> internal/scripting/engine.go
echo "// ExecuteScript runs a script with the configured bindings and timeout" >> internal/scripting/engine.go
echo "func (e *ScriptEngine) ExecuteScript(script string) error {" >> internal/scripting/engine.go
echo "	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)" >> internal/scripting/engine.go
echo "	defer cancel()" >> internal/scripting/engine.go
echo "" >> internal/scripting/engine.go
echo "	log.Debug().Msg(\"Executing script...\")" >> internal/scripting/engine.go
echo "" >> internal/scripting/engine.go
echo "	// Placeholder for actual script execution" >> internal/scripting/engine.go
echo "	select {" >> internal/scripting/engine.go
echo "	case <-time.After(50 * time.Millisecond):" >> internal/scripting/engine.go
echo "		log.Debug().Msg(\"Script executed successfully\")" >> internal/scripting/engine.go
echo "		return nil" >> internal/scripting/engine.go
echo "	case <-ctx.Done():" >> internal/scripting/engine.go
echo "		return fmt.Errorf(\"script execution timed out after %v\", e.timeout)" >> internal/scripting/engine.go
echo "	}" >> internal/scripting/engine.go
echo "}" >> internal/scripting/engine.go

# Create pcap_helpers.go
echo "package scripting" > internal/scripting/pcap_helpers.go
echo "" >> internal/scripting/pcap_helpers.go
echo "import (" >> internal/scripting/pcap_helpers.go
echo "	\"path/filepath\"" >> internal/scripting/pcap_helpers.go
echo "	\"sync\"" >> internal/scripting/pcap_helpers.go
echo "	\"time\"" >> internal/scripting/pcap_helpers.go
echo "" >> internal/scripting/pcap_helpers.go
echo "	\"github.com/kubeshark/kubeshark/internal/worker\"" >> internal/scripting/pcap_helpers.go
echo ")" >> internal/scripting/pcap_helpers.go
echo "" >> internal/scripting/pcap_helpers.go
echo "// PcapHelper provides functionality for managing PCAP files" >> internal/scripting/pcap_helpers.go
echo "type PcapHelper struct {" >> internal/scripting/pcap_helpers.go
echo "	manager    *worker.PcapManager" >> internal/scripting/pcap_helpers.go
echo "	retentions map[string]time.Time" >> internal/scripting/pcap_helpers.go
echo "	mutex      sync.RWMutex" >> internal/scripting/pcap_helpers.go
echo "}" >> internal/scripting/pcap_helpers.go
echo "" >> internal/scripting/pcap_helpers.go
echo "// NewPcapHelper creates a new PCAP helper" >> internal/scripting/pcap_helpers.go
echo "func NewPcapHelper(manager *worker.PcapManager) *PcapHelper {" >> internal/scripting/pcap_helpers.go
echo "	return &PcapHelper{" >> internal/scripting/pcap_helpers.go
echo "		manager:    manager," >> internal/scripting/pcap_helpers.go
echo "		retentions: make(map[string]time.Time)," >> internal/scripting/pcap_helpers.go
echo "		mutex:      sync.RWMutex{}," >> internal/scripting/pcap_helpers.go
echo "	}" >> internal/scripting/pcap_helpers.go
echo "}" >> internal/scripting/pcap_helpers.go
echo "" >> internal/scripting/pcap_helpers.go
echo "// RetainPcap marks a PCAP for retention for the specified duration in seconds" >> internal/scripting/pcap_helpers.go
echo "func (p *PcapHelper) RetainPcap(pcapName string, durationSec int) {" >> internal/scripting/pcap_helpers.go
echo "	p.mutex.Lock()" >> internal/scripting/pcap_helpers.go
echo "	defer p.mutex.Unlock()" >> internal/scripting/pcap_helpers.go
echo "	" >> internal/scripting/pcap_helpers.go
echo "	// Calculate the expiration time" >> internal/scripting/pcap_helpers.go
echo "	expiration := time.Now().Add(time.Duration(durationSec) * time.Second)" >> internal/scripting/pcap_helpers.go
echo "	" >> internal/scripting/pcap_helpers.go
echo "	// Store in our retention map" >> internal/scripting/pcap_helpers.go
echo "	p.retentions[pcapName] = expiration" >> internal/scripting/pcap_helpers.go
echo "}" >> internal/scripting/pcap_helpers.go
echo "" >> internal/scripting/pcap_helpers.go
echo "// IsRetained checks if a PCAP is currently retained" >> internal/scripting/pcap_helpers.go
echo "func (p *PcapHelper) IsRetained(pcapName string) bool {" >> internal/scripting/pcap_helpers.go
echo "	p.mutex.RLock()" >> internal/scripting/pcap_helpers.go
echo "	defer p.mutex.RUnlock()" >> internal/scripting/pcap_helpers.go
echo "	" >> internal/scripting/pcap_helpers.go
echo "	expiration, exists := p.retentions[pcapName]" >> internal/scripting/pcap_helpers.go
echo "	if !exists {" >> internal/scripting/pcap_helpers.go
echo "		return false" >> internal/scripting/pcap_helpers.go
echo "	}" >> internal/scripting/pcap_helpers.go
echo "	" >> internal/scripting/pcap_helpers.go
echo "	// Check if the retention has expired" >> internal/scripting/pcap_helpers.go
echo "	return time.Now().Before(expiration)" >> internal/scripting/pcap_helpers.go
echo "}" >> internal/scripting/pcap_helpers.go
echo "" >> internal/scripting/pcap_helpers.go
echo "// GetPcapPath returns the full path to a PCAP file based on stream ID" >> internal/scripting/pcap_helpers.go
echo "func (p *PcapHelper) GetPcapPath(streamID string) string {" >> internal/scripting/pcap_helpers.go
echo "	filename := streamID + \".pcap\"" >> internal/scripting/pcap_helpers.go
echo "	return filepath.Join(p.manager.GetPcapDir(), filename)" >> internal/scripting/pcap_helpers.go
echo "}" >> internal/scripting/pcap_helpers.go

# Create a simple pcap_helpers_test.go
echo "package scripting" > internal/scripting/pcap_helpers_test.go
echo "" >> internal/scripting/pcap_helpers_test.go
echo "import (" >> internal/scripting/pcap_helpers_test.go
echo "	\"testing\"" >> internal/scripting/pcap_helpers_test.go
echo "	\"time\"" >> internal/scripting/pcap_helpers_test.go
echo "" >> internal/scripting/pcap_helpers_test.go
echo "	\"github.com/kubeshark/kubeshark/internal/worker\"" >> internal/scripting/pcap_helpers_test.go
echo ")" >> internal/scripting/pcap_helpers_test.go
echo "" >> internal/scripting/pcap_helpers_test.go
echo "func TestPcapHelperRetention(t *testing.T) {" >> internal/scripting/pcap_helpers_test.go
echo "	// Create a PCAP manager for testing" >> internal/scripting/pcap_helpers_test.go
echo "	manager := worker.NewPcapManager(\"/tmp/pcaps\", 10*time.Second, 1024*1024)" >> internal/scripting/pcap_helpers_test.go
echo "" >> internal/scripting/pcap_helpers_test.go
echo "	// Create a PCAP helper" >> internal/scripting/pcap_helpers_test.go
echo "	helper := NewPcapHelper(manager)" >> internal/scripting/pcap_helpers_test.go
echo "" >> internal/scripting/pcap_helpers_test.go
echo "	// Test PCAP retention" >> internal/scripting/pcap_helpers_test.go
echo "	testPcapName := \"test_stream_123456.pcap\"" >> internal/scripting/pcap_helpers_test.go
echo "	helper.RetainPcap(testPcapName, 60) // Retain for 60 seconds" >> internal/scripting/pcap_helpers_test.go
echo "	" >> internal/scripting/pcap_helpers_test.go
echo "	// Verify retention" >> internal/scripting/pcap_helpers_test.go
echo "	if !helper.IsRetained(testPcapName) {" >> internal/scripting/pcap_helpers_test.go
echo "		t.Errorf(\"PCAP should be retained\")" >> internal/scripting/pcap_helpers_test.go
echo "	}" >> internal/scripting/pcap_helpers_test.go
echo "}" >> internal/scripting/pcap_helpers_test.go

# Clean caches
go clean -cache -modcache -testcache

echo "One-line fix completed! Try running 'make test' now."
