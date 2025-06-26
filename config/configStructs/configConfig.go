package configStructs

// ConfigConfig defines configuration-related settings
type ConfigConfig struct {
	// Path to the config file
	Path string `yaml:"path" json:"path"`

	// Whether to save changes automatically
	AutoSave bool `yaml:"autoSave" json:"autoSave" default:"true"`

	// Whether to watch for file changes
	WatchChanges bool `yaml:"watchChanges" json:"watchChanges" default:"true"`

	// Whether to regenerate the config
	Regenerate bool `yaml:"regenerate,omitempty" json:"regenerate,omitempty" default:"false" readonly:""`
}

const (
	RegenerateConfigName = "regenerate"
)
