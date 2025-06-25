package configStructs

// SecurityContext defines security settings for containers
type SecurityContext struct {
	Privileged      bool                  `yaml:"privileged" json:"privileged" default:"true"`
	AppArmorProfile AppArmorProfileConfig `yaml:"appArmorProfile" json:"appArmorProfile"`
	SeLinuxOptions  SeLinuxOptionsConfig  `yaml:"seLinuxOptions" json:"seLinuxOptions"`
	Capabilities    CapabilitiesConfig    `yaml:"capabilities" json:"capabilities"`
}

// AppArmorProfileConfig defines the configuration for AppArmor profiles
type AppArmorProfileConfig struct {
	Type             string `yaml:"type" json:"type"`
	LocalhostProfile string `yaml:"localhostProfile" json:"localhostProfile"`
}

// SeLinuxOptionsConfig defines the configuration for SELinux options
type SeLinuxOptionsConfig struct {
	Level string `yaml:"level" json:"level"`
	Role  string `yaml:"role" json:"role"`
	Type  string `yaml:"type" json:"type"`
	User  string `yaml:"user" json:"user"`
}

// CapabilitiesConfig defines the configuration for container capabilities
type CapabilitiesConfig struct {
	NetworkCapture     []string `yaml:"networkCapture" json:"networkCapture"  default:"[]"`
	ServiceMeshCapture []string `yaml:"serviceMeshCapture" json:"serviceMeshCapture"  default:"[]"`
	EBPFCapture        []string `yaml:"ebpfCapture" json:"ebpfCapture"  default:"[]"`
}
