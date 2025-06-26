package configStructs

// ConfigStruct is the root configuration structure
type ConfigStruct struct {
	// Config contains configuration for the config system
	Config ConfigConfig `yaml:"config" json:"config"`
	// Logs contains configuration for logging
	Logs LogsConfig `yaml:"logs" json:"logs"`
	// Scripting contains configuration for scripting
	Scripting ScriptingConfig `yaml:"scripting" json:"scripting"`
	// Tap contains configuration for tap system
	Tap TapConfig `yaml:"tap" json:"tap"`
}
